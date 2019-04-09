package ws

import (
	logging "github.com/op/go-logging"
	"github.com/tomochain/dex-server/types"
	"github.com/tomochain/dex-server/utils"
)

type WebsocketLogger struct {
	*logging.Logger
	mainLogger             *logging.Logger
	websocketMessageLogger *logging.Logger
}

func NewWebsocketLogger() *WebsocketLogger {
	return &WebsocketLogger{
		utils.StdoutLogger,
		utils.MainLogger,
		utils.WebsocketMessagesLogger,
	}
}

func (l *WebsocketLogger) LogMessageIn(msg *types.WebsocketMessage) {
	l.mainLogger.Infof("Receiving %v/%v message", msg.Channel, msg.Event.Type, utils.JSON(msg))
	l.websocketMessageLogger.Infof("Receiving %v/%v message", msg.Channel, msg.Event.Type, utils.JSON(msg))
}

func (l *WebsocketLogger) LogMessageOut(msg *types.WebsocketMessage) {
	l.mainLogger.Infof("Receiving %v/%v message", msg.Channel, msg.Event.Type, utils.JSON(msg))
	l.websocketMessageLogger.Infof("Sending %v/%v message", msg.Channel, msg.Event.Type, utils.JSON(msg))
}
