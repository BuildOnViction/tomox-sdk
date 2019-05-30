package crons

import (
	"log"
	"time"

	"github.com/robfig/cron"
	"github.com/tomochain/tomoxsdk/types"
	"github.com/tomochain/tomoxsdk/utils"
	"github.com/tomochain/tomoxsdk/ws"
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
		pairData, err := s.PairService.GetAllTokenPairData()

		if err != nil {
			log.Printf("%s", err)
			return
		}

		p := make([]types.PairAddresses, 0)
		duration := int64(1)
		unit := "hour"
		end := int64(time.Now().Unix())
		start := end - 24*60*60 // -1 day
		ticks, err := s.OHLCVService.GetOHLCV(p, duration, unit, start, end)

		tickResult := make(map[string][]*types.Tick)

		for _, tick := range ticks {
			tickResult[tick.Pair.PairName] = append(tickResult[tick.Pair.PairName], &types.Tick{
				Close:     tick.Close,
				Timestamp: tick.Timestamp,
				Pair:      tick.Pair,
			})
		}

		res := &types.MarketData{
			PairData:        pairData,
			SmallChartsData: tickResult,
		}

		id := utils.GetMarketsChannelID(ws.MarketsChannel)

		ws.GetMarketSocket().BroadcastMessage(id, res)
	}
}
