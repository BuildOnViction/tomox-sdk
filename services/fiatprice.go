package services

import (
	"fmt"

	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// FiatPriceService functions are responsible for interacting with daos and implements business logics.
type FiatPriceService struct {
	TokenDao       interfaces.TokenDao
	FiatPriceDao   interfaces.FiatPriceDao
	FiatPriceCache interfaces.FiatPriceCache
}

// NewFiatPriceService returns a new instance of TradeService
func NewFiatPriceService(
	tokenDao interfaces.TokenDao,
	fiatPriceDao interfaces.FiatPriceDao,
	fiatPriceCache interfaces.FiatPriceCache,
) *FiatPriceService {
	return &FiatPriceService{
		TokenDao:       tokenDao,
		FiatPriceDao:   fiatPriceDao,
		FiatPriceCache: fiatPriceCache,
	}
}

// InitFiatPrice will query Coingecko API and stores fiat price data in the last 1 day after booting up server
func (s *FiatPriceService) InitFiatPrice() {
	s.UpdateFiatPrice()

}

// UpdateFiatPrice will query Coingecko API and stores fiat price data in the last 30 minutes
func (s *FiatPriceService) UpdateFiatPrice() {
	logger.Info("UpdateFiatPrice starting...")
	// Fix ids with 4 coins
	symbols := []string{"bitcoin", "ethereum", "ripple", "tomochain"}
	// Fix fiat currency with USD
	vsCurrency := "usd"
	for _, symbol := range symbols {
		data, err := s.FiatPriceDao.GetCoinMarketChart(symbol, vsCurrency, "2")

		if err != nil {
			logger.Error(err)
			continue
		}

		items := data.Prices
		totalVolumes := data.TotalVolumes

		for i, item := range items {
			fiatPriceItem := &types.FiatPriceItem{
				Symbol:       symbol,
				Timestamp:    int64(item[0]),
				Price:        fmt.Sprintf("%f", item[1]),
				FiatCurrency: vsCurrency,
				TotalVolume:  fmt.Sprintf("%f", totalVolumes[i][1]),
			}

			_, err := s.FiatPriceDao.FindAndUpdate(
				fiatPriceItem.Symbol,
				fiatPriceItem.FiatCurrency,
				fiatPriceItem.Timestamp,
				fiatPriceItem,
			)

			if err != nil {
				logger.Error(err)
			}
			err = s.FiatPriceCache.AddFiat(symbol, fiatPriceItem)
			if err != nil {
				logger.Error(err)
			}
		}
	}
}

func (s *FiatPriceService) SyncFiatPrice() error {
	prices, err := s.FiatPriceDao.GetLatestQuotes()
	if err != nil {
		logger.Error(err)
		return err
	}

	for k, v := range prices {
		err := s.TokenDao.UpdateFiatPriceBySymbol(k, v)

		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// GetFiatPriceChart return price chart
func (s *FiatPriceService) GetFiatPriceChart() (map[string][]*types.FiatPriceItem, error) {
	result := make(map[string][]*types.FiatPriceItem)

	// Fix ids with 4 coins
	symbols := []string{"bitcoin", "ethereum", "ripple", "tomochain"}

	for _, symbol := range symbols {
		data, err := s.FiatPriceCache.GetFiatRange(symbol, 24)
		if err != nil {
			data, err = s.FiatPriceDao.Get24hChart(symbol, "usd")
			if err != nil {
				logger.Error(err)
			} else {
				result[symbol] = data
			}

		} else {
			result[symbol] = data
		}

	}

	return result, nil
}
