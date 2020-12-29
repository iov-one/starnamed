package executor

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestDomain_Transfer(t *testing.T) {
	// defines test prereqs
	init := func() (k keeper.Keeper, ctx sdk.Context, ex *Domain) {
		k, ctx, _ = keeper.NewTestKeeper(t, false)
		domain := types.Domain{
			Name:       "test",
			Admin:      keeper.BobKey,
			ValidUntil: 1,
			Type:       types.OpenDomain,
			Broker:     nil,
		}
		acc1 := types.Account{
			Domain:       "test",
			Name:         utils.StrPtr("1"),
			Owner:        keeper.BobKey,
			ValidUntil:   1,
			Resources:    nil,
			Certificates: nil,
			Broker:       nil,
			MetadataURI:  "",
		}
		acc2 := types.Account{
			Domain:       "test",
			Name:         utils.StrPtr("2"),
			Owner:        keeper.BobKey,
			ValidUntil:   1,
			Resources:    nil,
			Certificates: nil,
			Broker:       nil,
			MetadataURI:  "",
		}
		// add account not owned
		acc3 := types.Account{
			Domain: "test",
			Name:   utils.StrPtr("not-owned"),
			Owner:  keeper.CharlieKey,
		}
		NewDomain(ctx, k, domain).Create()
		NewAccount(ctx, k, acc1).Create()
		NewAccount(ctx, k, acc2).Create()
		NewAccount(ctx, k, acc3).Create()
		ex = NewDomain(ctx, k, domain)
		return
	}
	t.Run("success init", func(t *testing.T) {
		_, _, ex := init()
		//cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		cursor, err := ex.accounts.Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		account := new(types.Account)
		for ; cursor.Valid(); cursor.Next() {
			if err := cursor.Read(account); err != nil {
				t.Fatal(err)
			}
			n++
		}
		if n != 4 {
			t.Fatalf("expected accounts acc{1,2,3} and the empty account but got a total of %d", n)
		}
	})
	t.Run("transfer owned", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferOwned, keeper.AliceKey)
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		acc := new(types.Account)
		for ; cursor.Valid(); cursor.Next() {
			if err = cursor.Read(acc); err != nil {
				t.Fatal(err)
			}
			if !acc.Owner.Equals(keeper.AliceKey) && *acc.Name != "not-owned" {
				t.Fatal("owner mismatch")
			}
			if *acc.Name == "not-owned" && !acc.Owner.Equals(keeper.CharlieKey) {
				t.Fatal("a not owned account was transferred")
			}
		}
	})
	t.Run("transfer-flush", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferFlush, keeper.AliceKey)
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		acc := new(types.Account)
		emptyAccountExists := false
		for ; cursor.Valid(); cursor.Next() {
			if err = cursor.Read(acc); err != nil {
				t.Fatal(err)
			}
			// only empty account is expected
			if *acc.Name != types.EmptyAccountName {
				t.Fatalf("only empty account is expected to exist, got: %s", *acc.Name)
			}
			if *acc.Name == types.EmptyAccountName {
				emptyAccountExists = true
			}
		}
		if !emptyAccountExists {
			t.Fatal("empty account not found")
		}
	})
	t.Run("transfer-reset-none", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferResetNone, keeper.AliceKey)
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		acc := new(types.Account)
		for ; cursor.Valid(); cursor.Next() {
			if err = cursor.Read(acc); err != nil {
				t.Fatal(err)
			}
			switch *acc.Name {
			case types.EmptyAccountName:
				if !acc.Owner.Equals(keeper.AliceKey) {
					t.Fatal("owner mismatch")
				}
			case "1":
				if !acc.Owner.Equals(keeper.BobKey) {
					t.Fatal("owner mismatch")
				}
			case "2":
				if !acc.Owner.Equals(keeper.BobKey) {
					t.Fatal("owner mismatch")
				}
			case "not-owned":
				if !acc.Owner.Equals(keeper.CharlieKey) {
					t.Fatal("owner mismatch")
				}
			default:
				t.Fatalf("unexpected account found: %s", *acc.Name)
			}
		}
	})
}

func TestDomain_Renew(t *testing.T) {
	t.Run("success renew from config", func(t *testing.T) {
		testCtx, _ := testCtx.CacheContext()
		ex := NewDomain(testCtx, testKeeper, testDomain)
		ex.Renew()
		newDom := new(types.Domain)
		if err := testKeeper.DomainStore(testCtx).Read(testDomain.PrimaryKey(), newDom); err != nil {
			t.Fatal("domain does not exist anymore")
		}
		if newDom.ValidUntil != testDomain.ValidUntil+int64(testConfig.DomainRenewalPeriod.Seconds()) {
			t.Fatal("mismatched times")
		}
	})
	t.Run("success renew from account", func(t *testing.T) {
		testCtx, _ := testCtx.CacheContext()
		var accValidUntil int64 = 10000
		ex := NewDomain(testCtx, testKeeper, testDomain)
		ex.Renew(accValidUntil)
		newDom := new(types.Domain)
		if err := testKeeper.DomainStore(testCtx).Read(testDomain.PrimaryKey(), newDom); err != nil {
			t.Fatal("domain does not exist anymore")
		}
		if newDom.ValidUntil != accValidUntil {
			t.Fatal("mismatched times")
		}
	})
}

func TestDomain_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		testCtx, _ := testCtx.CacheContext()
		NewDomain(testCtx, testKeeper, testDomain).Delete()
		if err := testKeeper.DomainStore(testCtx).Read(testDomain.PrimaryKey(), &types.Domain{}); err == nil {
			t.Fatal("domain was not deleted")
		}
	})
}
