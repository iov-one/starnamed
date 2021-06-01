package escrow

import (
	"encoding/hex"

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

	//TODO: what is it for ?? : k.SetPreviousBlockTime(ctx, data.PreviousBlockTime)

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

	// TODO: what is it for ?
	/*previousBlockTime, found := k.GetPreviousBlockTime(ctx)
	if !found {
		previousBlockTime = types.DefaultPreviousBlockTime
	}*/

	return types.NewGenesisState(escrows)
}

//TODO: what was this for ?
/*
func PrepForZeroHeightGenesis(ctx sdk.Context, k keeper.Keeper) {
	k.IterateHTLCs(
		ctx,
		func(id tmbytes.HexBytes, h types.HTLC) (stop bool) {
			if h.State == types.Open {
				h.ExpirationHeight = h.ExpirationHeight - uint64(ctx.BlockHeight()) + 1
				k.SetHTLC(ctx, h, id)
			}
			return false
		},
	)
	// TODO: update asset supplies and previous block time
}*/
