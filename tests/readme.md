# Tests Folder - e2e

This folder contains the end-to-end tests for the project.
Will be used (Strange Love - interchain test framework)[https://github.com/strangelove-ventures/interchaintest] to run the tests.


## Test Cases:

Test Group |Test Name | Test File | Description | Status
--- | --- | --- | --- | ---
Cosmos | Chain Upgrade v0.11.7 -> v0.12 | [upgrade_test.go](./e2e/upgrade_test.go) | Test the upgrade of the chain from v0.11.7 to v0.12 without ibc transactions | Done
Cosmos | Chain Upgrade v0.11.7 -> v0.12 with ibc transactions | [upgrade_with_ibc_test.go](./e2e/upgrade_with_ibc_test.go) | Test the upgrade of the chain from v0.11.7 to v0.12 with ibc transactions | Done
Starname | Starname Module test | [starname_module_test.go](./e2e/starname_module_test.go) | Test the starname module - Domain, Account, Escrow| WIP




