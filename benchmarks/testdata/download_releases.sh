#!/bin/bash
set -o errexit -o nounset -o pipefail
command -v shellcheck > /dev/null && shellcheck "$0"

if [ $# -ne 1 ]; then
  echo "Usage: ./download_releases.sh RELEASE_TAG"
  exit 1
fi

tag="$1"

for contract in cw20_base cw1_whitelist; do
  url="https://github.com/CosmWasm/cw-plus/releases/download/$tag/${contract}.wasm"
  echo "Downloading $url ..."
  wget -O "${contract}.wasm" "$url"
done

<<<<<<<< HEAD:benchmarks/testdata/download_releases.sh
========
# create the zip variant
gzip -k hackatom.wasm
mv hackatom.wasm.gz hackatom.wasm.gzip

>>>>>>>> tags/v0.11.6:x/wasm/keeper/testdata/download_releases.sh
rm -f version.txt
echo "$tag" >version.txt