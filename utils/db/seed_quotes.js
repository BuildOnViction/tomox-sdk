const fs = require('fs')
const argv = require('yargs').argv
const process = require('process')
const utils = require('ethers').utils
const path = require('path')
const MongoClient = require('mongodb').MongoClient
const { getNetworkID } = require('../../utils/helpers')

const network = argv.network
const mongoUrl = argv.mongo_url
const networkID = getNetworkID(network)

const truffleBuildPath = path.join(`${process.env.AMP_DEX_PATH}`, `/build/contracts`)
const { quoteTokens, makeFees, takeFees, baseTokens, contractAddresses, decimals } = require('../../config')

let documents = []
let addresses = contractAddresses[networkID]

let client, db, response

const seed = async () => {
  try {
    client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
    db = client.db('proofdex')

    documents = quoteTokens.map((symbol) => ({
      symbol: symbol,
      contractAddress: utils.getAddress(addresses[symbol]),
      decimals: decimals[symbol],
      makeFee: makeFees[symbol],
      takeFee: takeFees[symbol],
      quote: true,
      createdAt: Date(),
      updatedAt: Date()
    }))

    response = await db.collection('tokens').insertMany(documents)
    client.close()
  } catch (e) {
    throw new Error(e.message)
  } finally {
    client.close()
  }
}

seed()