package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomochain/tomoxsdk/interfaces"
	"github.com/tomochain/tomoxsdk/utils/httputils"
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
	r.HandleFunc("/relayer/update", e.handleRelayerUpdate).Methods("PUT")
}

func (e *relayerEndpoint) handleRelayerUpdate(w http.ResponseWriter, r *http.Request) {
	err := e.relayerService.UpdateRelayer()

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	httputils.WriteMessage(w, http.StatusOK, "")
}
