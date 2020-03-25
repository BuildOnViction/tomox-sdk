package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type lendingTradeEndpoint struct {
	lendingTradeService interfaces.LendingTradeService
}

// ServeLendingTradeResource sets up the routing of trade endpoints and the corresponding handlers.
// TODO trim down to one single endpoint with the 3 following params: base, quote, address
func ServeLendingTradeResource(
	r *mux.Router,
	lendingTradeService interfaces.LendingTradeService,
) {
	e := &lendingTradeEndpoint{lendingTradeService}
	r.HandleFunc("/api/lending/trades", e.handleGetLendingTrades).Methods("GET")
	r.HandleFunc("/api/lending/trades/history", e.handleGetLendingTradesHistory).Methods("GET")
	ws.RegisterChannel(ws.LendingTradeChannel, e.lendingTradeWebsocket)
}
func (e *lendingTradeEndpoint) lendingTradeWebsocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	if err := json.Unmarshal(b, &ev); err != nil {
		logger.Error(err)
		return
	}
	socket := ws.GetLendingTradeSocket()
	if ev == nil {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}
	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload
	err := json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		if p == nil {
			socket.SendErrorMessage(c, errInvalidPayload)
			return
		}
		if p.Term == 0 {
			err := map[string]string{"Message": "Invalid base token"}
			socket.SendErrorMessage(c, err)
			return
		}

		if (p.LendingToken == common.Address{}) {
			err := map[string]string{"Message": "Invalid lending token"}
			socket.SendErrorMessage(c, err)
			return
		}

		e.lendingTradeService.Subscribe(c, p.Term, p.LendingToken)
	}

	if ev.Type == types.UNSUBSCRIBE {
		if p == nil {
			e.lendingTradeService.Unsubscribe(c)
			return
		}

		e.lendingTradeService.UnsubscribeChannel(c, p.Term, p.LendingToken)
	}
}

// handleGetLendingTradesHistory is responsible for handling user's trade history requests
func (e *lendingTradeEndpoint) handleGetLendingTradesHistory(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")
	term := v.Get("term")
	lt := v.Get("lendingToken")
	status := v.Get("status")
	fromParam := v.Get("from")
	toParam := v.Get("to")

	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	if sortBy == "" {
		sortBy = "time"
	}

	var lendingTradeSpec types.LendingTradeSpec

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}
	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	if lt != "" {
		if !common.IsHexAddress(lt) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lending token address")
			return
		} else {
			lendingTradeSpec.LendingToken = common.HexToAddress(lt).Hex()
		}
	}

	if toParam != "" {
		t, _ := strconv.Atoi(toParam)
		lendingTradeSpec.DateTo = int64(t)
	}
	if fromParam != "" {
		t, _ := strconv.Atoi(fromParam)
		lendingTradeSpec.DateFrom = int64(t)
	}
	if term != "" {
		_, err := strconv.Atoi(term)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
			return
		}
		lendingTradeSpec.Term = term
	}
	if status != "" {
		lendingTradeSpec.Status = status
	}
	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}
	if sortType != "asc" && sortType != "dec" {
		sortType = "dec"
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

	res, err := e.lendingTradeService.GetLendingTradesUserHistory(common.HexToAddress(addr), &lendingTradeSpec, sortDB, offset*size, size)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Trade{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)

}
func (e *lendingTradeEndpoint) handleGetLendingTrades(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	lendingToken := v.Get("lendingToken")
	term := v.Get("term")
	fromParam := v.Get("from")
	toParam := v.Get("to")
	status := v.Get("status")
	pageOffset := v.Get("pageOffset")
	pageSize := v.Get("pageSize")
	sortBy := v.Get("sortBy")
	sortType := v.Get("sortType")

	sortedList := make(map[string]string)
	sortedList["time"] = "createdAt"
	if sortBy == "" {
		sortBy = "time"
	}

	var lendingTradeSpec types.LendingTradeSpec
	if lendingToken != "" {
		if !common.IsHexAddress(lendingToken) {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken")
			return
		} else {
			lendingTradeSpec.LendingToken = common.HexToAddress(lendingToken).Hex()
		}
	}

	if term != "" {
		_, err := strconv.Atoi(term)
		if err != nil {
			httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
			return
		}
		lendingTradeSpec.Term = term
	}
	if status != "" {
		lendingTradeSpec.Status = status
	}
	if toParam != "" {
		t, _ := strconv.Atoi(toParam)
		lendingTradeSpec.DateTo = int64(t)
	}
	if fromParam != "" {
		t, _ := strconv.Atoi(fromParam)
		lendingTradeSpec.DateFrom = int64(t)
	}

	offset := 0
	size := types.DefaultLimit
	sortDB := []string{}
	if sortType != "asc" && sortType != "dec" {
		sortType = "dec"
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

	res, err := e.lendingTradeService.GetLendingTrades(&lendingTradeSpec, sortDB, offset*size, size)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if res == nil {
		r := types.LendingTradeRes{
			Total:         0,
			LendingTrades: []*types.LendingTrade{},
		}
		httputils.WriteJSON(w, http.StatusOK, r)
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
