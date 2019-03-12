#!/usr/bin/env bash
NETWORK="development"

echo "Query Tokens From dex-smart-contract"
node deployment/query_tokens.js $NETWORK

echo "Update config"
go run utils/seed-data/main.go seeds $NETWORK

echo "Drop existing collections"
node deployment/drop_collection.js --collection pairs --network=$NETWORK
node deployment/drop_collection.js --collection tokens --network=$NETWORK
node deployment/drop_collection.js --collection orders --network=$NETWORK
node deployment/drop_collection.js --collection trades --network=$NETWORK
node deployment/drop_collection.js --collection wallets --network=$NETWORK
node deployment/drop_collection.js --collection accounts --network=$NETWORK
node deployment/drop_collection.js --collection config --network=$NETWORK

echo "Create collections"
node deployment/create_accounts_collection.js --network=$NETWORK
node deployment/create_orders_collection.js --network=$NETWORK
node deployment/create_pairs_collection.js --network=$NETWORK
node deployment/create_tokens_collection.js --network=$NETWORK
node deployment/create_trades_collection.js --network=$NETWORK
node deployment/create_wallets_collection.js --network=$NETWORK
node deployment/create_config_collection.js --network=$NETWORK


echo "Seed data"
node deployment/seed_tokens.js --network=$NETWORK
node deployment/seed_pairs.js --network=$NETWORK
node deployment/seed_config.js --network=$NETWORK
node deployment/seed_wallets.js --network=$NETWORK
