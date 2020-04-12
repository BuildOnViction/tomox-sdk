package endpoints

import (
	"encoding/json"
	"log"
	"math/big"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

const (
	COLLATERAL_RATE = 1.5
)

type lendingorderEndpoint struct {
	lendingorderService interfaces.LendingOrderService
}

// ServeLendingOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeLendingOrderResource(r *mux.Router, lendingorderService interfaces.LendingOrderService) {
	e := &lendingorderEndpoint{lendingorderService}
	r.HandleFunc("/api/lending/orders", e.handleGetLendingOrders).Methods("GET")
	r.HandleFunc("/api/lending/repay", e.handleGetRepay).Methods("GET")
	r.HandleFunc("/api/lending/topup", e.handleGetTopup).Methods("GET")
	r.HandleFunc("/api/lending/recall", e.handleGetRecall).Methods("GET")
	r.HandleFunc("/api/lending/estimate", e.handleGetEstimateCollateral).Methods("GET")
	r.HandleFunc("/api/lending/nonce", e.handleGetLendingOrderNonce).Methods("GET")
	r.HandleFunc("/api/lending", e.handleNewLendingOrder).Methods("POST")
	r.HandleFunc("/api/lending/cancel", e.handleCancelLendingOrder).Methods("POST")
	r.HandleFunc("/api/lending/repay", e.handleRepayLendingOrder).Methods("POST")
	r.HandleFunc("/api/lending/topup", e.handleTopupLendingOrder).Methods("POST")
	r.HandleFunc("/api/lending/{hash}", e.handleLendingByHash).Methods("GET")

	ws.RegisterChannel(ws.LendingOrderChannel, e.ws)
}
func (e *lendingorderEndpoint) handleGetLendingOrders(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	lendingToken := v.Get("lendingToken")
	collateralToken := v.Get("collateralToken")
	term := v.Get("term")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")
	side := v.Get("lendingSide")
	status := v.Get("lendingStatus")
	lendingType := v.Get("lendingType")
	lendinghash := v.Get("hash")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	sortedList["lendingStatus"] = "status"
	sortedList["lendingType"] = "type"
	sortedList["lendingSide"] = "side"

	var lendingSpec types.LendingSpec
	if addr != "" {
		if !common.IsHexAddress(addr) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
			return
		}
		lendingSpec.UserAddress = common.HexToAddress(addr).Hex()

	}

	if lendingToken != "" {
		if !common.IsHexAddress(lendingToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken Address")
			return
		}
		lendingSpec.LendingToken = common.HexToAddress(lendingToken).Hex()
	}

	if collateralToken != "" {
		if !common.IsHexAddress(collateralToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
			return
		}
		lendingSpec.CollateralToken = common.HexToAddress(collateralToken).Hex()
	}

	if term != "" {
		_, err := strconv.Atoi(term)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
			return
		}
		lendingSpec.Term = term
	}

	if fromParam != "" {
		t, err := strconv.Atoi(fromParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time from value")
			return
		}
		lendingSpec.DateFrom = int64(t)
	}

	if toParam != "" {
		t, err := strconv.Atoi(toParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time to value")
			return
		}
		lendingSpec.DateTo = int64(t)
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}

	if sortType != "asc" && sortType != "dec" {
		sortType = "dec"
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
		lendingSpec.Side = side
	}
	if status != "" {
		lendingSpec.Status = status
	}
	if lendinghash != "" {
		lendingSpec.Hash = lendinghash
	}
	if lendingType != "" {
		lendingSpec.Type = lendingType
	}
	var err error
	var lendings *types.LendingRes

	lendings, err = e.lendingorderService.GetLendingOrders(lendingSpec, sortDB, offset*size, size)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if lendings == nil {
		r := types.LendingRes{
			Total:        0,
			LendingItems: []*types.LendingOrder{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, lendings)
}

func (e *lendingorderEndpoint) handleLendingByHash(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	lendingHash := vars["hash"]
	res, err := e.lendingorderService.GetByHash(common.HexToHash(lendingHash))

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, res)
}
func (e *lendingorderEndpoint) handleNewLendingOrder(w http.ResponseWriter, r *http.Request) {
	var o *types.LendingOrder
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	o.Hash = o.ComputeHash()
	err = e.lendingorderService.NewLendingOrder(o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusCreated, o)
}

func (e *lendingorderEndpoint) handleCancelLendingOrder(w http.ResponseWriter, r *http.Request) {
	o := &types.LendingOrder{}
	decoder := json.NewDecoder(r.Body)
	defer r.Body.Close()
	err := decoder.Decode(&o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	err = e.lendingorderService.CancelLendingOrder(o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, o.Hash)
}

func (e *lendingorderEndpoint) handleRepayLendingOrder(w http.ResponseWriter, r *http.Request) {
	o := &types.LendingOrder{}
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	err = e.lendingorderService.RepayLendingOrder(o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, o.Hash)
}
func (e *lendingorderEndpoint) handleTopupLendingOrder(w http.ResponseWriter, r *http.Request) {
	o := &types.LendingOrder{}
	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	err = e.lendingorderService.TopupLendingOrder(o)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, o.Hash)
}

