const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'

const create = async () => {
  const client = await MongoClient.connect(url, { useNewUrlParser: true })

  const db = client.db('proofdex')
  const response = await db.createCollection('balances', {
    validator:  {
      $jsonSchema: 'object',
      required: ['address'],
      properties:  {
        address: {
          bsonType: "string",
        },
        tokens: {
          bsonType: "object",
        },
        createdAt: {
          bsonType: "long"
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

