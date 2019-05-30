package testutils

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomoxsdk/types"
)

func GetTestZRXToken() types.Token {
	return types.Token{
		Name:            "ZRX",
		Symbol:          "ZRX",
		Decimals:        18,
		ContractAddress: common.HexToAddress("0x2034842261b82651885751fc293bba7ba5398156"),
	}
}

func GetTestWETHToken() types.Token {
	return types.Token{
		Name:            "WETH",
		Symbol:          "WETH",
		Decimals:        18,
		ContractAddress: common.HexToAddress("0x276e16ada4b107332afd776691a7fbbaede168ef"),
	}
}
