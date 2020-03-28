package services

import (
	"math/big"

	"github.com/tomochain/tomox-sdk/relayer"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// RelayerService struct
type RelayerService struct {
	relayer           interfaces.Relayer
	tokenDao          interfaces.TokenDao
	colateralTokenDao interfaces.TokenDao
	lendingTokenDao   interfaces.TokenDao
	pairDao           interfaces.PairDao
	lendingPairDao    interfaces.LendingPairDao
	relayerDao        interfaces.RelayerDao
}

// NewRelayerService returns a new instance of orderservice
func NewRelayerService(
	relaye interfaces.Relayer,
	tokenDao interfaces.TokenDao,
	colateralTokenDao interfaces.TokenDao,
	lendingTokenDao interfaces.TokenDao,
	pairDao interfaces.PairDao,
	lendingPairDao interfaces.LendingPairDao,
	relayerDao interfaces.RelayerDao,
) *RelayerService {
	return &RelayerService{
		relaye,
		tokenDao,
		colateralTokenDao,
		lendingTokenDao,
		pairDao,
		lendingPairDao,
		relayerDao,
	}
}

func (s *RelayerService) updatePairRelayer(relayerInfo *relayer.RInfo) error {
	currentPairs, err := s.pairDao.GetAllByCoinbase(relayerInfo.Address)
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
				RelayerAddress:     relayerInfo.Address,
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
			err := s.pairDao.DeleteByTokenAndCoinbase(currentPair.BaseTokenAddress, currentPair.QuoteTokenAddress, relayerInfo.Address)
			if err == nil {
				logger.Error(err)
			}
		}
	}
	return nil
}

