package types

type MarketData struct {
	PairData        []*PairData        `json:"pairData" bson:"pairData"`
	SmallChartsData map[string][]*Tick `json:"smallChartsData" bson:"smallChartsData"`
}
