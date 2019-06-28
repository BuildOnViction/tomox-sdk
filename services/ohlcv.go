package services

import (
	"math"
	"time"

	"github.com/globalsign/mgo/bson"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils"
	"github.com/tomochain/tomox-sdk/ws"
)

type OHLCVService struct {
	tradeDao interfaces.TradeDao
}

func NewOHLCVService(TradeDao interfaces.TradeDao) *OHLCVService {
	return &OHLCVService{TradeDao}
}

// Unsubscribe handles all the unsubscription messages for ticks corresponding to a pair
func (s *OHLCVService) Unsubscribe(conn *ws.Client) {
	ws.GetOHLCVSocket().Unsubscribe(conn)
}

// Unsubscribe handles all the unsubscription messages for ticks corresponding to a pair
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

// GetOHLCV fetches OHLCV data using
// pairName: can be "" for fetching data for all pairs
// duration: in integer
// unit: sec,min,hour,day,week,month,yr
// timeInterval: 0-2 entries (0 argument: latest data,1st argument: from timestamp, 2nd argument: to timestamp)
func (s *OHLCVService) GetOHLCV(pairs []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error) {
	res := make([]*types.Tick, 0)

	currentTimestamp := time.Now().Unix()

	modTime, intervalInSeconds := getModTime(currentTimestamp, duration, unit)

	start := time.Unix(modTime-intervalInSeconds, 0)
	end := time.Unix(currentTimestamp, 0)

	if len(timeInterval) >= 1 {
		end = time.Unix(timeInterval[1], 0)
		start = time.Unix(timeInterval[0], 0)
	}

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

func getModTime(ts, interval int64, unit string) (int64, int64) {
	var modTime, intervalInSeconds int64
	switch unit {
	case "sec":
		intervalInSeconds = interval
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "min":
		intervalInSeconds = interval * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "hour":
		intervalInSeconds = interval * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "day":
		intervalInSeconds = interval * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "week":
		intervalInSeconds = interval * 7 * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "month":
		d := time.Date(time.Now().Year(), time.Now().Month()+1, 1, 0, 0, 0, 0, time.UTC).Day()
		intervalInSeconds = interval * int64(d) * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))

	case "year":
		// Number of days in current year
		d := time.Date(time.Now().Year()+1, 1, 1, 0, 0, 0, 0, time.UTC).Sub(time.Date(time.Now().Year(), 0, 0, 0, 0, 0, 0, time.UTC)).Hours() / 24
		intervalInSeconds = interval * int64(d) * 24 * 60 * 60
		modTime = ts - int64(math.Mod(float64(ts), float64(intervalInSeconds)))
	}

	return modTime, intervalInSeconds
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
		"count":  bson.M{"$sum": one},
		"high":   bson.M{"$max": "$pricepoint"},
		"low":    bson.M{"$min": "$pricepoint"},
		"open":   bson.M{"$first": "$pricepoint"},
		"close":  bson.M{"$last": "$pricepoint"},
		"volume": bson.M{"$sum": bson.M{"$toDecimal": "$amount"}},
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
