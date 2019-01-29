const MongoClient = require('mongodb').MongoClient;
const argv = require('yargs').argv;
const { DB_NAME, mongoUrl } = require('./utils/config');
let client, db;

const drop = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true }
    );
    db = client.db(DB_NAME);
    await db.dropDatabase();

    client.close();
  } catch (e) {
    console.log(e.message);
  } finally {
    client.close();
  }
};

drop();
