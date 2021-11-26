package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	crud "github.com/iov-one/cosmos-sdk-crud"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// Check that the escrow implements crud.Object
var _ crud.Object = &Escrow{}

const (
	// StateIndex represents the state as a secondary key of an escrow
	// The state is serialized as a byte array from the big endian representation of the enumeration index
	StateIndex = 0x01
	// SellerIndex represents the seller as a secondary key of an escrow
	SellerIndex = 0x02
	// ObjectIndex represents the object as a secondary key of an escrow
	// The object is represented by its primary key
	ObjectIndex = 0x03
)

// NewEscrow constructs a new escrow instance
func NewEscrow(
	id string,
	seller sdk.AccAddress,
	price sdk.Coins,
	object TransferableObject,
	deadline uint64,
	brokerAddress string,
	brokerCommission sdk.Dec,
) Escrow {
	objectAny, err := codectypes.NewAnyWithValue(object)
	if err != nil {
		panic(err)
	}
	return Escrow{
		Id:               id,
		Seller:           seller.String(),
		Object:           objectAny,
		Price:            price,
		State:            EscrowState_Open,
		Deadline:         deadline,
		BrokerAddress:    brokerAddress,
		BrokerCommission: brokerCommission,
	}
}

// PrimaryKey implements crud.Object
func (e Escrow) PrimaryKey() []byte {
	return GetEscrowKey(e.Id)
}

// SecondaryKeys implements crud.Object
func (e Escrow) SecondaryKeys() []crud.SecondaryKey {
	// If this is an empty object, return an empty array
	if len(e.Id) == 0 {
		return make([]crud.SecondaryKey, 0)
	}
	sks := make([]crud.SecondaryKey, 3)
	sks[0] = crud.SecondaryKey{
		ID:    StateIndex,
		Value: sdk.Uint64ToBigEndian(uint64(e.State)),
	}
	seller, err := sdk.AccAddressFromBech32(e.Seller)
	if err == nil {
		sks[1] = crud.SecondaryKey{
			ID:    SellerIndex,
			Value: seller,
		}
	}
	sks[2] = crud.SecondaryKey{
		ID:    ObjectIndex,
		Value: e.GetObject().GetUniqueKey(),
	}
	return sks
}

// UnpackInterfaces make sure the Anys included in Escrow are unpacked (e.g the object field)
func (e *Escrow) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if e.Object != nil {
		var obj TransferableObject
		return unpacker.UnpackAny(e.Object, &obj)
	}

	return nil
}

// ValidateWithoutDeadlineAndObject validates the escrow without validating the deadline and the object
// If priceDenom is empty, does not validate the price denomination
func (e Escrow) ValidateWithoutDeadlineAndObject(priceDenom string) error {
	if err := ValidateID(e.Id); err != nil {
		return err
	}
	// Validate seller address
	err := ValidateAddress(e.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}

	// Validate price
	if err := ValidatePrice(e.Price, priceDenom); err != nil {
		return err
	}

	// Validate broker address
	err = ValidateAddress(e.BrokerAddress)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid broker address (%s)", err)
	}

	// Validate broker commission
	if err := ValidateCommission(e.BrokerCommission); err != nil {
		return err
	}

	// Validate state
	return ValidateState(e.State)
}

// Validate validates the escrow, if priceDenom is empty, does not validate the price denomination
func (e Escrow) Validate(priceDenom string, lastBlockTime uint64) error {
	// Validate all fields expect deadline and object
	if err := e.ValidateWithoutDeadlineAndObject(priceDenom); err != nil {
		return err
	}

	// Validate object valid and possessed by the seller
	// The seller address is already validated
	seller, _ := sdk.AccAddressFromBech32(e.Seller)
	if err := ValidateObject(e.GetObject(), seller); err != nil {
		return err
	}

	// Validate that the object has no conflict with the escrow deadline
	if err := ValidateObjectDeadlineBasic(e.GetObject(), e.Deadline); err != nil {
		return err
	}

	// Validate deadline
	return ValidateDeadline(e.Deadline, lastBlockTime)
}

// ValidateWithContext validates the escrow with a given context and custom data. It performs the same validation as Validate
// plus some context-aware validation
func (e *Escrow) ValidateWithContext(ctx sdk.Context, priceDenom string, lastBlockTime uint64, data CustomData) error {
	if err := e.Validate(priceDenom, lastBlockTime); err != nil {
		return err
	}
	//NOTE: this performs a repeated validation for the basic validation part of the object
	return ValidateObjectDeadline(ctx, e.GetObject(), e.Deadline, data)

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
