package bitcoin

import (
	"math/big"

	"github.com/tomochain/backend-matching-engine/swap/config"
)

func (t Transaction) ValueToTomo() string {
	valueSat := new(big.Int).SetInt64(t.ValueSat)
	valueBtc := new(big.Rat).Quo(new(big.Rat).SetInt(valueSat), satInBtc)
	return valueBtc.FloatString(config.TomoAmountPrecision)
}
