package starname

import (
	"encoding/json"
	"testing"

	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
)

func TestExportGenesis(t *testing.T) {
	expected := `{"domains":[{"name":"test","admin":"cosmos1ze7y9qwdddejmy7jlw4cymqqlt2wh05ytm076d","valid_until":100,"type":"open"}],"accounts":[{"domain":"test","name":"","owner":"cosmos1ze7y9qwdddejmy7jlw4cymqqlt2wh05ytm076d","valid_until":100},{"domain":"test","name":"test","owner":"cosmos1ze7y9qwdddejmy7jlw4cymqqlt2wh05ytm076d","valid_until":100}]}`
	k, ctx, _ := keeper.NewTestKeeper(t, true)
	accounts := k.AccountStore(ctx)
	domains := k.DomainStore(ctx)
	keeper.NewDomainExecutor(ctx, types.Domain{
		Name:       "test",
		Admin:      keeper.AliceKey,
		ValidUntil: 100,
		Type:       types.OpenDomain,
		Broker:     nil,
	}).WithAccounts(&accounts).WithDomains(&domains).Create()
	keeper.NewAccountExecutor(ctx, types.Account{
		Domain:      "test",
		Name:        utils.StrPtr("test"),
		Owner:       keeper.AliceKey,
		ValidUntil:  100,
		MetadataURI: "",
	}).WithAccounts(&accounts).Create()
	b, err := json.Marshal(ExportGenesis(ctx, k))
	if err != nil {
		t.Fatal(err)
	}
	if string(b) != expected {
		t.Fatalf("unexpected genesis state:\nGot: %s\nWanted: %s", b, expected)
	}
}
