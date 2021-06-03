package keeper

import (
	"context"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/escrow/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Escrow(c context.Context, request *types.QueryEscrowRequest) (*types.QueryEscrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	escrow, found := k.GetEscrow(ctx, request.Id)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrEscrowNotFound, request.Id)
	}

	return &types.QueryEscrowResponse{Escrow: &escrow}, nil
}
