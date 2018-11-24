const MongoClient = require('mongodb').MongoClient
const argv = require('yargs').argv
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'
const { DB_NAME } = require('./utils/config')
let client, db, response

const drop = async () => {
  try {
    client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
    db = client.db(DB_NAME)
    response = await db.dropCollection('tokens')
    console.log(response)
  } catch(e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

drop()