# DEX backend

Official decentralized exchange backend, forked from the Proof project.  
The matching-engine will be soon moved to blockchain services.

# Getting Started

## Requirements

- **mongoDB** version 3.6 or newer
- **rabbitmq** version 3.7.7 or newer
- **redis** version 4.0 or newer
- **dep** latest

## Booting up the server

**Install the dependencies**

You need to run `./install-requirements.sh` to install all required libraries

```bash
cd dex-server
export BACKEND=$GOPATH/src/github.com/tomochain/backend-matching-engine
mkdir -p $BACKEND
ln -sF $PWD $BACKEND
cd $BACKEND
dep ensure
```

**Start the development**

_If you need to generate genesis block, take a look at [seed-data](./utils/seed-data/README.md)_

```bash
# start dockers including mongo, redis, rabbitmq
yarn start-env
# generate seeds data
yarn generate-seeds
# seed the sample data
yarn seeds
# this will start the server in hot-reload mode
yarn hot
```

# API Endpoints

## Tokens

- `GET /tokens` : returns list of all the tokens from the database
- `GET /tokens/<addr>`: returns details of a token from db using token's contract address
- `POST /tokens`: Create/Insert token in DB. Sample input:

```
{
	"name":"HotPotCoin",
	"symbol":"HPC",
	"decimal":18,
	"contractAddress":"0x1888a8db0b7db59413ce07150b3373972bf818d3",
	"active":true,
	"quote":true
}
```

## Pairs

- `GET /pairs` : returns list of all the pairs from the database
- `GET /pairs/<baseToken>/<quoteToken>`: returns details of a pair from db using using contract address of its constituting tokens
- `GET /pairs/book/<pairName>`: Returns orderbook for the pair using pair name
- `POST /pairs`: Create/Insert pair in DB. Sample input:

```
{
    "baseToken":"5b3e82587b44576ba8000001",
    "quoteToken":"5b3e82607b44576ba8000002",
    "active":true,
    "quoteTokenSymbol":"hpc"
}
```

## Address

- `POST /address`: Create/Insert address and corresponding balance entry in DB. Sample input:

```
{
	"address":"0xefD7eB287CeeFCE8256Dd46e25F398acEA7C4b63"
}
```

## Balance

- `GET /account/<addr>`: Fetch the balance details from db of the given address.

## Order

- `GET /orders/<addr>`: Fetch all the orders placed by the given address

## Trade

- `GET /trades/history/<pair>`: Fetch complete trade history of given pair using pair name
- `GET /trades/<addr>`: Fetch all the trades in which the given address is either maker or taker
- `GET /trades/ticks`: Fetch ohlcv data. Query Params:

```
// Query Params for /trades/ticks
pairName: names of pair separated by comma(,) ex: "hpc/aut,abc/xyz". (Atleast 1 Required)
unit: sec,min,hour,day,week,month,yr. (default:hour)
duration: in int. (default: 24)
from: unix timestamp of from time.(default: start of timestamp)
to: unix timestamp of to time. (default: current timestamp)
```

# Types

## Orders

Orders contain the information that is required to register an order in the orderbook as a "Maker".

- **id** is the primary ID of the order (possibly deprecated)
- **orderType** is either BUY or SELL. It is currently not parsed by the server and compute directly from buyToken, sellToken, buyAmount, sellAmount
- **exchangeAddress** is the exchange smart contract address
- **maker** is the maker (usually sender) ethereum account address
- **buyToken** is the BUY token ethereum address
- **sellToken** is the SELL token ethereum address
- **buyAmount** is the BUY amount (in BUY_TOKEN units)
- **sellAmount** is the SELL amount (in SELL_TOKEN units)
- **expires** is the order expiration timestamp
- **nonce** is the nonce that corresponds to
- **feeMake** is the maker fee (not implemented yet)
- **feeTake** is the taker fee (not implemented yet)
- **pairID** is a hash of the corresponding
- **hash** is a hash of the order details (see details below)
- **signature** is a signature of the order hash. The signer must equal to the maker address for the order to be valid.
- **price** corresponds to the pricepoint computed by the matching engine (not parsed)
- **amount** corresponds to the amount computed by the matching engine (not parsed)

**Order Price and Amount**

There are two ways to describe the amount of tokens being bought/sold. The smart-contract requires (buyToken, sellToken, buyAmount, sellAmount) while the
orderbook requires (pairID, amount, price).

The conversion between both systems can be found in the engine.ComputeOrderPrice
function

**Order Hash**

The order hash is a sha-256 hash of the following elements:

- Exchange address
- Token Buy address
- Amount Buy
- Token Sell Address
- Amount Sell
- Expires
- Nonce
- Maker Address

## Trades

When an order matches another order in the orderbook, the "taker" is required
to sign a trade object that matches an order.

- **orderHash** is the hash of the matching order
- **amount** is the amount of tokens that will be traded
- **trade nonce** is a unique integer to distinguish successive but identical orders (note: can probably be renamed to nonce)
- **taker** is the taker ethereum account address
- **pairID** is a hash identifying the token pair that will be traded
- **hash** is a unique identifier hash of the trade details (see details below)
- **signature** is a signature of the trade hash

Trade Hash:

The trade hash is a sha-256 hash of the following elements:

- Order Hash
- Amount
- Taker Address
- Trade Nonce

The (Order, Trade) tuple can then be used to perform an on-chain transaction for this trade.

## Quote Tokens and Token Pairs

In the same way as traditional exchanges function with the idea of base
currencies and quote currencies, the decentralized exchange works with
base tokens and quote tokens under the following principles:

- Only the exchange operator can register a quote token
- Anybody can register a token pair (but the quote token needs to be registered)

Token pairs are identified by an ID (a hash of both token addresses)

# REST API

See [REST_API.md](REST_API.md)

# Websocket API

See [WEBSOCKET_API.md](WEBSOCKET_API.md)
