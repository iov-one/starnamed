package test

import (
	"encoding/hex"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/gogo/protobuf/proto"
	crud "github.com/iov-one/cosmos-sdk-crud"
	crudtypes "github.com/iov-one/cosmos-sdk-crud/types"
	"github.com/stretchr/testify/assert"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/iov-one/starnamed/app"
	"github.com/iov-one/starnamed/mock"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type EscrowGenerator struct {
	nextId       uint64
	nextObjectId uint64
	now          uint64
}

func (gen *EscrowGenerator) NextID() string {
	id := hex.EncodeToString(sdk.Uint64ToBigEndian(gen.nextId))
	gen.nextId++
	return id
}

func (gen *EscrowGenerator) NextObjectID() uint64 {
	id := gen.nextObjectId
	gen.nextObjectId++
	return id
}

func (gen *EscrowGenerator) NewTestObject(seller sdk.AccAddress) *types.TestObject {
	return &types.TestObject{
		Id:                  gen.NextObjectID(),
		Owner:               append([]byte(nil), seller...),
		NumAllowedTransfers: -1,
	}
}

func (gen *EscrowGenerator) NewNotPossessedTestObject() *types.TestObject {
	return gen.NewTestObject(nil)
}

func (gen *EscrowGenerator) NewErroredTestObject(nbTransferAllowed int64) *types.TestObject {
	testObj := gen.NewNotPossessedTestObject()
	testObj.NumAllowedTransfers = nbTransferAllowed
	return testObj
}

func (gen *EscrowGenerator) NewTimeConstrainedObject(seller sdk.AccAddress, expiration uint64) *types.TestTimeConstrainedObject {
	return &types.TestTimeConstrainedObject{
		Id:         gen.NextObjectID(),
		Owner:      seller,
		Expiration: expiration,
	}
}

func (gen *EscrowGenerator) NewTestEscrow(seller sdk.AccAddress, price sdk.Coins, deadline uint64) (types.Escrow, *types.TestObject) {
	obj := gen.NewTestObject(seller)
	packedObj, err := cdctypes.NewAnyWithValue(obj)
	if err != nil {
		panic(err)
	}
	return types.Escrow{
		Id:               gen.NextID(),
		Seller:           seller.String(),
		Object:           packedObj,
		Price:            price,
		Deadline:         deadline,
		BrokerCommission: sdk.ZeroDec(),
		BrokerAddress:    gen.NewAccAddress().String(),
	}, obj
}

func (gen *EscrowGenerator) NewRandomTestEscrow() (types.Escrow, *types.TestObject) {
	seller := gen.NewAccAddress()
	coins := sdk.NewCoins(sdk.NewCoin(Denom, sdk.NewInt(int64(rand.Uint64()>>1)+1)))
	return gen.NewTestEscrow(seller, coins, gen.now+1+uint64(rand.Uint32()%5000))
}
func (gen *EscrowGenerator) NewAccAddress() sdk.AccAddress {
	return rand.Bytes(20)
}

func (gen *EscrowGenerator) NewEmptyEscrowGenesis() *types.GenesisState {
	return &types.GenesisState{
		Escrows:       []types.Escrow{},
		LastBlockTime: 0,
		NextEscrowId:  0,
	}
}

func (gen *EscrowGenerator) NewEscrowGenesis(numEscrows int) *types.GenesisState {
	escrows := make([]types.Escrow, 0)

	now := gen.now
	for i := 0; i < numEscrows; i++ {
		escrow, _ := gen.NewRandomTestEscrow()
		escrows = append(escrows, escrow)
	}

	return &types.GenesisState{
		Escrows:       escrows,
		LastBlockTime: now,
		NextEscrowId:  gen.GetNextId(),
	}
}

func (gen *EscrowGenerator) NowAfter(date uint64) uint64 {
	return gen.now + date
}

func (gen *EscrowGenerator) GetNextId() uint64 {
	return gen.nextId
}

var TimeNow = time.Unix(2000, 0)
var Denom = "tiov"
var DenomAux = "tiov2"

func NewTestCodec() (*codec.ProtoCodec, *codec.LegacyAmino) {
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	types.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*types.TransferableObject)(nil),
		&types.TestObject{},
		&types.TestTimeConstrainedObject{},
	)
	//Register the test object implementation
	cdc := codec.NewProtoCodec(interfaceRegistry)

	legacyCdc := codec.NewLegacyAmino()
	types.RegisterLegacyAminoCodec(legacyCdc)
	return cdc, legacyCdc
}

func SetConfig() {

	config := sdk.GetConfig()
	config.SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)
	config.SetBech32PrefixForValidator(app.Bech32PrefixValAddr, app.Bech32PrefixValPub)
	config.SetBech32PrefixForConsensusNode(app.Bech32PrefixConsAddr, app.Bech32PrefixConsPub)
}

