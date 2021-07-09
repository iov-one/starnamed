package keeper

import (
	"context"
	"encoding/hex"
	"strings"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"

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

func (k Keeper) Escrows(c context.Context, request *types.QueryEscrowsRequest) (*types.QueryEscrowsResponse, error) {
	ctx := sdk.UnwrapSDKContext(c)

	var seller sdk.AccAddress
	if len(request.Seller) != 0 {
		var err error
		seller, err = sdk.AccAddressFromBech32(request.Seller)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "Invalid seller address")
		}
	}

	var state types.EscrowState
	var hasState bool
	if len(request.State) != 0 {
		hasState = true
		switch strings.ToLower(request.State) {
		case "open":
			state = types.EscrowState_Open
		case "expired":
			state = types.EscrowState_Expired
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The state is invalid, can be open or expired")
		}
	}

	var objectKey []byte
	if len(request.ObjectKey) != 0 {
		var err error
		objectKey, err = hex.DecodeString(request.ObjectKey)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "Object key must be an hex-encoded byte array")
		}
	}

	var end uint64
	if request.PaginationLength == 0 {
		end = 0
	} else {
		end = request.PaginationStart + request.PaginationLength
	}

	filter := func(query crud.QueryStatement) crud.ValidQuery {
		getStatement := func(query crud.QueryStatement, previous crud.FinalizedIndexStatement) crud.WhereStatement {
			if previous == nil {
				return query.Where()
			} else {
				return previous.And()
			}
		}

		previousStatement := crud.FinalizedIndexStatement(nil)
		//TODO: maybe optimize by filtering first by object if there is an object filter
		if seller != nil {
			previousStatement = getStatement(query, previousStatement).
				Index(types.SellerIndex).Equals(seller)
		}
		if hasState {
			previousStatement = getStatement(query, previousStatement).
				Index(types.StateIndex).Equals(sdk.Uint64ToBigEndian(uint64(state)))
		}
		if objectKey != nil {
			previousStatement = getStatement(query, previousStatement).
				Index(types.ObjectIndex).Equals(objectKey)
		}

		if previousStatement == nil {
			return query
		} else {
			return previousStatement
		}
	}

	escrows, err := k.QueryEscrowsWithRange(ctx, filter, request.PaginationStart, end)

	if err != nil {
		return nil, err
	}

	return &types.QueryEscrowsResponse{Escrows: escrows}, nil
}
