package crons

import (
	"log"

	"github.com/robfig/cron"
	"github.com/tomochain/dex-server/types"
)

// syncOrderBookCron will fetch new orders from TomoX RPC API periodically
func (s *CronService) syncTradesCron(c *cron.Cron) {
	pairs, err := s.PairService.GetAll()

	if err != nil {
		log.Println(err.Error())
	}

	for _, p := range pairs {
		c.AddFunc("*/2 * * * * *", s.syncTrades(p))
	}
}

func (s *CronService) syncTrades(p types.Pair) func() {
	return func() {
		err := s.Engine.SyncTrades(p)

		if err != nil {
			log.Println(err.Error())
		}
	}
}
