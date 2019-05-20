package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomodex/errors"
)

const (
	TypeStopMarketOrder = "SMO"
	TypeStopLimitOrder  = "SLO"
)

type StopOrder struct {
	ID              bson.ObjectId  `json:"id" bson:"_id"`
	UserAddress     common.Address `json:"userAddress" bson:"userAddress"`
	ExchangeAddress common.Address `json:"exchangeAddress" bson:"exchangeAddress"`
	BaseToken       common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken      common.Address `json:"quoteToken" bson:"quoteToken"`
	Status          string         `json:"status" bson:"status"`
	Side            string         `json:"side" bson:"side"`
	Type            string         `json:"type" bson:"type"`
	Hash            common.Hash    `json:"hash" bson:"hash"`
	Signature       *Signature     `json:"signature,omitempty" bson:"signature"`
	StopPrice       *big.Int       `json:"stopPrice" bson:"stopPrice"`
	LimitPrice      *big.Int       `json:"limitPrice" bson:"limitPrice"`
	Amount          *big.Int       `json:"amount" bson:"amount"`
	FilledAmount    *big.Int       `json:"filledAmount" bson:"filledAmount"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	MakeFee         *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee         *big.Int       `json:"takeFee" bson:"takeFee"`
	PairName        string         `json:"pairName" bson:"pairName"`
	CreatedAt       time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt" bson:"updatedAt"`
}

// ToOrder converts a stop order to a real order that will be pushed to TomoX
func (so *StopOrder) ToOrder() (*Order, error) {
	var o *Order

	switch so.Type {
	case TypeStopMarketOrder:
		o = &Order{
			UserAddress:     so.UserAddress,
			ExchangeAddress: so.ExchangeAddress,
			BaseToken:       so.BaseToken,
			QuoteToken:      so.QuoteToken,
			Status:          OrderStatusOpen,
			Side:            so.Side,
			Type:            TypeMarketOrder,
			Hash:            so.Hash,
			Signature:       so.Signature,
			PricePoint:      so.StopPrice,
			Amount:          so.Amount,
			FilledAmount:    big.NewInt(0),
			Nonce:           so.Nonce,
			MakeFee:         so.MakeFee,
			TakeFee:         so.TakeFee,
			PairName:        so.PairName,
		}

		break
	case TypeStopLimitOrder:
		o = &Order{
			UserAddress:     so.UserAddress,
			ExchangeAddress: so.ExchangeAddress,
			BaseToken:       so.BaseToken,
			QuoteToken:      so.QuoteToken,
			Status:          OrderStatusOpen,
			Side:            so.Side,
			Type:            TypeLimitOrder,
			Hash:            so.Hash,
			Signature:       so.Signature,
			PricePoint:      so.LimitPrice,
			Amount:          so.Amount,
			FilledAmount:    big.NewInt(0),
			Nonce:           so.Nonce,
			MakeFee:         so.MakeFee,
			TakeFee:         so.TakeFee,
			PairName:        so.PairName,
		}

		break
	default:
		return nil, errors.New("Unknown stop order type")
	}

	return o, nil
}
