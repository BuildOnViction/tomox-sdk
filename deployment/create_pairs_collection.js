const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')
const create = async () => {
  const client = await MongoClient.connect(
    mongoUrl,
    { useNewUrlParser: true },
  )
  const db = client.db(DB_NAME)
  console.log('Creating pairs collection')
  try {
    await db.createCollection('pairs', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['baseTokenAddress', 'quoteTokenAddress'],
          properties: {
            baseTokenSymbol: {
              bsonType: 'string',
              description: 'must be a a string and is not required',
            },
            baseTokenAddress: {
              bsonType: 'string',
              description: 'must be a string and is required',
            },
            baseTokenDecimals: {
              bsonType: 'int',
            },
            quoteTokenSymbol: {
              bsonType: 'string',
              description: 'must be a string and is required',
            },
            quoteTokenAddress: {
              bsonType: 'string',
              description: 'must be a string and is required',
            },
            quoteTokenDecimals: {
              bsonType: 'int',
            },
            active: {
              bsonType: 'bool',
              description: 'must be a boolean and is not required',
            },
            makeFee: {
              bsonType: 'string',
            },
            takeFee: {
              bsonType: 'string',
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
