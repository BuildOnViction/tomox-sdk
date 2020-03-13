package types

// LendingMarketData lending pair tick data
type LendingMarketData struct {
	PairData []*LendingTick `json:"pairData" bson:"pairData"`
}
