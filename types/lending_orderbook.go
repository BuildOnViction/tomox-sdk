package types

// LendingOrderBook for lending orderbook
type LendingOrderBook struct {
	Name   string              `json:"name"`
	Borrow []map[string]string `json:"borrow"`
	Lend   []map[string]string `json:"lend"`
}

// RawLendingOrderBook for lending orderbook
type RawLendingOrderBook struct {
	PairName string   `json:"pairName"`
	Orders   []*Order `json:"orders"`
}
