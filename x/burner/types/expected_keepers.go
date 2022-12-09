package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth/types"
)

type SupplyKeeper interface {
	GetAllBalances(sdk.Context, sdk.AccAddress) sdk.Coins
	BurnCoins(sdk.Context, string, sdk.Coins) error
}

type AccountKeeper interface {
	GetModuleAddress(string) sdk.AccAddress
	GetModuleAccount(sdk.Context, string) types.ModuleAccountI
}
