package units

import (
	"math/big"

	"github.com/tomochain/tomoxsdk/utils/math"
)

func Ethers(value int64) *big.Int {
	return math.Mul(big.NewInt(1e18), big.NewInt(value))
}
