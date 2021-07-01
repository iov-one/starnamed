package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/configuration"
)

// BankKeeper defines the expected bank keeper (noalias)
type BankKeeper interface {
	SendCoinsFromModuleToAccount(ctx sdk.Context, senderModule string, recipientAddr sdk.AccAddress, amt sdk.Coins) error
	SendCoinsFromAccountToModule(ctx sdk.Context, senderAddr sdk.AccAddress, recipientModule string, amt sdk.Coins) error
}

// AccountKeeper defines the expected account keeper (noalias)
type AccountKeeper interface {
	GetModuleAddress(name string) sdk.AccAddress
}

// FIXME this should not be a dependency of this module

// ConfigurationKeeper defines the expected configuration keeper
type ConfigurationKeeper interface {
	// GetFees gets the fees
	GetFees(ctx sdk.Context) *configuration.Fees
	// GetConfiguration returns the configuration
	GetConfiguration(ctx sdk.Context) configuration.Config
}
