package crons

import (
	"github.com/robfig/cron"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) getFiatPriceCron(c *cron.Cron) {
	c.AddFunc("*/5 * * * * *", s.getFiatPrice())
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getFiatPrice() func() {
	return func() {
		s.PriceBoardService.SyncFiatPrice()
	}
}
