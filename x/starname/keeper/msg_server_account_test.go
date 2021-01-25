package keeper

import (
	"bytes"
	"errors"
	"reflect"
	"testing"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/keeper/executor"
	"github.com/iov-one/starnamed/x/starname/types"
)

func Test_Close_handlerMsgAddAccountCertificate(t *testing.T) {
	cases := map[string]SubTest{
		"does not respect account expiration": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
					Type:       types.ClosedDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: nil,
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
		},
	}
	RunTests(t, cases)
}

func Test_Open_handlerMsgAddAccountCertificate(t *testing.T) {
	cases := map[string]SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
					Type:       types.OpenDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
		},
	}
	RunTests(t, cases)
}

func Test_Common_handlerMsgAddAccountCertificate(t *testing.T) {
	cases := map[string]SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "does not exist",
					Name:           "",
					Owner:          BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add expired domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: 0,
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "",
					Owner:          BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "does not exist",
					Owner:          nil,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg owner is not account owner": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"admin cannot add cert": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          BobKey,
					NewCertificate: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate exists": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
					AccountGracePeriod:  1000 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        AliceKey,
					Certificates: [][]byte{[]byte("test")},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: []byte("test"),
				})
				if !errors.Is(err, types.ErrCertificateExists) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"certificate size exceeded": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  4,
					AccountGracePeriod:  1000 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        AliceKey,
					Certificates: nil,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: []byte("12345"),
				})
				if !errors.Is(err, types.ErrCertificateSizeExceeded) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
				_, err = handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: []byte("1234"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"certificate limit reached": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  100,
					AccountGracePeriod:  1000 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        AliceKey,
					Certificates: [][]byte{[]byte("1")},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: []byte("12345"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
				_, err = handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: []byte("1234"),
				})
				if !errors.Is(err, types.ErrCertificateLimitReached) {
					t.Fatalf("handlerMsgAddAccountCertificate() expected error: %s, got: %s", types.ErrCertificateExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					CertificateCountMax: 2,
					CertificateSizeMax:  4,
					AccountGracePeriod:  1000 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgAddAccountCertificate(ctx, k, &types.MsgAddAccountCertificate{
					Domain:         "test",
					Name:           "test",
					Owner:          AliceKey,
					NewCertificate: []byte("test"),
				})
				if err != nil {
					t.Fatalf("handlerMsgAddAccountCertificate() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				expected := [][]byte{[]byte("test")}
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("test")}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				if !reflect.DeepEqual(account.Certificates, expected) {
					t.Fatalf("handlerMsgAddAccountCertificate: got: %#v, expected: %#v", account.Certificates, expected)
				}
			},
		},
	}
	RunTests(t, cases)
}

func Test_Closed_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]SubTest{
		"does not respect account valid until": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
					Type:       types.ClosedDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   0,
					Owner:        AliceKey,
					Certificates: [][]byte{[]byte("test")},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Domain: "test", Name: utils.StrPtr("test")}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				for _, cert := range account.Certificates {
					if bytes.Equal(cert, []byte("test")) {
						t.Fatalf("handlerMsgDeleteAccountCertificates() certificate not deleted")
					}
				}
				// success
			},
		},
	}

	RunTests(t, cases)
}

func Test_Open_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
					Type:       types.OpenDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   0,
					Owner:        AliceKey,
					Certificates: [][]byte{[]byte("test")},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             AliceKey,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
		},
	}

	RunTests(t, cases)
}
func Test_Common_handlerMsgDeleteAccountCertificate(t *testing.T) {
	cases := map[string]SubTest{
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "does not exist",
					DeleteCertificate: nil,
					Owner:             BobKey,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"msg signer is not account owner": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"domain admin cannot delete cert": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      BobKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: nil,
					Owner:             AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"certificate does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("does not exist"),
					Owner:             AliceKey,
				})
				if !errors.Is(err, types.ErrCertificateDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccountCertificate() expected error: %s, got: %s", types.ErrCertificateDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Admin:      AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add mock account
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Owner:        AliceKey,
					Certificates: [][]byte{[]byte("test")},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccountCertificate(ctx, k, &types.MsgDeleteAccountCertificate{
					Domain:            "test",
					Name:              "test",
					DeleteCertificate: []byte("test"),
					Owner:             AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccountCertificates() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// check if certificate is still present
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				for _, cert := range account.Certificates {
					if bytes.Equal(cert, []byte("test")) {
						t.Fatalf("handlerMsgDeleteAccountCertificates() certificate not deleted")
					}
				}
				// success
			},
		},
	}

	RunTests(t, cases)
}

