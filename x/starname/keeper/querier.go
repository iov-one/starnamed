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
	if err := keeper.DomainStore(ctx).Read(filter.PrimaryKey(), domain); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", name)
	}

	return &types.QueryDomainResponse{Domain: domain}, nil
}

func (q grpcQuerier) DomainAccounts(c context.Context, req *types.QueryDomainAccountsRequest) (*types.QueryDomainAccountsResponse, error) {
	if req.Domain == "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDomainName, "'%s'", req.Domain)
	}
	// TODO: pagination
	return queryDomainAccounts(sdk.UnwrapSDKContext(c), req.Domain, q.keeper)
}

func queryDomainAccounts(ctx sdk.Context, domain string, keeper *Keeper) (*types.QueryDomainAccountsResponse, error) {
	// TODO: pagination
	cursor, err := keeper.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte(domain)).Do()
	if err != nil {
		return nil, sdkerrors.Wrapf(types.ErrInvalidDomainName, "'%s'", domain)
	}

	accounts := make([]*types.Account, 0)
	for ; cursor.Valid(); cursor.Next() {
		account := new(types.Account)
		if err := cursor.Read(account); err != nil {
			return nil, sdkerrors.Wrapf(types.ErrInvalidDomainName, "'%s'", domain)
		}
		accounts = append(accounts, account)
	}

	return &types.QueryDomainAccountsResponse{Accounts: accounts}, nil
}
