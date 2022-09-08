#!/bin/bash

set -eo pipefail

SCRIPT_DIR=$(dirname $0)
GO_MOD_FILE=$SCRIPT_DIR/../go.mod
COSMOS_VERSION=$(grep github.com/cosmos/cosmos-sdk $GO_MOD_FILE|cut -d " " -f 2)
SWAGGER_URL=https://raw.githubusercontent.com/cosmos/cosmos-sdk/$COSMOS_VERSION/client/docs/swagger-ui/swagger.yaml

SWAGGER_DEST=$SCRIPT_DIR/../client/docs/swagger-ui/swagger.yaml

wget $SWAGGER_URL -q -O $SWAGGER_DEST 1> /dev/null