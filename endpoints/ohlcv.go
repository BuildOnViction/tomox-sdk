package endpoints

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/ws"
)

var (
	startTs = time.Date(1970, time.January, 1, 0, 0, 0, 0, time.UTC)
)

type OHLCVEndpoint struct {
	ohlcvService interfaces.OHLCVService
}

func ServeOHLCVResource(
	r *gin.Engine,
	ohlcvService interfaces.OHLCVService,
) {
	e := &OHLCVEndpoint{ohlcvService}
	// TODO: using GET and contruct param object from query params
	r.GET("/ohlcv", e.handleGetOHLCV)
	ws.RegisterChannel(ws.OHLCVChannel, e.ohlcvWebSocket)
}

func (e *OHLCVEndpoint) handleGetOHLCV(c *gin.Context) {
	var model types.TickRequest

	bt := c.Query("baseToken")
	qt := c.Query("quoteToken")
	pair := c.Query("pairName")
	unit := c.DefaultQuery("unit", "hour")
	duration := c.Query("duration")
	from := c.Query("from")
	to := c.Query("to")

	model.Units = unit

	if duration == "" {
		model.Duration = 24
	} else {
		d, _ := strconv.Atoi(duration)
		model.Duration = int64(d)
	}

	now := time.Now()

	if to == "" {
		model.To = now.Unix()
	} else {
		t, _ := strconv.Atoi(to)
		model.To = int64(t)
	}

	if from == "" {
		model.From = startTs.Unix()
	} else {
		f, _ := strconv.Atoi(from)
		model.From = int64(f)
	}

	if bt == "" {
		c.JSON(http.StatusBadRequest, GinError("baseToken Parameter missing"))
		return
	}

	if qt == "" {
		c.JSON(http.StatusBadRequest, GinError("quoteToken Parameter missing"))
		return
	}

	if !common.IsHexAddress(bt) {
		c.JSON(http.StatusBadRequest, GinError("Invalid base token address"))
		return
	}

	if !common.IsHexAddress(qt) {
		c.JSON(http.StatusBadRequest, GinError("Invalid quote token address"))
		return
	}

	model.Pair = []types.PairSubDoc{{
		BaseToken:  common.HexToAddress(bt),
		QuoteToken: common.HexToAddress(qt),
		Name:       pair,
	}}

	res, err := e.ohlcvService.GetOHLCV(model.Pair, model.Duration, model.Units, model.From, model.To)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (e *OHLCVEndpoint) ohlcvWebSocket(input interface{}, conn *ws.Conn) {

	mab, _ := json.Marshal(input)
	var event *types.WebsocketEvent

	err := json.Unmarshal(mab, &event)
	if err != nil {
		logger.Error(err)
	}

	socket := ws.GetOHLCVSocket()

	if event.Type != "subscription" {
		socket.SendErrorMessage(conn, "Invalid payload")
		return
	}

	dab, _ := json.Marshal(event.Payload)
	var msg *types.WebSocketSubscription

	err = json.Unmarshal(dab, &msg)
	if err != nil {
		logger.Error(err)
	}

	if (msg.Pair.BaseToken == common.Address{}) {
		socket.SendErrorMessage(conn, "Invalid base token")
		return
	}

	if (msg.Pair.QuoteToken == common.Address{}) {
		socket.SendErrorMessage(conn, "Invalid Quote Token")
		return
	}

	if msg.Params.From == 0 {
		msg.Params.From = startTs.Unix()
	}

	if msg.Params.To == 0 {
		msg.Params.To = time.Now().Unix()
	}

	if msg.Params.Duration == 0 {
		msg.Params.Duration = 24
	}

	if msg.Params.Units == "" {
		msg.Params.Units = "hour"
	}

	if msg.Event == types.SUBSCRIBE {
		e.ohlcvService.Subscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken, &msg.Params)
	}

	if msg.Event == types.UNSUBSCRIBE {
		e.ohlcvService.Unsubscribe(conn, msg.Pair.BaseToken, msg.Pair.QuoteToken, &msg.Params)
	}
}
