package ethereum

import (
	"github.com/ethereum/go-ethereum/common"
	"github.com/tomochain/tomox-sdk/utils"
)

var logger = utils.Logger

type EthereumConfig struct {
	url             string
	exchangeAddress common.Address
}

func NewEthereumConfig(url string, exchange common.Address) *EthereumConfig {
	return &EthereumConfig{
		url:             url,
		exchangeAddress: exchange,
	}
}

func (c *EthereumConfig) GetURL() string {
	return c.url
}

func (c *EthereumConfig) ExchangeAddress() common.Address {
	return c.exchangeAddress
}
