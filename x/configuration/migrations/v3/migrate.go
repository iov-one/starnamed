package v3

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/configuration/migrations/utils"
)

// MigrateStore performs in-place store migrations from version 2 to version 3
// This adds the auction completion fees
func MigrateStore(ctx sdk.Context, storeKey sdk.StoreKey, cdc codec.BinaryCodec) error {
	// The v2 -> v3 change is also a simple addition of fees
	return utils.MigrateNewFees(ctx, storeKey, cdc)
}