package services

import (
	"context"
	"encoding/json"
	"math/big"
	"strconv"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// LendingOrderService struct
type LendingOrderService struct {
	lendingDao        interfaces.LendingOrderDao
	engine            interfaces.Engine
	broker            *rabbitmq.Connection
	mutext            sync.RWMutex
	bulkLendingOrders map[string]map[common.Hash]*types.LendingOrder
}

// NewLendingOrderService returns a new instance of lending order service
func NewLendingOrderService(lendingDao interfaces.LendingOrderDao, engine interfaces.Engine, broker *rabbitmq.Connection) *LendingOrderService {
	bulkLendingOrders := make(map[string]map[common.Hash]*types.LendingOrder)
	return &LendingOrderService{
		lendingDao,
		engine,
		broker,
		sync.RWMutex{},
		bulkLendingOrders,
	}
}

// GetByHash get lending by hash
func (s *LendingOrderService) GetByHash(hash common.Hash) (*types.LendingOrder, error) {
	return s.lendingDao.GetByHash(hash)
}

// NewLendingOrder validates if the passed order is valid or not based on user's available
// funds and order data.
// If valid: LendingOrder is inserted in DB with order status as new and order is publiched
// on rabbitmq queue for matching engine to process the order
func (s *LendingOrderService) NewLendingOrder(o *types.LendingOrder) error {
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
	err = s.broker.PublishLendingOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelLendingOrder handles the cancellation order requests.
// Only Orders which are OPEN or NEW i.e. Not yet filled/partially filled
// can be cancelled
func (s *LendingOrderService) CancelLendingOrder(o *types.LendingOrder) error {
	return s.lendingDao.CancelLendingOrder(o)
}

// RepayLendingOrder repay
func (s *LendingOrderService) RepayLendingOrder(o *types.LendingOrder) error {
	return s.lendingDao.RepayLendingOrder(o)
}

// TopupLendingOrder topup
func (s *LendingOrderService) TopupLendingOrder(o *types.LendingOrder) error {

	return s.lendingDao.TopupLendingOrder(o)
}

// HandleLendingOrderResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *LendingOrderService) HandleLendingOrderResponse(res *types.EngineResponse) error {
	switch res.Status {
	case types.LENDING_ORDER_ADDED:
		s.handleLendingOrderAdded(res)
		break
	case types.LENDING_ORDER_CANCELLED:
		s.handleLendingOrderCancelled(res)
		break
	case types.LENDING_ORDER_REJECTED:
		s.handleLendingOrderRejected(res)
		break
	case types.LENDING_ORDER_PARTIALLY_FILLED:
		s.handleLendingOrderPartialFilled(res)
		break
	case types.LENDING_ORDER_FILLED:
		s.handleLendingOrderFilled(res)
		break
	case types.ERROR_STATUS:
		s.handleEngineError(res)
		break
	default:
		s.handleEngineUnknownMessage(res)
	}

	if res.Status != types.ERROR_STATUS {
		err := s.saveBulkLendingOrders(res)
		if err != nil {
			logger.Error("Save bulk order", err)
		}
	}

	return nil
}

// handleLendingOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *LendingOrderService) handleLendingOrderAdded(res *types.EngineResponse) {
	o := res.LendingOrder
	ws.SendLendingOrderMessage("LENDINNG_ORDER_ADDED", o.UserAddress, o)
}

func (s *LendingOrderService) handleLendingOrderPartialFilled(res *types.EngineResponse) {
}

func (s *LendingOrderService) handleLendingOrderFilled(res *types.EngineResponse) {
}

func (s *LendingOrderService) handleLendingOrderCancelled(res *types.EngineResponse) {
}

func (s *LendingOrderService) handleLendingOrderRejected(res *types.EngineResponse) {
}

// handleEngineError returns an websocket error message to the client and recovers orders on the
func (s *LendingOrderService) handleEngineError(res *types.EngineResponse) {
}

// handleEngineUnknownMessage returns a websocket messsage in case the engine resonse is not recognized
func (s *LendingOrderService) handleEngineUnknownMessage(res *types.EngineResponse) {
}

