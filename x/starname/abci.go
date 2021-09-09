package starname

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

// BeginBlocker stores the fees payed for the previous block,
// just before they are distributed by the cosmos-sdk/x/distribution's BeginBlocker
func BeginBlocker(ctx sdk.Context, keeper Keeper) {
	feeCollector := keeper.AuthKeeper.GetModuleAddress(authtypes.FeeCollectorName)
	feeCollected := keeper.SupplyKeeper.GetAllBalances(ctx, feeCollector)

	keeper.StoreBlockFees(ctx, feeCollected)
}
