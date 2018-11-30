const argv = require('yargs').argv;
const mongoUrl = argv.mongo_url || 'mongodb://localhost:27017';
const user = argv.user;
const pwd = argv.password;
const { DB_NAME } = require('./utils/config');
const MongoClient = require('mongodb').MongoClient;
let client, db;

const create = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true }
    );
    db = client.db(DB_NAME);
    db.addUser({
      username: user,
      password: pwd,
      options: {
        roles: [{ role: 'userAdminAnyDatabase', db: 'admin' }]
      }
    });

    client.close();
  } catch (e) {
    throw new Error(e.message);
  }
};

create();
