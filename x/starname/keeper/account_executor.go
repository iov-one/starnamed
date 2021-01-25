package executor

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

// NewAccount is the constuctor for an account executor
func NewAccount(ctx sdk.Context, account types.Account) *Account {
	return &Account{
		account: &account,
		ctx:     ctx,
	}
}

// Account defines an account executor
type Account struct {
	store   *crud.Store
	account *types.Account
	ctx     sdk.Context
	conf    *configuration.Config
}

// WithAccounts allows to specify a cached accounts store
func (a *Account) WithAccounts(store *crud.Store) *Account {
	a.store = store
	return a
}

// WithConfiguration allows to specify a cached config
func (a *Account) WithConfiguration(cfg configuration.Config) *Account {
	a.conf = &cfg
	return a
}

// Transfer transfers the account to the provided owner with information reset if reset is true
func (a *Account) Transfer(newOwner sdk.AccAddress, reset bool) {
	if a.account == nil {
		panic("cannot transfer non specified account")
	}
	// apply account changes
	// update owner
	a.account.Owner = newOwner
	// if reset is required then clear the account
	if reset {
		a.account.Certificates = nil
		a.account.Resources = nil
		a.account.MetadataURI = ""
	}
	// apply changes
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Update(a.account)
}

// UpdateMetadata updates account's metadata
func (a *Account) UpdateMetadata(newMetadata string) {
	if a.account == nil {
		panic("cannot update metadata on non specified account")
	}
	a.account.MetadataURI = newMetadata
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Update(a.account)
}

// ReplaceResources replaces account's resources
func (a *Account) ReplaceResources(newTargets []*types.Resource) {
	if a.account == nil {
		panic("cannot replace targets on non specified account")
	}
	a.account.Resources = newTargets
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Update(a.account)
}

// Renew renews an account
func (a *Account) Renew() {
	if a.account == nil {
		panic("cannot renew a non specified account")
	}
	renew := a.conf.AccountRenewalPeriod
	a.account.ValidUntil = utils.TimeToSeconds(
		utils.SecondsToTime(a.account.ValidUntil).Add(renew),
	)
	// update account in kv store
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Update(a.account)
}

// Create creates an account
func (a *Account) Create() {
	if a.account == nil {
		panic("cannot create a non specified account")
	}
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Create(a.account)
}

// Delete deletes the account
func (a *Account) Delete() {
	if a.account == nil {
		panic("cannot delete a non specified account")
	}
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Delete(a.account.PrimaryKey())
}

// DeleteCertificate deletes the certificate of the account at the provided index
func (a *Account) DeleteCertificate(index int) {
	if a.account == nil {
		panic("cannot delete certificate on a non specified account")
	}
	a.account.Certificates = append(a.account.Certificates[:index], a.account.Certificates[index+1:]...)
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Update(a.account)
}

// AddCertificate adds a certificate to the account
func (a *Account) AddCertificate(cert []byte) {
	if a.account == nil {
		panic("cannot add certificate on a non specified account")
	}
	a.account.Certificates = append(a.account.Certificates, cert)
	if a.store == nil {
		panic("store is missing")
	}
	(*a.store).Update(a.account)
}

// State returns the current state of the account
func (a *Account) State() types.Account {
	if a.account == nil {
		panic("cannot get state of a non specified account")
	}
	return *a.account
}
