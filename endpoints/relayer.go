package endpoints

import (
	"net/http"
	"net/url"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

type relayerEndpoint struct {
	relayerService interfaces.RelayerService
}

// ServeRelayerResource sets up the routing of order endpoints and the corresponding handlers.
func ServeRelayerResource(
	r *mux.Router,
	relayerService interfaces.RelayerService,
) {
	e := &relayerEndpoint{relayerService}
	r.HandleFunc("/api/relayer", e.handleRelayerUpdate).Methods("PUT")
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
