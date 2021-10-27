package mock

import sdk "github.com/cosmos/cosmos-sdk/types"

type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
	GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins
}

type supplyKeeper struct {
	sendCoinsFromAccountToModule func(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	getAllBalances               func(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}

func (s *supplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	return s.sendCoinsFromAccountToModule(ctx, addr, moduleName, coins)
}

func (s *supplyKeeper) GetAllBalances(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return s.getAllBalances(ctx, address)
}

type SupplyKeeperMock struct {
	s *supplyKeeper
}

func (s *SupplyKeeperMock) SetSendCoinsFromAccountToModule(f func(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error) {
	s.s.sendCoinsFromAccountToModule = f
}

func (s *SupplyKeeperMock) SetGetAllBalances(f func(ctx sdk.Context, address sdk.AccAddress) sdk.Coins) {
	s.s.getAllBalances = f
}

func (s *SupplyKeeperMock) Mock() SupplyKeeper {
	return s.s
}

func NewSupplyKeeper() *SupplyKeeperMock {
	mock := &SupplyKeeperMock{s: &supplyKeeper{}}
	// set default
	mock.SetSendCoinsFromAccountToModule(func(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
		return nil
	})

	mock.SetGetAllBalances(func(sdk.Context, sdk.AccAddress) sdk.Coins { return nil })
	return mock
}
