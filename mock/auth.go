package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

type accountKeeper struct {
}

func (a accountKeeper) GetModuleAddress(name string) sdk.AccAddress {

	return authtypes.NewModuleAddress(name)
}

type AccountKeeperMock struct {
	a *accountKeeper
}

func (a *AccountKeeperMock) Mock() AccountKeeper {
	return a.a
}

func NewAccountKeeper() *AccountKeeperMock {
	return &AccountKeeperMock{a: &accountKeeper{}}
}
