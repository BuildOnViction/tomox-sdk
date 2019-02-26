const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')
const create = async () => {
  const client = await MongoClient.connect(
    mongoUrl,
    { useNewUrlParser: true },
  )
  console.log('Creating accounts collection')
  const db = client.db(DB_NAME)
  try {
    await db.createCollection('accounts', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['address'],
          properties: {
            address: {
              bsonType: 'string',
            },
            tokenBalances: {
              bsonType: 'object',
            },
            isBlocked: {
              bsonType: 'bool',
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
