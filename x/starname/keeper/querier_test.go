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

	accounts := make([]*types.Account, 0) // in primary key order
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
		accounts = append(accounts, &account)
		//fmt.Printf("%s %x %x %x\n", account.GetStarname(), account.PrimaryKey(), account.SecondaryKeys()[0].Value, account.SecondaryKeys()[1].Value)
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
					//fmt.Printf("got  %s %x %x %x\n", got.GetStarname(), got.PrimaryKey(), got.SecondaryKeys()[0].Value, got.SecondaryKeys()[1].Value)
					//fmt.Printf("want %s %x %x %x\n", want.GetStarname(), want.PrimaryKey(), want.SecondaryKeys()[0].Value, want.SecondaryKeys()[1].Value)
					if err := CompareAccounts(got, want); err != nil {
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
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
					//fmt.Printf("got  %s %x %x %x\n", got.GetStarname(), got.PrimaryKey(), got.SecondaryKeys()[0].Value, got.SecondaryKeys()[1].Value)
					//fmt.Printf("want %s %x %x %x\n", want.GetStarname(), want.PrimaryKey(), want.SecondaryKeys()[0].Value, want.SecondaryKeys()[1].Value)
					if err := CompareAccounts(got, want); err != nil {
						t.Fatal(sdkErrors.Wrapf(err, want.GetStarname()))
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
