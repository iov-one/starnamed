package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewEscrow constructs a new escrow instance
func NewEscrow(
	id string,
	seller sdk.AccAddress,
	price sdk.Coins,
	object TransferableObject,
	deadline uint64,
) Escrow {
	objectAny, err := codectypes.NewAnyWithValue(object)
	if err != nil {
		panic(err)
	}
	return Escrow{
		Id:       id,
		Seller:   seller.String(),
		Object:   objectAny,
		Price:    price,
		State:    EscrowState_Open,
		Deadline: deadline,
	}
}

// UnpackInterfaces make sure the Anys included in Escrow are unpacked (e.g the object field)
func (msg *Escrow) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.Object != nil {
		var obj TransferableObject
		return unpacker.UnpackAny(msg.Object, &obj)
	}

	return nil
}

// ValidateWithoutDeadlineAndObject validates the escrow without validating the deadline and the object
func (e Escrow) ValidateWithoutDeadlineAndObject() error {
	if err := ValidateID(e.Id); err != nil {
		return err
	}
	// Validate seller address
	_, err := sdk.AccAddressFromBech32(e.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}

	// Validate price
	if err := ValidatePrice(e.Price); err != nil {
		return err
	}
	// Validate state
	return ValidateState(e.State)
}

// Validate validates the escrow
func (e Escrow) Validate(lastBlockTime uint64) error {
	// Validate all fields expect deadline and object
	if err := e.ValidateWithoutDeadlineAndObject(); err != nil {
		return err
	}

	// Validate object valid and possessed by the seller
	// The seller address is already validated
	seller, _ := sdk.AccAddressFromBech32(e.Seller)
	if err := ValidateObject(e.GetObject(), seller); err != nil {
		return err
	}

	// Validate deadline
	return ValidateDeadline(e.Deadline, lastBlockTime)
}

func (e *Escrow) GetObject() TransferableObject {
	return e.Object.GetCachedValue().(TransferableObject)
}

func (e *Escrow) SyncObject() {
	any, err := codectypes.NewAnyWithValue(e.GetObject())
	if err != nil {
		panic(sdkerrors.Wrap(err, "Cannot sync escrow object with its cached value"))
	}
	e.Object = any
}
