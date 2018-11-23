const utils = require('ethers').utils
const argv = require('yargs').argv
const MongoClient = require('mongodb').MongoClient
const { getNetworkID, getPriceMultiplier } = require('../../utils/helpers')

const network = argv.network
const mongoUrl = argv.mongo_url
const networkID = getNetworkID(network)

let client, db

const seed = async () => {
  try {
    client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
    db = client.db('proofdex')

    let pairs = []

    const tokens = await db.collection('tokens')
      .find(
        { quote: false },
        { symbol: 1, contractAddress: 1, decimals: 1 }
      )
      .toArray()

    const quotes = await db.collection('tokens')
      .find(
        { quote: true },
        { symbol: 1, contractAddress: 1, decimals: 1, makeFee: 1, takeFee: 1 }
      )
      .toArray()

    
    quotes.forEach(quote => {
      tokens.forEach(token => {
        pairs.push({
          baseTokenSymbol: token.symbol,
          baseTokenAddress: utils.getAddress(token.contractAddress),
          baseTokenDecimals: token.decimals,
          quoteTokenSymbol: quote.symbol,
          quoteTokenAddress: utils.getAddress(quote.contractAddress),
          quoteTokenDecimals: quote.decimals,
          priceMultiplier: getPriceMultiplier(token.decimals, quote.decimals),
          active: true,
          makeFee: quote.makeFee,
          takeFee: quote.takeFee,
          createdAt: Date(),
          updatedAt: Date()
        })
      })
    })

    const response = await db.collection('pairs').insertMany(pairs)
  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

seed()