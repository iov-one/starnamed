package types

import (
	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// MsgWithFeePayer is a sdk.Msg object that has a fee payer information.
type MsgWithFeePayer interface {
	sdk.Msg
	// GetFeePayer returns the address of the account that will pay fees for this message
	GetFeePayer() sdk.AccAddress
}

// TypeID is the type of the object type id
type TypeID uint64

// CustomData represents arbitrary data
type CustomData interface{}

// ObjectWithCustomFees is an object (that should be a TransferableObject in the context of this module) that
// retrieve the fees for several operations.
type ObjectWithCustomFees interface {
	// GetCreationFees returns the fees (as used in the starnamed/x/configuration package) of a creation operation
	// for this object.
	GetCreationFees() sdk.Dec
}

// ObjectWithTimeConstraint is an object that should be a TransferableObject in the context of this module) that can
// validate the deadline of an escrow. This validation is done upon creation and update.
type ObjectWithTimeConstraint interface {
	// ValidateDeadline returns an error if this object rejects the given deadline (a Unix timestamp), e.g. this object will
	// not be valid at this date.
	// ValidateDeadline is responsible for doing all the necessary checks, including those which may be included in
	// ValidateDeadlineBasic
	ValidateDeadline(ctx sdk.Context, deadline uint64, data CustomData) error
	// ValidateDeadlineBasic is like ValidateDeadline but without any state information and extra data. This should implement
	// all the state-less checks for this object, if any.
	ValidateDeadlineBasic(deadline uint64) error
}

// TransferableObject is the object type that is used in escrows.
// It is an object that can be marshalled, transferred and that has a unique type ID.
type TransferableObject interface {
	codec.ProtoMarshaler

	// GetObjectTypeID returns the unique TypeID of this object.
	GetObjectTypeID() TypeID
	// GetUniqueKey returns a byte array that is unique for this type of object
	GetUniqueKey() []byte

	// IsOwnedBy return true if this account is owned by the specified account.
	// If IsOwnedBy returns an error, then the boolean value should be ignored.
	IsOwnedBy(account sdk.AccAddress) (bool, error)
	// Transfer transfers this object from an account to another. It returns an error if the transfer is invalid and the
	// transfer should not happen in that case.
	// It may receive custom data set by the keeper.RegisterCustomData method.
	Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, data CustomData) error
}
