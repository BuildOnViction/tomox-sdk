package ws

import (
	"sync"

	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/types"
)

var lendingOhlcvSocket *LendingOhlcvSocket

var lockLendingOhlcv = &sync.Mutex{}

// LendingOhlcvSocket holds the map of subscribtions subscribed to OHLCV channels
// corresponding to the key/event they have subscribed to.
type LendingOhlcvSocket struct {
	subscriptions     map[string]map[*Client]bool
	subscriptionsList map[*Client][]string
}

// NewLendingOhlcvSocket create new instance
func NewLendingOhlcvSocket() *LendingOhlcvSocket {
	return &LendingOhlcvSocket{
		subscriptions:     make(map[string]map[*Client]bool),
		subscriptionsList: make(map[*Client][]string),
	}
}

// GetLendingOhlcvSocket return singleton instance of LendingOhlcvSocket type struct
func GetLendingOhlcvSocket() *LendingOhlcvSocket {
	if lendingOhlcvSocket == nil {
		lendingOhlcvSocket = NewLendingOhlcvSocket()
	}

	return lendingOhlcvSocket
}

// Subscribe handles the registration of connection to get
// streaming data over the socket for any pair.
func (s *LendingOhlcvSocket) Subscribe(channelID string, c *Client) error {
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
	lockLendingOhlcv.Lock()
	s.subscriptionsList[c] = append(s.subscriptionsList[c], channelID)
	lockLendingOhlcv.Unlock()
	return nil
}

// UnsubscribeChannelHandler returns function of type unsubscribe handler,
// it handles the unsubscription of pair in case of connection closing.
func (s *LendingOhlcvSocket) UnsubscribeChannelHandler(channelID string) func(c *Client) {
	return func(c *Client) {
		s.UnsubscribeChannel(channelID, c)
	}
}

//UnsubscribeHandler returns function of type unsubscribe handler
func (s *LendingOhlcvSocket) UnsubscribeHandler() func(c *Client) {
	return func(c *Client) {
		s.Unsubscribe(c)
	}
}

// UnsubscribeChannel is used to unsubscribe the connection from listening to the key
// subscribed to. It can be called on unsubscription message from user or due to some other reason by
// system
func (s *LendingOhlcvSocket) UnsubscribeChannel(channelID string, c *Client) {
	lockLendingOhlcv.Lock()
	if s.subscriptions[channelID][c] {
		s.subscriptions[channelID][c] = false
		delete(s.subscriptions[channelID], c)
	}
	lockLendingOhlcv.Unlock()
}

// Unsubscribe  returns function of type unsubscribe handler
func (s *LendingOhlcvSocket) Unsubscribe(c *Client) {
	channelIDs := s.subscriptionsList[c]
	if channelIDs == nil {
		return
	}

	for _, id := range s.subscriptionsList[c] {
		s.UnsubscribeChannel(id, c)
	}
}

// BroadcastLendingOhlcv Message streams message to all the subscriptions subscribed to the pair
func (s *LendingOhlcvSocket) BroadcastLendingOhlcv(channelID string, p interface{}) error {
	for c, status := range s.subscriptions[channelID] {
		if status {
			s.SendUpdateMessage(c, p)
		}
	}

	return nil
}

// SendMessage sends a websocket message on the trade channel
func (s *LendingOhlcvSocket) SendMessage(c *Client, msgType types.SubscriptionEvent, p interface{}) {
	c.SendMessage(LendingOhlcvChannel, msgType, p)
}

// SendInitMessage is responsible for sending message on trade channel at subscription
func (s *LendingOhlcvSocket) SendInitMessage(c *Client, p interface{}) {
	c.SendMessage(LendingOhlcvChannel, types.INIT, p)
}

// SendUpdateMessage is responsible for sending message on trade channel at subscription
func (s *LendingOhlcvSocket) SendUpdateMessage(c *Client, p interface{}) {
	c.SendMessage(LendingOhlcvChannel, types.UPDATE, p)
}

// SendErrorMessage sends an error message on the trade channel
func (s *LendingOhlcvSocket) SendErrorMessage(c *Client, p interface{}) {
	c.SendMessage(LendingOhlcvChannel, types.ERROR, p)
}
