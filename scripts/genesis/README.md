# Add A (Validator) Node To iov-mainnet-ibc

Execute the following commands to light-up a node for **iov-mainnet-ibc**.  Be sure to change the custom environment variables `USER_IOV`, the user that will run the node, and `SIGNER`, the key that will sign the `create-validator` transaction.  `USER_IOV` must exist before you can proceed.  If your node is not a validator then `SIGNER` is inconsequential.


```sh
sudo su -c bash # make life easier for the next ~100 lines

for i in basename chgrp chmod curl grep journalctl jq sed sha256sum systemctl wget ; do [[ $(command -v $i) ]] || { echo "❌ $i is not in PATH; PATH == $PATH; cannot proceed" ; exit -1 ; } ; echo "✅ $i" ; done # check for necessary apps

cd /etc/systemd/system

# custom variables - use values appropriate for your setup
export USER_IOV=iov # "iov" is not recommended
export SIGNER=dave*iov # signer for the create-validator tx

# assert that USER_IOV exists
id ${USER_IOV} && echo '✅ All good!' || echo "❌ ${USER_IOV} does not exist."

# constants
export CHAIN_ID=iov-mainnet-ibc # IBC enabled Starname (IOV) chain id
export DIR_STARNAMED=/opt/iovns/bin # directory for starnamed related artifacts

# create an environment file for the Starname Asset Name Service
cat <<__EOF_STARNAMED_ENV__ > starnamed.env
# operator variables
CHAIN_ID=${CHAIN_ID}
MONIKER=$(hostname)
SIGNER=${SIGNER}
USER_IOV=${USER_IOV}

# directories (without spaces to ease pain)
DIR_STARNAMED=${DIR_STARNAMED}
DIR_WORK=/home/${USER_IOV}/${CHAIN_ID}

# paths for starnamed and it's required libwasmvm.so
PATH=${PATH}:${DIR_STARNAMED}
LD_LIBRARY_PATH=${LD_LIBRARY_PATH}:${DIR_STARNAMED}

# artifacts
STARNAMED=https://github.com/iov-one/starnamed/releases/download/v0.10.12/starnamed-0.10.12-linux-amd64.tar.gz
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
ExecStart=${DIR_STARNAMED}/starnamed.sh
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

# download gitian built binary; starnamed is the Starname Asset Name Service daemon
mkdir -p ${DIR_STARNAMED} && cd ${DIR_STARNAMED}
wget -c ${STARNAMED} && sha256sum $(basename ${STARNAMED}) | grep 505aaddb5a8390576a6f0315627387a64a5682133ddc07b612fd2e05f1280bad && tar xvf $(basename ${STARNAMED}) && echo '✅ All good!' || echo '❌ BAD BINARY!'

# create starnamed.sh, a wrapper for starnamed
cat <<__EOF_STARNAMED_SH__ > starnamed.sh
#!/bin/bash

exec starnamed start \\
  --home=${DIR_WORK} \\
  --minimum-gas-prices='1.0uiov' \\
  --moniker='${MONIKER}' \\
  --p2p.laddr='tcp://0.0.0.0:16656' \\
  --p2p.seeds='0a550a22b027e05206436831a4ec74ccb80feca5@167.99.194.126:16656' \\
  --rpc.laddr='tcp://127.0.0.1:16657' \\
  --rpc.unsafe=true \\

__EOF_STARNAMED_SH__

chgrp ${USER_IOV} starnamed.sh
chmod a+x starnamed.sh

# initialize the Starname Asset Name Service
su - ${USER_IOV}
set -o allexport ; source /etc/systemd/system/starnamed.env ; set +o allexport # pick-up env vars

rm -rf ${DIR_WORK} && mkdir -p ${DIR_WORK} && cd ${DIR_WORK}

# initialize starnamed
starnamed init ${MONIKER} --chain-id ${CHAIN_ID} --home ${DIR_WORK} 2>&1 | jq -r .chain_id

# customize ${DIR_WORK}/{app.toml,config.toml} if you want

# get the genesis file
curl --fail https://gist.githubusercontent.com/davepuchyr/6bea7bf369064d118195e9b15ea08a0f/raw/genesis.json > config/genesis.json
sha256sum config/genesis.json | grep e20eb984b3a85eb3d2c76b94d1a30c4b3cfa47397d5da2ec60dca8bef6d40b17 && echo '✅ All good!' || echo "❌ BAD GENESIS FILE!"

exit # ${USER_IOV}

systemctl enable starnamed.service
journalctl -f -u starnamed.service & systemctl start starnamed.service # watch the chain sync

exit # root
```
