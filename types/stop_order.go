package types

import (
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomodex/app"
	"github.com/tomochain/tomodex/errors"
	"github.com/tomochain/tomodex/utils/math"
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

// TODO: Verify userAddress, baseToken, quoteToken, etc. conditions are working
func (so *StopOrder) Validate() error {
	if so.ExchangeAddress != common.HexToAddress(app.Config.Ethereum["exchange_address"]) {
		return errors.New("Order 'exchangeAddress' parameter is incorrect")
	}

	if (so.UserAddress == common.Address{}) {
		return errors.New("Order 'userAddress' parameter is required")
	}

	if so.Nonce == nil {
		return errors.New("Order 'nonce' parameter is required")
	}

	if (so.BaseToken == common.Address{}) {
		return errors.New("Order 'baseToken' parameter is required")
	}

	if (so.QuoteToken == common.Address{}) {
		return errors.New("Order 'quoteToken' parameter is required")
	}

	if so.MakeFee == nil {
		return errors.New("Order 'makeFee' parameter is required")
	}

	if so.TakeFee == nil {
		return errors.New("Order 'takeFee' parameter is required")
	}

	if so.Amount == nil {
		return errors.New("Order 'amount' parameter is required")
	}

	if so.StopPrice == nil {
		return errors.New("Order 'stopPrice' parameter is required")
	}

	if so.Side != BUY && so.Side != SELL {
		return errors.New("Order 'side' should be 'SELL' or 'BUY', but got: '" + so.Side + "'")
	}

	if so.Signature == nil {
		return errors.New("Order 'signature' parameter is required")
	}

	if math.IsSmallerThan(so.Nonce, big.NewInt(0)) {
		return errors.New("Order 'nonce' parameter should be positive")
	}

	if math.IsEqualOrSmallerThan(so.Amount, big.NewInt(0)) {
		return errors.New("Order 'amount' parameter should be strictly positive")
	}

	if math.IsEqualOrSmallerThan(so.StopPrice, big.NewInt(0)) {
		return errors.New("Order 'stopPrice' parameter should be strictly positive")
	}

	valid, err := so.VerifySignature()
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("Order 'signature' parameter is invalid")
	}

	return nil
}

// ComputeHash calculates the orderRequest hash
func (so *StopOrder) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(so.ExchangeAddress.Bytes())
	sha.Write(so.UserAddress.Bytes())
	sha.Write(so.BaseToken.Bytes())
	sha.Write(so.QuoteToken.Bytes())
	sha.Write(common.BigToHash(so.Amount).Bytes())
	sha.Write(common.BigToHash(so.StopPrice).Bytes())
	sha.Write(common.BigToHash(so.EncodedSide()).Bytes())
	sha.Write(common.BigToHash(so.Nonce).Bytes())
	sha.Write(common.BigToHash(so.TakeFee).Bytes())
	sha.Write(common.BigToHash(so.MakeFee).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (so *StopOrder) VerifySignature() (bool, error) {
	so.Hash = so.ComputeHash()

	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		so.Hash.Bytes(),
	)

	address, err := so.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != so.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

func (so *StopOrder) Process(p *Pair) error {
	if so.FilledAmount == nil {
		so.FilledAmount = big.NewInt(0)
	}

	// TODO: Handle this in Validate function
	if so.Type != TypeMarketOrder && so.Type != TypeLimitOrder {
		so.Type = TypeLimitOrder
	}

	if !math.IsEqual(so.MakeFee, p.MakeFee) {
		return errors.New("Invalid MakeFee")
	}

	if !math.IsEqual(so.TakeFee, p.TakeFee) {
		return errors.New("Invalid TakeFee")
	}

	so.PairName = p.Name()
	so.CreatedAt = time.Now()
	so.UpdatedAt = time.Now()
	return nil
}

func (so *StopOrder) QuoteAmount(p *Pair) *big.Int {
	pairMultiplier := p.PairMultiplier()
	return math.Div(math.Mul(so.Amount, so.StopPrice), pairMultiplier)
}

//TODO handle error case ?
func (so *StopOrder) EncodedSide() *big.Int {
	if so.Side == BUY {
		return big.NewInt(0)
	} else {
		return big.NewInt(1)
	}
}

func (so *StopOrder) PairCode() (string, error) {
	if so.PairName == "" {
		return "", errors.New("Pair name is required")
	}

	return so.PairName + "::" + so.BaseToken.Hex() + "::" + so.QuoteToken.Hex(), nil
}

type StopOrderBSONUpdate struct {
	*StopOrder
}
