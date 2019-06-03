package types

type MarketData struct {
	PairData        []*PairData        `json:"pairData" bson:"pairData"`
	SmallChartsData map[string][]*Tick `json:"smallChartsData" bson:"smallChartsData"`
}

type ChartItem [2]float32

type CoinsIDMarketChart struct {
	Prices       []*ChartItem `json:"prices"`
	MarketCaps   []*ChartItem `json:"market_caps"`
	TotalVolumes []*ChartItem `json:"total_volumes"`
}

type FiatPriceItem struct {
	Symbol       string `json:"symbol" bson:"symbol"`
	Price        string `json:"price" bson:"price"`
	Timestamp    string `json:"timestamp" bson:"timestamp"`
	FiatCurrency string `json:"fiatCurrency" bson:"fiatCurrency"`
}

type FiatPriceItemBSONUpdate struct {
	*FiatPriceItem
}
