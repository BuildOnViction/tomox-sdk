const { utils } = require('ethers');
const fs = require('fs');
const path = require('path');

// __dirname is running folder, __filename is current included file
const tokenContent = fs
  .readFileSync(path.resolve(__dirname, '../../seed-data/tokens.bson'))
  .toString();

const initMap = { '8888': {} };
const symbols = [];
const contractAddresses = tokenContent
  .split('\n')
  .filter(line => line.length)
  .reduce((map, line) => {
    try {
      const token = JSON.parse(line);
      // console.log(token);
      initMap['8888'][token.symbol] = token.contractAddress;
      symbols.push(token.symbol);
    } catch (e) {}
    return initMap;
  }, initMap);

const quoteTokens = ['WETH', 'DAI'];
const baseTokens = symbols.filter(symbol => !quoteTokens.includes(symbol));

const makeFees = {
  WETH: utils
    .bigNumberify(10)
    .pow(18)
    .div(250),
  DAI: utils
    .bigNumberify(10)
    .pow(18)
    .div(2),
  USDC: utils
    .bigNumberify(10)
    .pow(6)
    .div(2)
};

const takeFees = {
  WETH: utils
    .bigNumberify(10)
    .pow(18)
    .div(250),
  DAI: utils
    .bigNumberify(10)
    .pow(18)
    .div(2),
  USDC: utils
    .bigNumberify(10)
    .pow(6)
    .div(2)
};

const decimals = symbols.reduce((map, symbol) => {
  map[symbol] = 18;
  return map;
}, {});

module.exports = {
  DB_NAME: 'tomodex',
  addresses: [
    '0x28074f8D0fD78629CD59290Cac185611a8d60109',
    '0x6e6BB166F420DDd682cAEbf55dAfBaFda74f2c9c',
    '0x53ee745b3d30d692dc016450fef68a898c16fa44',
    '0xe0a1240b358dfa6c167edea09c763ae9f3b51ea0'
  ],
  keys: {
    '1': (process.env.TOMO_MAINNET_KEYS || '').split(','),
    '4': (process.env.TOMO_RINKEBY_KEYS || '').split(','),
    '8888': [
      '0x3411b45169aa5a8312e51357db68621031020dcf46011d7431db1bbb6d3922ce',
      '0x75c3e3150c0127af37e7e9df51430d36faa4c4660b6984c1edff254486d834e9'
    ]
  },
  quoteTokens,
  baseTokens,
  takeFees,
  makeFees,
  decimals,
  contractAddresses
};
