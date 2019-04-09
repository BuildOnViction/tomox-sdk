package bitcoin

import (
	"testing"

	"github.com/btcsuite/btcd/rpcclient"
	// "github.com/ethereum/go-ethereum/ethclient"
)

// Run this: bitcoind -rpcuser=user -rpcpassword=pass -rpcport=18334 -regtest
func TestGetBlock(t *testing.T) {

	connConfig := &rpcclient.ConnConfig{
		Host:         "localhost:18334",
		User:         "user",
		Pass:         "pass",
		HTTPPostMode: true,
		DisableTLS:   true,
	}

	client, err := rpcclient.New(connConfig, nil)
	if err != nil {
		t.Error(err)
	}
	defer client.Shutdown()

	blockCount, err := client.GetBlockCount()
	if err != nil {
		t.Error(err)
	}
	t.Logf("Block count: %d", blockCount)
}
