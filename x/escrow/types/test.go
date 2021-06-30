package types

import (
	"fmt"

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

func (m *TestObject) GetType() TypeID {
	return TypeIDTestObject
}

func (m *TestObject) GetObject() crud.Object {
	return m
}

func (m *TestObject) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	return m.Owner.Equals(account), nil
}

func (m *TestObject) Transfer(from sdk.AccAddress, to sdk.AccAddress) error {
	if !m.Owner.Equals(from) {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "the object %v does not belong to %v", m.Id, from.String())
	}
	m.Owner = to
	return nil
}

type NotPossessedTestObject struct {
	TestObject
}

func (m *NotPossessedTestObject) IsOwnedBy(sdk.AccAddress) (bool, error) {
	return true, nil
}

func (m *NotPossessedTestObject) Transfer(sdk.AccAddress, sdk.AccAddress) error {
	return nil
}

type ErroredTestObject struct {
	NotPossessedTestObject
	NbTransferAllowed int64
}

func (m *ErroredTestObject) Transfer(sdk.AccAddress, sdk.AccAddress) error {
	if m.NbTransferAllowed > 0 {
		m.NbTransferAllowed--
		return nil
	}
	return fmt.Errorf("this test object cannot be transfered")
}
