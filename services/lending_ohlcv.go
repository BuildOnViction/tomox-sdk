package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"math/big"
	"os"
	"sort"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

const (
	// LendingCachePath OHL cache file name
	LendingCachePath = "lending.cache"
)

// LendingOhlcvService ohlcv lending struct
type LendingOhlcvService struct {
	lendingTradeDao  interfaces.LendingTradeDao
	lendingTickCache *lendingTickCache
	lendingPairDao   interfaces.LendingPairDao
	mutex            sync.RWMutex
	tokenCache       map[common.Address]int
}

type lendingTickCache struct {
	tframes timeframes
	ticks   map[string]map[int64]*types.LendingTick
}

type lendingtickfile struct {
	Frame        timeframes         `json:"frame" bson:"frame"`
	LendingTicks types.LendingTicks `json:"ticks" bson:"ticks"`
}

// NewLendingOhlcvService init new ohlcv service
func NewLendingOhlcvService(lendingTradeDao interfaces.LendingTradeDao, lendingPairDao interfaces.LendingPairDao) *LendingOhlcvService {
	cache := &lendingTickCache{
		ticks: make(map[string]map[int64]*types.LendingTick),
	}
	return &LendingOhlcvService{
		lendingTradeDao:  lendingTradeDao,
		lendingPairDao:   lendingPairDao,
		lendingTickCache: cache,
		tokenCache:       make(map[common.Address]int),
	}
}

// Unsubscribe handles all the unsubscription messages for ticks corresponding to a pair
func (s *LendingOhlcvService) Unsubscribe(conn *ws.Client) {
	ws.GetLendingOhlcvSocket().Unsubscribe(conn)
}

// UnsubscribeChannel handles all the unsubscription messages for ticks corresponding to a pair
func (s *LendingOhlcvService) UnsubscribeChannel(conn *ws.Client, p *types.SubscriptionPayload) {
	id := utils.GetLendingOhlcvChannelID(p.Term, p.LendingToken, p.Units, p.Duration)
	ws.GetLendingOhlcvSocket().UnsubscribeChannel(id, conn)
}

