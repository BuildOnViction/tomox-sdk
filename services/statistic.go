package services

import (
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/big"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/daos"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
)

const (
	duration = 1
	unit     = "hour"
	sideBuy  = "BUY"
	sideSell = "SELL"
)

// TradeStatisticService struct with daos required, responsible for communicating with daos.
// TradeStatisticService functions are responsible for interacting with daos and implements business logics.
type TradeStatisticService struct {
	tokenDao        *daos.TokenDao
	tradeCache      *tradeCache
	tokenCache      map[common.Address]*tokenCache
	lastPairPrice   map[string]*big.Int
	mutex           sync.RWMutex
	tradeDispatcher *TradeDispatcherService
}

type tradeCache struct {
	lastTime int64
	// relayerAddress => pairAddress => userAddress => time => UserTrade
	relayerUserTrades map[common.Address]map[string]map[common.Address]map[int64]*types.UserTrade
}

type cachetradefile struct {
	LastTime          int64              `json:"lastTime"`
	UserTrades        []*types.UserTrade `json:"userTrades"`
	RelayerUserTrades []*types.UserTrade `json:"relayerUserTrades"`
}
type tokenCache struct {
	token    *types.Token
	timelife int64
}

// NewTradeStatisticService init new instance
func NewTradeStatisticService(tokenDao *daos.TokenDao, tradeDispatcher *TradeDispatcherService) *TradeStatisticService {

	cache := &tradeCache{
		relayerUserTrades: make(map[common.Address]map[string]map[common.Address]map[int64]*types.UserTrade),
	}
	return &TradeStatisticService{
		tokenDao:        tokenDao,
		tradeCache:      cache,
		tokenCache:      make(map[common.Address]*tokenCache),
		lastPairPrice:   make(map[string]*big.Int),
		tradeDispatcher: tradeDispatcher,
	}
}

func (s *TradeStatisticService) getTokenByAddress(token common.Address) (*types.Token, error) {
	now := time.Now().Unix()
	if tokenCache, ok := s.tokenCache[token]; ok {
		if now-tokenCache.timelife < cacheTimeLifeMax {
			return tokenCache.token, nil
		}
		delete(s.tokenCache, token)
	}
	t, err := s.tokenDao.GetByAddress(token)
	if err == nil && t != nil {
		s.tokenCache[token] = &tokenCache{
			token:    t,
			timelife: now,
		}
	}
	return t, err
}

// NotifyTrade handle trade insert/update db trigger
func (s *TradeStatisticService) NotifyTrade(trade *types.TradeChangeEvent) {
	s.notifyTrade(trade.FullDocument)
}

// NotifyTrade handle trade insert/update db trigger
func (s *TradeStatisticService) notifyTrade(trade *types.Trade) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	key := s.getPairString(trade.BaseToken, trade.QuoteToken)
	s.lastPairPrice[key] = trade.PricePoint
	s.updateRelayerUserTrade(trade)
	return nil
}

// Init init cache
// ensure add current time frame before trade notify come
func (s *TradeStatisticService) Init() {
	s.loadCache()
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	go func() {
		for {
			select {
			case <-ticker.C:
				err := s.commitCache()
				if err != nil {
					logger.Error(err)
				}
			case <-quit:
				ticker.Stop()
				return
			}
		}
	}()
	s.tradeDispatcher.SubscribeTrade(s.NotifyTrade)
	s.tradeDispatcher.SubscribeFetch(s.fetch)
}

func (s *TradeStatisticService) fetch(trade *types.Trade) {

	s.mutex.Lock()
	s.updateRelayerUserTrade(trade)
	s.mutex.Unlock()

}

func (s *TradeStatisticService) flattenRelayerUserTrades() []*types.UserTrade {
	var relayerUserTrades []*types.UserTrade
	for _, tradebyRelayer := range s.tradeCache.relayerUserTrades {
		for _, tradebyUserAddess := range tradebyRelayer {
			for _, tradeBytime := range tradebyUserAddess {
				for _, trade := range tradeBytime {
					relayerUserTrades = append(relayerUserTrades, trade)
				}
			}
		}
	}
	return relayerUserTrades
}

