package app

import (
	"encoding/base64"
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	storetypes "github.com/cosmos/cosmos-sdk/store/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	authz "github.com/cosmos/cosmos-sdk/x/authz/keeper"
	"github.com/cosmos/cosmos-sdk/x/feegrant"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"

	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"
)

type upgradeData struct {
	name                  string
	handler               upgradetypes.UpgradeHandler
	storeLoaderRegisterer func(*WasmApp, storetypes.UpgradeInfo)
}

func (app *WasmApp) RegisterUpgradeHandlers() {
	upgrades := []upgradeData{
		getIOVMainnetIBC2UpgradeHandler(app),
		getCosmosSDKv44UpgradeHandler(app),
	}

	upgradeInfo, err := app.upgradeKeeper.ReadUpgradeInfoFromDisk()
	if err != nil {
		panic(sdkerrors.Wrap(err, "cannot read upgrade info to register store loaders"))
	}

	for _, upgrade := range upgrades {
		// Register upgrade handler
		app.upgradeKeeper.SetUpgradeHandler(upgrade.name, upgrade.handler)

		// Run the store migrations when needed
		// This is taken from https://docs.cosmos.network/master/migrations/chain-upgrade-guide-044.html#chain-upgrade
		if upgradeInfo.Name == upgrade.name && !app.upgradeKeeper.IsSkipHeight(upgradeInfo.Height) &&
			upgrade.storeLoaderRegisterer != nil {
			upgrade.storeLoaderRegisterer(app, upgradeInfo)
		}
	}

}

func getIOVMainnetIBC2UpgradeHandler(app *WasmApp) upgradeData {
	const planName = "fix-cosmos-sdk-migrate-bug"
	multisigAccounts := []struct {
		address string
		pubkeys []string
	}{
		{"star1p0d75y4vpftsx9z35s93eppkky7kdh220vrk8n", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"A1h5IQfNkKMzMgbKjx6vzYpcfsMLnt9dszPOuc/MOtIB",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
		}},
		{"star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
			"AsV8BTthFjvoa0lWsaYNtJZnenAfl+ys/aTWRTFtgva3",
		}},
		{"star15u4kl3lalt8pm2g4m23erlqhylz76rfh50cuv8", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"Ah1xwTaRh8QCSn1g4r48WdqHDo3LNrK+pb1Y2qPf/csK",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
		}},
		{"star1hjf04872s9rlcdg2wqwvapwttvt3p4gjpp0xmc", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
			"AscithgRiVr9Wsy0U4hY1Y/1arq6DKJzyE9fs1y4UUS/",
		}},
		{"star1elad203jykd8la6wgfnvk43rzajyqpk0wsme9g", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"A96xuKUIkqK/Esfw31c9aC9iIZ+LczZimVio0S9gQ0py",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
		}},
		{"star16tm7scg0c2e04s0exk5rgpmws2wk4xkd84p5md", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"AjMwu9UZTpYnjUP4oAroRreac8hIzbACU54Homzn1hcY",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
		}},
		{"star1m7jkafh4gmds8r0w79y2wu2kvayqvrwt7cy7rf", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
			"Ay8+nUz7XlSA0/bXv37OixZcvoIA2eoY75IDTf0wFW2G",
		}},
		{"star1uyny88het6zaha4pmkwrkdyj9gnqkdfe4uqrwq", []string{
			"A5x4p4VzRwdN+1GLwBWprv/j7W2MY8i4qzAV3mPfCjM4",
			"AsW0xwRR6wn6Hh7xzwA2RqWg4DIPZhHr0ybOfC/rZuM/",
			"A2jF4pI8Pv2bbDd8RvKp+69F9NPPaqiDzSg4GlTwnC4R",
			"Ar9we3pS97HTuADdhs/i5raJoPgjU7TsuLc6Zi76cYUc",
			"A59bIrZzOLKNjWYsYj5dm8lSEhPWwrFttNBkGQv2GWzN",
			"AwOzGduZPxmjUMKASZGKPrUA7Drs9CvfJfXkgR/RSdyu",
		}},
	}

	handler := func(ctx sdk.Context, plan upgradetypes.Plan, vm module.VersionMap) (module.VersionMap, error) {

		if !plan.ShouldExecute(ctx) || plan.Name != planName {
			panic(fmt.Errorf("the %s upgrade handler has been called when it should not have been", planName))
		}

		for _, accountData := range multisigAccounts {
			address, err := sdk.AccAddressFromBech32(accountData.address)
			if err != nil {
				return vm, sdkerrors.Wrap(err, "error while decoding a multisig account address")
			}

			// Get the account
			account := app.accountKeeper.GetAccount(ctx, address)

			// Ensure that the account is a multisig account
			multisigPubkey, match := account.GetPubKey().(*multisig.LegacyAminoPubKey)
			if !match {
				return vm, fmt.Errorf("the accounts %s is not a multisig account", accountData.address)
			}

			// Populate the pubkeys array
			for n, pubkeyData := range accountData.pubkeys {

				pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkeyData)
				if err != nil {
					return vm, fmt.Errorf("error while decdoding the pubkey number %d of account %s", n, accountData.address)
				}
				pubkey := &secp256k1.PubKey{Key: pubkeyBytes}

				wrappedPubkey, err := cdctypes.NewAnyWithValue(pubkey)
				if err != nil {
					return vm, fmt.Errorf("error while packing the pubkey number %d of account %s", n, accountData.address)
				}
				multisigPubkey.PubKeys = append(multisigPubkey.PubKeys, wrappedPubkey)
			}

			// Write back the multi signature public key to the account
			if err := account.SetPubKey(multisigPubkey); err != nil {
				return vm, fmt.Errorf("error while writing the multisig pubkey back to the account %s", accountData.address)
			}

			// Write back the account to the store
			app.accountKeeper.SetAccount(ctx, account)
		}

		return vm, nil
	}
	return upgradeData{name: planName, handler: handler}
}

