package endpoints

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/utils/httputils"
)

type depositEndpoint struct {
	depositService interfaces.DepositService
}

func ServeDepositResource(
	r *mux.Router,
	depositService interfaces.DepositService,
) {

	e := &depositEndpoint{depositService}
	r.HandleFunc("/deposit/schema", e.handleGetSchema).Methods("GET")
	r.HandleFunc("/deposit/generate-address", e.handleGenerateAddress).Methods("GET")
	r.HandleFunc("/deposit/recovery-transaction", e.handleRecoveryTransaction).Methods("GET")
}

func (e *depositEndpoint) handleGetSchema(w http.ResponseWriter, r *http.Request) {
	schemaVersion := e.depositService.GetSchemaVersion()
	schema := map[string]interface{}{
		"version": schemaVersion,
	}
	httputils.WriteJSON(w, http.StatusOK, schema)
}

func (e *depositEndpoint) handleGenerateAddress(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	// addr := vars["address"]
	// if !common.IsHexAddress(addr) {
	// 	httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	// 	return
	// }

	// address := common.HexToAddress(addr)
	// a, err := e.depositService.GetByAddress(address)
	// if err != nil {
	// 	logger.Error(err)
	// 	httputils.WriteError(w, http.StatusInternalServerError, "")
	// 	return
	// }

	// httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *depositEndpoint) handleRecoveryTransaction(w http.ResponseWriter, r *http.Request) {
	// vars := mux.Vars(r)

	// a := vars["address"]
	// if !common.IsHexAddress(a) {
	// 	httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	// }

	// t := vars["token"]
	// if !common.IsHexAddress(a) {
	// 	httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
	// }

	// addr := common.HexToAddress(a)
	// tokenAddr := common.HexToAddress(t)

	// b, err := e.depositService.GetTokenBalance(addr, tokenAddr)
	// if err != nil {
	// 	logger.Error(err)
	// 	httputils.WriteError(w, http.StatusInternalServerError, "")
	// }

	// httputils.WriteJSON(w, http.StatusOK, b)
}
