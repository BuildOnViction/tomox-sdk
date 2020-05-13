package services

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type PairService struct {
	pairDao  interfaces.PairDao
	tokenDao interfaces.TokenDao
	tradeDao interfaces.TradeDao
	orderDao interfaces.OrderDao
	ohlcv    interfaces.OHLCVService
	eng      interfaces.Engine

	provider interfaces.EthereumProvider
}

// NewPairService returns a new instance of balance service
func NewPairService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	tradeDao interfaces.TradeDao,
	orderDao interfaces.OrderDao,
	ohlcv interfaces.OHLCVService,
	eng interfaces.Engine,
	provider interfaces.EthereumProvider,
) *PairService {

	return &PairService{pairDao, tokenDao, tradeDao, orderDao, ohlcv, eng, provider}
}

func (s *PairService) CreatePairs(addr common.Address) ([]*types.Pair, error) {
	quotes, err := s.tokenDao.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	base, err := s.tokenDao.GetByAddress(addr)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if base == nil {
		symbol, err := s.provider.Symbol(addr)
		if err != nil {
			logger.Error(err)
			return nil, ErrNoContractCode
		}

		decimals, err := s.provider.Decimals(addr)
		if err != nil {
			logger.Error(err)
			return nil, ErrNoContractCode
		}

		base = &types.Token{
			Symbol:   symbol,
			Address:  addr,
			Decimals: int(decimals),
			Active:   true,
			Listed:   false,
			Quote:    false,
		}

		err = s.tokenDao.Create(base)
		if err != nil {
			logger.Error(err)
			return nil, err
		}
	}

	pairs := []*types.Pair{}
	for _, q := range quotes {
		p, err := s.pairDao.GetByTokenAddress(addr, q.Address)
		if err != nil {
			logger.Error(err)
			return nil, err
		}

		if p == nil {
			p := types.Pair{
				QuoteTokenSymbol:   q.Symbol,
				QuoteTokenAddress:  q.Address,
				QuoteTokenDecimals: q.Decimals,
				BaseTokenSymbol:    base.Symbol,
				BaseTokenAddress:   base.Address,
				BaseTokenDecimals:  base.Decimals,
				Active:             true,
				Listed:             false,
				MakeFee:            q.MakeFee,
				TakeFee:            q.TakeFee,
			}

			err := s.pairDao.Create(&p)
			if err != nil {
				logger.Error(err)
				return nil, err
			}

			pairs = append(pairs, &p)
		}
	}

	return pairs, nil
}

// Create function is responsible for inserting new pair in DB.
// It checks for existence of tokens in DB first
func (s *PairService) Create(pair *types.Pair) error {
	p, err := s.pairDao.GetByTokenAddress(pair.BaseTokenAddress, pair.QuoteTokenAddress)
	if err != nil {
		return err
	}

	if p != nil {
		return ErrPairExists
	}

	bt, err := s.tokenDao.GetByAddress(pair.BaseTokenAddress)
	if err != nil {
		return err
	}

	if bt == nil {
		return ErrBaseTokenNotFound
	}

	st, err := s.tokenDao.GetByAddress(pair.QuoteTokenAddress)
	if err != nil {
		return err
	}

	if st == nil {
		return ErrQuoteTokenNotFound
	}

	if !st.Quote {
		return ErrQuoteTokenInvalid
	}

	pair.QuoteTokenSymbol = st.Symbol
	pair.QuoteTokenAddress = st.ContractAddress
	pair.QuoteTokenDecimals = st.Decimals
	pair.BaseTokenSymbol = bt.Symbol
	pair.BaseTokenAddress = bt.ContractAddress
	pair.BaseTokenDecimals = bt.Decimals
	err = s.pairDao.Create(pair)
	if err != nil {
		return err
	}

	return nil
}

// GetByID fetches details of a pair using its mongo ID
func (s *PairService) GetByID(id bson.ObjectId) (*types.Pair, error) {
	return s.pairDao.GetByID(id)
}

// GetByTokenAddress fetches details of a pair using contract address of
// its constituting tokens
func (s *PairService) GetByTokenAddress(bt, qt common.Address) (*types.Pair, error) {
	return s.pairDao.GetByTokenAddress(bt, qt)
}

// GetAll is reponsible for fetching all the pairs in the DB
func (s *PairService) GetAll() ([]types.Pair, error) {
	return s.pairDao.GetAll()
}

// GetAllByCoinbase get all pair by coinbase
func (s *PairService) GetAllByCoinbase(addr common.Address) ([]types.Pair, error) {
	return s.pairDao.GetAllByCoinbase(addr)
}

// GetTokenPairData get tick of a token pair
func (s *PairService) GetTokenPairData(bt, qt common.Address) (*types.PairData, error) {
	pairData := s.ohlcv.GetTokenPairData(bt, qt)
	if pairData == nil {
		return nil, nil
	}
	bidPrice, err := s.orderDao.GetBestBid(bt, qt)
	if err == nil && bidPrice != nil {
		pairData.BidPrice = bidPrice.Price
	}
	askPrice, err := s.orderDao.GetBestAsk(bt, qt)
	if err == nil && askPrice != nil {
		pairData.AskPrice = askPrice.Price
	}
	return pairData, nil
}

// GetAllTokenPairData get tick of all tokens
func (s *PairService) GetAllTokenPairData() ([]*types.PairData, error) {
	return s.ohlcv.GetAllTokenPairData()
}
