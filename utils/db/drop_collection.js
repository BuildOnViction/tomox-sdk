const MongoClient = require('mongodb').MongoClient;
const argv = require('yargs').argv;
const collection = argv.collection;
const { DB_NAME, mongoUrl } = require('./utils/config');
let client, db, response;

const drop = async () => {
  console.log('Dropping ' + collection + ' collections');
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true }
    );
    db = client.db(DB_NAME);
    response = await db.dropCollection(collection);
    console.log(response);
  } catch (e) {
    console.log(e.message);
  } finally {
    client.close();
  }
};

drop();
