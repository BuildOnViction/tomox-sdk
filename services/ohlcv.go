package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

const (
	intervalMin      = 60 * 60 * 24
	intervalMax      = 5 * 12 * 30 * 24 * 60 * 60 // 5 years
	yesterdaySec     = 24 * 60 * 60
	hourSec          = 60 * 60
	milisecond       = 1000
	baseFiat         = "USDT"
	tomo             = "TOMO"
	cacheTimeLifeMax = 15 * 50
)

type PairCache struct {
	pair     *types.Pair
	timelife int64
}
type TokenCache struct {
	token    *types.Token
	timelife int64
}

type PriceUsdt struct {
	price    *big.Int
	timelife int64
}
type OHLCVService struct {
	tradeDao           interfaces.TradeDao
	pairDao            interfaces.PairDao
	tokenDao           interfaces.TokenDao
	tickCache          *tickCache
	mutex              sync.RWMutex
	tokenCache         map[common.Address]*TokenCache
	pairCacheByAddress map[string]*PairCache
	pairCacheByName    map[string]*PairCache
	priceCacheByUsdt   map[common.Address]*PriceUsdt
}

type timeframe struct {
	FirstTime int64 `json:"firstTime" bson:"firstTime"`
	LastTime  int64 `json:"lastTime" bson:"lastTime"`
}
type timeframes []*timeframe

type tickCache struct {
	tframes timeframes
	ticks   map[string]map[int64]*types.Tick

	// relayerAddress => pairAddress => time => tick
	relayerTicks map[common.Address]map[string]map[int64]*types.Tick
}

type tickfile struct {
	Frame        timeframes         `json:"frame" bson:"frame"`
	Ticks        types.Ticks        `json:"ticks" bson:"ticks"`
	RelayerTicks types.RelayerTicks `json:"relayerticks" bson:"relayerticks"`
}
type durationtick struct {
	duration int64
	unit     string
	interval int64
}

var fiatToken *types.Token

// NewOHLCVService init new ohlcv service
func NewOHLCVService(TradeDao interfaces.TradeDao, pairDao interfaces.PairDao, tokenDao interfaces.TokenDao) *OHLCVService {
	fiatToken = new(types.Token)
	f, _ := tokenDao.GetBySymbol(baseFiat)

	if f != nil {
		fiatToken.Decimals = f.Decimals
	} else {
		fiatToken.Decimals = 6
	}

	cache := &tickCache{
		ticks:        make(map[string]map[int64]*types.Tick),
		relayerTicks: make(map[common.Address]map[string]map[int64]*types.Tick),
	}
	return &OHLCVService{
		tradeDao:           TradeDao,
		pairDao:            pairDao,
		tokenDao:           tokenDao,
		tickCache:          cache,
		tokenCache:         make(map[common.Address]*TokenCache),
		pairCacheByAddress: make(map[string]*PairCache),
		pairCacheByName:    make(map[string]*PairCache),
		priceCacheByUsdt:   make(map[common.Address]*PriceUsdt),
	}
}

// Unsubscribe handles all the unsubscription messages for ticks corresponding to a pair
func (s *OHLCVService) Unsubscribe(conn *ws.Client) {
	ws.GetOHLCVSocket().Unsubscribe(conn)
}

// UnsubscribeChannel handles all the unsubscription messages for ticks corresponding to a pair
func (s *OHLCVService) UnsubscribeChannel(conn *ws.Client, p *types.SubscriptionPayload) {
	id := utils.GetOHLCVChannelID(p.BaseToken, p.QuoteToken, p.Units, p.Duration)
	ws.GetOHLCVSocket().UnsubscribeChannel(id, conn)
}

// Subscribe handles all the subscription messages for ticks corresponding to a pair
// It calls the corresponding channel's subscription method and sends trade history back on the connection
func (s *OHLCVService) Subscribe(conn *ws.Client, p *types.SubscriptionPayload) {
	socket := ws.GetOHLCVSocket()

	ohlcv, err := s.GetOHLCV(
		[]types.PairAddresses{{BaseToken: p.BaseToken, QuoteToken: p.QuoteToken}},
		p.Duration,
		p.Units,
		p.From,
		p.To,
	)

	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	id := utils.GetOHLCVChannelID(p.BaseToken, p.QuoteToken, p.Units, p.Duration)
	err = socket.Subscribe(id, conn)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(conn, ohlcv)
}

func (s *OHLCVService) getConfig() []durationtick {
	return []durationtick{
		{
			duration: 1,
			unit:     "min",
			interval: intervalMax,
		},
		{
			duration: 5,
			unit:     "min",
			interval: intervalMax,
		},
		{
			duration: 15,
			unit:     "min",
			interval: intervalMax,
		},
		{
			duration: 30,
			unit:     "min",
			interval: intervalMax,
		},
		{
			duration: 1,
			unit:     "hour",
			interval: intervalMax,
		},
		{
			duration: 2,
			unit:     "hour",
			interval: intervalMax,
		},
		{
			duration: 4,
			unit:     "hour",
			interval: intervalMax,
		},
		{
			duration: 12,
			unit:     "hour",
			interval: intervalMax,
		},
		{
			duration: 1,
			unit:     "day",
			interval: intervalMax,
		},
		{
			duration: 1,
			unit:     "week",
			interval: intervalMax,
		},
		{
			duration: 1,
			unit:     "month",
			interval: intervalMax,
		},
		{
			duration: 3,
			unit:     "month",
			interval: intervalMax,
		},
		{
			duration: 6,
			unit:     "month",
			interval: intervalMax,
		},
		{
			duration: 9,
			unit:     "month",
			interval: intervalMax,
		},
		{
			duration: 1,
			unit:     "year",
			interval: intervalMax,
		},
	}
}