// ws function handles incoming websocket messages on the order channel
func (e *lendingorderEndpoint) ws(input interface{}, c *ws.Client) {
	msg := &types.WebsocketEvent{}

	bytes, _ := json.Marshal(input)
	if err := json.Unmarshal(bytes, &msg); err != nil {
		logger.Error(err)
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, err.Error())
	}

	switch msg.Type {
	case "NEW_LENDING_ORDER":
		e.handleWSNewLendingOrder(msg, c)
	case "CANCEL_LENDING_ORDER":
		e.handleWSCancelLendingOrder(msg, c)
	case "SUBSCRIBE":
		e.handleWSSubLendingOrder(msg, c)
	default:
		log.Print("Response with error")
	}
}

func (e *lendingorderEndpoint) handleWSSubLendingOrder(ev *types.WebsocketEvent, c *ws.Client) {
	var addr string
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	bytes, err := json.Marshal(ev.Payload)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, err.Error())
		return
	}

	logger.Debugf("Payload: %v#", ev.Payload)

	err = json.Unmarshal(bytes, &addr)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, err.Error())
		return
	}

	if !common.IsHexAddress(addr) {
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, errInvalidPayload)
		return
	}

	a := common.HexToAddress(addr)
	ws.RegisterLendingOrderConnection(a, c)
	ws.SendLendingOrderMessage(types.INIT, a, nil)
}

// handleWSNewLendingOrder handles NewOrder message. New order messages are transmitted to the order service after being unmarshalled
func (e *lendingorderEndpoint) handleWSNewLendingOrder(ev *types.WebsocketEvent, c *ws.Client) {
	o := &types.LendingOrder{}
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	bytes, err := json.Marshal(ev.Payload)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, err.Error())
		return
	}

	logger.Debugf("Payload: %v#", ev.Payload)

	err = json.Unmarshal(bytes, &o)
	if err != nil {
		logger.Error(err)
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, err.Error())
		return
	}
	if o == nil {
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, errInvalidPayload)
		return
	}
	if err := o.Validate(); err != nil {
		c.SendMessage(ws.LendingOrderChannel, types.ERROR, err.Error())
		return
	}

	o.Hash = o.ComputeHash()
	ws.RegisterLendingOrderConnection(o.UserAddress, c)

	err = e.lendingorderService.NewLendingOrder(o)
	if err != nil {
		logger.Error(err)
		c.SendLendingOrderErrorMessage(err, o.Hash)
		return
	}
}

// handleCancelLendingOrder handles CancelLendingOrder message.
func (e *lendingorderEndpoint) handleWSCancelLendingOrder(ev *types.WebsocketEvent, c *ws.Client) {
	bytes, err := json.Marshal(ev.Payload)
	o := &types.LendingOrder{}

	err = json.Unmarshal(bytes, &o)
	if err != nil {
		logger.Error(err)
		c.SendLendingOrderErrorMessage(err, o.Hash)
	}

	ws.RegisterLendingOrderConnection(o.UserAddress, c)

	orderErr := e.lendingorderService.CancelLendingOrder(o)
	if orderErr != nil {
		logger.Error(orderErr)
		c.SendLendingOrderErrorMessage(orderErr, o.Hash)
		return
	}
}

func (e *lendingorderEndpoint) handleGetLendingOrderNonce(w http.ResponseWriter, r *http.Request) {
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

	total, err := e.lendingorderService.GetLendingNonceByUserAddress(a)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, total)
}

