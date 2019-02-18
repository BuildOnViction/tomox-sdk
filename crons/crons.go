package crons

import (
	"github.com/robfig/cron"
	"github.com/tomochain/dex-server/services"
)

// CronService contains the services required to initialize crons
type CronService struct {
	ohlcvService *services.OHLCVService
	orderService *services.OrderService
}

// NewCronService returns a new instance of CronService
func NewCronService(
	ohlcvService *services.OHLCVService,
	orderService *services.OrderService,
) *CronService {
	return &CronService{
		ohlcvService: ohlcvService,
		orderService: orderService,
	}
}

// InitCrons is responsible for initializing all the crons in the system
func (s *CronService) InitCrons() {
	c := cron.New()
	s.tickStreamingCron(c)
	s.syncOrdersCron(c)
	c.Start()
}
