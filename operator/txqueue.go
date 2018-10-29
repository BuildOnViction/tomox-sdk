package operator

import (
	"encoding/json"
	"errors"
	"math/big"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	eth "github.com/ethereum/go-ethereum/core/types"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/rabbitmq"
	"github.com/tomochain/backend-matching-engine/types"
)

type TxQueue struct {
	Name             string
	Wallet           *types.Wallet
	TradeService     interfaces.TradeService
	OrderService     interfaces.OrderService
	EthereumProvider interfaces.EthereumProvider
	Exchange         interfaces.Exchange
	RabbitMQConn     *rabbitmq.Connection
}

// NewTxQueue
func NewTxQueue(
	n string,
	tr interfaces.TradeService,
	p interfaces.EthereumProvider,
	o interfaces.OrderService,
	w *types.Wallet,
	ex interfaces.Exchange,
	rabbitConn *rabbitmq.Connection,
) (*TxQueue, error) {

	txq := &TxQueue{
		Name:             n,
		TradeService:     tr,
		OrderService:     o,
		EthereumProvider: p,
		Wallet:           w,
		Exchange:         ex,
		RabbitMQConn:     rabbitConn,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return txq, nil
}

func (txq *TxQueue) GetTxSendOptions() *bind.TransactOpts {
	return bind.NewKeyedTransactor(txq.Wallet.PrivateKey)
}

func (txq *TxQueue) GetTxCallOptions() *ethereum.CallMsg {
	// address := txq.Exchange.GetAddress()
	return &ethereum.CallMsg{From: txq.Wallet.Address, To: nil}
}

// Length
func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)
	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	return q.Messages
}

// AddTradeToExecutionList adds a new trade to the execution list. If the execution list is empty (= contains 1 element
// after adding the transaction hash), the given order/trade pair gets executed. If the tranasction queue is full,
// we return an error. Ultimately we want to account send the transaction to another queue that is handled by another ethereum account
// func (op *Operator) QueueTrade(o *types.Order, t *types.Trade) error {
// TODO: Currently doesn't seem thread safe and fails unless called with a sleep time between each call.
func (txq *TxQueue) QueueTrade(o *types.Order, t *types.Trade) error {
	logger.Info("QUEUE LENGTH", txq.Length())
	if txq.Length() == 0 {
		_, err := txq.ExecuteTrade(o, t)
		if err != nil {
			logger.Error(err)
			logger.Info("This is an invalid trade")
			return err
		}
	}

	err := txq.PublishPendingTrade(o, t)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(o *types.Order, tr *types.Trade) (*eth.Transaction, error) {
	logger.Info("EXECUTE_TRADE: ", tr.Hash.Hex())

	// callOpts := txq.GetTxCallOptions()
	// gasLimit, err := txq.Exchange.CallTrade(o, tr, callOpts)
	// if err != nil {
	// 	logger.Error(err)
	// 	return nil, err
	// }
	// should get from config
	var err error
	gasLimit := 0

	if gasLimit < 120000 {
		logger.Warning("GAS LIMIT: ", gasLimit)
		err = txq.RabbitMQConn.PublishTradeInvalidMessage(o, tr)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		go txq.ExecuteNextTrade(tr)
		return nil, errors.New("Invalid Trade")
	}

	nonce, err := txq.EthereumProvider.GetPendingNonceAt(txq.Wallet.Address)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	txOpts := txq.GetTxSendOptions()
	txOpts.Nonce = big.NewInt(int64(nonce))
	// tx, err := txq.Exchange.Trade(o, tr, txOpts)
	// if err != nil {
	// 	logger.Error(err)
	// 	return nil, err
	// }

	var tx common.Address

	err = txq.TradeService.UpdateTradeTxHash(tr, tx.Hash())
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Could not update trade tx attribute")
	}

	err = txq.RabbitMQConn.PublishTradeSentMessage(o, tr)
	if err != nil {
		logger.Error(err)
		return nil, errors.New("Could not update")
	}

	go func() {
		_, err := txq.EthereumProvider.WaitMined(tx.Hash())
		if err != nil {
			logger.Error(err)
		}

		logger.Info("TRADE_MINED IN EXECUTE TRADE: ", tr.Hash.Hex())

		len := txq.Length()
		if len > 0 {
			msg, err := txq.PopPendingTrade()
			if err != nil {
				logger.Error(err)
				return
			}

			// very hacky
			if msg.Trade.Hash == tr.Hash {
				msg, err = txq.PopPendingTrade()
				if err != nil {
					logger.Error(err)
					return
				}

				if msg == nil {
					return
				}
			}

			logger.Info("NEXT_TRADE: ", msg.Trade.Hash.Hex())
			go txq.ExecuteTrade(msg.Order, msg.Trade)
		}
	}()

	return nil, nil
}

func (txq *TxQueue) ExecuteNextTrade(tr *types.Trade) error {
	len := txq.Length()
	logger.Info("LENGTH of the queue is ", len)
	if len > 0 {
		msg, err := txq.PopPendingTrade()
		if err != nil {
			logger.Error(err)
			return err
		}

		logger.Info("NEXT_TRADE: ", msg.Trade.Hash.Hex())
		go txq.ExecuteTrade(msg.Order, msg.Trade)
		return nil
	}

	return nil
}

func (txq *TxQueue) PublishPendingTrade(o *types.Order, t *types.Trade) error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)
	q := txq.RabbitMQConn.GetQueue(ch, name)
	msg := &types.PendingTradeMessage{o, t}

	bytes, err := json.Marshal(msg)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = txq.RabbitMQConn.Publish(ch, q, bytes)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PurgePendingTrades() error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)

	err := txq.RabbitMQConn.Purge(ch, name)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

// PopPendingTrade
func (txq *TxQueue) PopPendingTrade() (*types.PendingTradeMessage, error) {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.RabbitMQConn.GetChannel(name)
	q := txq.RabbitMQConn.GetQueue(ch, name)

	msg, _, _ := ch.Get(
		q.Name,
		true,
	)

	if len(msg.Body) == 0 {
		return nil, nil
	}

	pding := &types.PendingTradeMessage{}
	err := json.Unmarshal(msg.Body, &pding)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return pding, nil
}
