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

// LendingMarketsEndpoint lending market endpoint
type LendingMarketsEndpoint struct {
	LendingMarketsService interfaces.LendingMarketsService
	LendingOhlcvService   interfaces.LendingOhlcvService
}

// ServeLendingMarketsResource sets up the routing of token endpoints and the corresponding handlers.
func ServeLendingMarketsResource(
	r *mux.Router,
	lendingMarketsService interfaces.LendingMarketsService,
	lendingOhlcvService interfaces.LendingOhlcvService,
) {
	e := &LendingMarketsEndpoint{lendingMarketsService, lendingOhlcvService}
	r.HandleFunc("/api/lending/market/stats/all", e.handleGetAllLendingMarketStats).Methods("GET")
	r.HandleFunc("/api/lending/market/stats", e.handleGetLendingMarketStats).Methods("GET")

	ws.RegisterChannel(ws.LendingMarketsChannel, e.handleLendingMarketsWebSocket)
}

// handleGetAllLendingMarketStats get all market token data
func (e *LendingMarketsEndpoint) handleGetAllLendingMarketStats(w http.ResponseWriter, r *http.Request) {

	res, err := e.LendingOhlcvService.GetAllTokenPairData()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.LendingTick{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
	return

}

// handleGetLendingMarketStats get market specific token data
func (e *LendingMarketsEndpoint) handleGetLendingMarketStats(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	t := v.Get("term")
	lendingToken := v.Get("lendingToken")

	if lendingToken == "" {
		httputils.WriteError(w, http.StatusBadRequest, "lendingToken parameter missing")
		return
	}

	if t == "" {
		httputils.WriteError(w, http.StatusBadRequest, "term parameter missing")
		return
	}

	if !common.IsHexAddress(lendingToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid lendingToken Address")
		return
	}
	term, err := strconv.ParseUint(t, 10, 32)
	if err != nil {
		httputils.WriteError(w, http.StatusBadRequest, "term parameter is incorect")
		return
	}
	lendingTokenAddress := common.HexToAddress(lendingToken)
	res := e.LendingOhlcvService.GetTokenPairData(term, lendingTokenAddress)
	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, types.LendingTick{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *LendingMarketsEndpoint) handleLendingMarketsWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetLendingMarketSocket()

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		e.LendingMarketsService.Subscribe(c)
	}

	if ev.Type == types.UNSUBSCRIBE {
		e.LendingMarketsService.Unsubscribe(c)
	}
}
