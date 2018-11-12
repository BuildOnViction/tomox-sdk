package endpoints

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/rlp"
	"github.com/ethereum/go-ethereum/swarm/storage/feed"
	"github.com/ethereum/go-ethereum/swarm/storage/feed/lookup"
	"github.com/gin-gonic/gin"
	"github.com/tomochain/backend-matching-engine/ethereum"
	"github.com/tomochain/backend-matching-engine/interfaces"

	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/ws"
)

const TopicName = "Token"

type orderEndpoint struct {
	orderService interfaces.OrderService
	engine       interfaces.Engine
}

// ServeOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(
	r *gin.Engine,
	orderService interfaces.OrderService,
	engine interfaces.Engine,
) {
	e := &orderEndpoint{orderService, engine}

	r.GET("/orders/:address", e.handleGetOrders)
	r.GET("/orders/:address/:action", e.handleGetOrdersAction)

	// r.POST("/orders/:address/:action/encode", e.handleEncode)

	ws.RegisterChannel(ws.OrderChannel, e.ws)
}

func (e *orderEndpoint) getEngineProvider() *ethereum.EthereumProvider {
	return e.engine.Provider().(*ethereum.EthereumProvider)
}

func (e *orderEndpoint) handleGetOrdersAction(c *gin.Context) {
	action := c.Param("action")

	switch action {
	case "history":
		e.handleGetOrderHistory(c)
	case "current":
		e.handleGetPositions(c)
	default:
		e.handleGetOrdersFromPss(c)
	}

}

func (e *orderEndpoint) handleGetOrdersFromPss(c *gin.Context) {
	coin := c.Param("action")
	address := c.Param("address")
	topic, _ := feed.NewTopic(TopicName, []byte(coin))
	fd := &feed.Feed{
		Topic: topic,
		User:  common.HexToAddress(address),
	}

	lookupParams := feed.NewQueryLatest(fd, lookup.NoClue)

	bzzClient := e.getEngineProvider().BzzClient
	reader, err := bzzClient.QueryFeed(lookupParams, "")

	if err != nil {
		c.JSON(http.StatusBadRequest, GinError(fmt.Sprintf("Error retrieving feed updates: %s", err)))
		return
	}
	defer reader.Close()
	databytes, err := ioutil.ReadAll(reader)

	if databytes == nil || err != nil {
		c.JSON(http.StatusBadRequest, GinError(fmt.Sprintf("Error retrieving feed updates: %s", err)))
		return
	}

	// // try to decode
	var feeds []types.OrderFeed
	err = rlp.DecodeBytes(databytes, &feeds)
	if err != nil {
		c.JSON(http.StatusBadRequest, GinError(fmt.Sprintf("Error retrieving feed updates: %s", err)))
		return
	}

	var messages []*types.OrderRecord
	for _, feed := range feeds {
		message, _ := feed.GetBSON()
		messages = append(messages, message)
	}
	c.JSON(http.StatusOK, messages)

}

// func (e *orderEndpoint) handleEncode(c *gin.Context) {

// 	// deserialize order, in the future can use cargo to batch
// 	// return byte update to client to update, even it can not extend chunk size, the actual storage size is not chunk size
// 	// Topic "Tomo" will give us information like user type, how many slots corresponding to IDs
// 	var msg = &protocol.OrderbookMsg{}

// 	c.BindJSON(msg)

// 	coin := c.Param("action")
// 	addr := c.Param("address")
// 	rpcClient := e.getRPCClient()
// 	var messages []*protocol.OrderbookMsg
// 	rpcClient.Call(&messages, "orderbook_getOrders", coin, addr)

// 	log.Printf("order results: %s", messages)
// 	// on server, try update if fail then create
// 	if messages == nil {
// 		messages = []*protocol.OrderbookMsg{msg}
// 	} else {
// 		// find item if found then append, else update
// 		var found = false
// 		for i, message := range messages {
// 			if message.ID == msg.ID {
// 				found = true
// 				messages[i] = msg
// 				break
// 			}
// 		}
// 		if !found {
// 			messages = append(messages, msg)
// 		}
// 	}

// 	data, _ := rlp.EncodeToBytes(messages)
// 	c.Data(http.StatusOK, "application/octet-stream", data)
// }

func (e *orderEndpoint) handleGetOrders(c *gin.Context) {

	addr := c.Param("address")
	if !common.IsHexAddress(addr) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
	}

	address := common.HexToAddress(addr)
	orders, err := e.orderService.GetByUserAddress(address)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
	}

	c.JSON(http.StatusOK, orders)
}

func (e *orderEndpoint) handleGetPositions(c *gin.Context) {
	addr := c.Param("address")

	if !common.IsHexAddress(addr) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
	}

	address := common.HexToAddress(addr)
	orders, err := e.orderService.GetCurrentByUserAddress(address)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
	}

	c.JSON(http.StatusOK, orders)
}

func (e *orderEndpoint) handleGetOrderHistory(c *gin.Context) {
	addr := c.Param("address")

	if !common.IsHexAddress(addr) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
	}

	address := common.HexToAddress(addr)
	orders, err := e.orderService.GetHistoryByUserAddress(address)
	if err != nil {
		c.JSON(http.StatusInternalServerError, GinError(""))
	}

	c.JSON(http.StatusOK, orders)
}

// ws function handles incoming websocket messages on the order channel
func (e *orderEndpoint) ws(input interface{}, conn *ws.Conn) {
	msg, ok := input.(*types.WebsocketEvent)
	if ok {
		switch msg.Type {
		case "NEW_ORDER":
			e.handleNewOrder(msg, conn)
		case "CANCEL_ORDER":
			e.handleCancelOrder(msg, conn)
		case "SUBMIT_SIGNATURE":
			e.handleSubmitSignatures(msg, conn)
		default:
			log.Printf("Unknown event type: %s", msg.Type)
		}
	}
}

// handleSubmitSignatures handles NewTrade messages. New trade messages are transmitted to the corresponding order channel
// and received in the handleClientResponse.
func (e *orderEndpoint) handleSubmitSignatures(p *types.WebsocketEvent, conn *ws.Conn) {
	hash := common.HexToHash(p.Hash)
	// get order channel return the channel of the order by its hash, waiting for data to be updated
	ch := ws.GetOrderChannel(hash)

	if ch != nil {
		ch <- p
	}
}

// handleNewOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *orderEndpoint) handleNewOrder(msg *types.WebsocketEvent, conn *ws.Conn) {

	o := &types.Order{}

	bytes, err := json.Marshal(msg.Payload)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
		return
	}

	err = json.Unmarshal(bytes, &o)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
		return
	}

	o.Hash = o.ComputeHash()

	ch := make(chan *types.WebsocketEvent)
	ws.RegisterOrderConnection(o.Hash, &ws.OrderConnection{Conn: conn, ReadChannel: ch})
	// ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketUnsubscribeHandler(o.Hash))

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
		return
	}

	// send hash to client, then client will submit signature
}

// handleCancelOrder handles CancelOrder message.
func (e *orderEndpoint) handleCancelOrder(event *types.WebsocketEvent, conn *ws.Conn) {
	bytes, err := json.Marshal(event.Payload)
	oc := &types.OrderCancel{}

	err = oc.UnmarshalJSON(bytes)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
	}

	ws.RegisterOrderConnection(oc.Hash, &ws.OrderConnection{Conn: conn, Active: true})
	// ws.RegisterConnectionUnsubscribeHandler(
	// 	conn,
	// 	ws.OrderSocketUnsubscribeHandler(oc.Hash),
	// )

	err = e.orderService.CancelOrder(oc)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
		return
	}
}
