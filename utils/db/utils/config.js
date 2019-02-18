const { utils } = require('ethers')
const fs = require('fs')
const path = require('path')
const argv = require('yargs').argv
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'
const network = argv.network || 'local'
// __dirname is running folder, __filename is current included file
const tokenContent = fs
    .readFileSync(
        process.env.TOKEN_ADDRESSES ||
        path.resolve(
            __dirname,
            '../../../../dex-client/src/config/addresses.json',
        ),
    )
    .toString()

const contractAddresses = JSON.parse(tokenContent)

const symbols = Object.keys(contractAddresses['8888'])

const quoteTokens = []
const baseTokens = symbols.filter(symbol => !quoteTokens.includes(symbol))

const makeFees = {
    WETH: utils.bigNumberify(10).pow(18).div(250),
    DAI: utils.bigNumberify(10).pow(18).div(2),
}

const takeFees = {
    WETH: utils.bigNumberify(10).pow(18).div(250),
    DAI: utils.bigNumberify(10).pow(18).div(2),
}

const decimals = symbols.reduce((map, symbol) => {
    map[symbol] = 18
    return map
}, {})

const nativeCurrency = {
    symbol: 'TOMO',
    address: '0x0000000000000000000000000000000000000001',
    decimals: 18,
    makerFee: utils.bigNumberify(10).pow(18).div(250),
    takerFee: utils.bigNumberify(10).pow(18).div(250),
}

module.exports = {
    DB_NAME: 'tomodex',
    addresses: [
        '0x28074f8D0fD78629CD59290Cac185611a8d60109',
        '0x6e6BB166F420DDd682cAEbf55dAfBaFda74f2c9c',
        '0x53ee745b3d30d692dc016450fef68a898c16fa44',
        '0xe0a1240b358dfa6c167edea09c763ae9f3b51ea0',
    ],
    keys: {
        '1': (process.env.TOMO_MAINNET_KEYS || '').split(','),
        '4': (process.env.TOMO_RINKEBY_KEYS || '').split(','),
        '8888': [
            '0x7f4c1bacba63f05827f6d8fc0e22cf68c42005775a7f73abff7d819986bae77c',
            '0x2c52197df32aa00940685ae94aeb4b8b6f4c81e2c5f9d289ec76eb614adb9686',
        ],
    },
    quoteTokens,
    baseTokens,
    takeFees,
    makeFees,
    decimals,
    mongoUrl,
    network,
    contractAddresses,
    nativeCurrency,
}
