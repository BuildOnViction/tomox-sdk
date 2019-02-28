package crons

import (
	"github.com/robfig/cron"
	"github.com/tomochain/dex-server/services"
)

// CronService contains the services required to initialize crons
type CronService struct {
	OHLCVService      *services.OHLCVService
	PriceBoardService *services.PriceBoardService
	PairService       *services.PairService
}

// NewCronService returns a new instance of CronService
func NewCronService(
	ohlcvService *services.OHLCVService,
	priceBoardService *services.PriceBoardService,
	pairService *services.PairService,
) *CronService {
	return &CronService{
		OHLCVService:      ohlcvService,
		PriceBoardService: priceBoardService,
		PairService:       pairService,
	}
}

// InitCrons is responsible for initializing all the crons in the system
func (s *CronService) InitCrons() {
	c := cron.New()
	s.tickStreamingCron(c)
	s.getFiatPriceCron(c)
	s.startPriceBoardCron(c)
	c.Start()
}