// WatchChanges watch database
func (s *LendingOrderService) WatchChanges() {
	ct, sc, err := s.lendingDao.Watch()

	if err != nil {
		logger.Error("Failed to open change stream")
		return //exiting func
	}

	defer ct.Close()
	defer sc.Close()

	// Watch the event again in case there is error and function returned
	defer s.WatchChanges()

	ctx := context.Background()
	go func() {
		for {
			<-time.After(500 * time.Millisecond)
			s.processBulkLendingOrders()
		}
	}()
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
			ev := types.LendingOrderChangeEvent{}

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
				logger.Debugf("Lending Operation Type: %s", ev.OperationType)
				s.HandleDocumentType(ev)
			}
		}
	}
}

// HandleDocumentType handle order frome changing db
func (s *LendingOrderService) HandleDocumentType(ev types.LendingOrderChangeEvent) error {
	res := &types.EngineResponse{}

	if ev.FullDocument.Status == types.OrderStatusOpen {
		res.Status = types.LENDING_ORDER_ADDED
		res.LendingOrder = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusCancelled {
		res.Status = types.LENDING_ORDER_CANCELLED
		res.LendingOrder = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusRejected {
		res.Status = types.LENDING_ORDER_REJECTED
		res.LendingOrder = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusFilled {
		res.Status = types.LENDING_ORDER_FILLED
		res.LendingOrder = ev.FullDocument
	} else if ev.FullDocument.Status == types.OrderStatusPartialFilled {
		res.Status = types.LENDING_ORDER_PARTIALLY_FILLED
		res.LendingOrder = ev.FullDocument
	}

	if res.Status != "" {
		err := s.broker.PublishLendingOrderResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// process for lending form rabbitmq

// HandleLendingOrdersCreateCancel handle lending order api
func (s *LendingOrderService) HandleLendingOrdersCreateCancel(msg *rabbitmq.Message) error {
	switch msg.Type {
	case "NEW_LENDING_ORDER":
		err := s.handleNewLendingOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "CANCEL_LENDING_ORDER":
		err := s.handleCancelLendingOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	default:
		logger.Error("Unknown message", msg)
	}

	return nil
}

func (s *LendingOrderService) handleNewLendingOrder(bytes []byte) error {
	o := &types.LendingOrder{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}
	return s.lendingDao.AddNewLendingOrder(o)
}

func (s *LendingOrderService) handleCancelLendingOrder(bytes []byte) error {
	o := &types.LendingOrder{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}
	return s.lendingDao.CancelLendingOrder(o)
}

// GetLendingNonceByUserAddress return nonce of user order
func (s *LendingOrderService) GetLendingNonceByUserAddress(addr common.Address) (uint64, error) {
	return s.lendingDao.GetLendingNonce(addr)
}

func (s *LendingOrderService) saveBulkLendingOrders(res *types.EngineResponse) error {
	id := utils.GetLendingOrderBookChannelID(res.LendingOrder.Term, res.LendingOrder.LendingToken)

	s.mutext.Lock()
	defer s.mutext.Unlock()

	if _, ok := s.bulkLendingOrders[id]; ok {
		s.bulkLendingOrders[id][res.LendingOrder.Hash] = res.LendingOrder
	} else {
		s.bulkLendingOrders[id] = make(map[common.Hash]*types.LendingOrder)
		s.bulkLendingOrders[id][res.LendingOrder.Hash] = res.LendingOrder
	}
	return nil
}

func (s *LendingOrderService) processBulkLendingOrders() {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	for p, orders := range s.bulkLendingOrders {
		borrow := []map[string]string{}
		lend := []map[string]string{}

		if len(orders) <= 0 {
			continue
		}
		for _, o := range orders {
			side := o.Side
			amount, err := s.lendingDao.GetLendingOrderBookInterest(o.Term, o.LendingToken, o.Interest, side)
			if err != nil {
				logger.Error(err)
			}

			// case where the amount at the pricepoint is equal to 0
			if amount == nil {
				amount = big.NewInt(0)
			}

			update := map[string]string{
				"interest": strconv.FormatUint(o.Interest, 10),
				"amount":   amount.String(),
			}

			if side == types.BORROW {
				borrow = append(borrow, update)
			} else {
				lend = append(lend, update)
			}
		}
		ws.GetLendingOrderBookSocket().BroadcastMessage(p, &types.LendingOrderBook{
			Name:   p,
			Borrow: borrow,
			Lend:   lend,
		})
	}
	s.bulkLendingOrders = make(map[string]map[common.Hash]*types.LendingOrder)
}
