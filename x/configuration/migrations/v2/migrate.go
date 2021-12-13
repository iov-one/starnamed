package v2

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/configuration/types"
)

// MigrateStore performs in-place store migrations from version 1 to version 2
// This adds the escrow parameters to the configuration
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	confBytes := store.Get([]byte(types.ConfigKey))
	if confBytes == nil {
		return fmt.Errorf("no configuration available")
	}

	feesBytes := store.Get([]byte(types.FeeKey))
	if feesBytes == nil {
		return fmt.Errorf("no fees available")
	}

	defaultState := types.DefaultGenesisState()

	// Get the default state for the configuration
	config := defaultState.Config
	// Overwrite with the values present in the chain for old fields
	cdc.MustUnmarshal(confBytes, &config)

	// Get the default state for the fees
	fees := defaultState.Fees
	// Overwrite with the values present in the chain for old fields
	cdc.MustUnmarshal(feesBytes, &fees)

	// Write back the configuration and the fees
	store.Set([]byte(types.ConfigKey), cdc.MustMarshal(&config))
	store.Set([]byte(types.FeeKey), cdc.MustMarshal(&fees))

	return nil
}