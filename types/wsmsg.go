package types

import (
	"encoding/json"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
)

// SubscriptionEvent is an enum signifies whether the incoming message is of type Subscribe or unsubscribe
type SubscriptionEvent string

// Enum members for SubscriptionEvent
const (
	SUBSCRIBE   SubscriptionEvent = "SUBSCRIBE"
	UNSUBSCRIBE SubscriptionEvent = "UNSUBSCRIBE"
	Fetch       SubscriptionEvent = "fetch"
)

const TradeChannel = "trades"
const OrderbookChannel = "order_book"
const OrderChannel = "orders"
const OHLCVChannel = "ohlcv"

type WebSocketMessage struct {
	Channel string         `json:"channel"`
	Event   WebsocketEvent `json:"event"`
}

type WebsocketEvent struct {
	Type    string      `json:"type"`
	Hash    string      `json:"hash,omitempty"`
	Payload interface{} `json:"payload"`
}

type WebSocketSubscription struct {
	Event  SubscriptionEvent `json:"event"`
	Pair   PairSubDoc        `json:"pair"`
	Params `json:"params"`
}

// Params is a sub document used to pass parameters in Subscription messages
type Params struct {
	From     int64  `json:"from"`
	To       int64  `json:"to"`
	Duration int64  `json:"duration"`
	Units    string `json:"units"`
	TickID   string `json:"tickID"`
}

type SignaturePayload struct {
	Order   *Order            `json:"order"`
	Matches []*OrderTradePair `json:"matches"`
}

func NewOrderWebsocketMessage(o *Order) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "NEW_ORDER",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderAddedWebsocketMessage(o *Order, p *Pair, filled int64) *WebSocketMessage {
	o.Process(p)
	o.FilledAmount = big.NewInt(filled)
	o.Status = "OPEN"
	return &WebSocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "ORDER_ADDED",
			Hash:    o.Hash.Hex(),
			Payload: o,
		},
	}
}

func NewOrderCancelWebsocketMessage(oc *OrderCancel) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "CANCEL_ORDER",
			Hash:    oc.Hash.Hex(),
			Payload: oc,
		},
	}
}

func NewRequestSignaturesWebsocketMessage(hash common.Hash, m []*OrderTradePair, o *Order) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "REQUEST_SIGNATURE",
			Hash:    hash.Hex(),
			Payload: SignaturePayload{o, m},
		},
	}
}

func NewSubmitSignatureWebsocketMessage(hash string, m []*OrderTradePair, o *Order) *WebSocketMessage {
	return &WebSocketMessage{
		Channel: "orders",
		Event: WebsocketEvent{
			Type:    "SUBMIT_SIGNATURE",
			Hash:    hash,
			Payload: SignaturePayload{o, m},
		},
	}
}

func (w *WebSocketMessage) Print() {
	b, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		logger.Error(err)
	}

	logger.Info(string(b))
}

func (w *WebsocketEvent) Print() {
	b, err := json.MarshalIndent(w, "", "  ")
	if err != nil {
		logger.Error(err)
	}

	logger.Info(string(b))
}
