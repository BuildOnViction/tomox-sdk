package endpoints

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/swap"
	"github.com/tomochain/backend-matching-engine/types"
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
	v := r.URL.Query()
	chainStr := v.Get("chain")
	var chain types.Chain
	err := chain.Scan([]byte(chainStr))
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Chain is not correct")
		return
	}
	address, err := e.depositService.GenerateAddress(chain)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Can not generate Address")
		return
	}

	response := types.GenerateAddressResponse{
		ProtocolVersion: swap.ProtocolVersion,
		Chain:           chain.String(),
		Address:         address.String(),
		Signer:          e.depositService.SignerPublicKey(),
	}

	httputils.WriteJSON(w, http.StatusOK, response)
}

// return Address association for testing first
func (e *depositEndpoint) handleRecoveryTransaction(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("userAddress")
	chainStr := v.Get("chain")
	var chain types.Chain
	err := chain.Scan([]byte(chainStr))
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Chain is not correct")
		return
	}

	if addr == "" {
		httputils.WriteError(w, http.StatusBadRequest, "address Parameter missing")
		return
	}

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid User Address")
		return
	}

	address := common.HexToAddress(addr)

	association, err := e.depositService.GetAssociationByChainAddress(chain, address)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "Can not get address association")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, association)
}
