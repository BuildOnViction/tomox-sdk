const faker = require('faker')
const MongoClient = require('mongodb').MongoClient
const Long = require('mongodb').Long
const { DB_NAME, mongoUrl } = require('./utils/config')
const argv = require('yargs').argv
const network = argv.network || 'development'
const { getEthereumBlockNumber } = require('./utils/helpers')

const create = async () => {
  const ethereumBlockNumber = await getEthereumBlockNumber(network)
  const client = await MongoClient.connect(
    mongoUrl,
    { useNewUrlParser: true },
  )

  const db = client.db(DB_NAME)
  await db.createCollection('config')

  const index = await db
    .collection('config')
    .createIndex({ key: 1 }, { unique: true })

  const createdAt = new Date(faker.fake('{{date.recent}}'))
  const documents = [
    { key: 'schema_version', value: Long.fromInt(2), createdAt },
    { key: 'ethereum_last_block', value: Long.fromInt(ethereumBlockNumber), createdAt },
    { key: 'ethereum_address_index', value: Long.fromInt(0), createdAt },
  ]
  console.log(documents)
  try {
    await db.collection('config').insertMany(documents)
  } catch (e) {
    console.log(e)
  } finally {
    client.close()
  }
}

create()
