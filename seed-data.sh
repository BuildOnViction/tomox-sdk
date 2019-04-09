#!/usr/bin/env bash
NETWORK="development"

echo "Update config"
go run utils/seed-data/main.go seeds $NETWORK
