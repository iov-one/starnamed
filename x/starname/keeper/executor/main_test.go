package executor

import (
	"fmt"
	"os"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/mock"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"
)

var testCtx sdk.Context
var testKey = sdk.NewKVStoreKey("test")
var testCdc *codec.ProtoCodec
var _, testAddrs = utils.GeneratePrivKeyAddressPairs(2)
var aliceKey sdk.AccAddress = testAddrs[0]
var bobKey sdk.AccAddress = testAddrs[1]
var testConfig = &configuration.Config{
	Configurer:           "wasm1fjppc038udty5lquva2fc72967y4mchsu06slw", // bojack
	DomainRenewalPeriod:  10 * time.Second,
	AccountRenewalPeriod: 20 * time.Second,
}

var testKeeper keeper.Keeper
var testAccount = types.Account{
	Domain:     "a-super-domain",
	Name:       utils.StrPtr("a-super-account"),
	Owner:      aliceKey,
	ValidUntil: 10000,
	Resources: []*types.Resource{
		{
			URI:      "a-super-uri",
			Resource: "a-super-res",
		},
	},
	Certificates: [][]byte{[]byte("a-random-cert")},
	Broker:       nil,
	MetadataURI:  "metadata",
}

var testDomain = types.Domain{
	Name:       "a-super-domain",
	Admin:      bobKey,
	ValidUntil: 100,
	Type:       types.ClosedDomain,
}

func newTest() error {
	mockConfig := mock.NewConfiguration(nil, testConfig)
	// gen test store
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterInterface("cosmos.base.v1beta1.Msg",
		(*sdk.Msg)(nil),
		&types.MsgRegisterDomain{},
		&types.MsgTransferDomain{},
		&types.MsgTransferAccount{},
		&types.MsgAddAccountCertificates{},
		&types.MsgDeleteAccountCertificate{},
		&types.MsgDeleteAccount{},
		&types.MsgDeleteDomain{},
		&types.MsgRegisterAccount{},
		&types.MsgRenewDomain{},
		&types.MsgReplaceAccountResources{},
		&types.MsgReplaceAccountMetadata{},
	)
	testCdc = codec.NewProtoCodec(interfaceRegistry)
	mdb := db.NewMemDB()
	ms := store.NewCommitMultiStore(mdb)
	ms.MountStoreWithDB(testKey, sdk.StoreTypeIAVL, mdb)
	err := ms.LoadLatestVersion()
	if err != nil {
		return err
	}
	testCtx = sdk.NewContext(ms, tmproto.Header{Time: time.Now()}, true, log.NewNopLogger())
	testKeeper = keeper.NewKeeper(testCdc, testKey, mockConfig, nil, nil)
	testKeeper.AccountStore(testCtx).Create(&testAccount)
	testKeeper.DomainStore(testCtx).Create(&testDomain)
	testKeeper.AccountStore(testCtx).Create(&types.Account{
		Domain:      testDomain.Name,
		Name:        utils.StrPtr(types.EmptyAccountName),
		Owner:       testDomain.Admin,
		ValidUntil:  testDomain.ValidUntil,
		MetadataURI: "",
	})
	return nil
}

func TestMain(m *testing.M) {
	err := newTest()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
	os.Exit(m.Run())
}
