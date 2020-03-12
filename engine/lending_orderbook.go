package engine

import (
	"sync"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
)

// LendingOrderBook lending order book
type LendingOrderBook struct {
	rabbitMQConn    *rabbitmq.Connection
	orderLendingDao interfaces.LendingOrderDao
	mutex           *sync.Mutex
}

// NewLendingOrderBook new lending order book instance
func NewLendingOrderBook(
	rabbitMQConn *rabbitmq.Connection,
	orderLendingDao interfaces.LendingOrderDao,
) *LendingOrderBook {

	return &LendingOrderBook{
		rabbitMQConn:    rabbitMQConn,
		orderLendingDao: orderLendingDao,
		mutex:           &sync.Mutex{},
	}
}

// newLendingOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (ob *LendingOrderBook) newLendingOrder(o *types.LendingOrder) error {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	err := ob.orderLendingDao.AddNewLendingOrder(o)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// cancelLendingOrder is used to cancel the order from orderbook
func (ob *LendingOrderBook) cancelLendingOrder(o *types.LendingOrder) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()
	err := ob.orderLendingDao.CancelLendingOrder(o)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
