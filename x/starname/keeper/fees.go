package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/x/starname/types"
)

// CollectProductFee takes the product fee from the payer and sends it to the distribution module for validators and delegators
func (k Keeper) CollectProductFee(ctx sdk.Context, msg types.MsgWithFeePayer, withs ...interface{}) error {
	feeConf := k.ConfigurationKeeper.GetFees(ctx)
	feeCtrl := NewFeeController(ctx, feeConf)
	if len(withs) > 0 {
		for _, with := range withs {
			switch with.(type) {
			case *types.Domain:
				feeCtrl.WithDomain(with.(*types.Domain))
			case func(sdk.Context) crud.Store: // can't pass in k.AccountStore(ctx) since its a storeWrapper
				accounts := k.AccountStore(ctx)
				feeCtrl.WithAccounts(&accounts)
			default:
				panic(fmt.Sprintf("unexpected type %T", with))
			}
		}
	}
	fee := feeCtrl.GetFee(msg)
	return k.SupplyKeeper.SendCoinsFromAccountToModule(ctx, msg.FeePayer(), authtypes.FeeCollectorName, sdk.NewCoins(fee))
}
