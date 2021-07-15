package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
)

type MsgWithFeePayer interface {
	sdk.Msg
	GetFeePayer() sdk.AccAddress
}

type TypeID uint64

type ObjectWithCustomFees interface {
	GetCreationFees() sdk.Dec
}

type ObjectWithTimeConstraint interface {
	ValidateDeadline(deadline uint64) error
}

type TransferableObject interface {
	codec.ProtoMarshaler

	GetObjectTypeID() TypeID
	GetCRUDObject() crud.Object

	IsOwnedBy(account sdk.AccAddress) (bool, error)
	Transfer(from sdk.AccAddress, to sdk.AccAddress) error
}

type StoreHolder interface {
	GetCRUDStore(sdk.Context) crud.Store
}

// Checks that SimpleStoreHolder implements the StoreHolder interface
var _ StoreHolder = SimpleStoreHolder{}

type SimpleStoreHolder struct {
	retrieveStore func(sdk.Context) crud.Store
}

func NewSimpleStoreHolder(storeRetriever func(sdk.Context) crud.Store) SimpleStoreHolder {
	return SimpleStoreHolder{retrieveStore: storeRetriever}
}

func (s SimpleStoreHolder) GetCRUDStore(ctx sdk.Context) crud.Store {
	return s.retrieveStore(ctx)
}
