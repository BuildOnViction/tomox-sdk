const MongoClient = require('mongodb').MongoClient;
const argv = require('yargs').argv;
const { DB_NAME, mongoUrl } = require('./utils/config');
let client, db;

const query = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true }
    );
    db = client.db(DB_NAME);

    const response = await db
      .collection('associations')
      .find({
        chain: 'ethereum',
        address: '787DFF5A56CF30D676E45D8DE4518C03C335386E'
      })
      .toArray();
    console.log(response);
  } catch (e) {
    console.log(e.message);
  } finally {
    client.close();
  }
};

query();
