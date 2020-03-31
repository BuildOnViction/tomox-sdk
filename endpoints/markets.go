package endpoints

import (
	"encoding/json"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
	"github.com/tomochain/tomox-sdk/ws"
)

type MarketsEndpoint struct {
	MarketsService interfaces.MarketsService
	OHLCVService   interfaces.OHLCVService
	RelayerService interfaces.RelayerService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeMarketsResource(
	r *mux.Router,
	marketsService interfaces.MarketsService,
	ohlcvService interfaces.OHLCVService,
	relayerService interfaces.RelayerService,
) {
	e := &MarketsEndpoint{marketsService, ohlcvService, relayerService}
	r.HandleFunc("/api/market/stats/all", e.HandleGetAllMarketStats).Methods("GET")
	r.HandleFunc("/api/market/stats", e.HandleGetMarketStats).Methods("GET")

	ws.RegisterChannel(ws.MarketsChannel, e.handleMarketsWebSocket)
}

// HandleGetAllMarketStats get all market token data
func (e *MarketsEndpoint) HandleGetAllMarketStats(w http.ResponseWriter, r *http.Request) {

	ex := e.RelayerService.GetRelayerAddress(r)
	res, err := e.OHLCVService.GetAllTokenPairDataByCoinbase(ex)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.PairData{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
	return

}

// HandleGetMarketStats get market specific token data
func (e *MarketsEndpoint) HandleGetMarketStats(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	baseToken := v.Get("baseToken")
	quoteToken := v.Get("quoteToken")

	if quoteToken == "" {
		httputils.WriteError(w, http.StatusBadRequest, "quoteToken Parameter missing")
		return
	}

	if baseToken == "" {
		httputils.WriteError(w, http.StatusBadRequest, "baseToken Parameter missing")
		return
	}

	if !common.IsHexAddress(baseToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Base Token Address")
		return
	}

	if !common.IsHexAddress(quoteToken) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Quote Token Address")
		return
	}

	baseTokenAddress := common.HexToAddress(baseToken)
	quoteTokenAddress := common.HexToAddress(quoteToken)

	res := e.OHLCVService.GetTokenPairData(baseTokenAddress, quoteTokenAddress)
	if res == nil {
		httputils.WriteError(w, http.StatusInternalServerError, "Pair data not found")
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Pair{})
		return
	}
	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *MarketsEndpoint) handleMarketsWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetMarketSocket()

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		e.MarketsService.Subscribe(c)
	}

	if ev.Type == types.UNSUBSCRIBE {
		e.MarketsService.Unsubscribe(c)
	}
}
