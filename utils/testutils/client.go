package testutils

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/utils"
	"github.com/tomochain/backend-matching-engine/ws"
)

var wg = &sync.WaitGroup{}
var addr = flag.String("addr", "localhost:8080", "http service address")
var logger = utils.TerminalLogger

// Client simulates the client websocket handler that will be used to perform trading.
// requests and responses are respectively the outbound and incoming messages.
// requestLogs and responseLogs are arrays of messages that denote the history of received messages
// wallet is the ethereum account used for orders and trades.
// mutex is used to prevent concurrent writes on the websocket connection
type Client struct {
	// ethereumClient *ethclient.Client
	connection     *ws.Conn
	Requests       chan *types.WebSocketMessage
	Responses      chan *types.WebSocketMessage
	Logs           chan *ClientLogMessage
	Wallet         *types.Wallet
	RequestLogs    []types.WebSocketMessage
	ResponseLogs   []types.WebSocketMessage
	mutex          sync.Mutex
	NonceGenerator *rand.Rand
}

// The client log is mostly used for testing. It optionally takes orders, trade,
// error ids and transaction hashes. All these parameters are optional in order to
// allow the client log message to take in a lot of different types of messages
// An error id of -1 means that there was no error.
type ClientLogMessage struct {
	MessageType string                  `json:"messageType"`
	Orders      []*types.Order          `json:"order"`
	Trades      []*types.Trade          `json:"trade"`
	Matches     []*types.OrderTradePair `json:"matches"`
	Tx          *common.Hash            `json:"tx"`
	ErrorID     int8                    `json:"errorID"`
}

type Server interface {
	ServeHTTP(res http.ResponseWriter, req *http.Request)
}

// NewClient a default client struct connected to the given server
func NewClient(w *types.Wallet, s Server) *Client {
	flag.Parse()
	uri := url.URL{Scheme: "ws", Host: "localhost:8080", Path: "/socket"}

	d := wstest.NewDialer(s)
	c, _, err := d.Dial(uri.String(), nil)
	if err != nil {
		panic(err)
	}

	reqs := make(chan *types.WebSocketMessage)
	resps := make(chan *types.WebSocketMessage)
	logs := make(chan *ClientLogMessage)
	reqLogs := make([]types.WebSocketMessage, 0)
	respLogs := make([]types.WebSocketMessage, 0)

	source := rand.NewSource(time.Now().UnixNano())
	ng := rand.New(source)

	return &Client{
		connection:     ws.NewConnection(c),
		Wallet:         w,
		Requests:       reqs,
		Responses:      resps,
		RequestLogs:    reqLogs,
		ResponseLogs:   respLogs,
		Logs:           logs,
		NonceGenerator: ng,
	}
}

// send is used to prevent concurrent writes on the websocket connection
func (c *Client) send(v interface{}) error {
	c.mutex.Lock()
	defer c.mutex.Unlock()

	return c.connection.WriteJSON(v)
}

// start listening and handling incoming messages
func (c *Client) Start() {
	c.handleMessages()
	c.handleIncomingMessages()
}

// handleMessages waits for incoming messages and routes messages to the
// corresponding handler.
// requests are the messages that are written on the client and destined to
// the server. responses are the message that are
func (c *Client) handleMessages() {
	go func() {
		for {
			select {
			case msg := <-c.Requests:
				c.RequestLogs = append(c.RequestLogs, *msg)
				c.handleOrderChannelMessagesOut(*msg)

			case msg := <-c.Responses:
				c.ResponseLogs = append(c.ResponseLogs, *msg)

				switch msg.Channel {
				case "orders":
					go c.handleOrderChannelMessagesIn(msg.Payload)
				case "order_book":
					go c.handleOrderBookChannelMessages(msg.Payload)
				case "trades":
					go c.handleTradeChannelMessages(msg.Payload)
				case "ohlcv":
					go c.handleOHLCVMessages(msg.Payload)
				}
			}
		}
	}()
}

// handleIncomingMessages reads incomings JSON messages from the websocket connection and
// feeds them into the responses channel
func (c *Client) handleIncomingMessages() {
	message := new(types.WebSocketMessage)
	go func() {
		for {
			err := c.connection.ReadJSON(&message)
			if err != nil {
				if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
					log.Print(err)
				}
				break
			}

			c.Responses <- message
		}
	}()
}

// handleChannelMessagesOut
func (c *Client) handleOrderChannelMessagesOut(m types.WebSocketMessage) {
	logger.Infof("Sending %v", m.Payload.Type)
	err := c.send(m)
	if err != nil {
		log.Printf("Error: Could not send signed orders. Payload: %#v", m.Payload)
		return
	}
}

// handleChannelMessagesIn
func (c *Client) handleOrderChannelMessagesIn(p types.WebsocketEvent) {
	logger.Infof("Receiving: %v", p.Type)
	switch p.Type {
	case "ERROR":
		c.handleError(p)
	case "ORDER_ADDED":
		c.handleOrderAdded(p)
	case "ORDER_CANCELLED":
		c.handleOrderCancelled(p)
	case "ORDER_SUCCESS":
		c.handleOrderSuccess(p)
	case "ORDER_ERROR":
		c.handleOrderError(p)
	case "ORDER_PENDING":
		c.handleOrderPending(p)
	case "REQUEST_SIGNATURE":
		c.handleSignatureRequested(p)
	}
}

