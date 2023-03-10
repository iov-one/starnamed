package keeper

import (
	"reflect"
	"testing"
	"time"

	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestAccount_AddCertificate(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	cert := []byte("a-cert")
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
	ex.AddCertificate(cert)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.Certificates, append(testAccount.Certificates, cert)) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Create(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	acc := testAccount
	acc.Domain = "some-random-domain"
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccountExecutor(testCtx, acc).WithAccounts(&as)
	ex.Create()
	got := new(types.Account)
	as.Read(acc.PrimaryKey(), got)
	if err := CompareAccounts(got, &acc); err != nil {
		t.Fatal(err)
	}
}

func TestAccount_DeleteCertificate(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
	ex.DeleteCertificate(0)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if len(got.Certificates) != 0 {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Renew(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	renewalPeriod := 20 * time.Second
	setConfig := GetConfigSetter(testKeeper.ConfigurationKeeper).SetConfig
	setConfig(testCtx, configuration.Config{
		AccountRenewalPeriod: renewalPeriod,
	})
	as := testKeeper.AccountStore(testCtx)
	conf := testKeeper.ConfigurationKeeper.GetConfiguration(testCtx)
	ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as).WithConfiguration(conf)
	ex.Renew()
	newAcc := new(types.Account)
	if err := as.Read(testAccount.PrimaryKey(), newAcc); err != nil {
		t.Fatal("account was deleted")
	}
	if newAcc.ValidUntil != testAccount.ValidUntil+int64(renewalPeriod.Seconds()) {
		t.Fatal("time mismatch")
	}
}

func TestAccount_ReplaceResources(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	newRes := []*types.Resource{{
		URI:      "uri",
		Resource: "res",
	}}
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
	ex.ReplaceResources(newRes)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.Resources, newRes) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_State(t *testing.T) {
	// TODO
}

func TestAccount_Transfer(t *testing.T) {
	t.Run("no-reset", func(t *testing.T) {
		testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
		as := testKeeper.AccountStore(testCtx)
		ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
		ex.Transfer(charlieAddr, false)
		got := new(types.Account)
		as.Read(testAccount.PrimaryKey(), got)
		if !got.Owner.Equals(charlieAddr) {
			t.Fatal("unexpected owner")
		}
		if !reflect.DeepEqual(got.Resources, testAccount.Resources) {
			t.Fatal("unexpected resources")
		}
		if !reflect.DeepEqual(got.MetadataURI, testAccount.MetadataURI) {
			t.Fatal("unexpected metadata")
		}
		if !reflect.DeepEqual(got.Certificates, testAccount.Certificates) {
			t.Fatal("unexpected certs")
		}
	})
	t.Run("with-reset", func(t *testing.T) {
		testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
		as := testKeeper.AccountStore(testCtx)
		ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
		ex.Transfer(bobAddr, true)
		got := new(types.Account)
		as.Read(testAccount.PrimaryKey(), got)
		if !got.Owner.Equals(bobAddr) {
			t.Fatal("owner mismatch")
		}
		if got.MetadataURI != "" || got.Resources != nil || got.Certificates != nil {
			t.Fatal("reset not performed")
		}
	})
}

func TestAccount_UpdateMetadata(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	newMeta := "a new meta"
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
	ex.UpdateMetadata(newMeta)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.MetadataURI, newMeta) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Delete(t *testing.T) {
	testKeeper, testCtx, _ := NewTestExecutorKeeper(t, false)
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccountExecutor(testCtx, testAccount).WithAccounts(&as)
	ex.Delete()
	got := new(types.Account)
	if err := as.Read(testAccount.PrimaryKey(), got); err == nil {
		t.Fatal("account was not deleted")
	}
}
