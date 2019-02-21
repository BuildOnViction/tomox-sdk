#!/usr/bin/env bash
echo "Query Tokens From dex-smart-contract"
node utils/db/query_tokens.js tomochainTestnet

echo "Update config"
go run utils/seed-data/main.go seeds

echo "Drop existing collections"
node utils/db/drop_collection.js --collection pairs
node utils/db/drop_collection.js --collection tokens
node utils/db/drop_collection.js --collection orders
node utils/db/drop_collection.js --collection trades
node utils/db/drop_collection.js --collection wallets
node utils/db/drop_collection.js --collection accounts
node utils/db/drop_collection.js --collection config

echo "Create collections"
node utils/db/create_accounts_collection.js
node utils/db/create_orders_collection.js
node utils/db/create_pairs_collection.js
node utils/db/create_tokens_collection.js
node utils/db/create_trades_collection.js
node utils/db/create_wallets_collection.js
node utils/db/create_config_collection.js


echo "Seed data"
node utils/db/seed_tokens.js
node utils/db/seed_quotes.js
node utils/db/seed_pairs.js
node utils/db/seed_config.js
node utils/db/seed_wallets.js
