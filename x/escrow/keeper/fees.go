package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/escrow/types"
)

// FIXME: this module should not have a dependency on the configuration module

//CollectFees collect the fees for the given message
func (k *Keeper) CollectFees(ctx sdk.Context, msg types.MsgWithFeePayer) error {
	feesConfiguration := k.configurationKeeper.GetFees(ctx)

	defaultFee := feesConfiguration.FeeDefault

	specificFee := getFee(feesConfiguration, msg)

	if specificFee.LT(defaultFee) {
		specificFee = defaultFee
	}

	finalFeeAmount := specificFee.Quo(feesConfiguration.FeeCoinPrice).TruncateInt()
	finalFee := sdk.NewCoins(sdk.NewCoin(feesConfiguration.FeeCoinDenom, finalFeeAmount))
	return k.bankKeeper.SendCoinsFromAccountToModule(
		ctx,
		msg.GetFeePayer(),
		authtypes.FeeCollectorName,
		finalFee,
	)
}

func getFee(feesConfig *configuration.Fees, msg sdk.Msg) sdk.Dec {
	switch msg.(type) {
	case *types.MsgCreateEscrow:
		// TODO: delegate this
		return feesConfig.FeeDefault
	case *types.MsgUpdateEscrow:
		return feesConfig.UpdateEscrow
	case *types.MsgTransferToEscrow:
		return feesConfig.TransferToEscrow
	case *types.MsgRefundEscrow:
		return feesConfig.RefundEscrow
	default:
		// TODO: should this be an error ?
		return feesConfig.FeeDefault
	}
}
