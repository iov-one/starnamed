package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
)

// Assert that TestObject is a TransferableObject
var _ TransferableObject = &TestObject{}

const TypeIDTestObject = 1

func (m *TestObject) PrimaryKey() []byte {
	return sdk.Uint64ToBigEndian(m.Id)
}

func (m *TestObject) SecondaryKeys() []crud.SecondaryKey {
	return make([]crud.SecondaryKey, 0)
}

func (m *TestObject) GetObjectTypeID() TypeID {
	return TypeIDTestObject
}

func (m *TestObject) GetCRUDObject() crud.Object {
	return m
}

func (m *TestObject) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	// If Owner is nil, then this tests object belongs to anyone
	if m.Owner == nil {
		return true, nil
	}
	return m.Owner.Equals(account), nil
}

func (m *TestObject) Transfer(from sdk.AccAddress, to sdk.AccAddress) error {
	if m.NumAllowedTransfers == 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "this test object cannot be transferred")
	} else if m.NumAllowedTransfers > 0 {
		m.NumAllowedTransfers--
	}

	if owned, err := m.IsOwnedBy(from); !owned || err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "the object %v does not belong to %v", m.Id, from.String())
	}
	if m.Owner != nil {
		m.Owner = to
	}
	return nil
}
