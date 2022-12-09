# Add A Node To iov-mainnet-ibc

Execute the following commands to light-up a node for **iov-mainnet-ibc**.  Be sure to change the custom environment variable `USER_IOV`, the user that will run the node.  `USER_IOV` must exist before you can proceed.


## Sync'ing Using `statesync`

TDB


## Sync'ing From Block 1

The chain has been upgraded twice since its inception, which means that three different binaries are needed to sync the chain from block 1.  Sync'ing from block 1 is broken into three phases and will take days.

### Phase 1 Of 3

```sh
sudo su -c bash # make life easier for the next ~100 lines

for i in basename chgrp chmod curl grep journalctl jq sed sha256sum systemctl wget ; do [[ $(command -v $i) ]] || { echo "❌ $i is not in PATH; PATH == $PATH; cannot proceed" ; exit -1 ; } ; echo "✅ $i" ; done # check for necessary apps

cd /etc/systemd/system

# custom variables - use values appropriate for your setup
export USER_IOV=iov # "iov" is not recommended

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
  --p2p.persistent_peers='ca133187b37b59d2454812cfcf31b6211395adec@167.99.194.126:16656, 1c7e014b65f7a3ea2cf48bffce78f5cbcad2a0b7@13.37.85.253:26656, 8c64a2127cc07d4570756b61f83af60d34258398@13.37.61.32:26656, 9aabe0ac122f3104d8fc098e19c66714c6f1ace9@3.37.140.5:26656, faedef1969911d24bf72c56fc01326eb891fa3b7@63.250.53.45:16656, 94ac1c02b4e2ca3fb2706c91a68b8030ed3615a1@35.247.175.128:16656, be2235996b1c785a9f57eed25fd673ca111f0bae@52.52.89.64:26656, f63d15ab7ed55dc75f332d0b0d2b01d529d5cbcd@212.71.247.11:26656, f5597a7ed33bc99eb6ba7253eb8ac76af27b4c6d@138.201.20.147:26656' \\
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
journalctl -f -u starnamed.service & systemctl start starnamed.service # watch the chain sync until it intentionlly panics after block 4597999

exit # root
```

### Phase 2 Of 3

[Proposal 9](https://big-dipper.iov-mainnet-ibc.iov.one/proposals/9), a software upgrade, forced `starnamed` to intentionally panic after block 4,597,999.  The `starnamed` binary needs to be upgraded in order to continue sync'ing `iov-mainnet-ibc`.

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

journalctl -f -u starnamed.service & systemctl restart starnamed.service # watch the chain sync  until it intentionlly panics after block xxxxxxx


### Phase 3 Of 3

[Proposal x](https://big-dipper.iov-mainnet-ibc.iov.one/proposals/x), a software upgrade, forced `starnamed` to intentionally panic after block x,xxx,xxx.  The `starnamed` binary needs to be upgraded in order to continue sync'ing `iov-mainnet-ibc`.

exit # sudo su
```

## Contributors

Thanks to the current and former employees of IOV SAS and all the open-source developers in the Cosmos community!

* Adrien Duval [LeCodeurDuDimanche](https://github.com/LeCodeurDuDimanche)
* Jacob Gadikian [faddat](https://github.com/faddat)
* Frojdi Dymylja [fdymylja](https://github.com/fdymylja)
* Orkun Külçe [orkunkl](https://github.com/orkunkl)
* Ethan Frey [ethanfrey](https://github.com/ethanfrey)
* Simon Warta [webmaster128](https://github.com/webmaster128)
* Alex Peters [alpe](https://github.com/alpe)
* Aaron Craelius [aaronc](https://github.com/aaronc)
* Sunny Aggarwal [sunnya97](https://github.com/sunnya97)
* Cory Levinson [clevinson](https://github.com/clevinson)
* Sahith Narahari [sahith-narahari](https://github.com/sahith-narahari)
* Jehan Tremback [jtremback](https://github.com/jtremback)
* Shane Vitarana [shanev](https://github.com/shanev)
* Billy Rennekamp [okwme](https://github.com/okwme)
* Westaking [westaking](https://github.com/westaking)
* Marko [marbar3778](https://github.com/marbar3778)
* JayB [kogisin](https://github.com/kogisin)
* Rick Dudley [AFDudley](https://github.com/AFDudley)
* KamiD [KamiD](https://github.com/KamiD)
* Valery Litvin [litvintech](https://github.com/litvintech)
* Leonardo Bragagnolo [bragaz](https://github.com/bragaz)
