const MongoClient = require('mongodb').MongoClient;
const url = process.env.MONGODB_URL || 'mongodb://localhost:27017';
const { DB_NAME } = require('./utils/config');
const create = async () => {
  const client = await MongoClient.connect(
    url,
    { useNewUrlParser: true }
  );

  const db = client.db(DB_NAME);
  const response = await db.createCollection('accounts', {
    validator: {
      $jsonSchema: 'object',
      required: ['address'],
      properties: {
        address: {
          bsonType: 'string'
        },
        tokenBalances: {
          bsonType: 'object'
        },
        isBlocked: {
          bsonType: 'bool'
        }
      }
    }
  });

  console.log(response);

  client.close();
};

create();
