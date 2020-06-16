package types

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// UserTrade user trade
type UserTrade struct {
	UserAddress      common.Address `json:"userAddress"`
	Count            *big.Int       `json:"count"`
	Volume           *big.Int       `json:"volume"`
	VolumeByQuote    *big.Int       `json:"volumeByQuote"`
	VolumeAskByQuote *big.Int       `json:"volumeAskByQuote"`
	VolumeBidByQuote *big.Int       `json:"volumeBidByQuote"`

	VolumeAsk      *big.Int       `json:"volumeAsk"`
	VolumeBid      *big.Int       `json:"volumeBid"`
	TimeStamp      int64          `json:"timestamp"`
	RelayerAddress common.Address `json:"relayerAddress"`
	BaseToken      common.Address `json:"baseToken"`
	QuoteToken     common.Address `json:"quoteToken"`
}

// RelayerTrade relayer trade
type RelayerTrade struct {
	RelayerAddress common.Hash `json:"relayerAddress"`
	Count          *big.Int    `json:"count"`
	Volume         *big.Int    `json:"volume"`
}

// UserTradeSpec user trade filter
type UserTradeSpec struct {
}

// UserVolume user volume trade
type UserVolume struct {
	UserAddress common.Address `json:"userAddress"`
	Volume      *big.Int       `json:"volume"`
	Rank        int            `json:"rank"`
}

// TradeVolume trade volume info
type TradeVolume struct {
	Trader      *big.Int `json:"trader"`
	TotalVolume *big.Int `json:"totalVolume"`
}
