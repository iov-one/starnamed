package types

import (
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"

	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"

	"github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
)

const DomainAdminIndex crud.IndexID = 0x1
const DomainBrokerIndex crud.IndexID = 0x2
const AccountAdminIndex crud.IndexID = 0x1
const AccountDomainIndex crud.IndexID = 0x2
const AccountResourcesIndex crud.IndexID = 0x3
const AccountBrokerIndex crud.IndexID = 0x4

// Type IDs used by the escrow module
const (
	DomainTypeID  escrowtypes.TypeID = 0x1
	AccountTypeID escrowtypes.TypeID = 0x2
)

// Ensure that Account and Domain implement crud.Object, escrowtypes.TransferableObject and escrowtypes.ObjectWithTimeConstraint

// TODO: when creator domain type is created, let accounts be TransferableObjects
//var _ escrowtypes.TransferableObject = &Account{}
var _ escrowtypes.TransferableObject = &Domain{}

var _ crud.Object = &Account{}
var _ crud.Object = &Domain{}

var _ escrowtypes.ObjectWithTimeConstraint = &Account{}
var _ escrowtypes.ObjectWithTimeConstraint = &Domain{}

// Delimit the uri and resource in GetResourceKey() with an ineligible
// character since, technically, it'd be possible to have uri "d" and
// resource "ave" collide with uri "da" and resource "ve" without a
// delimiter.
const resourceDelimiter = "\t"

// GetResourceKey computes the index key for a given uri and resource
func GetResourceKey(uri, resource string) []byte {
	return []byte(strings.Join([]string{uri, resource}, resourceDelimiter))
}

// StarnameSeparator defines the starname separator identifier
const StarnameSeparator = "*"

// expectedTransferKeeper is the expected interface for accounts and domains transfer.
// It is used in the Transfer method of those objects, to cast the provided custom data
type expectedTransferKeeper interface {
	DoAccountTransfer(ctx sdk.Context, name string, domain string, currentOwner sdk.AccAddress, newOwner sdk.AccAddress, toReset bool) (*Account, *Domain, error)
	DoDomainTransfer(ctx sdk.Context, domain string, currentOwner sdk.AccAddress, newOwner sdk.AccAddress, transferFlag TransferFlag) error

	AccountStore(ctx sdk.Context) crud.Store
	DomainStore(ctx sdk.Context) crud.Store
}

// PrimaryKey fulfills part of the crud.Object interface
func (m *Domain) PrimaryKey() []byte {
	if m.Name == "" {
		return nil
	}
	return []byte(m.Name)
}

func (m *Domain) SecondaryKeys() []crud.SecondaryKey {
	var sks []crud.SecondaryKey
	// index by owner
	if !m.Admin.Empty() {
		idx := crud.SecondaryKey{DomainAdminIndex, m.Admin}
		sks = append(sks, idx)
	}
	// index by broker
	if !m.Broker.Empty() {
		idx := crud.SecondaryKey{DomainBrokerIndex, m.Broker}
		sks = append(sks, idx)
	}
	return sks
}

// Make Domain implement escrowtypes.TransferableObject

// GetObjectTypeID implements escrowtypes.TransferableObject
func (m *Domain) GetObjectTypeID() escrowtypes.TypeID {
	return DomainTypeID
}

// GetCRUDObject implements escrowtypes.TransferableObject
func (m *Domain) GetCRUDObject() crud.Object {
	return m
}

// IsOwnedBy implements escrowtypes.TransferableObject
func (m *Domain) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	return m.Admin.Equals(account), nil
}

// Transfer implements escrowtypes.TransferableObject
func (m *Domain) Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, data escrowtypes.CustomData) error {
	// Extract the custom data (which is the starname keeper)
	k, correct := data.(expectedTransferKeeper)
	if !correct {
		panic("Corrupted custom data for Domain : the data should be a starname keeper")
	}

	// Should the transfer remove the accounts associated with this domain
	// TODO: choose this wisely
	const transferType = TransferFlush

	return k.DoDomainTransfer(ctx, m.Name, from, to, transferType)
}

