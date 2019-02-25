const MongoClient = require('mongodb').MongoClient
const { DB_NAME, mongoUrl } = require('./utils/config')

const create = async () => {
  const client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
  console.log('Creating orders collection')
  const db = client.db(DB_NAME)
  try {
    await db.createCollection('orders', {
      validator: {
        $jsonSchema: {
          bsonType: 'object',
          required: [
            'baseToken',
            'quoteToken',
            'amount',
            'pricepoint',
            'userAddress',
            'exchangeAddress',
            'filledAmount',
          ],
          properties: {
            baseToken: {
              bsonType: 'string',
            },
            quoteToken: {
              bsonType: 'string',
            },
            filledAmount: {
              bsonType: 'string',
            },
            amount: {
              bsonType: 'string',
            },
            pricepoint: {
              bsonType: 'string',
            },
            makeFee: {
              bsonType: 'string',
            },
            takeFee: {
              bsonType: 'string',
            },
            side: {
              bsonType: 'string',
            },
            status: {
              bsonType: 'string',
            },
            exchangeAddress: {
              bsonType: 'string',
            },
            userAddress: {
              bsonType: 'string',
            },
            signature: {
              bsonType: 'object',
            },
            nonce: {
              bsonType: 'string',
            },
            pairName: {
              bsonType: 'string',
            },
            hash: {
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
