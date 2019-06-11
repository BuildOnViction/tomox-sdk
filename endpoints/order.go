package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomoxsdk/errors"
	"github.com/tomochain/tomoxsdk/interfaces"
	"github.com/tomochain/tomoxsdk/types"
	"github.com/tomochain/tomoxsdk/utils/httputils"
	"github.com/tomochain/tomoxsdk/ws"
)

type orderEndpoint struct {
	orderService   interfaces.OrderService
	accountService interfaces.AccountService
}

// ServeOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeOrderResource(
	r *mux.Router,
	orderService interfaces.OrderService,
	accountService interfaces.AccountService,
) {
	e := &orderEndpoint{orderService, accountService}

	r.HandleFunc("/orders/count", e.handleGetCountOrder).Methods("GET")
	r.HandleFunc("/orders/history", e.handleGetOrderHistory).Methods("GET")
	r.HandleFunc("/orders/positions", e.handleGetPositions).Methods("GET")
	r.HandleFunc("/orders", e.handleGetOrders).Methods("GET")
	r.HandleFunc("/orders", e.HandleNewOrder).Methods("POST")
	r.HandleFunc("/orders/cancel", e.HandleCancelOrder).Methods("POST")

	ws.RegisterChannel(ws.OrderChannel, e.ws)
}

func (e *orderEndpoint) handleGetCountOrder(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter Missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := common.HexToAddress(addr)

	total, err := e.orderService.GetOrderCountByUserAddress(a)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, total)
}

func (e *orderEndpoint) handleGetOrders(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	limit := v.Get("limit")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter Missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	var err error
	var orders []*types.Order
	address := common.HexToAddress(addr)

	if limit == "" {
		orders, err = e.orderService.GetByUserAddress(address)
	} else {
		lim, _ := strconv.Atoi(limit)
		orders, err = e.orderService.GetByUserAddress(address, lim)
	}

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if orders == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Order{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) handleGetPositions(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	limit := v.Get("limit")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	var err error
	var orders []*types.Order
	address := common.HexToAddress(addr)

	if limit == "" {
		orders, err = e.orderService.GetCurrentByUserAddress(address)
	} else {
		lim, _ := strconv.Atoi(limit)
		orders, err = e.orderService.GetCurrentByUserAddress(address, lim)
	}

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if orders == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Order{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) handleGetOrderHistory(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	limit := v.Get("limit")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	var err error
	var orders []*types.Order
	address := common.HexToAddress(addr)

	if limit == "" {
		orders, err = e.orderService.GetHistoryByUserAddress(address)
	} else {
		lim, _ := strconv.Atoi(limit)
		orders, err = e.orderService.GetHistoryByUserAddress(address, lim)
	}

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Internal Server Error")
		return
	}

	if orders == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Order{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) HandleNewOrder(w http.ResponseWriter, r *http.Request) {
	var o *types.Order
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	o.Hash = o.ComputeHash()

	acc, err := e.accountService.FindOrCreate(o.UserAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if acc.IsBlocked {
		httputils.WriteError(w, http.StatusForbidden, "Account is blocked")
		return
	}

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, o)
}

func (e *orderEndpoint) HandleCancelOrder(w http.ResponseWriter, r *http.Request) {
	oc := &types.OrderCancel{}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}

	_, err = oc.GetSenderAddress()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	err = e.orderService.CancelOrder(oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, oc.Hash)
}

// ws function handles incoming websocket messages on the order channel
func (e *orderEndpoint) ws(input interface{}, c *ws.Client) {
	msg := &types.WebsocketEvent{}

	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		logger.Error(err)
		c.SendMessage(ws.OrderChannel, types.ERROR, err.Error())
	}

	switch msg.Type {
	case "NEW_ORDER":
		e.handleWSNewOrder(msg, c)
	case "CANCEL_ORDER":
		e.handleWSCancelOrder(msg, c)
	default:
		log.Print("Response with error")
	}
}

// handleNewOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *orderEndpoint) handleWSNewOrder(ev *types.WebsocketEvent, c *ws.Client) {
	o := &types.Order{}

	bytes, err := json.Marshal(ev.Payload)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.OrderChannel, types.ERROR, err.Error())
		return
	}

	logger.Debugf("Payload: %v#", ev.Payload)

	err = json.Unmarshal(bytes, &o)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, o.Hash)
		return
	}

	o.Hash = o.ComputeHash()
	ws.RegisterOrderConnection(o.UserAddress, c)

	acc, err := e.accountService.FindOrCreate(o.UserAddress)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, o.Hash)
	}

	if acc.IsBlocked {
		c.SendMessage(ws.OrderChannel, types.ERROR, errors.New("Account is blocked"))
	}

	err = e.orderService.NewOrder(o)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, o.Hash)
		return
	}
}

// handleCancelOrder handles CancelOrder message.
func (e *orderEndpoint) handleWSCancelOrder(ev *types.WebsocketEvent, c *ws.Client) {
	bytes, err := json.Marshal(ev.Payload)
	oc := &types.OrderCancel{}

	err = json.Unmarshal(bytes, &oc)
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, oc.Hash)
	}

	addr, err := oc.GetSenderAddress()
	if err != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(err, oc.Hash)
	}

	ws.RegisterOrderConnection(addr, c)

	orderErr := e.orderService.CancelOrder(oc)
	if orderErr != nil {
		logger.Error(err)
		c.SendOrderErrorMessage(orderErr, oc.Hash)
		return
	}
}
