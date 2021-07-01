package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
)

type TypeID uint64

type TransferableObject interface {
	codec.ProtoMarshaler

	GetType() TypeID
	//TODO: simplify this to transfer only the used data (getPk, load(store), update(store))
	GetObject() crud.Object

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