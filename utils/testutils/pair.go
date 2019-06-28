package testutils

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/types"
)

func GetZRXWETHTestPair() *types.Pair {
	return &types.Pair{
		BaseTokenSymbol:    "ZRX",
		BaseTokenAddress:   common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
		BaseTokenDecimals:  18,
		QuoteTokenSymbol:   "WETH",
		PriceMultiplier:    big.NewInt(1e9),
		QuoteTokenAddress:  common.HexToAddress("0x276e16ada4b107332afd776691a7fbbaede168ef"),
		QuoteTokenDecimals: 18,
	}
}
