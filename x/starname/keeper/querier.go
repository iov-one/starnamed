package keeper

import (
	"context"
	"fmt"
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
	defaultStart     uint64 = 0
	defaultLimit     uint64 = 100 // TODO: read this from config.toml
	NumBlocksInAWeek uint64 = 100000
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

// NewQuerier provides a gRPC querier
func NewQuerier(keeper *Keeper) grpcQuerier {
	return grpcQuerier{keeper: keeper}
}

// Domain returns a types.Domain if the domain exists and nil on error
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

// DomainAccounts returns types.Accounts associated with a given domain and nil on error
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
	accounts := make([]*types.Account, 0, end-start)
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

// Starname returns the types.Account associated with a given starname and nil on error
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

// OwnerAccounts returns types.Accounts associated with a given owner and nil on error
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
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", owner.String())
	}
	accounts := make([]*types.Account, 0, end-start)
	for ; cursor.Valid(); cursor.Next() {
		account := new(types.Account)
		if err := cursor.Read(account); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to read")
		}
		accounts = append(accounts, account)
	}
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", owner.String())
	}
	return &types.QueryOwnerAccountsResponse{Accounts: accounts, Page: page}, nil
}

// OwnerDomains returns types.Domains associated with a given owner and nil on error
func (q grpcQuerier) OwnerDomains(c context.Context, req *types.QueryOwnerDomainsRequest) (*types.QueryOwnerDomainsResponse, error) {
	address, err := sdk.AccAddressFromBech32(req.Owner)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' isn't a vaild address", req.Owner)
	}
	start, end, count, err := getPagination(req.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed pagination")
	}
	return queryOwnerDomains(sdk.UnwrapSDKContext(c), q.keeper, address, start, end, count)
}

func queryOwnerDomains(ctx sdk.Context, keeper *Keeper, owner sdk.AccAddress, start, end uint64, count bool) (*types.QueryOwnerDomainsResponse, error) {
	query := func() crud.FinalizedIndexStatement {
		return keeper.DomainStore(ctx).Query().Where().Index(types.DomainAdminIndex).Equals(owner)
	}
	cursor, err := query().WithRange().Start(start).End(end).Do()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", owner.String())
	}
	domains := make([]*types.Domain, 0, end-start)
	for ; cursor.Valid(); cursor.Next() {
		domain := new(types.Domain)
		if err := cursor.Read(domain); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to read")
		}
		domains = append(domains, domain)
	}
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", owner.String())
	}
	return &types.QueryOwnerDomainsResponse{Domains: domains, Page: page}, nil
}

// ResourceAccounts returns types.Accounts associated with a given resource and nil on error
func (q grpcQuerier) ResourceAccounts(c context.Context, req *types.QueryResourceAccountsRequest) (*types.QueryResourceAccountsResponse, error) {
	if req.Uri == "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidResource, "'%s'", req.Uri)
	}
	if req.Resource == "" {
		return nil, sdkerrors.Wrapf(types.ErrInvalidResource, "'%s'", req.Resource)
	}
	start, end, count, err := getPagination(req.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed pagination")
	}
	return queryResourceAccounts(sdk.UnwrapSDKContext(c), q.keeper, req.Uri, req.Resource, start, end, count)
}

func queryResourceAccounts(ctx sdk.Context, keeper *Keeper, uri string, resource string, start, end uint64, count bool) (*types.QueryResourceAccountsResponse, error) {
	key := types.GetResourceKey(uri, resource)
	query := func() crud.FinalizedIndexStatement {
		return keeper.AccountStore(ctx).Query().Where().Index(types.AccountResourcesIndex).Equals(key)
	}
	cursor, err := query().WithRange().Start(start).End(end).Do()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s:%s' caused error", uri, resource)
	}
	accounts := make([]*types.Account, 0, end-start)
	for ; cursor.Valid(); cursor.Next() {
		account := new(types.Account)
		if err := cursor.Read(account); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to read")
		}
		accounts = append(accounts, account)
	}
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s:%s' caused error", uri, resource)
	}
	return &types.QueryResourceAccountsResponse{Accounts: accounts, Page: page}, nil
}

// BrokerAccounts returns types.Accounts associated with a given broker and nil on error
func (q grpcQuerier) BrokerAccounts(c context.Context, req *types.QueryBrokerAccountsRequest) (*types.QueryBrokerAccountsResponse, error) {
	address, err := sdk.AccAddressFromBech32(req.Broker)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' isn't a vaild address", req.Broker)
	}
	start, end, count, err := getPagination(req.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed pagination")
	}
	return queryBrokerAccounts(sdk.UnwrapSDKContext(c), q.keeper, address, start, end, count)
}

