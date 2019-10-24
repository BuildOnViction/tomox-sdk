package services

import (
	"fmt"
	"math/big"
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

type OHLCVService struct {
	tradeDao     interfaces.TradeDao
	pairDao      interfaces.PairDao
	fiatPriceDao interfaces.FiatPriceDao
	ticks        map[string]map[int64]*types.Tick
	mutex        sync.Mutex
}
type Duration struct {
	duration int64
	unit     string
	interval int64
}

func NewOHLCVService(TradeDao interfaces.TradeDao, pairDao interfaces.PairDao, fiatDao interfaces.FiatPriceDao) *OHLCVService {
	return &OHLCVService{
		tradeDao:     TradeDao,
		pairDao:      pairDao,
		fiatPriceDao: fiatDao,
		ticks:        make(map[string]map[int64]*types.Tick),
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

func (s *OHLCVService) getConfig() []Duration {
	return []Duration{
		Duration{
			duration: 1,
			unit:     "min",
			interval: 1 * 60 * 60,
		},
		Duration{
			duration: 5,
			unit:     "min",
			interval: 1 * 60 * 60,
		},
		Duration{
			duration: 15,
			unit:     "min",
			interval: 1 * 60 * 60,
		},
		Duration{
			duration: 30,
			unit:     "min",
			interval: 1 * 60 * 60,
		},
		Duration{
			duration: 1,
			unit:     "hour",
			interval: 7 * 1 * 60 * 60,
		},
		Duration{
			duration: 2,
			unit:     "hour",
			interval: 7 * 1 * 60 * 60,
		},
		Duration{
			duration: 4,
			unit:     "hour",
			interval: 7 * 1 * 60 * 60,
		},
		Duration{
			duration: 12,
			unit:     "hour",
			interval: 7 * 1 * 60 * 60,
		},
		Duration{
			duration: 1,
			unit:     "day",
			interval: 30 * 1 * 60 * 60,
		},
		Duration{
			duration: 1,
			unit:     "week",
			interval: 30 * 1 * 60 * 60,
		},
		Duration{
			duration: 1,
			unit:     "month",
			interval: 30 * 1 * 60 * 60,
		},
		Duration{
			duration: 3,
			unit:     "month",
			interval: 30 * 1 * 60 * 60,
		},
		Duration{
			duration: 6,
			unit:     "month",
			interval: 30 * 1 * 60 * 60,
		},
		Duration{
			duration: 9,
			unit:     "month",
			interval: 30 * 1 * 60 * 60,
		},
		Duration{
			duration: 1,
			unit:     "year",
			interval: 30 * 1 * 60 * 60,
		},
	}
}

// Init init data cache
func (s *OHLCVService) Init() {
	logger.Info("OHLCV init...")
	durations := s.getConfig()
	pairs, err := s.pairDao.GetAll()
	now := time.Now().Unix()
	if err == nil {
		for _, pair := range pairs {
			p := types.PairAddresses{
				Name:       pair.Name(),
				BaseToken:  pair.BaseTokenAddress,
				QuoteToken: pair.QuoteTokenAddress,
			}
			ps := []types.PairAddresses{p}
			for _, d := range durations {
				yes := now - d.interval
				ticks, err := s.getOHLCV(ps, d.duration, d.unit, time.Unix(yes, 0), time.Unix(now, 0))
				if err == nil {
					key := s.getTickKey(p.BaseToken, p.QuoteToken, d.duration, d.unit)
					for _, t := range ticks {
						t.Timestamp = t.Timestamp / 1000
						s.addTick(key, t)
					}

				}
			}

		}
	}
	logger.Info("OHLCV finish")
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

func (s *OHLCVService) commitCache(tick *types.Tick) {

}

// updateTick update lastest tick
func (s *OHLCVService) updateTick(key string, trade *types.Trade) error {
	tradeTime := trade.CreatedAt.Unix()
	baseToken, quoteToken, duration, unit, err := s.parseTickKey(key)
	if err != nil {
		return err
	}
	if baseToken.Hex() == trade.BaseToken.Hex() && quoteToken.Hex() == trade.QuoteToken.Hex() {
		modTime, _ := utils.GetModTime(tradeTime, duration, unit)
		if _, ok := s.ticks[key]; !ok {
			s.ticks[key] = make(map[int64]*types.Tick)
		}
		if tickByTime, ok := s.ticks[key]; ok {
			if last, ok := tickByTime[modTime]; ok {
				logger.Info("updateTick", key, modTime)
				last.Timestamp = modTime
				last.Close = trade.PricePoint
				if last.High.Cmp(trade.PricePoint) < 0 {
					last.High = trade.PricePoint
				}
				if last.Low.Cmp(trade.PricePoint) > 0 {
					last.Low = trade.PricePoint
				}
				last.Volume = last.Volume.Add(last.Volume, trade.Amount)
				last.Count = last.Count.Add(last.Count, big.NewInt(1))
				last.CloseTime = trade.CreatedAt
			} else {
				logger.Info("updateTick not found", key, modTime)
				tick := &types.Tick{
					Pair: types.PairID{
						PairName:   trade.PairName,
						BaseToken:  trade.BaseToken,
						QuoteToken: trade.QuoteToken,
					},
					OpenTime:  trade.CreatedAt,
					Open:      trade.PricePoint,
					Close:     trade.PricePoint,
					High:      trade.PricePoint,
					Low:       trade.PricePoint,
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

func (s *OHLCVService) addTick(key string, tick *types.Tick) {
	logger.Info("addTick", key, tick.Timestamp)
	if _, ok := s.ticks[key]; ok {

		s.ticks[key][tick.Timestamp] = tick
	} else {
		s.ticks[key] = make(map[int64]*types.Tick)
		s.ticks[key][tick.Timestamp] = tick
	}

}

func (s *OHLCVService) filterTick(key string, start, end int64) []*types.Tick {
	var res []*types.Tick
	if _, ok := s.ticks[key]; ok {
		for _, t := range s.ticks[key] {
			if t.Timestamp >= start && (t.Timestamp <= end || end == 0) {
				c := *t
				c.Timestamp = t.Timestamp * 1000
				res = append(res, &c)
			}
		}
	} else {
		logger.Info("keynull", key)
		return nil
	}
	sort.Slice(res, func(i, j int) bool {
		return res[i].Timestamp < res[j].Timestamp
	})
	return res
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
		count := big.NewInt(0)
		for _, t := range res {
			if high.Cmp(t.High) < 0 {
				high = t.High
			}
			if low.Cmp(t.Low) > 0 {
				low = t.Low
			}
			volume = volume.Add(volume, t.Volume)
			count = count.Add(count, t.Count)
		}
		return &types.Tick{
			Open:      first.Open,
			Close:     last.Close,
			High:      high,
			Low:       low,
			CloseTime: last.CloseTime,
			Count:     count,
			Volume:    volume,
		}
	}
	return nil
}

// NotifyTrade trigger if trade comming
func (s *OHLCVService) NotifyTrade(trade *types.Trade) {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	for key, _ := range s.ticks {
		s.updateTick(key, trade)
	}
}
func (s *OHLCVService) getOHLCV(pairs []types.PairAddresses, duration int64, unit string, start, end time.Time) ([]*types.Tick, error) {
	logger.Info("Get OHLCV", pairs[0].Name, duration, unit, start, end)
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
	s.mutex.Lock()
	defer s.mutex.Unlock()
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
		logger.Info("getOHLCV", start, end, "key", cacheKey)
		ts, err := s.getOHLCV(pairs, duration, unit, start, end)
		if err == nil {
			for _, t := range ts {
				logger.Info("getOHLCV result", "key", cacheKey, time.Unix(t.Timestamp/1000, 0))
			}
		}
	}
	for _, t := range ticks {
		logger.Info("result", cacheKey, t.Timestamp, t.Open, t.Close, t.Low, t.High, t.Pair.PairName)
	}
	return ticks, nil
}

func getMatchQuery(start, end time.Time, pairs ...types.PairAddresses) bson.M {
	match := bson.M{
		"createdAt": bson.M{
			"$gte": start,

			"$lt": end,
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
func (s *OHLCVService) GetTokenPairData(pairName string, baseTokenSymbol string, baseToken common.Address, quoteToken common.Address) *types.PairData {
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
		pairData.Volume = tick.Volume
		pairData.Close = tick.Close
		pairData.Count = tick.Count
		fiatItem, err := s.fiatPriceDao.GetLastPriceCurrentByTime(baseTokenSymbol, tick.CloseTime)
		if err == nil {
			pairData.CloseBaseUsd, _ = pairData.CloseBaseUsd.SetString(fiatItem.Price)
		}
		return pairData
	}
	return nil
}

// GetAllTokenPairData get tick of all tokens
func (s *OHLCVService) GetAllTokenPairData() ([]*types.PairData, error) {
	pairs, err := s.pairDao.GetActivePairs()
	if err != nil {
		return nil, err
	}
	pairsData := make([]*types.PairData, 0)
	for _, p := range pairs {
		pairData := s.GetTokenPairData(p.Name(), p.BaseTokenSymbol, p.BaseTokenAddress, p.QuoteTokenAddress)
		if pairData != nil {
			pairsData = append(pairsData, pairData)
		}

	}

	return pairsData, nil
}
