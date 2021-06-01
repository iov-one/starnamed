package types

import (
	fmt "fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewGenesisState constructs a new GenesisState instance
func NewGenesisState(escrows []Escrow) *GenesisState {
	return &GenesisState{
		Escrows: escrows,
	}
}

// DefaultGenesisState gets the raw genesis message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Escrows: []Escrow{},
	}
}

// ValidateGenesis validates the provided genesis state to ensure the expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	ids := map[string]bool{}
	for _, escrow := range data.Escrows {
		if ids[escrow.Id] {
			return fmt.Errorf("found duplicate escrow ID %s", escrow.Id)
		}

		if escrow.State != Open {
			return sdkerrors.Wrap(ErrEscrowNotOpen, escrow.Id)
		}

		if err := escrow.Validate(); err != nil {
			return err
		}

		ids[escrow.Id] = true
	}
	return nil
}
