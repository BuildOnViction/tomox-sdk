const { utils } = require('ethers')
const fs = require('fs')
const path = require('path')
const argv = require('yargs').argv
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'
const network = argv.network || 'development'
const { getNetworkID } = require('./helpers')
const networkId = getNetworkID(network)
// __dirname is running folder, __filename is current included file
const tokenContent = fs.readFileSync(path.resolve(__dirname, './addresses.json')).toString()

const contractAddresses = JSON.parse(tokenContent)

const symbols = Object.keys(contractAddresses[networkId])

const quoteTokens = ['TOMO', 'BTC', 'ETH', 'USDT']
// const quoteTokens = ['TOMO']

const supportedPairs = [
  'ETH/TOMO',
  'ETH/BTC',
  'BTC/USDT',
  'ETH/USDT',
  'TOMO/BTC',
  'TOMO/ETH',
]

const makeFees = {
  TOMO: utils.bigNumberify(10).pow(18).div(250),
  BTC: utils.bigNumberify(10).pow(18).div(250),
  ETH: utils.bigNumberify(10).pow(18).div(250),
  USDT: utils.bigNumberify(10).pow(18).div(250),
}

const takeFees = {
  TOMO: utils.bigNumberify(10).pow(18).div(250),
  BTC: utils.bigNumberify(10).pow(18).div(250),
  ETH: utils.bigNumberify(10).pow(18).div(250),
  USDT: utils.bigNumberify(10).pow(18).div(250),
}

const decimals = {
  'TOMO': 18,
  'BTC': 8,
  'ETH': 18,
  'USDT': 18,
}

const nativeCurrency = {
  symbol: 'TOMO',
  address: '0x0000000000000000000000000000000000000001',
  decimals: 18,
  makerFee: utils.bigNumberify(10).pow(18).div(250),
  takerFee: utils.bigNumberify(10).pow(18).div(250),
}

module.exports = {
  DB_NAME: 'tomodex',
  keys: {
    '1': (process.env.TOMO_MAINNET_KEYS || '').split(','),
    '4': (process.env.TOMO_RINKEBY_KEYS || '').split(','),
    '89': [
      '0x463D27C152040C4E49C5D9606BF3A27E7CE00ACBA25FF4E6A42DD486C27443DA',
    ],
    '8888': [
      '0x7f4c1bacba63f05827f6d8fc0e22cf68c42005775a7f73abff7d819986bae77c',
      '0x2c52197df32aa00940685ae94aeb4b8b6f4c81e2c5f9d289ec76eb614adb9686',
    ],
  },
  symbols,
  quoteTokens,
  supportedPairs,
  takeFees,
  makeFees,
  decimals,
  mongoUrl,
  network,
  contractAddresses,
  nativeCurrency,
}
