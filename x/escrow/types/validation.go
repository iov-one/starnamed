package types

import (
	"encoding/hex"

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
	if len(id) != HTLCIDLength {
		return sdkerrors.Wrapf(ErrInvalidID, "length of the escrow id must be %d", HTLCIDLength)
	}
	if _, err := hex.DecodeString(id); err != nil {
		return sdkerrors.Wrapf(ErrInvalidID, "id must be a hex encoded string")
	}
	return nil
}

//TODO:
//func ValidateStarname
// func ValidateDeadline
