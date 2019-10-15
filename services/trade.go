package services

import (
	"context"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo"
	"github.com/globalsign/mgo/bson"
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
	accountDao      interfaces.AccountDao
	notificationDao interfaces.NotificationDao
	broker          *rabbitmq.Connection
}

// NewTradeService returns a new instance of TradeService
func NewTradeService(
	orderdao interfaces.OrderDao,
	tradeDao interfaces.TradeDao,
	accountDao interfaces.AccountDao,
	notificationDao interfaces.NotificationDao,
	broker *rabbitmq.Connection,
) *TradeService {
	return &TradeService{
		OrderDao:        orderdao,
		tradeDao:        tradeDao,
		accountDao:      accountDao,
		notificationDao: notificationDao,
		broker:          broker,
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
	pipeline := []bson.M{}

	ct, err := s.tradeDao.GetCollection().Watch(pipeline, mgo.ChangeStreamOptions{FullDocument: mgo.UpdateLookup})

	if err != nil {
		logger.Error("Failed to open change stream")
		return //exiting func
	}

	defer ct.Close()

	// Watch the event again in case there is error and function returned
	defer s.WatchChanges()

	ctx := context.Background()

	//Handling change stream in a cycle
	for {
		select {
		case <-ctx.Done(): // if parent context was cancelled
			err := ct.Close() // close the stream
			if err != nil {
				logger.Error("Change stream closed")
			}
			return //exiting from the func
		default:
			ev := types.TradeChangeEvent{}

			//getting next item from the steam
			ok := ct.Next(&ev)

			//if data from the stream wasn't un-marshaled, we get ok == false as a result
			//so we need to call Err() method to get info why
			//it'll be nil if we just have no data
			if !ok {
				err := ct.Err()
				if err != nil {
					logger.Error(err)
				}
			}

			//if item from the stream un-marshaled successfully, do something with it
			if ok {
				logger.Debugf("Operation Type: %s", ev.OperationType)
				s.HandleDocumentType(ev)
			}
		}
	}
}

// HandleDocumentType handle trade insert/update db trigger
func (s *TradeService) HandleDocumentType(ev types.TradeChangeEvent) error {
	res := &types.EngineResponse{}

	switch ev.OperationType {
	case types.OPERATION_TYPE_INSERT:
		res.Status = types.TradeAdded
		res.Trade = ev.FullDocument
		break
	case types.OPERATION_TYPE_UPDATE:
	case types.OPERATION_TYPE_REPLACE:
		res.Status = types.TradeUpdated
		res.Trade = ev.FullDocument
		break
	default:
		break
	}

	if res.Status != "" {
		err := s.broker.PublishTradeResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
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
	s.broadcastTradeUpdate(trades)
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
