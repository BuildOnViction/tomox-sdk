package ethereum

import (
	"math/big"

	"github.com/tomochain/backend-matching-engine/swap/config"
)

// convert string value from txEnvelope to Native Token
func (t Transaction) ValueToTomo() string {
	valueEth := new(big.Rat)
	valueEth.Quo(new(big.Rat).SetInt(t.ValueWei), weiInEth)
	return valueEth.FloatString(config.TomoAmountPrecision)
}
