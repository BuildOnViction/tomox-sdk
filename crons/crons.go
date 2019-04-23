package crons

import (
	"github.com/robfig/cron"
	"github.com/tomochain/tomodex/engine"
	"github.com/tomochain/tomodex/services"
)

// CronService contains the services required to initialize crons
type CronService struct {
	OHLCVService      *services.OHLCVService
	PriceBoardService *services.PriceBoardService
	PairService       *services.PairService
	Engine            *engine.Engine
}

// NewCronService returns a new instance of CronService
func NewCronService(
	ohlcvService *services.OHLCVService,
	priceBoardService *services.PriceBoardService,
	pairService *services.PairService,
	engine *engine.Engine,
) *CronService {
	return &CronService{
		OHLCVService:      ohlcvService,
		PriceBoardService: priceBoardService,
		PairService:       pairService,
		Engine:            engine,
	}
}

// InitCrons is responsible for initializing all the crons in the system
func (s *CronService) InitCrons() {
	c := cron.New()
	s.tickStreamingCron(c)   // Cron to fetch OHLCV data
	s.getFiatPriceCron(c)    // Cron to query USD price from coinmarketcap.com and update "tokens" collection
	s.startPriceBoardCron(c) // Cron to fetch data for top price board
	s.startMarketsCron(c)    // Cron to fetch markets data
	c.Start()
}
