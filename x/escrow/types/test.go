package types

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
)

// Assert that TestObject is a TransferableObject
var _ TransferableObject = &TestObject{}

// Assert that TestTimeConstrainedObject is a TransferableObject and an ObjectWithTimeConstraint
var _ TransferableObject = &TestTimeConstrainedObject{}
var _ ObjectWithTimeConstraint = &TestTimeConstrainedObject{}

const TypeIDTestObject = 1
const TypeIDTestTimeConstrainedObject = 2

func (m *TestObject) PrimaryKey() []byte {
	return sdk.Uint64ToBigEndian(m.Id)
}

func (m *TestObject) SecondaryKeys() []crud.SecondaryKey {
	return make([]crud.SecondaryKey, 0)
}

func (m *TestObject) GetObjectTypeID() TypeID {
	return TypeIDTestObject
}

func (m *TestObject) GetUniqueKey() []byte {
	return m.PrimaryKey()
}

func (m *TestObject) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	// If Owner is nil, then this tests object belongs to anyone
	if m.Owner == nil {
		return true, nil
	}
	return m.Owner.Equals(account), nil
}

func (m *TestObject) Transfer(_ sdk.Context, from sdk.AccAddress, to sdk.AccAddress, data CustomData) error {
	store := data.(crud.Store)
	if err := store.Read(m.PrimaryKey(), m); err != nil {
		return sdkerrors.Wrap(err, "The object is not synchronized with the store")
	}

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
	if err := store.Update(m); err != nil {
		return sdkerrors.Wrap(err, "Cannot update the object after transfer")
	}
	return nil
}

type expectedData interface {
	GetDeadlineOrDefault(sdk.Context, TransferableObject, uint64) uint64
	GetCrudStore() crud.Store
}

func (m *TestTimeConstrainedObject) ValidateDeadlineBasic(time uint64) error {
	if time >= m.Expiration {
		return fmt.Errorf("this object is expired : %v >= %v", time, m.Expiration)
	}
	return nil
}

func (m *TestTimeConstrainedObject) ValidateDeadline(ctx sdk.Context, time uint64, data CustomData) error {
	if data == nil {
		return fmt.Errorf("no custom data found for TestTimeConstrainedObjects")
	}
	extractor, hasExtractor := data.(expectedData)
	if !hasExtractor {
		return fmt.Errorf("custom data not set properly for TestTimeConstrainedObject: invalid type %T", data)
	}
	oldExpiration := m.Expiration
	m.Expiration = extractor.GetDeadlineOrDefault(ctx, m, m.Expiration)
	err := m.ValidateDeadlineBasic(time)
	m.Expiration = oldExpiration
	return err
}

func (m *TestTimeConstrainedObject) PrimaryKey() []byte {
	return sdk.Uint64ToBigEndian(m.Id)
}

func (m *TestTimeConstrainedObject) SecondaryKeys() []crud.SecondaryKey {
	return make([]crud.SecondaryKey, 0)
}

func (m *TestTimeConstrainedObject) GetObjectTypeID() TypeID {
	return TypeIDTestTimeConstrainedObject
}

func (m *TestTimeConstrainedObject) GetUniqueKey() []byte {
	return m.PrimaryKey()
}

func (m *TestTimeConstrainedObject) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	return m.Owner.Equals(account), nil
}

func (m *TestTimeConstrainedObject) Transfer(_ sdk.Context, from sdk.AccAddress, to sdk.AccAddress, data CustomData) error {
	store := data.(expectedData).GetCrudStore()
	if err := store.Read(m.PrimaryKey(), m); err != nil {
		panic(err)
	}

	if owned, err := m.IsOwnedBy(from); !owned || err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrUnauthorized, "the object %v does not belong to %v", m.Id, from.String())
	}
	m.Owner = to

	if err := store.Update(m); err != nil {
		panic(err)
	}
	return nil
}
