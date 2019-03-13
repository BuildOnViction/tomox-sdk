const utils = require('ethers').utils
const faker = require('faker')
const MongoClient = require('mongodb').MongoClient
const { getPriceMultiplier } = require('./utils/helpers')
const { DB_NAME, mongoUrl, supportedPairs } = require('./utils/config')

let client, db

const getToken = (symbol, tokens) => {
  return tokens.find(t => {
    return t.symbol === symbol
  })
}

const seed = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true },
    )
    db = client.db(DB_NAME)

    let pairs = []

    const tokens = await db
      .collection('tokens')
      .find(
        {},
        { symbol: 1, contractAddress: 1, decimals: 1, makeFee: 1, takeFee: 1 },
      )
      .toArray()

    supportedPairs.forEach(pair => {
      pair = pair.split('/')
      const baseTokenSymbol = pair[0]
      const quoteTokenSymbol = pair[1]

      const baseToken = getToken(baseTokenSymbol, tokens)
      const quoteToken = getToken(quoteTokenSymbol, tokens)

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

    console.log(pairs)

    await db.collection('pairs').insertMany(pairs)
  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

seed()
