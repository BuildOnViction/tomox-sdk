package services

import (
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/dex-server/interfaces"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
	"github.com/tomochain/dex-server/ws"
	"gopkg.in/mgo.v2/bson"
)

// TradeService struct with daos required, responsible for communicating with daos.
// TradeService functions are responsible for interacting with daos and implements business logics.
type PriceBoardService struct {
	TradeDao interfaces.TradeDao
}

// NewTradeService returns a new instance of TradeService
func NewPriceBoardService(tradeDao interfaces.TradeDao) *PriceBoardService {
	return &PriceBoardService{TradeDao: tradeDao}
}

// Subscribe
func (s *PriceBoardService) Subscribe(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetPriceBoardSocket()

	// Fix the value at 1 day because we only care about 24h change
	duration := int64(1)
	unit := "day"

	data, err := s.GetPriceBoardData(
		[]types.PairAddresses{{BaseToken: bt, QuoteToken: qt}},
		duration,
		unit,
	)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	id := utils.GetPriceBoardChannelID(bt, qt)
	err = socket.Subscribe(id, c)
	if err != nil {
		logger.Error(err)
		socket.SendErrorMessage(c, err.Error())
		return
	}

	ws.RegisterConnectionUnsubscribeHandler(c, socket.UnsubscribeChannelHandler(id))
	socket.SendInitMessage(c, data)
}

// Unsubscribe
func (s *PriceBoardService) UnsubscribeChannel(c *ws.Client, bt, qt common.Address) {
	socket := ws.GetPriceBoardSocket()

	id := utils.GetPriceBoardChannelID(bt, qt)
	socket.UnsubscribeChannel(id, c)
}

// Unsubscribe
func (s *PriceBoardService) Unsubscribe(c *ws.Client) {
	socket := ws.GetPriceBoardSocket()
	socket.Unsubscribe(c)
}

func (s *PriceBoardService) GetPriceBoardData(pairs []types.PairAddresses, duration int64, unit string, timeInterval ...int64) ([]*types.Tick, error) {
	res := make([]*types.Tick, 0)

	currentTimestamp := time.Now().Unix()

	_, intervalInSeconds := getModTime(currentTimestamp, duration, unit)

	start := time.Unix(currentTimestamp-intervalInSeconds, 0)
	end := time.Unix(currentTimestamp, 0)

	if len(timeInterval) >= 1 {
		end = time.Unix(timeInterval[1], 0)
		start = time.Unix(timeInterval[0], 0)
	}

	match := make(bson.M)
	match = getMatchQuery(start, end, pairs...)
	match = bson.M{"$match": match}

	group := getGroupBson()
	group = bson.M{"$group": group}

	sort := bson.M{"$sort": bson.M{"timestamp": 1}}

	query := []bson.M{match, group, sort}

	utils.PrintJSON(query)
	res, err := s.TradeDao.Aggregate(query)
	if err != nil {
		return nil, err
	}

	if res == nil {
		return []*types.Tick{}, nil
	}

	return res, nil
}

// query for grouping of the documents into one
func getGroupBson() bson.M {
	var group bson.M

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
	groupID["pairName"] = "$pairName"
	groupID["baseToken"] = "$baseToken"
	groupID["quoteToken"] = "$quoteToken"
	group["_id"] = groupID

	return group
}

func (s *PriceBoardService) SyncFiatPrice() {
	//client := &http.Client{}
	//req, err := http.NewRequest("GET", "https://pro-api.coinmarketcap.com/v1/cryptocurrency/quotes/latest?symbol=ETH,TOMO&convert=USD", nil)
	//req.Header.Add("X-CMC_PRO_API_KEY", `a928d4ca-37cc-41b8-a9ea-5a65ae025aa5`)
	//resp, err := client.Do(req)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//body, err := ioutil.ReadAll(resp.Body)
	//if err != nil {
	//	log.Fatalln(err)
	//}
	//
	//log.Println(string(body))
}
