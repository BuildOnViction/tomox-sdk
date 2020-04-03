package services

import (
	"context"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/rabbitmq"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

// LendingTradeService struct with daos required, responsible for communicating with daos.
// LendingTradeService functions are responsible for interacting with daos and implements business logics.
type LendingTradeService struct {
	lendingDao          interfaces.LendingOrderDao
	lendingTradeDao     interfaces.LendingTradeDao
	notificationDao     interfaces.NotificationDao
	broker              *rabbitmq.Connection
	bulkLendingTrades   map[string][]*types.LendingTrade
	mutext              sync.RWMutex
	tradeNotifyCallback func(*types.LendingTrade)
}

// NewLendingTradeService returns a new instance of LendingTradeService
func NewLendingTradeService(
	lendingdao interfaces.LendingOrderDao,
	lendingTradeDao interfaces.LendingTradeDao,
	notificationDao interfaces.NotificationDao,
	broker *rabbitmq.Connection,
) *LendingTradeService {
	bulkLendingTrades := make(map[string][]*types.LendingTrade)
	return &LendingTradeService{
		lendingDao:          lendingdao,
		lendingTradeDao:     lendingTradeDao,
		notificationDao:     notificationDao,
		broker:              broker,
		bulkLendingTrades:   bulkLendingTrades,
		mutext:              sync.RWMutex{},
		tradeNotifyCallback: nil,
	}
}

// RegisterNotify register a only trade notify function
func (s *LendingTradeService) RegisterNotify(fn func(*types.LendingTrade)) {
	s.tradeNotifyCallback = fn
}

// Subscribe Subscribe lending trade channel
func (s *LendingTradeService) Subscribe(c *ws.Client, term uint64, lendingToken common.Address) {
	socket := ws.GetLendingTradeSocket()
	numTrades := types.DefaultLimit
	trades, err := s.GetLendingTradeByOrderBook(term, lendingToken, 0, 0, numTrades)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetLendingTradeChannelID(term, lendingToken)
	err = socket.Subscribe(id, c)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, trades)
}

