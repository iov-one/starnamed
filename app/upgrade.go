package app

import (
	"encoding/base64"
	"fmt"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/crypto/keys/multisig"
	"github.com/cosmos/cosmos-sdk/crypto/keys/secp256k1"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authkeeper "github.com/cosmos/cosmos-sdk/x/auth/keeper"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramskeeper "github.com/cosmos/cosmos-sdk/x/params/keeper"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	upgradetypes "github.com/cosmos/cosmos-sdk/x/upgrade/types"
)

type upgradeData struct {
	name    string
	handler upgradetypes.UpgradeHandler
}

func (app *WasmApp) RegisterUpgradeHandlers() {
	upgrades := []upgradeData{GetIOVMainnetIBC2UpgradeHandler()}

	for _, upgrade := range upgrades {
		app.upgradeKeeper.SetUpgradeHandler(upgrade.name, upgrade.handler)
	}
}

func GetIOVMainnetIBC2UpgradeHandler() upgradeData {
	const planName = "iov-mainnet-ibc-2"
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

	handler := func(ctx sdk.Context, plan upgradetypes.Plan) {

		if !plan.ShouldExecute(ctx) || plan.Name != planName {
			panic(fmt.Errorf("the %s upgrade handler has been called when it should not have been", planName))
		}

		encodingConfig := MakeEncodingConfig()
		appCodec, legacyAmino := encodingConfig.Marshaler, encodingConfig.Amino

		var paramKeeper = paramskeeper.NewKeeper(
			appCodec,
			legacyAmino,
			sdk.NewKVStoreKey(paramstypes.StoreKey),
			sdk.NewTransientStoreKey(paramstypes.TStoreKey),
		)

		var keeper = authkeeper.NewAccountKeeper(
			appCodec,
			sdk.NewKVStoreKey(authtypes.StoreKey),
			paramKeeper.Subspace(authtypes.ModuleName),
			authtypes.ProtoBaseAccount,
			maccPerms,
		)

		for _, accountData := range multisigAccounts {
			address, err := sdk.AccAddressFromBech32(accountData.address)
			if err != nil {
				panic(sdkerrors.Wrap(err, "error while decoding a multisig account address"))
			}

			// Get the account
			account := keeper.GetAccount(ctx, address)

			// Ensure that the account is a multisig account
			multisigPubkey, match := account.GetPubKey().(*multisig.LegacyAminoPubKey)
			if !match {
				panic(fmt.Errorf("the accounts %s is not a multisig account", accountData.address))
			}

			// Populate the pubkeys array
			for n, pubkeyData := range accountData.pubkeys {

				pubkeyBytes, err := base64.StdEncoding.DecodeString(pubkeyData)
				if err != nil {
					panic(fmt.Errorf("error while decdoding the pubkey number %n of account %s", n, accountData.address))
				}
				pubkey := &secp256k1.PubKey{Key: pubkeyBytes}

				wrappedPubkey, err := cdctypes.NewAnyWithValue(pubkey)
				if err != nil {
					panic(fmt.Errorf("error while packing the pubkey number %n of account %s", n, accountData.address))
				}
				multisigPubkey.PubKeys = append(multisigPubkey.PubKeys, wrappedPubkey)
			}

			// Write back the multi signature public key to the account
			if err := account.SetPubKey(multisigPubkey); err != nil {
				panic(fmt.Errorf("error while writing the multisig pubkey back to the account %s", accountData.address))
			}

			// Write back the account to the store
			keeper.SetAccount(ctx, account)
		}
	}
	return upgradeData{name: planName, handler: handler}
}
