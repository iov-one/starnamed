#!/bin/sh
#set -o errexit -o nounset -o pipefail

PASSWORD=${PASSWORD:-12341234}
STAKE=${STAKE_TOKEN:-ustake}
FEE=${FEE_TOKEN:-tiov}
CHAIN_ID=${CHAIN_ID:-testing}
MONIKER=${MONIKER:-node001}

rm -rf ~/.starnamed/config/gentx
starnamed tendermint unsafe-reset-all
starnamed init --chain-id "$CHAIN_ID" "$MONIKER" -o
# staking/governance token is hardcoded in config, change this
sed -i "s/\"stake\"/\"$STAKE\"/" "$HOME"/.starnamed/config/genesis.json
# this is essential for sub-1s block times (or header times go crazy)
sed -i 's/"time_iota_ms": "1000"/"time_iota_ms": "10"/' "$HOME"/.starnamed/config/genesis.json


if ! starnamed keys show validator; then
  (echo "$PASSWORD"; echo "$PASSWORD") | starnamed keys add validator
fi
# hardcode the validator account for this instance
echo "$PASSWORD" | starnamed add-genesis-account validator "1000000000$STAKE,1000000000$FEE"

# (optionally) add a few more genesis accounts
for addr in wasm1fjppc038udty5lquva2fc72967y4mchsu06slw "$@"; do
  echo $addr
  starnamed add-genesis-account "$addr" "1000000000$STAKE,1000000000$FEE"
done

# submit a genesis validator tx
## Workraround for https://github.com/cosmos/cosmos-sdk/issues/8251
(echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | starnamed gentx validator "250000000$STAKE" --chain-id="$CHAIN_ID" --amount="250000000$STAKE"
## should be:
# (echo "$PASSWORD"; echo "$PASSWORD"; echo "$PASSWORD") | starnamed gentx validator "250000000$STAKE" --chain-id="$CHAIN_ID"
starnamed collect-gentxs
