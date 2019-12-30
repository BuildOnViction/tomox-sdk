package abi

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const abiJSON = `[
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

const lendingJSON = `[
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"name": "COLLATERALS",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "term",
				"type": "uint256"
			}
		],
		"name": "addTerm",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "LENDINGRELAYER_LIST",
		"outputs": [
			{
				"name": "_tradeFee",
				"type": "uint16"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "coinbase",
				"type": "address"
			},
			{
				"name": "tradeFee",
				"type": "uint16"
			},
			{
				"name": "baseTokens",
				"type": "address[]"
			},
			{
				"name": "terms",
				"type": "uint256[]"
			},
			{
				"name": "collaterals",
				"type": "address[]"
			}
		],
		"name": "update",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "uint256"
			}
		],
		"name": "BASES",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "token",
				"type": "address"
			},
			{
				"name": "depositRate",
				"type": "uint256"
			}
		],
		"name": "addCollateral",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"name": "COLLATERAL_LIST",
		"outputs": [
			{
				"name": "_depositRate",
				"type": "uint256"
			},
			{
				"name": "_price",
				"type": "uint256"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": false,
		"inputs": [
			{
				"name": "token",
				"type": "address"
			}
		],
		"name": "addBaseToken",
		"outputs": [],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [],
		"name": "relayer",
		"outputs": [
			{
				"name": "",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"constant": true,
		"inputs": [
			{
				"name": "coinbase",
				"type": "address"
			}
		],
		"name": "getLendingRelayerByCoinbase",
		"outputs": [
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
				"type": "uint256[]"
			},
			{
				"name": "",
				"type": "address[]"
			}
		],
		"payable": false,
		"stateMutability": "view",
		"type": "function"
	},
	{
		"inputs": [
			{
				"name": "r",
				"type": "address"
			}
		],
		"payable": false,
		"stateMutability": "nonpayable",
		"type": "constructor"
	}
]
`

// GetRelayerAbi return ABI relayer
func GetRelayerAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(abiJSON))
}

// GetLendingAbi return ABI relayer
func GetLendingAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(lendingJSON))
}
