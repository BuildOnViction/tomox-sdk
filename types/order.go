package types

import (
	"encoding/json"
	"fmt"
	"math/big"
	"strconv"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/utils/math"
)

const (
	BUY             = "BUY"
	SELL            = "SELL"
	TypeMarketOrder = "MO"
	TypeLimitOrder  = "LO"

	OrderStatusOpen          = "OPEN"
	OrderStatusPartialFilled = "PARTIAL_FILLED"
	OrderStatusFilled        = "FILLED"
	OrderStatusCancelled     = "CANCELLED"
)

// Order contains the data related to an order sent by the user
type Order struct {
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
	PricePoint      *big.Int       `json:"pricepoint" bson:"price"`
	Amount          *big.Int       `json:"amount" bson:"quantity"`
	FilledAmount    *big.Int       `json:"filledAmount" bson:"filledAmount"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	MakeFee         *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee         *big.Int       `json:"takeFee" bson:"takeFee"`
	PairName        string         `json:"pairName" bson:"pairName"`
	CreatedAt       time.Time      `json:"createdAt" bson:"createdAt"`
	UpdatedAt       time.Time      `json:"updatedAt" bson:"updatedAt"`
	OrderID         uint64         `json:"orderID,omitempty" bson:"orderID"`
	NextOrder       []byte         `json:"-"`
	PrevOrder       []byte         `json:"-"`
	OrderList       []byte         `json:"-"`
	Key             string         `json:"key" bson:"key"`
}

func (o *Order) String() string {
	return fmt.Sprintf("Pair: %v, Pricepoint: %v, Hash: %v", o.PairName, o.PricePoint.String(), o.Hash.Hex())
}

// TODO: Verify userAddress, baseToken, quoteToken, etc. conditions are working
func (o *Order) Validate() error {
	if o.ExchangeAddress != common.HexToAddress(app.Config.Ethereum["exchange_address"]) {
		return errors.New("Order 'exchangeAddress' parameter is incorrect")
	}

	if (o.UserAddress == common.Address{}) {
		return errors.New("Order 'userAddress' parameter is required")
	}

	if o.Nonce == nil {
		return errors.New("Order 'nonce' parameter is required")
	}

	if (o.BaseToken == common.Address{}) {
		return errors.New("Order 'baseToken' parameter is required")
	}

	if (o.QuoteToken == common.Address{}) {
		return errors.New("Order 'quoteToken' parameter is required")
	}

	if o.MakeFee == nil {
		return errors.New("Order 'makeFee' parameter is required")
	}

	if o.TakeFee == nil {
		return errors.New("Order 'takeFee' parameter is required")
	}

	if o.Amount == nil {
		return errors.New("Order 'amount' parameter is required")
	}

	if o.PricePoint == nil {
		return errors.New("Order 'pricepoint' parameter is required")
	}

	if o.Side != BUY && o.Side != SELL {
		return errors.New("Order 'side' should be 'SELL' or 'BUY', but got: '" + o.Side + "'")
	}

	if o.Signature == nil {
		return errors.New("Order 'signature' parameter is required")
	}

	if math.IsSmallerThan(o.Nonce, big.NewInt(0)) {
		return errors.New("Order 'nonce' parameter should be positive")
	}

	if math.IsEqualOrSmallerThan(o.Amount, big.NewInt(0)) {
		return errors.New("Order 'amount' parameter should be strictly positive")
	}

	if math.IsEqualOrSmallerThan(o.PricePoint, big.NewInt(0)) {
		return errors.New("Order 'pricepoint' parameter should be strictly positive")
	}

	valid, err := o.VerifySignature()
	if err != nil {
		return err
	}

	if !valid {
		return errors.New("Order 'signature' parameter is invalid")
	}

	return nil
}

// ComputeHash calculates the orderRequest hash
func (o *Order) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.ExchangeAddress.Bytes())
	sha.Write(o.UserAddress.Bytes())
	sha.Write(o.BaseToken.Bytes())
	sha.Write(o.QuoteToken.Bytes())
	sha.Write(common.BigToHash(o.Amount).Bytes())
	sha.Write(common.BigToHash(o.PricePoint).Bytes())
	sha.Write(common.BigToHash(o.EncodedSide()).Bytes())
	sha.Write(common.BigToHash(o.Nonce).Bytes())
	sha.Write(common.BigToHash(o.MakeFee).Bytes())
	sha.Write(common.BigToHash(o.TakeFee).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature checks that the orderRequest signature corresponds to the address in the userAddress field
func (o *Order) VerifySignature() (bool, error) {
	o.Hash = o.ComputeHash()

	message := crypto.Keccak256(
		[]byte("\x19Ethereum Signed Message:\n32"),
		o.Hash.Bytes(),
	)

	address, err := o.Signature.Verify(common.BytesToHash(message))
	if err != nil {
		return false, err
	}

	if address != o.UserAddress {
		return false, errors.New("Recovered address is incorrect")
	}

	return true, nil
}

// Sign first calculates the order hash, then computes a signature of this hash
// with the given wallet
func (o *Order) Sign(w *Wallet) error {
	hash := o.ComputeHash()
	sig, err := w.SignHash(hash)
	if err != nil {
		return err
	}

	o.Hash = hash
	o.Signature = sig
	return nil
}

func (o *Order) Process(p *Pair) error {
	if o.FilledAmount == nil {
		o.FilledAmount = big.NewInt(0)
	}

	// TODO: Handle this in Validate function
	if o.Type != TypeMarketOrder && o.Type != TypeLimitOrder {
		o.Type = TypeLimitOrder
	}

	if !math.IsEqual(o.MakeFee, p.MakeFee) {
		return errors.New("Invalid MakeFee")
	}

	if !math.IsEqual(o.TakeFee, p.TakeFee) {
		return errors.New("Invalid TakeFee")
	}

	o.PairName = p.Name()
	o.CreatedAt = time.Now()
	o.UpdatedAt = time.Now()
	return nil
}

func (o *Order) Pair() (*Pair, error) {
	if (o.BaseToken == common.Address{}) {
		return nil, errors.New("Base token is not set")
	}

	if (o.QuoteToken == common.Address{}) {
		return nil, errors.New("Quote token is set")
	}

	return &Pair{
		BaseTokenAddress:  o.BaseToken,
		QuoteTokenAddress: o.QuoteToken,
	}, nil
}

func (o *Order) RemainingAmount() *big.Int {
	return math.Sub(o.Amount, o.FilledAmount)
}

func (o *Order) SellTokenSymbol() string {
	if o.Side == BUY {
		return o.QuoteTokenSymbol()
	}

	if o.Side == SELL {
		return o.BaseTokenSymbol()
	}

	return ""
}

//TODO handle error case
func (o *Order) SellToken() common.Address {
	if o.Side == BUY {
		return o.QuoteToken
	} else {
		return o.BaseToken
	}
}

func (o *Order) BuyToken() common.Address {
	if o.Side == BUY {
		return o.BaseToken
	} else {
		return o.QuoteToken
	}
}

func (o *Order) QuoteAmount(p *Pair) *big.Int {
	pairMultiplier := p.PairMultiplier()
	return math.Div(math.Mul(o.Amount, o.PricePoint), pairMultiplier)
}

// SellAmount
// If order is a "BUY", then sellToken = quoteToken
func (o *Order) SellAmount(p *Pair) *big.Int {
	pairMultiplier := p.PairMultiplier()

	if o.Side == BUY {
		return math.Div(math.Mul(o.Amount, o.PricePoint), pairMultiplier)
	} else {
		return o.Amount
	}
}

func (o *Order) RemainingSellAmount(p *Pair) *big.Int {
	pairMultiplier := p.PairMultiplier()

	if o.Side == BUY {
		remainingAmount := math.Sub(o.Amount, o.FilledAmount)
		return math.Div(math.Mul(remainingAmount, o.PricePoint), pairMultiplier)
	} else {
		return math.Sub(o.Amount, o.FilledAmount)
	}
}

func (o *Order) RequiredSellAmount(p *Pair) *big.Int {
	var requiredSellTokenAmount *big.Int

	pairMultiplier := p.PairMultiplier()

	if o.Side == BUY {
		requiredSellTokenAmount = math.Div(math.Mul(o.Amount, o.PricePoint), pairMultiplier)
	} else {
		requiredSellTokenAmount = o.Amount
	}

	return requiredSellTokenAmount
}

func (o *Order) TotalRequiredSellAmount(p *Pair) *big.Int {
	var requiredSellTokenAmount *big.Int

	pairMultiplier := p.PairMultiplier()

	if o.Side == BUY {
		sellAmount := math.Div(math.Mul(o.Amount, o.PricePoint), pairMultiplier)
		fee := math.Max(p.MakeFee, p.TakeFee)
		requiredSellTokenAmount = math.Add(sellAmount, fee)
	} else {
		requiredSellTokenAmount = o.Amount
	}

	return requiredSellTokenAmount
}

func (o *Order) BuyAmount(pairMultiplier *big.Int) *big.Int {
	if o.Side == SELL {
		return o.Amount
	} else {
		return math.Div(math.Mul(o.Amount, o.PricePoint), pairMultiplier)
	}
}

//TODO handle error case ?
func (o *Order) EncodedSide() *big.Int {
	if o.Side == BUY {
		return big.NewInt(0)
	} else {
		return big.NewInt(1)
	}
}

func (o *Order) BuyTokenSymbol() string {
	if o.Side == BUY {
		return o.BaseTokenSymbol()
	}

	if o.Side == SELL {
		return o.QuoteTokenSymbol()
	}

	return ""
}

func (o *Order) PairCode() (string, error) {
	if o.PairName == "" {
		return "", errors.New("Pair name is required")
	}

	return o.PairName + "::" + o.BaseToken.Hex() + "::" + o.QuoteToken.Hex(), nil
}

func (o *Order) BaseTokenSymbol() string {
	if o.PairName == "" {
		return ""
	}

	return o.PairName[:strings.IndexByte(o.PairName, '/')]
}

func (o *Order) QuoteTokenSymbol() string {
	if o.PairName == "" {
		return ""
	}

	return o.PairName[strings.IndexByte(o.PairName, '/')+1:]
}

// JSON Marshal/Unmarshal interface

// MarshalJSON implements the json.Marshal interface
func (o *Order) MarshalJSON() ([]byte, error) {
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
		"pricepoint":      o.PricePoint.String(),
		"makeFee":         o.MakeFee.String(),
		"takeFee":         o.TakeFee.String(),
		"createdAt":       o.CreatedAt.Format(time.RFC3339Nano),
		"updatedAt":       o.UpdatedAt.Format(time.RFC3339Nano),
		"orderID":         strconv.FormatUint(o.OrderID, 10),
		"key":             o.Key,
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

// UnmarshalJSON : write custom logic to unmarshal bytes to Order
func (o *Order) UnmarshalJSON(b []byte) error {
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

	if order["pricepoint"] != nil {
		o.PricePoint = math.ToBigInt(order["pricepoint"].(string))
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

	if order["orderID"] != nil {
		o.Status = order["status"].(string)
		orderID, err := strconv.ParseInt(order["orderID"].(string), 10, 64)
		if err != nil {
			logger.Error(err)
		}
		o.OrderID = uint64(orderID)
	}

	if order["key"] != nil {
		o.Key = order["key"].(string)
	}

	return nil
}

func (o *Order) GetBSON() (interface{}, error) {
	or := OrderRecord{
		PairName:        o.PairName,
		ExchangeAddress: o.ExchangeAddress.Hex(),
		UserAddress:     o.UserAddress.Hex(),
		BaseToken:       o.BaseToken.Hex(),
		QuoteToken:      o.QuoteToken.Hex(),
		Status:          o.Status,
		Side:            o.Side,
		Type:            o.Type,
		Hash:            o.Hash.Hex(),
		Quantity:        o.Amount.String(),
		Price:           o.PricePoint.String(),
		Nonce:           o.Nonce.String(),
		MakeFee:         o.MakeFee.String(),
		TakeFee:         o.TakeFee.String(),
		CreatedAt:       strconv.FormatInt(o.CreatedAt.Unix(), 10),
		UpdatedAt:       strconv.FormatInt(o.UpdatedAt.Unix(), 10),
		OrderID:         strconv.FormatUint(o.OrderID, 10),
		NextOrder:       common.Bytes2Hex(o.NextOrder),
		PrevOrder:       common.Bytes2Hex(o.PrevOrder),
		OrderList:       common.Bytes2Hex(o.OrderList),
		Key:             o.Key,
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

func (o *Order) SetBSON(raw bson.Raw) error {
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
		Price           string           `json:"price" bson:"price"`
		Quantity        string           `json:"quantity" bson:"quantity"`
		FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
		Nonce           string           `json:"nonce" bson:"nonce"`
		MakeFee         string           `json:"makeFee" bson:"makeFee"`
		TakeFee         string           `json:"takeFee" bson:"takeFee"`
		Signature       *SignatureRecord `json:"signature" bson:"signature"`
		CreatedAt       string           `json:"createdAt" bson:"createdAt"`
		UpdatedAt       string           `json:"updatedAt" bson:"updatedAt"`
		OrderID         string           `json:"orderID" bson:"orderID"`
		NextOrder       string           `json:"-"`
		PrevOrder       string           `json:"-"`
		OrderList       string           `json:"-"`
		Key             string           `json:"key" bson:"key"`
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

	if decoded.Quantity != "" {
		o.Amount = math.ToBigInt(decoded.Quantity)
	}

	if decoded.FilledAmount != "" {
		o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	}

	if decoded.Price != "" {
		o.PricePoint = math.ToBigInt(decoded.Price)
	}

	if decoded.Signature != nil {
		o.Signature = &Signature{
			V: byte(decoded.Signature.V),
			R: common.HexToHash(decoded.Signature.R),
			S: common.HexToHash(decoded.Signature.S),
		}
	}

	createdAt, err := strconv.ParseInt(decoded.CreatedAt, 10, 64)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	o.CreatedAt = time.Unix(createdAt, 0)

	updatedAt, err := strconv.ParseInt(decoded.UpdatedAt, 10, 64)
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	o.UpdatedAt = time.Unix(updatedAt, 0)

	orderID, err := strconv.ParseInt(decoded.OrderID, 10, 64)
	if err != nil {
		logger.Error(err)
	}
	o.OrderID = uint64(orderID)
	o.NextOrder = common.Hex2Bytes(decoded.NextOrder)
	o.PrevOrder = common.Hex2Bytes(decoded.PrevOrder)
	o.OrderList = common.Hex2Bytes(decoded.OrderList)
	o.Key = decoded.Key

	return nil
}

// OrderRecord is the object that will be saved in the database
type OrderRecord struct {
	ID              bson.ObjectId    `json:"id" bson:"_id"`
	UserAddress     string           `json:"userAddress" bson:"userAddress"`
	ExchangeAddress string           `json:"exchangeAddress" bson:"exchangeAddress"`
	BaseToken       string           `json:"baseToken" bson:"baseToken"`
	QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
	Status          string           `json:"status" bson:"status"`
	Side            string           `json:"side" bson:"side"`
	Type            string           `json:"type" bson:"type"`
	Hash            string           `json:"hash" bson:"hash"`
	Price           string           `json:"price" bson:"price"`
	Quantity        string           `json:"quantity" bson:"quantity"`
	FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
	Nonce           string           `json:"nonce" bson:"nonce"`
	MakeFee         string           `json:"makeFee" bson:"makeFee"`
	TakeFee         string           `json:"takeFee" bson:"takeFee"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`
	PairName        string           `json:"pairName" bson:"pairName"`
	CreatedAt       string           `json:"createdAt" bson:"createdAt"`
	UpdatedAt       string           `json:"updatedAt" bson:"updatedAt"`
	OrderID         string           `json:"orderID,omitempty" bson:"orderID"`
	NextOrder       string           `json:"nextOrder,omitempty" bson:"nextOrder"`
	PrevOrder       string           `json:"prevOrder,omitempty" bson:"prevOrder"`
	OrderList       string           `json:"orderList,omitempty" bson:"orderList"`
	Key             string           `json:"key" bson:"key"`
}

type OrderBSONUpdate struct {
	*Order
}

func (o OrderBSONUpdate) GetBSON() (interface{}, error) {
	now := time.Now()

	set := bson.M{
		"pairName":        o.PairName,
		"exchangeAddress": o.ExchangeAddress.Hex(),
		"userAddress":     o.UserAddress.Hex(),
		"baseToken":       o.BaseToken.Hex(),
		"quoteToken":      o.QuoteToken.Hex(),
		"status":          o.Status,
		"side":            o.Side,
		"type":            o.Type,
		"pricepoint":      o.PricePoint.String(),
		"amount":          o.Amount.String(),
		"nonce":           o.Nonce.String(),
		"makeFee":         o.MakeFee.String(),
		"takeFee":         o.TakeFee.String(),
		"updatedAt":       now,
	}

	if o.FilledAmount != nil {
		set["filledAmount"] = o.FilledAmount.String()
	}

	if o.Signature != nil {
		set["signature"] = bson.M{
			"V": o.Signature.V,
			"R": o.Signature.R.Hex(),
			"S": o.Signature.S.Hex(),
		}
	}

	setOnInsert := bson.M{
		"_id":       bson.NewObjectId(),
		"hash":      o.Hash.Hex(),
		"createdAt": now,
	}

	update := bson.M{
		"$set":         set,
		"$setOnInsert": setOnInsert,
	}

	return update, nil
}

type OrderData struct {
	Pair        PairID   `json:"id,omitempty" bson:"_id"`
	OrderVolume *big.Int `json:"orderVolume,omitempty" bson:"orderVolume"`
	OrderCount  *big.Int `json:"orderCount,omitempty" bson:"orderCount"`
	BestPrice   *big.Int `json:"bestPrice,omitempty" bson:"bestPrice"`
}

func (o *OrderData) AddressCode() string {
	code := o.Pair.BaseToken.Hex() + "::" + o.Pair.QuoteToken.Hex()
	return code
}

func (o *OrderData) ConvertedVolume(p *Pair, exchangeRate float64) float64 {
	valueAsToken := math.DivideToFloat(o.OrderVolume, p.BaseTokenMultiplier())
	value := valueAsToken / exchangeRate

	return value
}

func (o *OrderData) MarshalJSON() ([]byte, error) {
	orderData := map[string]interface{}{
		"pair": map[string]interface{}{
			"pairName":   o.Pair.PairName,
			"baseToken":  o.Pair.BaseToken.Hex(),
			"quoteToken": o.Pair.QuoteToken.Hex(),
		},
	}

	if o.OrderVolume != nil {
		orderData["orderVolume"] = o.OrderVolume.String()
	}

	if o.OrderCount != nil {
		orderData["orderCount"] = o.OrderCount.String()
	}

	if o.BestPrice != nil {
		orderData["bestPrice"] = o.BestPrice.String()
	}

	bytes, err := json.Marshal(orderData)
	return bytes, err
}

// UnmarshalJSON creates a trade object from a json byte string
func (o *OrderData) UnmarshalJSON(b []byte) error {
	orderData := map[string]interface{}{}
	err := json.Unmarshal(b, &orderData)

	if err != nil {
		return err
	}

	if orderData["pair"] != nil {
		pair := orderData["pair"].(map[string]interface{})
		o.Pair = PairID{
			PairName:   pair["pairName"].(string),
			BaseToken:  common.HexToAddress(pair["baseToken"].(string)),
			QuoteToken: common.HexToAddress(pair["quoteToken"].(string)),
		}
	}

	if orderData["orderVolume"] != nil {
		o.OrderVolume = math.ToBigInt(orderData["orderVolume"].(string))
	}

	if orderData["orderCount"] != nil {
		o.OrderCount = math.ToBigInt(orderData["orderCount"].(string))
	}

	if orderData["bestPrice"] != nil {
		o.BestPrice = math.ToBigInt(orderData["bestPrice"].(string))
	}

	return nil
}

func (o *OrderData) GetBSON() (interface{}, error) {
	type PairID struct {
		PairName   string `json:"pairName" bson:"pairName"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}

	count, err := bson.ParseDecimal128(o.OrderCount.String())
	if err != nil {
		return nil, err
	}

	volume, err := bson.ParseDecimal128(o.OrderVolume.String())
	if err != nil {
		return nil, err
	}

	bestPrice := o.BestPrice.String()
	if err != nil {
		return nil, err
	}

	return struct {
		ID          PairID          `json:"id,omitempty" bson:"_id"`
		OrderVolume bson.Decimal128 `json:"orderCount" bson:"orderCount"`
		OrderCount  bson.Decimal128 `json:"orderVolume" bson:"orderVolume"`
		BestPrice   string          `json:"bestPrice" bson:"bestPrice"`
	}{
		ID: PairID{
			o.Pair.PairName,
			o.Pair.BaseToken.Hex(),
			o.Pair.QuoteToken.Hex(),
		},
		OrderVolume: volume,
		OrderCount:  count,
		BestPrice:   bestPrice,
	}, nil
}