func (s *TradeStatisticService) commitCache() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	logger.Info("commit trade cache")
	relayerUserTrades := s.flattenRelayerUserTrades()
	cachefile := &cachetradefile{
		LastTime:          s.tradeCache.lastTime,
		RelayerUserTrades: relayerUserTrades,
	}
	cacheData, err := json.Marshal(cachefile)
	if err != nil {
		return err
	}
	file, err := os.Create("statistic.cache")
	defer file.Close()
	if err == nil {
		_, err = file.Write(cacheData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *TradeStatisticService) loadCache() error {
	file, err := os.Open("statistic.cache")
	defer file.Close()
	if err != nil {
		return err
	}
	stats, statsErr := file.Stat()
	if statsErr != nil {
		return statsErr
	}

	size := stats.Size()
	bytes := make([]byte, size)
	bufr := bufio.NewReader(file)
	_, err = bufr.Read(bytes)
	var cache cachetradefile
	err = json.Unmarshal(bytes, &cache)
	if err != nil {
		return err
	}
	for _, t := range cache.RelayerUserTrades {
		s.addRelayerUserTrade(t)
	}
	s.tradeCache.lastTime = cache.LastTime
	return nil
}

func (s *TradeStatisticService) getPairString(baseToken, quoteToken common.Address) string {
	return fmt.Sprintf("%s::%s", baseToken.Hex(), quoteToken.Hex())
}

func (s *TradeStatisticService) parsePairString(key string) (common.Address, common.Address, error) {
	tokens := strings.Split(key, "::")
	if len(tokens) == 2 {
		return common.HexToAddress(tokens[0]), common.HexToAddress(tokens[1]), nil
	}
	return common.Address{}, common.Address{}, errors.New("Invalid Key")
}

func (s *TradeStatisticService) getVolumeByQuote(baseToken, quoteToken common.Address, amount *big.Int, price *big.Int) *big.Int {
	token, err := s.getTokenByAddress(baseToken)
	if err == nil && token != nil {
		baseTokenDecimalBig := big.NewInt(int64(math.Pow10(token.Decimals)))
		p := new(big.Int).Mul(amount, price)
		return new(big.Int).Div(p, baseTokenDecimalBig)
	}
	return big.NewInt(0)
}

// updateRelayerTick update lastest tick, need to be lock
func (s *TradeStatisticService) updateRelayerUserTrade(trade *types.Trade) error {
	tradeTime := trade.CreatedAt.Unix()
	key := s.getPairString(trade.BaseToken, trade.QuoteToken)
	exchange := make(map[common.Address]bool)
	exchange[trade.MakerExchange] = true
	exchange[trade.TakerExchange] = true
	for addr := range exchange {
		if _, ok := s.tradeCache.relayerUserTrades[addr]; !ok {
			s.tradeCache.relayerUserTrades[addr] = make(map[string]map[common.Address]map[int64]*types.UserTrade)
		}
		if _, ok := s.tradeCache.relayerUserTrades[addr][key]; !ok {
			s.tradeCache.relayerUserTrades[addr][key] = make(map[common.Address]map[int64]*types.UserTrade)
		}
		if _, ok := s.tradeCache.relayerUserTrades[addr][key][trade.Maker]; !ok {
			s.tradeCache.relayerUserTrades[addr][key][trade.Maker] = make(map[int64]*types.UserTrade)
		}
		if _, ok := s.tradeCache.relayerUserTrades[addr][key][trade.Taker]; !ok {
			s.tradeCache.relayerUserTrades[addr][key][trade.Taker] = make(map[int64]*types.UserTrade)
		}
	}

	modTime, _ := utils.GetModTime(tradeTime, duration, unit)
	volumeByQuote := s.getVolumeByQuote(trade.BaseToken, trade.QuoteToken, trade.Amount, trade.PricePoint)

	for addr := range exchange {
		if trade.Taker.Hex() == trade.Maker.Hex() {
			if last, ok2 := s.tradeCache.relayerUserTrades[addr][key][trade.Taker][modTime]; ok2 {
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
				last.Volume = big.NewInt(0).Add(last.Volume, trade.Amount)
				last.VolumeByQuote = big.NewInt(0).Add(last.VolumeByQuote, volumeByQuote)
				last.VolumeAsk = big.NewInt(0).Add(last.VolumeAsk, trade.Amount)
				last.VolumeBid = big.NewInt(0).Add(last.VolumeBid, trade.Amount)
				last.VolumeAskByQuote = big.NewInt(0).Add(last.VolumeAskByQuote, volumeByQuote)
				last.VolumeBidByQuote = big.NewInt(0).Add(last.VolumeBidByQuote, volumeByQuote)

			} else {
				userTrade := &types.UserTrade{
					UserAddress:      trade.Maker,
					Count:            big.NewInt(1),
					Volume:           utils.CloneBigInt(trade.Amount),
					VolumeByQuote:    utils.CloneBigInt(volumeByQuote),
					VolumeAskByQuote: utils.CloneBigInt(volumeByQuote),
					VolumeAsk:        utils.CloneBigInt(trade.Amount),
					VolumeBid:        utils.CloneBigInt(trade.Amount),
					VolumeBidByQuote: utils.CloneBigInt(volumeByQuote),
					TimeStamp:        modTime,
					RelayerAddress:   addr,
					BaseToken:        trade.BaseToken,
					QuoteToken:       trade.QuoteToken,
				}
				s.tradeCache.relayerUserTrades[addr][key][trade.Taker][modTime] = userTrade
			}
		} else {

			if last, ok2 := s.tradeCache.relayerUserTrades[addr][key][trade.Taker][modTime]; ok2 {
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
				last.Volume = big.NewInt(0).Add(last.Volume, trade.Amount)
				last.VolumeByQuote = big.NewInt(0).Add(last.VolumeByQuote, volumeByQuote)
				if trade.TakerOrderSide == sideBuy {
					last.VolumeAsk = big.NewInt(0).Add(last.VolumeAsk, trade.Amount)
					last.VolumeAskByQuote = big.NewInt(0).Add(last.VolumeAskByQuote, volumeByQuote)
				} else {
					last.VolumeBid = big.NewInt(0).Add(last.VolumeBid, trade.Amount)
					last.VolumeBidByQuote = big.NewInt(0).Add(last.VolumeBidByQuote, volumeByQuote)
				}

			} else {
				userTrade := &types.UserTrade{
					UserAddress:    trade.Taker,
					Count:          big.NewInt(1),
					Volume:         utils.CloneBigInt(trade.Amount),
					VolumeByQuote:  utils.CloneBigInt(volumeByQuote),
					BaseToken:      trade.BaseToken,
					QuoteToken:     trade.QuoteToken,
					TimeStamp:      modTime,
					RelayerAddress: addr,
				}
				if trade.TakerOrderSide == sideBuy {
					userTrade.VolumeBid = utils.CloneBigInt(trade.Amount)
					userTrade.VolumeBidByQuote = utils.CloneBigInt(volumeByQuote)
					userTrade.VolumeAsk = big.NewInt(0)
					userTrade.VolumeAskByQuote = big.NewInt(0)
				} else {
					userTrade.VolumeBid = big.NewInt(0)
					userTrade.VolumeBidByQuote = big.NewInt(0)
					userTrade.VolumeAsk = utils.CloneBigInt(trade.Amount)
					userTrade.VolumeAskByQuote = utils.CloneBigInt(volumeByQuote)
				}
				s.tradeCache.relayerUserTrades[addr][key][trade.Taker][modTime] = userTrade
			}

			if last, ok2 := s.tradeCache.relayerUserTrades[addr][key][trade.Maker][modTime]; ok2 {
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
				last.Volume = big.NewInt(0).Add(last.Volume, trade.Amount)
				last.VolumeByQuote = big.NewInt(0).Add(last.VolumeByQuote, volumeByQuote)
				if trade.TakerOrderSide == sideSell {
					last.VolumeAsk = big.NewInt(0).Add(last.VolumeAsk, trade.Amount)
					last.VolumeAskByQuote = big.NewInt(0).Add(last.VolumeAskByQuote, volumeByQuote)
				} else {
					last.VolumeBid = big.NewInt(0).Add(last.VolumeBid, trade.Amount)
					last.VolumeBidByQuote = big.NewInt(0).Add(last.VolumeBidByQuote, volumeByQuote)
				}

			} else {
				userTrade := &types.UserTrade{
					UserAddress:    trade.Maker,
					Count:          big.NewInt(1),
					Volume:         utils.CloneBigInt(trade.Amount),
					VolumeByQuote:  utils.CloneBigInt(volumeByQuote),
					BaseToken:      trade.BaseToken,
					QuoteToken:     trade.QuoteToken,
					TimeStamp:      modTime,
					RelayerAddress: addr,
				}
				if trade.TakerOrderSide == sideSell {
					userTrade.VolumeBid = utils.CloneBigInt(trade.Amount)
					userTrade.VolumeBidByQuote = utils.CloneBigInt(volumeByQuote)
					userTrade.VolumeAsk = big.NewInt(0)
					userTrade.VolumeAskByQuote = big.NewInt(0)
				} else {
					userTrade.VolumeBid = big.NewInt(0)
					userTrade.VolumeBidByQuote = big.NewInt(0)
					userTrade.VolumeAsk = utils.CloneBigInt(trade.Amount)
					userTrade.VolumeAskByQuote = utils.CloneBigInt(volumeByQuote)
				}
				s.tradeCache.relayerUserTrades[addr][key][trade.Maker][modTime] = userTrade
			}

		}
	}
	return nil
}

func (s *TradeStatisticService) addRelayerUserTrade(userTrade *types.UserTrade) {
	key := s.getPairString(userTrade.BaseToken, userTrade.QuoteToken)
	if _, ok := s.tradeCache.relayerUserTrades[userTrade.RelayerAddress]; !ok {
		s.tradeCache.relayerUserTrades[userTrade.RelayerAddress] = make(map[string]map[common.Address]map[int64]*types.UserTrade)
	}
	if _, ok := s.tradeCache.relayerUserTrades[userTrade.RelayerAddress][key]; !ok {
		s.tradeCache.relayerUserTrades[userTrade.RelayerAddress][key] = make(map[common.Address]map[int64]*types.UserTrade)
	}
	if _, ok := s.tradeCache.relayerUserTrades[userTrade.RelayerAddress][key][userTrade.UserAddress]; !ok {
		s.tradeCache.relayerUserTrades[userTrade.RelayerAddress][key][userTrade.UserAddress] = make(map[int64]*types.UserTrade)
	}
	s.tradeCache.relayerUserTrades[userTrade.RelayerAddress][key][userTrade.UserAddress][userTrade.TimeStamp] = userTrade
}

func (s *TradeStatisticService) filterRelayerUserTrade() {

}
