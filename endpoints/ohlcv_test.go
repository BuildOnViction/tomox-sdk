package endpoints

import (
	"github.com/Proofsuite/amp-matching-engine/utils/testutils/mocks"
	"github.com/gorilla/mux"
)

func SetupOhlcvEndpointTest() (*mux.Router, *mocks.OHLCVService) {
	r := mux.NewRouter()
	ohlcvService := new(mocks.OHLCVService)

	ServeOHLCVResource(r, ohlcvService)

	return r, ohlcvService
}
