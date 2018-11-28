package types

import (
	"errors"
	"time"

	"github.com/ethereum/go-ethereum/common"
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
