package endpoints

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/errors"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
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

	r.HandleFunc("/api/orders/count", e.handleGetCountOrder).Methods("GET")
	r.HandleFunc("/api/orders/nonce", e.handleGetOrderNonce).Methods("GET")
	r.HandleFunc("/api/orders/history", e.handleGetOrderHistory).Methods("GET")
	r.HandleFunc("/api/orders/positions", e.handleGetPositions).Methods("GET")
	r.HandleFunc("/api/orders", e.handleGetOrders).Methods("GET")
	r.HandleFunc("/api/orders", e.handleNewOrder).Methods("POST")
	r.HandleFunc("/api/orders/cancel", e.handleCancelOrder).Methods("POST")
	r.HandleFunc("/api/orders/cancelAll", e.handleCancelAllOrders).Methods("POST")
	r.HandleFunc("/api/orders/balance/lock", e.handleGetLockedBalanceInOrder).Methods("GET")
	ws.RegisterChannel(ws.OrderChannel, e.ws)
}

func (e *orderEndpoint) handleGetLockedBalanceInOrder(w http.ResponseWriter, r *http.Request) {
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

	total, err := e.orderService.GetOrdersLockedBalanceByUserAddress(a)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, total)
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
	baseToken := v.Get("baseToken")
	quoteToken := v.Get("quoteToken")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")
	side := v.Get("orderSide")
	status := v.Get("orderStatus")
	orderType := v.Get("orderType")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	sortedList["orderStatus"] = "status"
	sortedList["orderType"] = "type"
	sortedList["orderSide"] = "side"

	var orderSpec types.OrderSpec
	if addr != "" {
		if !common.IsHexAddress(addr) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
			return
		}
		orderSpec.UserAddress = addr

	}

	if baseToken != "" {
		if !common.IsHexAddress(baseToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
			return
		}
		orderSpec.BaseToken = baseToken
	}

	if quoteToken != "" {
		if !common.IsHexAddress(quoteToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
			return
		}
		orderSpec.QuoteToken = quoteToken
	}

	if fromParam != "" {
		t, err := strconv.Atoi(fromParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time from value")
			return
		}
		orderSpec.DateFrom = int64(t)
	}

	if toParam != "" {
		t, err := strconv.Atoi(toParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time to value")
			return
		}
		orderSpec.DateTo = int64(t)
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}

	if sortType != "asc" && sortType != "desc" {
		sortType = "desc"
	}

	if sortBy == "" {
		sortBy = "time"
	}

	if val, ok := sortedList[sortBy]; ok {
		if sortType == "asc" {
			sortDB = append(sortDB, "+"+val)
		} else {
			sortDB = append(sortDB, "-"+val)
		}

	}

	if pageOffset != "" {
		t, err := strconv.Atoi(pageOffset)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid page offset")
			return
		}
		offset = t
	}
	if pageSize != "" {
		t, err := strconv.Atoi(pageSize)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid page size")
			return
		}
		size = t
	}
	if side != "" {
		orderSpec.Side = side
	}
	if status != "" {
		orderSpec.Status = status
	}
	if orderType != "" {
		orderSpec.OrderType = orderType
	}
	var err error
	var orders *types.OrderRes

	orders, err = e.orderService.GetOrders(orderSpec, sortDB, offset*size, size)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if orders == nil {
		r := types.OrderRes{
			Total:  0,
			Orders: []*types.Order{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
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
	baseToken := v.Get("baseToken")
	quoteToken := v.Get("quoteToken")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")
	side := v.Get("orderSide")
	status := v.Get("orderStatus")
	orderType := v.Get("orderType")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	sortedList["orderStatus"] = "status"
	sortedList["orderType"] = "type"
	sortedList["orderSide"] = "side"

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	var orderSpec types.OrderSpec
	orderSpec.UserAddress = addr
	if baseToken != "" {
		if !common.IsHexAddress(baseToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
			return
		}
		orderSpec.BaseToken = baseToken
	}

	if quoteToken != "" {
		if !common.IsHexAddress(quoteToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
			return
		}
		orderSpec.QuoteToken = quoteToken
	}

	if fromParam != "" {
		t, err := strconv.Atoi(fromParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time from value")
			return
		}
		orderSpec.DateFrom = int64(t)
	}

	if toParam != "" {
		t, err := strconv.Atoi(toParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time to value")
			return
		}
		orderSpec.DateTo = int64(t)
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}
	if sortType != "asc" && sortType != "dec" {
		sortType = "asc"
	}
	if sortBy != "" {
		if val, ok := sortedList[sortBy]; ok {
			if sortType == "asc" {
				sortDB = append(sortDB, "+"+val)
			} else {
				sortDB = append(sortDB, "-"+val)
			}

		}
	}
	if pageOffset != "" {
		t, err := strconv.Atoi(pageOffset)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid page offset")
			return
		}
		offset = t
	}
	if pageSize != "" {
		t, err := strconv.Atoi(pageSize)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid page size")
			return
		}
		size = t
	}
	if side != "" {
		orderSpec.Side = side
	}
	if status != "" {
		orderSpec.Status = status
	}
	if orderType != "" {
		orderSpec.OrderType = orderType
	}
	var err error
	var orders *types.OrderRes

	orders, err = e.orderService.GetOrders(orderSpec, sortDB, offset*size, size)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if orders == nil {
		r := types.OrderRes{
			Total:  0,
			Orders: []*types.Order{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, orders)
}

func (e *orderEndpoint) handleNewOrder(w http.ResponseWriter, r *http.Request) {
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

func (e *orderEndpoint) handleCancelOrder(w http.ResponseWriter, r *http.Request) {
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

// handleCancelAllOrder cancels all open/partial filled orders of an user address
func (e *orderEndpoint) handleCancelAllOrders(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "Address parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := common.HexToAddress(addr)

	err := e.orderService.CancelAllOrder(a)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
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
		logger.Error(orderErr)
		c.SendOrderErrorMessage(orderErr, oc.Hash)
		return
	}
}

func (e *orderEndpoint) handleGetOrderNonce(w http.ResponseWriter, r *http.Request) {
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

	total, err := e.orderService.GetOrderNonceByUserAddress(a)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	if total == nil {
		httputils.WriteError(w, http.StatusInternalServerError, "unknow error")
	}
	s := total.(string)
	s = strings.TrimPrefix(s, "0x")
	n, err := strconv.ParseUint(s, 16, 32)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, n)
}
