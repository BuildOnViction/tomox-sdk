const faker = require('faker')
const utils = require('ethers').utils
const MongoClient = require('mongodb').MongoClient

const { getNetworkID } = require('./utils/helpers')
const { DB_NAME, mongoUrl, network } = require('./utils/config')
const networkID = getNetworkID(network)

const {
  nativeCurrency,
  symbols,
  quoteTokens,
  contractAddresses,
  decimals,
  makeFees,
  takeFees,
} = require('./utils/config')

let documents = []
let addresses = contractAddresses[networkID]
let client, db

const seed = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true },
    )
    console.log('Seeding tokens collection')
    db = client.db(DB_NAME)

    documents = symbols.map(symbol => {

      if (quoteTokens.includes(symbol)) {
        return {
          symbol: symbol,
          contractAddress: utils.getAddress(addresses[symbol]),
          decimals: decimals[symbol],
          makeFee: makeFees[symbol].toString(),
          takeFee: takeFees[symbol].toString(),
          quote: true,
          createdAt: new Date(faker.fake('{{date.recent}}')),
        }
      }

      return {
        symbol: symbol,
        contractAddress: utils.getAddress(addresses[symbol]),
        decimals: decimals[symbol],
        quote: false,
        createdAt: new Date(faker.fake('{{date.recent}}')),
      }
    })

    // Add TOMO symbol
    if (quoteTokens.includes(nativeCurrency.symbol)) {
      documents.push({
        symbol: nativeCurrency.symbol,
        contractAddress: utils.getAddress(nativeCurrency.address),
        decimals: nativeCurrency.decimals,
        makeFee: makeFees[nativeCurrency.symbol].toString(),
        takeFee: takeFees[nativeCurrency.symbol].toString(),
        quote: true,
        createdAt: new Date(faker.fake('{{date.recent}}')),
      })
    } else {
      documents.push({
        symbol: nativeCurrency.symbol,
        contractAddress: utils.getAddress(nativeCurrency.address),
        decimals: nativeCurrency.decimals,
        quote: false,
        createdAt: new Date(faker.fake('{{date.recent}}')),
      })
    }

    if (documents && documents.length > 0) {
      await db.collection('tokens').insertMany(documents)
    }
    client.close()
  } catch (e) {
    throw new Error(e.message)
  } finally {
    client.close()
  }
}

seed()
