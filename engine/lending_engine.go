package engine

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/daos"
	"github.com/tomochain/tomox-sdk/types"
)

// LendingEngine matching lending engine
type LendingEngine struct {
	lendingDao       *daos.LendingOrderDao
	lendingTradeDao  *daos.LendingTradeDao
	openLendingOrder map[common.Hash]*types.LendingOrder
}

// NewLendingEngine initializes the engine singleton instance
func NewLendingEngine(
	lendingDao *daos.LendingOrderDao,
	lendingTradeDao *daos.LendingTradeDao,
) *LendingEngine {
	engine := &LendingEngine{
		lendingDao:       lendingDao,
		lendingTradeDao:  lendingTradeDao,
		openLendingOrder: make(map[common.Hash]*types.LendingOrder),
	}
	return engine
}

// HandleNewOrder handle new order comming
func (e *LendingEngine) HandleNewOrder(o *types.LendingOrder) error {
	o.Status = "OPEN"
	o.FilledAmount = big.NewInt(0)
	e.lendingDao.Create(o)
	if len(e.openLendingOrder) == 0 {
		e.openLendingOrder[o.Hash] = o
	} else {
		for _, order := range e.openLendingOrder {
			if order.Side != o.Side && order.Term == o.Term {
				if (order.Side == types.BORROW && order.Interest >= o.Interest) || (order.Side == types.LEND && order.Interest <= o.Interest) {
					var lendingTrade types.LendingTrade
					if order.Side == types.BORROW {
						lendingTrade.BorrowingOwner = order.UserAddress
						lendingTrade.InvestingOwner = o.UserAddress
						lendingTrade.InvestingHash = o.Hash
						lendingTrade.BorrowingHash = order.Hash
						lendingTrade.InvestingRelayer = o.RelayerAddress
						lendingTrade.BorrowingRelayer = order.RelayerAddress
						lendingTrade.CollateralToken = order.CollateralToken
						lendingTrade.LendingToken = o.LendingToken
					}
					if order.Side == types.LEND {
						lendingTrade.BorrowingOwner = o.UserAddress
						lendingTrade.InvestingOwner = order.UserAddress
						lendingTrade.InvestingHash = order.Hash
						lendingTrade.BorrowingHash = o.Hash
						lendingTrade.InvestingRelayer = order.RelayerAddress
						lendingTrade.BorrowingRelayer = o.RelayerAddress
						lendingTrade.LendingToken = order.LendingToken
					}
					lendingTrade.Term = order.Term
					lendingTrade.Interest = order.Interest
					lendingTrade.InvestingFee = big.NewInt(1)
					lendingTrade.BorrowingFee = big.NewInt(1)
					lendingTrade.CollateralPrice = big.NewInt(0)
					lendingTrade.LiquidationPrice = big.NewInt(0)
					lendingTrade.Status = types.TradeStatusSuccess
					remain := big.NewInt(0).Sub(order.Quantity, order.FilledAmount)
					if remain.Cmp(o.Quantity) <= 0 {
						o.FilledAmount.Add(o.FilledAmount, remain)
						e.openLendingOrder[o.Hash] = o
						lendingTrade.Amount = remain
					} else {
						o.FilledAmount = o.Quantity
						order.FilledAmount.Add(order.FilledAmount, o.Quantity)
						lendingTrade.Amount = o.Quantity
					}
					e.lendingTradeDao.Create(&lendingTrade)

				}

			}
			e.lendingDao.UpdateFilledAmount(order.Hash, order.Quantity)
		}
		if o.FilledAmount.Cmp(big.NewInt(0)) != 0 {
			e.lendingDao.UpdateFilledAmount(o.Hash, o.FilledAmount)
		}
	}

	return nil
}

// HandleCancelOrder handle new order comming
func (e *LendingEngine) HandleCancelOrder(order types.LendingOrder) error {
	return nil
}
