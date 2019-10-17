package relayer

import (
	"encoding/json"
	"os"
	"path/filepath"

	"io/ioutil"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	ethereum "github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
)

type Signer struct {
	Keystore   string `json:"keystore_path"`
	Passphrase string `json:"passphrase"`
	opts       *bind.TransactOpts
}

func (self *Signer) GetTransactOpts() *bind.TransactOpts {
	return self.opts
}

func (self *Signer) GetAddress() ethereum.Address {
	return self.opts.From
}

func (self *Signer) Sign(tx *types.Transaction) (*types.Transaction, error) {
	return self.opts.Signer(types.HomesteadSigner{}, self.GetAddress(), tx)
}

func NewSignerFile(file string, fileLocation string) *Signer {
	raw, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}
	signer := &Signer{}
	err = json.Unmarshal(raw, signer)
	if err != nil {
		panic(err)
	}

	keyio, err := os.Open(filepath.Join(fileLocation, signer.Keystore))
	if err != nil {
		panic(err)
	}
	logger.Debug("keyio: ", keyio)
	auth, err := bind.NewTransactor(keyio, signer.Passphrase)
	if err != nil {
		panic(err)
	}
	logger.Debug("auth: ", auth.From.Hex())
	signer.opts = auth

	return signer
}

func NewSigner() *Signer {
	signer := &Signer{}
	passParser, keyio := GetKeyStore()
	signer.Passphrase = passParser

	auth, err := bind.NewTransactor(keyio, signer.Passphrase)
	if err != nil {
		panic(err)
	}
	logger.Debug("auth: ", auth.From.Hex())
	signer.opts = auth

	return signer
}
