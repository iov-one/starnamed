
```sh

# Prerequisites

- Linux OS
- Golang 1.18
- Packages:
  - `curl`
  - `jq`
  - `wget`
  - `sha256sum`
  - `systemctl`
  - `journalctl`
  - `grep`
  - `sed`
  - `basename`
  - `chmod`
  - `chgrp`
  - `Git`


## Preparing the environment

- Checkput the starnamed repository

if you'll be syncing from block 1, you'll need to checkout multiple versions of the starnamed repository.  The following commands will checkout the versions of the starnamed repository that are needed to sync from block 1. If not you checkout the latest version of the starnamed repository. 


```sh
# Starnamed versions : v0.10.12, v0.10.13, v0.11.7
mkdir -p ~/starname-binaries

git clone https://github.com/iov-one/starnamed.git
cd starnamed
git checkout tags/v0.10.12
make build
cp build/starnamed ~/starname-binaries/starnamed-v0.10.12

git checkout tags/v0.10.13
make build
cp build/starnamed ~/starname-binaries/starnamed-v0.10.13

git checkout tags/v0.11.7
make build
cp build/starnamed ~/starname-binaries/starnamed-v0.11.7

ls -l ~/starname-binaries
```

if you are syncing from block one you should install the [cosmovisor binary](https://docs.cosmos.network/main/tooling/cosmovisor) .The following commands will install the cosmovisor binary.  If you are not syncing from block one, you can skip this step. The cosmosvisor is a tool that will help you to upgrade the starnamed binary without downtime.

```sh
mkdir -p /starnamed/cosmovisor/genesis/bin
cp ~/starname-binaries/starnamed-v0.10.12 /starnamed/cosmovisor/genesis/bin/starnamed
cp ~/starname-binaries/starnamed-v0.10.13 /starnamed/cosmovisor/upgrades/fix-cosmos-sdk-migrate-bug/bin/starnamed
cp ~/starname-binaries/starnamed-v0.11.7 /starnamed/cosmovisor/upgrades/starname-version-11/bin/starnamed
```

Create the starnamed user and group. Feel free to change the username and groupname.

```sh
sudo groupadd --system starnamed
sudo useradd -s /usr/sbin/nologin --system -g starnamed starnamed
```

Create the starnamed home directory and set the permissions.

```sh
sudo mkdir -p /starnamed
sudo chown -R starnamed:starnamed /starnamed
sudo chmod 750 /starnamed
```

# How to create a linux service

The following commands will create a linux service for starnamed.  You can use this service to start, stop, and restart starnamed.  You can also use this service to enable starnamed to start automatically when the server boots.

```sh 
sudo su

# Get libwasmvm.so
curl https://github.com/CosmWasm/wasmvm/raw/v0.13.0/api/libwasmvm.so > libwasmvm.so
mv -f libwasmvm.so /lib/libwasmvm.so


export USER_IOV=iov
export CHAIN_ID=iov-mainnet-ibc # IBC enabled Starname (IOV) chain id
export DIR_STARNAMED=/starnamed # directory for starnamed related artifacts

# create an environment file for the Starname Asset Name Service
cat <<__EOF_STARNAMED_ENV__ > starnamed.env
# operator variables
CHAIN_ID=${CHAIN_ID}
MONIKER=$(hostname)
SIGNER=${SIGNER}
USER_IOV=${USER_IOV}

# Cosmovisor variables
DAEMON_NAME=starnamed
DAEMON_HOME=${DIR_STARNAMED}
HOME=${DIR_STARNAMED}
DAEMON_RESTART_AFTER_UPGRADE=true

# directories (without spaces to ease pain)
DIR_STARNAMED=${DIR_STARNAMED}
DIR_WORK=${DIR_STARNAMED}

# paths for starnamed and it's required libwasmvm.so
PATH=${PATH}:${DIR_STARNAMED}
LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${DIR_STARNAMED}

__EOF_STARNAMED_ENV__

chgrp ${USER_IOV} starnamed.env
chmod g+r starnamed.env

set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport # pick-up env vars


# create starnamed.service
cat <<__EOF_STARNAMED_SERVICE__ > starnamed.service
[Unit]
Description=Starname Asset Name Service
After=network-online.target

[Service]
Type=simple
User=$(id ${USER_IOV} -u -n)
Group=$(id ${USER_IOV} -g -n)
EnvironmentFile=/etc/systemd/system/starnamed.env
ExecStart=cosmovisor run start --home ${DIR_STARNAMED}
LimitNOFILE=4096
#Restart=on-failure
#RestartSec=3
StandardError=journal
StandardOutput=journal
SyslogIdentifier=starnamed

[Install]
WantedBy=multi-user.target
__EOF_STARNAMED_SERVICE__

systemctl daemon-reload

# initialize the Starname Asset Name Service
su - ${USER_IOV}
set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport # pick-up env vars

# initialize starnamed
starnamed init ${MONIKER} --chain-id ${CHAIN_ID} --home ${DIR_WORK} 2>&1 | jq -r .chain_id

# Download the genesis file
curl --fail https://gist.githubusercontent.com/davepuchyr/6bea7bf369064d118195e9b15ea08a0f/raw/genesis.json > config/genesis.json
sha256sum config/genesis.json | grep e20eb984b3a85eb3d2c76b94d1a30c4b3cfa47397d5da2ec60dca8bef6d40b17 && echo '✅ All good!' || echo "❌ BAD GENESIS FILE!"

exit

systemctl enable starnamed.service
journalctl -f -u starnamed.service & systemctl start starnamed.service

```