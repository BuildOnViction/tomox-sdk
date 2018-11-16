package testutils

import (
	"encoding/json"
	"flag"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"net/url"
	"sync"
	"time"

	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/utils"
	"github.com/tomochain/backend-matching-engine/ws"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/websocket"
	"github.com/posener/wstest"
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
	connection     *ws.Client
	Requests       chan *types.WebsocketMessage
	Responses      chan *types.WebsocketMessage
	Logs           chan *ClientLogMessage
	Wallet         *types.Wallet
	RequestLogs    []types.WebsocketMessage
	ResponseLogs   []types.WebsocketMessage
	mutex          sync.Mutex
	NonceGenerator *rand.Rand
}

// The client log is mostly used for testing. It optionally takes orders, trade,
// error ids and transaction hashes. All these parameters are optional in order to
// allow the client log message to take in a lot of different types of messages
// An error id of -1 means that there was no error.
type ClientLogMessage struct {
	MessageType string         `json:"messageType"`
	Orders      []*types.Order `json:"order"`
	Trades      []*types.Trade `json:"trade"`
	Matches     *types.Matches `json:"matches"`
	Tx          *common.Hash   `json:"tx"`
	ErrorID     int8           `json:"errorID"`
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

	reqs := make(chan *types.WebsocketMessage)
	resps := make(chan *types.WebsocketMessage)
	logs := make(chan *ClientLogMessage)
	reqLogs := make([]types.WebsocketMessage, 0)
	respLogs := make([]types.WebsocketMessage, 0)

	source := rand.NewSource(time.Now().UnixNano())
	ng := rand.New(source)

	return &Client{
		connection:     ws.NewClient(c),
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
					go c.handleOrderChannelMessagesIn(msg.Event)
				case "order_book":
					go c.handleOrderBookChannelMessages(msg.Event)
				case "trades":
					go c.handleTradeChannelMessages(msg.Event)
				case "ohlcv":
					go c.handleOHLCVMessages(msg.Event)
				}
			}
		}
	}()
}

// handleIncomingMessages reads incomings JSON messages from the websocket connection and
// feeds them into the responses channel
func (c *Client) handleIncomingMessages() {
	message := new(types.WebsocketMessage)
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
func (c *Client) handleOrderChannelMessagesOut(m types.WebsocketMessage) {
	logger.Infof("Sending %v", m.Event.Type)
	err := c.send(m)
	if err != nil {
		log.Printf("Error: Could not send signed orders. Payload: %#v", m.Event)
		return
	}
}

// handleChannelMessagesIn
func (c *Client) handleOrderChannelMessagesIn(e types.WebsocketEvent) {
	logger.Infof("Receiving: %v", e.Type)

	switch e.Type {
	case "ERROR":
		c.handleError(e)
	case "ORDER_ADDED":
		c.handleOrderAdded(e)
	case "ORDER_CANCELLED":
		c.handleOrderCancelled(e)
	case "ORDER_SUCCESS":
		c.handleOrderSuccess(e)
	case "ORDER_ERROR":
		c.handleOrderError(e)
	case "ORDER_PENDING":
		c.handleOrderPending(e)
	}
}

func (c *Client) handleOrderBookChannelMessages(e types.WebsocketEvent) {
	switch e.Type {
	case "INIT":
		c.handleOrderBookInit(e)
	case "UPDATE":
		c.handleOrderBookUpdate(e)
	}
}

func (c *Client) handleTradeChannelMessages(e types.WebsocketEvent) {
	switch e.Type {
	case "INIT":
		c.handleTradesInit(e)
	case "UPDATE":
		c.handleTradesUpdate(e)
	}
}

func (c *Client) handleOHLCVMessages(e types.WebsocketEvent) {
	switch e.Type {
	case "INIT":
		c.handleOHLCVInit(e)
	case "UPDATE":
		c.handleOHLCVUpdate(e)
	}
}

// handleError handles incoming error mesasges (does not include tx errors)
func (c *Client) handleError(e types.WebsocketEvent) {
	utils.PrintJSON(e)
}

// handleOrderAdded handles incoming order added messages
func (c *Client) handleOrderAdded(e types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(e.Payload)
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

// handleOrderAdded handles incoming order cancelled messages
func (c *Client) handleOrderCancelled(e types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(e.Payload)
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

// handleOrderPending handles incoming pending messages (the order has been matched/partially matched
// and the execution tx is currently waiting to be mined)
func (c *Client) handleOrderPending(e types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(e.Payload)
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
func (c *Client) handleOrderSuccess(e types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(e.Payload)
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
func (c *Client) handleOrderError(e types.WebsocketEvent) {
	o := &types.Order{}

	bytes, err := json.Marshal(e.Payload)
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

func (c *Client) handleOrderBookInit(e types.WebsocketEvent) {

}

func (c *Client) handleOrderBookUpdate(e types.WebsocketEvent) {

}

func (c *Client) handleTradesInit(e types.WebsocketEvent) {

}

func (c *Client) handleTradesUpdate(e types.WebsocketEvent) {

}

func (c *Client) handleOHLCVInit(e types.WebsocketEvent) {

}

func (c *Client) handleOHLCVUpdate(e types.WebsocketEvent) {

}

func (c *Client) SetNonce(o *types.Order) {
	o.Nonce = big.NewInt(int64(c.NonceGenerator.Intn(1e8)))
}
