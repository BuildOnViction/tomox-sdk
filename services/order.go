package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// OrderService
type OrderService struct {
	orderDao          interfaces.OrderDao
	tokenDao          interfaces.TokenDao
	pairDao           interfaces.PairDao
	accountDao        interfaces.AccountDao
	tradeDao          interfaces.TradeDao
	notificationDao   interfaces.NotificationDao
	engine            interfaces.Engine
	validator         interfaces.ValidatorService
	broker            *rabbitmq.Connection
	orderByPricepoint map[string]map[common.Hash]*amountByTime
	mutext            sync.RWMutex
	orderPending      []*types.Order
	isFinishCache     bool
	bulkOrders        map[*types.PairAddresses]map[common.Hash]*types.Order
}

type amountByTime struct {
	filledAmount *big.Int
	amount       *big.Int
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(
	orderDao interfaces.OrderDao,
	tokenDao interfaces.TokenDao,
	pairDao interfaces.PairDao,
	accountDao interfaces.AccountDao,
	tradeDao interfaces.TradeDao,
	notificationDao interfaces.NotificationDao,
	engine interfaces.Engine,
	validator interfaces.ValidatorService,
	broker *rabbitmq.Connection,
) *OrderService {
	bulkOrders := make(map[*types.PairAddresses]map[common.Hash]*types.Order)
	orderByPricepoint := make(map[string]map[common.Hash]*amountByTime)
	return &OrderService{
		orderDao,
		tokenDao,
		pairDao,
		accountDao,
		tradeDao,
		notificationDao,
		engine,
		validator,
		broker,
		orderByPricepoint,
		sync.RWMutex{},
		[]*types.Order{},
		false,
		bulkOrders,
	}
}

func (s *OrderService) getOrderPricepointKey(baseToken, quoteToken common.Address, pricepoint *big.Int, side string) string {
	return fmt.Sprintf("%s::%s::%s::%s", baseToken.Hex(), quoteToken.Hex(), pricepoint.String(), side)
}

func (s *OrderService) update(order *types.Order) {
	key := s.getOrderPricepointKey(order.BaseToken, order.QuoteToken, order.PricePoint, order.Side)
	remain := big.NewInt(0)
	remain = remain.Sub(order.Amount, order.FilledAmount)
	if obyHash, ok := s.orderByPricepoint[key]; ok {
		if amountbytime, ok := obyHash[order.Hash]; ok {
			if order.FilledAmount.Cmp(amountbytime.filledAmount) > 0 {
				amountbytime.amount = remain
			} else {
				logger.Info("update not in order")
			}
		} else {
			obyHash[order.Hash] = &amountByTime{
				filledAmount: order.FilledAmount,
				amount:       remain,
			}
		}
	} else {
		s.orderByPricepoint[key] = make(map[common.Hash]*amountByTime)
		s.orderByPricepoint[key][order.Hash] = &amountByTime{
			filledAmount: order.FilledAmount,
			amount:       remain,
		}
	}
}
func (s *OrderService) updateOrderPricepoint(o *types.Order) {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	if !s.isFinishCache {
		s.orderPending = append(s.orderPending, o)
	} else {
		s.orderPending = append(s.orderPending, o)
		for _, order := range s.orderPending {
			s.update(order)
		}
		s.orderPending = s.orderPending[:0]
	}
}

// GetOrderBookPricePoint return amount remain for pricepoint
func (s *OrderService) GetOrderBookPricePoint(baseToken, quoteToken common.Address, pricepoint *big.Int, side string) (*big.Int, error) {
	key := s.getOrderPricepointKey(baseToken, quoteToken, pricepoint, side)
	s.mutext.RLock()
	defer s.mutext.RUnlock()
	if o, ok := s.orderByPricepoint[key]; ok {
		amount := big.NewInt(0)
		for _, am := range o {
			amount = amount.Add(amount, am.amount)
		}
		return amount, nil
	}
	return nil, errors.New("Cound not found pricepoint key")
}

// LoadCache init order data for caching
func (s *OrderService) LoadCache() {
	logger.Info("Order cache starting ...")
	orders, err := s.orderDao.GetOpenOrders()
	if err == nil {
		for _, order := range orders {
			s.update(order)
		}
	}
	s.isFinishCache = true
	logger.Info("Order cache finish")
}

// GetOrdersLockedBalanceByUserAddress get the total number of orders amount created by a user
func (s *OrderService) GetOrdersLockedBalanceByUserAddress(addr common.Address) (map[string]*big.Int, error) {
	mapAccountBalance := make(map[string]*big.Int)
	pairs, err := s.pairDao.GetActivePairs()
	if err != nil {
		return nil, err
	}
	tokens, err := s.tokenDao.GetAll()
	for _, t := range tokens {
		lockBalance, err := s.orderDao.GetUserLockedBalance(addr, t.ContractAddress, pairs)
		if err != nil {
			return nil, err
		}
		mapAccountBalance[t.Symbol] = lockBalance
	}
	return mapAccountBalance, nil
}

// GetOrderCountByUserAddress get the total number of orders created by a user
func (s *OrderService) GetOrderCountByUserAddress(addr common.Address) (int, error) {
	return s.orderDao.GetOrderCountByUserAddress(addr)
}

// GetByID fetches the details of an order using order's mongo ID
func (s *OrderService) GetByID(id bson.ObjectId) (*types.Order, error) {
	return s.orderDao.GetByID(id)
}

// GetByUserAddress fetches all the orders placed by passed user address
func (s *OrderService) GetByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetByUserAddress(a, bt, qt, from, to, limit...)
}