func Test_Closed_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]SubTest{
		"domain admin can delete": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					Type:       types.ClosedDomain,
					ValidUntil: types.MaxValidUntil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  BobKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"domain expired": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					Type:       types.ClosedDomain,
					ValidUntil: 2,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  BobKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 2,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
		},
		"account owner cannot delete": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					Type:       types.ClosedDomain,
					ValidUntil: types.MaxValidUntil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %s", err)
				}
			},
		},
	}
	RunTests(t, cases)
}

func Test_Open_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]SubTest{
		"domain admin cannot can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      RegexMatchNothing,
					ValidURI:           RegexMatchAll,
					AccountGracePeriod: 1000 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  BobKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 3,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
		},
		"no domain valid until check": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchNothing,
					ValidURI:            RegexMatchAll,
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					Type:       types.OpenDomain,
					ValidUntil: 2,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  BobKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"only account owner can delete before grace period": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      RegexMatchNothing,
					ValidURI:           RegexMatchAll,
					AccountGracePeriod: 10 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 5,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %v", err)
				}
				// anyone test
				_, err = handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("unexpected error: %v", err)
				}
				// account owner test
				_, err = handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"domain admin can delete after grace": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      RegexMatchNothing,
					ValidURI:           RegexMatchAll,
					AccountGracePeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"anyone can delete after grace": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:      RegexMatchNothing,
					ValidURI:           RegexMatchAll,
					AccountGracePeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					Type:       types.OpenDomain,
					ValidUntil: types.MaxValidUntil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 100,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// admin test
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  CharlieKey,
				})
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}
	RunTests(t, cases)
}

func Test_Common_handlerMsgDeleteAccount(t *testing.T) {
	cases := map[string]SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {

			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "does not exist",
					Name:   "does not exist",
					Owner:  AliceKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Admin: BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgDeleteAccount() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success domain owner": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Admin: AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  BobKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
		"success account owner": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Admin: AliceKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain: "test",
					Name:   utils.StrPtr("test"),
					Owner:  BobKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgDeleteAccount(ctx, k, &types.MsgDeleteAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgDeleteAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err == nil {
					t.Fatalf("handlerMsgDeleteAccount() account was not deleted")
				}
			},
		},
	}

	// run tests
	RunTests(t, cases)
}

func Test_ClosedDomain_handlerMsgRegisterAccount(t *testing.T) {
	testCases := map[string]SubTest{
		"only domain admin can register account": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidResource:        RegexMatchNothing,
					ValidURI:             RegexMatchAll,
					DomainRenewalPeriod:  10,
					AccountRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      BobKey,
					Registerer: BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
				_, err = handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test2",
					Owner:      BobKey,
					Registerer: AliceKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
		},
		"account valid until is set to max": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchNothing, // don't match anything
					ValidURI:            RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a closed domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      BobKey,
					Registerer: BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account test not found")
				}
				if account.ValidUntil != types.MaxValidUntil {
					t.Fatalf("unexpected account valid until %d", account.ValidUntil)
				}
			},
		},
		"account owner can be different than domain admin": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchNothing, // don't match anything
					ValidURI:            RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a closed domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Registerer: BobKey,
					Owner:      BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
				_, err = handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test2",
					Registerer: BobKey,
					Owner:      AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
		},
	}
	// run tests
	RunTests(t, testCases)
}