// Init init cache
// ensure add current time frame before trade notify come
func (s *OHLCVService) Init() {
	logger.Info("OHLCV init starting...")
	now := time.Now().Unix()
	datefrom := now - intervalMin
	s.loadCache()
	lastFrame := s.lastTimeFrame()
	if lastFrame != nil {
		logger.Info("last frame first time", time.Unix(lastFrame.FirstTime, 0))
		if now-lastFrame.LastTime < intervalMin {
			datefrom = lastFrame.LastTime
		}
	} else {
		// add start frame to list
		s.tickCache.tframes = append(s.tickCache.tframes, &timeframe{
			FirstTime: now - intervalMax,
			LastTime:  now - intervalMax,
		})
	}
	// add current frame to list
	s.tickCache.tframes = append(s.tickCache.tframes, &timeframe{
		FirstTime: now,
		LastTime:  now,
	})

	lastFrame = s.lastTimeFrame()
	logger.Info("init fetch", time.Unix(datefrom, 0), time.Unix(now, 0))
	s.fetch(datefrom, now, lastFrame)
	s.commitCache()
	go s.continueCache()
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

	logger.Info("OHLCV finished")
}

func (s *OHLCVService) getIntervelByUint(d int64, unit string) (int64, error) {
	durations := s.getConfig()
	for _, duration := range durations {
		if duration.duration == d && duration.unit == unit {
			return duration.interval, nil
		}
	}
	return 0, errors.New("unit not found")
}

// cache need to be locked
func (s *OHLCVService) truncate() {
	now := time.Now().Unix()
	for key, tickby := range s.tickCache.ticks {
		_, _, d, unit, err := s.parseTickKey(key)
		if err == nil {
			interval, e := s.getIntervelByUint(d, unit)
			if e == nil {
				for timeby := range tickby {
					if timeby < now-interval {
						delete(tickby, timeby)
					}
				}
			}
		}
	}
}
func (s *OHLCVService) fetch(fromdate int64, todate int64, frame *timeframe) {
	durations := s.getConfig()
	pageOffset := 0
	size := 1000
	now := time.Now().Unix()
	for {
		trades, err := s.tradeDao.GetTradeByTime(fromdate, todate, pageOffset*size, size)
		logger.Debug("FETCH DATA", pageOffset*size)
		if err != nil || len(trades) == 0 {
			break
		}
		sort.Slice(trades, func(i, j int) bool {
			return trades[i].CreatedAt.Unix() < trades[j].CreatedAt.Unix()
		})
		s.mutex.Lock()
		for i, trade := range trades {
			for _, d := range durations {
				key := s.getTickKey(trade.BaseToken, trade.QuoteToken, d.duration, d.unit)
				if trade.CreatedAt.Unix() > now-d.interval {
					s.updateTick(key, trade)
				}
			}

			if trade.MakerExchange.Hex() == trade.TakerExchange.Hex() {
				s.updateRelayerTick(trade.MakerExchange, s.getTickKey(trade.BaseToken, trade.QuoteToken, 1, "hour"), trade)
			} else {
				s.updateRelayerTick(trade.MakerExchange, s.getTickKey(trade.BaseToken, trade.QuoteToken, 1, "hour"), trade)
				s.updateRelayerTick(trade.Taker, s.getTickKey(trade.BaseToken, trade.QuoteToken, 1, "hour"), trade)
			}

			if i == 0 {
				s.updatefisttimeframe(trade.CreatedAt.Unix(), frame)
			}

		}
		s.mutex.Unlock()
		pageOffset = pageOffset + 1
	}
}

// ensure init cache finished before invoke
func (s *OHLCVService) continueCache() {
	if len(s.tickCache.tframes) > 1 {
		for i := len(s.tickCache.tframes) - 1; i > 0; i-- {
			currentframe := s.tickCache.tframes[i]
			preframe := s.tickCache.tframes[i-1]
			if currentframe.FirstTime > preframe.LastTime {
				logger.Debug("continue cache", time.Unix(preframe.LastTime, 0), time.Unix(currentframe.FirstTime, 0))
				s.fetch(preframe.LastTime, currentframe.FirstTime, currentframe)
			}

		}
	}
	logger.Debug("continueCache finished")
}
func (s *OHLCVService) lastTimeFrame() *timeframe {
	if len(s.tickCache.tframes) > 0 {
		return s.tickCache.tframes[len(s.tickCache.tframes)-1]
	}
	return nil
}
func (s *OHLCVService) updatefisttimeframe(firsttime int64, frame *timeframe) {
	logger.Info("updatefisttimeframe", time.Unix(firsttime, 0))
	if frame != nil {
		frame.FirstTime = firsttime
	}
}

func (s *OHLCVService) updatelasttimeframe(lasttime int64, frame *timeframe) {
	if frame != nil {
		frame.LastTime = lasttime
	}

}

func (s *OHLCVService) flatten() []*types.Tick {
	var ticks []*types.Tick
	for _, tickbytime := range s.tickCache.ticks {
		for _, tick := range tickbytime {
			ticks = append(ticks, tick)
		}
	}
	return ticks
}

func (s *OHLCVService) flattenRelayerTick() []*types.RelayerTick {
	var relayerTicks []*types.RelayerTick
	for addr, tickbyrekayer := range s.tickCache.relayerTicks {
		for _, tickbytime := range tickbyrekayer {
			for _, tick := range tickbytime {
				relayerTicks = append(relayerTicks, &types.RelayerTick{
					RelayerAddress: addr,
					Tick:           tick,
				})
			}
		}
	}
	return relayerTicks
}

