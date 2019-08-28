package services

import (
	"math/big"
	"strings"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/math"
)

// PairService struct with daos required, responsible for communicating with daos.
// PairService functions are responsible for interacting with daos and implements business logics.
type PairService struct {
	pairDao      interfaces.PairDao
	tokenDao     interfaces.TokenDao
	tradeDao     interfaces.TradeDao
	orderDao     interfaces.OrderDao
	fiatPriceDao interfaces.FiatPriceDao
	eng          interfaces.Engine

	provider interfaces.EthereumProvider
}

// NewPairService returns a new instance of balance service
func NewPairService(
	pairDao interfaces.PairDao,
	tokenDao interfaces.TokenDao,
	tradeDao interfaces.TradeDao,
	orderDao interfaces.OrderDao,
	fiatPriceDao interfaces.FiatPriceDao,
	eng interfaces.Engine,
	provider interfaces.EthereumProvider,
) *PairService {

	return &PairService{pairDao, tokenDao, tradeDao, orderDao, fiatPriceDao, eng, provider}
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

func (s *PairService) GetListedPairs() ([]types.Pair, error) {
	return s.pairDao.GetListedPairs()
}

func (s *PairService) GetUnlistedPairs() ([]types.Pair, error) {
	return s.pairDao.GetUnlistedPairs()
}

// GetTokenPairData get tick of a token pair
func (s *PairService) GetTokenPairData(bt, qt common.Address) (*types.PairData, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -1).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	tradeDataQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": start,
					"$lt":  end,
				},
				"status":     bson.M{"$in": []string{types.SUCCESS}},
				"baseToken":  bt.Hex(),
				"quoteToken": qt.Hex(),
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"count":     bson.M{"$sum": one},
				"open":      bson.M{"$first": "$pricepoint"},
				"openTime":  bson.M{"$first": "$createdAt"},
				"high":      bson.M{"$max": "$pricepoint"},
				"low":       bson.M{"$min": "$pricepoint"},
				"close":     bson.M{"$last": "$pricepoint"},
				"closeTime": bson.M{"$last": "$createdAt"},
				"volume":    bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	bidsQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":       "BUY",
				"baseToken":  bt.Hex(),
				"quoteToken": qt.Hex(),
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$quantity"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$max": "$price"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status":     bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":       "SELL",
				"baseToken":  bt.Hex(),
				"quoteToken": qt.Hex(),
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$quantity"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$min": "$price"},
			},
		},
	}

	tradeData, err := s.tradeDao.Aggregate(tradeDataQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bidsData, err := s.orderDao.Aggregate(bidsQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	asksData, err := s.orderDao.Aggregate(asksQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	if tradeData == nil || bidsData == nil || asksData == nil {
		return nil, nil
	}
	pairData := &types.PairData{
		Pair:         types.PairID{BaseToken: bt, QuoteToken: qt},
		Open:         big.NewInt(0),
		High:         big.NewInt(0),
		Low:          big.NewInt(0),
		Volume:       big.NewInt(0),
		Close:        big.NewInt(0),
		CloseBaseUsd: big.NewFloat(0),
		Count:        big.NewInt(0),
		OrderVolume:  big.NewInt(0),
		OrderCount:   big.NewInt(0),
		BidPrice:     big.NewInt(0),
		AskPrice:     big.NewInt(0),
		Price:        big.NewInt(0),
	}
	if len(tradeData) > 0 {
		t := tradeData[0]
		pairData.Open = t.Open
		pairData.High = t.High
		pairData.Low = t.Low
		pairData.Volume = t.Volume
		pairData.Close = t.Close
		pairData.Count = t.Count
		pairName := t.Pair.PairName
		ps := strings.Split(pairName, "/")
		if len(ps) == 2 {
			fiatItem, err := s.fiatPriceDao.GetLastPriceCurrentByTime(ps[1], t.CloseTime)
			if err == nil {
				pairData.CloseBaseUsd, _ = pairData.CloseBaseUsd.SetString(fiatItem.Price)
			}
		}

	}
	if len(bidsData) > 0 {
		b := bidsData[0]
		pairData.OrderVolume = b.OrderVolume
		pairData.OrderCount = b.OrderCount
		pairData.BidPrice = b.BestPrice
	}
	if len(asksData) > 0 {
		a := asksData[0]
		pairData.OrderVolume = math.Add(pairData.OrderVolume, a.OrderVolume)
		pairData.OrderCount = math.Add(pairData.OrderCount, a.OrderCount)
		pairData.AskPrice = a.BestPrice

		if math.IsNotEqual(pairData.BidPrice, big.NewInt(0)) && math.IsNotEqual(pairData.AskPrice, big.NewInt(0)) {
			pairData.Price = math.Avg(pairData.BidPrice, pairData.AskPrice)
		} else {
			pairData.Price = big.NewInt(0)
		}
	}

	return pairData, nil
}

// GetAllTokenPairData get tick of all tokens
func (s *PairService) GetAllTokenPairData() ([]*types.PairData, error) {
	now := time.Now()
	end := time.Unix(now.Unix(), 0)
	start := time.Unix(now.AddDate(0, 0, -1).Unix(), 0)
	one, _ := bson.ParseDecimal128("1")

	pairs, err := s.pairDao.GetActivePairs()
	if err != nil {
		return nil, err
	}

	tradeDataQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"createdAt": bson.M{
					"$gte": start,
					"$lt":  end,
				},
				"status": bson.M{"$in": []string{types.SUCCESS}},
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"count":     bson.M{"$sum": one},
				"open":      bson.M{"$first": "$pricepoint"},
				"openTime":  bson.M{"$first": "$createdAt"},
				"high":      bson.M{"$max": "$pricepoint"},
				"low":       bson.M{"$min": "$pricepoint"},
				"close":     bson.M{"$last": "$pricepoint"},
				"closeTime": bson.M{"$last": "$createdAt"},
				"volume":    bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
			},
		},
	}

	bidsQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "BUY",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$quantity"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$max": "$price"},
			},
		},
	}

	asksQuery := []bson.M{
		bson.M{
			"$match": bson.M{
				"status": bson.M{"$in": []string{"OPEN", "PARTIAL_FILLED"}},
				"side":   "SELL",
			},
		},
		bson.M{
			"$group": bson.M{
				"_id": bson.M{
					"pairName":   "$pairName",
					"baseToken":  "$baseToken",
					"quoteToken": "$quoteToken",
				},
				"orderCount": bson.M{"$sum": one},
				"orderVolume": bson.M{
					"$sum": bson.M{
						"$subtract": []bson.M{bson.M{"$toDecimal": "$quantity"}, bson.M{"$toDecimal": "$filledAmount"}},
					},
				},
				"bestPrice": bson.M{"$min": "$price"},
			},
		},
	}

	tradeData, err := s.tradeDao.Aggregate(tradeDataQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	bidsData, err := s.orderDao.Aggregate(bidsQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	asksData, err := s.orderDao.Aggregate(asksQuery)
	if err != nil {
		logger.Error(err)
		return nil, err
	}

	pairsData := make([]*types.PairData, 0)
	for _, p := range pairs {
		pairData := &types.PairData{
			Pair:         types.PairID{PairName: p.Name(), BaseToken: p.BaseTokenAddress, QuoteToken: p.QuoteTokenAddress},
			Open:         big.NewInt(0),
			High:         big.NewInt(0),
			Low:          big.NewInt(0),
			Volume:       big.NewInt(0),
			Close:        big.NewInt(0),
			CloseBaseUsd: big.NewFloat(0),
			Count:        big.NewInt(0),
			OrderVolume:  big.NewInt(0),
			OrderCount:   big.NewInt(0),
			BidPrice:     big.NewInt(0),
			AskPrice:     big.NewInt(0),
			Price:        big.NewInt(0),
		}

		for _, t := range tradeData {
			if t.AddressCode() == p.AddressCode() {
				pairData.Open = t.Open
				pairData.High = t.High
				pairData.Low = t.Low
				pairData.Volume = t.Volume
				pairData.Close = t.Close
				pairData.Count = t.Count
				fiatItem, err := s.fiatPriceDao.GetLastPriceCurrentByTime(p.QuoteTokenSymbol, t.CloseTime)
				if err == nil {
					pairData.CloseBaseUsd, _ = pairData.CloseBaseUsd.SetString(fiatItem.Price)
				}

			}
		}

		for _, o := range bidsData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = o.OrderVolume
				pairData.OrderCount = o.OrderCount
				pairData.BidPrice = o.BestPrice
			}
		}

		for _, o := range asksData {
			if o.AddressCode() == p.AddressCode() {
				pairData.OrderVolume = math.Add(pairData.OrderVolume, o.OrderVolume)
				pairData.OrderCount = math.Add(pairData.OrderCount, o.OrderCount)
				pairData.AskPrice = o.BestPrice

				if math.IsNotEqual(pairData.BidPrice, big.NewInt(0)) && math.IsNotEqual(pairData.AskPrice, big.NewInt(0)) {
					pairData.Price = math.Avg(pairData.BidPrice, pairData.AskPrice)
				} else {
					pairData.Price = big.NewInt(0)
				}
			}
		}

		pairsData = append(pairsData, pairData)
	}

	return pairsData, nil
}
