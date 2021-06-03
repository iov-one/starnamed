package types

import (
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// HTLCIDLength is the length for the hash lock in hex string
	HTLCIDLength = 64
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
	//TODO: make that have sense
	if len(id) != HTLCIDLength {
		return sdkerrors.Wrapf(ErrInvalidID, "length of the escrow id must be %d", HTLCIDLength)
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
	if state != Open {
		return sdkerrors.Wrap(ErrEscrowNotOpen, strconv.FormatUint(uint64(state), 10))
	}
	return nil
}

func ValidateDeadline(deadline uint64) error {
	//TODO: check deadline
	return nil
}
