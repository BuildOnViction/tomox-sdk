package services

import (
	"math/big"

	"github.com/tomochain/tomox-sdk/relayer"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// RelayerService struct
type RelayerService struct {
	relayer        interfaces.Relayer
	tokenDao       interfaces.TokenDao
	pairDao        interfaces.PairDao
	lendingPairDao interfaces.LendingPairDao
}

// NewRelayerService returns a new instance of orderservice
func NewRelayerService(
	relaye interfaces.Relayer,
	tokenDao interfaces.TokenDao,
	pairDao interfaces.PairDao,
	lendingPairDao interfaces.LendingPairDao,

) *RelayerService {
	return &RelayerService{
		relaye,
		tokenDao,
		pairDao,
		lendingPairDao,
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

func (s *RelayerService) updateLendingPair(relayerInfo *relayer.LendingRInfo) error {
	currentPairs, err := s.lendingPairDao.GetAll()
	logger.Info("UpdatePairRelayer starting...")
	if err != nil {
		return err
	}

	for _, newpair := range relayerInfo.LendingPairs {
		found := false
		for _, currentPair := range currentPairs {
			if newpair.Term == currentPair.Term && newpair.LendingToken == currentPair.LendingTokenAddress {
				found = true
				break
			}
		}
		if !found {
			lendingTokenData := relayerInfo.Tokens[newpair.LendingToken]
			pair := &types.LendingPair{
				Term:                 newpair.Term,
				LendingTokenAddress:  newpair.LendingToken,
				LendingTokenDecimals: int(lendingTokenData.Decimals),
				LendingTokenSymbol:   lendingTokenData.Symbol,
			}
			logger.Info("Create Pair:", pair.Term, pair.LendingTokenAddress.Hex())
			err := s.lendingPairDao.Create(pair)
			if err != nil {
				return err
			}
		}
	}

	for _, currentPair := range currentPairs {
		found := false
		for _, newpair := range relayerInfo.LendingPairs {
			if currentPair.Term == newpair.Term && currentPair.LendingTokenAddress == newpair.LendingToken {
				found = true
			}
		}
		if !found {
			logger.Info("Delete Pair:", currentPair.Term, currentPair.LendingTokenAddress.Hex())
			err := s.lendingPairDao.DeleteByLendingKey(currentPair.Term, currentPair.LendingTokenAddress)
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
		token := &types.Token{
			Symbol:          v.Symbol,
			ContractAddress: ntoken,
			Decimals:        int(v.Decimals),
			MakeFee:         big.NewInt(int64(relayerInfo.MakeFee)),
			TakeFee:         big.NewInt(int64(relayerInfo.TakeFee)),
		}
		if !found {
			logger.Info("Create Token:", token.ContractAddress.Hex())
			err = s.tokenDao.Create(token)
			if err != nil {
				logger.Error(err)
			}
		} else {
			logger.Info("Update Token:", token.ContractAddress.Hex())
			err = s.tokenDao.UpdateByToken(ntoken, token)
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

	relayerLendingInfo, err := s.relayer.GetLending()
	if err != nil {
		return err
	}
	s.updateLendingPair(relayerLendingInfo)

	return nil
}
