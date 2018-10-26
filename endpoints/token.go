package endpoints

import (
	"log"
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/services"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/tomochain/backend-matching-engine/ws"
)

type tokenEndpoint struct {
	tokenService interfaces.TokenService
}

// ServeTokenResource sets up the routing of token endpoints and the corresponding handlers.
func ServeTokenResource(
	r *gin.Engine,
	tokenService interfaces.TokenService,
) {
	e := &tokenEndpoint{tokenService}

	r.GET("/tokens/:address", e.handleGetToken)
	r.GET("/tokens", e.handleGetTokens)
	r.POST("/tokens", e.handleCreateTokens)

	ws.RegisterChannel(ws.TokenChannel, e.ws)
}

func (e *tokenEndpoint) handleCreateTokens(c *gin.Context) {
	t := &types.Token{}

	err := c.BindJSON(t)

	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, GinError("Invalid payload"))
	}

	err = e.tokenService.Create(t)
	if err != nil {
		if err == services.ErrTokenExists {
			c.JSON(http.StatusBadRequest, GinError(""))
			return
		}

		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return

	}

	c.JSON(http.StatusCreated, t)
}

func (e *tokenEndpoint) handleGetTokens(c *gin.Context) {
	res, err := e.tokenService.GetAll()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (e *tokenEndpoint) handleGetQuoteTokens(c *gin.Context) {
	res, err := e.tokenService.GetQuoteTokens()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
	}

	c.JSON(http.StatusOK, res)
}

func (e *tokenEndpoint) handleGetBaseTokens(c *gin.Context) {
	res, err := e.tokenService.GetBaseTokens()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, res)
}

func (e *tokenEndpoint) handleGetToken(c *gin.Context) {

	a := c.Param("address")

	if a == "base" {
		e.handleGetBaseTokens(c)
		return
	} else if a == "quote" {
		e.handleGetQuoteTokens(c)
		return
	}

	if !common.IsHexAddress(a) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
	}

	tokenAddress := common.HexToAddress(a)
	res, err := e.tokenService.GetByAddress(tokenAddress)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, res)
}

// ws function handles incoming websocket messages on the order channel
func (e *tokenEndpoint) ws(input interface{}, conn *ws.Conn) {
	// it means that we can handle not only WebSocketPayload but other Payloads as well
	msg, ok := input.(*types.WebSocketPayload)
	if ok {
		switch msg.Type {
		case "GET_TOKENS":
			e.handleGetTokensWS(msg, conn)
			log.Printf("Data: %+v", msg)
		default:
			log.Print("Response with error")
		}
	}
}

// handleSubmitSignatures handles NewTrade messages. New trade messages are transmitted to the corresponding order channel
// and received in the handleClientResponse.
func (e *tokenEndpoint) handleGetTokensWS(p *types.WebSocketPayload, conn *ws.Conn) {
	res, err := e.tokenService.GetAll()
	if err != nil {
		logger.Error(err)
		ws.SendMessage(conn, ws.TokenChannel, ws.ERROR, GinError(""))
		return
	}

	ws.SendMessage(conn, ws.TokenChannel, ws.UPDATE, res)
}
