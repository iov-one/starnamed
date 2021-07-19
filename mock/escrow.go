package mock

import (
	crud "github.com/iov-one/cosmos-sdk-crud"

	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"
)

type EscrowKeeper interface {
	AddStore(id escrowtypes.TypeID, store crud.Store)
	AddStoreHolder(id escrowtypes.TypeID, holder escrowtypes.StoreHolder)
	RegisterCustomData(id escrowtypes.TypeID, data escrowtypes.CustomData)
}

type escrowKeeper struct {
}

func (s escrowKeeper) AddStore(escrowtypes.TypeID, crud.Store) {
}

func (s escrowKeeper) AddStoreHolder(escrowtypes.TypeID, escrowtypes.StoreHolder) {
}

func (s escrowKeeper) RegisterCustomData(id escrowtypes.TypeID, data escrowtypes.CustomData) {

}

type EscrowKeeperMock struct {
	e *escrowKeeper
}

func (e *EscrowKeeperMock) Mock() EscrowKeeper {
	return e.e
}

func NewEscrowKeeper() *EscrowKeeperMock {
	return &EscrowKeeperMock{e: &escrowKeeper{}}
}
