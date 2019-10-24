package engine

import (
	"encoding/json"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/ethereum"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
)

// Engine contains daos required for engine to work
type Engine struct {
	orderbooks   map[string]*OrderBook
	rabbitMQConn *rabbitmq.Connection
	orderDao     interfaces.OrderDao
	tradeDao     interfaces.TradeDao
	pairDao      interfaces.PairDao
	provider     *ethereum.EthereumProvider
}

var logger = utils.Logger

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
		ob := NewOrderBook(rabbitMQConn, orderDao, tradeDao, p)

		obs[p.Code()] = ob
	}

	engine := &Engine{obs, rabbitMQConn, orderDao, tradeDao, pairDao, provider}
	return engine
}

// Provider : implement engine interface
func (e *Engine) Provider() interfaces.EthereumProvider {
	return e.provider
}

func (e *Engine) getObs() (map[string]*OrderBook, error) {
	pairs, err := e.pairDao.GetAll()

	if err != nil {
		return nil, err
	}
	obs := map[string]*OrderBook{}
	for _, p := range pairs {
		ob := NewOrderBook(e.rabbitMQConn, e.orderDao, e.tradeDao, p)

		obs[p.Code()] = ob
	}
	return obs, nil
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
	obs, err := e.getObs()
	if err != nil {
		return errors.New("Orderbook error")
	}
	ob := obs[code]
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

	obs, err := e.getObs()
	if err != nil {
		return errors.New("Orderbook error")
	}
	ob := obs[code]

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
