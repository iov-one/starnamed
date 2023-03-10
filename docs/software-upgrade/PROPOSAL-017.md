# Software Upgrade Proposal 17 At Block :12180387 #

This version contains significant changes to the chain it self.
- Rebase from wasmd project
- Update to cosmos v0.45.9
- Update the IBC v3
- Starname modules:
  - Add the escrow module
  - Add the burner module

```sh
$ starnamed query gov proposal 17 --chain-id iov-mainnet-ibc --output json | jq .content.plan
{
  "name": "starname-version-11",
  "time": "2022-12-15T12:00:00Z",
  "height": "0",
  "info": "https://github.com/iov-one/starnamed/blob/v0.11.6/upgrades/v011/v011_binaries.json",
  "upgraded_client_state": null
}
```

Note that the intentional panic will occur at block 12180387 at approximately 12am UTC on December 15, 2022.

**Validators and node operators must upgrade their `starnamed` binary to `v0.11.6` and restart it.**

## Upgrading if you followed the "standard" procedure for lighting-up your starnamed node ##
