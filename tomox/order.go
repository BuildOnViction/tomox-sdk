package tomox

import (
	"fmt"
	"math/big"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

const (
	OrderStatusOpen          = "OPEN"
	OrderStatusPartialFilled = "PARTIAL_FILLED"
	OrderStatusFilled        = "FILLED"
	OrderStatusCancelled     = "CANCELLED"
)

// Signature struct
type Signature struct {
	V byte
	R common.Hash
	S common.Hash
}

type SignatureRecord struct {
	V byte   `json:"V" bson:"V"`
	R string `json:"R" bson:"R"`
	S string `json:"S" bson:"S"`
}

// OrderItem : info that will be store in database
type OrderItem struct {
	Quantity        *big.Int       `json:"quantity,omitempty"`
	Price           *big.Int       `json:"price,omitempty"`
	ExchangeAddress common.Address `json:"exchangeAddress,omitempty"`
	UserAddress     common.Address `json:"userAddress,omitempty"`
	BaseToken       common.Address `json:"baseToken,omitempty"`
	QuoteToken      common.Address `json:"quoteToken,omitempty"`
	Status          string         `json:"status,omitempty"`
	Side            string         `json:"side,omitempty"`
	Type            string         `json:"type,omitempty"`
	Hash            common.Hash    `json:"hash,omitempty"`
	Signature       *Signature     `json:"signature,omitempty"`
	FilledAmount    *big.Int       `json:"filledAmount,omitempty"`
	Nonce           *big.Int       `json:"nonce,omitempty"`
	PairName        string         `json:"pairName,omitempty"`
	CreatedAt       time.Time      `json:"createdAt,omitempty"`
	UpdatedAt       time.Time      `json:"updatedAt,omitempty"`
	OrderID         uint64         `json:"orderID,omitempty"`
	// *OrderMeta
	NextOrder []byte `json:"-"`
	PrevOrder []byte `json:"-"`
	OrderList []byte `json:"-"`
	Key       string `json:"key"`
}

type OrderItemBSON struct {
	Quantity        string           `json:"quantity,omitempty" bson:"quantity"`
	Price           string           `json:"price,omitempty" bson:"price"`
	ExchangeAddress string           `json:"exchangeAddress,omitempty" bson:"exchangeAddress"`
	UserAddress     string           `json:"userAddress,omitempty" bson:"userAddress"`
	BaseToken       string           `json:"baseToken,omitempty" bson:"baseToken"`
	QuoteToken      string           `json:"quoteToken,omitempty" bson:"quoteToken"`
	Status          string           `json:"status,omitempty" bson:"status"`
	Side            string           `json:"side,omitempty" bson:"side"`
	Type            string           `json:"type,omitempty" bson:"type"`
	Hash            string           `json:"hash,omitempty" bson:"hash"`
	Signature       *SignatureRecord `json:"signature,omitempty" bson:"signature"`
	FilledAmount    string           `json:"filledAmount,omitempty" bson:"filledAmount"`
	Nonce           string           `json:"nonce,omitempty" bson:"nonce"`
	PairName        string           `json:"pairName,omitempty" bson:"pairName"`
	CreatedAt       time.Time        `json:"createdAt,omitempty" bson:"createdAt"`
	UpdatedAt       time.Time        `json:"updatedAt,omitempty" bson:"updatedAt"`
	OrderID         string           `json:"orderID,omitempty" bson:"orderID"`
	NextOrder       string           `json:"nextOrder,omitempty" bson:"nextOrder"`
	PrevOrder       string           `json:"prevOrder,omitempty" bson:"prevOrder"`
	OrderList       string           `json:"orderList,omitempty" bson:"orderList"`
	Key             string           `json:"key" bson:"key"`
}

type Order struct {
	Item *OrderItem
	Key  []byte `json:"orderID"`
}

func (order *Order) String() string {

	return fmt.Sprintf("orderID : %s, price: %s, quantity :%s, relayerID: %s",
		new(big.Int).SetBytes(order.Key), order.Item.Price, order.Item.Quantity, order.Item.ExchangeAddress.Hex())
}