func (s *OHLCVService) commitCache() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	logger.Info("commit ohlcv cache")
	s.truncate()
	ticks := s.flatten()
	tickbyrelayer := s.flattenRelayerTick()
	tickfile := &tickfile{
		Frame:        s.tickCache.tframes,
		Ticks:        ticks,
		RelayerTicks: tickbyrelayer,
	}
	tickData, err := json.Marshal(tickfile)
	if err != nil {
		return err
	}
	file, err := os.Create("ohlcv.cache")
	defer file.Close()
	if err == nil {
		_, err = file.Write(tickData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *OHLCVService) loadCache() error {
	file, err := os.Open("ohlcv.cache")
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
	var tickf tickfile
	err = json.Unmarshal(bytes, &tickf)
	if err != nil {
		return err
	}
	for _, t := range tickf.Ticks {
		s.addTick(t)
	}
	for _, t := range tickf.RelayerTicks {
		s.addRelayerTick(t)
	}
	s.tickCache.tframes = tickf.Frame
	return nil
}

func (s *OHLCVService) getTickKey(baseToken, quoteToken common.Address, duration int64, unit string) string {
	return fmt.Sprintf("%s::%s::%s::%s", baseToken.Hex(), quoteToken.Hex(), strconv.FormatInt(duration, 10), unit)
}

func (s *OHLCVService) parseTickKey(key string) (common.Address, common.Address, int64, string, error) {
	keys := strings.Split(key, "::")
	if len(keys) != 4 {
		return common.Address{}, common.Address{}, 0, "", errors.New("invalid key")
	}
	baseToken := common.HexToAddress(keys[0])
	quoteToken := common.HexToAddress(keys[1])
	duration, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		return common.Address{}, common.Address{}, 0, "", errors.New("invalid key")
	}
	unit := keys[3]
	return baseToken, quoteToken, duration, unit, nil
}

func (s *OHLCVService) getVolumeByQuote(baseToken, quoteToken common.Address, amount *big.Int, price *big.Int) *big.Int {
	token, err := s.getTokenByAddress(baseToken)
	if err == nil && token != nil {
		baseTokenDecimalBig := big.NewInt(int64(math.Pow10(token.Decimals)))
		p := new(big.Int).Mul(amount, price)
		return new(big.Int).Div(p, baseTokenDecimalBig)
	}
	return big.NewInt(0)
}

func (s *OHLCVService) getVolumeByUsdt(baseToken, quoteToken common.Address, baseVolume *big.Int, quoteVolume *big.Int) *big.Int {
	token, err := s.getTokenByAddress(quoteToken)
	if err == nil {
		quoteTokenPrice, err := s.getTokenPriceByUsdt(quoteToken)
		if err == nil {
			quoteTokenDecimal := big.NewInt(int64(math.Pow10(token.Decimals)))
			volumeByUsdt := new(big.Int).Mul(quoteVolume, quoteTokenPrice)
			volumeByUsdt.Div(volumeByUsdt, quoteTokenDecimal)
			return volumeByUsdt
		}
	}
	return big.NewInt(0)
}

// updateTick update lastest tick, need to be lock
func (s *OHLCVService) updateTick(key string, trade *types.Trade) error {
	tradeTime := trade.CreatedAt.Unix()
	baseToken, quoteToken, duration, unit, err := s.parseTickKey(key)
	if err != nil {
		return err
	}
	if baseToken.Hex() == trade.BaseToken.Hex() && quoteToken.Hex() == trade.QuoteToken.Hex() {
		modTime, _ := utils.GetModTime(tradeTime, duration, unit)
		if _, ok := s.tickCache.ticks[key]; !ok {
			s.tickCache.ticks[key] = make(map[int64]*types.Tick)
		}
		if tickByTime, ok1 := s.tickCache.ticks[key]; ok1 {
			if last, ok2 := tickByTime[modTime]; ok2 {
				last.Timestamp = modTime
				last.Close = trade.PricePoint
				if last.High.Cmp(trade.PricePoint) < 0 {
					last.High = trade.PricePoint
				}
				if last.Low.Cmp(trade.PricePoint) > 0 {
					last.Low = trade.PricePoint
				}
				last.Volume = big.NewInt(0).Add(last.Volume, trade.Amount)
				volumeByQuote := s.getVolumeByQuote(trade.BaseToken, trade.QuoteToken, trade.Amount, trade.PricePoint)
				last.VolumeByQuote = big.NewInt(0).Add(last.VolumeByQuote, volumeByQuote)
				last.VolumeUsdt = big.NewInt(0).Add(last.VolumeUsdt, s.getVolumeByUsdt(trade.BaseToken, trade.QuoteToken, trade.Amount, volumeByQuote))
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
				last.CloseTime = trade.CreatedAt
			} else {
				volumeByQuote := s.getVolumeByQuote(trade.BaseToken, trade.QuoteToken, trade.Amount, trade.PricePoint)
				tick := &types.Tick{
					Pair: types.PairID{
						PairName:   trade.PairName,
						BaseToken:  trade.BaseToken,
						QuoteToken: trade.QuoteToken,
					},
					OpenTime:      trade.CreatedAt,
					Open:          trade.PricePoint,
					Close:         trade.PricePoint,
					High:          trade.PricePoint,
					Low:           trade.PricePoint,
					Volume:        trade.Amount,
					VolumeByQuote: volumeByQuote,
					VolumeUsdt:    s.getVolumeByUsdt(trade.BaseToken, trade.QuoteToken, trade.Amount, volumeByQuote),
					Count:         big.NewInt(1),
					Timestamp:     modTime,
					Duration:      duration,
					Unit:          unit,
				}
				tickByTime[modTime] = tick
			}
		}
	}

	return nil
}

