package endpoints

import (
	"net/http"

	"github.com/ethereum/go-ethereum/common"
	"github.com/gin-gonic/gin"

	"github.com/tomochain/backend-matching-engine/interfaces"
	"github.com/tomochain/backend-matching-engine/services"
	"github.com/tomochain/backend-matching-engine/types"
)

type pairEndpoint struct {
	pairService interfaces.PairService
}

// ServePairResource sets up the routing of pair endpoints and the corresponding handlers.
func ServePairResource(
	r *gin.Engine,
	p interfaces.PairService,
) {
	e := &pairEndpoint{p}

	// we also suport pair channel for websocket

	r.POST("/pairs", e.handleCreatePair)
	r.GET("/pairs", e.handleGetAllPairs)
	r.GET("/pairs/:baseToken/:quoteToken", e.HandleGetPair)

}

func (e *pairEndpoint) handleCreatePair(c *gin.Context) {
	p := &types.Pair{}
	err := c.BindJSON(p)
	if err != nil {
		c.JSON(http.StatusBadRequest, GinError("Invalid payload"))
		return
	}

	err = p.Validate()
	if err != nil {
		c.JSON(http.StatusBadRequest, GinError(err.Error()))
		return
	}

	err = e.pairService.Create(p)
	if err != nil {
		switch err {
		case services.ErrPairExists:
			c.JSON(http.StatusBadRequest, GinError("Pair exists"))
			return
		case services.ErrBaseTokenNotFound:
			c.JSON(http.StatusBadRequest, GinError("Base token not found"))
			return
		case services.ErrQuoteTokenNotFound:
			c.JSON(http.StatusBadRequest, GinError("Quote token not found"))
			return
		case services.ErrQuoteTokenInvalid:
			c.JSON(http.StatusBadRequest, GinError("Quote token invalid (token is not registered as quote"))
			return
		default:
			logger.Error(err)
			c.JSON(http.StatusInternalServerError, GinError(""))
			return
		}
	}

	c.JSON(http.StatusCreated, p)
}

func (e *pairEndpoint) handleGetAllPairs(c *gin.Context) {
	res, err := e.pairService.GetAll()
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, res)
}

func (e *pairEndpoint) HandleGetPair(c *gin.Context) {

	baseToken := c.Param("baseToken")
	quoteToken := c.Param("quoteToken")

	if !common.IsHexAddress(baseToken) {
		c.JSON(http.StatusBadRequest, "Invalid Address")
	}

	if !common.IsHexAddress(quoteToken) {
		c.JSON(http.StatusBadRequest, "Invalid Address")
	}

	baseTokenAddress := common.HexToAddress(baseToken)
	quoteTokenAddress := common.HexToAddress(quoteToken)
	res, err := e.pairService.GetByTokenAddress(baseTokenAddress, quoteTokenAddress)
	if err != nil {
		logger.Error(err)
		c.JSON(http.StatusInternalServerError, "")
		return
	}

	c.JSON(http.StatusOK, res)
}