func (o *OrderData) SetBSON(raw bson.Raw) error {
	type PairIDRecord struct {
		PairName   string `json:"pairName" bson:"pairName"`
		BaseToken  string `json:"baseToken" bson:"baseToken"`
		QuoteToken string `json:"quoteToken" bson:"quoteToken"`
	}

	decoded := new(struct {
		Pair        PairIDRecord    `json:"pair,omitempty" bson:"_id"`
		OrderCount  bson.Decimal128 `json:"orderCount" bson:"orderCount"`
		OrderVolume bson.Decimal128 `json:"orderVolume" bson:"orderVolume"`
		BestPrice   string          `json:"bestPrice" bson:"bestPrice"`
	})

	err := raw.Unmarshal(decoded)
	if err != nil {
		return err
	}

	o.Pair = PairID{
		PairName:   decoded.Pair.PairName,
		BaseToken:  common.HexToAddress(decoded.Pair.BaseToken),
		QuoteToken: common.HexToAddress(decoded.Pair.QuoteToken),
	}

	orderCount := decoded.OrderCount.String()
	orderVolume := decoded.OrderVolume.String()
	bestPrice := decoded.BestPrice

	o.OrderCount = math.ToBigInt(orderCount)
	o.OrderVolume = math.ToBigInt(orderVolume)
	o.BestPrice = math.ToBigInt(bestPrice)

	return nil
}

type updateDesc struct {
	UpdatedFields map[string]interface{} `bson:"updatedFields"`
	RemovedFields []string               `bson:"removedFields"`
}

type evNamespace struct {
	DB   string `bson:"db"`
	Coll string `bson:"coll"`
}

type M bson.M

type OrderChangeEvent struct {
	ID                interface{} `bson:"_id"`
	OperationType     string      `bson:"operationType"`
	FullDocument      *Order      `bson:"fullDocument,omitempty"`
	Ns                evNamespace `bson:"ns"`
	DocumentKey       M           `bson:"documentKey"`
	UpdateDescription *updateDesc `bson:"updateDescription,omitempty"`
}

const (
	OPERATION_TYPE_INSERT  = "insert"
	OPERATION_TYPE_UPDATE  = "update"
	OPERATION_TYPE_REPLACE = "replace"
	OPERATION_TYPE_DELETE  = "delete"
)