// updateRelayerTick update lastest tick, need to be lock
func (s *OHLCVService) updateRelayerTick(relayerAddress common.Address, key string, trade *types.Trade) error {
	tradeTime := trade.CreatedAt.Unix()
	baseToken, quoteToken, duration, unit, err := s.parseTickKey(key)
	if err != nil {
		return err
	}
	if baseToken.Hex() == trade.BaseToken.Hex() && quoteToken.Hex() == trade.QuoteToken.Hex() {
		modTime, _ := utils.GetModTime(tradeTime, duration, unit)
		if _, ok := s.tickCache.relayerTicks[relayerAddress]; !ok {
			s.tickCache.relayerTicks[relayerAddress] = make(map[string]map[int64]*types.Tick)
		}
		if _, ok := s.tickCache.relayerTicks[relayerAddress][key]; !ok {
			s.tickCache.relayerTicks[relayerAddress][key] = make(map[int64]*types.Tick)
		}

		if tickByTime, ok1 := s.tickCache.relayerTicks[relayerAddress][key]; ok1 {
			if last, ok2 := tickByTime[modTime]; ok2 {
				last.Timestamp = modTime
				last.Close = trade.PricePoint
				if last.High.Cmp(trade.PricePoint) < 0 {
					last.High = trade.PricePoint
				}
				if last.Low.Cmp(trade.PricePoint) > 0 {
					last.Low = trade.PricePoint
				}
				last.Volume = big.NewInt(0).Add(last.Volume, trade.Amount)
				volumeByQuote := s.getVolumeByQuote(trade.BaseToken, trade.QuoteToken, trade.Amount, trade.PricePoint)
				last.VolumeByQuote = big.NewInt(0).Add(last.VolumeByQuote, volumeByQuote)
				last.VolumeUsdt = big.NewInt(0).Add(last.VolumeUsdt, s.getVolumeByUsdt(trade.BaseToken, trade.QuoteToken, trade.Amount, volumeByQuote))
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
				last.CloseTime = trade.CreatedAt
			} else {
				volumeByQuote := s.getVolumeByQuote(trade.BaseToken, trade.QuoteToken, trade.Amount, trade.PricePoint)
				tick := &types.Tick{
					Pair: types.PairID{
						PairName:   trade.PairName,
						BaseToken:  trade.BaseToken,
						QuoteToken: trade.QuoteToken,
					},
					OpenTime:      trade.CreatedAt,
					Open:          trade.PricePoint,
					Close:         trade.PricePoint,
					High:          trade.PricePoint,
					Low:           trade.PricePoint,
					Volume:        trade.Amount,
					VolumeByQuote: volumeByQuote,
					VolumeUsdt:    s.getVolumeByUsdt(trade.BaseToken, trade.QuoteToken, trade.Amount, volumeByQuote),
					Count:         big.NewInt(1),
					Timestamp:     modTime,
					Duration:      duration,
					Unit:          unit,
				}
				tickByTime[modTime] = tick
			}
		}
	}

	return nil
}
func (s *OHLCVService) addTick(tick *types.Tick) {
	tick.VolumeByQuote = big.NewInt(0)
	tick.VolumeUsdt = big.NewInt(0)
	key := s.getTickKey(tick.Pair.BaseToken, tick.Pair.QuoteToken, tick.Duration, tick.Unit)
	if _, ok := s.tickCache.ticks[key]; ok {

		s.tickCache.ticks[key][tick.Timestamp] = tick
	} else {
		s.tickCache.ticks[key] = make(map[int64]*types.Tick)
		s.tickCache.ticks[key][tick.Timestamp] = tick
	}

}

func (s *OHLCVService) addRelayerTick(relayerTick *types.RelayerTick) {
	relayerTick.Tick.VolumeByQuote = big.NewInt(0)
	relayerTick.Tick.VolumeUsdt = big.NewInt(0)
	key := s.getTickKey(relayerTick.Tick.Pair.BaseToken, relayerTick.Tick.Pair.QuoteToken, relayerTick.Tick.Duration, relayerTick.Tick.Unit)
	if _, ok := s.tickCache.relayerTicks[relayerTick.RelayerAddress]; ok {
		if _, ok := s.tickCache.relayerTicks[relayerTick.RelayerAddress][key]; ok {
			s.tickCache.relayerTicks[relayerTick.RelayerAddress][key][relayerTick.Tick.Timestamp] = relayerTick.Tick
		} else {
			s.tickCache.relayerTicks[relayerTick.RelayerAddress][key] = make(map[int64]*types.Tick)
			s.tickCache.relayerTicks[relayerTick.RelayerAddress][key][relayerTick.Tick.Timestamp] = relayerTick.Tick
		}
	} else {
		s.tickCache.relayerTicks[relayerTick.RelayerAddress] = make(map[string]map[int64]*types.Tick)
		if _, ok := s.tickCache.relayerTicks[relayerTick.RelayerAddress][key]; ok {
			s.tickCache.relayerTicks[relayerTick.RelayerAddress][key][relayerTick.Tick.Timestamp] = relayerTick.Tick
		} else {
			s.tickCache.relayerTicks[relayerTick.RelayerAddress][key] = make(map[int64]*types.Tick)
			s.tickCache.relayerTicks[relayerTick.RelayerAddress][key][relayerTick.Tick.Timestamp] = relayerTick.Tick
		}
	}

}

func (s *OHLCVService) filterTick(key string, start, end int64) []*types.Tick {
	var res []*types.Tick
	if _, ok := s.tickCache.ticks[key]; ok {
		for _, t := range s.tickCache.ticks[key] {
			if (t.Timestamp >= start || start == 0) && (t.Timestamp <= end || end == 0) {
				c := *t
				c.Timestamp = t.Timestamp * 1000
				res = append(res, &c)
			}
		}
	} else {
		return nil
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Timestamp < res[j].Timestamp
	})
	return res
}

func (s *OHLCVService) filterRelayerTick(relayerAddress common.Address, key string, start, end int64) []*types.Tick {
	var res []*types.Tick
	if _, ok := s.tickCache.relayerTicks[relayerAddress]; ok {
		if _, ok := s.tickCache.relayerTicks[relayerAddress][key]; ok {
			for _, t := range s.tickCache.relayerTicks[relayerAddress][key] {
				if (t.Timestamp >= start || start == 0) && (t.Timestamp <= end || end == 0) {
					c := *t
					c.Timestamp = t.Timestamp * 1000
					res = append(res, &c)
				}
			}
		} else {
			return nil
		}
	} else {
		return nil
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Timestamp < res[j].Timestamp
	})
	return res
}

