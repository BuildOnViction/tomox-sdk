package ws

import (
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/tomochain/tomox-sdk/types"
)

type Client struct {
	*websocket.Conn
	mu   sync.Mutex
	send chan types.WebsocketMessage
}

var unsubscribeHandlers map[*Client][]func(*Client)

func NewClient(c *websocket.Conn) *Client {
	conn := &Client{Conn: c, mu: sync.Mutex{}, send: make(chan types.WebsocketMessage)}

	if unsubscribeHandlers == nil {
		unsubscribeHandlers = make(map[*Client][]func(*Client))
	}

	if unsubscribeHandlers[conn] == nil {
		unsubscribeHandlers[conn] = make([]func(*Client), 0)
	}

	return conn
}

func (c *Client) writeMessage(m types.WebsocketMessage) {
	c.SetWriteDeadline(time.Now().Add(writeWait))
	err := c.WriteJSON(m)
	if err != nil {
		logger.Info("writeMessage closing connection:", err)
		c.closeConnection()
		return
	}

}

// SendMessage constructs the message with proper structure to be sent over websocket
func (c *Client) SendMessage(channel string, msgType types.SubscriptionEvent, payload interface{}, h ...common.Hash) {
	e := types.WebsocketEvent{
		Type:    msgType,
		Payload: payload,
	}

	if len(h) > 0 {
		e.Hash = h[0].Hex()
	}

	logger.Debug("SendMessage", channel, msgType)

	m := types.WebsocketMessage{
		Channel: channel,
		Event:   e,
	}

	c.mu.Lock()
	defer c.mu.Unlock()

	c.writeMessage(m)
}

// SendPingMessage check conntection
func (c *Client) SendPingMessage() error {
	c.mu.Lock()
	defer c.mu.Unlock()
	c.SetWriteDeadline(time.Now().Add(writeWait))
	return c.WriteMessage(websocket.PingMessage, nil)
}

func (c *Client) closeConnection() {
	for _, unsub := range unsubscribeHandlers[c] {
		unsub(c)
	}

	c.Close()
}

func (c *Client) SendOrderErrorMessage(err error, h common.Hash) {
	p := map[string]interface{}{
		"message": err.Error(),
		"hash":    h.Hex(),
	}

	e := types.WebsocketEvent{
		Type:    "ERROR",
		Payload: p,
	}

	m := types.WebsocketMessage{
		Channel: OrderChannel,
		Event:   e,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeMessage(m)
}

// SendLendingOrderErrorMessage send error lending transaction
func (c *Client) SendLendingOrderErrorMessage(err error, h common.Hash) {
	p := map[string]interface{}{
		"message": err.Error(),
		"hash":    h.Hex(),
	}

	e := types.WebsocketEvent{
		Type:    "ERROR",
		Payload: p,
	}

	m := types.WebsocketMessage{
		Channel: LendingOrderChannel,
		Event:   e,
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	c.writeMessage(m)
}
