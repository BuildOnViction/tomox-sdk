package ws

import (
	"sync"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var lendingPriceBoardSocket *LendingPriceBoardSocket

// LendingPriceBoardSocket holds the map of subscriptions subscribed to price board channels
// corresponding to the key/event they have subscribed to.
type LendingPriceBoardSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
	subsMutex         sync.RWMutex
	subsListMutex     sync.RWMutex
}

func NewLendingPriceBoardSocket() *LendingPriceBoardSocket {
	return &LendingPriceBoardSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetLendingPriceBoardSocket return singleton instance of LendingPriceBoardSocket type struct
func GetLendingPriceBoardSocket() *LendingPriceBoardSocket {
	if lendingPriceBoardSocket == nil {
		lendingPriceBoardSocket = NewLendingPriceBoardSocket()
	}

	return lendingPriceBoardSocket
}

// Subscribe handles the subscription of connection to get
// streaming data over the socker for any pair.
func (s *LendingPriceBoardSocket) Subscribe(channelID string, c *Client) error {
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

// UnsubscribeChannelHandler unsubscribes a connection from a certain lending price board channel id
func (s *LendingPriceBoardSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

func (s *LendingPriceBoardSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// UnsubscribeChannel removes a websocket connection from the price board channel updates
func (s *LendingPriceBoardSocket) UnsubscribeChannel(channelID string, c *Client) {
	s.subsMutex.Lock()
	defer s.subsMutex.Unlock()
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
}

func (s *LendingPriceBoardSocket) Unsubscribe(c *Client) {
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

func (s *LendingPriceBoardSocket) getSubscriptions() map[string]map[*Client]bool {
	s.subsMutex.RLock()
	defer s.subsMutex.RUnlock()
	return s.subscriptions
}

// BroadcastMessage streams message to all the subscriptions subscribed to the pair
func (s *LendingPriceBoardSocket) BroadcastMessage(channelID string, p interface{}) error {
	subs := s.getSubscriptions()
	for c, status := range subs[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the price board channel
func (s *LendingPriceBoardSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(LendingPriceBoardChannel, msgType, p)
}

// SendInitMessage sends INIT message on price board channel on subscription event
func (s *LendingPriceBoardSocket) SendInitMessage(c *Client, data interface{}) {
	c.SendMessage(LendingPriceBoardChannel, types.INIT, data)
}

// SendUpdateMessage sends UPDATE message on price board channel as new data is created
func (s *LendingPriceBoardSocket) SendUpdateMessage(c *Client, data interface{}) {
	c.SendMessage(LendingPriceBoardChannel, types.UPDATE, data)
}

// SendErrorMessage sends error message on price board channel
func (s *LendingPriceBoardSocket) SendErrorMessage(c *Client, data interface{}) {
	c.SendMessage(LendingPriceBoardChannel, types.ERROR, data)
}
