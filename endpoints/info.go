package endpoints

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomoxsdk/app"
	"github.com/tomochain/tomoxsdk/interfaces"
	"github.com/tomochain/tomoxsdk/utils/httputils"
)

type infoEndpoint struct {
	walletService interfaces.WalletService
	tokenService  interfaces.TokenService
}

func ServeInfoResource(
	r *mux.Router,
	walletService interfaces.WalletService,
	tokenService interfaces.TokenService,
) {

	e := &infoEndpoint{walletService, tokenService}
	r.HandleFunc("/info", e.handleGetInfo)
	r.HandleFunc("/info/exchange", e.handleGetExchangeInfo)
	r.HandleFunc("/info/fees", e.handleGetFeeInfo)
}

func (e *infoEndpoint) handleGetInfo(w http.ResponseWriter, r *http.Request) {
	ex := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	quotes, err := e.tokenService.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
	}

	fees := []map[string]string{}
	for _, q := range quotes {
		fees = append(fees, map[string]string{
			"quote":   q.Symbol,
			"makeFee": q.MakeFee.String(),
			"takeFee": q.TakeFee.String(),
		})
	}

	res := map[string]interface{}{
		"exchangeAddress": ex.Hex(),
		"fees":            fees,
	}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetExchangeInfo(w http.ResponseWriter, r *http.Request) {
	ex := common.HexToAddress(app.Config.Ethereum["exchange_address"])

	res := map[string]string{"exchangeAddress": ex.Hex()}

	httputils.WriteJSON(w, http.StatusOK, res)
}

func (e *infoEndpoint) handleGetFeeInfo(w http.ResponseWriter, r *http.Request) {
	quotes, err := e.tokenService.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
	}

	fees := []map[string]string{}
	for _, q := range quotes {
		fees = append(fees, map[string]string{
			"quote":   q.Symbol,
			"makeFee": q.MakeFee.String(),
			"takeFee": q.TakeFee.String(),
		})
	}

	httputils.WriteJSON(w, http.StatusOK, fees)
}
