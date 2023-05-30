package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/escrow/keeper"
)

type Migrator = struct {
	RegistrationErrMSG string
	MigrationHandler   func(ctx sdk.Context) error
}

// Migrate the store from version 1 to version 2.
// This migrations its only to enforce the new escrow consensus version
// Because of the change of logic in the escrow module, we need to enforce
func migrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	return nil
}

func NewMigrator(keeper keeper.Keeper) Migrator {
	return Migrator{
		RegistrationErrMSG: "Failed to register migration handler for escrow module, from version 1 to 2",
		MigrationHandler:   func(ctx sdk.Context) error { return migrateStore(ctx, keeper.GetStoreKey(), keeper.GetCdc()) },
	}
}
