package types

import (
	"strings"

	"github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
)

const DomainAdminIndex byte = 0x1
const AccountAdminIndex byte = 0x1
const AccountDomainIndex byte = 0x2
const AccountResourcesIndex byte = 0x3

// StarnameSeparator defines the starname separator identifier
const StarnameSeparator = "*"

// PrimaryKey fulfills part of the crud.Object interface
func (m *Domain) PrimaryKey() []byte {
	if m.Name == "" {
		return nil
	}
	return []byte(m.Name)
}

func (m *Domain) SecondaryKeys() []crud.SecondaryKey {
	if m.Admin.Empty() {
		return nil
	}
	sk := crud.SecondaryKey{DomainAdminIndex, m.Admin}
	return []crud.SecondaryKey{sk}
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

/* TODO: FIXME
type accountCodecAlias struct {
	Underlying *Account
	NameNil    bool
}

func (m *Account) MarshalCRUD() interface{} {
	return accountCodecAlias{
		Underlying: m,
		NameNil:    m.Name == nil,
	}
}

func (m *Account) UnmarshalCRUD(cdc codec.Marshaler, b []byte) {
	trg := new(accountCodecAlias)
	cdc.MustUnmarshalBinaryBare(b, trg)
	*m = *trg.Underlying
	if m.Name == nil && !trg.NameNil {
		m.Name = utils.StrPtr("")
	}
}
*/

func (m *Account) PrimaryKey() []byte {
	if len(m.Domain) == 0 || m.Name == nil {
		return nil
	}
	j := strings.Join([]string{m.Domain, *m.Name}, "*")
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
	// index by resources
	for _, res := range m.Resources {
		// exclude empty resources
		if res.Resource == "" || res.URI == "" {
			continue
		}
		resKey := strings.Join([]string{res.URI, res.Resource}, "")
		// append resource
		sk = append(sk, crud.SecondaryKey{AccountResourcesIndex, []byte(resKey)})
	}
	// return keys
	return sk
}

// Certificate defines a certificate
type Certificate []byte
