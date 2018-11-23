const MongoClient = require('mongodb').MongoClient
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017'

let client, db

const query = async () => {
  try {
    client = await MongoClient.connect(url, { useNewUrlParser: true })
    db = client.db('proofdex')

    const pairs = await db.collection('pairs').find().toArray()
    const pair = pairs[0]
    const query = {
      "status": { $in: [ "OPEN", "PARTIALLY_FILLED" ]},
      "baseToken": pair.baseTokenAddress,
      "quoteToken": pair.quoteTokenAddress
    }

    const response = await db.collection('orders').find(query).toArray()
    console.log(response)

  } catch (e) {
    console.log(e.message)
  } finally {
    client.close()
  }
}

query()