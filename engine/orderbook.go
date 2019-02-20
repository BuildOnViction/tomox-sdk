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

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/rabbitmq"
	"github.com/tomochain/dex-server/types"
)

type OrderBook struct {
	rabbitMQConn *rabbitmq.Connection
	orderDao     interfaces.OrderDao
	tradeDao     interfaces.TradeDao
	pair         *types.Pair
	mutex        *sync.Mutex
}

// newOrder calls buyOrder/sellOrder based on type of order recieved and
// publishes the response back to rabbitmq
func (ob *OrderBook) newOrder(o *types.Order) (err error) {
	// Attain lock on engineResource, so that recovery or cancel order function doesn't interfere
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	err = ob.orderDao.AddNewOrder(o)

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

	err := ob.orderDao.CancelOrder(o)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// cancelTrades revertTrades and reintroduces the taker orders in the orderbook
func (ob *OrderBook) invalidateMakerOrders(matches types.Matches) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	orders := matches.MakerOrders
	trades := matches.Trades
	tradeAmounts := matches.TradeAmounts()
	makerOrderHashes := []common.Hash{}
	takerOrderHashes := []common.Hash{}

	for i, _ := range orders {
		makerOrderHashes = append(makerOrderHashes, trades[i].MakerOrderHash)
		takerOrderHashes = append(takerOrderHashes, trades[i].TakerOrderHash)
	}

	takerOrders, err := ob.orderDao.UpdateOrderFilledAmounts(takerOrderHashes, tradeAmounts)
	if err != nil {
		logger.Error(err)
		return err
	}

	makerOrders, err := ob.orderDao.UpdateOrderStatusesByHashes("INVALIDATED", makerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	//TODO in the case the trades are not in the database they should not be created.
	cancelledTrades, err := ob.tradeDao.UpdateTradeStatusesByOrderHashes("CANCELLED", takerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := &types.EngineResponse{
		Status:            "TRADES_CANCELLED",
		InvalidatedOrders: &makerOrders,
		CancelledTrades:   &cancelledTrades,
	}

	err = ob.rabbitMQConn.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
	}

	for _, o := range takerOrders {
		err := ob.rabbitMQConn.PublishNewOrderMessage(o)
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func (ob *OrderBook) invalidateTakerOrders(matches types.Matches) error {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	makerOrders := matches.MakerOrders
	takerOrder := matches.TakerOrder
	trades := matches.Trades
	tradeAmounts := matches.TradeAmounts()

	makerOrderHashes := []common.Hash{}
	for i, _ := range trades {
		makerOrderHashes = append(makerOrderHashes, trades[i].MakerOrderHash)
	}

	makerOrders, err := ob.orderDao.UpdateOrderFilledAmounts(makerOrderHashes, tradeAmounts)
	if err != nil {
		logger.Error(err)
		return err
	}

	invalidatedOrders, err := ob.orderDao.UpdateOrderStatusesByHashes("INVALIDATED", takerOrder.Hash)
	if err != nil {
		logger.Error(err)
		return err
	}

	cancelledTrades, err := ob.tradeDao.UpdateTradeStatusesByOrderHashes("CANCELLED", makerOrderHashes...)
	if err != nil {
		logger.Error(err)
		return err
	}

	res := &types.EngineResponse{
		Status:            "TRADES_CANCELLED",
		InvalidatedOrders: &invalidatedOrders,
		CancelledTrades:   &cancelledTrades,
	}

	err = ob.rabbitMQConn.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
		return err
	}

	for _, o := range makerOrders {
		err := ob.rabbitMQConn.PublishNewOrderMessage(o)
		if err != nil {
			logger.Error(err)
		}
	}

	return nil
}

func (ob *OrderBook) InvalidateOrder(o *types.Order) (*types.EngineResponse, error) {
	ob.mutex.Lock()
	defer ob.mutex.Unlock()

	o.Status = "ERROR"
	err := ob.orderDao.UpdateOrderStatus(o.Hash, "ERROR")
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	res := &types.EngineResponse{
		Status:  "INVALIDATED",
		Order:   o,
		Matches: nil,
	}

	return res, nil
}
