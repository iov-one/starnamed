package mock

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
)

type SupplyKeeper interface {
	SendCoinsFromAccountToModule(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	SendCoinsFromModuleToAccount(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
	GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins
}

type supplyKeeper struct {
	sendCoinsFromAccountToModule func(sdk.Context, sdk.AccAddress, string, sdk.Coins) error
	sendCoinsFromModuleToAccount func(sdk.Context, string, sdk.AccAddress, sdk.Coins) error
	getAllBalances               func(ctx sdk.Context, address sdk.AccAddress) sdk.Coins
}

func (s *supplyKeeper) SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
	return s.sendCoinsFromAccountToModule(ctx, addr, moduleName, coins)
}

func (s *supplyKeeper) SendCoinsFromModuleToAccount(ctx sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
	return s.sendCoinsFromModuleToAccount(ctx, moduleName, addr, coins)
}

func (s *supplyKeeper) GetAllBalances(ctx sdk.Context, address sdk.AccAddress) sdk.Coins {
	return s.getAllBalances(ctx, address)
}

type SupplyKeeperMock struct {
	s *supplyKeeper
}

func (s *SupplyKeeperMock) SetSendCoinsFromAccountToModule(f func(sdk.Context, sdk.AccAddress, string, sdk.Coins) error) {
	s.s.sendCoinsFromAccountToModule = f
}

func (s *SupplyKeeperMock) SetSendCoinsFromModuleToAccount(f func(sdk.Context, string, sdk.AccAddress, sdk.Coins) error) {
	s.s.sendCoinsFromModuleToAccount = f
}
func (s *SupplyKeeperMock) SetGetAllBalances(f func(ctx sdk.Context, address sdk.AccAddress) sdk.Coins) {
	s.s.getAllBalances = f
}

func (s *SupplyKeeperMock) Mock() SupplyKeeper {
	return s.s
}

func (s *SupplyKeeperMock) WithDefaultsBalances(balances map[string]sdk.Coins) *SupplyKeeperMock {
	send := func(src, dest sdk.AccAddress, coins sdk.Coins) error {
		if !balances[src.String()].IsAllGTE(coins) {
			return sdkerrors.ErrInsufficientFunds
		}
		balanceSrc := balances[src.String()].Sub(coins)
		balances[src.String()] = balanceSrc
		balanceDest := balances[dest.String()].Add(coins...)
		balances[dest.String()] = balanceDest
		return nil
	}

	// set default
	s.SetSendCoinsFromAccountToModule(func(_ sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
		return send(addr, authtypes.NewModuleAddress(moduleName), coins)
	})

	s.SetSendCoinsFromModuleToAccount(func(_ sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
		return send(authtypes.NewModuleAddress(moduleName), addr, coins)
	})

	s.SetGetAllBalances(func(_ sdk.Context, addr sdk.AccAddress) sdk.Coins {
		return balances[addr.String()]
	})
	return s
}

func NewSupplyKeeper() *SupplyKeeperMock {
	mock := &SupplyKeeperMock{s: &supplyKeeper{}}
	// set no-ops
	mock.SetSendCoinsFromAccountToModule(func(_ sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error {
		return nil
	})

	mock.SetSendCoinsFromModuleToAccount(func(_ sdk.Context, moduleName string, addr sdk.AccAddress, coins sdk.Coins) error {
		return nil
	})

	mock.SetGetAllBalances(func(sdk.Context, sdk.AccAddress) sdk.Coins { return nil })
	return mock
}
