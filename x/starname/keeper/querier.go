package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/starname/types"
)

var _ types.QueryServer = &grpcQuerier{}

type grpcQuerier struct {
	keeper *Keeper
}

// NewQuerier provides a gRPC querier
// TODO: this needs proper tests and doc
func NewQuerier(keeper *Keeper) grpcQuerier {
	return grpcQuerier{keeper: keeper}
}

func (q grpcQuerier) Domain(c context.Context, req *types.QueryDomainRequest) (*types.QueryDomainResponse, error) {
	if req.Name == "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDomainName, "'%s'", req.Name)
	}
	return queryDomain(sdk.UnwrapSDKContext(c), req.Name, q.keeper)
}

func queryDomain(ctx sdk.Context, name string, keeper *Keeper) (*types.QueryDomainResponse, error) {
	domain := new(types.Domain)
	filter := &types.Domain{Name: name}
	ok := keeper.DomainStore(ctx).Read(filter.PrimaryKey(), domain)
	if !ok {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", name)
	}

	return &types.QueryDomainResponse{Domain: domain}, nil
}
