package endpoints

import (
	"testing"

	"github.com/Proofsuite/amp-matching-engine/app"
	"github.com/Proofsuite/amp-matching-engine/ethereum"
	"github.com/ethereum/go-ethereum/rpc"
)

var (
	rpcClient *rpc.Client
)

func TestGetOrderFromPss(t *testing.T) {
	err := app.LoadConfig("../config", "test")
	if err != nil {
		panic(err)
	}

	provider := ethereum.NewWebsocketProvider()
	var orderResult interface{}
	err = provider.RPCClient.Call(&orderResult, "orderbook_getOrders", "Tomo", "0x28074f8d0fd78629cd59290cac185611a8d60109")
	t.Logf("Order :%+v", orderResult)
	if err != nil {
		t.Error(err)
	}
}
