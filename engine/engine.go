package engine

import (
	"encoding/json"
	"fmt"

	"github.com/tomochain/dex-server/errors"
	"github.com/tomochain/dex-server/ethereum"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/rabbitmq"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
)

// Engine contains daos required for engine to work
type Engine struct {
	orderbooks   map[string]*OrderBook
	rabbitMQConn *rabbitmq.Connection
	provider     *ethereum.EthereumProvider
}

var logger = utils.EngineLogger

// NewEngine initializes the engine singleton instance
func NewEngine(
	rabbitMQConn *rabbitmq.Connection,
	orderDao interfaces.OrderDao,
	tradeDao interfaces.TradeDao,
	pairDao interfaces.PairDao,
	provider *ethereum.EthereumProvider,
) *Engine {
	pairs, err := pairDao.GetAll()

	if err != nil {
		panic(err)
	}

	obs := map[string]*OrderBook{}
	for _, p := range pairs {
		ob := NewOrderBook(rabbitMQConn, orderDao, tradeDao, &p)

		obs[p.Code()] = ob
	}

	engine := &Engine{obs, rabbitMQConn, provider}
	return engine
}

// Provider : implement engine interface
func (e *Engine) Provider() interfaces.EthereumProvider {
	return e.provider
}

// HandleOrders parses incoming rabbitmq order messages and redirects them to the appropriate
// engine function
func (e *Engine) HandleOrders(msg *rabbitmq.Message) error {
	switch msg.Type {
	case "NEW_ORDER":
		err := e.handleNewOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "CANCEL_ORDER":
		err := e.handleCancelOrder(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "INVALIDATE_MAKER_ORDERS":
		err := e.handleInvalidateMakerOrders(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	case "INVALIDATE_TAKER_ORDERS":
		err := e.handleInvalidateTakerOrders(msg.Data)
		if err != nil {
			logger.Error(err)
			return err
		}
	default:
		logger.Error("Unknown message", msg)
	}

	return nil
}

func (e *Engine) handleNewOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.newOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleCancelOrder(bytes []byte) error {
	o := &types.Order{}
	err := json.Unmarshal(bytes, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.cancelOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleInvalidateMakerOrders(bytes []byte) error {
	m := types.Matches{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := m.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.invalidateMakerOrders(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) handleInvalidateTakerOrders(bytes []byte) error {
	m := types.Matches{}
	err := json.Unmarshal(bytes, &m)
	if err != nil {
		logger.Error(err)
		return err
	}

	code, err := m.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		logger.Error(err)
		return err
	}

	err = ob.invalidateTakerOrders(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) SyncOrderBook(p types.Pair) error {
	logger.Debugf("*#####%s", p.Code())
	ob := e.orderbooks[p.Code()]

	if ob.topic == "" {
		return errors.New(fmt.Sprintf("Orderbook topic is missing: %s", p.Name()))
	}

	orders, err := ob.orderDao.GetNewOrders(ob.topic)

	if err != nil {
		logger.Error(err)
		return err
	}

	for _, o := range orders {
		res := &types.EngineResponse{}
		if o.Side == "SELL" {
			res, err = ob.sellOrder(o)
			if err != nil {
				logger.Error(err)
				return err
			}

		} else if o.Side == "BUY" {
			res, err = ob.buyOrder(o)
			if err != nil {
				logger.Error(err)
				return err
			}
		}

		// Note: Plug the option for orders like FOC, Limit here (if needed)
		err = ob.rabbitMQConn.PublishEngineResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	//err = ob.orderDao.SyncNewOrders(orders)
	//
	//if err != nil {
	//	logger.Error(err)
	//	return err
	//}

	//for _, o := range orders {
	//	switch o.Status {
	//	case "OPEN":
	//		res := &types.EngineResponse{
	//			Status:  types.ORDER_ADDED,
	//			Order:   o,
	//			Matches: nil,
	//		}
	//
	//		// Note: Plug the option for orders like FOC, Limit here (if needed)
	//		err = e.rabbitMQConn.PublishEngineResponse(res)
	//		if err != nil {
	//			logger.Error(err)
	//			return err
	//		}
	//
	//		return nil
	//
	//	case "CANCELLED":
	//		res := &types.EngineResponse{
	//			Status:  types.ORDER_CANCELLED,
	//			Order:   o,
	//			Matches: nil,
	//		}
	//
	//		err = e.rabbitMQConn.PublishEngineResponse(res)
	//		if err != nil {
	//			logger.Error(err)
	//			return err
	//		}
	//
	//		return nil
	//
	//	default:
	//		res := &types.EngineResponse{
	//			Status:  types.ERROR_STATUS,
	//			Order:   o,
	//			Matches: nil,
	//		}
	//
	//		err = e.rabbitMQConn.PublishEngineResponse(res)
	//		if err != nil {
	//			logger.Error(err)
	//			return err
	//		}
	//
	//		return nil
	//	}
	//}

	return nil
}
