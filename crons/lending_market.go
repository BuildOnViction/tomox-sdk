package crons

import (
	"github.com/robfig/cron"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) startLendingMarketsCron(c *cron.Cron) {
	c.AddFunc("*/3 * * * * *", s.getLendingMarketsData())
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getLendingMarketsData() func() {
	return func() {
		tick, err := s.lendingOhlcvService.GetAllTokenPairData()
		if err != nil {
			tick = types.LendingTicks{}
		}

		data := &types.LendingMarketData{
			PairData: tick,
		}
		id := utils.GetLendingMarketsChannelID(ws.LendingMarketsChannel)
		ws.GetLendingMarketSocket().BroadcastMessage(id, data)
	}
}
