package services

import (
	"context"
	"fmt"
	"log"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	lru "github.com/hashicorp/golang-lru"
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
	stopOrderDao      interfaces.StopOrderDao
	tokenDao          interfaces.TokenDao
	pairDao           interfaces.PairDao
	accountDao        interfaces.AccountDao
	tradeDao          interfaces.TradeDao
	notificationDao   interfaces.NotificationDao
	engine            interfaces.Engine
	validator         interfaces.ValidatorService
	broker            *rabbitmq.Connection
	orderCache        *lru.Cache
	orderbyPricepoint map[string]map[common.Hash]*big.Int
	mutext            sync.RWMutex
	orderPending      []*types.Order
	isFinishCache     bool
}

// NewOrderService returns a new instance of orderservice
func NewOrderService(
	orderDao interfaces.OrderDao,
	stopOrderDao interfaces.StopOrderDao,
	tokenDao interfaces.TokenDao,
	pairDao interfaces.PairDao,
	accountDao interfaces.AccountDao,
	tradeDao interfaces.TradeDao,
	notificationDao interfaces.NotificationDao,
	engine interfaces.Engine,
	validator interfaces.ValidatorService,
	broker *rabbitmq.Connection,
) *OrderService {
	orderCache, _ := lru.New(2000)
	orderbyPricepoint := make(map[string]map[common.Hash]*big.Int)
	return &OrderService{
		orderDao,
		stopOrderDao,
		tokenDao,
		pairDao,
		accountDao,
		tradeDao,
		notificationDao,
		engine,
		validator,
		broker,
		orderCache,
		orderbyPricepoint,
		sync.RWMutex{},
		[]*types.Order{},
		false,
	}
}

func (s *OrderService) getOrderPricepointKey(baseToken, quoteToken common.Address, pricepoint *big.Int, side string) string {
	return fmt.Sprintf("%s::%s::%s::%s", baseToken.Hex(), quoteToken.Hex(), pricepoint.String(), side)
}
func (s *OrderService) updateOrderPricepoint(o *types.Order) {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	if !s.isFinishCache {
		s.orderPending = append(s.orderPending, o)
	} else {
		s.orderPending = append(s.orderPending, o)
		for _, order := range s.orderPending {
			key := s.getOrderPricepointKey(order.BaseToken, order.QuoteToken, order.PricePoint, order.Side)
			remain := big.NewInt(0)
			remain = remain.Sub(order.Amount, order.FilledAmount)
			if o, ok := s.orderbyPricepoint[key]; ok {
				o[order.Hash] = remain
			} else {
				s.orderbyPricepoint[key] = make(map[common.Hash]*big.Int)
				s.orderbyPricepoint[key][order.Hash] = remain
			}
		}
		s.orderPending = s.orderPending[:0]
	}

}