// Get24hTick get 24h tick of token
func (s *OHLCVService) Get24hTick(baseToken, quoteToken common.Address) *types.Tick {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	return s.get24hTick(baseToken, quoteToken)
}
func (s *OHLCVService) get24hTick(baseToken, quoteToken common.Address) *types.Tick {
	var res []*types.Tick
	now := time.Now()
	begin := now.AddDate(0, 0, -1).Unix()
	key := s.getTickKey(baseToken, quoteToken, 1, "min")
	res = s.filterTick(key, begin, 0)

	if len(res) >= 1 {
		first := res[0]
		last := res[len(res)-1]
		high := first.High
		low := first.Low
		volume := big.NewInt(0)
		volumebyquote := big.NewInt(0)
		volumebyUsdt := big.NewInt(0)
		count := big.NewInt(0)
		for _, t := range res {
			if high.Cmp(t.High) < 0 {
				high = t.High
			}
			if low.Cmp(t.Low) > 0 {
				low = t.Low
			}
			volume = volume.Add(volume, t.Volume)
			volumebyquote = volumebyquote.Add(volumebyquote, t.VolumeByQuote)
			volumebyUsdt = volumebyUsdt.Add(volumebyUsdt, t.VolumeUsdt)
			count = count.Add(count, t.Count)
		}
		return &types.Tick{
			Open:          first.Open,
			Close:         last.Close,
			High:          high,
			Low:           low,
			CloseTime:     last.CloseTime,
			Count:         count,
			Volume:        volume,
			VolumeByQuote: volumebyquote,
			VolumeUsdt:    volumebyUsdt,
			Timestamp:     last.Timestamp,
		}
	}
	return nil
}

func (s *OHLCVService) get24hRelayerTick(relayerAddress common.Address, baseToken, quoteToken common.Address) *types.Tick {
	var res []*types.Tick
	now := time.Now()
	begin := now.AddDate(0, 0, -1).Unix()
	key := s.getTickKey(baseToken, quoteToken, 1, "hour")
	res = s.filterRelayerTick(relayerAddress, key, begin, 0)

	if len(res) >= 1 {
		first := res[0]
		last := res[len(res)-1]
		high := first.High
		low := first.Low
		volume := big.NewInt(0)
		volumebyquote := big.NewInt(0)
		volumebyUsdt := big.NewInt(0)
		count := big.NewInt(0)
		for _, t := range res {
			if high.Cmp(t.High) < 0 {
				high = t.High
			}
			if low.Cmp(t.Low) > 0 {
				low = t.Low
			}
			volume = volume.Add(volume, t.Volume)
			volumebyquote = volumebyquote.Add(volumebyquote, t.VolumeByQuote)
			volumebyUsdt = volumebyUsdt.Add(volumebyUsdt, t.VolumeUsdt)
			count = count.Add(count, t.Count)
		}
		return &types.Tick{
			Open:          first.Open,
			Close:         last.Close,
			High:          high,
			Low:           low,
			CloseTime:     last.CloseTime,
			Count:         count,
			Volume:        volume,
			VolumeByQuote: volumebyquote,
			VolumeUsdt:    volumebyUsdt,
			Timestamp:     last.Timestamp,
		}
	}
	return nil
}

// NotifyTrade trigger if trade comming
func (s *OHLCVService) NotifyTrade(trade *types.Trade) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, d := range s.getConfig() {
		key := s.getTickKey(trade.BaseToken, trade.QuoteToken, d.duration, d.unit)
		s.updateTick(key, trade)
	}
	if trade.MakerExchange.Hex() == trade.TakerExchange.Hex() {
		s.updateRelayerTick(trade.MakerExchange, s.getTickKey(trade.BaseToken, trade.QuoteToken, 1, "hour"), trade)
	} else {
		s.updateRelayerTick(trade.MakerExchange, s.getTickKey(trade.BaseToken, trade.QuoteToken, 1, "hour"), trade)
		s.updateRelayerTick(trade.Taker, s.getTickKey(trade.BaseToken, trade.QuoteToken, 1, "hour"), trade)
	}
	lastFrame := s.lastTimeFrame()
	s.updatelasttimeframe(trade.CreatedAt.Unix(), lastFrame)
}

func (s *OHLCVService) getOHLCV(pairs []types.PairAddresses, duration int64, unit string, start, end time.Time) ([]*types.Tick, error) {
	res := make([]*types.Tick, 0)
	match := make(bson.M)
	match = getMatchQuery(start, end, pairs...)
	match = bson.M{"$match": match}

	addFields := make(bson.M)
	group, addFields := getGroupAddFieldBson("$createdAt", unit, duration)
	group = bson.M{"$group": group}

	sort := bson.M{"$sort": bson.M{"timestamp": 1}}

	query := []bson.M{match, group, addFields, sort}

	res, err := s.tradeDao.Aggregate(query)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return []*types.Tick{}, nil
	}

	return res, nil
}

// GetOHLCV fetches OHLCV data using
// pairName: can be "" for fetching data for all pairs
// duration: in integer
// unit: sec,min,hour,day,week,month,yr
// timeInterval: 0-2 entries (0 argument: latest data,1st argument: from timestamp, 2nd argument: to timestamp)
func (s *OHLCVService) GetOHLCV(pairs []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	currentTimestamp := time.Now().Unix()
	modTime, intervalInSeconds := utils.GetModTime(currentTimestamp, duration, unit)
	start := time.Unix(modTime-intervalInSeconds, 0)
	end := time.Unix(currentTimestamp, 0)
	if len(timeInterval) >= 1 {
		end = time.Unix(timeInterval[1], 0)
		start = time.Unix(timeInterval[0], 0)
	}
	p := pairs[0]
	cacheKey := s.getTickKey(p.BaseToken, p.QuoteToken, duration, unit)
	ticks := s.filterTick(cacheKey, start.Unix(), end.Unix())
	if ticks == nil {
		return s.getOHLCV(pairs, duration, unit, start, end)
	}
	return ticks, nil
}

