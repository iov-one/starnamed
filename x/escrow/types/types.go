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
	CanTransfer(from sdk.AccAddress, to sdk.AccAddress) error
	Transfer(from sdk.AccAddress, to sdk.AccAddress) error
}

type StoreHolder interface {
	GetCRUDStore() crud.Store
}
