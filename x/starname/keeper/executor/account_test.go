package executor

import (
	"reflect"
	"testing"

	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestAccount_AddCertificate(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	cert := []byte("a-cert")
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.AddCertificate(cert)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.Certificates, append(testAccount.Certificates, cert)) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Create(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	acc := testAccount
	acc.Domain = "some-random-domain"
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.Create()
	got := new(types.Account)
	as.Read(acc.PrimaryKey(), got)
	if !reflect.DeepEqual(*got, acc) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_DeleteCertificate(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.DeleteCertificate(0)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if len(got.Certificates) != 0 {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Renew(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.Renew()
	newAcc := new(types.Account)
	if err := as.Read(testAccount.PrimaryKey(), newAcc); err != nil {
		t.Fatal("account was deleted")
	}
	if newAcc.ValidUntil != testAccount.ValidUntil+int64(testConfig.AccountRenewalPeriod.Seconds()) {
		t.Fatal("time mismatch")
	}
}

func TestAccount_ReplaceResources(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	newRes := []*types.Resource{{
		URI:      "uri",
		Resource: "res",
	}}
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.ReplaceResources(newRes)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.Resources, newRes) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_State(t *testing.T) {

}

func TestAccount_Transfer(t *testing.T) {
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	t.Run("no-reset", func(t *testing.T) {
		// dmjp testCtx, _ := testCtx.CacheContext()

		ex.Transfer(keeper.CharlieKey, false)
		got := new(types.Account)
		as.Read(testAccount.PrimaryKey(), got)
		if !got.Owner.Equals(keeper.CharlieKey) {
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
		// dmjp testCtx, _ := testCtx.CacheContext()

		ex.Transfer(keeper.BobKey, true)
		got := new(types.Account)
		as.Read(testAccount.PrimaryKey(), got)
		if !got.Owner.Equals(keeper.BobKey) {
			t.Fatal("owner mismatch")
		}
		if got.MetadataURI != "" || got.Resources != nil || got.Certificates != nil {
			t.Fatal("reset not performed")
		}
	})
}

func TestAccount_UpdateMetadata(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	newMeta := "a new meta"
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.UpdateMetadata(newMeta)
	got := new(types.Account)
	as.Read(testAccount.PrimaryKey(), got)
	if !reflect.DeepEqual(got.MetadataURI, newMeta) {
		t.Fatal("unexpected result")
	}
}

func TestAccount_Delete(t *testing.T) {
	testCtx, _ := testCtx.CacheContext()
	as := testKeeper.AccountStore(testCtx)
	ex := NewAccount(testCtx, testAccount).WithStore(&as)
	ex.Delete()
	got := new(types.Account)
	if err := as.Read(testAccount.PrimaryKey(), got); err == nil {
		t.Fatal("account was not deleted")
	}
}
