package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/crypto/secp256k1"
	"github.com/tomochain/backend-matching-engine/swap/errors"
	"github.com/tyler-smith/go-bip32"
)

// NewAddressGenerator : generate new address from master key : cfg.Ethereum.MasterPublicKey
func NewAddressGenerator(masterPublicKeyString string) (*AddressGenerator, error) {
	deserializedMasterPublicKey, err := bip32.B58Deserialize(masterPublicKeyString)
	if err != nil {
		return nil, errors.Wrap(err, "Error deserializing master public key")
	}

	if deserializedMasterPublicKey.IsPrivate {
		return nil, errors.New("Key is not a master public key")
	}

	return &AddressGenerator{deserializedMasterPublicKey}, nil
}

func (g *AddressGenerator) Generate(index uint64) (*common.Address, error) {
	if g.masterPublicKey == nil {
		return nil, errors.New("No master public key set")
	}

	accountKey, err := g.masterPublicKey.NewChildKey(uint32(index))
	if err != nil {
		return nil, errors.Wrap(err, "Error creating new child key")
	}

	x, y := secp256k1.DecompressPubkey(accountKey.Key)

	uncompressed := make([]byte, 64)
	copy(uncompressed[0:32], x.Bytes())
	copy(uncompressed[32:], y.Bytes())

	keccak := crypto.Keccak256(uncompressed)
	address := common.BytesToAddress(keccak[12:]) // Encode lower 160 bits/20 bytes
	return &address, nil
}
