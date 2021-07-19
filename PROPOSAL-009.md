# Software Upgrade Proposal 9 #

The public keys for multisig accounts were not migrated from the exported state of `iov-mainnet-2` into `iov-mainnet-ibc`. In order for multisig accounts to work on `iov-mainnet-ibc` the multisig public keys need to be injected into the app state.  https://big-dipper.iov-mainnet-ibc.iov.one/proposals/9 aims to do that.  It will use the cosmos-sdk `upgrade` module to intentionally cause the `starnamed` binary to panic according to the proposal's plan:

```sh
$ starnamed query gov proposal 9 --chain-id iov-mainnet-ibc --output json | jq .content.plan
{
  "name": "fix-cosmos-sdk-migrate-bug",
  "time": "0001-01-01T00:00:00Z",
  "height": "4598000",
  "info": "v0.10.13-9c4ac24779bfb4d52bb5920db1fa0cde0ec811b6"
}
```

Note that the intentional panic will occur at block 4598000 at approximately 10am on July 20, 2021.

**Validators and node operators must upgrade their `starnamed` binary to `v0.10.13` and restart it.**

## Upgrading if you followed the "standard" procedure for lighting-up your starnamed node ##

If you followed the procedure [here](README.md) then you can use the following script to upgrade the `starnamed` binary **after** block 4598000 exists.

```sh
# make life easier for the next ~20 lines
sudo su -c bash

# replace the old starnamed artifact with the new one in starnamed.env
sed --in-place 's/0.10.12/0.10.13/g' /etc/systemd/system/starnamed.env

# pick-up env vars
set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport

# cd to where all the action is going to happen
cd ${DIR_STARNAMED}

# move the old binaries out of the way of the new ones
mv -v starnamed starnamed-v0.10.12
mv -v libwasmvm.so libwasmvm.so-v0.10.12

# get the new starnamed binary
wget -c ${STARNAMED} && sha256sum $(basename ${STARNAMED}) | grep a3555955a1d001449d7e05793852ea23064905614bab0bb8cefc250758ea81bf && tar xvf $(basename ${STARNAMED}) && echo '✅ All good!' || echo '❌ BAD BINARY!'

journalctl -f -u starnamed.service & systemctl restart starnamed.service # wait for more than 2/3rds of the voting power to come online

exit # sudo su
```

## Upgrading if you didn't follow the "standard" procedure ##

If you didn't follow the procedure [here](README.md) then you're on your own.  I don't know if `cosmvisor` will work, though, in theory, it should.  If you're a validator for the Cosmos Hub then what we're doing here is exactly the same procedure that introduced the Gravity DEX into the Cosmos Hub.
