package endpoints

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

type lendingTokenEndpoint struct {
	collateralTokenService interfaces.TokenService
	lendingTokenService    interfaces.TokenService
}

// ServeLendingTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeLendingTokenResource(
	r *mux.Router,
	collateralTokenService interfaces.TokenService,
	lendingTokenService interfaces.TokenService,
) {
	e := &lendingTokenEndpoint{collateralTokenService, lendingTokenService}
	r.HandleFunc("/api/lending/collateraltoken/{address}", e.handleGetCollateralToken).Methods("GET")
	r.HandleFunc("/api/lending/collateraltoken", e.handleGetCollateralTokens).Methods("GET")
	r.HandleFunc("/api/lending/lendingtoken/{address}", e.handleGetLendingToken).Methods("GET")
	r.HandleFunc("/api/lending/lendingtoken", e.handleGetLendingTokens).Methods("GET")
}

func (e *lendingTokenEndpoint) handleGetCollateralTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.collateralTokenService.GetAll()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *lendingTokenEndpoint) handleGetCollateralToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]

	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	tokenAddress := common.HexToAddress(a)
	res, err := e.collateralTokenService.GetByAddress(tokenAddress)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *lendingTokenEndpoint) handleGetLendingTokens(w http.ResponseWriter, r *http.Request) {
	res, err := e.lendingTokenService.GetAll()
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *lendingTokenEndpoint) handleGetLendingToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]

	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	tokenAddress := common.HexToAddress(a)
	res, err := e.lendingTokenService.GetByAddress(tokenAddress)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}
