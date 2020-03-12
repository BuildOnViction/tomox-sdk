package ws

import (
	"sync"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var lockLendingOderbook = &sync.Mutex{}
var lendingOrderbookSocket *LendingOrderBookSocket

// LendingOrderBookSocket holds the map of subscriptions subscribed to orderbook channels
// corresponding to the key/event they have subscribed to.
type LendingOrderBookSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

// NewLendingOrderBookSocket new lending order book instance
func NewLendingOrderBookSocket() *LendingOrderBookSocket {
	return &LendingOrderBookSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetLendingOrderBookSocket return singleton instance of LendingOrderBookSocket type struct
func GetLendingOrderBookSocket() *LendingOrderBookSocket {
	if lendingOrderbookSocket == nil {
		lendingOrderbookSocket = NewLendingOrderBookSocket()
	}

	return lendingOrderbookSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
// pair := utils.GetPairKey(bt, qt)
func (s *LendingOrderBookSocket) Subscribe(channelID string, c *Client) error {
	if c == nil {
		return errors.New("No connection found")
	}

	if s.subscriptions[channelID] == nil {
		s.subscriptions[channelID] = make(map[*Client]bool)
	}

	s.subscriptions[channelID][c] = true

	if s.subscriptionsList[c] == nil {
		s.subscriptionsList[c] = []string{}
	}

	s.subscriptionsList[c] = append(s.subscriptionsList[c], channelID)

	return nil
}

// UnsubscribeChannelHandler unsubscribes a connection from a certain orderbook channel id
func (s *LendingOrderBookSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

// UnsubscribeHandler unsubscribe lending orderbook handler
func (s *LendingOrderBookSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// UnsubscribeChannel removes a websocket connection from the orderbook channel updates
func (s *LendingOrderBookSocket) UnsubscribeChannel(channelID string, c *Client) {
	lockLendingOderbook.Lock()
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
	lockLendingOderbook.Unlock()
}

// Unsubscribe unsubscribe
func (s *LendingOrderBookSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage streams message to all the subscribtions subscribed to the pair
func (s *LendingOrderBookSocket) BroadcastMessage(channelID string, p interface{}) error {

	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the orderbook channel
func (s *LendingOrderBookSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(LendingOrderBookChannel, msgType, p)
}

// SendInitMessage sends INIT message on orderbook channel on subscription event
func (s *LendingOrderBookSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(LendingOrderBookChannel, types.INIT, data)
}

// SendUpdateMessage sends UPDATE message on orderbook channel as new data is created
func (s *LendingOrderBookSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(LendingOrderBookChannel, types.UPDATE, data)
}

// SendErrorMessage sends error message on orderbook channel
func (s *LendingOrderBookSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(LendingOrderBookChannel, types.ERROR, data)
}
