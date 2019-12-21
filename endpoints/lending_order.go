package endpoints

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type lendingorderEndpoint struct {
	lendingorderService interfaces.LendingOrderService
}

// ServeLendingOrderResource sets up the routing of order endpoints and the corresponding handlers.
func ServeLendingOrderResource(r *mux.Router, lendingorderService interfaces.LendingOrderService) {
	e := &lendingorderEndpoint{lendingorderService}

	r.HandleFunc("/api/lending/nonce", e.handleGetLendingOrderNonce).Methods("GET")
	r.HandleFunc("/api/lending", e.handleNewLendingOrder).Methods("POST")
	r.HandleFunc("/api/lending/cancel", e.handleCancelLendingOrder).Methods("POST")
	ws.RegisterChannel(ws.LendingOrderChannel, e.ws)
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
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusCreated, o)
}

func (e *lendingorderEndpoint) handleCancelLendingOrder(w http.ResponseWriter, r *http.Request) {
	oc := &types.LendingOrderCancel{}

	decoder := json.NewDecoder(r.Body)

	defer r.Body.Close()

	err := decoder.Decode(&oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, "Invalid payload")
		return
	}
	logger.Info("handle cancel order nonce", oc.Nonce)
	err = e.lendingorderService.CancelLendingOrder(oc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, oc.Hash)
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
	default:
		log.Print("Response with error")
	}
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
	oc := &types.LendingOrderCancel{}

	err = json.Unmarshal(bytes, &oc)
	if err != nil {
		logger.Error(err)
		c.SendLendingOrderErrorMessage(err, oc.Hash)
	}

	ws.RegisterLendingOrderConnection(oc.UserAddress, c)

	orderErr := e.lendingorderService.CancelLendingOrder(oc)
	if orderErr != nil {
		logger.Error(orderErr)
		c.SendLendingOrderErrorMessage(orderErr, oc.Hash)
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
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, total)
}
