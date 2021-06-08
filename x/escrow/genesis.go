package escrow

import (
	"encoding/hex"

	"github.com/iov-one/starnamed/x/escrow/keeper"

	"github.com/iov-one/starnamed/x/escrow/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

//TODO: add escrow.nextID to the genesis state

// InitGenesis stores the genesis state
func InitGenesis(ctx sdk.Context, k keeper.Keeper, data types.GenesisState) {
	if err := types.ValidateGenesis(data); err != nil {
		panic(err.Error())
	}

	//TODO: check if we should use initial_date instead
	k.SetLastBlockTime(ctx, data.LastBlockTime)

	//var incomingSupplies sdk.Coins
	//var outgoingSupplies sdk.Coins
	for _, escrow := range data.GetEscrows() {
		_, err := hex.DecodeString(escrow.Id)
		if err != nil {
			panic(err.Error())
		}

		//TODO: manage this
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

	return types.NewGenesisState(escrows, lastBlockTime)
}

func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	// TODO: update previous block time
	// TODO: check what we need to
	//TODO: check how to do this init
	k.SetLastBlockTime(ctx, uint64(ctx.BlockTime().Unix()))
}
