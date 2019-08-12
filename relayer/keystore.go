package relayer

import "strings"

const keyString = `{"address":"008e3986e8aa4b9b5202857cbdac4d15caa0a4f1","crypto":{"cipher":"aes-128-ctr","ciphertext":"b3c4632042ba7137da9168ddd73dc8cca89292a25c5afe773d326f865260efe8","cipherparams":{"iv":"c858a9f0425296ff24b0943e20167da7"},"kdf":"scrypt","kdfparams":{"dklen":32,"n":262144,"p":1,"r":8,"salt":"0df33a43541bde4aa953574d9110825c0da5b095f337e3c2a5023b7ed0f5f63d"},"mac":"97c89141083866402c067bf1e9dc4f7ea3b32ba1d4d050e235fd63735c3626a6"},"id":"5a75ed08-bb03-4ec5-a93d-c9c4ba4752d6","version":3}`
const passParser = "123456"

// GetKeyStoreReader return reader for keystore
func GetKeyStoreReader() *strings.Reader {
	return strings.NewReader(keyString)
}

// GetKeyStore return passparser and keystore reader
func GetKeyStore() (string, *strings.Reader) {
	return passParser, GetKeyStoreReader()
}