func Test_OpenDomain_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]SubTest{
		"account valid until is now plus config account renew": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:        RegexMatchNothing, // don't match anything
					ValidURI:             RegexMatchAll,     // match all
					DomainRenewalPeriod:  10 * time.Second,
					AccountRenewalPeriod: 10 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a closed domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      BobKey,
					Registerer: BobKey,
				})
				if err != nil {
					t.Fatalf("handlerRegisterAccount() got error: %s", err)
				}
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account test not found")
				}
				expected := utils.TimeToSeconds(time.Unix(11, 0))
				if account.ValidUntil != expected {
					t.Fatalf("unexpected account valid until %d, expected %d", account.ValidUntil, expected)
				}
			},
		},
	}
	RunTests(t, testCases)
}

func Test_Common_handleMsgRegisterAccount(t *testing.T) {
	testCases := map[string]SubTest{
		"fail resource": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {

				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchNothing, // don't match anything
					ValidURI:            RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "won't work",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidResource) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
			AfterTest: nil,
		},
		"fail invalid uri": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidURI:            RegexMatchNothing, // don't match anything
					ValidResource:       RegexMatchAll,     // match all
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
					Resources: []*types.Resource{
						{
							URI:      "invalid blockchain id",
							Resource: "valid blockchain address",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidResource) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
			AfterTest: nil,
		},
		"fail invalid account name": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchAll,     // match all
					ValidURI:            RegexMatchAll,     // match all
					ValidAccountName:    RegexMatchNothing, // match nothing
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "this won't match",
					Owner:  AliceKey,
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrInvalidAccountName) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrInvalidAccountName, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain name does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchAll, // match all
					ValidURI:            RegexMatchAll, // match all
					ValidAccountName:    RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "this does not exist",
					Name:   "works",
					Owner:  AliceKey,
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"fail only owner of domain with superuser can register accounts": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchAll, // match all
					ValidURI:            RegexMatchAll, // match all
					ValidAccountName:    RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: 2,
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey, // invalid owner
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"fail domain has expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchAll, // match all
					ValidURI:            RegexMatchAll, // match all
					ValidAccountName:    RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: 0, // domain is expired
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"fail account exists": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchAll, // match all
					ValidURI:            RegexMatchAll, // match all
					ValidAccountName:    RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// add an account that we are gonna try to overwrite
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("exists"),
					Owner:        AliceKey,
					ValidUntil:   0,
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain: "test",
					Name:   "exists",
					Owner:  BobKey,
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if !errors.Is(err, types.ErrAccountExists) {
					t.Fatalf("handleMsgRegisterAccount() expected error: %s, got: %s", types.ErrAccountExists, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set regexp match nothing in resources
				// get set config function
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				// set configs with a domain regexp that matches nothing
				setConfig(ctx, configuration.Config{
					ValidResource:       RegexMatchAll, // match all
					ValidURI:            RegexMatchAll, // match all
					ValidAccountName:    RegexMatchAll, // match nothing
					DomainRenewalPeriod: 10,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// add a domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: time.Now().Add(100000 * time.Hour).Unix(),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handleMsgRegisterAccount(ctx, k, &types.MsgRegisterAccount{
					Domain:     "test",
					Name:       "test",
					Owner:      BobKey,
					Registerer: BobKey,
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
					Broker: nil,
				})
				if err != nil {
					t.Fatalf("handleMsgRegisterAccount() got error: %s", err)
				}
			},
			AfterTest: nil, // TODO fill with matching data
		},
	}
	// run tests
	RunTests(t, testCases)
}