// GetOrderBookPricePoint return amount remain for pricepoint
func (s *OrderService) GetOrderBookPricePoint(baseToken, quoteToken common.Address, pricepoint *big.Int, side string) (*big.Int, error) {
	key := s.getOrderPricepointKey(baseToken, quoteToken, pricepoint, side)
	s.mutext.RLock()
	defer s.mutext.RUnlock()
	if o, ok := s.orderbyPricepoint[key]; ok {
		amount := big.NewInt(0)
		for _, am := range o {
			amount = amount.Add(amount, am)
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
			s.updateOrderPricepoint(order)
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

	err = s.validator.ValidateAvailableBalance(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	s.orderCache.Add(o.Hash, o)

	err = s.broker.PublishNewOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// NewOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: Order is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *OrderService) NewStopOrder(so *types.StopOrder) error {
	if err := so.Validate(); err != nil {
		logger.Error(err)
		return err
	}

	ok, err := so.VerifySignature()
	if err != nil {
		logger.Error(err)
	}

	if !ok {
		return errors.New("Invalid Signature")
	}

	p, err := s.pairDao.GetByTokenAddress(so.BaseToken, so.QuoteToken)
	if err != nil {
		logger.Error(err)
		return err
	}

	if p == nil {
		return errors.New("Pair not found")
	}

	/*
		if math.IsStrictlySmallerThan(so.QuoteAmount(p), p.MinQuoteAmount()) {
			return errors.New("Order amount too low")
		}
	*/

	// Fill token and pair data
	err = so.Process(p)
	if err != nil {
		logger.Error(err)
		return err
	}

	//err = s.validator.ValidateAvailableBalance(so)
	//if err != nil {
	//	logger.Error(err)
	//	return err
	//}

	err = s.broker.PublishNewStopOrderMessage(so)
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
	var o *types.Order
	var err error
	o, err = s.orderDao.GetByHash(oc.OrderHash)
	if err != nil || o == nil {
		order, ok := s.orderCache.Get(oc.OrderHash)
		if !ok {
			return errors.New("No order with corresponding hash")
		} else {
			o = order.(*types.Order)
		}
	}

	if o == nil {
		return errors.New("No order with corresponding hash")
	}

	if o.Status == types.ORDER_FILLED || o.Status == types.ERROR_STATUS || o.Status == types.ORDER_CANCELLED {
		return fmt.Errorf("Cannot cancel order. Status is %v", o.Status)
	}

	o.Nonce = oc.Nonce
	o.Signature = oc.Signature

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

// CancelStopOrder handles the cancellation stop order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *OrderService) CancelStopOrder(oc *types.OrderCancel) error {
	o, err := s.stopOrderDao.GetByHash(oc.OrderHash)
	if err != nil {
		logger.Error(err)
		return err
	}

	if o == nil {
		return errors.New("No stop order with corresponding hash")
	}

	if o.Status == types.ORDER_FILLED || o.Status == types.ERROR_STATUS || o.Status == types.ORDER_CANCELLED {
		return fmt.Errorf("cannot cancel order. Status is %v", o.Status)
	}

	err = s.broker.PublishCancelStopOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
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
	logger.Info("BroadcastOrderBookUpdate add")
	s.broadcastOrderBookUpdate([]*types.Order{o})
	s.broadcastRawOrderBookUpdate([]*types.Order{o})
}

func (s *OrderService) handleOrderPartialFilled(res *types.EngineResponse) {
	logger.Info("BroadcastOrderBookUpdate PartialFilled")
	s.updateOrderPricepoint(res.Order)
	s.broadcastOrderBookUpdate([]*types.Order{res.Order})
	s.broadcastRawOrderBookUpdate([]*types.Order{res.Order})
}

func (s *OrderService) handleOrderFilled(res *types.EngineResponse) {
	logger.Info("BroadcastOrderBookUpdate Filled")
	s.updateOrderPricepoint(res.Order)
	s.broadcastOrderBookUpdate([]*types.Order{res.Order})
	s.broadcastRawOrderBookUpdate([]*types.Order{res.Order})
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
	s.broadcastOrderBookUpdate([]*types.Order{res.Order})
	s.broadcastRawOrderBookUpdate([]*types.Order{res.Order})
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

func (s *OrderService) WatchChanges() {
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

// HandleDocumentType handle order frome changing db
func (s *OrderService) HandleDocumentType(ev types.OrderChangeEvent) error {
	res := &types.EngineResponse{}

	switch ev.OperationType {
	case types.OPERATION_TYPE_INSERT:
		if ev.FullDocument.Status == types.OrderStatusOpen || ev.FullDocument.Status == types.OrderStatusNew {
			res.Status = types.ORDER_ADDED
			res.Order = ev.FullDocument
		}
		break
	case types.OPERATION_TYPE_UPDATE:
	case types.OPERATION_TYPE_REPLACE:
		if ev.FullDocument.Status == types.OrderStatusCancelled {
			res.Status = types.ORDER_CANCELLED
			res.Order = ev.FullDocument
		} else if ev.FullDocument.Status == types.OrderStatusFilled {
			res.Status = types.ORDER_FILLED
			res.Order = ev.FullDocument
		} else if ev.FullDocument.Status == types.OrderStatusPartialFilled {
			res.Status = types.ORDER_PARTIALLY_FILLED
			res.Order = ev.FullDocument
		}
		break
	default:
		break
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

func (s *OrderService) GetTriggeredStopOrders(baseToken, quoteToken common.Address, lastPrice *big.Int) ([]*types.StopOrder, error) {
	return s.stopOrderDao.GetTriggeredStopOrders(baseToken, quoteToken, lastPrice)
}

func (s *OrderService) UpdateStopOrder(h common.Hash, so *types.StopOrder) error {
	return s.stopOrderDao.UpdateByHash(h, so)
}

func (s *OrderService) triggerStopOrders(trades []*types.Trade) {
	for _, trade := range trades {
		stopOrders, err := s.GetTriggeredStopOrders(trade.BaseToken, trade.QuoteToken, trade.PricePoint)

		if err != nil {
			logger.Error(err)
			continue
		}

		for _, stopOrder := range stopOrders {
			err := s.handleStopOrder(stopOrder)

			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

func (s *OrderService) handleStopOrder(so *types.StopOrder) error {
	o, err := so.ToOrder()

	if err != nil {
		logger.Error(err)
		return err
	}

	err = s.NewOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	so.Status = types.StopOrderStatusDone
	err = s.UpdateStopOrder(so.Hash, so)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// GetOrderNonceByUserAddress return nonce of user order
func (s *OrderService) GetOrderNonceByUserAddress(addr common.Address) (interface{}, error) {
	return s.orderDao.GetOrderNonce(addr)
}