func getMatchQuery(start, end time.Time, pairs ...types.PairAddresses) bson.M {
	match := bson.M{
		"createdAt": bson.M{
			"$gte": start,
			"$lt":  end,
		},
		"status": bson.M{"$in": []string{types.SUCCESS}},
	}

	if len(pairs) >= 1 {
		or := make([]bson.M, 0)

		for _, pair := range pairs {
			or = append(or, bson.M{
				"$and": []bson.M{
					{
						"baseToken":  pair.BaseToken.Hex(),
						"quoteToken": pair.QuoteToken.Hex(),
					},
				},
			},
			)
		}

		match["$or"] = or
	}

	return match
}

// query for grouping of the documents and addition of required fields using aggregate pipeline
func getGroupAddFieldBson(key, units string, duration int64) (bson.M, bson.M) {
	var group, addFields bson.M

	t := time.Unix(0, 0)
	var date interface{}
	if key == "now" {
		date = time.Now()
	} else {
		date = key
	}

	one, _ := bson.ParseDecimal128("1")
	group = bson.M{
		"count":     bson.M{"$sum": one},
		"high":      bson.M{"$max": "$pricepoint"},
		"low":       bson.M{"$min": "$pricepoint"},
		"open":      bson.M{"$first": "$pricepoint"},
		"close":     bson.M{"$last": "$pricepoint"},
		"volume":    bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
		"openTime":  bson.M{"$first": "$createdAt"},
		"closeTime": bson.M{"$last": "$createdAt"},
	}

	groupID := make(bson.M)
	switch units {
	case "sec":
		groupID = bson.M{
			"year":   bson.M{"$year": date},
			"day":    bson.M{"$dayOfMonth": date},
			"month":  bson.M{"$month": date},
			"hour":   bson.M{"$hour": date},
			"minute": bson.M{"$minute": date},
			"second": bson.M{
				"$subtract": []interface{}{
					bson.M{"$second": date},
					bson.M{"$mod": []interface{}{bson.M{"$second": date}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{
			"timestamp": bson.M{
				"$subtract": []interface{}{bson.M{
					"$dateFromParts": bson.M{
						"year":   "$_id.year",
						"month":  "$_id.month",
						"day":    "$_id.day",
						"hour":   "$_id.hour",
						"minute": "$_id.minute",
						"second": "$_id.second"}}, t}}}}

	case "min":
		groupID = bson.M{
			"year":  bson.M{"$year": date},
			"day":   bson.M{"$dayOfMonth": date},
			"month": bson.M{"$month": date},
			"hour":  bson.M{"$hour": date},
			"minute": bson.M{
				"$subtract": []interface{}{
					bson.M{"$minute": date},
					bson.M{"$mod": []interface{}{bson.M{"$minute": date}, duration}},
				}}}

		addFields = bson.M{"$addFields": bson.M{"timestamp": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":   "$_id.year",
			"month":  "$_id.month",
			"day":    "$_id.day",
			"hour":   "$_id.hour",
			"minute": "$_id.minute",
		}}, t}}}}

	case "hour":
		groupID = bson.M{
			"year":  bson.M{"$year": date},
			"day":   bson.M{"$dayOfMonth": date},
			"month": bson.M{"$month": date},
			"hour": bson.M{
				"$subtract": []interface{}{
					bson.M{"$hour": date},
					bson.M{"$mod": []interface{}{bson.M{"$hour": date}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"timestamp": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
			"day":   "$_id.day",
			"hour":  "$_id.hour",
		}}, t}}}}

	case "day":
		groupID = bson.M{
			"year":  bson.M{"$year": date},
			"month": bson.M{"$month": date},
			"day": bson.M{
				"$subtract": []interface{}{
					bson.M{"$dayOfMonth": date},
					bson.M{"$mod": []interface{}{bson.M{"$dayOfMonth": date}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"timestamp": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
			"day":   "$_id.day",
		}}, t}}}}

	case "week":
		groupID = bson.M{
			"year": bson.M{"$isoWeekYear": date},
			"isoWeek": bson.M{
				"$subtract": []interface{}{
					bson.M{"$isoWeek": date},
					bson.M{"$mod": []interface{}{bson.M{"$isoWeek": date}, duration}}}}}

		addFields = bson.M{"$addFields": bson.M{"timestamp": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"isoWeekYear": "$_id.year",
			"isoWeek":     "$_id.isoWeek",
		}}, t}}}}

	case "month":
		groupID = bson.M{
			"year": bson.M{"$year": date},
			"month": bson.M{
				"$subtract": []interface{}{
					bson.M{
						"$multiply": []interface{}{
							bson.M{"$ceil": bson.M{"$divide": []interface{}{
								bson.M{"$month": date},
								duration}},
							},
							duration},
					}, duration - 1}}}

		addFields = bson.M{"$addFields": bson.M{"timestamp": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year":  "$_id.year",
			"month": "$_id.month",
		}}, t}}}}

	case "year":
		groupID = bson.M{
			"year": bson.M{
				"$subtract": []interface{}{
					bson.M{"$year": date},
					bson.M{"$mod": []interface{}{bson.M{"$year": date}, duration}},
				},
			},
		}

		addFields = bson.M{"$addFields": bson.M{"timestamp": bson.M{"$subtract": []interface{}{bson.M{"$dateFromParts": bson.M{
			"year": "$_id.year"}}, t}}}}

	}

	groupID["pairName"] = "$pairName"
	groupID["baseToken"] = "$baseToken"
	groupID["quoteToken"] = "$quoteToken"
	group["_id"] = groupID

	return group, addFields
}

