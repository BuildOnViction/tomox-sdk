package relayer

import (
	"context"
	"errors"
	"math/big"
	"os"
	"strconv"

	ether "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/rpc"
	relayerAbi "github.com/tomochain/tomox-sdk/relayer/abi"
	"github.com/tomochain/tomox-sdk/utils"
)

var logger = utils.Logger

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

// LendingPairToken lending pari
type LendingPairToken struct {
	Term         uint64
	LendingToken common.Address
}

// TokenInfo token info
type TokenInfo struct {
	Name     string
	Symbol   string
	Decimals uint8
	address  common.Address
}

// RInfo struct
type RInfo struct {
	Address common.Address
	Tokens  map[common.Address]*TokenInfo
	Pairs   []*PairToken
	MakeFee uint16
	TakeFee uint16
}

// LendingRInfo lending relayer info
type LendingRInfo struct {
	Address         common.Address
	ColateralTokens map[common.Address]*TokenInfo
	LendingTokens   map[common.Address]*TokenInfo
	LendingPairs    []*LendingPairToken
	Fee             uint16
}
type Corrateral struct {
	Name    string         `json:"name"`
	Address common.Address `json:"address"`
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

	msg := ether.CallMsg{To: &contractAddr, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Error(err)
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

func (b *Blockchain) setBaseTokenInfo() *TokenInfo {
	return &TokenInfo{
		Name:     "TOMO",
		Symbol:   "TOMO",
		Decimals: 18,
	}
}

func (b *Blockchain) GetRelayers(contractAddress common.Address) ([]*RInfo, error) {
	count, _ := b.GetRelayerCount(contractAddress)

	var rInfos []*RInfo
	logger.Debug("Relayer count", count.String())
	for i := int64(0); i < int64(count.Uint64()); i++ {
		coinbase, _ := b.GetRelayerCoinBaseByIndex(i, contractAddress)
		rInfo, _ := b.GetRelayer(coinbase, contractAddress)
		rInfos = append(rInfos, rInfo)
	}

	return rInfos, nil
}

func (b *Blockchain) GetRelayerCoinBaseByIndex(idx int64, contractAddress common.Address) (common.Address, error) {
	abiRelayer, err := relayerAbi.GetRelayerAbi()
	if err != nil {
		return common.Address{}, err
	}
	input, err := abiRelayer.Pack("RELAYER_COINBASES", big.NewInt(idx))
	if err != nil {
		return common.Address{}, err
	}

	msg := ether.CallMsg{To: &contractAddress, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Error(err)
		return common.Address{}, err
	}

	if method, ok := abiRelayer.Methods["RELAYER_COINBASES"]; ok {
		contractData, _ := method.Outputs.UnpackValues(result)
		return contractData[0].(common.Address), nil
	} else {
		return common.Address{}, errors.New("Can not get coinbase")
	}

	return common.Address{}, nil
}

func (b *Blockchain) GetRelayerCount(contractAddress common.Address) (*big.Int, error) {
	abiRelayer, err := relayerAbi.GetRelayerAbi()
	if err != nil {
		return nil, err
	}

	input, err := abiRelayer.Pack("RelayerCount")
	if err != nil {
		return nil, err
	}

	msg := ether.CallMsg{To: &contractAddress, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if method, ok := abiRelayer.Methods["RelayerCount"]; ok {
		contractData, _ := method.Outputs.UnpackValues(result)
		return contractData[0].(*big.Int), nil
	} else {
		return nil, errors.New("Can not get relayer information")
	}

	return nil, nil
}

// GetRelayer return all tokens in smart contract
func (b *Blockchain) GetRelayer(coinAddress common.Address, contractAddress common.Address) (*RInfo, error) {
	abiRelayer, err := relayerAbi.GetRelayerAbi()
	if err != nil {
		return nil, err
	}
	abiToken, err := relayerAbi.GetTokenAbi()
	if err != nil {
		return nil, err
	}
	input, err := abiRelayer.Pack("getRelayerByCoinbase", coinAddress)
	if err != nil {
		return nil, err
	}

	msg := ether.CallMsg{To: &contractAddress, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Debug("relayer coinbase:", coinAddress.Hex())

	relayerInfo := RInfo{
		Tokens:  make(map[common.Address]*TokenInfo),
		Address: coinAddress,
	}
	if method, ok := abiRelayer.Methods["getRelayerByCoinbase"]; ok {
		contractData, err := method.Outputs.UnpackValues(result)
		if err == nil {
			if len(contractData) == 6 {
				relayerInfo.MakeFee = contractData[3].(uint16)
				relayerInfo.TakeFee = contractData[3].(uint16)
				fromTokens := contractData[4].([]common.Address)
				toTokens := contractData[5].([]common.Address)
				setToken := utils.Union(fromTokens, toTokens)
				for _, t := range setToken {
					if utils.IsNativeTokenByAddress(t) {
						tokenInfo := b.setBaseTokenInfo()
						relayerInfo.Tokens[t] = tokenInfo
					} else {
						tokenInfo, err := b.GetTokenInfo(t, &abiToken)
						if err != nil {
							return nil, err
						}
						relayerInfo.Tokens[t] = tokenInfo
						logger.Debug("Token data:", tokenInfo.Name, tokenInfo.Symbol)
					}

				}
				if len(fromTokens) == len(toTokens) {
					for i, v := range fromTokens {
						base := v
						quote := toTokens[i]

						pairToken := &PairToken{
							BaseToken:  base,
							QuoteToken: quote,
						}
						relayerInfo.Pairs = append(relayerInfo.Pairs, pairToken)
					}
				}

			}
		}
	} else {
		return &relayerInfo, errors.New("Can not get relayer information")
	}

	return &relayerInfo, nil
}

func (b *Blockchain) GetLendingRelayers(contractAddress common.Address) ([]*LendingRInfo, error) {
	count, _ := b.GetRelayerCount(contractAddress)

	var rLInfos []*LendingRInfo
	logger.Debug("Lending Relayer count", count.String())
	for i := int64(0); i < int64(count.Uint64()); i++ {
		coinbase, _ := b.GetRelayerCoinBaseByIndex(i, contractAddress)
		rLInfo, _ := b.GetLendingRelayer(coinbase, contractAddress)
		rLInfos = append(rLInfos, rLInfo)
	}

	return rLInfos, nil
}

// GetLendingRelayer return all lending pair in smart contract
func (b *Blockchain) GetLendingRelayer(coinAddress common.Address, contractAddress common.Address) (*LendingRInfo, error) {
	logger.Debug("GetLendingRelayer:", coinAddress.Hex(), contractAddress.Hex())
	abiRelayer, err := relayerAbi.GetLendingAbi()
	if err != nil {
		return nil, err
	}

	input, err := abiRelayer.Pack("getLendingRelayerByCoinbase", coinAddress)
	if err != nil {
		return nil, err
	}
	abiToken, err := relayerAbi.GetTokenAbi()
	if err != nil {
		return nil, err
	}
	msg := ether.CallMsg{To: &contractAddress, Data: input}
	result, err := b.ethclient.CallContract(context.Background(), msg, nil)
	if err != nil {
		logger.Error(err)
		return nil, err
	}
	logger.Debug("lending relayer coinbase:", coinAddress.Hex())

	lendingRInfo := LendingRInfo{
		ColateralTokens: make(map[common.Address]*TokenInfo),
		LendingTokens:   make(map[common.Address]*TokenInfo),
		Address:         coinAddress,
	}

	if method, ok := abiRelayer.Methods["getLendingRelayerByCoinbase"]; ok {
		contractData, err := method.Outputs.UnpackValues(result)
		if err == nil {
			if len(contractData) == 4 {
				lendingRInfo.Fee = contractData[0].(uint16)
				termList := contractData[2].([]*big.Int)
				lendingTokenList := contractData[1].([]common.Address)
				setLendingToken := utils.Union(lendingTokenList, lendingTokenList)
				for _, t := range setLendingToken {
					if utils.IsNativeTokenByAddress(t) {
						tokenInfo := b.setBaseTokenInfo()
						lendingRInfo.LendingTokens[t] = tokenInfo
					} else {
						tokenInfo, err := b.GetTokenInfo(t, &abiToken)
						if err != nil {
							return nil, err
						}
						lendingRInfo.LendingTokens[t] = tokenInfo
						logger.Debug("Token data:", tokenInfo.Name, tokenInfo.Symbol)
					}

				}
				if len(termList) == len(lendingTokenList) {
					for i, v := range termList {
						t, err := strconv.ParseUint(v.String(), 10, 64)
						if err != nil {
							return &lendingRInfo, err
						}
						pairToken := &LendingPairToken{
							Term:         t,
							LendingToken: lendingTokenList[i],
						}
						lendingRInfo.LendingPairs = append(lendingRInfo.LendingPairs, pairToken)
					}
				}

			}
		}
	} else {
		return &lendingRInfo, errors.New("Can not get relayer information")
	}

	for i := 0; i < len(lendingRInfo.LendingPairs); i++ {
		input, err = abiRelayer.Pack("COLLATERALS", big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}

		msg = ether.CallMsg{To: &contractAddress, Data: input}
		result, err = b.ethclient.CallContract(context.Background(), msg, nil)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
		var unpackResult interface{}
		err = abiRelayer.Unpack(&unpackResult, "COLLATERALS", result)
		if err == nil {
			t := unpackResult.(common.Address)
			if utils.IsNativeTokenByAddress(t) {
				tokenInfo := b.setBaseTokenInfo()
				lendingRInfo.ColateralTokens[t] = tokenInfo
			} else {
				tokenInfo, err := b.GetTokenInfo(t, &abiToken)
				if err != nil {
					return nil, err
				}
				lendingRInfo.ColateralTokens[t] = tokenInfo
			}
		}

	}
	return &lendingRInfo, nil
}
