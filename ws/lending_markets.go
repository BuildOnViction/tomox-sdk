package ws

import (
	"sync"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var lendingMarketsSocket *LendingMarketsSocket
var lockLendingMarket = &sync.Mutex{}

// LendingMarketsSocket holds the map of subscriptions subscribed to markets channels
// corresponding to the key/event they have subscribed to.
type LendingMarketsSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

// NewLendingMarketsSocket new lending market socket
func NewLendingMarketsSocket() *LendingMarketsSocket {
	return &LendingMarketsSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetLendingMarketSocket return singleton instance of LendingMarketsSocket type struct
func GetLendingMarketSocket() *LendingMarketsSocket {
	if lendingMarketsSocket == nil {
		lendingMarketsSocket = NewLendingMarketsSocket()
	}
	return lendingMarketsSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
func (s *LendingMarketsSocket) Subscribe(channelID string, c *Client) error {
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

// UnsubscribeChannelHandler unsubscribes a connection from a certain markets channel id
func (s *LendingMarketsSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

// UnsubscribeHandler unsubscribes a connection from a certain markets channel id
func (s *LendingMarketsSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// UnsubscribeChannel removes a websocket connection from the markets channel updates
func (s *LendingMarketsSocket) UnsubscribeChannel(channelID string, c *Client) {
	lockLendingMarket.Lock()
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
	lockLendingMarket.Unlock()
}

// Unsubscribe Unsubscribe a connection from a certain markets channel id
func (s *LendingMarketsSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage streams message to all the subscriptions subscribed to the pair
func (s *LendingMarketsSocket) BroadcastMessage(channelID string, p interface{}) error {

	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the markets channel
func (s *LendingMarketsSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(LendingMarketsChannel, msgType, p)
}

// SendInitMessage sends INIT message on markets channel on subscription event
func (s *LendingMarketsSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(LendingMarketsChannel, types.INIT, data)
}

// SendUpdateMessage sends UPDATE message on markets channel as new data is created
func (s *LendingMarketsSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(LendingMarketsChannel, types.UPDATE, data)
}

// SendErrorMessage sends error message on markets channel
func (s *LendingMarketsSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(LendingMarketsChannel, types.ERROR, data)
}
