package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

type LendingPriceBoardService struct {
	lendingPairService  interfaces.LendingPairService
	lendingOhlcvService interfaces.LendingOhlcvService
}

// NewLendingPriceBoardService returns a new instance of LendingPriceBoardService
func NewLendingPriceBoardService(
	lendingPairService interfaces.LendingPairService,
	lendingOhlcvService interfaces.LendingOhlcvService,
) *LendingPriceBoardService {
	return &LendingPriceBoardService{
		lendingPairService:  lendingPairService,
		lendingOhlcvService: lendingOhlcvService,
	}
}

func (s *LendingPriceBoardService) Subscribe(c *ws.Client, term uint64, lendingToken common.Address) {
	socket := ws.GetLendingPriceBoardSocket()
	_, err := s.lendingPairService.GetByLendingID(term, lendingToken)
	if err != nil {
		socket.SendErrorMessage(c, err.Error())
		return
	}
	tick := s.GetLendingPriceBoardData(term, lendingToken)
	id := utils.GetLendingChannelID(term, lendingToken)
	err = socket.Subscribe(id, c)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}
	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, tick)
}

func (s *LendingPriceBoardService) UnsubscribeChannel(c *ws.Client, term uint64, lendingToken common.Address) {
	socket := ws.GetLendingPriceBoardSocket()

	id := utils.GetLendingChannelID(term, lendingToken)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe unsubscribe registered socket
func (s *LendingPriceBoardService) Unsubscribe(c *ws.Client) {
	socket := ws.GetLendingPriceBoardSocket()
	socket.Unsubscribe(c)
}

// GetLendingPriceBoardData get data of 24h change tokens
func (s *LendingPriceBoardService) GetLendingPriceBoardData(term uint64, lendingToken common.Address) *types.LendingTick {
	tick := s.lendingOhlcvService.GetTokenPairData(term, lendingToken)
	if tick != nil {
		return tick
	}
	return &types.LendingTick{}

}
