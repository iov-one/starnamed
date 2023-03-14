**WIP**

# Node Operator Guide


This guide is for node operators who wish to run a node on the **iov-mainnet-ibc** chain.  It is assumed that the reader has a basic understanding of the Linux operating system and the command line.  The guide is written for Ubuntu 20.04 LTS, but should be applicable to other Linux distributions.

You can use the knowledge gained from this guide to run a node on any chain that uses the Cosmos SDK.  The only difference is that the chain ID, the binary name, and the genesis file will be different and the Cosmos SDK version may be different.


# Table Of Contents:

- Suggestions
- Requirements
- Prerequisites
- Setup the environment
  - Create the starnamed user and group
- Sync:
  - Sync'ing Using `statesync`
  - Sync with the data backup recovery
  - Sync'ing with reprocessing the chain
    - Cosmovisor
    - Setup
    - Starnamed start
- Installing and Configuring Cosmovisor
- Troubleshooting Common Issues
- Best Practices for Running a Production Node


# Suggestions

Before you begin, please consider the following:

- Read the [Cosmos SDK documentation](https://docs.cosmos.network/main/intro/overview) to understand the Cosmos SDK and the Cosmos Hub.
- The recommended architecture for a Cosmos SDK node is to run the node with the sentry node architecture.  The sentry node architecture is described in the [Sentinel Architecture](https://forum.cosmos.network/t/sentry-node-architecture-overview/454) document.

# Requirements

minimum requirements:
- 1/2 CPU
- 1 GB RAM
- Disk space:
  - if you are syncing from block 1, you will need 240 GB of disk space
  - if you are syncing using `statesync`, you will need 50 GB of disk space
  - if you are syncing using the data backup recovery, you will need the same amount of disk space as the data backup + 100 GB. That can change based in the compression ratio of the data backup.

recommended requirements:
- 3/4 CPU
- 2 GB RAM

**NOTE**: If you want to reprocess the chain is recommended to have at 4 CPU and a fast SSD disk. That will speed up the process. Reduce the log level from `info` to `error` will also speed up the process. 
I got 400tx per second with 4 CPU and a fast SSD disk.  


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

if you are syncing from block one you should install the [cosmovisor binary](https://docs.cosmos.network/main/tooling/cosmovisor) .The following commands will install the cosmovisor binary.  If you are not syncing from block one, you can skip this step. The cosmosvisor is a tool that will help you to upgrade the starnamed binary without downtime.


# Setup the environment

## Create the starnamed user and group

```bash
sudo adduser --system --group --home /home/starnamed starnamed
sudo mkdir -p /home/starnamed/.starnamed
sudo chown -R starnamed:starnamed /home/starnamed
```

# Sync

You have three options to sync the chain: using `statesync`, using the data backup recovery, or reprocessing the chain.
The `statesync` is the fastest option, this will connect to other node and download the state of the chain, but you will need to trust the other node. The data backup recovery is the second fastest option, this will download the data backup and restore the chain state. The reprocessing the chain is the slowest option, this will reprocess the chain from block one.


## TLTR:

This is a simple overview of the sync process. For more details, please read the following sections.

### Sync'ing Using `statesync` or the data backup recovery
The `statesync` and the data backup recovery are the recommended options. Has the same approch:
- download or build the latest version of the starname binary;
- init the node with the ``starnamed init`` command;
- download the genesis file;
- if you are using the data backup recovery, download the data backup;
- configure the node;
  - ${HOME}/.starnamed/config/config.toml
  - ${HOME}/.starnamed/config/app.toml
- create a systemd service;
- start the node.

### Sync'ing with reprocessing the chain
The reprocessing the chain is the slowest option, this will reprocess the chain from block one. The reprocessing the chain has two options, manually or using cosmovisor. The cosmovisor is a tool that will help you to upgrade the starnamed binary without downtime. The cosmovisor is the recommended option. Has the same approch:

- install the cosmovisor binary;
- download or build the latest version of the starname binary;
- init the node with the ``starnamed init`` command;
- download the genesis file;
- configure the node;
  - ${HOME}/.starnamed/config/config.toml
  - ${HOME}/.starnamed/config/app.toml
- create a systemd service;
- start the node.

## Sync'ing Using `statesync` or the data backup recovery

