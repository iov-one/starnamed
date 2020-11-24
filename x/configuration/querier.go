package configuration

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/configuration/types"
)

type grpcQuerier struct {
	keeper *Keeper
}

// NewQuerier provides a gRPC querier
// TODO: this needs proper tests and doc
func NewQuerier(keeper *Keeper) grpcQuerier {
	return grpcQuerier{keeper: keeper}
}

func (q grpcQuerier) Config(c context.Context, req *types.QueryConfigRequest) (*types.QueryConfigResponse, error) {
	res := queryConfig(sdk.UnwrapSDKContext(c), *q.keeper)
	return &types.QueryConfigResponse{
		Config: res,
	}, nil
}

func (q grpcQuerier) Fees(c context.Context, req *types.QueryFeesRequest) (*types.QueryFeesResponse, error) {
	res := queryFees(sdk.UnwrapSDKContext(c), *q.keeper)
	return &types.QueryFeesResponse{
		Fees: res,
	}, nil
}

func queryConfig(ctx sdk.Context, keeper Keeper) *types.Config {
	config := keeper.GetConfiguration(ctx) // panics on failure
	return &config
}

func queryFees(ctx sdk.Context, keeper Keeper) *types.Fees {
	fees := keeper.GetFees(ctx) // panics on failure
	return fees
}
