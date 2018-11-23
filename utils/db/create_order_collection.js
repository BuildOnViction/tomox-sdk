const tokens = require('../tokens.json')
const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'

const create = async () => {
  const client = await MongoClient.connect(url, { useNewUrlParser: true })

  const db = client.db('proofdex')
  const response = await db.createCollection('orders', {
    validator:  {
      $jsonSchema: 'object',
      required: [
         'baseToken',
         'quoteToken',
         'amount',
         'pricepoint',
         'userAddress'
         'exchangeAddress',
         'filledAmount',
         'amount',
        ],
      properties:  {
        baseToken: {
          bsonType: "string",
        },
        quoteToken: {
          bsonType: "string",
        },
        filledAmount: {
          bsonType: "long"
        },
        amount: {
          bsonType: "long"
        },
        pricepoint: {
          bsonType: "long"
        },
        makeFee: {
          bsonType: "long"
        },
        takeFee: {
          bsonType: "long"
        },
        side: {
          bsonType: "string"
        },
        exchangeAddress: {
          bsonType: "string"
        },
        userAddress: {
          bsonType: "string"
        },
        signature: {
          bsonType: "object"
        },
        nonce: {
          bsonType: 'string'
        }
        pairName: {
          bsonType: "string"
        },
        hash: {
          bsonType: "string"
        },
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
}

create()