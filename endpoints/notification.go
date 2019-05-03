package endpoints

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"github.com/tomochain/tomodex/interfaces"
	"github.com/tomochain/tomodex/types"
	"github.com/tomochain/tomodex/ws"
)

type NotificationEndpoint struct {
	NotificationService interfaces.NotificationService
}

// ServeNotificationResource sets up the routing of notification endpoints and the corresponding handlers.
func ServeNotificationResource(
	r *mux.Router,
	notificationService interfaces.NotificationService,
) {
	e := &NotificationEndpoint{notificationService}

	ws.RegisterChannel(ws.NotificationChannel, e.handleNotificationWebSocket)
}

func (e *NotificationEndpoint) handleNotificationWebSocket(input interface{}, c *ws.Client) {
	b, _ := json.Marshal(input)
	var ev *types.WebsocketEvent

	err := json.Unmarshal(b, &ev)
	if err != nil {
		logger.Error(err)
	}
}
