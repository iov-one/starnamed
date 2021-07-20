package mock

import (
	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"
)

type EscrowKeeper interface {
	RegisterCustomData(id escrowtypes.TypeID, data escrowtypes.CustomData)
}

type escrowKeeper struct {
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
