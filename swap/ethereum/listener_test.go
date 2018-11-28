package ethereum

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestGetBlock(t *testing.T) {
	client, _ := ethclient.Dial("http://localhost:8545/")
	d := time.Now().Add(5 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	blockNumberInt := big.NewInt(7000)
	block, err := client.BlockByNumber(ctx, blockNumberInt)
	t.Logf("Block info: %#v, err: %v", block.Transactions(), err)
}
