package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/iov-one/starnamed/x/starname/types"
)

// CollectFees collects the fees of a msg and sends them
// to the distribution module to validators and stakers
func (k Keeper) CollectFees(ctx sdk.Context, msg types.MsgWithFeePayer, fee sdk.Coin) error {
	// transfer fee to distribution
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), authtypes.FeeCollectorName, sdk.NewCoins(fee))
}
