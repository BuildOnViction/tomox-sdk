package types

import (
	"encoding/json"
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

// MarshalJSON implements the json.Marshal interface
func (o *StopOrder) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{
		"exchangeAddress": o.ExchangeAddress,
		"userAddress":     o.UserAddress,
		"baseToken":       o.BaseToken,
		"quoteToken":      o.QuoteToken,
		"side":            o.Side,
		"type":            o.Type,
		"status":          o.Status,
		"pairName":        o.PairName,
		"amount":          o.Amount.String(),
		"stopPrice":       o.StopPrice.String(),
		"limitPrice":      o.LimitPrice.String(),
		"makeFee":         o.MakeFee.String(),
		"takeFee":         o.TakeFee.String(),
		"createdAt":       o.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt":       o.UpdatedAt.Format(time.RFC3339Nano),
	}

	if o.FilledAmount != nil {
		order["filledAmount"] = o.FilledAmount.String()
	}

	if o.Hash.Hex() != "" {
		order["hash"] = o.Hash.Hex()
	}

	if o.Nonce != nil {
		order["nonce"] = o.Nonce.String()
	}

	if o.Signature != nil {
		order["signature"] = map[string]interface{}{
			"V": o.Signature.V,
			"R": o.Signature.R,
			"S": o.Signature.S,
		}
	}

	return json.Marshal(order)
}

// UnmarshalJSON : write custom logic to unmarshal bytes to StopOrder
func (o *StopOrder) UnmarshalJSON(b []byte) error {
	order := map[string]interface{}{}

	err := json.Unmarshal(b, &order)
	if err != nil {
		return err
	}

	if order["id"] != nil && bson.IsObjectIdHex(order["id"].(string)) {
		o.ID = bson.ObjectIdHex(order["id"].(string))
	}

	if order["pairName"] != nil {
		o.PairName = order["pairName"].(string)
	}

	if order["exchangeAddress"] != nil {
		o.ExchangeAddress = common.HexToAddress(order["exchangeAddress"].(string))
	}

	if order["userAddress"] != nil {
		o.UserAddress = common.HexToAddress(order["userAddress"].(string))
	}

	if order["baseToken"] != nil {
		o.BaseToken = common.HexToAddress(order["baseToken"].(string))
	}

	if order["quoteToken"] != nil {
		o.QuoteToken = common.HexToAddress(order["quoteToken"].(string))
	}

	if order["stopPrice"] != nil {
		o.StopPrice = math.ToBigInt(order["stopPrice"].(string))
	}

	if order["limitPrice"] != nil {
		o.LimitPrice = math.ToBigInt(order["limitPrice"].(string))
	}

	if order["amount"] != nil {
		o.Amount = math.ToBigInt(order["amount"].(string))
	}

	if order["filledAmount"] != nil {
		o.FilledAmount = math.ToBigInt(order["filledAmount"].(string))
	}

	if order["nonce"] != nil {
		o.Nonce = math.ToBigInt(order["nonce"].(string))
	}

	if order["makeFee"] != nil {
		o.MakeFee = math.ToBigInt(order["makeFee"].(string))
	}

	if order["takeFee"] != nil {
		o.TakeFee = math.ToBigInt(order["takeFee"].(string))
	}

	if order["hash"] != nil {
		o.Hash = common.HexToHash(order["hash"].(string))
	}

	if order["side"] != nil {
		o.Side = order["side"].(string)
	}

	if order["type"] != nil {
		o.Type = order["type"].(string)
	}

	if order["status"] != nil {
		o.Status = order["status"].(string)
	}

	if order["signature"] != nil {
		signature := order["signature"].(map[string]interface{})
		o.Signature = &Signature{
			V: byte(signature["V"].(float64)),
			R: common.HexToHash(signature["R"].(string)),
			S: common.HexToHash(signature["S"].(string)),
		}
	}

	if order["createdAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, order["createdAt"].(string))
		o.CreatedAt = t
	}

	if order["updatedAt"] != nil {
		t, _ := time.Parse(time.RFC3339Nano, order["updatedAt"].(string))
		o.UpdatedAt = t
	}

	return nil
}

