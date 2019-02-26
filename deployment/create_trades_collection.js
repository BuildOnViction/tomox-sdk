const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')

const create = async () => {
  const client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
  console.log('Creating trades collection')
  const db = client.db(DB_NAME)
  try {
    await db.createCollection('trades', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: [
            'baseToken',
            'quoteToken',
            'amount',
            'maker',
            'taker',
          ],
          properties: {
            amount: {
              bsonType: 'string',
            },
            pricepoint: {
              bsonType: 'string',
            },
            status: {
              bsonType: 'string',
            },
            maker: {
              bsonType: 'string',
            },
            taker: {
              bsonType: 'string',
            },
            takerOrderHash: {
              bsonType: 'string',
            },
            makerOrderHash: {
              bsonType: 'string',
            },
            hash: {
              bsonType: 'string',
            },
            txHash: {
              bsonType: 'string',
            },
            pairName: {
              bsonType: 'string',
            },
            baseToken: {
              bsonType: 'string',
            },
            quoteToken: {
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
