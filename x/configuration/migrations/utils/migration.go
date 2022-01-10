package utils

import (
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/configuration/types"
)

// MigrateNewConfigurationAttributes allows to add newly defined configuration attributes to the store state, using default values
// from genesis state
func MigrateNewConfigurationAttributes(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	confBytes := store.Get([]byte(types.ConfigKey))
	if confBytes == nil {
		return fmt.Errorf("no configuration available")
	}

	defaultState := types.DefaultGenesisState()

	// Get the default state for the configuration
	config := defaultState.Config
	// Overwrite with the values present in the chain for old fields
	cdc.MustUnmarshal(confBytes, &config)

	// Write back the configuration
	store.Set([]byte(types.ConfigKey), cdc.MustMarshal(&config))
	return nil
}

// MigrateNewFees allows to add newly defined fees to the store state, using default values
// from genesis state
func MigrateNewFees(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	store := ctx.KVStore(storeKey)

	feesBytes := store.Get([]byte(types.FeeKey))
	if feesBytes == nil {
		return fmt.Errorf("no fees available")
	}

	defaultState := types.DefaultGenesisState()

	// Get the default state for the fees
	fees := defaultState.Fees
	// Overwrite with the values present in the chain for old fields
	cdc.MustUnmarshal(feesBytes, &fees)

	// Write back the fees
	store.Set([]byte(types.FeeKey), cdc.MustMarshal(&fees))
	return nil
}
