package keeper

import (
	"bytes"
	"errors"
	"sort"
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/iov-one/starnamed/x/starname/types"
)

// domains for populateAccounts()
var domains = []string{"a", "b"}

// account names for populateAccounts()
var names = []string{"", "4", "3", "2", "1"}

// account owners for populateAccounts()
var alice, bob = genTestAddress()
var owners = []sdk.AccAddress{alice, bob}

// resources for populateAccounts()
var resources = []types.Resource{
	{URI: "u0", Resource: "r0"},
	{URI: "u1", Resource: "r1"},
	{URI: "u2", Resource: "r2"},
}

// populateAccounts uses vars domains, names, and owners to create crud objects
func populateAccounts(t *testing.T, keeper Keeper, ctx sdk.Context) (accounts []*types.Account, accountsByOwner map[string][]*types.Account, accountsByDomain map[string][]*types.Account, accountsByResource map[string][]*types.Account) {
	accounts = make([]*types.Account, 0)
	accountsByOwner = make(map[string][]*types.Account)
	accountsByDomain = make(map[string][]*types.Account)
	accountsByResource = make(map[string][]*types.Account)

	n := len(owners)
	m := len(resources)
	for i, domain := range domains {
		domain := domain
		for j, name := range names {
			name := name
			owner := owners[(i+j)%n] // pseudo random owner
			bech32 := owner.String()
			resource := resources[(i+j)%m] // pseudo random resource
			account := types.Account{
				Owner:     owner,
				Domain:    domain,
				Name:      &name,
				Resources: []*types.Resource{&resource},
			}
			if err := keeper.AccountStore(ctx).Create(&account); err != nil {
				t.Fatal(err)
			}
			accounts = append(accounts, &account)
			accountsByOwner[bech32] = append(accountsByOwner[bech32], &account)
			accountsByDomain[domain] = append(accountsByDomain[domain], &account)
			key := string(types.GetResourceKey(resource.URI, resource.Resource))
			accountsByResource[key] = append(accountsByResource[key], &account)
		}
	}

	// sort test vectors on primary key
	sort.Slice(accounts, func(i, j int) bool {
		return bytes.Compare(accounts[i].PrimaryKey(), accounts[j].PrimaryKey()) < 0
	})
	DebugAccounts("accounts", accounts)
	for owner, slice := range accountsByOwner {
		sort.Slice(slice, func(i, j int) bool {
			return bytes.Compare(slice[i].PrimaryKey(), slice[j].PrimaryKey()) < 0
		})
		DebugAccounts(owner, slice)
	}
	for domain, slice := range accountsByDomain {
		sort.Slice(slice, func(i, j int) bool {
			return bytes.Compare(slice[i].PrimaryKey(), slice[j].PrimaryKey()) < 0
		})
		DebugAccounts(domain, slice)
	}
	for key, slice := range accountsByResource {
		sort.Slice(slice, func(i, j int) bool {
			return bytes.Compare(slice[i].PrimaryKey(), slice[j].PrimaryKey()) < 0
		})
		DebugAccounts(key, slice)
	}

	return accounts, accountsByOwner, accountsByDomain, accountsByResource
}

func TestDomain(t *testing.T) {
	name := "test"
	domain := types.Domain{
		Name:       name,
		Admin:      aliceAddr,
		Broker:     bobAddr,
		ValidUntil: 69,
		Type:       "open",
	}
	keeper, ctx, _ := NewTestKeeper(t, false)
	if err := keeper.DomainStore(ctx).Create(&domain); err != nil {
		t.Fatalf("failed to create domain '%s'", name)
	}

	tests := map[string]struct {
		request  types.QueryDomainRequest
		wantErr  error
		validate func(*types.QueryDomainResponse)
	}{
		"query invalid domain": {
			request: types.QueryDomainRequest{
				Name: "",
			},
			wantErr:  sdkerrors.Wrapf(types.ErrInvalidDomainName, "''"),
			validate: nil,
		},
		"query valid domain": {
			request: types.QueryDomainRequest{
				Name: name,
			},
			wantErr: nil,
			validate: func(response *types.QueryDomainResponse) {
				if response.Domain == nil {
					t.Fatal("nil Domain")
				}
				got := response.Domain
				if err := CompareDomains(got, &domain); err != nil {
					DebugDomain(got)
					DebugDomain(&domain)
					t.Fatal(sdkErrors.Wrapf(err, name))
				}
			},
		},
		"query non-existent domain": {
			request: types.QueryDomainRequest{
				Name: name[1:],
			},
			wantErr:  sdkerrors.Wrap(types.ErrDomainDoesNotExist, ""),
			validate: nil,
		},
	}

	for _, test := range tests {
		res, err := NewQuerier(&keeper).Domain(sdk.WrapSDKContext(ctx), &test.request)
		if test.wantErr != nil && !errors.Is(err, test.wantErr) {
			t.Fatalf("wanted err: %s, got: %s", test.wantErr, err)
		}

		if test.validate != nil {
			test.validate(res)
		}
	}
}

