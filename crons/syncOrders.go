package crons

import (
	"github.com/robfig/cron"
)

// syncOrdersCron will fetch new orders from TomoX RPC API periodically
func (s *CronService) syncOrdersCron(c *cron.Cron) {
	c.AddFunc("*/2 * * * * *", s.syncOrders())
}

func (s *CronService) syncOrders() func() {
	return func() {
		s.orderService.SyncOrders()
	}
}
