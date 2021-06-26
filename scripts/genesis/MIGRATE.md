# Migrating iov-mainnet-2 To IBC-Enabled iov-mainnet-ibc

This document assumes that you followed the procedure [here](https://docs.starname.me/for-validators/mainnet) when you joined **iov-mainnet-2** as a validator and that you're going to use the same node for your validator on **iov-mainnet-ibc**.  It will detail the procedure necessary to export the state from **iov-mainnet-2**, migrate the state to **iov-mainnet-ibc**, and verify the migrated genesis file.  It requires that certain apps, including `node` and `yarn`, are installed on your system, and user `$USER_IOV` exists.

The migration process consists of three parts:
1. Making your validator halt at block 4,294,679.
1. Setting up your validator for `iov-mainnet-ibc`.
1. Exporting the state of `iov-mainnet-2` at block 4,294,679, migrating state, and verifying the migrated genesis file.

The first two steps of the process above can be done at any time.  The last one can only be done when block 4,294,679 on **iov-mainnet-2** is committed, which should happen around 8am UTC on June 29th, 2021.


## Step 1 of 3: Making your validator halt at block 4,294,679.

This one is easy and quick.  On your validator node execute

```sh
# become root
sudo su -c bash

# pick-up env vars
set -o allexport ; source /etc/systemd/system/starname.env ; set +o allexport

# inject the halt height into iovnsd.sh
sed --in-place 's/iovnsd start/iovnsd start --halt-height 4294679/' ${DIR_IOVNS}/iovnsd.sh

# restart the service
systemctl restart starname.service

exit
```


## Step 2 of 3: Setting up your `iov-mainnet-ibc` validator.

On your validator node, follow the procedure [here](README.md).  Expect `❌ BAD GENESIS FILE!` near the end of the procedure since the hash of the migrated genesis file is not yet known.  Don't bother starting starnamed.service.

**Make sure the private key of your validator for `iov-mainnet-2` matches your key for `iov-mainnet-ibc`.  If you use an HSM you know how to do that.  If you don't use an HSM then simply execute the following on your validator node**

```sh
# become root
sudo su -c bash

# pick-up iov-mainnet-2 env vars
set -o allexport ; source /etc/systemd/system/starname.env ; set +o allexport # iov-mainnet-2
export OLD_DIR_WORK=${DIR_WORK}

# pick-up iov-mainnet-ibc env vars
set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport # iov-mainnet-ibc

# copy your validator's private key
cp -av ${OLD_DIR_WORK}/config/priv_validator_key.json ${DIR_WORK}/config

exit
```


## Step 3 of 3: Exporting the state of `iov-mainnet-2` at block 4,294,679, migrating state, and verifying the migrated genesis file.

Once **iov-mainnet-2** reaches height 4,294,679 your validator will gracefully stop if Step 1 above was executed sucessfully.  Now you can export state by executing the following on your validator node

```sh
# become ${USER_IOV}
su - ${USER_IOV}

# pick-up iov-mainnet-2 env vars
set -o allexport ; source /etc/systemd/system/starname.env ; set +o allexport # iov-mainnet-2

${DIR_IOVNS}/iovnsd export --home=${DIR_WORK} > iov-mainnet-2.json # add --height 4294679 if you want but it should be that by virtue of the halt-height
```

**Wait for the word that https://github.com/iov-one/starnamed.git is tagged `iov-mainnet-ibc-genesis`.**

Next you have to migrate the state.  Here's where `node` and `yarn` enter the mix

```sh
# pick-up iov-mainnet-ibc env vars
set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport # iov-mainnet-ibc

git clone --recursive https://github.com/iov-one/starnamed.git
diff iov-mainnet-2.json starnamed/scripts/genesis/data/iov-mainnet-2.json && echo '✅ All good!' || echo '❌ Exported state mismatch!'
cd starnamed
git submodule foreach git checkout master
cd scripts/genesis
yarn
# export available ports for starnamed for `yarn test` (customize them if necessary)
export PORT_P2P=tcp://127.0.0.1:16656
export PORT_RPC=tcp://127.0.0.1:16657
export PORT_GRPC=tcp://127.0.0.1:16090
yarn test # takes 3 minutes on my laptop; if the tests fail then change
node -r esm genesis.js iov-mainnet-ibc # takes 2 minutes
cd data/iov-mainnet-ibc/config
git diff && cp -av genesis.json ${DIR_WORK}/config && echo '✅ All good!' || echo '❌ BAD genesis file!'
exit # $USER_IOV
```

That's it for the migration!


<hr style="width:100%;"></hr>

## Light-Up Your `iov-mainnet-ibc` Validator.

All that's left to do now is start your validator

```sh
sudo systemctl restart starnamed.service
```


## Get Ease Of Mind Prior To Block 4,294,679.

You can test your new validator's setup even before block 4,294,679 by doing the following

```sh
# become ${USER_IOV}
su - ${USER_IOV}

# pick-up iov-mainnet-ibc env vars
set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport # iov-mainnet-ibc

# use the pre-block 4,294,679 sample iov-mainnet-ibc genesis file
cd starnamed/scripts/genesis/data/iov-mainnet-ibc/config && cp -av genesis.json ${DIR_WORK}/config

# become root
sudo su -c bash

# breifly start starnamed.service; it won't produce blocks because the network (seed node) is offline
systemctl restart starnamed.service
while ! journalctl -u starnamed.service --no-pager | grep 'This node is a ' ; do sleep 1 ; done
systemctl stop starnamed.service

# determine if your node is a validator
journalctl -u starnamed.service --no-pager | grep 'This node is a validator' && echo '✅ You are golden!' || echo '❌ BAD validator private key!'

exit # root

# reset all
starnamed unsafe-reset-all --home ${DIR_WORK} # *** THIS IS CRUCIAL ***

exit # ${USER_IOV}
```

## Audit The Migration Process. ##

`node -r esm genesis.js iov-mainnet-ibc` is where the migration from **iov-mainnet-2**'s state to **iov-mainnet-ibc**'s genesis file takes place. [Check it out](genesis.js).
