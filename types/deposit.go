package types

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto/sha3"
	"gopkg.in/mgo.v2/bson"
)

type AssetCode string

const (
	AssetCodeETH AssetCode = "ETH"
)

type Chain string

// Scan implements database/sql.Scanner interface
func (s *Chain) Scan(src interface{}) error {
	value, ok := src.([]byte)
	if !ok {
		return errors.New("Cannot convert value to Chain")
	}
	*s = Chain(value)
	return nil
}

func (s *Chain) String() string {
	return string(*s)
}

func (s *Chain) Bytes() []byte {
	return []byte(*s)
}

const (
	SwapSchemaVersion uint64 = 2
	ChainEthereum     Chain  = "ethereum"
)

type AddressAssociation struct {
	// Chain is the name of the payment origin chain
	Chain Chain `json:"chain"`
	// BIP-44
	AddressIndex       uint64         `json:"addressIndex"`
	Address            common.Address `json:"address"`
	TomochainPublicKey common.Address `json:"tomochainPublicKey"`
	CreatedAt          time.Time      `json:"createdAt"`
}

// AddressAssociationRecord is the object that will be saved in the database
type AddressAssociationRecord struct {
	ID                bson.ObjectId `json:"id" bson:"_id"`
	Chain             string        `json:"chain" bson:"chain"`
	Address           string        `json:"address" bson:"address"`
	TxEnvelope        string        `json:"txEnvelope" bson: "txEnvelope"`
	AssociatedAddress string        `json:"associatedAddress" bson:"associatedAddress"`
	CreatedAt         time.Time     `json:"createdAt" bson:"createdAt"`
	UpdatedAt         time.Time     `json:"updatedAt" bson:"updatedAt"`
}

type AddressAssociationFeed struct {
	Address            []byte `json:"address"`
	AddressIndex       uint64 `json:"addressIndex"`
	Chain              Chain  `json:"chain"`
	CreatedAt          uint64 `json:"createdAt"`
	TomochainPublicKey []byte `json:"tomochainPublicKey"`
}

func (aaf *AddressAssociationFeed) GetJSON() (*AddressAssociation, error) {
	// convert back to JSON object
	timestamp := time.Unix(int64(aaf.CreatedAt), 0)
	aa := &AddressAssociation{
		Chain:              aaf.Chain,
		Address:            common.BytesToAddress(aaf.Address),
		AddressIndex:       aaf.AddressIndex,
		TomochainPublicKey: common.BytesToAddress(aaf.TomochainPublicKey),
		CreatedAt:          timestamp,
	}

	return aa, nil
}

type GenerateAddressResponse struct {
	ProtocolVersion int    `json:"protocolVersion"`
	Chain           string `json:"chain"`
	Address         string `json:"address"`
	Signer          string `json:"signer"`
}

type DepositTransaction struct {
	Chain         Chain
	TransactionID string
	AssetCode     AssetCode
	// CRITICAL REQUIREMENT: Amount in the base unit of currency.
	// For 10 satoshi this should be equal 0.0000001
	// For 1 BTC      this should be equal 1.0000000
	// For 1 Finney   this should be equal 0.0010000
	// For 1 ETH      this should be equal 1.0000000
	// Currently, the length of Amount string shouldn't be longer than 17 characters.
	Amount             string
	TomochainPublicKey string
}

type AssociationTransaction struct {
	Source          string   `json:"source"`
	Signature       []byte   `json:"signature"`
	Hash            []byte   `json:"hash"`
	TransactionType string   `json:"transactionType"`
	Params          []string `json:"params"`
}

// ComputeHash calculates the orderRequest hash
func (o *AssociationTransaction) ComputeHash() []byte {
	sha := sha3.NewKeccak256()

	sha.Write([]byte(o.Source))
	sha.Write([]byte(o.TransactionType))
	for _, param := range o.Params {
		sha.Write([]byte(param))
	}
	return sha.Sum(nil)
}