func queryBrokerAccounts(ctx sdk.Context, keeper *Keeper, broker sdk.AccAddress, start, end uint64, count bool) (*types.QueryBrokerAccountsResponse, error) {
	query := func() crud.FinalizedIndexStatement {
		return keeper.AccountStore(ctx).Query().Where().Index(types.AccountBrokerIndex).Equals(broker)
	}
	cursor, err := query().WithRange().Start(start).End(end).Do()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", broker.String())
	}
	accounts := make([]*types.Account, 0, end-start)
	for ; cursor.Valid(); cursor.Next() {
		account := new(types.Account)
		if err := cursor.Read(account); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to read")
		}
		accounts = append(accounts, account)
	}
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", broker.String())
	}
	return &types.QueryBrokerAccountsResponse{Accounts: accounts, Page: page}, nil
}

// BrokerDomains returns types.Domains associated with a given broker and nil on error
func (q grpcQuerier) BrokerDomains(c context.Context, req *types.QueryBrokerDomainsRequest) (*types.QueryBrokerDomainsResponse, error) {
	address, err := sdk.AccAddressFromBech32(req.Broker)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' isn't a vaild address", req.Broker)
	}
	start, end, count, err := getPagination(req.Pagination)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "failed pagination")
	}
	return queryBrokerDomains(sdk.UnwrapSDKContext(c), q.keeper, address, start, end, count)
}

func queryBrokerDomains(ctx sdk.Context, keeper *Keeper, broker sdk.AccAddress, start, end uint64, count bool) (*types.QueryBrokerDomainsResponse, error) {
	query := func() crud.FinalizedIndexStatement {
		return keeper.DomainStore(ctx).Query().Where().Index(types.DomainBrokerIndex).Equals(broker)
	}
	cursor, err := query().WithRange().Start(start).End(end).Do()
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", broker.String())
	}
	domains := make([]*types.Domain, 0, end-start)
	for ; cursor.Valid(); cursor.Next() {
		domain := new(types.Domain)
		if err := cursor.Read(domain); err != nil {
			return nil, sdkerrors.Wrap(err, "failed to read")
		}
		domains = append(domains, domain)
	}
	page, err := getPageResponse(count, query())
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "'%s' caused error", broker.String())
	}
	return &types.QueryBrokerDomainsResponse{Domains: domains, Page: page}, nil
}

//Yield return an estimation of the delegators annualized yield based on the last 100k blocks
func (q grpcQuerier) Yield(ctx context.Context, _ *types.QueryYieldRequest) (*types.QueryYieldResponse, error) {
	var response types.QueryYieldResponse

	apy, err := calculateYield(sdk.UnwrapSDKContext(ctx), q.keeper)
	if err != nil {
		return nil, err
	}
	response.Apy = apy
	return &response, err
}

func calculateYield(ctx sdk.Context, keeper *Keeper) (sdk.Dec, error) {

	totalFees, numBlocks := keeper.GetBlockFeesSum(ctx, NumBlocksInAWeek)

	if numBlocks != NumBlocksInAWeek {
		return sdk.ZeroDec(), fmt.Errorf("not enough data to estimate yield: current height %v is smaller than %v",
			ctx.BlockHeight(), NumBlocksInAWeek)
	}

	rewardPool := sdk.NewDecCoinsFromCoins(totalFees...)

	totalDelegatedTokens := keeper.StakingKeeper.GetLastTotalPower(ctx) // in iov

	// Voting power is returned in tokens while fees are in a sub-unit (iov vs uiov)
	multiplier := 1e6
	totalDelegatedTokens = totalDelegatedTokens.Mul(sdk.NewInt(int64(multiplier))) // in uiov

	// Compute yield for numBlocks blocks
	yieldForPeriod := rewardPool.QuoDec(sdk.NewDecFromInt(totalDelegatedTokens))

	var apy sdk.Dec
	if len(yieldForPeriod) == 0 {
		apy = sdk.ZeroDec()
	} else {
		const WeeksPerYear = 52
		// TODO: manage multiple tokens for fees
		apy = yieldForPeriod.MulDec(sdk.NewDec(int64(WeeksPerYear)))[0].Amount
	}

	// Give it a nice percentage format
	return apy.Mul(sdk.NewDec(100)), nil
}
