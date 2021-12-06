package configuration

import (
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/configuration/types"
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
	// set default configs
	config := types.Config{
		Configurer:             "star1d3lhm5vtta78cm7c7ytzqh7z5pcgktmautntqv", // msig1
		ValidDomainName:        "^[-_a-z0-9]{4,16}$",
		ValidAccountName:       "^[-_\\.a-z0-9]{1,64}$",
		ValidURI:               "^[-a-z0-9A-Z:]+$",
		ValidResource:          "^[a-z0-9A-Z]+$",
		DomainRenewalPeriod:    31557600 * 1e9,
		DomainRenewalCountMax:  2,
		DomainGracePeriod:      2592000 * 1e9,
		AccountRenewalPeriod:   31557600 * 1e9,
		AccountRenewalCountMax: 3,
		AccountGracePeriod:     2592000 * 1e9,
		ResourcesMax:           3,
		CertificateSizeMax:     10000,
		CertificateCountMax:    3,
		MetadataSizeMax:        86400,
		EscrowCommission:       sdk.NewDecFromInt(sdk.NewInt(1)).QuoInt(sdk.NewInt(100)), // 1%
		EscrowBroker:           "star1d3lhm5vtta78cm7c7ytzqh7z5pcgktmautntqv",            // msig1
		EscrowMaxPeriod:        7890000 * 1e9,                                            // 3 months
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
