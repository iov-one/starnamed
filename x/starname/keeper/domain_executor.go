package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

// DomainExecutor defines the domain keeper executor
type DomainExecutor struct {
	domain   *types.Domain
	ctx      sdk.Context
	domains  *crud.Store
	accounts *crud.Store
	conf     *configuration.Config
}

// NewDomainExecutor returns is domain's constructor
func NewDomainExecutor(ctx sdk.Context, dom types.Domain) *DomainExecutor {
	return &DomainExecutor{
		ctx:    ctx,
		domain: &dom,
	}
}

// WithAccounts allows to specify a cached config
func (d *DomainExecutor) WithAccounts(store *crud.Store) *DomainExecutor {
	d.accounts = store
	return d
}

// WithDomains allows to specify a cached config
func (d *DomainExecutor) WithDomains(store *crud.Store) *DomainExecutor {
	d.domains = store
	return d
}

// WithConfiguration allows to specify a cached config
func (d *DomainExecutor) WithConfiguration(cfg configuration.Config) *DomainExecutor {
	d.conf = &cfg
	return d
}

// Renew renews a domain based on the configuration or accValidUntil
func (d *DomainExecutor) Renew(accValidUntil ...int64) {
	if d.domain == nil {
		panic("cannot execute renew state change on non present domain")
	}
	if d.domains == nil {
		panic("domains is missing")
	}
	// if account valid until is specified then the renew is coming from accounts
	if len(accValidUntil) != 0 {
		d.domain.ValidUntil = accValidUntil[0]
		(*d.domains).Update(d.domain)
		return
	}
	// get configuration
	if d.conf == nil {
		panic("conf is missing")
	}
	renewDuration := d.conf.DomainRenewalPeriod
	// update domain valid until
	d.domain.ValidUntil = utils.TimeToSeconds(
		utils.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// set domain
	(*d.domains).Update(d.domain)
	// update empty account
	account, cursor := d.getEmptyNameAccount()
	account.ValidUntil = d.domain.ValidUntil
	(*cursor).Update(account)
}

// Delete deletes a domain from the kvstore
func (d *DomainExecutor) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present domain")
	}
	if d.accounts == nil {
		panic("accounts is missing")
	}
	cursor, err := (*d.accounts).Query().Where().Index(types.AccountDomainIndex).Equals([]byte(d.domain.Name)).Do()
	if err != nil {
		panic(err)
	}
	for ; cursor.Valid(); cursor.Next() {
		if err = cursor.Delete(); err != nil {
			panic(err)
		}
	}
	if d.domains == nil {
		panic("domains is missing")
	}
	(*d.domains).Delete(d.domain.PrimaryKey())
}

// Transfer transfers a domain given a flag and an owner
func (d *DomainExecutor) Transfer(flag types.TransferFlag, newOwner sdk.AccAddress) {
	if d.domain == nil {
		panic("cannot execute transfer state on non defined domain")
	}
	if d.accounts == nil {
		panic("accounts is missing")
	}
	if d.domains == nil {
		panic("domains is missing")
	}
	// transfer domain
	var oldOwner = d.domain.Admin // cache it for future uses
	d.domain.Admin = newOwner
	(*d.domains).Update(d.domain)
	// transfer empty account
	account, _ := d.getEmptyNameAccount()
	executor := NewAccountExecutor(d.ctx, *account).WithAccounts(d.accounts)
	executor.Transfer(newOwner, false)
	// transfer accounts of the domain based on the transfer flag
	switch flag {
	// reset none is simply skipped as empty account is already transferred during domain transfer
	case types.TransferResetNone:
		return
	// transfer flush, deletes all domain's accounts except the empty one since it was transferred in the first step
	case types.TransferFlush:
		cursor, err := (*d.accounts).Query().Where().Index(types.AccountDomainIndex).Equals(d.domain.PrimaryKey()).Do()
		if err != nil {
			panic(err)
		}
		for ; cursor.Valid(); cursor.Next() {
			account.Reset()
			if err = cursor.Read(account); err != nil {
				panic(err)
			}
			ex := NewAccountExecutor(d.ctx, *account).WithAccounts(d.accounts)
			// reset the empty account...
			if *account.Name == types.EmptyAccountName {
				ex.Transfer(newOwner, true)
				continue
			}
			// ...delete all others
			ex.Delete()
		}
	// transfer owned transfers only accounts owned by the old owner
	case types.TransferOwned:
		cursor, err := (*d.accounts).Query().Where().
			Index(types.AccountDomainIndex).Equals(d.domain.PrimaryKey()).And().
			Index(types.AccountAdminIndex).Equals(oldOwner.Bytes()).Do()
		if err != nil {
			panic(err)
		}
		for ; cursor.Valid(); cursor.Next() {
			account.Reset()
			if err = cursor.Read(account); err != nil {
				panic(err)
			}
			ex := NewAccountExecutor(d.ctx, *account).WithAccounts(d.accounts)
			// transfer accounts without reset
			ex.Transfer(newOwner, false)
		}
	}
}

// Create creates a new domain
func (d *DomainExecutor) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	if d.accounts == nil {
		panic("accounts is missing")
	}
	if d.domains == nil {
		panic("domains is missing")
	}
	(*d.domains).Create(d.domain)
	emptyAccount := &types.Account{
		Domain:       d.domain.Name,
		Name:         utils.StrPtr(types.EmptyAccountName),
		Owner:        d.domain.Admin,
		ValidUntil:   d.domain.ValidUntil,
		Resources:    nil,
		Certificates: nil,
		Broker:       nil,
		MetadataURI:  "",
	}
	(*d.accounts).Create(emptyAccount)
}

// Gets the empty name account and cursor
func (d *DomainExecutor) getEmptyNameAccount() (*types.Account, *crud.Cursor) {
	if d.accounts == nil {
		panic("domains is missing")
	}
	cursor, err := (*d.accounts).Query().Where().Index(types.AccountDomainIndex).Equals(d.domain.PrimaryKey()).Do()
	if err != nil {
		panic(err)
	}
	account := new(types.Account)
	for ; cursor.Valid(); cursor.Next() {
		account.Reset()
		if err = cursor.Read(account); err != nil {
			panic(err)
		}
		if account.Name != nil && *account.Name == types.EmptyAccountName {
			return account, &cursor
		}
	}
	panic(fmt.Sprintf("failed to get empty account in domain %s", d.domain.Name))
}
