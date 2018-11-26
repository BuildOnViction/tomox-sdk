const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'
const { DB_NAME } = require('./utils/config')

const create = async () => {
  const client = await MongoClient.connect(url, { useNewUrlParser: true })

  const db = client.db(DB_NAME)
  let response = await db.createCollection('config')

  const index = await db.collection('config').createIndex( { "key": 1 }, { unique: true } )

  const documents = [
    {key:'schema_version', value: '2'},
    {key:'ethereum_last_block', value: '0'},
    {key:'ethereum_address_index', value: '0'},
  ]

  response = await db.collection('config').insertMany(documents)
  console.log(response)
}

create()

