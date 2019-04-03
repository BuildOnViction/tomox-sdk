package services

import (
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
	"github.com/tomochain/dex-server/ws"
)

// MarketsService struct with daos required, responsible for communicating with daos.
// MarketsService functions are responsible for interacting with daos and implements business logics.
type MarketsService struct {
}

// NewTradeService returns a new instance of TradeService
func NewMarketsService() *MarketsService {
	return &MarketsService{}
}

// Subscribe
func (s *MarketsService) Subscribe(c *ws.Client) {
	socket := ws.GetMarketSocket()

	// Fix the value at 1 day because we only care about 24h change
	duration := int64(1)
	unit := "day"

	data, err := s.GetMarketsData(
		duration,
		unit,
	)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetMarketsChannelID(ws.MarketsChannel)
	err = socket.Subscribe(id, c)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, data)
}

// Unsubscribe
func (s *MarketsService) UnsubscribeChannel(c *ws.Client) {
	socket := ws.GetMarketSocket()

	id := utils.GetMarketsChannelID(ws.MarketsChannel)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe
func (s *MarketsService) Unsubscribe(c *ws.Client) {
	socket := ws.GetMarketSocket()
	socket.Unsubscribe(c)
}

func (s *MarketsService) GetMarketsData(duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error) {
	res := make([]*types.Tick, 0)

	return res, nil
}
