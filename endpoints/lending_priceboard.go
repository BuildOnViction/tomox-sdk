package endpoints

import (
	"encoding/json"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomox-sdk/interfaces"
	"github.com/tomochain/tomox-sdk/types"
	"github.com/tomochain/tomox-sdk/ws"
)

type LendingPriceBoardEndpoint struct {
	lendingPriceBoardService interfaces.LendingPriceBoardService
}

// ServeLendingPriceBoardResource sets up the routing of token endpoints and the corresponding handlers.
func ServeLendingPriceBoardResource(
	r *mux.Router,
	lendingPriceBoardService interfaces.LendingPriceBoardService,
) {
	e := &LendingPriceBoardEndpoint{lendingPriceBoardService}

	ws.RegisterChannel(ws.LendingPriceBoardChannel, e.handleLendingPriceBoardWebSocket)
}

func (e *LendingPriceBoardEndpoint) handleLendingPriceBoardWebSocket(input interface{}, c *ws.Client) {
	socket := ws.GetLendingPriceBoardSocket()
	errInvalidPayload := map[string]string{"Message": "Invalid payload"}
	if input == nil {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
		return
	}
	if ev == nil {
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}

	if ev.Type != types.SUBSCRIBE && ev.Type != types.UNSUBSCRIBE {
		logger.Info("Event Type", ev.Type)
		socket.SendErrorMessage(c, errInvalidPayload)
		return
	}

	b, _ = json.Marshal(ev.Payload)
	var p *types.SubscriptionPayload

	err = json.Unmarshal(b, &p)
	if err != nil {
		logger.Error(err)
		msg := map[string]string{"Message": "Internal server error"}
		socket.SendErrorMessage(c, msg)
		return
	}

	if ev.Type == types.SUBSCRIBE {
		if p == nil {
			socket.SendErrorMessage(c, errInvalidPayload)
			return
		}
		if (p.LendingToken == common.Address{}) {
			msg := map[string]string{"Message": "Invalid lending token"}
			socket.SendErrorMessage(c, msg)
			return
		}

		if p.Term == 0 {
			msg := map[string]string{"Message": "Invalid term"}
			socket.SendErrorMessage(c, msg)
			return
		}

		e.lendingPriceBoardService.Subscribe(c, p.Term, p.LendingToken)
	}

	if ev.Type == types.UNSUBSCRIBE {
		if p == nil {
			e.lendingPriceBoardService.Unsubscribe(c)
			return
		}

		e.lendingPriceBoardService.UnsubscribeChannel(c, p.Term, p.LendingToken)
	}
}