func (c *Client) handleOrderBookChannelMessages(p types.WebsocketEvent) {
	switch p.Type {
	case "INIT":
		c.handleOrderBookInit(p)
	case "UPDATE":
		c.handleOrderBookUpdate(p)
	}
}

func (c *Client) handleTradeChannelMessages(p types.WebsocketEvent) {
	switch p.Type {
	case "INIT":
		c.handleTradesInit(p)
	case "UPDATE":
		c.handleTradesUpdate(p)
	}
}

func (c *Client) handleOHLCVMessages(p types.WebsocketEvent) {
	switch p.Type {
	case "INIT":
		c.handleOHLCVInit(p)
	case "UPDATE":
		c.handleOHLCVUpdate(p)
	}
}

// handleError handles incoming error mesasges (does not include tx errors)
func (c *Client) handleError(p types.WebsocketEvent) {
	utils.PrintJSON(p)
}

// handleOrderAdded handles incoming order added messages
func (c *Client) handleOrderAdded(p types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(p.Data)
	if err != nil {
		log.Print(err)
	}

	err = o.UnmarshalJSON(bytes)
	if err != nil {
		log.Print(err)
	}

	l := &ClientLogMessage{
		MessageType: "ORDER_ADDED",
		Orders:      []*types.Order{o},
	}

	c.Logs <- l
}

// handleOrderAdded handles incoming order canceled messages
func (c *Client) handleOrderCancelled(p types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(p.Data)
	if err != nil {
		log.Print(err)
	}

	err = o.UnmarshalJSON(bytes)
	if err != nil {
		log.Print(err)
	}

	l := &ClientLogMessage{
		MessageType: "ORDER_CANCELLED",
		Orders:      []*types.Order{o},
	}

	c.Logs <- l
}

// handleSignatureRequested handles incoming signature requested messages.
// It follows up by signing given data and sending back a SUBMIT_SIGNATURES messages
func (c *Client) handleSignatureRequested(p types.WebsocketEvent) {
	data := &types.SignaturePayload{}
	bytes, err := json.Marshal(p.Data)
	if err != nil {
		logger.Error(err)
	}

	err = json.Unmarshal(bytes, data)
	if err != nil {
		logger.Error(err)
	}

	for _, m := range data.Matches {
		t := m.Trade
		c.SetTradeNonce(t)

		err := c.Wallet.SignTrade(t)
		if err != nil {
			logger.Error(err)
		}
	}

	//sign and return the remaining part of the previous order.
	if data.Order != nil {
		c.SetNonce(data.Order)
		err = c.Wallet.SignOrder(data.Order)
		if err != nil {
			logger.Error(err)
		}
	}

	l := &ClientLogMessage{
		MessageType: "REQUEST_SIGNATURE",
		Orders:      []*types.Order{data.Order},
		Matches:     data.Matches,
	}

	c.Logs <- l
	req := types.NewSubmitSignatureWebsocketMessage(p.Hash, data.Matches, data.Order)
	c.Requests <- req
}

// handleOrderPending handles incoming pending messages (the order has been matched/partially matched
// and the execution tx is currently waiting to be mined)
func (c *Client) handleOrderPending(p types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(p.Data)
	if err != nil {
		log.Print(err)
	}

	err = o.UnmarshalJSON(bytes)
	if err != nil {
		log.Print(err)
	}

	l := &ClientLogMessage{
		MessageType: "ORDER_PENDING",
		Orders:      []*types.Order{o},
	}

	c.Logs <- l
}

// handleOrderSuccess handles incoming tx success messages
func (c *Client) handleOrderSuccess(p types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(p.Data)
	if err != nil {
		log.Print(err)
	}

	err = o.UnmarshalJSON(bytes)
	if err != nil {
		log.Print(err)
	}

	l := &ClientLogMessage{
		MessageType: "ORDER_SUCCESS",
		Orders:      []*types.Order{o},
	}

	c.Logs <- l
}

// handleOrderError handles incoming tx error messages (a tx has been sent but the
// the transaction was reverted)
func (c *Client) handleOrderError(p types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(p.Data)
	if err != nil {
		log.Print(err)
	}

	err = o.UnmarshalJSON(bytes)
	if err != nil {
		log.Print(err)
	}

	l := &ClientLogMessage{
		MessageType: "ORDER_ERROR",
		Orders:      []*types.Order{o},
	}

	c.Logs <- l
}

func (c *Client) handleOrderBookInit(p types.WebsocketEvent) {

}

func (c *Client) handleOrderBookUpdate(p types.WebsocketEvent) {

}

func (c *Client) handleTradesInit(p types.WebsocketEvent) {

}

func (c *Client) handleTradesUpdate(p types.WebsocketEvent) {

}

func (c *Client) handleOHLCVInit(p types.WebsocketEvent) {

}

func (c *Client) handleOHLCVUpdate(p types.WebsocketEvent) {

}

func (c *Client) SetNonce(o *types.Order) {
	o.Nonce = big.NewInt(int64(c.NonceGenerator.Intn(1e8)))
}

func (c *Client) SetTradeNonce(t *types.Trade) {
	t.TradeNonce = big.NewInt(int64(c.NonceGenerator.Intn(1e8)))
}

func (c *ClientLogMessage) Print() {
	b, err := json.MarshalIndent(c, "", "  ")
	if err != nil {
		log.Print(err)
	}

	fmt.Print(string(b))
}
