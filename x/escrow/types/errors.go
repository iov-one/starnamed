package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrEscrowNotOpen         = sdkerrors.Register(ModuleName, 1, "The escrow should be in open state")
	ErrIdNotUnique           = sdkerrors.Register(ModuleName, 2, "The generated escrow ID is not unique")
	ErrInvalidID             = sdkerrors.Register(ModuleName, 3, "The escrow ID is not valid")
	ErrEscrowNotFound        = sdkerrors.Register(ModuleName, 4, "This escrow does not exists")
	ErrEscrowExpired         = sdkerrors.Register(ModuleName, 5, "This escrow is expired")
	ErrInvalidAmount         = sdkerrors.Register(ModuleName, 6, "The provided amount is invalid")
	ErrTransferAmountTooLow  = sdkerrors.Register(ModuleName, 7, "The transfer amount is lower than the defined price")
	ErrUnknownTypeID         = sdkerrors.Register(ModuleName, 8, "The object type ID is unknown")
	ErrUnknownObject         = sdkerrors.Register(ModuleName, 9, "The object does not exist")
	ErrInvalidAccount        = sdkerrors.Register(ModuleName, 10, "This account cannot be used in an escrow")
	ErrPastDeadline          = sdkerrors.Register(ModuleName, 11, "The deadline has passed")
	ErrEmptyUpdate           = sdkerrors.Register(ModuleName, 12, "No fields have been filled : nothing to update")
	ErrInvalidCommissionRate = sdkerrors.Register(ModuleName, 13, "The broker commission must be a number between 0 and 1")
	ErrInvalidPrice          = sdkerrors.Register(ModuleName, 14, "The price is invalid")
	ErrInvalidDeadline       = sdkerrors.Register(ModuleName, 15, "The deadline is invalid")
)
