package ethereum

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/ethclient"
)

func TestGetBlock(t *testing.T) {
	client, _ := ethclient.Dial("http://localhost:8545")
	d := time.Now().Add(15 * time.Second)
	ctx, cancel := context.WithDeadline(context.Background(), d)
	defer cancel()

	networkID, err := client.NetworkID(ctx)
	blockNumberInt := big.NewInt(0)
	block, err := client.BlockByNumber(ctx, blockNumberInt)
	t.Logf("NetworkID: %s, Block info: %#v, err: %v", networkID.String(), block.Transactions(), err)
}
