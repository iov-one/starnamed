package keeper

import (
	"github.com/iov-one/starnamed/x/escrow/types"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewQuerier creates a new escrow Querier instance
func NewQuerier(k Keeper, legacyQuerierCdc *codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		switch path[0] {
		case types.QueryEscrow:
			return queryEscrow(ctx, req, k, legacyQuerierCdc)
		default:
			return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "unknown %s query path: %s", types.ModuleName, path[0])
		}
	}
}

func queryEscrow(ctx sdk.Context, req abci.RequestQuery, k Keeper, legacyQuerierCdc *codec.LegacyAmino) ([]byte, error) {
	var params types.QueryEscrowParams
	if err := legacyQuerierCdc.UnmarshalJSON(req.Data, &params); err != nil {
		return nil, sdkerrors.Wrap(sdkerrors.ErrJSONUnmarshal, err.Error())
	}

	escrow, found := k.GetEscrow(ctx, params.Id)
	if !found {
		return nil, sdkerrors.Wrap(types.ErrEscrowNotFound, params.Id)
	}

	bz, err := legacyQuerierCdc.MarshalJSON(&escrow)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "Cannot marshall the queried escrow")
	}
	return bz, nil
}
