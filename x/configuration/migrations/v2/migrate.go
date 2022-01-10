package v2

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/configuration/migrations/utils"
)

// MigrateStore performs in-place store migrations from version 1 to version 2
// This adds the escrow parameters to the configuration and the fees corresponding to new transactions
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	if err := utils.MigrateNewConfigurationAttributes(ctx, storeKey, cdc); err != nil {
		return err
	}
	return utils.MigrateNewFees(ctx, storeKey, cdc);
}