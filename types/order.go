package types

import (
	"encoding/json"
	"errors"
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"github.com/tomochain/backend-matching-engine/app"
	"github.com/tomochain/backend-matching-engine/utils"
	"github.com/tomochain/backend-matching-engine/utils/math"
	"gopkg.in/mgo.v2/bson"
)

// Order contains the data related to an order sent by the user
type Order struct {
	ID              bson.ObjectId  `json:"id" bson:"_id"`
	UserAddress     common.Address `json:"userAddress" bson:"userAddress"`
	ExchangeAddress common.Address `json:"exchangeAddress" bson:"exchangeAddress"`
	BuyToken        common.Address `json:"buyToken" bson:"buyToken"`
	SellToken       common.Address `json:"sellToken" bson:"sellToken"`
	BaseToken       common.Address `json:"baseToken" bson:"baseToken"`
	QuoteToken      common.Address `json:"quoteToken" bson:"quoteToken"`
	BuyAmount       *big.Int       `json:"buyAmount" bson:"buyAmount"`
	SellAmount      *big.Int       `json:"sellAmount" bson:"sellAmount"`
	Status          string         `json:"status" bson:"status"`
	Side            string         `json:"side" bson:"side"`
	Hash            common.Hash    `json:"hash" bson:"hash"`
	Signature       *Signature     `json:"signature,omitempty" bson:"signature"`
	PricePoint      *big.Int       `json:"pricepoint" bson:"pricepoint"`
	Amount          *big.Int       `json:"amount" bson:"amount"`
	FilledAmount    *big.Int       `json:"filledAmount" bson:"filledAmount"`
	Nonce           *big.Int       `json:"nonce" bson:"nonce"`
	Expires         *big.Int       `json:"expires" bson:"expires"`
	MakeFee         *big.Int       `json:"makeFee" bson:"makeFee"`
	TakeFee         *big.Int       `json:"takeFee" bson:"takeFee"`
	PairName        string         `json:"pairName" bson:"pairName"`

	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

// Validate : validate a valid order
func (o *Order) Validate() error {

	if o.ExchangeAddress != common.HexToAddress(app.Config.Ethereum["exchange_address"]) {
		return errors.New("Incorrect exchange address")
	}

	if math.IsSmallerThan(o.BuyAmount, big.NewInt(0)) {
		return errors.New("Buy amount should be positive")
	}

	if math.IsSmallerThan(o.SellAmount, big.NewInt(0)) {
		return errors.New("Sell amount should be positive")
	}

	if math.IsSmallerThan(o.Nonce, big.NewInt(0)) {
		return errors.New("Nonce should be positive")
	}

	return nil
}

// ComputeHash: calculates the orderRequest hash, should calculate at server and need client to sign
func (o *Order) ComputeHash() common.Hash {
	sha := sha3.NewKeccak256()
	sha.Write(o.ExchangeAddress.Bytes())
	sha.Write(o.UserAddress.Bytes())
	sha.Write(o.SellToken.Bytes())
	sha.Write(o.BuyToken.Bytes())
	sha.Write(common.BigToHash(o.SellAmount).Bytes())
	sha.Write(common.BigToHash(o.BuyAmount).Bytes())
	sha.Write(common.BigToHash(o.MakeFee).Bytes())
	sha.Write(common.BigToHash(o.TakeFee).Bytes())
	sha.Write(common.BigToHash(o.Expires).Bytes())
	sha.Write(common.BigToHash(o.Nonce).Bytes())
	return common.BytesToHash(sha.Sum(nil))
}

// VerifySignature: checks that the orderRequest signature corresponds to the address in the userAddress field
// If client send the correct signature then update to swarm, need to have a FeedID in order at mongo
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

	if o.BuyToken == p.BaseTokenAddress {
		o.Side = "BUY"
		o.Amount = o.BuyAmount
		o.PricePoint = math.Div(math.Mul(o.SellAmount, p.PriceMultiplier), o.BuyAmount)
	} else if o.BuyToken == p.QuoteTokenAddress {
		o.Side = "SELL"
		o.Amount = o.SellAmount
		o.PricePoint = math.Div(math.Mul(o.BuyAmount, p.PriceMultiplier), o.SellAmount)
	} else {
		return errors.New("Could not determine o side")
	}

	o.BaseToken = p.BaseTokenAddress
	o.QuoteToken = p.QuoteTokenAddress
	o.PairName = p.Name()
	return nil
}

