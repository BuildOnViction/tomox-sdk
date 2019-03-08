const ethers = require('ethers')
const { utils } = require('ethers')

const getNetworkID = networkName => {
  return {
    ethereum: '1',
    rinkeby: '4',
    tomochain: '88',
    tomochainTestnet: '89',
    development: '8888',
  }[networkName]
}

const getEthereumBlockNumber = async networkName => {
  let httpProvider

  // Get latest block number
  switch (networkName) {
    case 'development':
      httpProvider = ethers.getDefaultProvider('ropsten')
      break
    case 'tomochainTestnet':
      httpProvider = ethers.getDefaultProvider('ropsten')
      break
    case 'tomochain':
      httpProvider = ethers.getDefaultProvider('homestead')
      break
    default:
      httpProvider = ethers.getDefaultProvider('ropsten')
      break
  }

  return await httpProvider.getBlockNumber()
}

const getPriceMultiplier = (baseTokenDecimals, quoteTokenDecimals) => {
  let defaultPricepointMultiplier = utils.bigNumberify(1e9)
  let decimalsPricepointMultiplier = utils.bigNumberify(
    (10 ** (Math.abs(baseTokenDecimals - quoteTokenDecimals))).toString(),
  )

  return defaultPricepointMultiplier.mul(decimalsPricepointMultiplier)
}

module.exports = {
  getNetworkID,
  getEthereumBlockNumber,
  getPriceMultiplier,
}