func TestDomainAccounts(t *testing.T) {
	keeper, ctx, _ := NewTestKeeper(t, false)
	_, _, accountsByDomain, _ := populateAccounts(t, keeper, ctx)
	tests := map[string]struct {
		request  types.QueryDomainAccountsRequest
		wantErr  func(error) error
		validate func(*types.QueryDomainAccountsResponse)
	}{
		"query invalid domain": {
			request: types.QueryDomainAccountsRequest{
				Domain:     "",
				Pagination: nil,
			},
			wantErr: func(err error) error {
				if !errors.Is(err, sdkerrors.Wrapf(types.ErrInvalidDomainName, "''")) {
					t.Fatal("wrong error")
				}
				return nil
			},
			validate: nil,
		},
		"query valid domain without pagination": {
			request: types.QueryDomainAccountsRequest{
				Domain:     domains[0],
				Pagination: nil,
			},
			wantErr: nil,
			validate: func(response *types.QueryDomainAccountsResponse) {
				if response.Accounts == nil {
					t.Fatalf("wanted non-nil accounts")
				}
				wants := accountsByDomain[domains[0]]
				if len(response.Accounts) != len(wants) {
					t.Fatalf("wanted %d accounts, got %d", len(wants), len(response.Accounts))
				}
				for i, got := range response.Accounts {
					want := wants[i]
					if err := CompareAccounts(got, want); err != nil {
						DebugAccount(got)
						DebugAccount(want)
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
					}
				}
			},
		},
		"query valid domain with pagination": {
			request: types.QueryDomainAccountsRequest{
				Domain: domains[0],
				Pagination: &query.PageRequest{
					Key:        nil,
					Offset:     1,
					Limit:      2,
					CountTotal: true,
				},
			},
			wantErr: nil,
			validate: func(response *types.QueryDomainAccountsResponse) {
				if response.Page == nil {
					t.Fatal("wanted non-nil Page")
				}
				wants := accountsByDomain[domains[0]]
				if response.Page.Total != uint64(len(wants)) {
					t.Fatalf("wanted %d accounts, got %d", len(wants), response.Page.Total)
				}
				if response.Accounts == nil {
					t.Fatalf("wanted non-nil accounts")
				}
				if len(response.Accounts) != 2 {
					t.Fatalf("wanted %d accounts, got %d", 2, len(response.Accounts))
				}
				limited := wants[1:3] // slice to offset and limit
				DebugAccounts("limited", limited)
				for i, got := range response.Accounts {
					want := limited[i]
					if err := CompareAccounts(got, want); err != nil {
						DebugAccount(got)
						DebugAccount(want)
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
					}
				}
			},
		},
	}

	for _, test := range tests {
		res, err := NewQuerier(&keeper).DomainAccounts(sdk.WrapSDKContext(ctx), &test.request)
		if test.wantErr != nil && test.wantErr(err) != nil {
			t.Fatalf("failed err test on: %s", err)
		}
		if test.validate != nil {
			test.validate(res)
		}
	}
}

