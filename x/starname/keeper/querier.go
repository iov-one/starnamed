package keeper

import (
	"context"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/types"
)

var _ types.QueryServer = &grpcQuerier{}

type grpcQuerier struct {
	keeper *Keeper
}

const (
	defaultStart uint64 = 0
	defaultLimit uint64 = 100
)

func getPagination(pageRequest *query.PageRequest) (uint64, uint64, bool, error) {
	start := defaultStart
	end := start + defaultLimit
	count := false
	if pageRequest != nil {
		if pageRequest.Key != nil {
			return 0, 0, false, sdkerrors.Wrapf(types.ErrInvalidRequest, "pagination by key is not implemented")
		}
		start = pageRequest.GetOffset()
		limit := pageRequest.GetLimit()
		if limit == 0 {
			limit = defaultLimit
		}
		end = start + limit
		count = pageRequest.GetCountTotal()
	}
	return start, end, count, nil
}

func getPageResponse(count bool, statement crud.FinalizedIndexStatement) (*query.PageResponse, error) {
	if !count {
		return nil, nil
	}
	cursor, err := statement.Do()
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed to compute total count")
	}
	page := new(query.PageResponse)
	for ; cursor.Valid(); cursor.Next() {
		page.Total++
	}
	return page, nil
}

func bech32FromBytes(bytes []byte) string {
	return sdk.MustBech32ifyAddressBytes(sdk.GetConfig().GetBech32AccountAddrPrefix(), bytes)
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
	start, end, count, err := getPagination(req.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed pagination")
	}
	return queryDomainAccounts(sdk.UnwrapSDKContext(c), q.keeper, req.Domain, start, end, count)
}

func queryDomainAccounts(ctx sdk.Context, keeper *Keeper, domain string, start, end uint64, count bool) (*types.QueryDomainAccountsResponse, error) {
	query := func() crud.FinalizedIndexStatement {
		return keeper.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte(domain))
	}
	cursor, err := query().WithRange().Start(start).End(end).Do()
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
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", domain)
	}
	return &types.QueryDomainAccountsResponse{Accounts: accounts, Page: page}, nil
}

func (q grpcQuerier) Starname(c context.Context, req *types.QueryStarnameRequest) (*types.QueryStarnameResponse, error) {
	if req.Starname == "" || !strings.Contains(req.Starname, types.StarnameSeparator) {
		return nil, sdkerrors.Wrapf(types.ErrInvalidAccountName, "'%s'", req.Starname)
	}
	return queryStarname(sdk.UnwrapSDKContext(c), q.keeper, req.Starname)
}

func queryStarname(ctx sdk.Context, keeper *Keeper, starname string) (*types.QueryStarnameResponse, error) {
	parts := strings.Split(starname, types.StarnameSeparator)
	filter := types.Account{Domain: parts[1], Name: utils.StrPtr(parts[0])}
	account := new(types.Account)
	if err := keeper.AccountStore(ctx).Read(filter.PrimaryKey(), account); err != nil {
		return nil, sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "not found: %s", starname)
	}
	return &types.QueryStarnameResponse{Account: account}, nil
}

func (q grpcQuerier) OwnerAccounts(c context.Context, req *types.QueryOwnerAccountsRequest) (*types.QueryOwnerAccountsResponse, error) {
	address, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' isn't a vaild address", req.Owner)
	}
	start, end, count, err := getPagination(req.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed pagination")
	}
	return queryOwnerAccounts(sdk.UnwrapSDKContext(c), q.keeper, address, start, end, count)
}

func queryOwnerAccounts(ctx sdk.Context, keeper *Keeper, owner sdk.AccAddress, start, end uint64, count bool) (*types.QueryOwnerAccountsResponse, error) {
	query := func() crud.FinalizedIndexStatement {
		return keeper.AccountStore(ctx).Query().Where().Index(types.AccountAdminIndex).Equals(owner)
	}
	cursor, err := query().WithRange().Start(start).End(end).Do()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", bech32FromBytes(owner))
	}
	accounts := make([]*types.Account, 0)
	for ; cursor.Valid(); cursor.Next() {
		account := new(types.Account)
		if err := cursor.Read(account); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to read")
		}
		accounts = append(accounts, account)
	}
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", bech32FromBytes(owner))
	}
	return &types.QueryOwnerAccountsResponse{Accounts: accounts, Page: page}, nil
}
