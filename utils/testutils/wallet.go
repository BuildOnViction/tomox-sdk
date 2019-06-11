package testutils

import "github.com/tomochain/tomoxsdk/types"

func GetTestWallet() *types.Wallet {
	return types.NewWalletFromPrivateKey("3411b45169aa5a8312e51357db68621031020dcf46011d7431db1bbb6d3922ce")
	// return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712660")
}

func GetTestWallet1() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712661")
}

func GetTestWallet2() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712662")
}

func GetTestWallet3() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712663")
}

func GetTestWallet4() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712664")
}

func GetTestWallet5() *types.Wallet {
	return types.NewWalletFromPrivateKey("7c78c6e2f65d0d84c44ac0f7b53d6e4dd7a82c35f51b251d387c2a69df712665")
}
