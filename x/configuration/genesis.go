package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/configuration/types"
)


// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetConfig(ctx, data.Config)
	k.SetFees(ctx, &data.Fees)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	return &types.GenesisState{
		Config: k.GetConfiguration(ctx),
		Fees:   *k.GetFees(ctx),
	}
}
