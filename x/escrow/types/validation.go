package types

import (
	"encoding/hex"
	"strconv"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// EscrowIDLength is the length for the hex encoded escrow id (uint64 => 8 bytes => 16 chars)
	EscrowIDLength = 16
)

// ValidatePrice verifies whether the given amount is valid and has the correct denomination
// If denom is empty, does not validate the denomination
func ValidatePrice(price sdk.Coins, denom string) error {
	if !(price.IsValid() && price.IsAllPositive()) {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidCoins, "the price must be valid and positive")
	}
	if len(denom) != 0 && (price.Len() != 1 || price[0].Denom != denom) {
		return sdkerrors.Wrap(ErrInvalidPrice, "the price must be in "+denom+", price: "+price.String())
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
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "The object does not belong to %s", seller)
	}
	return nil
}

func ValidateObjectDeadline(transferableObj TransferableObject, deadline uint64) error {
	if obj, hasCustomCheck := transferableObj.(ObjectWithTimeConstraint); hasCustomCheck {
		return sdkerrors.Wrap(obj.ValidateDeadline(deadline), "the deadline has not been validated by the object")
	}
	return nil
}

func ValidateState(state EscrowState) error {
	if state != EscrowState_Open && state != EscrowState_Expired {
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

func ValidateAddress(addr string) error {
	_, err := sdk.AccAddressFromBech32(addr)
	return err
}

func ValidateCommission(commission sdk.Dec) error {
	if commission.LT(sdk.ZeroDec()) || commission.GT(sdk.OneDec()) {
		return ErrInvalidCommissionRate
	}
	return nil
}
