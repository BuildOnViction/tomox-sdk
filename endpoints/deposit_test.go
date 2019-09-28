package endpoints

import (
    "math/big"
    "testing"
)

func TestConvert(t *testing.T) {
    transactionAmount := "100000000000000000000"
    amount := new(big.Int)
    amount.SetString(transactionAmount, 10)

    t.Logf("Got amount :%d", amount)

}