func (o *Order) PairCode() (string, error) {
	if o.PairName == "" {
		return "", errors.New("Pair name is required")
	}

	return o.PairName + "::" + o.BaseToken.Hex() + "::" + o.QuoteToken.Hex(), nil
}

// GetKVPrefix returns the key value store(redis) prefix to be used
// by matching engine correspondind to a particular order.
func (o *Order) GetKVPrefix() string {
	return o.BaseToken.Hex() + "::" + o.QuoteToken.Hex()
}

// GetOBKeys returns the keys corresponding to an order
// orderbook price point key
// orderbook list key corresponding to order price.
func (o *Order) GetOBKeys() (ss, list string) {
	var k string
	if o.Side == "BUY" {
		k = "BUY"
	} else if o.Side == "SELL" {
		k = "SELL"
	}

	ss = o.GetKVPrefix() + "::" + k
	list = o.GetKVPrefix() + "::" + k + "::" + utils.UintToPaddedString(o.PricePoint.Int64())
	return
}

// GetOBMatchKey returns the orderbook price point key
// aginst which the order needs to be matched
func (o *Order) GetOBMatchKey() (ss string) {
	var k string
	if o.Side == "BUY" {
		k = "SELL"
	} else if o.Side == "SELL" {
		k = "BUY"
	}

	ss = o.GetKVPrefix() + "::" + k
	return
}

// JSON Marshal/Unmarshal interface

