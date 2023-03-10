#!/usr/bin/env bash

set -eo pipefail

WORKDIR=$(pwd)

mkdir -p ./tmp-swagger-gen
proto_dirs=$(find ./proto -path -prune -o -name '*.proto' -print0 | xargs -0 -n1 dirname | sort | uniq)
for dir in $proto_dirs; do
  dir_temp=./tmp-swagger-gen/$(echo $dir|  sed  's/\.\///g' | sed 's/\//-/g')
  mkdir -p $dir_temp

  # generate swagger files (filter query files)
  query_file=$(find "${dir}" -maxdepth 1 -name 'query.proto')
  if [[ ! -z "$query_file" ]]; then
    protoc  \
    -I "proto" \
    -I "third_party/proto" \
    "$query_file" \
    --openapi_out $dir_temp
  fi
done

# combine swagger files
# uses nodejs package `swagger-combine`.
# all the individual swagger files need to be configured in `config.json` for merging
mkdir -p ./client/docs/swagger-ui/
swagger-combine ./client/docs/config.json -o ./client/docs/swagger-ui/swagger-starname.yaml -f yaml --continueOnConflictingPaths true --includeDefinitions true

# clean swagger files
rm -rf ./tmp-swagger-gen

# Update the Cosmos SDK yaml file
scripts/fetch-cosmos-swagger.sh $WORKDIR

# Update static assets
statik -src client/docs/swagger-ui/ -dest client/docs/ -f -m

