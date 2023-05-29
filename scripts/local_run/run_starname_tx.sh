#!/bin/bash

# This script connect to local node and perform some basic txs.

# Variables:

ACCOUNT_NAME=${1:-"user"}
FOLDER_PATH=${1:-"./tmp/node0"}
DEFAULT_DENOM=${2:-"stake"}
DEFAULT_DENOM_SECONDARY="tusd"
NODE_NAME=${3:-"node0"}
CHAIN_ID="localnet"
DOMAIN_NAME="localnode"

# Register a starname domain:
# example: ./build/starnamed --home tmp/node0/ tx starname domain-register -y -d localnode --from validator --keyring-backend test --chain-id=localnet

echo "Registering a starname domain..."
./build/starnamed --home $FOLDER_PATH/ tx starname domain-register -y -d $DOMAIN_NAME --from $ACCOUNT_NAME --keyring-backend test --chain-id=$CHAIN_ID 

# delay 6 seconds
sleep 6

# Register an escrow domain:
# example: ./build/starnamed --home tmp/node0/  -y - --from validator --keyring-backend test --chain-id=localnet tx starname dec -d localnode --price 100tusd --expiration 2023-05-25T02:04:05Z

echo "Registering an escrow domain..."
./build/starnamed --home $FOLDER_PATH/ tx starname domain-escrow-create -y -d $DOMAIN_NAME --price 100$DEFAULT_DENOM_SECONDARY --expiration $(date -u +"%Y-%m-%dT%H:%M:%SZ" -d "$(date -u +"%Y-%m-%dT%H:%M:%SZ") + 30 minutes") --from $ACCOUNT_NAME --keyring-backend test --chain-id=$CHAIN_ID 

# delay 6 seconds
sleep 6

# The validator account buys the escrow domain:
# example:  ./build/ starnamed tx escrow transfer [id] [amount] [flags]

echo "Buying the escrow domain..."
./build/starnamed --home $FOLDER_PATH/ tx escrow transfer "0000000000000001" 100$DEFAULT_DENOM_SECONDARY --from validator --keyring-backend test --chain-id=$CHAIN_ID -y --gas=500000