func Test_Closed_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]SubTest{
		"account cannot be renewed since its max": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod: 1,
					AccountGracePeriod:   5,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// set mock domain
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.ClosedDomain,
					Admin: BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// set mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Unix(1000, 0)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if !errors.Is(err, types.ErrInvalidDomainType) {
					t.Fatalf("handlerMsgRenewAccount() want err: %s, got: %s", types.ErrInvalidDomainType, err)
				}
			},
		},
	}

	RunTests(t, cases)
}
func Test_Open_handlerMsgRenewAccount(t *testing.T) {
	cases := map[string]SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod: 1,
					AccountGracePeriod:   5,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "does not exist",
					Name:   "",
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod: 1,
					AccountGracePeriod:   5,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// set mock domain
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "does not exist",
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgRenewAccount() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"success domain grace period not updated": {
			TestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod:   1 * time.Second,
					AccountRenewalCountMax: 200000,
					AccountGracePeriod:     5 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// set mock domain
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Type:  types.OpenDomain,
					Admin: BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// set mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Unix(1, 0)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				want := ctx.BlockTime().Add(k.ConfigurationKeeper.GetConfiguration(ctx).AccountRenewalPeriod)
				if account.ValidUntil != utils.TimeToSeconds(want) {
					t.Fatalf("handlerMsgRenewAccount() want: %d, got: %d", want.Unix(), account.ValidUntil)
				}
			},
		},
		"success domain valid until updated": {
			BeforeTestBlockTime: 1,
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					AccountRenewalPeriod:   1 * time.Second,
					AccountRenewalCountMax: 200000,
					AccountGracePeriod:     5 * time.Second,
					DomainGracePeriod:      2 * time.Second,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// set mock domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Type:       types.OpenDomain,
					Admin:      BobKey,
					ValidUntil: 2,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// set mock account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Unix(1, 0)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			TestBlockTime: 10,
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgRenewAccount(ctx, k, &types.MsgRenewAccount{
					Domain: "test",
					Name:   "test",
				})
				if err != nil {
					t.Fatalf("handlerMsgRenewAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domain := new(types.Domain)
				if err := k.DomainStore(ctx).Read((&types.Domain{Name: "test"}).PrimaryKey(), domain); err != nil {
					t.Fatal("domain not found")
				}
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				if domain.ValidUntil != account.ValidUntil {
					t.Fatalf("handlerMsgRenewAccount() want: %d, got: %d", domain.ValidUntil, account.ValidUntil)
				}
			},
		},
	}

	RunTests(t, cases)
}

func Test_Closed_handlerMsgReplaceAccountResources(t *testing.T) {
	cases := map[string]SubTest{
		"fail does not respect account valid until": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
					ResourcesMax:  5,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
					Type:       types.ClosedDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() got error: %s", err)
				}
			},
		},
	}

	RunTests(t, cases)
}

func Test_Open_handlerMsgReplaceAccountResources(t *testing.T) {
	cases := map[string]SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
					ResourcesMax:  3,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
					Type:       types.OpenDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: AliceKey,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
	}

	RunTests(t, cases)
}
func Test_Common_handlerMsgReplaceAccountResources(t *testing.T) {
	cases := map[string]SubTest{
		"invalid blockchain resource": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchNothing,
					ValidResource: RegexMatchNothing,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "invalid",
							Resource: "invalid",
						},
					},
					Owner: AliceKey,
				})
				if !errors.Is(err, types.ErrInvalidResource) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
		},
		"resource limit exceeded": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
					ResourcesMax:  2,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
						{
							URI:      "valid1",
							Resource: "valid1",
						},
						{
							URI:      "valid2",
							Resource: "valid2",
						},
					},
					Owner: AliceKey,
				})
				if !errors.Is(err, types.ErrResourceLimitExceeded) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
				_, err = handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "invalid",
							Resource: "invalid",
						},
						{
							URI:      "invalid2",
							Resource: "invalid2",
						},
					},
					Owner: AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
				_, err = handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "invalid",
							Resource: "invalid",
						},
					},
					Owner: AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrInvalidResource, err)
				}
			},
		},
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
				})
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "does not exist",
					Name:   "",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Admin: BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "does not exist",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// set config to match all
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					ValidURI:      RegexMatchAll,
					ValidResource: RegexMatchAll,
					ResourcesMax:  5,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountResources(ctx, k, &types.MsgReplaceAccountResources{
					Domain: "test",
					Name:   "test",
					NewResources: []*types.Resource{
						{
							URI:      "valid",
							Resource: "valid",
						},
					},
					Owner: AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountResources() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				expected := []*types.Resource{{
					URI:      "valid",
					Resource: "valid",
				}}
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				if !reflect.DeepEqual(expected, account.Resources) {
					t.Fatalf("handlerMsgReplaceAccountResources() expected: %+v, got %+v", expected, account.Resources)
				}
			},
		},
	}
	// run tests
	RunTests(t, cases)
}

