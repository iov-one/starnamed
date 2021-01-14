package keeper

import (
	"errors"
	"fmt"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkErrors "github.com/cosmos/cosmos-sdk/types/errors"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/query"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestDomainAccounts(t *testing.T) {
	domain := "test"
	admin := aliceAddr
	owners := []sdk.AccAddress{aliceAddr, bobAddr, aliceAddr, bobAddr, aliceAddr}
	keeper, ctx, _ := NewTestKeeper(t, false)
	if err := keeper.DomainStore(ctx).Create(&types.Domain{Name: domain, Admin: admin}); err != nil {
		t.Fatalf("failed to create domain '%s'", domain)
	}

	for i, owner := range owners {
		name := fmt.Sprintf("%d", i)
		account := types.Account{
			Domain: domain,
			Name:   utils.StrPtr(name),
			Owner:  owner,
		}
		if err := keeper.AccountStore(ctx).Create(&account); err != nil {
			t.Fatalf("failed to create account '%s'", name)
		}
	}

	// populate accounts in the IAVL order, which is a function of insertion order
	accounts := make([]*types.Account, 0)
	cursor, err := keeper.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte(domain)).Do()
	if err != nil {
		t.Fatal(err)
	}
	for ; cursor.Valid(); cursor.Next() {
		account := new(types.Account)
		if err := cursor.Read(account); err != nil {
			t.Fatal(err)
		}
		accounts = append(accounts, account)
	}

	tests := map[string]struct {
		request  types.QueryDomainAccountsRequest
		wantErr  error
		validate func(*types.QueryDomainAccountsResponse)
	}{
		"query invalid domain": {
			request: types.QueryDomainAccountsRequest{
				Domain:     "",
				Pagination: nil,
			},
			wantErr:  sdkerrors.Wrapf(types.ErrInvalidDomainName, "''"),
			validate: nil,
		},
		"query valid domain without pagination": {
			request: types.QueryDomainAccountsRequest{
				Domain:     domain,
				Pagination: nil,
			},
			wantErr: nil,
			validate: func(response *types.QueryDomainAccountsResponse) {
				if len(response.Accounts) != len(accounts) {
					t.Fatalf("wanted %d accounts, got %d", len(response.Accounts), len(accounts))
				}
				for i, got := range response.Accounts {
					want := accounts[i]
					if err := CompareAccounts(got, want); err != nil {
						t.Fatal(sdkErrors.Wrapf(err, "%s%s%s", *want.Name, types.StarnameSeparator, want.Domain))
					}
				}
			},
		},
		"query valid domain with pagination": {
			request: types.QueryDomainAccountsRequest{
				Domain: domain,
				Pagination: &query.PageRequest{
					Key:        nil,
					Offset:     1,
					Limit:      2,
					CountTotal: true,
				},
			},
			wantErr: nil,
			validate: func(response *types.QueryDomainAccountsResponse) {
				if len(response.Accounts) != 2 {
					t.Fatalf("wanted %d accounts, got %d", 2, len(response.Accounts))
				}
				limited := accounts[1:3] // slice to offset and limit
				for i, got := range response.Accounts {
					want := limited[i]
					if err := CompareAccounts(got, want); err != nil {
						t.Fatal(sdkErrors.Wrapf(err, "%s%s%s", *want.Name, types.StarnameSeparator, want.Domain))
					}
				}
			},
		},
	}

	for _, test := range tests {
		res, err := NewQuerier(&keeper).DomainAccounts(sdk.WrapSDKContext(ctx), &test.request)
		if test.wantErr != nil && !errors.Is(err, test.wantErr) {
			t.Fatalf("wanted err: %s, got: %s", test.wantErr, err)
		}

		if test.validate != nil {
			test.validate(res)
		}
	}
}
