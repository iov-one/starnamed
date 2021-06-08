package types

import (
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// EscrowIDLength is the length for the hash lock in hex string
	EscrowIDLength = 64
)

// ValidatePrice verifies whether the given amount is legal
func ValidatePrice(price sdk.Coins) error {
	if !(price.IsValid() && price.IsAllPositive()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "the price must be valid and positive")
	}
	return nil
}

// ValidateID verifies whether the given ID lock is legal
func ValidateID(id string) error {
	if len(id) != EscrowIDLength {
		return sdkerrors.Wrapf(ErrInvalidID, "length of the escrow id must be %d", EscrowIDLength)
	}
	if _, err := hex.DecodeString(id); err != nil {
		return sdkerrors.Wrapf(ErrInvalidID, "id must be a hex encoded string")
	}
	return nil
}

func ValidateObject(object TransferableObject, seller sdk.AccAddress) error {
	ownedBySeller, err := object.IsOwnedBy(seller)
	if err != nil {
		return err
	}
	if !ownedBySeller {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "The object does not belong to %v", seller)
	}
	return nil
}

func ValidateState(state EscrowState) error {
	if state != EscrowState_Open {
		return sdkerrors.Wrap(ErrEscrowNotOpen, strconv.FormatUint(uint64(state), 10))
	}
	return nil
}

func ValidateDeadline(deadline uint64, lastBlockTime uint64) error {
	if deadline < lastBlockTime {
		return ErrPastDeadline
	}
	return nil
}
