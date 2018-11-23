const MongoClient = require('mongodb').MongoClient
const argv = require('yargs').argv
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017'

let client, db, response

const drop = async () => {
  try {
    client = await MongoClient.connect(mongoUrl, { useNewUrlParser: true })
    db = client.db('proofdex')
    response = await db.dropDatabase()

    client.close()
    console.log(response)
  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

drop()