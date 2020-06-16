package services

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type TradeService struct {
	OrderDao        interfaces.OrderDao
	tradeDao        interfaces.TradeDao
	notificationDao interfaces.NotificationDao
	broker          *rabbitmq.Connection
	ohlcvService    *OHLCVService
	bulkTrades      map[types.PairAddresses][]*types.Trade
	mutext          sync.RWMutex
	tradeDispatcher *TradeDispatcherService
}

// NewTradeService returns a new instance of TradeService
func NewTradeService(
	orderdao interfaces.OrderDao,
	tradeDao interfaces.TradeDao,
	tradeDispatcher *TradeDispatcherService,
	ohlcvService *OHLCVService,
	notificationDao interfaces.NotificationDao,
	broker *rabbitmq.Connection,
) *TradeService {
	bulkTrades := make(map[types.PairAddresses][]*types.Trade)
	return &TradeService{
		OrderDao:        orderdao,
		tradeDao:        tradeDao,
		notificationDao: notificationDao,
		broker:          broker,
		ohlcvService:    ohlcvService,
		bulkTrades:      bulkTrades,
		mutext:          sync.RWMutex{},
		tradeDispatcher: tradeDispatcher,
	}
}

// Subscribe
func (s *TradeService) Subscribe(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetTradeSocket()

	numTrades := types.DefaultLimit
	trades, err := s.GetSortedTrades(bt, qt, 0, 0, numTrades)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetTradeChannelID(bt, qt)
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
func (s *TradeService) UnsubscribeChannel(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetTradeSocket()

	id := utils.GetTradeChannelID(bt, qt)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe
func (s *TradeService) Unsubscribe(c *ws.Client) {
	socket := ws.GetTradeSocket()
	socket.Unsubscribe(c)
}

// GetByPairName fetches all the trades corresponding to a pair using pair's name
func (s *TradeService) GetByPairName(p string) ([]*types.Trade, error) {
	return s.tradeDao.GetByPairName(p)
}

// GetByPairAddress fetches all the trades corresponding to a pair using pair's token address
func (s *TradeService) GetAllTradesByPairAddress(bt, qt common.Address) ([]*types.Trade, error) {
	return s.tradeDao.GetAllTradesByPairAddress(bt, qt)
}

func (s *TradeService) GetSortedTradesByUserAddress(a, bt, qt common.Address, from, to int64, limit ...int) ([]*types.Trade, error) {
	return s.tradeDao.GetSortedTradesByUserAddress(a, bt, qt, from, to, limit...)
}

func (s *TradeService) GetSortedTrades(bt, qt common.Address, from, to int64, n int) ([]*types.Trade, error) {
	return s.tradeDao.GetSortedTrades(bt, qt, from, to, n)
}

// GetByUserAddress fetches all the trades corresponding to a user address
func (s *TradeService) GetByUserAddress(a common.Address) ([]*types.Trade, error) {
	return s.tradeDao.GetByUserAddress(a)
}

// GetByHash fetches all trades corresponding to a trade hash
func (s *TradeService) GetByHash(h common.Hash) (*types.Trade, error) {
	return s.tradeDao.GetByHash(h)
}

func (s *TradeService) GetByMakerOrderHash(h common.Hash) ([]*types.Trade, error) {
	return s.tradeDao.GetByMakerOrderHash(h)
}

func (s *TradeService) GetByTakerOrderHash(h common.Hash) ([]*types.Trade, error) {
	return s.tradeDao.GetByTakerOrderHash(h)
}

func (s *TradeService) GetByOrderHashes(hashes []common.Hash) ([]*types.Trade, error) {
	return s.tradeDao.GetByOrderHashes(hashes)
}

func (s *TradeService) WatchChanges() {
	go func() {
		for {
			<-time.After(500 * time.Millisecond)
			s.processBulkTrades()
		}
	}()
	s.tradeDispatcher.SubscribeTrade(s.HandleDocumentType)
}

func (s *TradeService) processBulkTrades() {
	s.mutext.Lock()
	defer s.mutext.Unlock()

	bulkPairs := make(map[types.PairAddresses]bool)
	for pair, trades := range s.bulkTrades {
		bulkPairs[pair] = true
		if len(trades) > 0 {
			id := utils.GetTradeChannelID(pair.BaseToken, pair.QuoteToken)
			ws.GetTradeSocket().BroadcastMessage(id, trades)
		}
	}
	s.bulkTrades = make(map[types.PairAddresses][]*types.Trade)
	pairs := make([]types.PairAddresses, 0)
	for p, val := range bulkPairs {
		if val {
			pairs = append(pairs, p)
		}
	}
	if len(pairs) > 0 {

		s.broadcastTickUpdate(pairs)
	}
}

// HandleDocumentType handle trade insert/update db trigger
func (s *TradeService) HandleDocumentType(ev *types.TradeChangeEvent) {
	res := &types.EngineResponse{}

	if ev.OperationType == types.OPERATION_TYPE_INSERT {
		res.Status = types.TradeAdded
		res.Trade = ev.FullDocument
	}
	if ev.OperationType == types.OPERATION_TYPE_UPDATE || ev.OperationType == types.OPERATION_TYPE_REPLACE {
		res.Status = types.TradeUpdated
		res.Trade = ev.FullDocument
	}

	if res.Status != "" {
		err := s.broker.PublishTradeResponse(res)
		if err != nil {
			logger.Error(err)
		}
	}

}

// HandleTradeResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *TradeService) HandleTradeResponse(res *types.EngineResponse) error {
	switch res.Status {
	case types.TradeAdded:
		s.HandleOperationInsert(res.Trade)
		break
	case types.TradeUpdated:
		s.HandleOperationUpdate(res.Trade)
		break
	}

	return nil
}

// HandleOperationInsert sent WS messages to client when a trade is created with status "PENDING""
func (s *TradeService) HandleOperationInsert(trade *types.Trade) error {
	m := &types.Matches{Trades: []*types.Trade{trade}}

	to, err := s.OrderDao.GetByHash(trade.TakerOrderHash)

	if err != nil {
		logger.Error(err)
		return errors.New("Can not find taker order")
	}

	m.TakerOrder = to

	mo, err := s.OrderDao.GetByHash(trade.MakerOrderHash)

	if err != nil {
		logger.Error(err)
		return errors.New("Can not find maker order")
	}

	m.MakerOrders = []*types.Order{mo}

	s.HandleTradeSuccess(m)

	return nil
}

// HandleOperationUpdate sent WS messages to client when a trade is updated with status "SUCCESS" or "ERROR"
func (s *TradeService) HandleOperationUpdate(trade *types.Trade) error {
	logger.Debug(trade.Status)
	m := &types.Matches{Trades: []*types.Trade{trade}}

	to, err := s.OrderDao.GetByHash(trade.TakerOrderHash)

	if err != nil {
		logger.Error(err)
		return errors.New("Can not find taker order")
	}

	m.TakerOrder = to

	mo, err := s.OrderDao.GetByHash(trade.MakerOrderHash)

	if err != nil {
		logger.Error(err)
		return errors.New("Can not find maker order")
	}

	m.MakerOrders = []*types.Order{mo}

	if trade.Status == types.TradeStatusSuccess {
		s.HandleTradeSuccess(m)
	}

	return nil
}

// HandleTradeSuccess handle order match success
func (s *TradeService) HandleTradeSuccess(m *types.Matches) {
	trades := m.Trades
	for _, t := range trades {
		maker := t.Maker
		taker := t.Taker

		s.saveBulkTrades(t)

		ws.SendOrderMessage("ORDER_SUCCESS", maker, types.OrderSuccessPayload{Matches: m})
		ws.SendOrderMessage("ORDER_SUCCESS", taker, types.OrderSuccessPayload{Matches: m})
		s.notificationDao.Create(&types.Notification{
			Recipient: taker,
			Message: types.Message{
				MessageType: "ORDER_SUCCESS",
				Description: t.Hash.Hex(),
			},
			Type:   types.TypeLog,
			Status: types.StatusUnread,
		})
		s.notificationDao.Create(&types.Notification{
			Recipient: maker,
			Message: types.Message{
				MessageType: "ORDER_SUCCESS",
				Description: t.Hash.Hex(),
			},
			Type:   types.TypeLog,
			Status: types.StatusUnread,
		})
	}

	s.ohlcvService.NotifyTrade(trades[0])
}

func (s *TradeService) saveBulkTrades(t *types.Trade) {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	pair := types.PairAddresses{BaseToken: t.BaseToken, QuoteToken: t.QuoteToken}
	s.bulkTrades[pair] = append(s.bulkTrades[pair], t)
}

func (s *TradeService) broadcastTickUpdate(pairs []types.PairAddresses) {
	for unit, durations := range app.Config.TickDuration {
		for _, duration := range durations {

			ticks, err := s.ohlcvService.GetOHLCV(pairs, duration, unit)
			if err != nil {
				logger.Error("Get ticks", err)
				return
			}

			for _, tick := range ticks {
				baseTokenAddress := tick.Pair.BaseToken
				quoteTokenAddress := tick.Pair.QuoteToken
				id := utils.GetTickChannelID(baseTokenAddress, quoteTokenAddress, unit, duration)
				ws.GetOHLCVSocket().BroadcastOHLCV(id, tick)
			}
		}
	}
}

func (s *TradeService) broadcastTradeUpdate(trades []*types.Trade) {
	p, err := trades[0].Pair()
	if err != nil {
		logger.Error(err)
		return
	}

	id := utils.GetTradeChannelID(p.BaseTokenAddress, p.QuoteTokenAddress)
	ws.GetTradeSocket().BroadcastMessage(id, trades)
}

// GetTrades filter trade
func (s *TradeService) GetTrades(tradeSpec *types.TradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.TradeRes, error) {
	return s.tradeDao.GetTrades(tradeSpec, sortedBy, pageOffset, pageSize)
}

// GetTradesUserHistory get trade by history
func (s *TradeService) GetTradesUserHistory(a common.Address, tradeSpec *types.TradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.TradeRes, error) {
	return s.tradeDao.GetTradesUserHistory(a, tradeSpec, sortedBy, pageOffset, pageSize)
}
