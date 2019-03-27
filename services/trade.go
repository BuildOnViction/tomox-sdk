package services

import (
	"context"
	"github.com/globalsign/mgo"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/rabbitmq"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
	"github.com/tomochain/dex-server/ws"

	"github.com/ethereum/go-ethereum/common"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type TradeService struct {
	tradeDao interfaces.TradeDao
	broker   *rabbitmq.Connection
}

// NewTradeService returns a new instance of TradeService
func NewTradeService(
	TradeDao interfaces.TradeDao,
	broker *rabbitmq.Connection,
) *TradeService {
	return &TradeService{
		TradeDao,
		broker,
	}
}

// Subscribe
func (s *TradeService) Subscribe(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetTradeSocket()

	numTrades := 40
	trades, err := s.GetSortedTrades(bt, qt, numTrades)
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

func (s *TradeService) GetSortedTradesByUserAddress(a common.Address, limit ...int) ([]*types.Trade, error) {
	return s.tradeDao.GetSortedTradesByUserAddress(a, limit...)
}

func (s *TradeService) GetSortedTrades(bt, qt common.Address, n int) ([]*types.Trade, error) {
	return s.tradeDao.GetSortedTrades(bt, qt, n)
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

func (s *TradeService) UpdatePendingTrade(t *types.Trade, txh common.Hash) (*types.Trade, error) {
	t.Status = types.PENDING
	t.TxHash = txh

	updated, err := s.tradeDao.FindAndModify(t.Hash, t)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *TradeService) UpdateSuccessfulTrade(t *types.Trade) (*types.Trade, error) {
	t.Status = types.SUCCESS

	updated, err := s.tradeDao.FindAndModify(t.Hash, t)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	return updated, nil
}

func (s *TradeService) UpdateTradeTxHash(tr *types.Trade, txh common.Hash) error {
	tr.TxHash = txh

	err := s.tradeDao.UpdateByHash(tr.Hash, tr)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}

func (s *TradeService) WatchChanges() {
	s.tradeDao.WatchChanges(s.handleChangeStream)
}

func (s *TradeService) handleChangeStream(ctx context.Context, ct *mgo.ChangeStream) {
	for {
		select {
		case <-ctx.Done(): // if parent context was cancelled
			err := ct.Close() // close the stream
			if err != nil {
				logger.Error("Change stream closed")
			}
			return //exiting from the func
		default:
			ev := types.ChangeEvent{}

			//getting next item from the steam
			ok := ct.Next(&ev)

			//if data from the stream wasn't un-marshaled, we get ok == false as a result
			//so we need to call Err() method to get info why
			//it'll be nil if we just have no data
			if !ok {
				err := ct.Err()
				if err != nil {
					//if err is not nil, it means something bad happened, let's finish our func
					logger.Error(err)
					return
				}
			}

			//if item from the stream un-marshaled successfully, do something with it
			if ok {
				logger.Debug(ev.OperationType)
				s.HandleDocumentType(ev)
			}
		}
	}
}

func (s *TradeService) HandleDocumentType(ev types.ChangeEvent) error {
	res := &types.EngineResponse{}

	switch ev.OperationType {
	case types.OPERATION_TYPE_INSERT:
		if ev.FullDocument.Status == types.TradeStatusPending {

		} else if ev.FullDocument.Status == types.TradeStatusSuccess {

		} else if ev.FullDocument.Status == types.TradeStatusError {

		} else {

		}
		break
	case types.OPERATION_TYPE_UPDATE:
		break
	case types.OPERATION_TYPE_REPLACE:
		break
	default:
		break
	}

	utils.PrintJSON(res)
	err := s.broker.PublishEngineResponse(res)
	if err != nil {
		logger.Error(err)
		return err
	}

	return nil
}
