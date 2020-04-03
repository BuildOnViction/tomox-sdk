package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"

	"github.com/tomochain/tomox-sdk/ws"
)

// LendingOrderBookService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type LendingOrderBookService struct {
	lendingOrderDao interfaces.LendingOrderDao
}

// NewLendingOrderBookService returns a new instance of balance service
func NewLendingOrderBookService(
	lendingOrderDao interfaces.LendingOrderDao,
) *LendingOrderBookService {
	return &LendingOrderBookService{lendingOrderDao}
}

// GetLendingOrderBook fetches orderbook from engine and returns it as an map[string]interface
func (s *LendingOrderBookService) GetLendingOrderBook(term uint64, lendingToken common.Address) (*types.LendingOrderBook, error) {
	borrow, lend, err := s.lendingOrderDao.GetLendingOrderBook(term, lendingToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ob := &types.LendingOrderBook{
		Name:   utils.GetLendingOrderBookChannelID(term, lendingToken),
		Lend:   lend,
		Borrow: borrow,
	}

	return ob, nil
}

func (s *LendingOrderBookService) GetLendingOrderBookInDb(term uint64, lendingToken common.Address) (*types.LendingOrderBook, error) {
	borrow, lend, err := s.lendingOrderDao.GetLendingOrderBookInDb(term, lendingToken)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	ob := &types.LendingOrderBook{
		Name:   utils.GetLendingOrderBookChannelID(term, lendingToken),
		Lend:   lend,
		Borrow: borrow,
	}

	return ob, nil
}

// SubscribeLendingOrderBook is responsible for handling incoming orderbook subscription messages
// It makes an entry of connection in pairSocket corresponding to pair,unit and duration
func (s *LendingOrderBookService) SubscribeLendingOrderBook(c *ws.Client, term uint64, lendingToken common.Address) {
	socket := ws.GetLendingOrderBookSocket()

	ob, err := s.GetLendingOrderBook(term, lendingToken)
	if err != nil {
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetLendingOrderBookChannelID(term, lendingToken)
	err = socket.Subscribe(id, c)
	if err != nil {
		msg := map[string]string{"Message": err.Error()}
		socket.SendErrorMessage(c, msg)
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, ob)
}

// UnsubscribeLendingOrderBook is responsible for handling incoming orderbook unsubscription messages
func (s *LendingOrderBookService) UnsubscribeLendingOrderBook(c *ws.Client) {
	socket := ws.GetLendingOrderBookSocket()
	socket.Unsubscribe(c)
}

// UnsubscribeLendingOrderBookChannel  is responsible for handling incoming orderbook unsubscription messages
func (s *LendingOrderBookService) UnsubscribeLendingOrderBookChannel(c *ws.Client, term uint64, lendingToken common.Address) {
	socket := ws.GetLendingOrderBookSocket()
	id := utils.GetLendingOrderBookChannelID(term, lendingToken)
	socket.UnsubscribeChannel(id, c)
}
