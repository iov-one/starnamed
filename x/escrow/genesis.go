package escrow

import (
	"github.com/iov-one/starnamed/x/escrow/keeper"

	"github.com/iov-one/starnamed/x/escrow/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetLastBlockTime(ctx, data.LastBlockTime)
	k.ImportNextID(ctx, data.NextEscrowId)

	for _, escrow := range data.GetEscrows() {
		k.SaveEscrow(ctx, escrow)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var escrows []types.Escrow
	k.IterateEscrows(
		ctx,
		func(e types.Escrow) (stop bool) {
			escrows = append(escrows, e)
			return false
		},
	)

	lastBlockTime := k.GetLastBlockTime(ctx)
	nextID := k.GetNextIDForExport(ctx)

	return types.NewGenesisState(escrows, lastBlockTime, nextID)
}
