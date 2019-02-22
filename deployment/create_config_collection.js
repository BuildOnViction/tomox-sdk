const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')
const create = async () => {
  const client = await MongoClient.connect(
    mongoUrl,
    { useNewUrlParser: true },
  )
  console.log('Creating config collection')
  const db = client.db(DB_NAME)
  try {
    await db.createCollection('config', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: ['key'],
          properties: {
            key: {
              bsonType: 'string',
            },
            value: {
              bsonType: [
                'int',
                'long',
                'string',
                'array',
                'bool',
                'date',
                'object',
              ],
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
