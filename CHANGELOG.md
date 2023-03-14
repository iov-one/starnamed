Versions:

- v0.10.12
- v0.10.13
- v0.10.14
- v0.10.15
- v0.10.16
- v0.10.17
- v0.10.18
- v0.11.0
- v0.11.1
- v0.11.2
- v0.11.3
- v0.11.4
- v0.11.5
- v0.11.6
- v0.11.7

# Changelog

## [Unreleased](https://github.com/iov-one/starnamed/tree/HEAD)
[Full Changelog](v0.11.6...main)


## [v0.11.6](https://github.com/iov-one/starnamed/releases/tag/v0.11.6) **Stable starnamed v0.11.X**
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.11.5...v0.11.6)
This verison contains significant changes to the chain it self. And is a mandatory upgrade. for more information please see the [PROPOSAL-017-doc](./docs/software-upgrade/PROPOSAL-017.md)
* Bump to cosmos-sdk v45.9
* Bump to wasmd v0.29.2
* Add the upgrde handler at app/upgrade.go "starname-version-11"



## [v0.11.x -> v0.11.5](https://github.com/iov-one/starnamed/releases/tag/v0.11.5) 
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.11.0...v0.11.5)
* Bump to cosmos-sdk v0.44.5
* Module add: escrow
* Module add: off-chain
* Module add: burner


## [v0.11.0](https://github.com/iov-one/starnamed/releases/tag/v0.11.0) 
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.18...v0.11.0)
* Base64 encoded return data on wasm raw query REST endpoint


## [v0.10.18](https://github.com/iov-one/starnamed/releases/tag/v0.10.18) **Stable starnamed v0.10.13**
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.17...v0.10.18)
* Fix a bug where GetBlockFees would always return an empty fee for the current height


## [v0.10.17](https://github.com/iov-one/starnamed/releases/tag/v0.10.17)
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.16...v0.10.17)
* Compatible with v0.10.13 through v0.10.16


## [v0.10.16](https://github.com/iov-one/starnamed/releases/tag/v0.10.16)
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.15...v0.10.16)
* Fixing the consensus braking bug 


## [v0.10.15](https://github.com/iov-one/starnamed/releases/tag/v0.10.15)
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.14...v0.10.15)
* Change `QueryYieldResponse.apy` to `QueryYieldResponse.yield`
* Pick-up the latest tweaks to stargatenet


## [v0.10.14](https://github.com/iov-one/starnamed/releases/tag/v0.10.14)
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.13...v0.10.14)
* Add real-time yield estimate endpoint `/starname/v1beta1/yield`


## [v0.10.13](https://github.com/iov-one/starnamed/releases/tag/v0.10.13)  **Stable starnamed v0.10.13**
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.12...v0.10.13)
* Software upgrade - The public keys for multisig accounts were not migrated from the exported state of `iov-mainnet-2` into `iov-mainnet-ibc`. In order for multisig accounts to work on `iov-mainnet-ibc` the multisig public keys need to be injected into the app state. https://big-dipper.iov-mainnet-ibc.iov.one/proposals/9 aims to do that.
* More info [software-upgrade_docs](./docs/software-upgrade/PROPOSAL-009.md)
* Add the upgrade handler at app/upgrade.go "fix-cosmos-sdk-migrate-bug"



## [v0.10.12](https://github.com/iov-one/starnamed/releases/tag/v0.10.12)
[Full Changelog](https://github.com/iov-one/starnamed/compare/v0.10.11...v0.10.12)
* Add libwasmvm.so to the gitian build
* Reintroduce `-gcflags` and `-asmflags` to the gitian build
They were removed in https://github.com/iov-one/starnamed/pull/59.

