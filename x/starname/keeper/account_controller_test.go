package keeper

import (
	"errors"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/mock"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestAccount_transferable(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
	// create mock domains and accounts
	// create open domain
	ds := k.DomainStore(ctx)
	as := k.AccountStore(ctx)
	ds.Create(&types.Domain{
		Name:       "open",
		Admin:      AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
		Type:       types.OpenDomain,
	})
	// creat open domain account
	as.Create(&types.Account{
		Domain: "open",
		Name:   utils.StrPtr("test"),
		Owner:  BobKey,
	})
	// create closed domain
	ds.Create(&types.Domain{
		Name:       "closed",
		Admin:      AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
		Type:       types.ClosedDomain,
	})
	// create closed domain account
	as.Create(&types.Account{
		Domain: "closed",
		Name:   utils.StrPtr("test"),
		Owner:  BobKey,
	})
	// run tests
	t.Run("closed domain", func(t *testing.T) {
		dc := NewDomainController(ctx, "closed").WithDomains(&ds)
		acc := NewAccountController(ctx, "closed", "test").WithAccounts(&as).WithDomainController(dc)
		// test success
		err := acc.
			TransferableBy(AliceKey).
			Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		// test failure
		err = acc.TransferableBy(BobKey).Validate()
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
	})
	t.Run("open domain", func(t *testing.T) {
		dc := NewDomainController(ctx, "open").WithDomains(&ds)
		acc := NewAccountController(ctx, "open", "test").WithAccounts(&as).WithDomainController(dc)
		err := acc.TransferableBy(BobKey).Validate()
		// test success
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		// test failure
		err = acc.TransferableBy(AliceKey).Validate()
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
	})
}

func TestAccount_Renewable(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
	ctx = ctx.WithBlockTime(time.Unix(1, 0))
	setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
	setConfig(ctx, configuration.Config{
		AccountRenewalCountMax: 1,
		AccountRenewalPeriod:   10 * time.Second,
	})
	domains := k.DomainStore(ctx)
	accounts := k.AccountStore(ctx)
	NewDomainExecutor(ctx, types.Domain{
		Name:       "open",
		Admin:      AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
		Type:       types.OpenDomain,
	}).WithDomains(&domains).WithAccounts(&accounts).Create()
	NewAccountExecutor(ctx, types.Account{
		Domain:     "open",
		Name:       utils.StrPtr("test"),
		ValidUntil: time.Unix(18, 0).Unix(),
		Owner:      BobKey,
	}).WithAccounts(&accounts).Create()
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)

	// 18(AccountValidUntil) + 10 (AccountRP) = 28 newValidUntil
	// no need to test closed domain since its not renewable
	t.Run("open domain", func(t *testing.T) {
		// 7(time) + 2(AccountRCM) * 10(AccountRP) = 27 maxValidUntil
		acc := NewAccountController(ctx.WithBlockTime(time.Unix(7, 0)), "open", "test").WithAccounts(&accounts).WithConfiguration(conf)
		err := acc.Renewable().Validate()
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
		// 100(time) + 2(AccountRCM) * 10(AccountRP) = 120 maxValidUntil
		acc = NewAccountController(ctx.WithBlockTime(time.Unix(100, 0)), "open", "test").WithAccounts(&accounts).WithConfiguration(conf)
		if err := acc.Renewable().Validate(); err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
}

func TestAccount_existence(t *testing.T) {
	k, ctx, _ := NewTestKeeper(t, true)
	accounts := k.AccountStore(ctx)
	// insert mock account
	accounts.Create(&types.Account{
		Domain:     "test",
		Name:       utils.StrPtr("test"),
		Owner:      AliceKey,
		ValidUntil: time.Now().Add(100 * time.Hour).Unix(),
	})
	// run MustExist test
	t.Run("must exist success", func(t *testing.T) {
		acc := NewAccountController(ctx, "test", "test").WithAccounts(&accounts)
		err := acc.MustExist().Validate()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
	})
	t.Run("must exist fail", func(t *testing.T) {
		acc := NewAccountController(ctx, "test", "does not exist").WithAccounts(&accounts)
		err := acc.MustExist().Validate()
		if !errors.Is(err, types.ErrAccountDoesNotExist) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountDoesNotExist, err)
		}
	})
	// run MustNotExist test
	t.Run("must not exist success", func(t *testing.T) {
		acc := NewAccountController(ctx, "test", "does not exist").WithAccounts(&accounts)
		err := acc.MustNotExist().Validate()
		if err != nil {
			t.Errorf("got error: %s", err)
		}
	})
	t.Run("must not exist fail", func(t *testing.T) {
		acc := NewAccountController(ctx, "test", "test").WithAccounts(&accounts)
		err := acc.MustNotExist().Validate()
		if !errors.Is(err, types.ErrAccountExists) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountExists, err)
		}
	})
}

