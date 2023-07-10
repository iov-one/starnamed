#!/bin/sh

# VARS:
CURRENT_DIR=$(pwd)
TEMP_DIR=$(mktemp -d)


# Create a temporary directory to check out the code, build the docker image 

function cleanup {
    rm -rf "$TEMP_DIR"
}

# pull and retag the docker image

docker pull ghcr.io/strangelove-ventures/heighliner/starname:v0.11.7
docker tag ghcr.io/strangelove-ventures/heighliner/starname:v0.11.7 iov-one/starname:v0.11.7

# verify if this the folder with the Dockerfile, if not up one level

if [ ! -f "Dockerfile" ]; then
    cd ..
fi

# Build the current version of the code into a docker image

# docker build -t iov-one/starnamed:v0.12.0 .
cd "$CURRENT_DIR"

# Run the tests
go test -v ./e2e/...