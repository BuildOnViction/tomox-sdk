package abi

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const abiJson = `[
{
	"constant": true,
	"inputs": [
	  {
		"name": "coinbase",
		"type": "address"
	  }
	],
	"name": "getRelayerByCoinbase",
	"outputs": [
	  {
		"name": "",
		"type": "uint256"
	  },
	  {
		"name": "",
		"type": "address"
	  },
	  {
		"name": "",
		"type": "uint256"
	  },
	  {
		"name": "",
		"type": "uint16"
	  },
	  {
		"name": "",
		"type": "address[]"
	  },
	  {
		"name": "",
		"type": "address[]"
	  }
	],
	"payable": false,
	"stateMutability": "view",
	"type": "function"
  }
]`

// GetRelayerAbi return ABI relayer
func GetRelayerAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(abiJson))
}
