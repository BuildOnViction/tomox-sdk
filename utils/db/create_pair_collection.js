const tokens = require('../tokens.json')
const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'

const create = async () => {
  const client = await MongoClient.connect(url, { useNewUrlParser: true })
  const db = client.db('proofdex')

  const response = await db.createCollection('pairs', {
    validator:  {
      $jsonSchema: 'object',
      required: [ 'baseTokenAddress', 'quoteTokenAddress'],
      properties:  {
        name: {
          bsonType: "string",
          description: "must be a string and is not required"
        },
        baseToken: {
          bsonType: "objectId",
        }
        baseTokenSymbol: {
          bsonType: "string",
          description: "must be a a string and is not required"
        },
        baseTokenAddress: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        quoteToken: {
          bsonType: "objectId"
        }
        quoteTokenSymbol: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        quoteTokenAddress: {
          bsonType: "string",
          description: "must be a string and is required"
        },
        active: {
          bsonType: "bool",
          description: "must be a boolean and is not required"
        },
        makerFee: {
          bsonType: "double",
        },
        takerFee: {
          bsonType: "double"
        },
        createdAt: {
          bsonType: "date"
        },
        updatedAt: {
          bsonType: "date"
        }
      }
      }
    })

  console.log(response)
}

create()

