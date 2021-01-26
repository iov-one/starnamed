package keeper

import (
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestDomain_Transfer(t *testing.T) {
	// defines test prereqs
	init := func() (k Keeper, ctx sdk.Context, ex *DomainExecutor) {
		k, ctx, _ = NewTestExecutorKeeper(t, false)
		domain := types.Domain{
			Name:       "test",
			Admin:      BobKey,
			ValidUntil: 1,
			Type:       types.OpenDomain,
			Broker:     nil,
		}
		acc1 := types.Account{
			Domain:       "test",
			Name:         utils.StrPtr("1"),
			Owner:        BobKey,
			ValidUntil:   1,
			Resources:    nil,
			Certificates: nil,
			Broker:       nil,
			MetadataURI:  "",
		}
		acc2 := types.Account{
			Domain:       "test",
			Name:         utils.StrPtr("2"),
			Owner:        BobKey,
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
			Owner:  CharlieKey,
		}
		domains := k.DomainStore(ctx)
		accounts := k.AccountStore(ctx)
		NewDomainExecutor(ctx, domain).WithDomains(&domains).WithAccounts(&accounts).Create()
		NewAccountExecutor(ctx, acc1).WithAccounts(&accounts).Create()
		NewAccountExecutor(ctx, acc2).WithAccounts(&accounts).Create()
		NewAccountExecutor(ctx, acc3).WithAccounts(&accounts).Create()
		ex = NewDomainExecutor(ctx, domain).WithDomains(&domains).WithAccounts(&accounts)
		return
	}
	t.Run("success init", func(t *testing.T) {
		k, ctx, _ := init()
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		n := 0
		account := new(types.Account)
		for ; cursor.Valid(); cursor.Next() {
			account.Reset() // stictly not necessary but is required in most cases, so consider it best practice
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
		ex.Transfer(types.TransferOwned, AliceKey)
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		acc := new(types.Account)
		for ; cursor.Valid(); cursor.Next() {
			acc.Reset()
			if err = cursor.Read(acc); err != nil {
				t.Fatal(err)
			}
			if !acc.Owner.Equals(AliceKey) && *acc.Name != "not-owned" {
				t.Fatal("owner mismatch")
			}
			if *acc.Name == "not-owned" && !acc.Owner.Equals(CharlieKey) {
				t.Fatal("a not owned account was transferred")
			}
		}
	})
	t.Run("transfer-flush", func(t *testing.T) {
		k, ctx, ex := init()
		ex.Transfer(types.TransferFlush, AliceKey)
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		acc := new(types.Account)
		emptyAccountExists := false
		for ; cursor.Valid(); cursor.Next() {
			acc.Reset()
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
		ex.Transfer(types.TransferResetNone, AliceKey)
		cursor, err := k.AccountStore(ctx).Query().Where().Index(types.AccountDomainIndex).Equals([]byte("test")).Do()
		if err != nil {
			t.Fatal(err)
		}
		acc := new(types.Account)
		for ; cursor.Valid(); cursor.Next() {
			acc.Reset()
			if err = cursor.Read(acc); err != nil {
				t.Fatal(err)
			}
			switch *acc.Name {
			case types.EmptyAccountName:
				if !acc.Owner.Equals(AliceKey) {
					t.Fatal("owner mismatch")
				}
			case "1":
				if !acc.Owner.Equals(BobKey) {
					t.Fatal("owner mismatch")
				}
			case "2":
				if !acc.Owner.Equals(BobKey) {
					t.Fatal("owner mismatch")
				}
			case "not-owned":
				if !acc.Owner.Equals(CharlieKey) {
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
		testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
		renewalPeriod := time.Duration(20)
		setConfig := GetConfigSetter(testKeeper.ConfigurationKeeper).SetConfig
		setConfig(testCtx, configuration.Config{
			DomainRenewalPeriod: renewalPeriod * time.Second,
		})
		domains := testKeeper.DomainStore(testCtx)
		accounts := testKeeper.AccountStore(testCtx)
		conf := testKeeper.ConfigurationKeeper.GetConfiguration(testCtx)
		ex := NewDomainExecutor(testCtx, testDomain).WithDomains(&domains).WithAccounts(&accounts).WithConfiguration(conf)
		ex.Renew()
		newDom := new(types.Domain)
		if err := testKeeper.DomainStore(testCtx).Read(testDomain.PrimaryKey(), newDom); err != nil {
			t.Fatal("domain does not exist anymore")
		}
		if newDom.ValidUntil != testDomain.ValidUntil+int64(renewalPeriod) {
			t.Fatal("mismatched times")
		}
	})
	t.Run("success renew from account, not config", func(t *testing.T) {
		testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
		var accValidUntil int64 = 10000
		domains := testKeeper.DomainStore(testCtx)
		accounts := testKeeper.AccountStore(testCtx)
		ex := NewDomainExecutor(testCtx, testDomain).WithDomains(&domains).WithAccounts(&accounts)
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
		testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
		domains := testKeeper.DomainStore(testCtx)
		accounts := testKeeper.AccountStore(testCtx)
		NewDomainExecutor(testCtx, testDomain).WithDomains(&domains).WithAccounts(&accounts).Delete()
		if err := testKeeper.DomainStore(testCtx).Read(testDomain.PrimaryKey(), &types.Domain{}); err == nil {
			t.Fatal("domain was not deleted")
		}
	})
}
