#!/bin/sh
# set -eo pipefail

# This script generates the swagger docs for the cosmos-sdk + chain modules
# Usage: ./scripts/generate_swagger_docs.sh <root path of the chain repo>

# 0 - Create paths

SWAGGER_PATH=$1/client/docs/swagger-ui


mkdir -p $SWAGGER_PATH

# 1 - Downlaod the swagger.yaml file from the cosmos-sdk repo at the version specified in the go.mod file

GO_MOD_FILE=$1/go.mod
COSMOS_VERSION=$(grep github.com/cosmos/cosmos-sdk $GO_MOD_FILE|cut -d " " -f 2 | head -n 1)
SWAGGER_URL=https://raw.githubusercontent.com/cosmos/cosmos-sdk/$COSMOS_VERSION/client/docs/swagger-ui/swagger.yaml

SWAGGER_DEST=$1/client/docs/swagger-ui/cosmos-swagger.yaml

wget $SWAGGER_URL -q -O $SWAGGER_DEST 1> /dev/null

# 2 - Use docker to read the ./proto folder and generate the chain-swagger.yaml file for the chain

DOCKER_IMAGE=quay.io/goswagger/swagger
DOCKER_TAG=0.26.1