// MarshalJSON implements the json.Marshal interface
func (o *Order) MarshalJSON() ([]byte, error) {
	order := map[string]interface{}{
		"exchangeAddress": o.ExchangeAddress,
		"userAddress":     o.UserAddress,
		"buyToken":        o.BuyToken,
		"sellToken":       o.SellToken,
		"baseToken":       o.BaseToken,
		"quoteToken":      o.QuoteToken,
		"side":            o.Side,
		"status":          o.Status,
		"pairName":        o.PairName,
		"buyAmount":       o.BuyAmount.String(),
		"sellAmount":      o.SellAmount.String(),
		"makeFee":         o.MakeFee.String(),
		"takeFee":         o.TakeFee.String(),
		"expires":         o.Expires.String(),
		// NOTE: Currently removing this to simplify public API, might reinclude
		// later. An alternative would be to create additional simplified type
		"createdAt": o.CreatedAt.Format(time.RFC3339Nano),
		// "updatedAt": o.UpdatedAt.Format(time.RFC3339Nano),
	}

	// NOTE: Currently removing this to simplify public API, will reinclude
	// if needed. An alternative would be to create additional simplified type
	// if o.ID != bson.ObjectId("") {
	// 	order["id"] = o.ID
	// }

	if o.Amount != nil {
		order["amount"] = o.Amount.String()
	}

	if o.FilledAmount != nil {
		order["filledAmount"] = o.FilledAmount.String()
	}

	if o.PricePoint != nil {
		order["pricepoint"] = o.PricePoint.String()
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

	if order["buyToken"] != nil {
		o.BuyToken = common.HexToAddress(order["buyToken"].(string))
	}

	if order["sellToken"] != nil {
		o.SellToken = common.HexToAddress(order["sellToken"].(string))
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

	if order["buyAmount"] != nil {
		o.BuyAmount = math.ToBigInt(order["buyAmount"].(string))
	}

	if order["sellAmount"] != nil {
		o.SellAmount = math.ToBigInt(order["sellAmount"].(string))
	}

	if order["expires"] != nil {
		o.Expires = math.ToBigInt(order["expires"].(string))
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

// OrderFeed is the object that will be retrieved from swarm feed
type OrderFeed struct {
	Amount          *big.Int `json:"amount"`
	BaseToken       []byte   `json:"baseToken"`
	BuyAmount       *big.Int `json:"buyAmount"`
	BuyToken        []byte   `json:"buyToken"`
	ExchangeAddress []byte   `json:"exchangeAddress"`
	Expires         *big.Int `json:"expires"`
	FilledAmount    *big.Int `json:"filledAmount"`
	Hash            []byte   `json:"hash"`
	ID              []byte   `json:"id"`
	MakeFee         *big.Int `json:"makeFee"`
	Nonce           *big.Int `json:"nonce"`
	PairName        string   `json:"pairName"`
	PricePoint      *big.Int `json:"pricepoint"`
	QuoteToken      []byte   `json:"quoteToken"`
	SellAmount      *big.Int `json:"sellAmount"`
	SellToken       []byte   `json:"sellToken"`
	Side            string   `json:"side"`
	Signature       []byte   `json:"signature"`
	Status          string   `json:"status"`
	TakeFee         *big.Int `json:"takeFee"`
	Timestamp       uint     `json:"timestamp"`
	UserAddress     []byte   `json:"userAddress"`
}

func (o *OrderFeed) GetBSON() (*OrderRecord, error) {

	timestamp := time.Unix(int64(o.Timestamp), 0)

	or := &OrderRecord{
		ID:              bson.ObjectId(o.ID),
		PairName:        o.PairName,
		ExchangeAddress: common.BytesToAddress(o.ExchangeAddress).Hex(),
		UserAddress:     common.BytesToAddress(o.UserAddress).Hex(),
		BuyToken:        common.BytesToAddress(o.BuyToken).Hex(),
		SellToken:       common.BytesToAddress(o.SellToken).Hex(),
		BaseToken:       common.BytesToAddress(o.BaseToken).Hex(),
		QuoteToken:      common.BytesToAddress(o.QuoteToken).Hex(),
		BuyAmount:       o.BuyAmount.String(),
		SellAmount:      o.SellAmount.String(),
		Status:          o.Status,
		Side:            o.Side,
		Hash:            common.BytesToHash(o.Hash).Hex(),
		Nonce:           o.Nonce.String(),
		Expires:         o.Expires.String(),
		MakeFee:         o.MakeFee.String(),
		TakeFee:         o.TakeFee.String(),
		CreatedAt:       timestamp,
		UpdatedAt:       timestamp,
	}

	if o.PricePoint != nil {
		or.PricePoint = o.PricePoint.String()
	}

	if o.Amount != nil {
		or.Amount = o.Amount.String()
	}

	if o.FilledAmount != nil {
		or.FilledAmount = o.FilledAmount.String()
	}

	if o.Signature != nil {
		signature, err := NewSignature(o.Signature)
		if err != nil {
			return or, err
		}
		or.Signature = signature.GetRecord()
	}

	return or, nil
}

// OrderRecord is the object that will be saved in the database
type OrderRecord struct {
	ID              bson.ObjectId    `json:"id" bson:"_id"`
	UserAddress     string           `json:"userAddress" bson:"userAddress"`
	ExchangeAddress string           `json:"exchangeAddress" bson:"exchangeAddress"`
	BuyToken        string           `json:"buyToken" bson:"buyToken"`
	SellToken       string           `json:"sellToken" bson:"sellToken"`
	BaseToken       string           `json:"baseToken" bson:"baseToken"`
	QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
	BuyAmount       string           `json:"buyAmount" bson:"buyAmount"`
	SellAmount      string           `json:"sellAmount" bson:"sellAmount"`
	Status          string           `json:"status" bson:"status"`
	Side            string           `json:"side" bson:"side"`
	Hash            string           `json:"hash" bson:"hash"`
	PricePoint      string           `json:"pricepoint" bson:"pricepoint"`
	Amount          string           `json:"amount" bson:"amount"`
	FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
	Nonce           string           `json:"nonce" bson:"nonce"`
	Expires         string           `json:"expires" bson:"expires"`
	MakeFee         string           `json:"makeFee" bson:"makeFee"`
	TakeFee         string           `json:"takeFee" bson:"takeFee"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`

	PairName  string    `json:"pairName" bson:"pairName"`
	CreatedAt time.Time `json:"createdAt" bson:"createdAt"`
	UpdatedAt time.Time `json:"updatedAt" bson:"updatedAt"`
}

func (o *Order) GetBSON() (*OrderRecord, error) {
	or := &OrderRecord{
		ID:              o.ID,
		PairName:        o.PairName,
		ExchangeAddress: o.ExchangeAddress.Hex(),
		UserAddress:     o.UserAddress.Hex(),
		BuyToken:        o.BuyToken.Hex(),
		SellToken:       o.SellToken.Hex(),
		BaseToken:       o.BaseToken.Hex(),
		QuoteToken:      o.QuoteToken.Hex(),
		BuyAmount:       o.BuyAmount.String(),
		SellAmount:      o.SellAmount.String(),
		Status:          o.Status,
		Side:            o.Side,
		Hash:            o.Hash.Hex(),
		Nonce:           o.Nonce.String(),
		Expires:         o.Expires.String(),
		MakeFee:         o.MakeFee.String(),
		TakeFee:         o.TakeFee.String(),
		CreatedAt:       o.CreatedAt,
		UpdatedAt:       o.UpdatedAt,
	}

	if o.PricePoint != nil {
		or.PricePoint = o.PricePoint.String()
	}

	if o.Amount != nil {
		or.Amount = o.Amount.String()
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
		BuyToken        string           `json:"buyToken" bson:"buyToken"`
		SellToken       string           `json:"sellToken" bson:"sellToken"`
		BaseToken       string           `json:"baseToken" bson:"baseToken"`
		QuoteToken      string           `json:"quoteToken" bson:"quoteToken"`
		BuyAmount       string           `json:"buyAmount" bson:"buyAmount"`
		SellAmount      string           `json:"sellAmount" bson:"sellAmount"`
		Status          string           `json:"status" bson:"status"`
		Side            string           `json:"side" bson:"side"`
		Hash            string           `json:"hash" bson:"hash"`
		PricePoint      string           `json:"pricepoint" bson:"pricepoint"`
		Amount          string           `json:"amount" bson:"amount"`
		FilledAmount    string           `json:"filledAmount" bson:"filledAmount"`
		Nonce           string           `json:"nonce" bson:"nonce"`
		Expires         string           `json:"expires" bson:"expires"`
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
	o.BuyToken = common.HexToAddress(decoded.BuyToken)
	o.SellToken = common.HexToAddress(decoded.SellToken)
	o.BaseToken = common.HexToAddress(decoded.BaseToken)
	o.QuoteToken = common.HexToAddress(decoded.QuoteToken)

	o.BuyAmount = math.ToBigInt(decoded.BuyAmount)
	o.SellAmount = math.ToBigInt(decoded.SellAmount)
	o.FilledAmount = math.ToBigInt(decoded.FilledAmount)

	o.Nonce = math.ToBigInt(decoded.Nonce)
	o.Expires = math.ToBigInt(decoded.Expires)
	o.MakeFee = math.ToBigInt(decoded.MakeFee)
	o.TakeFee = math.ToBigInt(decoded.TakeFee)
	o.Status = decoded.Status
	o.Side = decoded.Side
	o.Hash = common.HexToHash(decoded.Hash)

	if decoded.Amount != "" {
		o.Amount = math.ToBigInt(decoded.Amount)
	}

	if decoded.FilledAmount != "" {
		o.FilledAmount = math.ToBigInt(decoded.FilledAmount)
	}

	if decoded.PricePoint != "" {
		o.PricePoint = math.ToBigInt(decoded.PricePoint)
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

func (o *Order) Print() {
	b, err := json.MarshalIndent(o, "", "  ")
	if err != nil {
		fmt.Println("Error: ", err)
	}

	fmt.Print("\n", string(b))
}
