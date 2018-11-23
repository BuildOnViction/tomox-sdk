const MongoClient = require('mongodb').MongoClient
const argv = require('yargs').argv
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'
const collection = argv.collection

let client, db

const query = async () => {
  try {
    client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
    db = client.db('proofdex')

    const response = await db.collection(collection).find().toArray()
    console.log(response)

  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

query()