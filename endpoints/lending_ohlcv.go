package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

// LendingOhlcvEndpoint struct for lending ohlcv endpoint
type LendingOhlcvEndpoint struct {
	lendingOhlcvService interfaces.LendingOhlcvService
}

// ServeLendingOhlcvResource handle lending ohlcv api
func ServeLendingOhlcvResource(
	r *mux.Router,
	lendingOhlcvService interfaces.LendingOhlcvService,
) {
	e := &LendingOhlcvEndpoint{lendingOhlcvService}
	r.HandleFunc("/api/lending/ohlcv", e.handleGetLendingOhlcv).Methods("GET")
	ws.RegisterChannel(ws.LendingOhlcvChannel, e.ohlcvWebSocket)
}

func (e *LendingOhlcvEndpoint) handleGetLendingOhlcv(w http.ResponseWriter, r *http.Request) {
	var p types.OHLCVParams

	v := r.URL.Query()
	t := v.Get("term")
	lendingToken := v.Get("lendingToken")
	from := v.Get("from")
	to := v.Get("to")
	timeInterval := v.Get("timeInterval")

	if timeInterval == "" {
		httputils.WriteError(w, http.StatusBadRequest, "timeInterval Parameter is missing")
		return
	}
	unit, duration := processTimeInterval(timeInterval)

	p.Units = unit
	p.Duration = int64(duration)

	now := time.Now()

	if to == "" {
		p.To = now.Unix()
	} else {
		t, _ := strconv.Atoi(to)
		p.To = int64(t)
	}

	if from == "" {
		p.From = now.AddDate(-1, 0, 0).Unix()
	} else {
		f, _ := strconv.Atoi(from)
		p.From = int64(f)
	}

	if t == "" {
		httputils.WriteError(w, http.StatusBadRequest, "term Parameter is missing")
		return
	}

	if lendingToken == "" {
		httputils.WriteError(w, http.StatusBadRequest, "lendingToken Parameter is missing")
		return
	}

	if !common.IsHexAddress(lendingToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid lending token address")
		return
	}
	term, err := strconv.ParseUint(t, 10, 64)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid term")
		return
	}
	res, err := e.lendingOhlcvService.GetOHLCV(term, common.HexToAddress(lendingToken), p.Duration, p.Units, p.From, p.To)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []*types.Tick{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *LendingOhlcvEndpoint) ohlcvWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	socket := ws.GetLendingOhlcvSocket()
	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
		return
	}
	if ev == nil {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		b, _ = json.Marshal(ev.Payload)
		var p *types.SubscriptionPayload

		err = json.Unmarshal(b, &p)
		if err != nil {
			logger.Error(err)
			return
		}
		if p == nil {
			socket.SendErrorMessage(c, errInvalidPayload)
			return
		}
		if p.Term == 0 {
			socket.SendErrorMessage(c, "Invalid term")
			return
		}

		if (p.LendingToken == common.Address{}) {
			socket.SendErrorMessage(c, "Invalid Lending Token")
			return
		}

		now := time.Now()

		if p.From == 0 {
			p.From = now.AddDate(-1, 0, 0).Unix()
		}

		if p.To == 0 {
			p.To = now.Unix()
		}

		if p.Duration == 0 {
			p.Duration = 24
		}

		if p.Units == "" {
			p.Units = "hour"
		}

		e.lendingOhlcvService.Subscribe(c, p)
	}

	if ev.Type == types.UNSUBSCRIBE {
		e.lendingOhlcvService.Unsubscribe(c)
	}
}
