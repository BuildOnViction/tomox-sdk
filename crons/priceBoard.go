package crons

import (
	"log"

	"github.com/ethereum/go-ethereum/common"
	"github.com/robfig/cron"
	"github.com/tomochain/dex-server/types"
)

// tickStreamingCron takes instance of cron.Cron and adds tickStreaming
// crons according to the durations mentioned in config/app.yaml file
func (s *CronService) startPriceBoardCron(c *cron.Cron) {
	c.AddFunc("*/3 * * * * *", s.getPriceBoardData())
}

// tickStream function fetches latest tick based on unit and duration for each pair
// and broadcasts the tick to the client subscribed to pair's respective channel
func (s *CronService) getPriceBoardData() func() {
	return func() {
		p := make([]types.PairAddresses, 0)
		p = []types.PairAddresses{{
			BaseToken:  common.HexToAddress("0x0e11C49B66b3d277b4292d8d86fD620b4342043A"),
			QuoteToken: common.HexToAddress("0x0000000000000000000000000000000000000001"),
		}}

		// Fix the value at 1 day because we only care about 24h change
		duration := int64(1)
		unit := "day"

		_, err := s.priceBoardService.GetPriceBoardData(p, duration, unit)
		if err != nil {
			log.Printf("%s", err)
			return
		}

		//for _, tick := range ticks {
		//	baseTokenAddress := tick.Pair.BaseToken
		//	quoteTokenAddress := tick.Pair.QuoteToken
		//	id := utils.GetTickChannelID(baseTokenAddress, quoteTokenAddress, unit, duration)
		//	ws.GetPriceBoardSocket().BroadcastMessage(id, tick)
		//}
	}
}