// Subscribe handles all the subscription messages for ticks corresponding to a pair
// It calls the corresponding channel's subscription method and sends trade history back on the connection
func (s *LendingOhlcvService) Subscribe(conn *ws.Client, p *types.SubscriptionPayload) {
	socket := ws.GetLendingOhlcvSocket()

	ohlcv, err := s.GetOHLCV(
		p.Term,
		p.LendingToken,
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

	id := utils.GetLendingOhlcvChannelID(p.Term, p.LendingToken, p.Units, p.Duration)
	err = socket.Subscribe(id, conn)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(conn, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(conn, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(conn, ohlcv)
}

func (s *LendingOhlcvService) getConfig() []durationtick {
	return []durationtick{
		{
			duration: 1,
			unit:     "min",
			interval: 24 * 60 * 60,
		},
		{
			duration: 5,
			unit:     "min",
			interval: 24 * 60 * 60,
		},
		{
			duration: 15,
			unit:     "min",
			interval: 24 * 60 * 60,
		},
		{
			duration: 30,
			unit:     "min",
			interval: 24 * 60 * 60,
		},
		{
			duration: 1,
			unit:     "hour",
			interval: 7 * 24 * 60 * 60,
		},
		{
			duration: 2,
			unit:     "hour",
			interval: 7 * 24 * 60 * 60,
		},
		{
			duration: 4,
			unit:     "hour",
			interval: 7 * 24 * 60 * 60,
		},
		{
			duration: 12,
			unit:     "hour",
			interval: 7 * 1 * 60 * 60,
		},
		{
			duration: 1,
			unit:     "day",
			interval: 30 * 24 * 60 * 60,
		},
		{
			duration: 1,
			unit:     "week",
			interval: 30 * 24 * 60 * 60,
		},
		{
			duration: 1,
			unit:     "month",
			interval: 12 * 30 * 24 * 60 * 60,
		},
		{
			duration: 3,
			unit:     "month",
			interval: 12 * 30 * 24 * 60 * 60,
		},
		{
			duration: 6,
			unit:     "month",
			interval: 2 * 12 * 30 * 24 * 60 * 60,
		},
		{
			duration: 9,
			unit:     "month",
			interval: 2 * 12 * 30 * 24 * 60 * 60,
		},
		{
			duration: 1,
			unit:     "year",
			interval: 2 * 12 * 30 * 24 * 60 * 60,
		},
	}
}

// Init init cache
// ensure add current time frame before trade notify come
func (s *LendingOhlcvService) Init() {
	logger.Info("Lending OHLCV init starting...")
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
		s.lendingTickCache.tframes = append(s.lendingTickCache.tframes, &timeframe{
			FirstTime: now - intervalMax,
			LastTime:  now - intervalMax,
		})
	}
	// add current frame to list
	s.lendingTickCache.tframes = append(s.lendingTickCache.tframes, &timeframe{
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

func (s *LendingOhlcvService) getIntervelByUint(d int64, unit string) (int64, error) {
	durations := s.getConfig()
	for _, duration := range durations {
		if duration.duration == d && duration.unit == unit {
			return duration.interval, nil
		}
	}
	return 0, errors.New("unit not found")
}

// cache need to be locked
func (s *LendingOhlcvService) truncate() {
	now := time.Now().Unix()
	for key, tickby := range s.lendingTickCache.ticks {
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
func (s *LendingOhlcvService) fetch(fromdate int64, todate int64, frame *timeframe) {
	durations := s.getConfig()
	pageOffset := 0
	size := 1000
	now := time.Now().Unix()
	for {
		trades, err := s.lendingTradeDao.GetLendingTradeByTime(fromdate, todate, pageOffset*size, size)
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
				key := s.getTickKey(trade.Term, trade.LendingToken, d.duration, d.unit)
				if trade.CreatedAt.Unix() > now-d.interval {
					s.updateTick(key, trade)
				}
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
func (s *LendingOhlcvService) continueCache() {
	if len(s.lendingTickCache.tframes) > 1 {
		for i := len(s.lendingTickCache.tframes) - 1; i > 0; i-- {
			currentframe := s.lendingTickCache.tframes[i]
			preframe := s.lendingTickCache.tframes[i-1]
			if currentframe.FirstTime > preframe.LastTime {
				logger.Debug("continue cache", time.Unix(preframe.LastTime, 0), time.Unix(currentframe.FirstTime, 0))
				s.fetch(preframe.LastTime, currentframe.FirstTime, currentframe)
			}

		}
	}
	logger.Debug("continueCache finished")
}
func (s *LendingOhlcvService) lastTimeFrame() *timeframe {
	if len(s.lendingTickCache.tframes) > 0 {
		return s.lendingTickCache.tframes[len(s.lendingTickCache.tframes)-1]
	}
	return nil
}
func (s *LendingOhlcvService) updatefisttimeframe(firsttime int64, frame *timeframe) {
	logger.Info("updatefisttimeframe", time.Unix(firsttime, 0))
	if frame != nil {
		frame.FirstTime = firsttime
	}
}

func (s *LendingOhlcvService) updatelasttimeframe(lasttime int64, frame *timeframe) {
	if frame != nil {
		frame.LastTime = lasttime
	}

}

func (s *LendingOhlcvService) flatten() []*types.LendingTick {
	var ticks []*types.LendingTick
	for _, tickbytime := range s.lendingTickCache.ticks {
		for _, tick := range tickbytime {
			ticks = append(ticks, tick)
		}
	}
	return ticks
}

func (s *LendingOhlcvService) commitCache() error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	logger.Info("commit ohlcv cache")
	s.truncate()
	ticks := s.flatten()
	lendingtickfile := &lendingtickfile{
		Frame:        s.lendingTickCache.tframes,
		LendingTicks: ticks,
	}
	tickData, err := json.Marshal(lendingtickfile)
	if err != nil {
		return err
	}
	file, err := os.Create(LendingCachePath)
	defer file.Close()
	if err == nil {
		_, err = file.Write(tickData)
		if err != nil {
			return err
		}
	}
	return nil
}

func (s *LendingOhlcvService) loadCache() error {
	file, err := os.Open(LendingCachePath)
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
	var tickf lendingtickfile
	err = json.Unmarshal(bytes, &tickf)
	if err != nil {
		return err
	}
	for _, t := range tickf.LendingTicks {
		s.addTick(t)
	}
	s.lendingTickCache.tframes = tickf.Frame
	return nil
}

func (s *LendingOhlcvService) getTickKey(term uint64, lendingToken common.Address, duration int64, unit string) string {
	return fmt.Sprintf("%d::%s::%s::%s", term, lendingToken.Hex(), strconv.FormatInt(duration, 10), unit)
}

func (s *LendingOhlcvService) parseTickKey(key string) (uint64, common.Address, int64, string, error) {
	keys := strings.Split(key, "::")
	if len(keys) != 4 {
		return 0, common.Address{}, 0, "", errors.New("invalid key")
	}
	term, err := strconv.ParseUint(keys[0], 10, 64)
	if err != nil {
		return 0, common.Address{}, 0, "", errors.New("invalid key")
	}
	lendingToken := common.HexToAddress(keys[1])
	duration, err := strconv.ParseInt(keys[2], 10, 64)
	if err != nil {
		return 0, common.Address{}, 0, "", errors.New("invalid key")
	}
	unit := keys[3]
	return term, lendingToken, duration, unit, nil
}

// updateTick update lastest tick, need to be lock
func (s *LendingOhlcvService) updateTick(key string, trade *types.LendingTrade) error {
	tradeTime := trade.CreatedAt.Unix()
	term, lendingToken, duration, unit, err := s.parseTickKey(key)
	if err != nil {
		return err
	}
	if term == trade.Term && lendingToken.Hex() == trade.LendingToken.Hex() {
		modTime, _ := utils.GetModTime(tradeTime, duration, unit)
		if _, ok := s.lendingTickCache.ticks[key]; !ok {
			s.lendingTickCache.ticks[key] = make(map[int64]*types.LendingTick)
		}
		if tickByTime, ok1 := s.lendingTickCache.ticks[key]; ok1 {
			if last, ok2 := tickByTime[modTime]; ok2 {
				last.Timestamp = modTime
				last.Close = trade.Interest
				if last.High < trade.Interest {
					last.High = trade.Interest
				}
				if last.Low > trade.Interest {
					last.Low = trade.Interest
				}
				last.Volume = big.NewInt(0).Add(last.Volume, trade.Amount)
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
			} else {
				tick := &types.LendingTick{
					LendingID: types.LendingID{
						Term:         trade.Term,
						LendingToken: trade.LendingToken,
					},
					Open:      trade.Interest,
					Close:     trade.Interest,
					High:      trade.Interest,
					Low:       trade.Interest,
					Volume:    trade.Amount,
					Count:     big.NewInt(1),
					Timestamp: modTime,
					Duration:  duration,
					Unit:      unit,
				}
				tickByTime[modTime] = tick
			}
		}
	}

	return nil
}

func (s *LendingOhlcvService) addTick(tick *types.LendingTick) {
	key := s.getTickKey(tick.LendingID.Term, tick.LendingID.LendingToken, tick.Duration, tick.Unit)
	if _, ok := s.lendingTickCache.ticks[key]; ok {

		s.lendingTickCache.ticks[key][tick.Timestamp] = tick
	} else {
		s.lendingTickCache.ticks[key] = make(map[int64]*types.LendingTick)
		s.lendingTickCache.ticks[key][tick.Timestamp] = tick
	}

}

func (s *LendingOhlcvService) filterTick(key string, start, end int64) []*types.LendingTick {
	var res []*types.LendingTick
	if _, ok := s.lendingTickCache.ticks[key]; ok {
		for _, t := range s.lendingTickCache.ticks[key] {
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

// Get24hTick get 24h tick of token
func (s *LendingOhlcvService) Get24hTick(term uint64, lendingToken common.Address) *types.LendingTick {
	return s.get24hTick(term, lendingToken)
}
func (s *LendingOhlcvService) get24hTick(term uint64, lendingToken common.Address) *types.LendingTick {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	var res []*types.LendingTick
	now := time.Now()
	begin := now.AddDate(0, 0, -1).Unix()
	key := s.getTickKey(term, lendingToken, 1, "min")
	res = s.filterTick(key, begin, 0)

	if len(res) >= 1 {
		first := res[0]
		last := res[len(res)-1]
		high := first.High
		low := first.Low
		volume := big.NewInt(0)
		count := big.NewInt(0)
		for _, t := range res {
			if high < t.High {
				high = t.High
			}
			if low > t.Low {
				low = t.Low
			}
			volume = volume.Add(volume, t.Volume)
			count = count.Add(count, t.Count)
		}
		return &types.LendingTick{
			Open:      first.Open,
			Close:     last.Close,
			High:      high,
			Low:       low,
			Count:     count,
			Volume:    volume,
			Timestamp: last.Timestamp,
			LendingID: types.LendingID{
				Term:         term,
				LendingToken: lendingToken,
			},
			Duration: 24,
			Unit:     "hour",
		}
	}
    return &types.LendingTick{
        LendingID: types.LendingID{
            Term:         term,
            LendingToken: lendingToken,
        },
        Duration: 24,
        Unit:     "hour",
        Volume:    big.NewInt(0),
    }
}

// NotifyTrade trigger if trade comming
func (s *LendingOhlcvService) NotifyTrade(trade *types.LendingTrade) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for _, d := range s.getConfig() {
		key := s.getTickKey(trade.Term, trade.LendingToken, d.duration, d.unit)
		s.updateTick(key, trade)
	}
	lastFrame := s.lastTimeFrame()
	s.updatelasttimeframe(trade.CreatedAt.Unix(), lastFrame)
}

// GetOHLCV fetches OHLCV data using
// pairName: can be "" for fetching data for all pairs
// duration: in integer
// unit: sec,min,hour,day,week,month,yr
// timeInterval: 0-2 entries (0 argument: latest data,1st argument: from timestamp, 2nd argument: to timestamp)
func (s *LendingOhlcvService) GetOHLCV(term uint64, lendingToken common.Address, duration int64, unit string, timeInterval ...int64) ([]*types.LendingTick, error) {
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
	cacheKey := s.getTickKey(term, lendingToken, duration, unit)
	ticks := s.filterTick(cacheKey, start.Unix(), end.Unix())
	if ticks == nil {
	}
	return ticks, nil
}

// GetTokenPairData get tick of pair
func (s *LendingOhlcvService) getTokenPairData(term uint64, lendingToken common.Address) *types.LendingTick {
	return s.get24hTick(term, lendingToken)
}

// GetTokenPairData get tick of pair
func (s *LendingOhlcvService) GetTokenPairData(term uint64, lendingToken common.Address) *types.LendingTick {
	p, err := s.lendingPairDao.GetByLendingID(term, lendingToken)
	tick := s.getTokenPairData(term, lendingToken)
	if err == nil {
		tick.LendingID.Name = utils.GetLendingPairName(term, p.LendingTokenSymbol)
	}
	return tick
}

// GetAllTokenPairData get tick of all tokens
func (s *LendingOhlcvService) GetAllTokenPairData() ([]*types.LendingTick, error) {
	pairs, err := s.lendingPairDao.GetAll()
	if err != nil {
		return nil, err
	}
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	pairsData := make([]*types.LendingTick, 0)
	for _, p := range pairs {
		pairData := s.getTokenPairData(p.Term, p.LendingTokenAddress)
		if pairData != nil {
			pairData.LendingID = types.LendingID{
				Term:         p.Term,
				LendingToken: p.LendingTokenAddress,
				Name:         utils.GetLendingPairName(p.Term, p.LendingTokenSymbol),
			}
			pairsData = append(pairsData, pairData)
		}

	}
	return pairsData, nil
}
