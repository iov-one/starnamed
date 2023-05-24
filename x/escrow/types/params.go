package types

import (
	"fmt"

	paramtypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Default parameter values
const (
	DefaultModuleEnabled  bool    = true
	DefaultCustomTokenFee float64 = 0.5
)

// Parameter keys
var (
	KeyModuleEnabled       = []byte("ModuleEnabled")
	KeyAllowedCustomTokens = []byte("AllowedCustomTokens")
	KeyCustomTokenFee      = []byte("CustomTokenFee")
)

var _ paramtypes.ParamSet = &Params{}

// ParamKeyTable for auth module
func ParamKeyTable() paramtypes.KeyTable {
	return paramtypes.NewKeyTable().RegisterParamSet(&Params{})
}

// ParamSetPairs implements the ParamSet interface and returns all the key/value pairs
// pairs of auth module's parameters.
// nolint
func (p *Params) ParamSetPairs() paramtypes.ParamSetPairs {
	return paramtypes.ParamSetPairs{
		paramtypes.NewParamSetPair(KeyModuleEnabled, &p.ModuleEnabled, validateIsBool),
	}
}

// DefaultParams returns a default set of parameters.
func DefaultParams() Params {
	return Params{
		ModuleEnabled: DefaultModuleEnabled,
	}
}

func validateIsBool(i interface{}) error {
	_, ok := i.(bool)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T, expected boolean", i)
	}
	return nil
}

func validateIsSliceOfStrings(i interface{}) error {
	_, ok := i.([]string)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T, expected []string", i)
	}
	return nil
}

func validateIfIsPositiveFloat64(i interface{}) error {
	f, ok := i.(float64)
	if !ok {
		return fmt.Errorf("invalid parameter type: %T, expected float64", i)
	}
	// Check if the float is positive
	if f < 0 {
		return fmt.Errorf("invalid parameter value: %f, expected positive float64", f)
	}
	return nil
}

// Validate checks that the parameters have valid values.
func (p Params) Validate() error {
	if err := validateIsBool(p.ModuleEnabled); err != nil {
		return err
	}
	return nil
}
