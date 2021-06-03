package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrEscrowNotOpen        = sdkerrors.Register(ModuleName, 1, "The escrow should be in open state")
	ErrIdNotUnique          = sdkerrors.Register(ModuleName, 2, "The generated escrow ID is not unique")
	ErrInvalidID            = sdkerrors.Register(ModuleName, 3, "The escrow ID is not valid")
	ErrEscrowNotFound       = sdkerrors.Register(ModuleName, 4, "This escrow does not exists")
	ErrInvalidAmount        = sdkerrors.Register(ModuleName, 5, "The provided amount is invalid")
	ErrTransferAmountTooLow = sdkerrors.Register(ModuleName, 6, "The transfer amount is lower than the defined price")
	ErrUnknownTypeID        = sdkerrors.Register(ModuleName, 7, "The object type ID is unknown")
	ErrUnknownObject        = sdkerrors.Register(ModuleName, 8, "The object does not exist")
	ErrObjectNotAvailable   = sdkerrors.Register(ModuleName, 9, "This object is not available for this account")
)
