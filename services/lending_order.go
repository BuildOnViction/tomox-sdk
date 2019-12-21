package services

import (
	"context"
	"encoding/json"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
)

// LendingOrderService struct
type LendingOrderService struct {
	lendingDao interfaces.LendingOrderDao
	engine     interfaces.Engine
	broker     *rabbitmq.Connection
	mutext     sync.RWMutex
}

// NewLendingOrderService returns a new instance of lending order service
func NewLendingOrderService(lendingDao interfaces.LendingOrderDao, engine interfaces.Engine, broker *rabbitmq.Connection) *LendingOrderService {
	return &LendingOrderService{
		lendingDao,
		engine,
		broker,
		sync.RWMutex{},
	}
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
func (s *LendingOrderService) CancelLendingOrder(oc *types.LendingOrderCancel) error {
	var err error
	var o *types.LendingOrder
	o.Nonce = oc.Nonce
	o.Signature = oc.Signature
	o.Status = oc.Status
	o.UserAddress = oc.UserAddress
	o.RelayerAddress = oc.RelayerAddress

	err = s.broker.PublishCancelLendingOrderMessage(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// HandleEngineResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *LendingOrderService) HandleEngineResponse(res *types.EngineResponse) error {
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
	}

	return nil
}

// handleLendingOrderAdded returns a websocket message informing the client that his order has been added
// to the orderbook (but currently not matched)
func (s *LendingOrderService) handleLendingOrderAdded(res *types.EngineResponse) {
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
				logger.Debugf("Operation Type: %s", ev.OperationType)
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
		err := s.broker.PublishOrderResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// process for lending form rabbitmq

// HandleLendingOrders handle lending order
func (s *LendingOrderService) HandleLendingOrders(msg *rabbitmq.Message) error {
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
