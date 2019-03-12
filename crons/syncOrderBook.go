package crons

import (
	"log"

	"github.com/robfig/cron"
	"github.com/tomochain/dex-server/types"
)

// syncOrderBookCron will fetch new orders from TomoX RPC API periodically
func (s *CronService) syncOrderBookCron(c *cron.Cron) {
	pairs, err := s.PairService.GetAll()

	if err != nil {
		log.Println(err.Error())
	}

	for _, p := range pairs {
		c.AddFunc("*/2 * * * * *", s.syncOrderBook(p))
	}
}

func (s *CronService) syncOrderBook(p types.Pair) func() {
	return func() {
		err := s.Engine.SyncOrderBook(p)

		if err != nil {
			log.Println(err.Error())
		}
	}
}
