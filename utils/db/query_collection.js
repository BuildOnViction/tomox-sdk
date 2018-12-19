const MongoClient = require('mongodb').MongoClient;
const argv = require('yargs').argv;
const collection = argv.collection;
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
      .collection(collection)
      .find()
      .toArray();
    console.log(response);
  } catch (e) {
    console.log(e.message);
  } finally {
    client.close();
  }
};

query();
