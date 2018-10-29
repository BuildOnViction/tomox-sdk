package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/services"

	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/ws"
)

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

	ws.RegisterChannel(ws.OrderChannel, e.ws)
}

func (e *orderEndpoint) handleGetOrdersAction(c *gin.Context) {
	// vars := mux.Vars(r)
	action := c.Param("action")

	if action == "history" {
		e.handleGetOrderHistory(c)
	} else if action == "current" {
		e.handleGetPositions(c)
	} else {
		e.handleGetOrdersFromPss(c)
	}

}

func (e *orderEndpoint) handleGetOrdersFromPss(c *gin.Context) {
	coin := c.Param("action")
	addr := c.Param("address")
	orderService := e.orderService.(*services.OrderService)
	rpcClient := orderService.Provider().RPCClient
	var orderResult interface{}
	rpcClient.Call(&orderResult, "orderbook_getOrders", coin, addr)

	c.JSON(http.StatusOK, orderResult)

}

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
			log.Print("Response with error")
		}
	}
}

// handleSubmitSignatures handles NewTrade messages. New trade messages are transmitted to the corresponding order channel
// and received in the handleClientResponse.
func (e *orderEndpoint) handleSubmitSignatures(p *types.WebsocketEvent, conn *ws.Conn) {
	hash := common.HexToHash(p.Hash)
	ch := ws.GetOrderChannel(hash)

	if ch != nil {
		ch <- p
	}
}

// handleNewOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *orderEndpoint) handleNewOrder(msg *types.WebsocketEvent, conn *ws.Conn) {
	ch := make(chan *types.WebsocketEvent)
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
	ws.RegisterOrderConnection(o.Hash, &ws.OrderConnection{Conn: conn, ReadChannel: ch})
	ws.RegisterConnectionUnsubscribeHandler(conn, ws.OrderSocketUnsubscribeHandler(o.Hash))

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
		return
	}
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
	ws.RegisterConnectionUnsubscribeHandler(
		conn,
		ws.OrderSocketUnsubscribeHandler(oc.Hash),
	)

	err = e.orderService.CancelOrder(oc)
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.OrderChannel, ws.ERROR, err.Error())
		return
	}
}