func (o *StopOrder) GetBSON() (interface{}, error) {
	or := StopOrderRecord{
		PairName:        o.PairName,
		ExchangeAddress: o.ExchangeAddress.Hex(),
		UserAddress:     o.UserAddress.Hex(),
		BaseToken:       o.BaseToken.Hex(),
		QuoteToken:      o.QuoteToken.Hex(),
		Status:          o.Status,
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash.Hex(),
		Amount:          o.Amount.String(),
		StopPrice:       o.StopPrice.String(),
		LimitPrice:      o.LimitPrice.String(),
		Nonce:           o.Nonce.String(),
		MakeFee:         o.MakeFee.String(),
		TakeFee:         o.TakeFee.String(),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	if o.ID.Hex() == "" {
		or.ID = bson.NewObjectId()
	} else {
		or.ID = o.ID
	}

	if o.FilledAmount != nil {
		or.FilledAmount = o.FilledAmount.String()
	}

	if o.Signature != nil {
		or.Signature = &SignatureRecord{
			V: o.Signature.V,
			R: o.Signature.R.Hex(),
			S: o.Signature.S.Hex(),
		}
	}

	return or, nil
}

func (o *StopOrder) SetBSON(raw bson.Raw) error {
	decoded := new(struct {
		ID              bson.ObjectId    `json:"id,omitempty" bson:"_id"`
		PairName        string           `json:"pairName" bson:"pairName"`
		ExchangeAddress string           `json:"exchangeAddress" bson:"exchangeAddress"`
		UserAddress     string           `json:"userAddress" bson:"userAddress"`
		BaseToken       string           `json:"baseToken" bson:"baseToken"`
		QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
		Status          string           `json:"status" bson:"status"`
		Side            string           `json:"side" bson:"side"`
		Type            string           `json:"type" bson:"type"`
		Hash            string           `json:"hash" bson:"hash"`
		StopPrice       string           `json:"stopPrice" bson:"stopPrice"`
		LimitPrice      string           `json:"limitPrice" bson:"limitPrice"`
		Amount          string           `json:"amount" bson:"amount"`
		FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
		Nonce           string           `json:"nonce" bson:"nonce"`
		MakeFee         string           `json:"makeFee" bson:"makeFee"`
		TakeFee         string           `json:"takeFee" bson:"takeFee"`
		Signature       *SignatureRecord `json:"signature" bson:"signature"`
		CreatedAt       time.Time        `json:"createdAt" bson:"createdAt"`
		UpdatedAt       time.Time        `json:"updatedAt" bson:"updatedAt"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		logger.Error(err)
		return err
	}

	o.ID = decoded.ID
	o.PairName = decoded.PairName
	o.ExchangeAddress = common.HexToAddress(decoded.ExchangeAddress)
	o.UserAddress = common.HexToAddress(decoded.UserAddress)
	o.BaseToken = common.HexToAddress(decoded.BaseToken)
	o.QuoteToken = common.HexToAddress(decoded.QuoteToken)
	o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	o.Nonce = math.ToBigInt(decoded.Nonce)
	o.MakeFee = math.ToBigInt(decoded.MakeFee)
	o.TakeFee = math.ToBigInt(decoded.TakeFee)
	o.Status = decoded.Status
	o.Side = decoded.Side
	o.Type = decoded.Type
	o.Hash = common.HexToHash(decoded.Hash)

	if decoded.Amount != "" {
		o.Amount = math.ToBigInt(decoded.Amount)
	}

	if decoded.FilledAmount != "" {
		o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	}

	if decoded.StopPrice != "" {
		o.StopPrice = math.ToBigInt(decoded.StopPrice)
	}

	if decoded.LimitPrice != "" {
		o.LimitPrice = math.ToBigInt(decoded.LimitPrice)
	}

	if decoded.Signature != nil {
		o.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	o.CreatedAt = decoded.CreatedAt
	o.UpdatedAt = decoded.UpdatedAt

	return nil
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
	if so.Type != TypeStopMarketOrder && so.Type != TypeStopLimitOrder {
		so.Type = TypeStopLimitOrder
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

// StopOrderRecord is the object that will be saved in the database
type StopOrderRecord struct {
	ID              bson.ObjectId    `json:"id" bson:"_id"`
	UserAddress     string           `json:"userAddress" bson:"userAddress"`
	ExchangeAddress string           `json:"exchangeAddress" bson:"exchangeAddress"`
	BaseToken       string           `json:"baseToken" bson:"baseToken"`
	QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
	Status          string           `json:"status" bson:"status"`
	Side            string           `json:"side" bson:"side"`
	Type            string           `json:"type" bson:"type"`
	Hash            string           `json:"hash" bson:"hash"`
	StopPrice       string           `json:"stopPrice" bson:"stopPrice"`
	LimitPrice      string           `json:"limitPrice" bson:"limitPrice"`
	Amount          string           `json:"amount" bson:"amount"`
	FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
	Nonce           string           `json:"nonce" bson:"nonce"`
	MakeFee         string           `json:"makeFee" bson:"makeFee"`
	TakeFee         string           `json:"takeFee" bson:"takeFee"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`

	PairName  string    `json:"pairName" bson:"pairName"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

type StopOrderBSONUpdate struct {
	*StopOrder
}
