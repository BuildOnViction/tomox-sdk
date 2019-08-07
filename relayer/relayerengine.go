package relayer

import (
	"context"
	"log"
	"os"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/tomochain/tomox-sdk/utils"
)

// Blockchain struct
type Blockchain struct {
	client    *rpc.Client
	ethclient *ethclient.Client
	signer    *Signer
}

// PairToken pare token
type PairToken struct {
	BaseToken  common.Address
	QuoteToken common.Address
}

// TokenInfo token info
type TokenInfo struct {
	Name     string
	Symbol   string
	Decimals uint8
}

// RInfo struct
type RInfo struct {
	Tokens  map[common.Address]*TokenInfo
	Pairs   []*PairToken
	MakeFee uint16
	TakeFee uint16
}

// NewBlockchain init
func NewBlockchain(client *rpc.Client,
	ethclient *ethclient.Client,
	signer *Signer) *Blockchain {

	return &Blockchain{
		client:    client,
		ethclient: ethclient,
		signer:    signer,
	}
}

func (b *Blockchain) abiFrom(abiPath string) (*abi.ABI, error) {
	file, err := os.Open(abiPath)
	if err != nil {
		return nil, err
	}
	parsed, err := abi.JSON(file)
	if err != nil {
		return nil, err
	}
	return &parsed, nil
}

// RunContract run smart contract
func (b *Blockchain) RunContract(contractAddr common.Address, abi *abi.ABI, method string, args ...interface{}) (interface{}, error) {
	input, err := abi.Pack(method)
	if err != nil {
		return nil, err
	}
	from := b.signer.GetAddress()
	msg := ether.CallMsg{From: from, To: &contractAddr, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Println(err)
	}
	var unpackResult interface{}
	err = abi.Unpack(&unpackResult, method, result)
	if err != nil {
		return nil, err
	}
	return unpackResult, nil
}

// GetTokenInfoEx return token info
func (b *Blockchain) GetTokenInfoEx(token common.Address, abiPath string) (*TokenInfo, error) {
	abi, err := b.abiFrom(abiPath)
	if err != nil {
		return nil, err
	}
	return b.GetTokenInfo(token, abi)
}

// GetTokenInfo return token info
func (b *Blockchain) GetTokenInfo(token common.Address, abi *abi.ABI) (*TokenInfo, error) {

	result, err := b.RunContract(token, abi, "name")
	if err != nil {
		return nil, err
	}
	name := result.(string)
	result, err = b.RunContract(token, abi, "symbol")
	if err != nil {
		return nil, err
	}
	symbol := result.(string)
	result, err = b.RunContract(token, abi, "decimals")
	if err != nil {
		return nil, err
	}
	decimals := result.(uint8)

	return &TokenInfo{
		Name:     name,
		Symbol:   symbol,
		Decimals: decimals,
	}, nil
}
func (b *Blockchain) isBaseToken(info *TokenInfo) bool {
	if info.Symbol == "TOMO" {
		return true
	}
	return false
}
func (b *Blockchain) setBaseAddress() common.Address {
	return common.HexToAddress("0x0000000000000000000000000000000000000001")
}

// GetRelayer return all tokens in smart contract
func (b *Blockchain) GetRelayer(coinAddress common.Address, abiFile string, tokenAbiFile string, contractAddress common.Address) (*RInfo, error) {
	abi, err := b.abiFrom(abiFile)
	if err != nil {
		return nil, err
	}
	input, err := abi.Pack("getRelayerByCoinbase", coinAddress)
	if err != nil {
		return nil, err
	}

	from := b.signer.GetAddress()
	msg := ether.CallMsg{From: from, To: &contractAddress, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		log.Println(err)
	}
	log.Println("data: ", result)

	relayerInfo := RInfo{
		Tokens: make(map[common.Address]*TokenInfo),
	}
	if method, ok := abi.Methods["getRelayerByCoinbase"]; ok {
		if len(result)%32 != 0 {

		}
		contractData, err := method.Outputs.UnpackValues(result)
		if err == nil {
			if len(contractData) == 5 {
				relayerInfo.MakeFee = contractData[2].(uint16)
				relayerInfo.TakeFee = contractData[2].(uint16)
				fromTokens := contractData[3].([]common.Address)
				toTokens := contractData[4].([]common.Address)
				setToken := utils.Union(fromTokens, toTokens)
				for _, t := range setToken {
					tokenInfo, err := b.GetTokenInfoEx(t, tokenAbiFile)
					if err != nil {
						return nil, err
					}
					relayerInfo.Tokens[t] = tokenInfo
				}
				if len(fromTokens) == len(toTokens) {
					for i, v := range fromTokens {
						base := v
						quote := toTokens[i]
						if b.isBaseToken(relayerInfo.Tokens[base]) {
							base = b.setBaseAddress()
						}
						if b.isBaseToken(relayerInfo.Tokens[quote]) {
							quote = b.setBaseAddress()
						}
						pairToken := &PairToken{
							BaseToken:  base,
							QuoteToken: quote,
						}
						relayerInfo.Pairs = append(relayerInfo.Pairs, pairToken)
					}
				}
				for k, v := range relayerInfo.Tokens {
					if b.isBaseToken(v) {
						relayerInfo.Tokens[b.setBaseAddress()] = v
						delete(relayerInfo.Tokens, k)
						break
					}
				}
			}
		}
	}

	return &relayerInfo, nil
}