func getCosmosSDKv44UpgradeHandler(app *WasmApp) upgradeData {
	const planName = "cosmos-sdk-v0.44-upgrade"
	handler := func(ctx sdk.Context, plan upgradetypes.Plan, fromVersionMap module.VersionMap) (module.VersionMap, error) {
		// Overwrite the version map :
		// Set up modules that were already present in previous version (but were not registered as the version map didn't
		// exist prior to v0.44.3). All those modules are at there first registered version (1).
		// If we keep the version map as is (empty) the upgrade handler will use the DefaultGenesis state for those modules
		fromVersionMap = map[string]uint64{
			// Cosmos sdk modules
			"auth":         1,
			"bank":         1,
			"capability":   1,
			"crisis":       1,
			"distribution": 1,
			"evidence":     1,
			"genutil":      1,
			"gov":          1,
			"ibc":          1,
			"mint":         1,
			"params":       1,
			"slashing":     1,
			"staking":      1,
			"transfer":     1,
			"upgrade":      1,
			"vesting":      1,

			// Custom modules
			"burner":        1, // the burner module has no state but it implements AppModule so its better to put it here
			"configuration": 2, // the configuration module will be updated to version 2 (adding the escrow conf)
			"starname":      1,
			"wasm":          1,

			// The escrow is a newly introduced module, as well as the feegrant and authz modules so we do not include them
		}
		return app.mm.RunMigrations(ctx, app.configurator, fromVersionMap)
	}

	// Set the store loader for the 3 new modules : authz, feegrant and escrow
	setStoreLoader := func(app *WasmApp, info storetypes.UpgradeInfo) {
		storeUpgrades := storetypes.StoreUpgrades{
			Added: []string{authz.StoreKey, feegrant.StoreKey, escrowtypes.StoreKey},
		}

		// configure store loader that checks if version == upgradeHeight and applies store upgrades
		app.SetStoreLoader(upgradetypes.UpgradeStoreLoader(info.Height, &storeUpgrades))
	}

	return upgradeData{
		name:                  planName,
		handler:               handler,
		storeLoaderRegisterer: setStoreLoader,
	}
}
