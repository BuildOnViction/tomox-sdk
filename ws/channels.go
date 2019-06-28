package ws

import (
	"fmt"

	"github.com/tomochain/tomox-sdk/errors"
)

const (
	TradeChannel        = "trades"
	RawOrderBookChannel = "raw_orderbook"
	OrderChannel        = "orders"
	OrderBookChannel    = "orderbook"
	TokenChannel        = "tokens"
	OHLCVChannel        = "ohlcv"
	PriceBoardChannel   = "price_board"
	DepositChannel      = "deposit"
	MarketsChannel      = "markets"
	NotificationChannel = "notification"
)

var socketChannels map[string]func(interface{}, *Client)

func RegisterChannel(channel string, fn func(interface{}, *Client)) error {
	if channel == "" {
		return errors.New("Channel can not be an empty string")
	}

	if fn == nil {
		return errors.New("Handler should not be nil")
	}

	ch := getChannels()
	if ch[channel] != nil {
		return fmt.Errorf("Channel already registered")
	}

	ch[channel] = fn
	return nil
}

func getChannels() map[string]func(interface{}, *Client) {
	if socketChannels == nil {
		socketChannels = make(map[string]func(interface{}, *Client))
	}

	return socketChannels
}
