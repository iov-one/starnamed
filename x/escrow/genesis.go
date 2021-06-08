package escrow

import (
	"github.com/iov-one/starnamed/x/escrow/keeper"

	"github.com/iov-one/starnamed/x/escrow/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	k.SetLastBlockTime(ctx, data.LastBlockTime)

	for _, escrow := range data.GetEscrows() {
		//TODO: check other things, in addition to the validation done in ValidateGenesis?
		k.SaveEscrow(ctx, escrow)
	}
}

// ExportGenesis outputs the genesis state
func ExportGenesis(ctx sdk.Context, k keeper.Keeper) *types.GenesisState {
	var escrows []types.Escrow
	k.IterateEscrows(
		ctx,
		func(_ tmbytes.HexBytes, e types.Escrow) (stop bool) {
			escrows = append(escrows, e)
			return false
		},
	)

	lastBlockTime := k.GetLastBlockTime(ctx)
	nextID := k.GetNextIDForExport(ctx)

	return types.NewGenesisState(escrows, lastBlockTime, nextID)
}

//PrepForZeroHeightGenesis TODO: figure out is this is actually needed or if it is legacy (and if it is, does anything needs to be done in replacement)
func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	// TODO: update previous block time
	//TODO: check how to do this init
	k.SetLastBlockTime(ctx, uint64(ctx.BlockTime().Unix()))
}
