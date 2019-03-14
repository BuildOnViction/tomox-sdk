package operator

import (
	"encoding/json"
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/dex-server/errors"
	"github.com/tomochain/dex-server/utils/math"
	"math/big"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/streadway/amqp"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/rabbitmq"
	"github.com/tomochain/dex-server/types"
)

type TxQueue struct {
	Name             string
	Wallet           *types.Wallet
	TradeService     interfaces.TradeService
	OrderService     interfaces.OrderService
	EthereumProvider interfaces.EthereumProvider
	Exchange         interfaces.Exchange
	Broker           *rabbitmq.Connection
	AccountService   interfaces.AccountService
}

type TxQueueOrder struct {
	userAddress common.Address
	baseToken   common.Address
	quoteToken  common.Address
	amount      *big.Int
	pricepoint  *big.Int
	side        *big.Int
	salt        *big.Int
	feeMake     *big.Int
	feeTake     *big.Int
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
	accountService interfaces.AccountService,
) (*TxQueue, error) {
	txq := &TxQueue{
		Name:             n,
		TradeService:     tr,
		OrderService:     o,
		EthereumProvider: p,
		Wallet:           w,
		Exchange:         ex,
		Broker:           rabbitConn,
		AccountService:   accountService,
	}

	err := txq.PurgePendingTrades()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)

	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	err = txq.Broker.ConsumeQueuedTrades(ch, &q, txq.ExecuteTrade)
	if err != nil {
		logger.Error(err)
	}

	return txq, nil
}

func (txq *TxQueue) GetChannel() *amqp.Channel {
	name := "TX_QUEUES" + txq.Name
	return txq.Broker.GetChannel(name)
}

func (txq *TxQueue) GetTxSendOptions() *bind.TransactOpts {
	return bind.NewKeyedTransactor(txq.Wallet.PrivateKey)
}

func (txq *TxQueue) GetTxCallOptions() *ethereum.CallMsg {
	address := txq.Exchange.GetAddress()

	return &ethereum.CallMsg{From: txq.Wallet.Address, To: &address}
}

// Length
func (txq *TxQueue) Length() int {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q, err := ch.QueueInspect(name)
	if err != nil {
		logger.Error(err)
	}

	return q.Messages
}

// ExecuteTrade send a trade execution order to the smart contract interface. After sending the
// trade message, the trade is updated on the database and is published to the operator subscribers
// (order service)
func (txq *TxQueue) ExecuteTrade(m *types.Matches, tag uint64) error {
	logger.Infof("Executing trades: %+v", m)

	makerOrders := m.MakerOrders
	trades := m.Trades
	takerOrder := m.TakerOrder

	orderValues := [][10]*big.Int{}
	orderAddresses := [][4]common.Address{}
	vValues := [][2]uint8{}
	rsValues := [][4][32]byte{}
	amounts := []*big.Int{}

	for i := range makerOrders {
		mo := makerOrders[i]
		to := takerOrder
		t := trades[i]

		orderValues = append(orderValues, [10]*big.Int{mo.Amount, mo.PricePoint, mo.EncodedSide(), mo.Nonce, to.Amount, to.PricePoint, to.EncodedSide(), to.Nonce, mo.MakeFee, mo.TakeFee})
		orderAddresses = append(orderAddresses, [4]common.Address{mo.UserAddress, to.UserAddress, mo.BaseToken, to.QuoteToken})
		vValues = append(vValues, [2]uint8{mo.Signature.V, to.Signature.V})
		rsValues = append(rsValues, [4][32]byte{mo.Signature.R, mo.Signature.S, to.Signature.R, to.Signature.S})
		amounts = append(amounts, t.Amount)
	}

	for i := range orderAddresses {
		mOrder := TxQueueOrder{
			userAddress: orderAddresses[i][0],
			baseToken:   orderAddresses[i][2],
			quoteToken:  orderAddresses[i][3],
			amount:      orderValues[i][0],
			pricepoint:  orderValues[i][1],
			side:        orderValues[i][2],
			salt:        orderValues[i][3],
			feeMake:     orderValues[i][8],
			feeTake:     orderValues[i][9],
		}

		tOrder := TxQueueOrder{
			userAddress: orderAddresses[i][1],
			baseToken:   orderAddresses[i][2],
			quoteToken:  orderAddresses[i][3],
			amount:      orderValues[i][4],
			pricepoint:  orderValues[i][5],
			side:        orderValues[i][6],
			salt:        orderValues[i][7],
			feeMake:     orderValues[i][8],
			feeTake:     orderValues[i][9],
		}

		baseTokenAmount := amounts[i]
		quoteTokenAmount := math.Div(math.Div(math.Mul(amounts[i], mOrder.pricepoint), big.NewInt(1e18)), big.NewInt(1e18))

		if math.IsEqual(mOrder.side, big.NewInt(0)) {
			err := txq.AccountService.Transfer(mOrder.quoteToken, mOrder.userAddress, tOrder.userAddress, quoteTokenAmount)
			logger.Error(err)

			err = txq.AccountService.Transfer(tOrder.baseToken, tOrder.userAddress, mOrder.userAddress, baseTokenAmount)
			logger.Error(err)
		} else {
			err := txq.AccountService.Transfer(mOrder.baseToken, mOrder.userAddress, tOrder.userAddress, baseTokenAmount)
			logger.Error(err)

			err = txq.AccountService.Transfer(tOrder.quoteToken, tOrder.userAddress, mOrder.userAddress, quoteTokenAmount)
			logger.Error(err)
		}
	}

	updatedTrades := []*types.Trade{}
	for _, t := range m.Trades {
		updated, err := txq.TradeService.UpdatePendingTrade(t, common.HexToHash("0xf331B044e6E48F4FD154a1B02f3Fb4C344114180"))
		if err != nil {
			logger.Error(err)
		}

		updatedTrades = append(updatedTrades, updated)
	}

	m.Trades = updatedTrades
	err := txq.Broker.PublishTradeSentMessage(m)
	if err != nil {
		logger.Error(err)
		return errors.New("Could not update")
	}

	err = txq.HandleTxSuccess(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) HandleTradeInvalid(m *types.Matches) error {
	logger.Errorf("Trade invalid: %v", m)

	err := txq.Broker.PublishTradeInvalidMessage(m)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) HandleTxError(m *types.Matches) error {
	logger.Errorf("Transaction Error: %v", m)

	errType := "Transaction error"
	err := txq.Broker.PublishTxErrorMessage(m, errType)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) HandleTxSuccess(m *types.Matches) error {
	logger.Infof("Transaction success: %v", m)

	err := txq.Broker.PublishTradeSuccessMessage(m)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) HandleError(m *types.Matches) error {
	logger.Errorf("Operator Error: %v", m)

	errType := "Server error"
	err := txq.Broker.PublishErrorMessage(m, errType)
	if err != nil {
		logger.Error(err)
	}

	return nil
}

func (txq *TxQueue) PublishPendingTrades(m *types.Matches) error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)
	q := txq.Broker.GetQueue(ch, name)

	b, err := json.Marshal(m)
	if err != nil {
		return errors.New("Failed to marshal trade object")
	}

	err = txq.Broker.Publish(ch, q, b)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (txq *TxQueue) PurgePendingTrades() error {
	name := "TX_QUEUES:" + txq.Name
	ch := txq.Broker.GetChannel(name)

	err := txq.Broker.Purge(ch, name)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
