const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')
const create = async () => {
  const client = await MongoClient.connect(
    mongoUrl,
    { useNewUrlParser: true },
  )
  console.log('Creating tokens collection')
  const db = client.db(DB_NAME)
  try {
    await db.createCollection('tokens', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['symbol', 'contractAddress', 'decimals'],
          properties: {
            symbol: {
              bsonType: 'string',
              description: 'must be a string and is required',
            },
            contractAddress: {
              bsonType: 'string',
            },
            quote: {
              bsonType: 'bool',
            },
            decimals: {
              bsonType: 'int',
            },
            usd: {
              bsonType: 'string'
            },
            createdAt: {
              bsonType: 'date',
            },
            updatedAt: {
              bsonType: 'date',
            },
          },
        },
      },
    })
  } catch (e) {
    console.log(e)
  } finally {
    client.close()
  }
}

create()
