package types

import (
	"bytes"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewGenesisState constructs a new GenesisState instance
func NewGenesisState(escrows []Escrow, lastBlockTime, nextEscrowID uint64) *GenesisState {
	return &GenesisState{
		Escrows:       escrows,
		LastBlockTime: lastBlockTime,
		NextEscrowId:  nextEscrowID,
	}
}

// DefaultGenesisState gets the raw genesis message for testing
func DefaultGenesisState() *GenesisState {
	return &GenesisState{
		Escrows:       []Escrow{},
		LastBlockTime: 0,
		NextEscrowId:  1,
	}
}

// ValidateGenesis validates the provided genesis state to ensure the expected invariants holds.
func ValidateGenesis(data GenesisState) error {
	ids := map[string]bool{}
	for _, escrow := range data.Escrows {
		// Escrow id must be unique
		if ids[escrow.Id] {
			return fmt.Errorf("found duplicate escrow ID %s", escrow.Id)
		}

		// The escrow id must be issued before data.NextEscrowId
		if bytes.Compare(GetEscrowKey(escrow.Id), sdk.Uint64ToBigEndian(data.NextEscrowId)) >= 0 {
			return fmt.Errorf("found escrow ID greater than next escrow ID : %v", escrow.Id)
		}

		if escrow.Deadline <= data.LastBlockTime && escrow.State != EscrowState_Expired {
			return sdkerrors.Wrap(ErrEscrowExpired, escrow.Id)
		}

		// Validate the escrow fields
		if err := escrow.ValidateWithoutDeadline(); err != nil {
			return err
		}

		// Mark the escrow as seen
		ids[escrow.Id] = true
	}
	return nil
}
