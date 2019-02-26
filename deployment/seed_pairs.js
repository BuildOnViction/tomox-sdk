const utils = require('ethers').utils
const faker = require('faker')
const argv = require('yargs').argv
const MongoClient = require('mongodb').MongoClient
const { getNetworkID, getPriceMultiplier } = require('./utils/helpers')
const { DB_NAME, mongoUrl, network, nativeCurrency } = require('./utils/config')
const networkID = getNetworkID(network)

let client, db

const seed = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true },
    )
    db = client.db(DB_NAME)

    let pairs = []

    const baseTokens = await db
      .collection('tokens')
      .find({ quote: false }, { symbol: 1, contractAddress: 1, decimals: 1 })
      .toArray()

    const quoteTokens = await db
      .collection('tokens')
      .find(
        { quote: true },
        { symbol: 1, contractAddress: 1, decimals: 1, makeFee: 1, takeFee: 1 },
      )
      .toArray()

    quoteTokens.forEach(quoteToken => {
      baseTokens.forEach(baseToken => {
        pairs.push({
          baseTokenSymbol: baseToken.symbol,
          baseTokenAddress: utils.getAddress(baseToken.contractAddress),
          baseTokenDecimals: baseToken.decimals,
          quoteTokenSymbol: quoteToken.symbol,
          quoteTokenAddress: utils.getAddress(quoteToken.contractAddress),
          quoteTokenDecimals: quoteToken.decimals,
          priceMultiplier: getPriceMultiplier(
            baseToken.decimals,
            quoteToken.decimals,
          ).toString(),
          active: true,
          makeFee: quoteToken.makeFee,
          takeFee: quoteToken.takeFee,
          createdAt: new Date(faker.fake('{{date.recent}}')),
        })
      })
    })

    baseTokens.forEach(baseToken => {
      pairs.push({
        baseTokenSymbol: baseToken.symbol,
        baseTokenAddress: utils.getAddress(baseToken.contractAddress),
        baseTokenDecimals: baseToken.decimals,
        quoteTokenSymbol: nativeCurrency.symbol,
        quoteTokenAddress: nativeCurrency.address,
        quoteTokenDecimals: nativeCurrency.decimals,
        active: true,
        makeFee: nativeCurrency.makerFee.toString(),
        takeFee: nativeCurrency.takerFee.toString(),
        createdAt: new Date(faker.fake('{{date.recent}}')),
      })
    })

    console.log(pairs)

    await db.collection('pairs').insertMany(pairs)
  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

seed()
