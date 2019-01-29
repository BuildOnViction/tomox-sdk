const fs = require('fs')
const argv = require('yargs').argv
const faker = require('faker')
const process = require('process')
const utils = require('ethers').utils
const path = require('path')
const MongoClient = require('mongodb').MongoClient
const { getNetworkID } = require('./utils/helpers')
const { DB_NAME, mongoUrl, network } = require('./utils/config')
const networkID = getNetworkID(network)

const truffleBuildPath = path.join(
  `${process.env.TOMO_DEX_PATH}`,
  `/build/contracts`,
)
const {
  quoteTokens,
  makeFees,
  takeFees,
  baseTokens,
  contractAddresses,
  decimals,
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
    console.log('Seeding quotes tokens')
    db = client.db(DB_NAME)

    documents = quoteTokens.map(symbol => ({
      symbol: symbol,
      contractAddress: utils.getAddress(addresses[symbol]),
      decimals: decimals[symbol],
      makeFee: makeFees[symbol].toString(),
      takeFee: takeFees[symbol].toString(),
      quote: true,
      createdAt: new Date(faker.fake('{{date.recent}}')),
    }))

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