// UnsubscribeChannel unsubscribe lending channel
func (s *LendingTradeService) UnsubscribeChannel(c *ws.Client, term uint64, lendingToken common.Address) {
	socket := ws.GetLendingTradeSocket()

	id := utils.GetLendingTradeChannelID(term, lendingToken)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe unsubscribe lending channel
func (s *LendingTradeService) Unsubscribe(c *ws.Client) {
	socket := ws.GetLendingTradeSocket()
	socket.Unsubscribe(c)
}

// GetLendingTradeByOrderBook get sorted lending trade from term and lending tokens
func (s *LendingTradeService) GetLendingTradeByOrderBook(tern uint64, lendingToken common.Address, from, to int64, n int) ([]*types.LendingTrade, error) {
	return s.lendingTradeDao.GetLendingTradeByOrderBook(tern, lendingToken, from, to, n)
}

// WatchChanges watch changing trade database
func (s *LendingTradeService) WatchChanges() {

	ct, sc, err := s.lendingTradeDao.Watch()

	if err != nil {
		logger.Error("Failed to open change stream")
		return //exiting func
	}

	defer ct.Close()
	defer sc.Close()

	// Watch the event again in case there is error and function returned
	defer s.WatchChanges()

	ctx := context.Background()

	go func() {
		for {
			<-time.After(500 * time.Millisecond)
			s.processBulkLendingTrades()
		}
	}()

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
			ev := types.LendingTradeChangeEvent{}

			//getting next item from the steam
			ok := ct.Next(&ev)

			//if data from the stream wasn't un-marshaled, we get ok == false as a result
			//so we need to call Err() method to get info why
			//it'll be nil if we just have no data
			if !ok {
				err := ct.Err()
				if err != nil {
					logger.Error(err)
					return
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

func (s *LendingTradeService) processBulkLendingTrades() {
	s.mutext.Lock()
	defer s.mutext.Unlock()

	bulkPairs := make(map[string]bool)
	for id, trades := range s.bulkLendingTrades {
		bulkPairs[id] = true
		if len(trades) > 0 {
			ws.GetLendingTradeSocket().BroadcastMessage(id, trades)
		}
	}
	s.bulkLendingTrades = make(map[string][]*types.LendingTrade)
	pairs := make([]string, 0)
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
func (s *LendingTradeService) HandleDocumentType(ev types.LendingTradeChangeEvent) error {
	res := &types.EngineResponse{}

	if ev.OperationType == types.OPERATION_TYPE_INSERT {
		res.Status = types.TradeAdded
		res.LendingTrade = ev.FullDocument
	}
	if ev.OperationType == types.OPERATION_TYPE_UPDATE || ev.OperationType == types.OPERATION_TYPE_REPLACE {
		res.Status = types.TradeUpdated
		res.LendingTrade = ev.FullDocument
	}

	if res.Status != "" {
		err := s.broker.PublishLendingTradeResponse(res)
		if err != nil {
			logger.Error(err)
			return err
		}
	}

	return nil
}

// HandleLendingTradeResponse listens to messages incoming from the engine and handles websocket
// responses and database updates accordingly
func (s *LendingTradeService) HandleLendingTradeResponse(res *types.EngineResponse) error {
	switch res.Status {
	case types.TradeAdded:
		s.HandleOperationInsert(res.LendingTrade)
		break
	case types.TradeUpdated:
		s.HandleOperationUpdate(res.LendingTrade)
		break
	}

	return nil
}

// HandleOperationInsert sent WS messages to client when a trade is created with status "PENDING""
func (s *LendingTradeService) HandleOperationInsert(trade *types.LendingTrade) error {
	m := &types.LendingMatches{LendingTrades: []*types.LendingTrade{trade}}
	borrower, err := s.lendingDao.GetByHash(trade.BorrowingOrderHash)
	if err != nil {
		logger.Error(err)
		return errors.New("Can not find borrower order")
	}
	m.Borrowing = borrower
	mo, err := s.lendingDao.GetByHash(trade.InvestingOrderHash)
	if err != nil {
		logger.Error(err)
		return errors.New("Can not find maker order")
	}
	m.Investing = []*types.LendingOrder{mo}
	s.HandleTradeSuccess(m)
	return nil
}

// HandleOperationUpdate sent WS messages to client when a trade is updated with status "SUCCESS" or "ERROR"
func (s *LendingTradeService) HandleOperationUpdate(trade *types.LendingTrade) error {
	m := &types.LendingMatches{LendingTrades: []*types.LendingTrade{trade}}
	borrower, err := s.lendingDao.GetByHash(trade.BorrowingOrderHash)
	if err != nil {
		logger.Error(err)
		return errors.New("Can not find borrower order")
	}
	m.Borrowing = borrower
	mo, err := s.lendingDao.GetByHash(trade.InvestingOrderHash)
	if err != nil {
		logger.Error(err)
		return errors.New("Can not find maker order")
	}
	m.Investing = []*types.LendingOrder{mo}
	if trade.Status == types.TradeStatusSuccess {
		s.HandleTradeSuccess(m)
	}

	return nil
}

// HandleTradeSuccess handle order match success
func (s *LendingTradeService) HandleTradeSuccess(m *types.LendingMatches) {
	trades := m.LendingTrades
	for _, t := range trades {
		investor := t.Investor
		borrower := t.Borrower

		s.saveBulkTrades(t)

		ws.SendLendingOrderMessage("LENDING_ORDER_SUCCESS", investor, types.LendingOrderSuccessPayload{LendingMatches: m})
		ws.SendLendingOrderMessage("LENDING_ORDER_SUCCESS", borrower, types.LendingOrderSuccessPayload{LendingMatches: m})
		s.notificationDao.Create(&types.Notification{
			Recipient: investor,
			Message: types.Message{
				MessageType: "LENDING_ORDER_SUCCESS",
				Description: t.Hash.Hex(),
			},
			Type:   types.TypeLog,
			Status: types.StatusUnread,
		})
		s.notificationDao.Create(&types.Notification{
			Recipient: borrower,
			Message: types.Message{
				MessageType: "LENDING_ORDER_SUCCESS",
				Description: t.Hash.Hex(),
			},
			Type:   types.TypeLog,
			Status: types.StatusUnread,
		})
	}
	if s.tradeNotifyCallback != nil {
		s.tradeNotifyCallback(trades[0])
	}
}

func (s *LendingTradeService) saveBulkTrades(t *types.LendingTrade) {
	s.mutext.Lock()
	defer s.mutext.Unlock()
	id := utils.GetLendingTradeChannelID(t.Term, t.LendingToken)
	s.bulkLendingTrades[id] = append(s.bulkLendingTrades[id], t)
}

func (s *LendingTradeService) broadcastTickUpdate(pairs []string) {
}

// GetLendingTradesUserHistory get lending trade by history
func (s *LendingTradeService) GetLendingTradesUserHistory(a common.Address, lendingtradeSpec *types.LendingTradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.LendingTradeRes, error) {
	return s.lendingTradeDao.GetLendingTradesUserHistory(a, lendingtradeSpec, sortedBy, pageOffset, pageSize)
}

// GetLendingTrades get lending trade
func (s *LendingTradeService) GetLendingTrades(lendingtradeSpec *types.LendingTradeSpec, sortedBy []string, pageOffset int, pageSize int) (*types.LendingTradeRes, error) {
	return s.lendingTradeDao.GetLendingTrades(lendingtradeSpec, sortedBy, pageOffset, pageSize)
}

//GetLendingTradeByTime get lending trade by range time
func (s *LendingTradeService) GetLendingTradeByTime(dateFrom, dateTo int64, pageOffset int, pageSize int) ([]*types.LendingTrade, error) {
	return s.lendingTradeDao.GetLendingTradeByTime(dateFrom, dateTo, pageOffset, pageSize)
}
