package escrow_test

import (
	"encoding/hex"

	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/tendermint/tendermint/libs/rand"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// Assert that TestObject is a TransferableObject
var _ types.TransferableObject = &TestObject{}

const TypeID_TestObject = 1

type TestObject struct {
	types.TestObject
}

func (m *TestObject) PrimaryKey() []byte {
	return m.Owner
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

type EscrowGenerator struct {
	NextId       uint64
	NextObjectId uint64
	Now          uint64
}

func (gen *EscrowGenerator) NextID() string {
	id := hex.EncodeToString(sdk.Uint64ToBigEndian(gen.NextId))
	gen.NextId++
	return id
}

func (gen *EscrowGenerator) NextObjectID() uint64 {
	id := gen.NextObjectId
	gen.NextObjectId++
	return id
}

func (gen *EscrowGenerator) NewTestEscrow(seller sdk.AccAddress, price sdk.Coins, deadline uint64) (types.Escrow, *TestObject) {
	obj := &TestObject{types.TestObject{
		Id:    gen.NextObjectID(),
		Owner: seller,
	}}
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

func (gen *EscrowGenerator) NewEmptyEscrowGenesis() *types.GenesisState {
	return &types.GenesisState{
		Escrows:       []types.Escrow{},
		LastBlockTime: 0,
		NextEscrowId:  0,
	}
}

func (gen *EscrowGenerator) NewEscrowGenesis(numEscrows int) *types.GenesisState {
	escrows := make([]types.Escrow, 0)

	now := gen.Now
	for i := 0; i < numEscrows; i++ {
		seller := rand.Bytes(20)
		coins := sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(rand.Int64()+1)))
		escrow, _ := gen.NewTestEscrow(seller, coins, now+uint64(i*500))
		escrows = append(escrows, escrow)
	}

	return &types.GenesisState{
		Escrows:       escrows,
		LastBlockTime: now,
		NextEscrowId:  0,
	}
}