func Test_Closed_handlerMsgReplaceAccountMetadata(t *testing.T) {
	cases := map[string]SubTest{
		"account expiration not respected": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 100,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
					Type:       types.ClosedDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
		},
	}

	RunTests(t, cases)
}

func Test_Open_handlerMsgReplaceAccountMetadata(t *testing.T) {
	cases := map[string]SubTest{
		"account expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 100,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
					Type:       types.OpenDomain,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: 0,
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  AliceKey,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
	}

	RunTests(t, cases)
}
func Test_Common_handlerMsgReplaceAccountMetadata(t *testing.T) {
	cases := map[string]SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "does not exist",
					Name:   "",
					Owner:  AliceKey,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:  "test",
					Admin: BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "",
					NewMetadataURI: "",
					Owner:          nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "does not exist",
					Owner:  nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"signer is not owner of account": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"domain admin cannot replace metadata": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				GetConfigSetter(k.ConfigurationKeeper).SetConfig(ctx, configuration.Config{})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain: "test",
					Name:   "test",
					Owner:  BobKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"metadata size exceeded": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 2,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          AliceKey,
				})
				if !errors.Is(err, types.ErrMetadataSizeExceeded) {
					t.Fatalf("handlerMsgReplaceAccountMetadata() got error: %s", err)
				}
				_, err = handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "12",
					Owner:          AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountMetadata() got error: %s", err)
				}
			},
		},
		"success": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				setConfig := GetConfigSetter(k.ConfigurationKeeper).SetConfig
				setConfig(ctx, configuration.Config{
					MetadataSizeMax: 100,
				})
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// create domain
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Admin:      BobKey,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// create account
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					ValidUntil: utils.TimeToSeconds(time.Now().Add(1000 * time.Hour)),
					Owner:      AliceKey,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgReplaceAccountMetadata(ctx, k, &types.MsgReplaceAccountMetadata{
					Domain:         "test",
					Name:           "test",
					NewMetadataURI: "https://test.com",
					Owner:          AliceKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgReplaceAccountMetadata() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				expected := "https://test.com"
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("account not found")
				}
				if !reflect.DeepEqual(expected, account.MetadataURI) {
					t.Fatalf("handlerMsgSetMetadataURI expected: %+v, got %+v", expected, account.MetadataURI)
				}
			},
		},
	}
	// run tests
	RunTests(t, cases)
}

func Test_Closed_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]SubTest{
		"only domain admin can transfer": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// domain owned by alice
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// account owned by bob
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					Owner:      BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    AliceKey,
					NewOwner: CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
				_, err = handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    BobKey,
					NewOwner: CharlieKey,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					panic("unexpected account deletion")
				}
				if !account.Owner.Equals(CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", CharlieKey, account.Owner)
				}
			},
		},
		"domain admin can reset account content": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// domain owned by alice
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// account owned by bob
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        BobKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					MetadataURI:  "lol",
					Certificates: [][]byte{[]byte("test")},
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    AliceKey,
					NewOwner: CharlieKey,
					ToReset:  true,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					panic("unexpected account deletion")
				}
				if account.Resources != nil {
					panic("resources not deleted")
				}
				if account.Certificates != nil {
					panic("certificates not deleted")
				}
				if account.MetadataURI != "" {
					panic("metadata not deleted")
				}
			},
		},
	}

	RunTests(t, testCases)
}

