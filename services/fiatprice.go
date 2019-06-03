package services

import (
	"github.com/tomochain/tomoxsdk/interfaces"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type FiatPriceService struct {
	TokenDao     interfaces.TokenDao
	FiatPriceDao interfaces.FiatPriceDao
}

// NewTradeService returns a new instance of TradeService
func NewFiatPriceService(
	tokenDao interfaces.TokenDao,
	fiatPriceDao interfaces.FiatPriceDao,
) *FiatPriceService {
	return &FiatPriceService{
		TokenDao:     tokenDao,
		FiatPriceDao: fiatPriceDao,
	}
}

func (s *FiatPriceService) SyncFiatPrice() {
	prices, err := s.FiatPriceDao.GetLatestQuotes()

	if err != nil {
		logger.Error(err)
		return
	}

	for k, v := range prices {
		err := s.TokenDao.UpdateFiatPriceBySymbol(k, v)

		if err != nil {
			logger.Error(err)
		}
	}
}
