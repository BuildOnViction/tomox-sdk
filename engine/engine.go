package engine

import (
	"encoding/json"
	"errors"
	"math/big"
	"sync"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/backend-matching-engine/ethereum"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/rabbitmq"
	"github.com/tomochain/backend-matching-engine/redis"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/utils"
)

// Engine contains daos and redis connection required for engine to work
type Engine struct {
	orderbooks   map[string]*OrderBook
	redisConn    *redis.RedisConnection
	rabbitMQConn *rabbitmq.Connection
	provider     *ethereum.EthereumProvider
}

var logger = utils.EngineLogger

// NewEngine initializes the engine singleton instance
func NewEngine(
	redisConn *redis.RedisConnection,
	rabbitMQConn *rabbitmq.Connection,
	pairDao interfaces.PairDao,
	provider *ethereum.EthereumProvider,
) *Engine {

	pairs, err := pairDao.GetAll()
	if err != nil {
		panic(err)
	}

	obs := map[string]*OrderBook{}
	for _, p := range pairs {
		ob := &OrderBook{
			redisConn:    redisConn,
			rabbitMQConn: rabbitMQConn,
			pair:         &p,
			mutex:        &sync.Mutex{},
		}

		obs[p.Code()] = ob
	}

	engine := &Engine{obs, redisConn, rabbitMQConn, provider}
	return engine
}

// Provider : implement engine interface
func (e *Engine) Provider() interfaces.EthereumProvider {
	return e.provider
}

// HandleOrders parses incoming rabbitmq order messages and redirects them to the appropriate
// engine function
func (e *Engine) HandleOrders(msg *rabbitmq.Message) error {
	o := &types.Order{}
	err := json.Unmarshal(msg.Data, o)
	if err != nil {
		logger.Error(err)
		return err
	}

	if msg.Type == "NEW_ORDER" {
		err := e.newOrder(o, msg.HashID)
		if err != nil {
			logger.Error(err)
			return err
		}
	} else if msg.Type == "ADD_ORDER" {
		err := e.addOrder(o)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

func (e *Engine) addOrder(o *types.Order) error {
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.addOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) newOrder(o *types.Order, hashID common.Hash) error {
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.newOrder(o, hashID)
	if err != nil {
		logger.Error(err)
		return err
	}

	// e.provider.RPCClient.Call()

	// also publish this order to PSS :D
	// TODO: call websocket to push the order to swarm feed with signature as well

	return nil
}

func (e *Engine) RecoverOrders(matches []*types.OrderTradePair) error {
	//TODO for now we assume all order/trades have the same token pair
	o := matches[0].Order
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.RecoverOrders(matches)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

//Cancel order is currently not sent through a queue. Not sure i agree with this mechanism
func (e *Engine) CancelOrder(o *types.Order) (*types.EngineResponse, error) {
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return nil, errors.New("Orderbook error")
	}

	res, err := ob.CancelOrder(o)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return res, nil
}

func (e *Engine) DeleteOrders(orders ...types.Order) error {
	//we assume all the orders correspond to the same pair
	code, err := orders[0].PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.deleteOrders(orders...)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) DeleteOrder(o *types.Order) error {
	//we assume all the orders correspond to the same pair
	code, err := o.PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.deleteOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (e *Engine) CancelTrades(orders []*types.Order, amounts []*big.Int) error {
	//we assume all orders are for the same pair
	code, err := orders[0].PairCode()
	if err != nil {
		logger.Error(err)
		return err
	}

	ob := e.orderbooks[code]
	if ob == nil {
		return errors.New("Orderbook error")
	}

	err = ob.CancelTrades(orders, amounts)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