func Test_Open_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]SubTest{
		"domain admin cannot transfer": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// domain owned by alice
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// account owned by bob
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					Owner:      BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    AliceKey,
					NewOwner: CharlieKey,
					ToReset:  false,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}

				_, err = handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    BobKey,
					NewOwner: CharlieKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					panic("unexpected account deletion")
				}
				if !account.Owner.Equals(CharlieKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", CharlieKey, account.Owner)
				}
			},
		},
		"domain admin cannot reset account content": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				// domain owned by alice
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				// account owned by bob
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        BobKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					MetadataURI:  "lol",
					Certificates: [][]byte{[]byte("test")},
					Resources: []*types.Resource{
						{
							URI:      "works",
							Resource: "works",
						},
					},
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// alice is domain owner and should transfer account owned by bob to alice
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    AliceKey,
					NewOwner: CharlieKey,
					ToReset:  true,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					t.Fatal("unexpected account deletion")
				}
				if account.Resources == nil {
					t.Fatal("resources deleted")
				}
				if account.Certificates == nil {
					t.Fatal("certificates deleted")
				}
				if account.MetadataURI == "" {
					t.Fatal("metadata not deleted")
				}
			},
		},
	}
	RunTests(t, testCases)
}

func Test_Common_handlerAccountTransfer(t *testing.T) {
	testCases := map[string]SubTest{
		"domain does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				// do nothing
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "does not exist",
					Name:     "does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrDomainDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrDomainDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"domain has expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "expired domain",
					Admin:      BobKey,
					ValidUntil: 0,
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "expired domain",
					Name:     "",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrDomainExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrDomainExpired, err)
				}
			},
			AfterTest: nil,
		},
		"account does not exist": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "this account does not exist",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrAccountDoesNotExist) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrAccountDoesNotExist, err)
				}
			},
			AfterTest: nil,
		},
		"account expired": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        BobKey,
					ValidUntil:   0,
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    nil,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrAccountExpired) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrAccountExpired, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has super user only domain admin can transfer accounts": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.ClosedDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        BobKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    BobKey,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"if domain has no super user then only account owner can transfer accounts": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain:       "test",
					Name:         utils.StrPtr("test"),
					Owner:        AliceKey,
					ValidUntil:   utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Resources:    nil,
					Certificates: nil,
					Broker:       nil,
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    BobKey,
					NewOwner: nil,
				})
				if !errors.Is(err, types.ErrUnauthorized) {
					t.Fatalf("handlerAccountTransfer() expected error: %s, got: %s", types.ErrUnauthorized, err)
				}
			},
			AfterTest: nil,
		},
		"success domain without superuser": {
			BeforeTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				domains := k.DomainStore(ctx)
				accounts := k.AccountStore(ctx)
				executor.NewDomain(ctx, types.Domain{
					Name:       "test",
					Admin:      BobKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
					Type:       types.OpenDomain,
					Broker:     nil,
				}).WithDomains(&domains).WithAccounts(&accounts).Create()
				executor.NewAccount(ctx, types.Account{
					Domain:     "test",
					Name:       utils.StrPtr("test"),
					Owner:      AliceKey,
					ValidUntil: utils.TimeToSeconds(ctx.BlockTime().Add(1000 * time.Hour)),
				}).WithAccounts(&accounts).Create()
			},
			Test: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				_, err := handlerMsgTransferAccount(ctx, k, &types.MsgTransferAccount{
					Domain:   "test",
					Name:     "test",
					Owner:    AliceKey,
					NewOwner: BobKey,
				})
				if err != nil {
					t.Fatalf("handlerMsgTransferAccount() got error: %s", err)
				}
			},
			AfterTest: func(t *testing.T, k Keeper, ctx sdk.Context, mocks *Mocks) {
				account := new(types.Account)
				if err := k.AccountStore(ctx).Read((&types.Account{Name: utils.StrPtr("test"), Domain: "test"}).PrimaryKey(), account); err != nil {
					panic("unexpected account deletion")
				}
				if account.Resources != nil {
					t.Fatalf("handlerAccountTransfer() account resources were not deleted")
				}
				if account.Certificates != nil {
					t.Fatalf("handlerAccountTransfer() account certificates were not deleted")
				}
				if !account.Owner.Equals(BobKey) {
					t.Fatalf("handlerAccounTransfer() expected new owner: %s, got: %s", BobKey, account.Owner)
				}
			},
		},
	}
	RunTests(t, testCases)
}
