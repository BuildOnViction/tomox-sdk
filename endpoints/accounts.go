package endpoints

import (
	"net/http"

	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/types"
	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"
)

type accountEndpoint struct {
	accountService interfaces.AccountService
}

func ServeAccountResource(
	r *gin.Engine,
	accountService interfaces.AccountService,
) {

	e := &accountEndpoint{accountService}
	r.POST("/account", e.handleCreateAccount)
	r.GET("/account/:address", e.handleGetAccount)
	r.GET("/account/:address/:token", e.handleGetAccountTokenBalance)
}

func (e *accountEndpoint) handleCreateAccount(c *gin.Context) {
	a := &types.Account{}
	err := c.BindJSON(a)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, GinError("Invalid payload"))
		return
	}

	err = a.Validate()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusBadRequest, GinError("Invalid payload"))
		return
	}

	err = e.accountService.Create(a)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusCreated, a)
}

func (e *accountEndpoint) handleGetAccount(c *gin.Context) {

	addr := c.Param("address")
	if !common.IsHexAddress(addr) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
		return
	}

	address := common.HexToAddress(addr)
	a, err := e.accountService.GetByAddress(address)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, a)
}

func (e *accountEndpoint) handleGetAccountTokenBalance(c *gin.Context) {

	a := c.Param("address")
	if !common.IsHexAddress(a) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Address"))
		return
	}

	t := c.Param("token")
	if !common.IsHexAddress(a) {
		c.JSON(http.StatusBadRequest, GinError("Invalid Token Address"))
		return
	}

	addr := common.HexToAddress(a)
	tokenAddr := common.HexToAddress(t)

	b, err := e.accountService.GetTokenBalance(addr, tokenAddr)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, GinError(""))
		return
	}

	c.JSON(http.StatusOK, b)
}
