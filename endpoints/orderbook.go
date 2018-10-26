package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/utils/httputils"
	"github.com/tomochain/backend-matching-engine/ws"
)

type OrderBookEndpoint struct {
	orderBookService interfaces.OrderBookService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeOrderBookResource(
	r *mux.Router,
	orderBookService interfaces.OrderBookService,
) {
	e := &OrderBookEndpoint{orderBookService}
	r.HandleFunc("/orderbook/{baseToken}/{quoteToken}/raw", e.handleGetRawOrderBook)
	r.HandleFunc("/orderbook/{baseToken}/{quoteToken}/", e.handleGetOrderBook)
	ws.RegisterChannel(ws.LiteOrderBookChannel, e.orderBookWebSocket)
	ws.RegisterChannel(ws.RawOrderBookChannel, e.rawOrderBookWebSocket)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) handleGetOrderBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bt := vars["baseToken"]
	qt := vars["quoteToken"]

	if !common.IsHexAddress(bt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

// orderBookEndpoint
func (e *OrderBookEndpoint) handleGetRawOrderBook(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bt := vars["baseToken"]
	qt := vars["quoteToken"]

	if !common.IsHexAddress(bt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	if !common.IsHexAddress(qt) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	baseTokenAddress := common.HexToAddress(bt)
	quoteTokenAddress := common.HexToAddress(qt)
	ob, err := e.orderBookService.GetRawOrderBook(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		httputils.WriteError(w, http.StatusInternalServerError, "")
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) rawOrderBookWebSocket(input interface{}, conn *ws.Conn) {
	mab, _ := json.Marshal(input)
	var payload *types.WebSocketPayload

	err := json.Unmarshal(mab, &payload)
	if err != nil {
		logger.Error(err)
		return
	}

	socket := ws.GetRawOrderBookSocket()

	if payload.Type != "subscription" {
		logger.Error("Payload is not of subscription type")
		socket.SendErrorMessage(conn, "Payload is not of subscription type")
		return
	}

	dab, _ := json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(dab, &msg)
	if err != nil {
		logger.Error(err)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid Base Token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid Quote Token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeRawOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnSubscribeRawOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}

// liteOrderBookWebSocket
func (e *OrderBookEndpoint) orderBookWebSocket(input interface{}, conn *ws.Conn) {
	bytes, _ := json.Marshal(input)
	var payload *types.WebSocketPayload
	err := json.Unmarshal(bytes, &payload)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetOrderBookSocket()
	if payload.Type != "subscription" {
		message := map[string]string{"Message": "Invalid subscription payload"}
		socket.SendErrorMessage(conn, message)
		return
	}

	bytes, _ = json.Marshal(payload.Data)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(bytes, &msg)
	if err != nil {
		logger.Error(err)
		message := map[string]string{"Message": "Internal server error"}
		socket.SendErrorMessage(conn, message)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid base token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		message := map[string]string{"Message": "Invalid quote token"}
		socket.SendErrorMessage(conn, message)
		return
	}

	if msg.Event == types.SUBSCRIBE {
		e.orderBookService.SubscribeOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.orderBookService.UnSubscribeOrderBook(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken)
	}
}