func (s *RelayerService) updateLendingPair(relayerInfo *relayer.LendingRInfo) error {
	currentPairs, err := s.lendingPairDao.GetAllByCoinbase(relayerInfo.Address)
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
			lendingTokenData := relayerInfo.LendingTokens[newpair.LendingToken]
			pair := &types.LendingPair{
				Term:                 newpair.Term,
				LendingTokenAddress:  newpair.LendingToken,
				LendingTokenDecimals: int(lendingTokenData.Decimals),
				LendingTokenSymbol:   lendingTokenData.Symbol,
				RelayerAddress:       relayerInfo.Address,
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

func (s *RelayerService) updateRelayers(relayerInfos []*relayer.RInfo) error {
	currentRelayers, err := s.relayerDao.GetAll()
	if err != nil {
		return err
	}

	found := false
	for _, r := range relayerInfos {
		found = false
		for _, v := range currentRelayers {
			if v.Address.Hex() == r.Address.Hex() {
				found = true
				break
			}
		}
		if !found {
			domain := r.Address.Hex() + ".devnet.tomochain.com"
			relayer := &types.Relayer{
				Domain:  domain,
				Address: r.Address,
			}
			logger.Info("Create relayer:", r.Address.Hex())
			err = s.relayerDao.Create(relayer)
			if err != nil {
				logger.Error(err)
			}
		}
	}

	for _, r := range currentRelayers {
		found = false
		for _, v := range relayerInfos {
			if v.Address.Hex() == r.Address.Hex() {
				found = true
				break
			}
		}
		if !found {
			logger.Info("Delete relayer:", r.Address.Hex)
			err = s.relayerDao.DeleteByAddress(r.Address)
			if err != nil {
				logger.Error(err)
			}
		}
	}
	return nil
}

func (s *RelayerService) updateTokenRelayer(relayerInfo *relayer.RInfo) error {
	currentTokens, err := s.tokenDao.GetAllByCoinbase(relayerInfo.Address)
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
			RelayerAddress:  relayerInfo.Address,
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
			err = s.tokenDao.UpdateByTokenAndCoinbase(ntoken, relayerInfo.Address, token)
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
				err = s.tokenDao.DeleteByTokenAndCoinbase(ctoken.ContractAddress, relayerInfo.Address)
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
	return nil
}

func (s *RelayerService) updateCollateralTokenRelayer(relayerInfo *relayer.LendingRInfo) error {
	currentTokens, err := s.colateralTokenDao.GetAllByCoinbase(relayerInfo.Address)
	if err != nil {
		return err
	}

	for ntoken, v := range relayerInfo.ColateralTokens {
		found := false
		for _, ctoken := range currentTokens {
			if ntoken.Hex() == ctoken.ContractAddress.Hex() {
				found = true
			}
		}
		token := &types.Token{
			Symbol:          v.Symbol,
			ContractAddress: ntoken,
			RelayerAddress:  relayerInfo.Address,
			Decimals:        int(v.Decimals),
			MakeFee:         big.NewInt(int64(relayerInfo.Fee)),
			TakeFee:         big.NewInt(int64(relayerInfo.Fee)),
		}
		if !found {
			logger.Info("Create collateral token:", token.ContractAddress.Hex())
			err = s.colateralTokenDao.Create(token)
			if err != nil {
				logger.Error(err)
			}
		} else {
			logger.Info("Update collateral token:", token.ContractAddress.Hex())
			err = s.colateralTokenDao.UpdateByTokenAndCoinbase(ntoken, relayerInfo.Address, token)
		}
		for _, ctoken := range currentTokens {
			found = false
			for ntoken, v = range relayerInfo.ColateralTokens {

				if ctoken.ContractAddress.Hex() == ntoken.Hex() {
					found = true
				}
			}
			if !found {
				logger.Info("Delete collateral token:", ctoken.ContractAddress.Hex)
				err = s.colateralTokenDao.DeleteByTokenAndCoinbase(ctoken.ContractAddress, relayerInfo.Address)
				if err != nil {
					logger.Error(err)
				}
			}
		}
	}
	return nil
}

func (s *RelayerService) updateLendingTokenRelayer(relayerInfo *relayer.LendingRInfo) error {
	currentTokens, err := s.lendingTokenDao.GetAllByCoinbase(relayerInfo.Address)
	if err != nil {
		return err
	}

	for ntoken, v := range relayerInfo.LendingTokens {
		found := false
		for _, ctoken := range currentTokens {
			if ntoken.Hex() == ctoken.ContractAddress.Hex() {
				found = true
			}
		}
		token := &types.Token{
			Symbol:          v.Symbol,
			ContractAddress: ntoken,
			RelayerAddress:  relayerInfo.Address,
			Decimals:        int(v.Decimals),
			MakeFee:         big.NewInt(int64(relayerInfo.Fee)),
			TakeFee:         big.NewInt(int64(relayerInfo.Fee)),
		}
		if !found {
			logger.Info("Create lending token:", token.ContractAddress.Hex())
			err = s.lendingTokenDao.Create(token)
			if err != nil {
				logger.Error(err)
			}
		} else {
			logger.Info("Update lending token:", token.ContractAddress.Hex())
			err = s.lendingTokenDao.UpdateByTokenAndCoinbase(ntoken, relayerInfo.Address, token)
		}
		for _, ctoken := range currentTokens {
			found = false
			for ntoken, v = range relayerInfo.LendingTokens {

				if ctoken.ContractAddress.Hex() == ntoken.Hex() {
					found = true
				}
			}
			if !found {
				logger.Info("Delete lending token:", ctoken.ContractAddress.Hex)
				err = s.lendingTokenDao.DeleteByTokenAndCoinbase(ctoken.ContractAddress, relayerInfo.Address)
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
	s.updateCollateralTokenRelayer(relayerLendingInfo)
	s.updateLendingTokenRelayer(relayerLendingInfo)
	return nil
}

func (s *RelayerService) UpdateRelayers() error {
	relayerInfos, err := s.relayer.GetRelayers()
	if err != nil {
		return err
	}
	s.updateRelayers(relayerInfos)
	for _, relayerInfo := range relayerInfos {
		s.updateTokenRelayer(relayerInfo)
		s.updatePairRelayer(relayerInfo)
	}

	relayerLendingInfos, err := s.relayer.GetLendings()
	if err != nil {
		return err
	}
	for _, relayerLendingInfo := range relayerLendingInfos {
		s.updateLendingPair(relayerLendingInfo)
		s.updateCollateralTokenRelayer(relayerLendingInfo)
		s.updateLendingTokenRelayer(relayerLendingInfo)
	}
	return nil
}
