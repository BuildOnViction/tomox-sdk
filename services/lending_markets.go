package services

import (
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// LendingMarketsService struct with daos required, responsible for communicating with daos.
// LendingMarketsService functions are responsible for interacting with daos and implements business logics.
type LendingMarketsService struct {
	LendingPairDao      interfaces.LendingPairDao
	LendingOhlcvService interfaces.LendingOhlcvService
}

// NewLendingMarketsService returns a new instance of TradeService
func NewLendingMarketsService(
	lendingPairDao interfaces.LendingPairDao,
	lendingOhlcvService interfaces.LendingOhlcvService,
) *LendingMarketsService {
	return &LendingMarketsService{
		LendingPairDao:      lendingPairDao,
		LendingOhlcvService: lendingOhlcvService,
	}
}

// Subscribe market
func (s *LendingMarketsService) Subscribe(c *ws.Client) {
	socket := ws.GetLendingMarketSocket()
	id := utils.GetLendingMarketsChannelID(ws.LendingMarketsChannel)

	tick, err := s.LendingOhlcvService.GetAllTokenPairData()
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	data := &types.LendingMarketData{
		PairData: tick,
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, data)
}

// UnsubscribeChannel UnsubscribeChannel lending market socket
func (s *LendingMarketsService) UnsubscribeChannel(c *ws.Client) {
	socket := ws.GetLendingMarketSocket()

	id := utils.GetLendingMarketsChannelID(ws.LendingMarketsChannel)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe unsubscribe lending market socket
func (s *LendingMarketsService) Unsubscribe(c *ws.Client) {
	socket := ws.GetLendingMarketSocket()
	socket.Unsubscribe(c)
}
