#!/bin/bash

node -r esm genesis testnet


rm ~/.starnamed/config/*.json
cp data/stargatenet/config/*.json ~/.starnamed/config/

starnamed unsafe-reset-all

starnamed start