func TestAccount_requireAccount(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		k, ctx, _ := NewTestKeeper(t, true)
		accounts := k.AccountStore(ctx)
		alice, _ := mock.Addresses()
		accounts.Create(&types.Account{
			Domain: "test",
			Name:   utils.StrPtr("test"),
			Owner:  alice,
		})
		ctrl := NewAccountController(ctx, "test", "test").WithAccounts(&accounts)
		err := ctrl.requireAccount()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("does not exist", func(t *testing.T) {
		k, ctx, _ := NewTestKeeper(t, true)
		as := k.AccountStore(ctx)
		ctrl := NewAccountController(ctx, "test", "test").WithAccounts(&as)
		err := ctrl.requireAccount()
		if !errors.Is(err, types.ErrAccountDoesNotExist) {
			t.Fatalf("want: %s, got: %s", types.ErrAccountDoesNotExist, err)
		}
	})
}

func TestAccount_certNotExist(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &AccountController{
			account: &types.Account{
				Certificates: [][]byte{[]byte("test-cert")},
			},
		}
		err := acc.CertificateNotExist([]byte("does not exist")).Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("cert exists", func(t *testing.T) {
		acc := &AccountController{
			account: &types.Account{
				Certificates: [][]byte{[]byte("test-cert"), []byte("exists")},
			},
		}
		i := new(int)
		err := acc.CertificateExists([]byte("exists"), i).Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
		if *i != 1 {
			t.Fatalf("unexpected index pointer: %d", *i)
		}
	})
}

func TestAccount_notExpired(t *testing.T) {
	closedDomain := (&DomainController{}).WithDomain(types.Domain{
		Type: types.ClosedDomain,
	})
	openDomain := (&DomainController{}).WithDomain(types.Domain{
		Type: types.OpenDomain,
	})
	t.Run("success", func(t *testing.T) {
		acc := (&AccountController{
			account: &types.Account{
				ValidUntil: 10,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(0, 0)),
		}).WithDomainController(openDomain)
		err := acc.NotExpired().Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("expired", func(t *testing.T) {
		acc := (&AccountController{
			account: &types.Account{
				ValidUntil: 10,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(11, 0)),
		}).WithDomainController(openDomain)
		err := acc.NotExpired().Validate()
		if !errors.Is(err, types.ErrAccountExpired) {
			t.Fatalf("want error: %s, got: %s", types.ErrAccountExpired, err)
		}
	})
	t.Run("success account expired but in closed domain", func(t *testing.T) {
		acc := (&AccountController{
			account: &types.Account{
				ValidUntil: 1,
			},
			ctx: sdk.Context{}.WithBlockTime(time.Unix(20, 0)),
		}).WithDomainController(closedDomain)
		err := acc.NotExpired().Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
}

func TestAccount_ownedBy(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		alice, _ := mock.Addresses()
		acc := &AccountController{
			account: &types.Account{Owner: alice},
		}
		err := acc.OwnedBy(alice).Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("bad owner", func(t *testing.T) {
		alice, bob := mock.Addresses()
		acc := &AccountController{
			account: &types.Account{Owner: alice},
		}
		err := acc.OwnedBy(bob).Validate()
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrUnauthorized)
		}
	})
}

func TestAccount_validName(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		acc := &AccountController{
			account: &types.Account{Name: utils.StrPtr("valid")},
			conf:    &configuration.Config{ValidAccountName: "^(.*?)?"},
		}
		err := acc.ValidName().Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("success", func(t *testing.T) {
		acc := &AccountController{
			name: "not valid",
			conf: &configuration.Config{ValidAccountName: "$^"},
		}
		err := acc.ValidName().Validate()
		if !errors.Is(err, types.ErrInvalidAccountName) {
			t.Fatalf("unexpected error: %s, wanted: %s", err, types.ErrInvalidAccountName)
		}
	})
}

func TestAccountRegistrableBy(t *testing.T) {
	closedDomain := (&DomainController{}).WithDomain(types.Domain{
		Type:  types.ClosedDomain,
		Admin: AliceKey,
	})
	openDomain := (&DomainController{}).WithDomain(types.Domain{
		Type: types.OpenDomain,
	})
	t.Run("success in closed domain", func(t *testing.T) {
		acc := (&AccountController{}).WithDomainController(closedDomain)
		err := acc.RegistrableBy(AliceKey).Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
	t.Run("fail in closed domain", func(t *testing.T) {
		acc := (&AccountController{}).WithDomainController(closedDomain)
		err := acc.RegistrableBy(BobKey).Validate()
		if !errors.Is(err, types.ErrUnauthorized) {
			t.Fatalf("want: %s, got: %s", types.ErrUnauthorized, err)
		}
	})
	t.Run("success other domain type", func(t *testing.T) {
		acc := (&AccountController{}).WithDomainController(openDomain)
		err := acc.RegistrableBy(AliceKey).Validate()
		if err != nil {
			t.Fatalf("got error: %s", err)
		}
	})
}
