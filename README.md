# Starnamed - A Cosmos SDK blockchain

[![Go Report Card](https://goreportcard.com/badge/github.com/iov-one/starnamed)](https://goreportcard.com/report/github.com/iov-one/starnamed)
[![GoDoc](https://godoc.org/github.com/iov-one/starnamed?status.svg)](https://pkg.go.dev/github.com/iov-one/starnamed)
[![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)]()
![Discord](https://img.shields.io/discord/906536286653923408?label=discord&logo=discord)

## The project

Starname is a blockchain platform that allows users to register and manage their own decentralized domain names, which can be used for various purposes such as sending and receiving cryptocurrency, creating unique identifiers for digital assets, and more. Starname is also a decentralized identity platform that enables users to have a universal username for the blockchain world.

## This repository

This repository contains the source code for the Starname blockchain. It is has as the base the code of [Wasmd](https://github.com/CosmWasm/wasmd) repository, which is extended version of [Cosmos SDK](https://github.com/cosmos/cosmos-sdk) repository. Allow the starname blockchain to be a Cosmos SDK blockchain with the ability to run smart contracts written in [Wasm](https://webassembly.org/). 

All the knowledge about how to maintain and use a cosmos-sdk blockchain is applicable to the starname blockchain. The only difference is that the starname blockchain can run smart contracts written in Wasm and the custom modules that are required for the starname blockchain.

## Installation
To install the starnamed CLI, follow these steps:

- Binary installation :
  - (Linux Only) Download the latest release from the v0.11.6 release page and use.

```bash
# Note: execute: sudo echo "enable sudo" to avoid the password prompt 
wget -O starnamed_temp https://github.com/iov-one/starnamed/releases/download/v0.11.7/starnamed.linux.amd64 && chmod +x starnamed_temp && sudo mv starnamed_temp /usr/local/bin/starnamed && starnamed version --long
```

- Source installation : (Any other OS)
  - Install Go 1.18 or higher.
  - Clone the starnamed repository.
  - Run make install to install the starnamed CLI.

- Using the web app:
  - Go to https://app.starname.me/ and use the web app.



## Usage

Once you have installed the starnamed CLI, you can use it to manage your Starname account and perform various operations on the Starname blockchain. Some examples of commands that you can use are:

- Create your key pair:
  - `starnamed keys add <your_key_name>` *This command will generate a key pair and store it in your local keyring. You will need to use this key pair to sign transactions that you send to the Starname blockchain.* **Copy the mnemonic phrase and keep it safe. You will need it to recover your key pair if you lose it.**

- Adding tokens to your account:
  - you can swap at the Osmosis swap pool (https://app.osmosis.zone/swap) and transfer the tokens to your account using the IBC protocol. The osmosis app has all the necessary information to do that.

- Using the starname module:
  - `starnamed tx starname --help` *This command will show you all the available commands that you can use to manage your Starname account.*
    - Create domain: `starnamed tx starname domain-register [flags]`
    - Create account: `starnamed tx starname account-register [flags]`
    - Create a escrow: **To use this command you must own a domain or account.** 
      - `starnamed tx starname domain-escrow-create [flags]`
      - `starnamed tx starname account-escrow-create [flags]`
    - You can also, all information to the account and domain, renew, transfer, add a certificate or delete the account or domain.


## Documentation

- [Starname Documentation](https://docs.starname.me/) (Can contain outdated information, that'll be updated)
- [Node operator](./docs/node-operator.md)
- [Swagger API documentation]( https://iov-one.github.io/starnamed/)
- [Developer documentation](./docs/developer.md)
- [Testnet documentation](./docs/testnet.md) 

## External links
- [Cosmos SDK Documentation](https://docs.cosmos.network/)
- [Cosmos SDK Tutorials](https://tutorials.cosmos.network/)

## Usage


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
* Davi Mello  [dsmello](https://github.com/dsmello)