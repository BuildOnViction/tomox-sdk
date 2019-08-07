package abi

import (
	"strings"

	"github.com/tomochain/tomochain/accounts/abi"
)

const abiJson = `
[
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "name": "max_relayer",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "max_token",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "min_deposit",
        "type": "uint256"
      }
    ],
    "name": "ConfigEvent",
    "type": "event"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "coinbase",
        "type": "address"
      }
    ],
    "name": "depositMore",
    "outputs": [],
    "payable": true,
    "stateMutability": "payable",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "maxRelayer",
        "type": "uint256"
      },
      {
        "name": "maxToken",
        "type": "uint256"
      },
      {
        "name": "minDeposit",
        "type": "uint256"
      }
    ],
    "name": "reconfigure",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "coinbase",
        "type": "address"
      }
    ],
    "name": "refund",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
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
        "name": "fromTokens",
        "type": "address[]"
      },
      {
        "name": "toTokens",
        "type": "address[]"
      }
    ],
    "name": "register",
    "outputs": [],
    "payable": true,
    "stateMutability": "payable",
    "type": "function"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "name": "deposit",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "tradeFee",
        "type": "uint16"
      },
      {
        "indexed": false,
        "name": "fromTokens",
        "type": "address[]"
      },
      {
        "indexed": false,
        "name": "toTokens",
        "type": "address[]"
      }
    ],
    "name": "RegisterEvent",
    "type": "event"
  },
  {
    "constant": false,
    "inputs": [
      {
        "name": "coinbase",
        "type": "address"
      }
    ],
    "name": "resign",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
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
        "name": "new_owner",
        "type": "address"
      },
      {
        "name": "new_coinbase",
        "type": "address"
      }
    ],
    "name": "transfer",
    "outputs": [],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "function"
  },
  {
    "inputs": [
      {
        "name": "maxRelayers",
        "type": "uint256"
      },
      {
        "name": "maxTokenList",
        "type": "uint256"
      },
      {
        "name": "minDeposit",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "nonpayable",
    "type": "constructor"
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
        "name": "fromTokens",
        "type": "address[]"
      },
      {
        "name": "toTokens",
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
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "name": "deposit",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "tradeFee",
        "type": "uint16"
      },
      {
        "indexed": false,
        "name": "fromTokens",
        "type": "address[]"
      },
      {
        "indexed": false,
        "name": "toTokens",
        "type": "address[]"
      }
    ],
    "name": "UpdateEvent",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "name": "deposit",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "tradeFee",
        "type": "uint16"
      },
      {
        "indexed": false,
        "name": "fromTokens",
        "type": "address[]"
      },
      {
        "indexed": false,
        "name": "toTokens",
        "type": "address[]"
      }
    ],
    "name": "TransferEvent",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "name": "deposit_release_time",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "deposit_amount",
        "type": "uint256"
      }
    ],
    "name": "ResignEvent",
    "type": "event"
  },
  {
    "anonymous": false,
    "inputs": [
      {
        "indexed": false,
        "name": "success",
        "type": "bool"
      },
      {
        "indexed": false,
        "name": "remaining_time",
        "type": "uint256"
      },
      {
        "indexed": false,
        "name": "deposit_amount",
        "type": "uint256"
      }
    ],
    "name": "RefundEvent",
    "type": "event"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "CONTRACT_OWNER",
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
    "name": "getRelayerByCoinbase",
    "outputs": [
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
  },
  {
    "constant": true,
    "inputs": [
      {
        "name": "owner",
        "type": "address"
      }
    ],
    "name": "getRelayerByOwner",
    "outputs": [
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
    "constant": true,
    "inputs": [],
    "name": "MaximumRelayers",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "MaximumTokenList",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "MinimumDeposit",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
      }
    ],
    "payable": false,
    "stateMutability": "view",
    "type": "function"
  },
  {
    "constant": true,
    "inputs": [],
    "name": "RelayerCount",
    "outputs": [
      {
        "name": "",
        "type": "uint256"
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
