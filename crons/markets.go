package crons

import (
	"log"

	"github.com/robfig/cron"
	"github.com/tomochain/dex-server/utils"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) startMarketsCron(c *cron.Cron) {
	c.AddFunc("*/3 * * * * *", s.getMarketsData())
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getMarketsData() func() {
	return func() {
		res, err := s.PairService.GetAllTokenPairData()

		if err != nil {
			log.Printf("%s", err)
			return
		}

		utils.PrintJSON(res)
	}
}
