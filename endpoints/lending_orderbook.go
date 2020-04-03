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

// LendingOrderBookEndpoint struct for lending
type LendingOrderBookEndpoint struct {
	lendingOrderBookService interfaces.LendingOrderBookService
}

// ServeLendingOrderBookResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeLendingOrderBookResource(
	r *mux.Router,
	lendingOrderBookService interfaces.LendingOrderBookService,
) {
	e := &LendingOrderBookEndpoint{lendingOrderBookService}
	r.HandleFunc("/api/lending/orderbook", e.HandleGetLendingOrderBook).Methods("GET")
	r.HandleFunc("/api/lending/orderbook/db", e.HandleGetLendingOrderBookInDb).Methods("GET")
	ws.RegisterChannel(ws.LendingOrderBookChannel, e.lendingOrderBookWebSocket)
}

func (e *LendingOrderBookEndpoint) HandleGetLendingOrderBook(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	t := v.Get("term")
	lendingToken := v.Get("lendingToken")
	term, err := strconv.ParseUint(t, 10, 32)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Term parameter is incorect")
		return
	}

	if lendingToken == "" {
		httputils.WriteError(w, http.StatusBadRequest, "lendingToken Parameter missing")
		return
	}

	if !common.IsHexAddress(lendingToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	lendingTokenAddress := common.HexToAddress(lendingToken)
	ob, err := e.lendingOrderBookService.GetLendingOrderBook(term, lendingTokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

func (e *LendingOrderBookEndpoint) HandleGetLendingOrderBookInDb(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	t := v.Get("term")
	lendingToken := v.Get("lendingToken")
	term, err := strconv.ParseUint(t, 10, 32)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Term parameter is incorect")
		return
	}

	if lendingToken == "" {
		httputils.WriteError(w, http.StatusBadRequest, "lendingToken Parameter missing")
		return
	}

	if !common.IsHexAddress(lendingToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	lendingTokenAddress := common.HexToAddress(lendingToken)
	ob, err := e.lendingOrderBookService.GetLendingOrderBookInDb(term, lendingTokenAddress)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, ob)
}

func (e *LendingOrderBookEndpoint) lendingOrderBookWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
		return
	}
	socket := ws.GetLendingOrderBookSocket()
	if ev == nil {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}
	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload

	err = json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		msg := map[string]string{"Message": "Internal server error"}
		socket.SendErrorMessage(c, msg)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		if p == nil {
			socket.SendErrorMessage(c, errInvalidPayload)
			return
		}

		if (p.LendingToken == common.Address{}) {
			msg := map[string]string{"Message": "Invalid lending token"}
			socket.SendErrorMessage(c, msg)
			return
		}

		e.lendingOrderBookService.SubscribeLendingOrderBook(c, p.Term, p.LendingToken)
	}

	if ev.Type == types.UNSUBSCRIBE {
		if p == nil {
			e.lendingOrderBookService.UnsubscribeLendingOrderBook(c)
			return
		}

		e.lendingOrderBookService.UnsubscribeLendingOrderBookChannel(c, p.Term, p.LendingToken)
	}
}
