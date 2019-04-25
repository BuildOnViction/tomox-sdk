package endpoints

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/justinas/alice"
	"github.com/tomochain/tomodex/interfaces"
	"github.com/tomochain/tomodex/middlewares"
	"github.com/tomochain/tomodex/types"
	"github.com/tomochain/tomodex/utils/httputils"
)

type AccountEndpoint struct {
	AccountService interfaces.AccountService
}

func ServeAccountResource(
	r *mux.Router,
	accountService interfaces.AccountService,
) {

	e := &AccountEndpoint{AccountService: accountService}

	r.Handle(
		"/account/create", http.HandlerFunc(e.handleCreateAccount),
	).Methods("POST")

	r.Handle(
		"/account/favorite",
		alice.New(middlewares.VerifySignature).Then(http.HandlerFunc(e.handleGetFavoriteTokens)),
	).Methods("GET")

	r.Handle(
		"/account/favorite/add",
		alice.New(middlewares.VerifySignature).Then(http.HandlerFunc(e.handleAddFavoriteToken)),
	).Methods("POST")

	r.Handle(
		"/account/favorite/remove",
		alice.New(middlewares.VerifySignature).Then(http.HandlerFunc(e.handleDeleteFavoriteToken)),
	).Methods("POST")

	r.Handle(
		"/account/{address}", http.HandlerFunc(e.handleGetAccount),
	).Methods("GET")

	r.Handle(
		"/account/{address}/{token}", http.HandlerFunc(e.handleGetAccountTokenBalance),
	).Methods("GET")
}

func (e *AccountEndpoint) handleCreateAccount(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	addr := v.Get("address")

	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	a := common.HexToAddress(addr)
	existingAccount, err := e.AccountService.GetByAddress(a)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	if existingAccount != nil {
		httputils.WriteJSON(w, http.StatusOK, "Account already exists")
		return
	}

	acc := &types.Account{Address: a}
	err = e.AccountService.Create(acc)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusCreated, acc)
}

func (e *AccountEndpoint) handleGetAccount(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr := vars["address"]
	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	address := common.HexToAddress(addr)
	a, err := e.AccountService.GetByAddress(address)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *AccountEndpoint) handleGetAccountTokenBalance(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	a := vars["address"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
	}

	t := vars["token"]
	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
	}

	addr := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	b, err := e.AccountService.GetTokenBalance(addr, tokenAddr)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusNotFound, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, b)
}

func (e *AccountEndpoint) handleGetFavoriteTokens(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)

	addr := vars["address"]
	if !common.IsHexAddress(addr) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	address := common.HexToAddress(addr)
	a, err := e.AccountService.GetFavoriteTokens(address)
	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusInternalServerError, "")
		return
	}

	httputils.WriteJSON(w, http.StatusOK, a)
}

func (e *AccountEndpoint) handleAddFavoriteToken(w http.ResponseWriter, r *http.Request) {
	v := r.URL.Query()
	a := v.Get("address")
	t := v.Get("token")

	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	if !common.IsHexAddress(t) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
		return
	}

	address := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	err := e.AccountService.AddFavoriteToken(address, tokenAddr)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, tokenAddr)
}

func (e *AccountEndpoint) handleDeleteFavoriteToken(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	a := vars["address"]
	t := vars["token"]

	if !common.IsHexAddress(a) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Address")
		return
	}

	if !common.IsHexAddress(t) {
		httputils.WriteError(w, http.StatusBadRequest, "Invalid Token Address")
		return
	}

	address := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	err := e.AccountService.DeleteFavoriteToken(address, tokenAddr)

	if err != nil {
		logger.Error(err)
		httputils.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	httputils.WriteJSON(w, http.StatusOK, tokenAddr)
}
