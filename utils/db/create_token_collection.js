const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'
const { DB_NAME } = require('./utils/config')
const create = async () => {
  const client = await MongoClient.connect(url, { useNewUrlParser: true })
  const db = client.db(DB_NAME)

  const response = await db.createCollection('tokens', {
    validator:  {
      $jsonSchema: 'object',
      required: [ 'symbol', 'contractAddress', 'decimals'],
      properties:  {
        name: {
          bsonType: "string",
          description: "must be a string and is required"
        }
      }
    }
  }

  console.log(response)
}

create()
