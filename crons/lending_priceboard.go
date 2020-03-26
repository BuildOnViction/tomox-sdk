package crons

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/robfig/cron"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) startLendingPriceBoardCron(c *cron.Cron) {
	pairs, err := s.lendingPairService.GetAll()

	if err != nil {
		log.Println(err.Error())
	}

	for _, p := range pairs {
		c.AddFunc("*/3 * * * * *", s.getLendingPriceBoardData(p.Term, p.LendingTokenAddress))
	}
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getLendingPriceBoardData(term uint64, lendingToken common.Address) func() {
	return func() {
		id := utils.GetLendingChannelID(term, lendingToken)
		tick := s.lendingPriceBoardService.GetLendingPriceBoardData(term, lendingToken)
		ws.GetLendingPriceBoardSocket().BroadcastMessage(id, tick)
	}
}