// GetOrders filter orders
func (s *OrderService) GetOrders(orderSpec types.OrderSpec, sort []string, offset int, size int) (*types.OrderRes, error) {
	return s.orderDao.GetOrders(orderSpec, sort, offset, size)
}

// GetByHash fetches all trades corresponding to a trade hash
func (s *OrderService) GetByHash(hash common.Hash) (*types.Order, error) {
	return s.orderDao.GetByHash(hash)
}

func (s *OrderService) GetByHashes(hashes []common.Hash) ([]*types.Order, error) {
	return s.orderDao.GetByHashes(hashes)
}

// // GetByAddress fetches the detailed document of a token using its contract address
// func (s *OrderService) GetTokenByAddress(addr common.Address) (*types.Token, error) {
// 	return s.tokenDao.GetByAddress(addr)
// }

// GetCurrentByUserAddress function fetches list of open/partial orders from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetCurrentByUserAddress(addr common.Address, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetCurrentByUserAddress(addr, limit...)
}

// GetHistoryByUserAddress function fetches list of orders which are not in open/partial order status
// from order collection based on user address.
// Returns array of Order type struct
func (s *OrderService) GetHistoryByUserAddress(addr, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Order, error) {
	return s.orderDao.GetHistoryByUserAddress(addr, bt, qt, from, to, limit...)
}

// NewOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewOrder(o *types.Order) error {
	if err := o.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	ok, err := o.VerifySignature()
	if err != nil {
		logger.Error(err)
	}

	if !ok {
		return errors.New("Invalid Signature")
	}

	p, err := s.pairDao.GetByTokenAddress(o.BaseToken, o.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	if p == nil {
		return errors.New("Pair not found")
	}

	/*
		if math.IsStrictlySmallerThan(o.QuoteAmount(p), p.MinQuoteAmount()) {
			return errors.New("Order amount too low")
		}
	*/

	// Fill token and pair data
	err = o.Process(p)
	if err != nil {
		logger.Error(err)
		return err
	}
	if o.Type == types.TypeLimitOrder {
		err = s.validator.ValidateAvailableBalance(o)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	err = s.broker.PublishNewOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelOrder(oc *types.OrderCancel) error {
	var err error
	var o *types.Order

	o, err = s.orderDao.GetByHash(oc.OrderHash)
	if err != nil || o == nil {
		return errors.New("No order with corresponding hash")
	}
	if o.Status == types.ORDER_FILLED || o.Status == types.ERROR_STATUS || o.Status == types.ORDER_CANCELLED {
		return fmt.Errorf("Cannot cancel order. Status is %v", o.Status)
	}

	o.Nonce = oc.Nonce
	o.Signature = oc.Signature
	o.OrderID = oc.OrderID
	o.Status = oc.Status
	o.UserAddress = oc.UserAddress
	o.ExchangeAddress = oc.ExchangeAddress

	err = s.broker.PublishCancelOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelAllOrder(a common.Address) error {
	orders, err := s.orderDao.GetOpenOrdersByUserAddress(a)

	if err != nil {
		logger.Error(err)
		return err
	}

	if len(orders) == 0 {
		return nil
	}

	for _, o := range orders {
		err = s.broker.PublishCancelOrderMessage(o)

		if err != nil {
			logger.Error(err)
			continue
		}
	}

	return nil
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *OrderService) HandleEngineResponse(res *types.EngineResponse) error {
	switch res.Status {
	case types.ORDER_ADDED:
		s.handleEngineOrderAdded(res)
		break
	case types.ORDER_CANCELLED:
		s.handleOrderCancelled(res)
		break
	case types.ORDER_REJECTED:
		s.handleOrderRejected(res)
		break
	case types.ORDER_PARTIALLY_FILLED:
		s.handleOrderPartialFilled(res)
		break
	case types.ORDER_FILLED:
		s.handleOrderFilled(res)
		break
	case types.ERROR_STATUS:
		s.handleEngineError(res)
		break
	default:
		s.handleEngineUnknownMessage(res)
	}

	if res.Status != types.ERROR_STATUS {
		err := s.saveBulkOrders(res)
		if err != nil {
			logger.Error("Save bulk order", err)
		}
	}

	return nil
}

func (s *OrderService) saveBulkOrders(res *types.EngineResponse) error {
	p, err := res.Order.Pair()
	if err != nil {
		return err
	}
	pa := types.PairAddresses{
		BaseToken:  p.BaseTokenAddress,
		QuoteToken: p.QuoteTokenAddress,
	}

	s.mutext.Lock()
	defer s.mutext.Unlock()

	if _, ok := s.bulkOrders[&pa]; ok {
		s.bulkOrders[&pa][res.Order.Hash] = res.Order
	} else {
		s.bulkOrders[&pa] = make(map[common.Hash]*types.Order)
		s.bulkOrders[&pa][res.Order.Hash] = res.Order
	}
	return nil
}

// handleEngineOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *OrderService) handleEngineOrderAdded(res *types.EngineResponse) {
	o := res.Order

	// Save notification
	notifications, err := s.notificationDao.Create(&types.Notification{
		Recipient: o.UserAddress,
		Message: types.Message{
			MessageType: "ORDER_ADDED",
			Description: o.Hash.Hex(),
		},
		Type:   types.TypeLog,
		Status: types.StatusUnread,
	})

	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_ADDED", o.UserAddress, o)
	ws.SendNotificationMessage("ORDER_ADDED", o.UserAddress, notifications)
	s.updateOrderPricepoint(o)
}

func (s *OrderService) handleOrderPartialFilled(res *types.EngineResponse) {
	logger.Info("BroadcastOrderBookUpdate PartialFilled")
	s.updateOrderPricepoint(res.Order)
}

func (s *OrderService) handleOrderFilled(res *types.EngineResponse) {
	logger.Info("BroadcastOrderBookUpdate Filled")
	s.updateOrderPricepoint(res.Order)
}

func (s *OrderService) handleOrderCancelled(res *types.EngineResponse) {
	o := res.Order

	// Save notification
	notifications, err := s.notificationDao.Create(&types.Notification{
		Recipient: o.UserAddress,
		Message: types.Message{
			MessageType: "ORDER_CANCELLED",
			Description: o.Hash.Hex(),
		},
		Type:   types.TypeLog,
		Status: types.StatusUnread,
	})

	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_CANCELLED", o.UserAddress, o)
	ws.SendNotificationMessage("ORDER_CANCELLED", o.UserAddress, notifications)
	logger.Info("BroadcastOrderBookUpdate Cancelled")
}

func (s *OrderService) handleOrderRejected(res *types.EngineResponse) {
	o := res.Order

	// Save notification
	notifications, err := s.notificationDao.Create(&types.Notification{
		Recipient: o.UserAddress,
		Message: types.Message{
			MessageType: "ORDER_REJECTED",
			Description: o.Hash.Hex(),
		},
		Type:   types.TypeLog,
		Status: types.StatusUnread,
	})

	if err != nil {
		logger.Error(err)
	}

	ws.SendOrderMessage("ORDER_REJECTED", o.UserAddress, o)
	ws.SendNotificationMessage("ORDER_REJECTED", o.UserAddress, notifications)
	logger.Info("BroadcastOrderBookUpdate rejected")
}

// handleEngineError returns an websocket error message to the client and recovers orders on the
func (s *OrderService) handleEngineError(res *types.EngineResponse) {
	o := res.Order
	ws.SendOrderMessage("ERROR", o.UserAddress, nil)
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine resonse is not recognized
func (s *OrderService) handleEngineUnknownMessage(res *types.EngineResponse) {
	log.Print("Receiving unknown engine message")
	utils.PrintJSON(res)
}

func (s *OrderService) broadcastOrderBookUpdate(orders []*types.Order) {
	bids := []map[string]string{}
	asks := []map[string]string{}

	p, err := orders[0].Pair()
	if err != nil {
		logger.Error()
		return
	}

	for _, o := range orders {
		pp := o.PricePoint
		side := o.Side
		amount, err := s.GetOrderBookPricePoint(o.BaseToken, o.QuoteToken, pp, side)
		if err != nil {
			logger.Error(err)
		}

		// case where the amount at the pricepoint is equal to 0
		if amount == nil {
			amount = big.NewInt(0)
		}

		update := map[string]string{
			"pricepoint": pp.String(),
			"amount":     amount.String(),
		}

		if side == "BUY" {
			bids = append(bids, update)
		} else {
			asks = append(asks, update)
		}
	}

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetOrderBookSocket().BroadcastMessage(id, &types.OrderBook{
		PairName: orders[0].PairName,
		Bids:     bids,
		Asks:     asks,
	})
}

func (s *OrderService) broadcastRawOrderBookUpdate(orders []*types.Order) {
	p, err := orders[0].Pair()
	if err != nil {
		logger.Error(err)
		return
	}

	id := utils.GetOrderBookChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetRawOrderBookSocket().BroadcastMessage(id, orders)
}

// WatchChanges wath change record
func (s *OrderService) WatchChanges() {
	go func() {
		for {
			<-time.After(500 * time.Millisecond)
			s.processBulkOrders()
		}
	}()
	s.watchChanges()
}
func (s *OrderService) watchChanges() {
	ct, sc, err := s.orderDao.Watch()

	if err != nil {
		logger.Error("Failed to open change stream")
		return //exiting func
	}

	defer ct.Close()
	defer sc.Close()

	// Watch the event again in case there is error and function returned
	defer s.WatchChanges()

	ctx := context.Background()

	//Handling change stream in a cycle
	for {
		select {
		case <-ctx.Done(): // if parent context was cancelled
			err := ct.Close() // close the stream
			if err != nil {
				logger.Error("Change stream closed")
			}
			return //exiting from the func
		default:
			ev := types.OrderChangeEvent{}

			//getting next item from the steam
			ok := ct.Next(&ev)

			//if data from the stream wasn't un-marshaled, we get ok == false as a result
			//so we need to call Err() method to get info why
			//it'll be nil if we just have no data
			if !ok {
				err := ct.Err()
				if err != nil {
					logger.Error(err)
					return
				}
			}

			//if item from the stream un-marshaled successfully, do something with it
			if ok {
				logger.Debugf("Operation Type: %s", ev.OperationType)
				s.HandleDocumentType(ev)
			}
		}
	}
}

func (s *OrderService) processBulkOrders() {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	for p, orders := range s.bulkOrders {
		bids := []map[string]string{}
		asks := []map[string]string{}

		if len(orders) <= 0 {
			continue
		}

		var pairName string
		for _, o := range orders {
			pp := o.PricePoint
			side := o.Side
			pairName = o.PairName

			pair, _ := o.Pair()
			amount, err := s.orderDao.GetOrderBookPricePoint(pair, pp, side)
			// amount, err := s.GetOrderBookPricePoint(o.BaseToken, o.QuoteToken, pp, side)
			if err != nil {
				logger.Error(err)
			}

			// case where the amount at the pricepoint is equal to 0
			if amount == nil {
				amount = big.NewInt(0)
			}

			update := map[string]string{
				"pricepoint": pp.String(),
				"amount":     amount.String(),
			}

			if side == "BUY" {
				bids = append(bids, update)
			} else {
				asks = append(asks, update)
			}
		}

		id := utils.GetOrderBookChannelID(p.BaseToken, p.QuoteToken)
		ws.GetOrderBookSocket().BroadcastMessage(id, &types.OrderBook{
			PairName: pairName,
			Bids:     bids,
			Asks:     asks,
		})
	}
	s.bulkOrders = make(map[*types.PairAddresses]map[common.Hash]*types.Order)
}

// HandleDocumentType handle order frome changing db
func (s *OrderService) HandleDocumentType(ev types.OrderChangeEvent) error {
	res := &types.EngineResponse{}

	if ev.FullDocument.Status == types.OrderStatusOpen {
		res.Status = types.ORDER_ADDED
		res.Order = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusCancelled {
		res.Status = types.ORDER_CANCELLED
		res.Order = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusRejected {
		res.Status = types.ORDER_REJECTED
		res.Order = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusFilled {
		res.Status = types.ORDER_FILLED
		res.Order = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusPartialFilled {
		res.Status = types.ORDER_PARTIALLY_FILLED
		res.Order = ev.FullDocument
	}

	if res.Status != "" {
		err := s.broker.PublishOrderResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// GetOrderNonceByUserAddress return nonce of user order
func (s *OrderService) GetOrderNonceByUserAddress(addr common.Address) (interface{}, error) {
	return s.orderDao.GetOrderNonce(addr)
}
