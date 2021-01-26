package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/iov-one/starnamed/x/starname/types"
)

// CollectProductFee takes the product fee from the payer and sends it to the distribution module for validators and delegators
func (k Keeper) CollectProductFee(ctx sdk.Context, msg types.MsgWithFeePayer, domain ...*types.Domain) error {
	feeConf := k.ConfigurationKeeper.GetFees(ctx)
	feeCtrl := NewFeeController(ctx, feeConf)
	if len(domain) > 0 {
		feeCtrl.WithDomain(domain[0])
	}
	fee := feeCtrl.GetFee(msg)
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), authtypes.FeeCollectorName, sdk.NewCoins(fee))
}
