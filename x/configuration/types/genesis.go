package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// NewGenesisState is GenesisState constructor
func NewGenesisState(conf Config, fees *Fees) GenesisState {
	return GenesisState{
		Config: conf,
		Fees:   *fees,
	}
}

// ValidateGenesis makes sure that the genesis state is valid
func ValidateGenesis(data GenesisState) error {
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
func DefaultGenesisState() GenesisState {
	// set default configs
	config := Config{
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
		EscrowBroker:           "star1nrnx8mft8mks3l2akduxdjlf8rwqs8r9l36a78", 					 // to IOV msig account
		EscrowMaxPeriod:        7890000 * 1e9,                                 					 // 3 months
	}
	// set fees
	// add domain module fees
	feeCoinDenom := "tiov" // set coin denom used for fees
	// generate new fees
	fees := NewFees()
	// set default fees
	fees.SetDefaults(feeCoinDenom)
	// return genesis
	return GenesisState{
		Config: config,
		Fees:   *fees,
	}
}
