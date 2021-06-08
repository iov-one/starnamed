package types

import (
	"fmt"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewGenesisState constructs a new GenesisState instance
func NewGenesisState(escrows []Escrow, lastBlockTime uint64) *GenesisState {
	return &GenesisState{
		Escrows:       escrows,
		LastBlockTime: lastBlockTime,
	}
}

// DefaultGenesisState gets the raw genesis message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Escrows:       []Escrow{},
		LastBlockTime: 0,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	ids := map[string]bool{}
	for _, escrow := range data.Escrows {
		if ids[escrow.Id] {
			return fmt.Errorf("found duplicate escrow ID %s", escrow.Id)
		}

		if escrow.State != EscrowState_Open {
			return sdkerrors.Wrap(ErrEscrowNotOpen, escrow.Id)
		}

		//TODO: duplicate check (but with more explicit error message) with escrow.Validate() => validateDeadline()
		if escrow.Deadline <= data.LastBlockTime {
			return sdkerrors.Wrap(ErrEscrowExpired, escrow.Id)
		}

		if err := escrow.Validate(data.LastBlockTime); err != nil {
			return err
		}

		ids[escrow.Id] = true
	}
	return nil
}
