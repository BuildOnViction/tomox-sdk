const MongoClient = require('mongodb').MongoClient
const { utils, Wallet } = require('ethers')
const { getNetworkID } = require('./utils/helpers')
const { DB_NAME, keys, mongoUrl, network } = require('./utils/config')
const networkID = getNetworkID(network)
const walletKeys = keys[networkID]

let client, db, documents

const seed = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true },
    )
    db = client.db(DB_NAME)
    documents = []

    walletKeys.forEach(key => {
      let walletRecord = {}
      wallet = new Wallet(key)

      walletRecord.privateKey = wallet.privateKey.slice(2)
      walletRecord.address = utils.getAddress(wallet.address)
      walletRecord.admin = true
      walletRecord.operator = true
      documents.push(walletRecord)
    })

    await db.collection('wallets').insertMany(documents)
  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

seed()
