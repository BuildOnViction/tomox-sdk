package endpoints

import (
	"math/big"
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

type relayerEndpoint struct {
	relayerService      interfaces.RelayerService
	ohlcvService        interfaces.OHLCVService
	lendingOhlcvService interfaces.LendingOhlcvService
}

// ServeRelayerResource sets up the routing of order endpoints and the corresponding handlers.
func ServeRelayerResource(
	r *mux.Router,
	relayerService interfaces.RelayerService,
	ohlcvService interfaces.OHLCVService,
	lendingOhlcvService interfaces.LendingOhlcvService,
) {
	e := &relayerEndpoint{relayerService, ohlcvService, lendingOhlcvService}
	r.HandleFunc("/api/relayer", e.handleRelayerUpdate).Methods("PUT")
	r.HandleFunc("/api/relayer/volume", e.handleGetVolume).Methods("GET")
	r.HandleFunc("/api/relayer/lending", e.handleGetLendingVolume).Methods("GET")
}

func (e *relayerEndpoint) handleRelayerUpdate(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	authKey := v.Get("authKey")
	relayerName := v.Get("relayerName")
	relayerUrl := v.Get("relayerUrl")
	address := v.Get("relayerAddress")
	relayerAddress := common.HexToAddress(address)

	if app.Config.ApiAuthKey != authKey {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid auth key")
		return
	}

	relayer, err := e.relayerService.GetByAddress(relayerAddress)
	if relayer == nil {
		err = e.relayerService.UpdateRelayer(relayerAddress)
	}

	if relayerUrl != "" {
		u, err := url.Parse(relayerUrl)
		if err == nil {
			relayerUrl = u.Host
		}
	}

	if relayerName != "" {
		err = e.relayerService.UpdateNameByAddress(relayerAddress, relayerName, relayerUrl)
	}

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.WriteMessage(w, http.StatusOK, "OK")
}

// HandleGetVolume get volume relayer
func (e *relayerEndpoint) handleGetVolume(w http.ResponseWriter, r *http.Request) {
	type res struct {
		RelayerAddress common.Address `json:"relayerAddress"`
		TotalVolume    *big.Int       `json:"totalVolume"`
	}
	var result res
	ex := e.relayerService.GetRelayerAddress(r)
	result.RelayerAddress = ex
	volume, err := e.ohlcvService.GetVolumeByCoinbase(ex)
	result.TotalVolume = volume
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, result)
	return

}

func (e *relayerEndpoint) handleGetLendingVolume(w http.ResponseWriter, r *http.Request) {
	type res struct {
		RelayerAddress common.Address `json:"relayerAddress"`
		TotalVolume    *big.Int       `json:"totalLendingVolume"`
		VolumeType     string         `json:"volumeType"`
	}
	var result res
	result.VolumeType = "USDT"
	ex := e.relayerService.GetRelayerAddress(r)
	result.RelayerAddress = ex
	volume, err := e.lendingOhlcvService.GetLendingVolumeByCoinbase(ex)
	result.TotalVolume = volume
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteJSON(w, http.StatusOK, result)
	return

}
