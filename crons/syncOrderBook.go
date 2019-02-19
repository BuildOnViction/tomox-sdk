package crons

import (
	"github.com/robfig/cron"
)

// syncOrderBookCron will fetch new orders from TomoX RPC API periodically
func (s *CronService) syncOrderBookCron(c *cron.Cron) {
	c.AddFunc("*/2 * * * * *", s.syncOrderBook())
}

func (s *CronService) syncOrderBook() func() {
	return func() {
		s.orderService.SyncOrderBook()
	}
}
