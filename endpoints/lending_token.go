package endpoints

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/app"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/utils/httputils"
)

type lendingTokenEndpoint struct {
	collateralTokenService interfaces.TokenService
	lendingTokenService    interfaces.TokenService
	lendingPairservice     interfaces.LendingPairService
}

// ServeLendingTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeLendingTokenResource(
	r *mux.Router,
	collateralTokenService interfaces.TokenService,
	lendingTokenService interfaces.TokenService,
	lendingPairservice interfaces.LendingPairService,
) {
	e := &lendingTokenEndpoint{collateralTokenService, lendingTokenService, lendingPairservice}
	r.HandleFunc("/api/lending/collateraltoken/{address}", e.handleGetCollateralToken).Methods("GET")
	r.HandleFunc("/api/lending/collateraltoken", e.handleGetCollateralTokens).Methods("GET")
	r.HandleFunc("/api/lending/lendingtoken/{address}", e.handleGetLendingToken).Methods("GET")
	r.HandleFunc("/api/lending/lendingtoken", e.handleGetLendingTokens).Methods("GET")
	r.HandleFunc("/api/lending/terms", e.handleGetTerms).Methods("GET")
}
func (e *lendingTokenEndpoint) handleGetTerms(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	relayerAddress := v.Get("relayerAddress")
	if relayerAddress == "" {
		relayerAddress = app.Config.Tomochain["exchange_address"]
	}
	ex := common.HexToAddress(relayerAddress)
	res, err := e.lendingPairservice.GetAllByCoinbase(ex)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
	terms := []uint64{}
	for _, p := range res {
		terms = append(terms, p.Term)
	}
	t := new(struct {
		Terms []uint64 `json:"terms"`
	})
	t.Terms = terms
	httputils.WriteJSON(w, http.StatusOK, t)
}
func (e *lendingTokenEndpoint) handleGetCollateralTokens(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	relayerAddress := v.Get("relayerAddress")
	if relayerAddress == "" {
		relayerAddress = app.Config.Tomochain["exchange_address"]
	}
	ex := common.HexToAddress(relayerAddress)
	res, err := e.collateralTokenService.GetAllByCoinbase(ex)
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
	v := r.URL.Query()
	relayerAddress := v.Get("relayerAddress")
	if relayerAddress == "" {
		relayerAddress = app.Config.Tomochain["exchange_address"]
	}
	ex := common.HexToAddress(relayerAddress)
	res, err := e.lendingTokenService.GetAllByCoinbase(ex)
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
