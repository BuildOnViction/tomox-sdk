const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'
const { DB_NAME } = require('./utils/config')

const create = async () => {
  const client = await MongoClient.connect(
    url,
    { useNewUrlParser: true },
  )

  const db = client.db(DB_NAME)
  await db.createCollection('associations')

  const index = await db
    .collection('associations')
    .createIndex({ chain: 1, address: 1 }, { unique: true })

  const documents = [
    {
      chain: 'ethereum',
      address: '787dff5a56cf30d676e45d8de4518c03c335386e'.toUpperCase(),
      associatedAddress: '0x59B8515E7fF389df6926Cd52a086B0f1f46C630A',
    },
  ]

  await db.collection('associations').insertMany(documents)

  client.close()
}

create()
