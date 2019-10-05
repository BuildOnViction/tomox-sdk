package services

import (
	"math/big"

	"github.com/tomochain/tomox-sdk/relayer"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// RelayerService struct
type RelayerService struct {
	relayer  interfaces.Relayer
	tokenDao interfaces.TokenDao
	pairDao  interfaces.PairDao
}

// NewRelayerService returns a new instance of orderservice
func NewRelayerService(
	relaye interfaces.Relayer,
	tokenDao interfaces.TokenDao,
	pairDao interfaces.PairDao,
) *RelayerService {
	return &RelayerService{
		relaye,
		tokenDao,
		pairDao,
	}
}

func (s *RelayerService) updatePairRelayer(relayerInfo *relayer.RInfo) error {
	currentPairs, err := s.pairDao.GetAll()
	logger.Info("UpdatePairRelayer starting...")
	if err != nil {
		return err
	}

	for _, newpair := range relayerInfo.Pairs {
		found := false
		for _, currentPair := range currentPairs {
			if newpair.BaseToken == currentPair.BaseTokenAddress && newpair.QuoteToken == currentPair.QuoteTokenAddress {
				found = true
			}
		}
		if !found {
			pairBaseData := relayerInfo.Tokens[newpair.BaseToken]
			pairQuoteData := relayerInfo.Tokens[newpair.QuoteToken]
			pair := &types.Pair{
				BaseTokenSymbol:    pairBaseData.Symbol,
				BaseTokenAddress:   newpair.BaseToken,
				BaseTokenDecimals:  int(pairBaseData.Decimals),
				QuoteTokenSymbol:   pairQuoteData.Symbol,
				QuoteTokenAddress:  newpair.QuoteToken,
				QuoteTokenDecimals: int(pairQuoteData.Decimals),
				Active:             true,
				MakeFee:            big.NewInt(int64(relayerInfo.MakeFee)),
				TakeFee:            big.NewInt(int64(relayerInfo.TakeFee)),
			}
			logger.Info("Create Pair:", pair.BaseTokenAddress.Hex(), pair.QuoteTokenAddress.Hex())
			err := s.pairDao.Create(pair)
			if err != nil {
				return err
			}
		}
	}

	for _, currentPair := range currentPairs {
		found := false
		for _, newpair := range relayerInfo.Pairs {
			if currentPair.BaseTokenAddress == newpair.BaseToken && currentPair.QuoteTokenAddress == newpair.QuoteToken {
				found = true
			}
		}
		if !found {
			logger.Info("Delete Pair:", currentPair.BaseTokenAddress.Hex(), currentPair.QuoteTokenAddress.Hex())
			err := s.pairDao.DeleteByToken(currentPair.BaseTokenAddress, currentPair.QuoteTokenAddress)
			if err == nil {
				logger.Error(err)
			}
		}
	}
	return nil
}
func (s *RelayerService) updateTokenRelayer(relayerInfo *relayer.RInfo) error {
	currentTokens, err := s.tokenDao.GetAll()
	if err != nil {
		return err
	}

	for ntoken, v := range relayerInfo.Tokens {
		found := false
		for _, ctoken := range currentTokens {
			if ntoken.Hex() == ctoken.ContractAddress.Hex() {
				found = true
			}
		}
		if !found {
			token := &types.Token{
				Symbol:          v.Symbol,
				ContractAddress: ntoken,
				Decimals:        int(v.Decimals),
				MakeFee:         big.NewInt(int64(relayerInfo.MakeFee)),
				TakeFee:         big.NewInt(int64(relayerInfo.TakeFee)),
			}
			logger.Info("Create Token:", token.ContractAddress.Hex())
			err = s.tokenDao.Create(token)
			if err != nil {
				logger.Error(err)
			}
		}
		for _, ctoken := range currentTokens {
			found = false
			for ntoken, v = range relayerInfo.Tokens {

				if ctoken.ContractAddress.Hex() == ntoken.Hex() {
					found = true
				}
			}
			if !found {
				logger.Info("Delete Token:", ctoken.ContractAddress.Hex)
				err = s.tokenDao.DeleteByToken(ctoken.ContractAddress)
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
	return nil
}

// UpdateRelayer get the total number of orders amount created by a user
func (s *RelayerService) UpdateRelayer() error {
	relayerInfo, err := s.relayer.GetRelayer()
	if err != nil {
		return err
	}
	s.updateTokenRelayer(relayerInfo)
	s.updatePairRelayer(relayerInfo)
	return nil
}
