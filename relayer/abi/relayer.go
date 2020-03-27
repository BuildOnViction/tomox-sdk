package abi

import (
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
)

const abiJSON = `[{"constant":false,"inputs":[{"name":"coinbase","type":"address"},{"name":"fromToken","type":"address"},{"name":"toToken","type":"address"}],"name":"listToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"MaximumRelayers","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"RELAYER_LIST","outputs":[{"name":"_deposit","type":"uint256"},{"name":"_tradeFee","type":"uint16"},{"name":"_index","type":"uint256"},{"name":"_owner","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"}],"name":"depositMore","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[{"name":"","type":"uint256"}],"name":"RELAYER_COINBASES","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"RESIGN_REQUESTS","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"coinbase","type":"address"}],"name":"getRelayerByCoinbase","outputs":[{"name":"","type":"uint256"},{"name":"","type":"address"},{"name":"","type":"uint256"},{"name":"","type":"uint16"},{"name":"","type":"address[]"},{"name":"","type":"address[]"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"},{"name":"tradeFee","type":"uint16"},{"name":"fromTokens","type":"address[]"},{"name":"toTokens","type":"address[]"}],"name":"update","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"maxRelayer","type":"uint256"},{"name":"maxToken","type":"uint256"},{"name":"minDeposit","type":"uint256"}],"name":"reconfigure","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"}],"name":"cancelSelling","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"},{"name":"fromToken","type":"address"},{"name":"toToken","type":"address"}],"name":"deListToken","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"},{"name":"price","type":"uint256"}],"name":"sellRelayer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"RelayerCount","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":true,"inputs":[{"name":"","type":"address"}],"name":"RELAYER_ON_SALE_LIST","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"}],"name":"resign","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"},{"name":"new_owner","type":"address"}],"name":"transfer","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"MinimumDeposit","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"},{"name":"tradeFee","type":"uint16"},{"name":"fromTokens","type":"address[]"},{"name":"toTokens","type":"address[]"}],"name":"register","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":true,"inputs":[],"name":"MaximumTokenList","outputs":[{"name":"","type":"uint256"}],"payable":false,"stateMutability":"view","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"}],"name":"buyRelayer","outputs":[],"payable":true,"stateMutability":"payable","type":"function"},{"constant":false,"inputs":[{"name":"coinbase","type":"address"}],"name":"refund","outputs":[],"payable":false,"stateMutability":"nonpayable","type":"function"},{"constant":true,"inputs":[],"name":"CONTRACT_OWNER","outputs":[{"name":"","type":"address"}],"payable":false,"stateMutability":"view","type":"function"},{"inputs":[{"name":"tomoxListing","type":"address"},{"name":"maxRelayers","type":"uint256"},{"name":"maxTokenList","type":"uint256"},{"name":"minDeposit","type":"uint256"}],"payable":false,"stateMutability":"nonpayable","type":"constructor"},{"anonymous":false,"inputs":[{"indexed":false,"name":"max_relayer","type":"uint256"},{"indexed":false,"name":"max_token","type":"uint256"},{"indexed":false,"name":"min_deposit","type":"uint256"}],"name":"ConfigEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"deposit","type":"uint256"},{"indexed":false,"name":"tradeFee","type":"uint16"},{"indexed":false,"name":"fromTokens","type":"address[]"},{"indexed":false,"name":"toTokens","type":"address[]"}],"name":"RegisterEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"deposit","type":"uint256"},{"indexed":false,"name":"tradeFee","type":"uint16"},{"indexed":false,"name":"fromTokens","type":"address[]"},{"indexed":false,"name":"toTokens","type":"address[]"}],"name":"UpdateEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"owner","type":"address"},{"indexed":false,"name":"deposit","type":"uint256"},{"indexed":false,"name":"tradeFee","type":"uint16"},{"indexed":false,"name":"fromTokens","type":"address[]"},{"indexed":false,"name":"toTokens","type":"address[]"}],"name":"TransferEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"deposit_release_time","type":"uint256"},{"indexed":false,"name":"deposit_amount","type":"uint256"}],"name":"ResignEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"success","type":"bool"},{"indexed":false,"name":"remaining_time","type":"uint256"},{"indexed":false,"name":"deposit_amount","type":"uint256"}],"name":"RefundEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"is_on_sale","type":"bool"},{"indexed":false,"name":"coinbase","type":"address"},{"indexed":false,"name":"price","type":"uint256"}],"name":"SellEvent","type":"event"},{"anonymous":false,"inputs":[{"indexed":false,"name":"success","type":"bool"},{"indexed":false,"name":"coinbase","type":"address"},{"indexed":false,"name":"price","type":"uint256"}],"name":"BuyEvent","type":"event"}]`

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
		  "name": "_liquidationRate",
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
		  "name": "",
		  "type": "uint256"
		}
	  ],
	  "name": "ALL_COLLATERALS",
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
		},
		{
		  "name": "liquidationRate",
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
  ]`

// GetRelayerAbi return ABI relayer
func GetRelayerAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(abiJSON))
}

// GetLendingAbi return ABI relayer
func GetLendingAbi() (abi.ABI, error) {
	return abi.JSON(strings.NewReader(lendingJSON))
}