func TestStarname(t *testing.T) {
	name := "test"
	domain := "domain"
	account := types.Account{
		Domain:     domain,
		Name:       &name,
		Owner:      aliceAddr,
		Broker:     bobAddr,
		ValidUntil: 69,
		Resources: []*types.Resource{
			&types.Resource{URI: "uri0", Resource: "resource0"},
			&types.Resource{URI: "uri1", Resource: "resource1"},
		},
		Certificates: [][]byte{
			[]byte("cert0"),
			[]byte("cert1"),
		},
		MetadataURI: "metadata",
	}
	keeper, ctx, _ := NewTestKeeper(t, false)
	if err := keeper.AccountStore(ctx).Create(&account); err != nil {
		t.Fatalf("failed to create account '%s'", account.GetStarname())
	}

	tests := map[string]struct {
		request  types.QueryStarnameRequest
		wantErr  error
		validate func(*types.QueryStarnameResponse)
	}{
		"query invalid Starname": {
			request: types.QueryStarnameRequest{
				Starname: "",
			},
			wantErr:  sdkerrors.Wrapf(types.ErrInvalidAccountName, "''"),
			validate: nil,
		},
		"query valid starname": {
			request: types.QueryStarnameRequest{
				Starname: account.GetStarname(),
			},
			wantErr: nil,
			validate: func(response *types.QueryStarnameResponse) {
				if response.Account == nil {
					t.Fatal("nil Account")
				}
				got := response.Account
				if err := CompareAccounts(got, &account); err != nil {
					DebugAccount(got)
					DebugAccount(&account)
					t.Fatal(sdkErrors.Wrapf(err, name))
				}
			},
		},
		"query non-existent starname": {
			request: types.QueryStarnameRequest{
				Starname: account.GetStarname()[1:],
			},
			wantErr:  sdkerrors.Wrap(types.ErrAccountDoesNotExist, ""),
			validate: nil,
		},
	}

	for _, test := range tests {
		res, err := NewQuerier(&keeper).Starname(sdk.WrapSDKContext(ctx), &test.request)
		if test.wantErr != nil && !errors.Is(err, test.wantErr) {
			t.Fatalf("wanted err: %s, got: %s", test.wantErr, err)
		}

		if test.validate != nil {
			test.validate(res)
		}
	}
}

func TestOwnerAccounts(t *testing.T) {
	keeper, ctx, _ := NewTestKeeper(t, false)
	_, accountsByOwner, _, _ := populateAccounts(t, keeper, ctx)
	tests := map[string]struct {
		request  types.QueryOwnerAccountsRequest
		wantErr  func(error) error
		validate func(*types.QueryOwnerAccountsResponse)
	}{
		"query invalid owner": {
			request: types.QueryOwnerAccountsRequest{
				Owner:      "bogus",
				Pagination: nil,
			},
			wantErr: func(err error) error {
				if strings.Index(err.Error(), "bogus") == -1 {
					t.Fatal("wrong error")
				}
				return nil
			},
			validate: nil,
		},
		"query valid owner without pagination": {
			request: types.QueryOwnerAccountsRequest{
				Owner:      owners[0].String(),
				Pagination: nil,
			},
			wantErr: nil,
			validate: func(response *types.QueryOwnerAccountsResponse) {
				if response.Accounts == nil {
					t.Fatalf("wanted non-nil accounts")
				}
				wants := accountsByOwner[owners[0].String()]
				if len(response.Accounts) != len(wants) {
					t.Fatalf("wanted %d accounts, got %d", len(wants), len(response.Accounts))
				}
				for i, got := range response.Accounts {
					want := wants[i]
					if err := CompareAccounts(got, want); err != nil {
						DebugAccount(got)
						DebugAccount(want)
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
					}
				}
			},
		},
		"query valid owner with pagination": {
			request: types.QueryOwnerAccountsRequest{
				Owner: owners[0].String(),
				Pagination: &query.PageRequest{
					Key:        nil,
					Offset:     1,
					Limit:      2,
					CountTotal: true,
				},
			},
			wantErr: nil,
			validate: func(response *types.QueryOwnerAccountsResponse) {
				if response.Page == nil {
					t.Fatal("wanted non-nil Page")
				}
				wants := accountsByOwner[owners[0].String()]
				if response.Page.Total != uint64(len(wants)) {
					t.Fatalf("wanted %d accounts, got %d", len(wants), response.Page.Total)
				}
				if response.Accounts == nil {
					t.Fatalf("wanted non-nil accounts")
				}
				if len(response.Accounts) != 2 {
					t.Fatalf("wanted %d accounts, got %d", 2, len(response.Accounts))
				}
				limited := wants[1:3] // slice to offset and limit
				DebugAccounts("limited", limited)
				for i, got := range response.Accounts {
					want := limited[i]
					if err := CompareAccounts(got, want); err != nil {
						DebugAccount(got)
						DebugAccount(want)
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
					}
				}
			},
		},
	}

	for _, test := range tests {
		res, err := NewQuerier(&keeper).OwnerAccounts(sdk.WrapSDKContext(ctx), &test.request)
		if test.wantErr != nil && test.wantErr(err) != nil {
			t.Fatalf("failed err test on: %s", err)
		}
		if test.validate != nil {
			test.validate(res)
		}
	}
}