func NewTestKeeper(coinHolders []sdk.AccAddress, isModuleEnabled bool) (keeper.Keeper, sdk.Context, crud.Store, map[string]sdk.Coins, sdk.StoreKey, types.ConfigurationKeeper) {
	cdc, legacyCdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	escrowStoreKey := sdk.NewKVStoreKey(types.StoreKey)                // domain module store key
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	paramStoreKey := sdk.NewKVStoreKey(paramstypes.StoreKey)           // SDK parameter store key
	tParamStoreKey := sdk.NewKVStoreKey("t" + paramstypes.StoreKey)    // SDK parameter transient store key

	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(escrowStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount config module
	ms.MountStoreWithDB(paramStoreKey, sdk.StoreTypeIAVL, mdb)         // mount params stores
	ms.MountStoreWithDB(tParamStoreKey, sdk.StoreTypeIAVL, mdb)

	// test no errors
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}
	// create a crud store for the crud objects
	crudStore := crudtypes.NewStore(cdc, ms.GetKVStore(escrowStoreKey), []byte{4})

	// create mock supply keeper with money on accounts
	bankMocker := mock.NewSupplyKeeper()
	balances := make(map[string]sdk.Coins)
	if coinHolders != nil {
		for _, holder := range coinHolders {
			balances[holder.String()] = sdk.NewCoins(
				sdk.NewCoin(Denom, sdk.NewInt(1000000)),
				sdk.NewCoin(DenomAux, sdk.NewInt(1000000)),
			)
		}
		bankMocker.WithDefaultsBalances(balances)
	}
	// Create mock auth keeper
	authMocker := mock.NewAccountKeeper()
	// Create config keeper
	configKeeper := configuration.NewKeeper(cdc, configurationStoreKey, nil)
	// create context
	ctx := sdk.NewContext(ms, tmproto.Header{Time: TimeNow}, true, log.NewNopLogger())
	// Create param subspace
	paramsSubspace := paramstypes.NewSubspace(cdc, legacyCdc, paramStoreKey, tParamStoreKey, types.ModuleName)
	if isModuleEnabled {
		paramsSubspace = paramsSubspace.WithKeyTable(types.ParamKeyTable())
		paramsSubspace.SetParamSet(ctx, &types.Params{ModuleEnabled: true})
	}

	// Set default fees
	defaultFees := configuration.NewFees()
	defaultFees.SetDefaults(Denom)
	configKeeper.SetFees(ctx, defaultFees)

	defaultConfig := configuration.DefaultGenesisState().Config
	configKeeper.SetConfig(ctx, defaultConfig)

	// register blocked addresses
	blockedAddr := make(map[string]bool)
	blockedAddr[authtypes.NewModuleAddress(types.ModuleName).String()] = true

	k := keeper.NewKeeper(cdc, escrowStoreKey, paramsSubspace, authMocker.Mock(), bankMocker.Mock(), configKeeper, blockedAddr)
	k.RegisterCustomData(types.TypeIDTestObject, crudStore)
	k.RegisterCustomData(types.TypeIDTestTimeConstrainedObject, wrapStoreForTimeContrainedObjects(crudStore))
	k.SetLastBlockTime(ctx, uint64(ctx.BlockTime().Unix()))
	return k, ctx, crudStore, balances, escrowStoreKey, configKeeper
}

type storeWrapper struct {
	store crud.Store
}

func wrapStoreForTimeContrainedObjects(store crud.Store) *storeWrapper {
	return &storeWrapper{store}
}

func (s storeWrapper) GetCrudStore() crud.Store {
	return s.store
}
func (s storeWrapper) GetDeadlineOrDefault(_ sdk.Context, _ types.TransferableObject, def uint64) uint64 {
	return def
}

func MustPackToAny(val proto.Message) *cdctypes.Any {
	any, err := cdctypes.NewAnyWithValue(val)
	if err != nil {
		panic(errors.Wrap(err, "error while converting a value to an any"))
	}
	return any
}

func CheckError(t *testing.T, name string, shouldFail bool, err error) {
	if (err != nil) != shouldFail {
		var verb string
		if shouldFail {
			verb = "should"
		} else {
			verb = "shouldn't"
		}
		t.Fatalf("The test %v %v have failed : %s", name, verb, err)
	}
}

// EvaluateTest executes the given test and checks that the execution is matching with the expected behaviour.
// The name of the test describe the expected behaviour: if the name contains "panic" then the test should panic,
// otherwise it should not. If the name contains "invalid" then the test function should return an error, otherwise it should not.
func EvaluateTest(t *testing.T, name string, test func(t *testing.T) error) {
	shouldPanic := strings.Contains(name, "panic")
	shouldFail := strings.Contains(name, "invalid")

	getWrappedTestFunc := func(t *testing.T) func() {
		return func() {
			err := test(t)
			CheckError(t, name, shouldFail, err)
		}
	}

	t.Run(name, func(t *testing.T) {
		if shouldPanic {
			assert.Panics(t, getWrappedTestFunc(t), name)
		} else {
			assert.NotPanics(t, getWrappedTestFunc(t), name)
		}
	})
}

func NewEscrowGenerator(now uint64) *EscrowGenerator {
	return &EscrowGenerator{now: now}
}

func DeleteEscrow(ctx sdk.Context, storeKey sdk.StoreKey, id string) {
	cdc, _ := NewTestCodec()
	str := crudtypes.NewStore(cdc, ctx.KVStore(storeKey), keeper.EscrowStoreKey)
	if err := str.Delete(types.GetEscrowKey(id)); err != nil {
		panic(errors.Wrap(err, "error while deleting an escrow in a test"))
	}

}
