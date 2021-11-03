package starname

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// EndBlocker refresh the fees sliding sum in order to make yield queries faster
func EndBlocker(ctx sdk.Context, k Keeper) {
	// Refresh the value of the fees sum
	k.RefreshBlockSumCache(ctx)
}