func (e *lendingorderEndpoint) handleGetTopup(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	lendingToken := v.Get("lendingToken")
	collateralToken := v.Get("collateralToken")
	term := v.Get("term")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")
	lendinghash := v.Get("hash")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	sortedList["term"] = "term"

	var topupSpec types.TopupSpec
	if addr != "" {
		if !common.IsHexAddress(addr) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
			return
		}
		topupSpec.UserAddress = common.HexToAddress(addr).Hex()

	}

	if lendingToken != "" {
		if !common.IsHexAddress(lendingToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken address")
			return
		}
		topupSpec.LendingToken = common.HexToAddress(lendingToken).Hex()
	}

	if collateralToken != "" {
		if !common.IsHexAddress(collateralToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid colateralToken address")
			return
		}
		topupSpec.CollateralToken = common.HexToAddress(collateralToken).Hex()
	}

	if term != "" {
		_, err := strconv.Atoi(term)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
			return
		}
		topupSpec.Term = term
	}

	if fromParam != "" {
		t, err := strconv.Atoi(fromParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time from value")
			return
		}
		topupSpec.DateFrom = int64(t)
	}

	if toParam != "" {
		t, err := strconv.Atoi(toParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time to value")
			return
		}
		topupSpec.DateTo = int64(t)
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}

	if sortType != "asc" && sortType != "dec" {
		sortType = "dec"
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

	if lendinghash != "" {
		topupSpec.Hash = lendinghash
	}

	var err error
	var lendings *types.LendingRes

	lendings, err = e.lendingorderService.GetTopup(topupSpec, sortDB, offset*size, size)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if lendings == nil {
		r := types.LendingRes{
			Total:        0,
			LendingItems: []*types.LendingOrder{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, lendings)
}

func (e *lendingorderEndpoint) handleGetRepay(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	lendingToken := v.Get("lendingToken")
	term := v.Get("term")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")
	lendinghash := v.Get("hash")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	sortedList["term"] = "term"

	var repaySpec types.RepaySpec
	if addr != "" {
		if !common.IsHexAddress(addr) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
			return
		}
		repaySpec.UserAddress = common.HexToAddress(addr).Hex()

	}

	if lendingToken != "" {
		if !common.IsHexAddress(lendingToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken address")
			return
		}
		repaySpec.LendingToken = common.HexToAddress(lendingToken).Hex()
	}

	if term != "" {
		_, err := strconv.Atoi(term)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
			return
		}
		repaySpec.Term = term
	}

	if fromParam != "" {
		t, err := strconv.Atoi(fromParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time from value")
			return
		}
		repaySpec.DateFrom = int64(t)
	}

	if toParam != "" {
		t, err := strconv.Atoi(toParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time to value")
			return
		}
		repaySpec.DateTo = int64(t)
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}

	if sortType != "asc" && sortType != "dec" {
		sortType = "dec"
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

	if lendinghash != "" {
		repaySpec.Hash = lendinghash
	}

	var err error
	var lendings *types.LendingRes

	lendings, err = e.lendingorderService.GetRepay(repaySpec, sortDB, offset*size, size)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if lendings == nil {
		r := types.LendingRes{
			Total:        0,
			LendingItems: []*types.LendingOrder{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, lendings)
}
func (e *lendingorderEndpoint) handleGetRecall(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	lendingToken := v.Get("lendingToken")
	collateralToken := v.Get("collateralToken")
	term := v.Get("term")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")
	lendinghash := v.Get("hash")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	sortedList["term"] = "term"

	var recallSpec types.RecallSpec
	if addr != "" {
		if !common.IsHexAddress(addr) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
			return
		}
		recallSpec.UserAddress = common.HexToAddress(addr).Hex()

	}

	if lendingToken != "" {
		if !common.IsHexAddress(lendingToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken address")
			return
		}
		recallSpec.LendingToken = common.HexToAddress(lendingToken).Hex()
	}

	if collateralToken != "" {
		if !common.IsHexAddress(collateralToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid colateralToken address")
			return
		}
		recallSpec.CollateralToken = common.HexToAddress(collateralToken).Hex()
	}

	if term != "" {
		_, err := strconv.Atoi(term)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
			return
		}
		recallSpec.Term = term
	}

	if fromParam != "" {
		t, err := strconv.Atoi(fromParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time from value")
			return
		}
		recallSpec.DateFrom = int64(t)
	}

	if toParam != "" {
		t, err := strconv.Atoi(toParam)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid time to value")
			return
		}
		recallSpec.DateTo = int64(t)
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}

	if sortType != "asc" && sortType != "dec" {
		sortType = "dec"
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

	if lendinghash != "" {
		recallSpec.Hash = lendinghash
	}

	var err error
	var lendings *types.LendingRes

	lendings, err = e.lendingorderService.GetRecall(recallSpec, sortDB, offset*size, size)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	if lendings == nil {
		r := types.LendingRes{
			Total:        0,
			LendingItems: []*types.LendingOrder{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, lendings)
}

func (e *lendingorderEndpoint) handleGetEstimateCollateral(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	lendingToken := v.Get("lendingToken")
	collateralToken := v.Get("collateralToken")
	amount := v.Get("amount")
	type res struct {
		CollateralToken          common.Address `json:"collateralToken"`
		LendingToken             common.Address `json:"lendingToken"`
		LendingAmount            *big.Int       `json:"lendingAmount"`
		Rate                     float64        `json:"rate"`
		CollateralPrice          *big.Float     `json:"collateralPrice"`
		EstimateCollateralAmount *big.Float     `json:"estimateCollateralAmount"`
	}
	var result res
	result.Rate = COLLATERAL_RATE
	if lendingToken != "" {
		if !common.IsHexAddress(lendingToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken address")
			return
		}
		result.LendingToken = common.HexToAddress(lendingToken)
	}

	if collateralToken != "" {
		if !common.IsHexAddress(collateralToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid colateralToken address")
			return
		}
		result.CollateralToken = common.HexToAddress(collateralToken)
	}
	lendingAmount, ok := new(big.Int).SetString(amount, 10)
	if !ok {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid lending amount")
		return
	}
	result.LendingAmount = lendingAmount
	collateralAmount, colalteralPrice, err := e.lendingorderService.EstimateCollateral(common.HexToAddress(collateralToken), common.HexToAddress(lendingToken), lendingAmount)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Cant not estimate collateral amount")
		return
	}
	result.EstimateCollateralAmount = collateralAmount.Mul(collateralAmount, big.NewFloat(result.Rate))
	result.CollateralPrice = colalteralPrice
	httputils.WriteJSON(w, http.StatusOK, result)
}
