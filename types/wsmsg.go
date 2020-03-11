package types

import (
	//"encoding/json"
	"fmt"
	"math/big"

	//"strconv"

	"github.com/ethereum/go-ethereum/common"
)

// SubscriptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubscriptionEvent string

// Enum members for SubscriptionEvent
const (
	SUBSCRIBE   SubscriptionEvent = "SUBSCRIBE"
	UNSUBSCRIBE SubscriptionEvent = "UNSUBSCRIBE"
	Fetch       SubscriptionEvent = "fetch"

	UPDATE        SubscriptionEvent = "UPDATE"
	ERROR         SubscriptionEvent = "ERROR"
	SUCCESS_EVENT SubscriptionEvent = "SUCCESS"
	INIT          SubscriptionEvent = "INIT"
	CANCEL        SubscriptionEvent = "CANCEL"

	// status

	ORDER_ADDED            = "ORDER_ADDED"
	ORDER_FILLED           = "ORDER_FILLED"
	ORDER_PARTIALLY_FILLED = "ORDER_PARTIALLY_FILLED"
	ORDER_CANCELLED        = "ORDER_CANCELLED"
	ORDER_REJECTED         = "ORDER_REJECTED"
	ERROR_STATUS           = "ERROR"

	TradeAdded   = "TRADE_ADDED"
	TradeUpdated = "TRADE_UPDATED"
	// channel
	TradeChannel     = "trades"
	OrderbookChannel = "orderbook"
	OrderChannel     = "orders"
	OHLCVChannel     = "ohlcv"

	// Lending
	LENDING_ORDER_ADDED            = "LENDING_ORDER_ADDED"
	LENDING_ORDER_FILLED           = "LENDING_ORDER_FILLED"
	LENDING_ORDER_PARTIALLY_FILLED = "LENDING_ORDER_PARTIALLY_FILLED"
	LENDING_ORDER_CANCELLED        = "LENDING_ORDER_CANCELLED"
	LENDING_ORDER_REJECTED         = "LENDING_ORDER_REJECTED"
)

type WebsocketMessage struct {
	Channel string         `json:"channel"`
	Event   WebsocketEvent `json:"event"`
}

func (ev *WebsocketMessage) String() string {
	return fmt.Sprintf("%v/%v", ev.Channel, ev.Event.String())
}

type WebsocketEvent struct {
	Type    SubscriptionEvent `json:"type"`
	Hash    string            `json:"hash,omitempty"`
	Payload interface{}       `json:"payload"`
}

func (ev *WebsocketEvent) String() string {
	return fmt.Sprintf("%v", ev.Type)
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	PairID   string `json:"pair"`
}

type OrderPendingPayload struct {
	Matches *Matches `json:"matches"`
}

type OrderSuccessPayload struct {
	Matches *Matches `json:"matches"`
}

type OrderMatchedPayload struct {
	Matches *Matches `json:"matches"`
}

type SubscriptionPayload struct {
	PairName     string         `json:"pairName,omitempty"`
	QuoteToken   common.Address `json:"quoteToken,omitempty"`
	BaseToken    common.Address `json:"baseToken,omitempty"`
	From         int64          `json"from"`
	To           int64          `json:"to"`
	Duration     int64          `json:"duration"`
	Units        string         `json:"units"`
	Term         uint64         `json:"term"`
	LendingToken common.Address `json:"lendingToken,omitempty"`
}

/*
func (s *SubscriptionPayload) UnmarshalJSON(b []byte) error {
	payload := map[string]interface{}{}
	err := json.Unmarshal(b, &payload)
	if err != nil {
		return err
	}
	if payload["pairName"] != nil {
		s.PairName = payload["pairName"].(string)
	}

	if payload["quoteToken"] != nil {
		s.QuoteToken = common.HexToAddress(payload["quoteToken"].(string))
	}
	if payload["baseToken"] != nil {
		s.BaseToken = common.HexToAddress(payload["baseToken"].(string))
	}
	if payload["term"] != nil {
		s.Term, _ = strconv.ParseUint(payload["term"].(string), 10, 64)
	}
	if payload["lendingToken"] != nil {
		s.LendingToken = common.HexToAddress(payload["lendingToken"].(string))
	}
	if payload["units"] != nil {
		s.Units = payload["units"].(string)
	}
	if payload["duration"] != nil {
		s.Duration, _ = strconv.ParseInt(payload["duration"].(string), 10, 64)
	}
	return nil

}
*/

func NewOrderWebsocketMessage(o *Order) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "NEW_ORDER",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderAddedWebsocketMessage(o *Order, p *Pair, filled int64) *WebsocketMessage {
	o.Process(p)
	o.FilledAmount = big.NewInt(filled)
	o.Status = "OPEN"
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "ORDER_ADDED",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderCancelWebsocketMessage(oc *OrderCancel) *WebsocketMessage {
	return &WebsocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "CANCEL_ORDER",
			Hash:    oc.Hash.Hex(),
			Payload: oc,
		},
	}
}
