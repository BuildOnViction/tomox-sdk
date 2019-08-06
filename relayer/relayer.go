package relayer

import (
	"fmt"
	"path/filepath"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
)

const rcpURL string = "https://testnet.tomochain.com/"
const coinBase string = "0xF9D87abd60435b70415CcC1FAAcA4F8B91786eDb"

// Relayer get token
type Relayer struct {
	rpcURL         string
	coinBase       common.Address
	relayerAddress common.Address
}

// NewRelayer init relayer
func NewRelayer(rpcURL string,
	coinBase common.Address,
	relayerAddress common.Address) *Relayer {

	return &Relayer{
		rpcURL:         rpcURL,
		coinBase:       coinBase,
		relayerAddress: relayerAddress,
	}
}

// GetRelayer get relayer information
func (r *Relayer) GetRelayer() (*RInfo, error) {
	fileLocation := "/home/nghiatt/go/src/github.com/tomochain/tomoxsdk/relayer"
	relayerAbi := filepath.Join(fileLocation, "./abi/relayer_registration.abi")
	tokenAbi := filepath.Join(fileLocation, "./abi/token.abi")
	path := filepath.Join(fileLocation, "./signer_config.json")
	signer := NewSigner(path, fileLocation)
	client, err := rpc.Dial(rcpURL)
	if err != nil {
		fmt.Println(err)
	}
	ethclient := ethclient.NewClient(client)
	bc := NewBlockchain(client, ethclient, signer)
	return bc.GetRelayer(r.coinBase, relayerAbi, tokenAbi, r.relayerAddress)

}
