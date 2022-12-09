package burner

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/burner/types"
)

//TODO: we could add a test for this function

//EndBlocker burns all the coins owned by the burner module
func EndBlocker(ctx sdk.Context, supplyKeeper types.SupplyKeeper, accountKeeper types.AccountKeeper) {
	moduleAcc := accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	if balance := supplyKeeper.GetAllBalances(ctx, moduleAcc.GetAddress()); !balance.IsZero() {
		if err := supplyKeeper.BurnCoins(ctx, types.ModuleName, balance); err != nil {
			panic(fmt.Sprintf("Error while burning tokens of the burner module account: %s", err.Error()))
		}
	}
}
