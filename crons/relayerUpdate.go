package crons

import (
	"github.com/robfig/cron"
)

func (s *CronService) startRelayerUpdate(c *cron.Cron) {
	s.RelayService.UpdateRelayers()
	c.AddFunc("*/600 * * * * *", s.updateRelayer())
}

func (s *CronService) updateRelayer() func() {
	return func() {
		s.RelayService.UpdateRelayers()
	}
}
