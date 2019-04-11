package types

type MarketData struct {
	PairData        []*PairData        `json:"pairData,omitempty" bson:"pairData"`
	SmallChartsData map[string][]*Tick `json:"smallChartsData,omitempty" bson:"smallChartsData"`
}
