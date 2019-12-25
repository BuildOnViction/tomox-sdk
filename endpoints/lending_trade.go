package endpoints

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
)

type lendingTradeEndpoint struct {
	lendingTradeService interfaces.LendingTradeService
}

// ServeLendingTradeResource sets up the routing of trade endpoints and the corresponding handlers.
// TODO trim down to one single endpoint with the 3 following params: base, quote, address
func ServeLendingTradeResource(
	r *mux.Router,
	lendingTradeService interfaces.LendingTradeService,
) {
	e := &lendingTradeEndpoint{lendingTradeService}
	ws.RegisterChannel(ws.TradeChannel, e.lendingTradeWebsocket)
}
func (e *lendingTradeEndpoint) lendingTradeWebsocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	if err := json.Unmarshal(b, &ev); err != nil {
		logger.Error(err)
		return
	}
	socket := ws.GetLendingTradeSocket()
	if ev == nil {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}
	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		err := map[string]string{"Message": "Invalid payload"}
		socket.SendErrorMessage(c, err)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload
	err := json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		if p == nil {
			socket.SendErrorMessage(c, errInvalidPayload)
			return
		}
		if p.Term == 0 {
			err := map[string]string{"Message": "Invalid base token"}
			socket.SendErrorMessage(c, err)
			return
		}

		if (p.LendingToken == common.Address{}) {
			err := map[string]string{"Message": "Invalid lending token"}
			socket.SendErrorMessage(c, err)
			return
		}

		e.lendingTradeService.Subscribe(c, p.Term, p.LendingToken)
	}

	if ev.Type == types.UNSUBSCRIBE {
		if p == nil {
			e.lendingTradeService.Unsubscribe(c)
			return
		}

		e.lendingTradeService.UnsubscribeChannel(c, p.Term, p.LendingToken)
	}
}
