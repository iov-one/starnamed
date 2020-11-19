package configuration

import (
	"github.com/CosmWasm/wasmd/x/configuration/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	ptypes "github.com/gogo/protobuf/types"
)

// NewGenesisState is GenesisState constructor
func NewGenesisState(conf types.Config, fees *types.Fees) types.GenesisState {
	return types.GenesisState{
		Config: conf,
		Fees:   *fees,
	}
}

// ValidateGenesis makes sure that the genesis state is valid
func ValidateGenesis(data types.GenesisState) error {
	conf := data.Config
	if err := conf.Validate(); err != nil {
		return err
	}
	if err := data.Fees.Validate(); err != nil {
		return err
	}
	return nil
}

// DefaultGenesisState returns the default genesis state
// TODO this needs to be updated, although it will be imported from iovns chain
func DefaultGenesisState() types.GenesisState {
	// get owner
	owner, err := sdk.AccAddressFromBech32("star1kxqay5tndu3w08ps5c27pkrksnqqts0zxeprzx")
	if err != nil {
		panic("invalid default owner provided")
	}
	// set default configs
	config := types.Config{
		Configurer:             owner,
		ValidDomainName:        "^[-_a-z0-9]{4,16}$",
		ValidAccountName:       "[-_\\.a-z0-9]{1,64}$",
		ValidURI:               "[-a-z0-9A-Z:]+$",
		ValidResource:          "^[a-z0-9A-Z]+$",
		DomainRenewalPeriod:    ptypes.Duration{Seconds: 180000000000},
		DomainRenewalCountMax:  2,
		DomainGracePeriod:      ptypes.Duration{Seconds: 60000000000},
		AccountRenewalPeriod:   ptypes.Duration{Seconds: 180000000000},
		AccountRenewalCountMax: 3,
		AccountGracePeriod:     ptypes.Duration{Seconds: 60000000000},
		ResourcesMax:           3,
		CertificateSizeMax:     10000,
		CertificateCountMax:    3,
		MetadataSizeMax:        86400,
	}
	// set fees
	// add domain module fees
	feeCoinDenom := "tiov" // set coin denom used for fees
	// generate new fees
	fees := types.NewFees()
	// set default fees
	fees.SetDefaults(feeCoinDenom)
	// return genesis
	return types.GenesisState{
		Config: config,
		Fees:   *fees,
	}
}

// InitGenesis sets the initial state of the configuration module
func InitGenesis(ctx sdk.Context, k Keeper, data types.GenesisState) {
	k.SetConfig(ctx, data.Config)
	k.SetFees(ctx, &data.Fees)
}

// ExportGenesis saves the state of the configuration module
func ExportGenesis(ctx sdk.Context, k Keeper) *types.GenesisState {
	return &types.GenesisState{
		Config: k.GetConfiguration(ctx),
		Fees:   *k.GetFees(ctx),
	}
}
