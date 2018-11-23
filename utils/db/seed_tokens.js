const fs = require('fs')
const path = require('path')
const process = require('process')
const argv = require('yargs').argv
const utils = require('ethers').utils
const MongoClient = require('mongodb').MongoClient
const { getNetworkID } = require('../../utils/helpers')

const network = argv.network
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'
const networkID = getNetworkID(network)

const truffleBuildPath = path.join(`${process.env.AMP_DEX_PATH}`, `/build/contracts`)
const { quoteTokens, baseTokens, contractAddresses, decimals } = require('../../config')

let documents = []
let addresses = contractAddresses[networkID]
let client, db, response



const seed = async () => {
  try {
    client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
    db = client.db('proofdex')

    documents = baseTokens.map((symbol) => ({
      symbol: symbol,
      contractAddress: utils.getAddress(addresses[symbol]),
      decimals: decimals[symbol],
      quote: false,
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