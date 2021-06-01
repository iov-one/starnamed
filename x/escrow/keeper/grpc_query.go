package keeper

import (
	"context"
	"encoding/hex"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/iov-one/starnamed/x/escrow/types"
)

var _ types.QueryServer = Keeper{}

func (k Keeper) Escrow(c context.Context, request *types.QueryEscrowRequest) (*types.QueryEscrowResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	id, err := hex.DecodeString(request.Id)
	if err != nil {
		return nil, status.Errorf(codes.InvalidArgument, "invalid escrow id %s", request.Id)
	}
	escrow, found := k.GetEscrow(ctx, id)

	if !found {
		return nil, sdkerrors.Wrap(types.ErrEscrowNotFound, string(id))
	}

	return &types.QueryEscrowResponse{Escrow: &escrow}, nil
}
