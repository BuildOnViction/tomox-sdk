## TOMOX-LENDING TEST GUILD

Using command line to test DEX, Tomox-lending

## Installation
#### 1. Tomoxlending SDK

- Install [tomoxlending](https://github.com/tomochain/tomox-sdk) branch "tomox-lending"

- Edit config.yaml

```bash
coingecko_api_url: https://api.coingecko.com/api/v3
db_name: tomodex
env: dev
error_file: config/errors.yaml
log_level: DEBUG
tomochain:
  exchange_address: 0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e
  exchange_contract_address: 0x0342d186212b04E69eA682b3bed8e232b6b3361a
  lending_contract_address: 0x4d7eA2cE949216D6b120f3AA10164173615A2b6C
  http_url: http://localhost:8501
  ws_url: ws://localhost:9501
jwt_signing_key: QfCAH04Cob7b71QCqy738vw5XGSnFZ9d
jwt_verification_key: QfCAH04Cob7b71QCqy738vw5XGSnFZ9d
mongo_url: localhost:27017
rabbitmq_url: amqp://guest:guest@localhost:5672/
server_port: 8080
simulated: false
supported_currencies: ETH,TOMO,BTC,USDT
tick_duration:
  day:
  - 1
  hour:
  - 1
  - 4
  - 12
  min:
  - 1
  - 5
  - 15
  - 30
  month:
  - 1
  - 3
  - 6
  - 9
  week:
  - 1
  year:
  - 1

```
#### 2. Tomoxjs
Provide command line to test

[https://github.com/tomochain/tomoxjs]

## Usage
#### Create Lending Order
```bash
./tomoxjs lending-create -l 0x45c25041b8e6CBD5c963E7943007187C3673C7c9 -c 0xC2fa1BA90b15E3612E0067A0020192938784D9C5 -s LEND \
-r 0x0D3ab14BBaD3D99F4203bd7a11aCB94882050E7e -q 5 -i 10 -term 30 -t LO -n 0
```
Param

-l: lending token address

-c: collateral token address

-s: side (LEND/BORROW)

-r: relayer address

-q: collateral quantity

-i: interest

-term: duration

-t: type order (LO/MO)

-n: nonce

#### Cancel Lending Order

```bash
./tomoxjs lending-cancel -i 0x094c31e84f43c00a21b483f2f45c901c35bd73a067d9aeedbc9a747c34f15722 -n 1
```
-i: lending order hash

