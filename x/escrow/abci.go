package escrow

import (
	"github.com/iov-one/starnamed/x/escrow/keeper"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// BeginBlocker handles block beginning logic for escrows
func BeginBlocker(ctx sdk.Context, k keeper.Keeper) {
	ctx = ctx.WithLogger(ctx.Logger().With("handler", "beginBlock").With("module", "starname/x/escrow"))

	// Automatically refund all expired escrows
	currentDate := uint64(ctx.BlockTime().Unix())
	k.MarkExpiredEscrows(ctx, currentDate)

	k.SetLastBlockTime(ctx, currentDate)
}
