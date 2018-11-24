const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'
const { DB_NAME } = require('./utils/config')
(async () => {
  const client = await MongoClient.connect(url, { useNewUrlParser: true })

  const db = client.db(DB_NAME)
  const response = await db.createCollection('trades', {
    validator:  {
      $jsonSchema: 'object',
      properties:  {
        orderHash: {
          bsonType: "string",
        },
        amount: {
          bsonType: "long",
        },
        price: {
          bsonType: "long"
        },
        type: {
          bsonType: "string"
        },
        tradeNonce: {
          bsonType: "string"
        },
        maker: {
          bsonType: "string"
        },
        taker: {
          bsonType: "string"
        },
        takerOrderId: {
          bsonType: "string"
        },
        makerOrderId: {
          bsonType: "string"
        },
        signature: {
          bsonType: "object"
        },
        hash: {
          bsonType: "string"
        },
        pairName: {
          bsonType: "string"
        },
        baseToken: {
          bsontype: "string"
        },
        quoteToken: {
          bsonType: "string"
        }
        createdAt: {
          bsonType: "string"
        },
        updatedAt: {
          bsonType: "string"
        }
      }
      }
    })

  console.log(response)
})()

