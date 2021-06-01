package types

import sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

var (
	ErrEscrowNotOpen  = sdkerrors.Register(ModuleName, 1, "The escrow should be in open state")
	ErrIdNotUnique    = sdkerrors.Register(ModuleName, 2, "The generated escrow ID is not unique")
	ErrInvalidID      = sdkerrors.Register(ModuleName, 3, "The escrow ID is not valid")
	ErrEscrowNotFound = sdkerrors.Register(ModuleName, 3, "The escrow does not exists")
)
