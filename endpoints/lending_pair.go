package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

type lendingPairEndpoint struct {
	lendingPairService interfaces.LendingPairService
	relayerService     interfaces.RelayerService
}

// ServeLendingPairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServeLendingPairResource(
	r *mux.Router,
	p interfaces.LendingPairService,
	rl interfaces.RelayerService,
) {
	e := &lendingPairEndpoint{p, rl}
	r.HandleFunc("/api/lending/pairs", e.HandleGetLendingPairs).Methods("GET")
}

func (e *lendingPairEndpoint) HandleGetLendingPairs(w http.ResponseWriter, r *http.Request) {
	ex := e.relayerService.GetRelayerAddress(r)
	res, err := e.lendingPairService.GetAllByCoinbase(ex)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	if res == nil {
		httputils.WriteJSON(w, http.StatusOK, []types.Pair{})
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
