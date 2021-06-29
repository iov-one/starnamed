package keeper_test

import (
	"encoding/hex"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
	"github.com/gogo/protobuf/proto"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/log"
	"github.com/tendermint/tendermint/libs/rand"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"

	"github.com/iov-one/starnamed/mock"
	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/types"
)

// Assert that TestObject is a TransferableObject
var _ types.TransferableObject = &TestObject{}

const TypeID_TestObject = 1

type TestObject struct {
	types.TestObject
}

func (m *TestObject) PrimaryKey() []byte {
	return sdk.Uint64ToBigEndian(m.Id)
}

func (m *TestObject) SecondaryKeys() []crud.SecondaryKey {
	return make([]crud.SecondaryKey, 0)
}

func (m *TestObject) GetType() types.TypeID {
	return TypeID_TestObject
}

func (m *TestObject) GetObject() crud.Object {
	return m
}

func (m *TestObject) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	return m.Owner.Equals(account), nil
}

func (m *TestObject) Transfer(from sdk.AccAddress, to sdk.AccAddress) error {
	if !m.Owner.Equals(from) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "the object %v does not belong to %v", m.Id, from.String())
	}
	m.Owner = to
	return nil
}

type NotPossessedTestObject struct {
	TestObject
}

func (m *NotPossessedTestObject) IsOwnedBy(sdk.AccAddress) (bool, error) {
	return true, nil
}

func (m *NotPossessedTestObject) Transfer(sdk.AccAddress, sdk.AccAddress) error {
	return nil
}

type ErroredTestObject struct {
	NotPossessedTestObject
	nbTransferAllowed int64
}

func (m *ErroredTestObject) Transfer(sdk.AccAddress, sdk.AccAddress) error {
	if m.nbTransferAllowed > 0 {
		m.nbTransferAllowed--
		return nil
	}
	return fmt.Errorf("this test object cannot be transfered")
}

func NewEscrowGenerator(now uint64) *EscrowGenerator {
	return &EscrowGenerator{now: now}
}

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

func (gen *EscrowGenerator) NewTestObject(seller sdk.AccAddress) *TestObject {
	return &TestObject{types.TestObject{
		Id:    gen.NextObjectID(),
		Owner: append([]byte(nil), seller...),
	}}
}

func (gen *EscrowGenerator) NewNotPossessedTestObject() *NotPossessedTestObject {
	// this seller address does not matter and is required for the crud store to work properly
	return &NotPossessedTestObject{*gen.NewTestObject(gen.NewAccAddress())}
}

func (gen *EscrowGenerator) NewErroredTestObject(nbTransferAllowed int64) *ErroredTestObject {
	return &ErroredTestObject{*gen.NewNotPossessedTestObject(), nbTransferAllowed}
}

func (gen *EscrowGenerator) NewTestEscrow(seller sdk.AccAddress, price sdk.Coins, deadline uint64) (types.Escrow, *TestObject) {
	obj := gen.NewTestObject(seller)
	packedObj, err := cdctypes.NewAnyWithValue(obj)
	if err != nil {
		panic(err)
	}
	return types.Escrow{
		Id:       gen.NextID(),
		Seller:   seller.String(),
		Object:   packedObj,
		Price:    price,
		Deadline: deadline,
	}, obj
}

func (gen *EscrowGenerator) NewRandomTestEscrow() (types.Escrow, *TestObject) {
	seller := gen.NewAccAddress()
	coins := sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(rand.Int64()+1)))
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
		NextEscrowId:  0,
	}
}

func (gen *EscrowGenerator) NowAfter(date uint64) uint64 {
	return gen.now + date
}

var TestTimeNow = time.Unix(2000, 0)

func NewTestCodec() *codec.ProtoCodec {
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	types.RegisterInterfaces(interfaceRegistry)
	interfaceRegistry.RegisterImplementations((*types.TransferableObject)(nil),
		&TestObject{},
	)
	//Register the test object implementation
	cdc := codec.NewProtoCodec(interfaceRegistry)
	return cdc
}

func NewTestKeeper(coinHolders []sdk.AccAddress) (keeper.Keeper, sdk.Context, crud.Store, map[string]sdk.Coins) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	escrowStoreKey := sdk.NewKVStoreKey(types.StoreKey) // domain module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(escrowStoreKey, sdk.StoreTypeIAVL, mdb) // mount domain module
	// test no errors
	if err := ms.LoadLatestVersion(); err != nil {
		panic(err)
	}
	// create a crud store for the crud objects
	crudStore := crud.NewStore(cdc, ms.GetKVStore(escrowStoreKey), []byte{1})

	// create mock supply keeper with money on accounts
	bankMocker := mock.NewSupplyKeeper()
	balances := make(map[string]sdk.Coins)
	if coinHolders != nil {
		for _, holder := range coinHolders {
			balances[holder.String()] = sdk.NewCoins(
				sdk.NewCoin("tiov", sdk.NewInt(1000000)),
				sdk.NewCoin("tiov2", sdk.NewInt(1000000)),
			)
		}
		bankMocker.WithDefaultsBalances(balances)
	}
	// Create mock auth keeper
	authMocker := mock.NewAccountKeeper()
	// create context
	ctx := sdk.NewContext(ms, tmproto.Header{Time: TestTimeNow}, true, log.NewNopLogger())
	// Create param subspace
	paramsSubspace := paramstypes.NewSubspace(cdc, nil, escrowStoreKey, sdk.NewKVStoreKey("t"+types.StoreKey), types.ModuleName)

	// register blocked addresses
	blockedAddr := make(map[string]bool)
	blockedAddr[authtypes.NewModuleAddress(types.ModuleName).String()] = true

	k := keeper.NewKeeper(cdc, escrowStoreKey, paramsSubspace, authMocker.Mock(), bankMocker.Mock(), blockedAddr)
	k.AddStore(TypeID_TestObject, crudStore)
	k.SetLastBlockTime(ctx, uint64(ctx.BlockTime().Unix()))
	return k, ctx, crudStore, balances
}

func MustPackToAny(val proto.Message) *cdctypes.Any {
	any, err := cdctypes.NewAnyWithValue(val)
	if err != nil {
		panic(sdkerrors.Wrap(err, "error while converting a value to an any"))
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

func EvaluateTest(t *testing.T, name string, test func() error) {
	shouldPanic := strings.Contains(name, "panic")
	shouldFail := strings.Contains(name, "invalid")

	testFuncWrapped := func() {
		err := test()
		CheckError(t, name, shouldFail, err)
	}

	if shouldPanic {
		assert.Panics(t, testFuncWrapped, name)
	} else {
		assert.NotPanics(t, testFuncWrapped, name)
	}
}

type BaseKeeperSuite struct {
	suite.Suite
	keeper    keeper.Keeper
	msgServer types.MsgServer
	ctx       sdk.Context
	generator *EscrowGenerator
	store     crud.Store
}

func (s *BaseKeeperSuite) Setup() {
	s.keeper, s.ctx, s.store, _ = NewTestKeeper(nil)
	s.keeper.ImportNextID(s.ctx, 1)
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.generator = NewEscrowGenerator(uint64(s.ctx.BlockTime().Unix()))
}
