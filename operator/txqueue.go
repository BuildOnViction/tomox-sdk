package operator

import (
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
)

type TxQueue struct {
	Name             string
	Wallet           *types.Wallet
	TradeService     interfaces.TradeService
	OrderService     interfaces.OrderService
	EthereumProvider interfaces.EthereumProvider
	Broker           *rabbitmq.Connection
	AccountService   interfaces.AccountService
	TokenService     interfaces.TokenService
}

func (txq *TxQueue) triggerStopOrders(trades []*types.Trade) {
	for _, trade := range trades {
		stopOrders, err := txq.OrderService.GetTriggeredStopOrders(trade.BaseToken, trade.QuoteToken, trade.PricePoint)

		if err != nil {
			logger.Error(err)
			continue
		}

		for _, stopOrder := range stopOrders {
			err := txq.handleStopOrder(stopOrder)

			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

func (txq *TxQueue) handleStopOrder(so *types.StopOrder) error {
	o, err := so.ToOrder()

	if err != nil {
		logger.Error(err)
		return err
	}

	err = txq.OrderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		return err
	}

	so.Status = types.StopOrderStatusDone
	err = txq.OrderService.UpdateStopOrder(so.Hash, so)

	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
