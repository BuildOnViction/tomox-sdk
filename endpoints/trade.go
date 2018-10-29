package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/ws"
)

type tradeEndpoint struct {
	tradeService interfaces.TradeService
}

// ServeTradeResource sets up the routing of trade endpoints and the corresponding handlers.
func ServeTradeResource(
	r *gin.Engine,
	tradeService interfaces.TradeService,
) {
	e := &tradeEndpoint{tradeService}
	r.GET("/trades/pair", e.HandleGetTradeHistory)
	r.GET("/trades", e.HandleGetTrades)
	ws.RegisterChannel(ws.TradeChannel, e.tradeWebSocket)
}

// history is reponsible for handling pair's trade history requests
func (e *tradeEndpoint) HandleGetTradeHistory(c *gin.Context) {
	bt := c.Query("baseToken")
	qt := c.Query("quoteToken")
	l := c.Query("length")

	if bt == "" {
		c.JSON(http.StatusBadRequest, GinError("baseToken Parameter missing"))
		return
	}

	if qt == "" {
		c.JSON(http.StatusBadRequest, GinError("quoteToken Parameter missing"))
		return
	}

	if !common.IsHexAddress(bt) {
		c.JSON(http.StatusBadRequest, GinError("Invalid base token address"))
		return
	}

	if !common.IsHexAddress(qt) {
		c.JSON(http.StatusBadRequest, GinError("Invalid quote token address"))
		return
	}
	length := 20
	if l != "" {
		length, _ = strconv.Atoi(l)
	}
	baseToken := common.HexToAddress(bt)
	quoteToken := common.HexToAddress(qt)
	res, err := e.tradeService.GetSortedTradesByDate(baseToken, quoteToken, length)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}
	if res == nil {
		c.JSON(http.StatusOK, []types.Trade{})
		return
	}

	c.JSON(http.StatusOK, res)
}

// get is reponsible for handling user's trade history requests
func (e *tradeEndpoint) HandleGetTrades(c *gin.Context) {
	addr := c.Query("address")

	if addr == "" {
		c.JSON(http.StatusBadRequest, GinError("address Parameter missing"))
		return
	}

	if !common.IsHexAddress(addr) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
		return
	}

	address := common.HexToAddress(addr)
	res, err := e.tradeService.GetByUserAddress(address)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (e *tradeEndpoint) tradeWebSocket(input interface{}, conn *ws.Conn) {
	bytes, _ := json.Marshal(input)
	var event *types.WebsocketEvent
	if err := json.Unmarshal(bytes, &event); err != nil {
		logger.Error(err)
	}

	socket := ws.GetTradeSocket()
	if event.Type != "subscription" {
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(conn, err)
		return
	}

	bytes, _ = json.Marshal(event.Payload)
	var msg *types.WebSocketSubscription
	err := json.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Error(err)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		err := map[string]string{"Message": "Invalid base token"}
		socket.SendErrorMessage(conn, err)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		err := map[string]string{"Message": "Invalid quote token"}
		socket.SendErrorMessage(conn, err)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.tradeService.Subscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.tradeService.Unsubscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