// GetTokenPairData get tick of pair tokens
func (s *OHLCVService) getTokenPairData(pairName string, baseTokenSymbol string, baseToken common.Address, quoteToken common.Address) *types.PairData {
	tick := s.get24hTick(baseToken, quoteToken)
	if tick != nil {
		pairData := &types.PairData{
			Pair:         types.PairID{PairName: pairName, BaseToken: baseToken, QuoteToken: quoteToken},
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
		pairData.Open = tick.Open
		pairData.High = tick.High
		pairData.Low = tick.Low
		pairData.Volume = tick.VolumeByQuote
		pairData.Close = tick.Close
		pairData.Count = tick.Count
		price, err := s.getLastPriceCurrentByTime(baseTokenSymbol, time.Unix(tick.Timestamp/milisecond, 0))
		if err == nil {
			pairData.CloseBaseUsd = price
		}
		return pairData
	}
	return nil
}

// GetTokenPairData get tick of pair tokens
func (s *OHLCVService) GetTokenPairData(baseToken common.Address, quoteToken common.Address) *types.PairData {
	p, err := s.pairDao.GetByTokenAddress(baseToken, quoteToken)
	if err != nil {
		return nil
	}
	return s.getTokenPairData(p.Name(), p.BaseTokenSymbol, p.BaseTokenAddress, p.QuoteTokenAddress)

}

func (s *OHLCVService) GetAllTokenPairDataByCoinbase(addr common.Address) ([]*types.PairData, error) {
	pairs, err := s.pairDao.GetActivePairsByCoinbase(addr)
	if err != nil {
		return nil, err
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	pairsData := make([]*types.PairData, 0)
	for _, p := range pairs {
		pairData := s.getTokenPairData(p.Name(), p.BaseTokenSymbol, p.BaseTokenAddress, p.QuoteTokenAddress)
		if pairData != nil {
			pairsData = append(pairsData, pairData)
		} else {
			emptyPairData := &types.PairData{
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

			pairsData = append(pairsData, emptyPairData)

		}

	}

	return pairsData, nil
}

// GetAllTokenPairData get tick of all tokens
func (s *OHLCVService) GetAllTokenPairData() ([]*types.PairData, error) {
	pairs, err := s.pairDao.GetActivePairs()
	if err != nil {
		return nil, err
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	pairsData := make([]*types.PairData, 0)
	for _, p := range pairs {
		pairData := s.getTokenPairData(p.Name(), p.BaseTokenSymbol, p.BaseTokenAddress, p.QuoteTokenAddress)
		if pairData != nil {
			pairsData = append(pairsData, pairData)
		} else {
			emptyPairData := &types.PairData{
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

			pairsData = append(pairsData, emptyPairData)

		}

	}

	return pairsData, nil
}

// GetPairPrice get lastest price by time
func (s *OHLCVService) GetPairPrice(pairName string, timestamp int64) (int64, error) {
	pair, err := s.pairDao.GetByName(pairName)
	if err == nil && pair != nil {

	}
	return 0, nil
}

//GetFiatPriceChart get fiat chart
func (s *OHLCVService) GetFiatPriceChart() (map[string][]*types.FiatPriceItem, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	symbols := []string{"BTC", "ETH", "BNB", "TOMO"}
	now := time.Now().Unix()
	yesterday := now - yesterdaySec
	res := make(map[string][]*types.FiatPriceItem)
	totalVolume := make(map[string]*big.Int)
	for _, symbol := range symbols {
		var fiats []*types.FiatPriceItem
		totalVolume[symbol] = big.NewInt(0)
		pairs, err := s.pairDao.GetActivePairs()
		if err != nil {
			continue
		}
		for _, pair := range pairs {
			if pair.BaseTokenSymbol == symbol {
				tick := s.get24hTick(pair.BaseTokenAddress, pair.QuoteTokenAddress)
				if tick != nil {
					quoteTokenDecimal := int64(math.Pow10(pair.QuoteTokenDecimals))
					volume := big.NewInt(0).Div(tick.VolumeByQuote, big.NewInt(quoteTokenDecimal))
					totalVolume[symbol].Add(totalVolume[symbol], volume)
				}
			}
		}
		for step := yesterday; step <= now; step = step + hourSec {
			price, err := s.getLastPriceCurrentByTime(symbol, time.Unix(step, 0))
			if err == nil {
				fiat := &types.FiatPriceItem{
					Symbol:       symbol,
					Price:        price.String(),
					Timestamp:    step,
					FiatCurrency: baseFiat,
					TotalVolume:  totalVolume[symbol].String(),
				}
				fiats = append(fiats, fiat)
			}
		}
		sort.Slice(fiats, func(i, j int) bool {
			return fiats[i].Timestamp > fiats[j].Timestamp
		})
		if len(fiats) > 0 {
			res[symbol] = fiats
		}
	}

	return res, nil
}

func (s *OHLCVService) getLastPricePairAtTime(pairName string, createAt time.Time) (*big.Float, error) {
	pairs := strings.Split(pairName, "/")
	if len(pairs) != 2 {
		return nil, errors.New("Invalid pair name")
	}
	if pairs[0] == pairs[1] {
		return big.NewFloat(1), nil
	}
	pair, err := s.getCachePairByName(pairName)
	if err == nil && pair != nil {
		durations := s.getConfig()
		for _, d := range durations {
			mod, _ := utils.GetModTime(createAt.Unix(), d.duration, d.unit)
			key := s.getTickKey(pair.BaseTokenAddress, pair.QuoteTokenAddress, d.duration, d.unit)
			if tradeTick, ok := s.tickCache.ticks[key]; ok {
				if tick, ok := tradeTick[mod]; ok {
					quoteTokenDecimal := int64(math.Pow10(pair.QuoteTokenDecimals))
					return big.NewFloat(0).Quo(new(big.Float).SetInt(tick.Close), big.NewFloat(float64(quoteTokenDecimal))), nil
				}
			}
		}
	}
	return nil, errors.New("Price not found")
}

func (s *OHLCVService) getLastPriceCurrentByTime(symbol string, createAt time.Time) (*big.Float, error) {
	USD := symbol + "/" + baseFiat
	price, err := s.getLastPricePairAtTime(USD, createAt)
	if err != nil {

		var symbolpricebytomo *big.Float
		var err error
		symbolpricebytomo, err = s.getLastPricePairAtTime(symbol+"/"+tomo, createAt)
		if err != nil {
			symbolpricebytomo, err = s.getLastPricePairAtTime(tomo+"/"+symbol, createAt)
			if err != nil {
				return nil, errors.New("Price not found")
			}
			symbolpricebytomo = new(big.Float).Quo(big.NewFloat(1), symbolpricebytomo)
		}
		tomopricebybase, err := s.getLastPricePairAtTime(tomo+"/"+baseFiat, createAt)
		if err != nil {
			return nil, errors.New("Price not found")
		}

		return big.NewFloat(0).Mul(symbolpricebytomo, tomopricebybase), nil
	}
	return price, err
}

// GetLastPriceCurrentByTime get last trade price
func (s *OHLCVService) GetLastPriceCurrentByTime(symbol string, createAt time.Time) (*big.Float, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.getLastPriceCurrentByTime(symbol, createAt)
}

func (s *OHLCVService) getCachePairByName(pairName string) (*types.Pair, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	now := time.Now().Unix()
	if pairCache, ok := s.pairCacheByName[pairName]; ok {
		if now-pairCache.timelife < cacheTimeLifeMax {
			return pairCache.pair, nil
		}
		delete(s.pairCacheByName, pairName)
	}
	pair, err := s.pairDao.GetByName(pairName)
	if err == nil && pair != nil {
		s.pairCacheByName[pairName] = &PairCache{
			pair:     pair,
			timelife: now,
		}
	}
	return pair, err
}

func (s *OHLCVService) getCachePairByAddress(baseToken, quoteToken common.Address) (*types.Pair, error) {
	now := time.Now().Unix()
	pairName := utils.GetPairKey(baseToken, quoteToken)
	if pairCache, ok := s.pairCacheByAddress[pairName]; ok {
		if now-pairCache.timelife < cacheTimeLifeMax {
			return pairCache.pair, nil
		}
		delete(s.pairCacheByAddress, pairName)
	}
	pair, err := s.pairDao.GetByTokenAddress(baseToken, quoteToken)
	if err == nil && pair != nil {
		s.pairCacheByAddress[pairName] = &PairCache{
			pair:     pair,
			timelife: now,
		}
	}
	return pair, err
}

func (s *OHLCVService) getTokenByAddress(token common.Address) (*types.Token, error) {
	now := time.Now().Unix()
	if tokenCache, ok := s.tokenCache[token]; ok {
		if now-tokenCache.timelife < cacheTimeLifeMax {
			return tokenCache.token, nil
		}
		delete(s.tokenCache, token)
	}
	t, err := s.tokenDao.GetByAddress(token)
	if err == nil && t != nil {
		s.tokenCache[token] = &TokenCache{
			token:    t,
			timelife: now,
		}
	}
	return t, err
}

func (s *OHLCVService) getTokenPriceByUsdt(token common.Address) (*big.Int, error) {
	now := time.Now().Unix()
	if tokenPrice, ok := s.priceCacheByUsdt[token]; ok {
		if now-tokenPrice.timelife < cacheTimeLifeMax {
			return tokenPrice.price, nil
		}
		delete(s.priceCacheByUsdt, token)
	}
	t, err := s.getTokenByAddress(token)
	if err != nil || t == nil {
		return big.NewInt(0), errors.New("cant not get token price by usdt")
	}
	fiatTokenDecimal := int64(math.Pow10(fiatToken.Decimals))
	if t.Symbol == baseFiat {
		s.priceCacheByUsdt[token] = &PriceUsdt{
			price:    big.NewInt(fiatTokenDecimal),
			timelife: now,
		}
		return big.NewInt(fiatTokenDecimal), nil
	}
	price, err := s.getLastPriceCurrentByTime(t.Symbol, time.Now())
	if price == nil || err != nil {
		return big.NewInt(0), err
	}
	priceDecimals := new(big.Float).Mul(price, big.NewFloat(float64(fiatTokenDecimal)))
	priceDecimalsInt := new(big.Int)
	priceDecimals.Int(priceDecimalsInt)
	if err == nil {
		s.priceCacheByUsdt[token] = &PriceUsdt{
			price:    priceDecimalsInt,
			timelife: now,
		}
	}
	return priceDecimalsInt, err
}

// GetVolumeByUsdt convert to USDT volume
func (s *OHLCVService) GetVolumeByUsdt(token common.Address, volume *big.Int) *big.Int {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return s.getVolumeByUsdt(common.Address{}, token, new(big.Int), volume)
}

// GetVolumeByCoinbase get total volume exchange
func (s *OHLCVService) GetVolumeByCoinbase(addr common.Address) (*big.Int, error) {
	pairs, err := s.pairDao.GetActivePairsByCoinbase(addr)
	if err != nil {
		return nil, err
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	totalVolume := big.NewInt(0)
	for _, p := range pairs {
		tick := s.get24hRelayerTick(addr, p.BaseTokenAddress, p.QuoteTokenAddress)
		if tick != nil {
			totalVolume = totalVolume.Add(totalVolume, tick.VolumeUsdt)
		}

	}

	return totalVolume, nil
}
