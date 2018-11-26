package types

import (
	"errors"
	"time"
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

const (
	SwapSchemaVersion uint64 = 2
	ChainEthereum     Chain  = "ethereum"
)

type AddressAssociation struct {
	// Chain is the name of the payment origin chain
	Chain Chain `json:"chain"`
	// BIP-44
	AddressIndex       uint32    `json:"address_index"`
	Address            string    `json:"address"`
	TomochainPublicKey string    `json:"tomochain_public_key"`
	CreatedAt          time.Time `json:"created_at"`
}
