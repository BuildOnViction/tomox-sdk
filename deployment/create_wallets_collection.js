const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')
const create = async () => {
  const client = await MongoClient.connect(
    mongoUrl,
    { useNewUrlParser: true },
  )
  console.log('Creating wallets collection')
  const db = client.db(DB_NAME)
  try {
    await db.createCollection('wallets', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['address', 'privateKey'],
          properties: {
            address: {
              bsonType: 'string',
            },
            privateKey: {
              bsonType: 'string',
            },
            admin: {
              bsonType: 'bool',
            },
            operator: {
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