func TestResourceAccounts(t *testing.T) {
	keeper, ctx, _ := NewTestKeeper(t, false)
	_, _, _, accountsByResource := populateAccounts(t, keeper, ctx)
	tests := map[string]struct {
		request  types.QueryResourceAccountsRequest
		wantErr  func(error) error
		validate func(*types.QueryResourceAccountsResponse)
	}{
		"query non-existent resource": {
			request: types.QueryResourceAccountsRequest{
				Uri:        "bogus",
				Resource:   "equally bogus",
				Pagination: nil,
			},
			wantErr: nil,
			validate: func(response *types.QueryResourceAccountsResponse) {
				if response.Accounts == nil {
					t.Fatalf("wanted non-nil accounts")
				}
				if len(response.Accounts) != 0 {
					t.Fatalf("wanted 0 accounts")
				}
			},
		},
		"query invalid uri": {
			request: types.QueryResourceAccountsRequest{
				Uri:        "",
				Resource:   "valid",
				Pagination: nil,
			},
			wantErr: func(err error) error {
				if !errors.Is(err, sdkerrors.Wrapf(types.ErrInvalidResource, "''")) {
					t.Fatal("wrong error")
				}
				return nil
			},
			validate: nil,
		},
		"query invalid resource": {
			request: types.QueryResourceAccountsRequest{
				Uri:        "valid",
				Resource:   "",
				Pagination: nil,
			},
			wantErr: func(err error) error {
				if !errors.Is(err, sdkerrors.Wrapf(types.ErrInvalidResource, "''")) {
					t.Fatal("wrong error")
				}
				return nil
			},
			validate: nil,
		},
		"query valid resource without pagination": {
			request: types.QueryResourceAccountsRequest{
				Uri:        resources[0].URI,
				Resource:   resources[0].Resource,
				Pagination: nil,
			},
			wantErr: nil,
			validate: func(response *types.QueryResourceAccountsResponse) {
				if response.Accounts == nil {
					t.Fatal("wanted non-nil accounts")
				}
				key := string(types.GetResourceKey(resources[0].URI, resources[0].Resource))
				wants := accountsByResource[key]
				if len(response.Accounts) != len(wants) {
					t.Fatalf("wanted %d accounts, got %d", len(wants), len(response.Accounts))
				}
				for i, got := range response.Accounts {
					want := wants[i]
					if err := CompareAccounts(got, want); err != nil {
						DebugAccount(got)
						DebugAccount(want)
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
					}
				}
			},
		},
		"query valid owner with pagination": {
			request: types.QueryResourceAccountsRequest{
				Uri:      resources[0].URI,
				Resource: resources[0].Resource,
				Pagination: &query.PageRequest{
					Key:        nil,
					Offset:     1,
					Limit:      2,
					CountTotal: true,
				},
			},
			wantErr: nil,
			validate: func(response *types.QueryResourceAccountsResponse) {
				if response.Page == nil {
					t.Fatal("wanted non-nil Page")
				}
				key := string(types.GetResourceKey(resources[0].URI, resources[0].Resource))
				wants := accountsByResource[key]
				if response.Page.Total != uint64(len(wants)) {
					t.Fatalf("wanted %d accounts, got %d", len(wants), response.Page.Total)
				}
				if response.Accounts == nil {
					t.Fatalf("wanted non-nil accounts")
				}
				if len(response.Accounts) != 2 {
					t.Fatalf("wanted %d accounts, got %d", 2, len(response.Accounts))
				}
				limited := wants[1:3] // slice to offset and limit
				DebugAccounts("limited", limited)
				for i, got := range response.Accounts {
					want := limited[i]
					if err := CompareAccounts(got, want); err != nil {
						DebugAccount(got)
						DebugAccount(want)
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
					}
				}
			},
		},
	}

	for _, test := range tests {
		res, err := NewQuerier(&keeper).ResourceAccounts(sdk.WrapSDKContext(ctx), &test.request)
		if test.wantErr != nil && test.wantErr(err) != nil {
			t.Fatalf("failed err test on: %s", err)
		}
		if test.validate != nil {
			test.validate(res)
		}
	}
}
