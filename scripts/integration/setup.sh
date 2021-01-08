#!/bin/bash
set -o errexit -o nounset -o pipefail

source .env

BINARY=wasmd
PASSWORD=${PASSWORD:-1234567890}
STAKE=${DENOM_STAKE:-ustake}
FEE=${DENOM_FEE:-tiov}
CHAIN_ID=${CHAIN:-testing}
MONIKER=${MONIKER:-node001}

${BINARY} init --chain-id "$CHAIN_ID" "$MONIKER"
# staking/governance token is hardcoded in config, change this
## OSX requires: -i.
sed -i. "s/\"stake\"/\"$STAKE\"/" "$HOME"/.${BINARY}/config/genesis.json
if ! ${BINARY} keys show validator --keyring-backend test; then
  (echo "$PASSWORD"; echo "$PASSWORD") | ${BINARY} keys add validator --keyring-backend test
fi
# hardcode the validator account for this instance
echo "$PASSWORD" | ${BINARY} add-genesis-account validator "1000000000$STAKE,1000000000$FEE" --keyring-backend test
# (optionally) add a few more genesis accounts
for addr in bojack w1 w2 w3 msig1; do
  echo $addr
  ${BINARY} add-genesis-account "$addr" "1000000000$STAKE,1000000000$FEE" --keyring-backend test
done
# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | ${BINARY} gentx validator "250000000$STAKE" --chain-id="$CHAIN_ID" --amount="250000000$STAKE" --keyring-backend test
## should be:
# (echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | ${BINARY} gentx validator "250000000$STAKE" --chain-id="$CHAIN_ID"
${BINARY} collect-gentxs
