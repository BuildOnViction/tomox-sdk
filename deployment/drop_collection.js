const MongoClient = require('mongodb').MongoClient
const argv = require('yargs').argv
const collection = argv.collection
const { DB_NAME, mongoUrl } = require('./utils/config')
let client, db

const drop = async () => {
  console.log('Dropping ' + collection + ' collections')
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true },
    )
    db = client.db(DB_NAME)
    await db.dropCollection(collection)
  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

drop()
