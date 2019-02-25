package services

import (
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
	"github.com/tomochain/dex-server/ws"

	"github.com/ethereum/go-ethereum/common"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type PriceBoardService struct {
	tradeDao interfaces.TradeDao
}

// NewTradeService returns a new instance of TradeService
func NewPriceBoardService(TradeDao interfaces.TradeDao) *PriceBoardService {
	return &PriceBoardService{TradeDao}
}

// Subscribe
func (s *PriceBoardService) Subscribe(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetPriceBoardSocket()

	numTrades := 40
	trades, err := s.GetSortedTrades(bt, qt, numTrades)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetPriceBoardChannelID(bt, qt)
	err = socket.Subscribe(id, c)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, trades)
}

// Unsubscribe
func (s *PriceBoardService) UnsubscribeChannel(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetPriceBoardSocket()

	id := utils.GetPriceBoardChannelID(bt, qt)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe
func (s *PriceBoardService) Unsubscribe(c *ws.Client) {
	socket := ws.GetPriceBoardSocket()
	socket.Unsubscribe(c)
}

func (s *PriceBoardService) GetSortedTrades(bt, qt common.Address, n int) ([]*types.Trade, error) {
	return s.tradeDao.GetSortedTrades(bt, qt, n)
}
