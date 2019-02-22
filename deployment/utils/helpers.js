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

const getPriceMultiplier = (baseTokenDecimals, quoteTokenDecimals) => {
    let defaultPricepointMultiplier = utils.bigNumberify(1e9)
    let decimalsPricepointMultiplier = utils.bigNumberify(
        (10 ** (baseTokenDecimals - quoteTokenDecimals)).toString(),
    )

    return defaultPricepointMultiplier.mul(decimalsPricepointMultiplier)
}

module.exports = {
    getNetworkID,
    getPriceMultiplier,
}
