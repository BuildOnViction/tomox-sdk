package services

import (
	"fmt"
	"time"

	"github.com/tomochain/tomoxsdk/interfaces"
	"github.com/tomochain/tomoxsdk/types"
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

// InitFiatPrice will query Coingecko API and stores fiat price data in the last 1 day after booting up server
func (s *FiatPriceService) InitFiatPrice() {
	// Fix ids with 4 coins
	ids := []string{"bitcoin", "ethereum", "ripple", "tomochain"}
	// Fix fiat currency with USD
	vsCurrency := "usd"

	for _, id := range ids {
		data, err := s.FiatPriceDao.GetCoinMarketChart(id, vsCurrency, "1")

		if err != nil {
			logger.Error(err)
			continue
		}

		items := data.Prices

		for _, item := range items {
			fiatPriceItem := &types.FiatPriceItem{
				Symbol:       id,
				Timestamp:    fmt.Sprintf("%f", item[0]),
				Price:        fmt.Sprintf("%f", item[1]),
				FiatCurrency: vsCurrency,
			}

			_, err := s.FiatPriceDao.FindAndModify(
				fiatPriceItem.Symbol,
				fiatPriceItem.FiatCurrency,
				fiatPriceItem.Timestamp,
				fiatPriceItem,
			)

			if err != nil {
				logger.Error(err)
				continue
			}
		}
	}
}

// UpdateFiatPrice will query Coingecko API and stores fiat price data in the last 30 minutes
func (s *FiatPriceService) UpdateFiatPrice() {
	// Fix ids with 4 coins
	ids := []string{"bitcoin", "ethereum", "ripple", "tomochain"}
	// Fix fiat currency with USD
	vsCurrency := "usd"

	now := time.Now()
	then := now.Add(time.Duration(-30) * time.Minute)

	for _, id := range ids {
		data, err := s.FiatPriceDao.GetCoinMarketChartRange(id, vsCurrency, then.Unix(), now.Unix())

		if err != nil {
			logger.Error(err)
			continue
		}

		items := data.Prices

		for _, item := range items {
			fiatPriceItem := &types.FiatPriceItem{
				Symbol:       id,
				Timestamp:    fmt.Sprintf("%f", item[0]),
				Price:        fmt.Sprintf("%f", item[1]),
				FiatCurrency: vsCurrency,
			}

			_, err := s.FiatPriceDao.FindAndModify(
				fiatPriceItem.Symbol,
				fiatPriceItem.FiatCurrency,
				fiatPriceItem.Timestamp,
				fiatPriceItem,
			)

			if err != nil {
				logger.Error(err)
				continue
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

func (s *FiatPriceService) GetFiatPriceChart() (map[string][]string, error) {
	// Fix ids with 4 coins
	ids := []string{"bitcoin", "ethereum", "ripple", "tomochain"}

	for _, id := range ids {
		logger.Debug(id)
	}

	return map[string][]string{}, nil
}
