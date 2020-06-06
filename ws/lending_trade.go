package ws

import (
	"sync"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var lendingTradeSocket *LendingTradeSocket

// LendingTradeSocket holds the map of connections subscribed to pair channels
// corresponding to the key/event they have subscribed to.
type LendingTradeSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
	subsMutex         sync.RWMutex
	subsListMutex     sync.RWMutex
}

// NewLendingTradeSocket init lending socket instance
func NewLendingTradeSocket() *LendingTradeSocket {
	return &LendingTradeSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetLendingTradeSocket get current lending socket
func GetLendingTradeSocket() *LendingTradeSocket {
	if lendingTradeSocket == nil {
		lendingTradeSocket = NewLendingTradeSocket()
	}

	return lendingTradeSocket
}

// Subscribe registers a new websocket connections to the trade channel updates
func (s *LendingTradeSocket) Subscribe(channelID string, c *Client) error {
	s.subsMutex.Lock()
	s.subsListMutex.Lock()
	defer s.subsMutex.Unlock()
	defer s.subsListMutex.Unlock()

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

// UnsubscribeChannelHandler unsubscribes a connection from a certain trade channel id
func (s *LendingTradeSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

// UnsubscribeHandler removes a websocket connection from the trade channel updates
func (s *LendingTradeSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// UnsubscribeChannel removes a websocket connection from the trade channel updates
func (s *LendingTradeSocket) UnsubscribeChannel(channelID string, c *Client) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

// Unsubscribe removes a websocket connection from the trade channel updates
func (s *LendingTradeSocket) Unsubscribe(c *Client) {
	s.subsListMutex.RLock()
	defer s.subsListMutex.RUnlock()
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastMessage broadcasts trade message to all subscribed sockets
func (s *LendingTradeSocket) BroadcastMessage(channelID string, p interface{}) {
	go func() {
		s.subsMutex.RLock()
		defer s.subsMutex.RUnlock()
		for conn, active := range lendingTradeSocket.subscriptions[channelID] {
			if active {
				s.SendUpdateMessage(conn, p)
			}
		}
	}()
}

// SendMessage sends a websocket message on the trade channel
func (s *LendingTradeSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(LendingTradeChannel, msgType, p)
}

// SendInitMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *LendingTradeSocket) SendInitMessage(c *Client, p interface{}) {
	c.SendMessage(LendingTradeChannel, types.INIT, p)
}

// SendUpdateMessage is responsible for sending message on trade ohlcv channel at subscription
func (s *LendingTradeSocket) SendUpdateMessage(c *Client, p interface{}) {
	c.SendMessage(LendingTradeChannel, types.UPDATE, p)
}

// SendErrorMessage sends an error message on the trade channel
func (s *LendingTradeSocket) SendErrorMessage(c *Client, p interface{}) {
	c.SendMessage(LendingTradeChannel, types.ERROR, p)
}
