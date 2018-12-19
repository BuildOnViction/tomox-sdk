const MongoClient = require('mongodb').MongoClient;
const { utils, Wallet } = require('ethers');
const { getNetworkID } = require('./utils/helpers');
const {
  DB_NAME,
  keys,
  mongoUrl,
  network,
  baseTokens,
  contractAddresses
} = require('./utils/config');
const networkID = getNetworkID(network);
const walletKeys = keys[networkID];

let addresses = contractAddresses[networkID];

let client, db, documents, response;

const seed = async () => {
  try {
    client = await MongoClient.connect(
      mongoUrl,
      { useNewUrlParser: true }
    );
    db = client.db(DB_NAME);
    documents = [];

    walletKeys.forEach(key => {
      let accountRecord = {
        isBlocked: false,
        tokenBalances: {}
      };
      account = new Wallet(key);

      accountRecord.address = utils.getAddress(account.address);
      baseTokens.forEach(symbol => {
        const contractAddress = utils.getAddress(addresses[symbol]);
        accountRecord.tokenBalances[contractAddress] = {
          address: contractAddress,
          allowance: '10000000000000000000000000000',
          balance: '10000000000000000000000000000',
          lockedBalance: '0',
          symbol: symbol
        };
      });

      documents.push(accountRecord);
    });

    response = await db.collection('accounts').insertMany(documents);
    console.log(response);
  } catch (e) {
    console.log(e.message);
  } finally {
    client.close();
  }
};

seed();
