package engine

// The orderbook currently uses the four following data structures to store engine
// state in mongo
// 1. Pricepoints set
// 2. Pricepoints volume set
// 3. Pricepoints hashes set
// 4. Orders map

// 1. The pricepoints set is an ordered set that store all pricepoints.
// Keys: ~ pair addresses + side (BUY or SELL)
// Values: pricepoints set (sorted set but all ranks are actually 0)

// 2. The pricepoints volume set is an order set that store the volume for a given pricepoint
// Keys: pair addresses + side + pricepoint
// Values: volume for corresponding (pair, pricepoint)

// 3. The pricepoints hashes set is an ordered set that stores a set of hashes ranked by creation time for a given pricepoint
// Keys: pair addresses + side + pricepoint
// Values: hashes of orders with corresponding pricepoint

// 4. The orders hashmap is a mapping that stores serialized orders
// Keys: hash
// Values: serialized order

import (
	"sync"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
)

type OrderBook struct {
	rabbitMQConn *rabbitmq.Connection
	orderDao     interfaces.OrderDao
	tradeDao     interfaces.TradeDao
	pair         *types.Pair
	mutex        *sync.Mutex
	topic        string
}

func NewOrderBook(
	rabbitMQConn *rabbitmq.Connection,
	orderDao interfaces.OrderDao,
	tradeDao interfaces.TradeDao,
	p types.Pair,
) *OrderBook {

	return &OrderBook{
		rabbitMQConn: rabbitMQConn,
		orderDao:     orderDao,
		tradeDao:     tradeDao,
		pair:         &p,
		mutex:        &sync.Mutex{},
	}
}

// newOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (ob *OrderBook) newOrder(o *types.Order) error {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	topic := ob.pair.EncodedTopic()

	err := ob.orderDao.AddNewOrder(o, topic)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// CancelOrder is used to cancel the order from orderbook
func (ob *OrderBook) cancelOrder(o *types.Order) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	topic := ob.pair.EncodedTopic()

	err := ob.orderDao.CancelOrder(o, topic)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