// Make Domain implement escrowtypes.ObjectWithTimeConstraint

// ValidateDeadline implements escrowtypes.TransferableObject
func (m *Domain) ValidateDeadline(time uint64) error {
	if uint64(m.ValidUntil) <= time {
		return ErrDomainExpired
	}
	return nil
}

// DomainType defines the type of the domain
type DomainType string

const (
	// OpenDomain is the domain type in which an account owner is the only entity that can perform actions on the account
	OpenDomain DomainType = "open"
	// ClosedDomain is the domain type in which the domain owner has control over accounts too
	ClosedDomain = "closed"
)

func ValidateDomainType(typ DomainType) error {
	switch typ {
	case OpenDomain, ClosedDomain:
		return nil
	default:
		return errors.Wrapf(ErrInvalidDomainType, "invalid domain type: %s", typ)
	}
}

func (m *Account) GetStarname() string {
	if len(m.Domain) == 0 || m.Name == nil {
		return "invalid Domain or Name"
	}
	return strings.Join([]string{*m.Name, m.Domain}, StarnameSeparator)
}

func (m *Account) PrimaryKey() []byte {
	if len(m.Domain) == 0 || m.Name == nil {
		return nil
	}
	j := strings.Join([]string{m.Domain, *m.Name}, StarnameSeparator)
	return []byte(j)
}

func (m *Account) SecondaryKeys() []crud.SecondaryKey {
	var sk []crud.SecondaryKey
	// index by owner
	if !m.Owner.Empty() {
		ownerIndex := crud.SecondaryKey{AccountAdminIndex, m.Owner}
		sk = append(sk, ownerIndex)
	}
	// index by domain
	if len(m.Domain) != 0 {
		domainIndex := crud.SecondaryKey{AccountDomainIndex, []byte(m.Domain)}
		sk = append(sk, domainIndex)
	}
	// index by broker
	if !m.Broker.Empty() {
		brokerIndex := crud.SecondaryKey{AccountBrokerIndex, m.Broker}
		sk = append(sk, brokerIndex)
	}
	// index by resources
	for _, res := range m.Resources {
		// exclude empty resources
		if res.Resource == "" || res.URI == "" {
			continue
		}
		key := GetResourceKey(res.URI, res.Resource)
		// append resource
		sk = append(sk, crud.SecondaryKey{AccountResourcesIndex, key})
	}
	// return keys
	return sk
}

// Make Account implement escrowtypes.TransferableObject
// TODO: when creator domain type is created, let accounts be TransferableObjects
/*
// GetType implements escrowtypes.TransferableObject
func (m *Account) GetObjectTypeID() escrowtypes.TypeID {
	return AccountTypeID
}

// GetCRUDObject implements escrowtypes.TransferableObject
func (m *Account) GetCRUDObject() crud.Object {
	return m
}

// IsOwnedBy implements escrowtypes.TransferableObject
func (m *Account) IsOwnedBy(account sdk.AccAddress) (bool, error) {
	return m.Owner.Equals(account), nil
}

// Transfer implements escrowtypes.TransferableObject
func (m *Account) Transfer(ctx sdk.Context, from sdk.AccAddress, to sdk.AccAddress, data escrowtypes.CustomData) error {
	// Extract the custom data (which is the starname keeper)
	k, correct := data.(expectedTransferKeeper)
	if !correct {
		panic("Corrupted custom data for Account : the data should be a starname keeper")
	}

	// Should the transfer reset the account metadata, resources and certificates
	const shouldReset = false

	_, _, err := k.DoAccountTransfer(ctx, *m.Name, m.Domain, from, to, shouldReset)

	return err
}
*/
// Make Account implement escrowtypes.ObjectWithTimeConstraint

// ValidateDeadline implements escrowtypes.TransferableObject
func (m *Account) ValidateDeadline(time uint64) error {
	if uint64(m.ValidUntil) <= time {
		return ErrAccountExpired
	}
	return nil
}

// Certificate defines a certificate
type Certificate []byte
