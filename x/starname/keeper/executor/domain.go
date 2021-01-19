package executor

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
)

// Domain defines the domain keeper executor
type Domain struct {
	domain   *types.Domain
	ctx      sdk.Context
	domains  crud.Store
	accounts crud.Store
	k        keeper.Keeper
}

// NewDomain returns is domain's constructor
func NewDomain(ctx sdk.Context, k keeper.Keeper, dom types.Domain) *Domain {
	return &Domain{
		k:        k,
		ctx:      ctx,
		domains:  k.DomainStore(ctx),
		accounts: k.AccountStore(ctx),
		domain:   &dom,
	}
}

// Renew renews a domain based on the configuration
func (d *Domain) Renew(accValidUntil ...int64) {
	if d.domain == nil {
		panic("cannot execute renew state change on non present domain")
	}
	// if account valid until is specified then the renew is coming from accounts
	if len(accValidUntil) != 0 {
		d.domain.ValidUntil = accValidUntil[0]
		d.domains.Update(d.domain)
		return
	}
	// get configuration
	renewDuration := d.k.ConfigurationKeeper.GetDomainRenewDuration(d.ctx)
	// update domain valid until
	d.domain.ValidUntil = utils.TimeToSeconds(
		utils.SecondsToTime(d.domain.ValidUntil).Add(renewDuration), // time(domain.ValidUntil) + renew duration
	)
	// set domain
	d.domains.Update(d.domain)
	// update empty account
	account, cursor := d.getEmptyNameAccount()
	account.ValidUntil = d.domain.ValidUntil
	(*cursor).Update(account)
}

// Delete deletes a domain from the kvstore
func (d *Domain) Delete() {
	if d.domain == nil {
		panic("cannot execute delete state change on non present domain")
	}
	cursor, err := d.accounts.Query().Where().Index(types.AccountDomainIndex).Equals([]byte(d.domain.Name)).Do()
	if err != nil {
		panic(err)
	}
	for ; cursor.Valid(); cursor.Next() {
		if err = cursor.Delete(); err != nil {
			panic(err)
		}
	}
	d.domains.Delete(d.domain.PrimaryKey())
}

// Transfer transfers a domain given a flag and an owner
func (d *Domain) Transfer(flag types.TransferFlag, newOwner sdk.AccAddress) {
	if d.domain == nil {
		panic("cannot execute transfer state on non defined domain")
	}

	// transfer domain
	var oldOwner = d.domain.Admin // cache it for future uses
	d.domain.Admin = newOwner
	d.domains.Update(d.domain)
	// transfer empty account
	account, _ := d.getEmptyNameAccount()
	executor := NewAccount(d.ctx, d.k, *account)
	executor.Transfer(newOwner, false)
	// transfer accounts of the domain based on the transfer flag
	switch flag {
	// reset none is simply skipped as empty account is already transferred during domain transfer
	case types.TransferResetNone:
		return
	// transfer flush, deletes all domain's accounts except the empty one since it was transferred in the first step
	case types.TransferFlush:
		cursor, err := d.accounts.Query().Where().Index(types.AccountDomainIndex).Equals([]byte(d.domain.Name)).Do()
		if err != nil {
			panic(err)
		}
		for ; cursor.Valid(); cursor.Next() {
			account.Reset()
			if err = cursor.Read(account); err != nil {
				panic(err)
			}
			ex := NewAccount(d.ctx, d.k, *account)
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
		cursor, err := d.accounts.Query().Where().
			Index(types.AccountDomainIndex).Equals([]byte(d.domain.Name)).And().
			Index(types.AccountAdminIndex).Equals(oldOwner.Bytes()).Do()
		if err != nil {
			panic(err)
		}
		for ; cursor.Valid(); cursor.Next() {
			account.Reset()
			if err = cursor.Read(account); err != nil {
				panic(err)
			}
			ex := NewAccount(d.ctx, d.k, *account)
			// transfer accounts without reset
			ex.Transfer(newOwner, false)
		}
	}
}

// Create creates a new domain
func (d *Domain) Create() {
	if d.domain == nil {
		panic("cannot create non specified domain")
	}
	d.domains.Create(d.domain)
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
	d.accounts.Create(emptyAccount)
}

// Gets the empty name account and cursor
func (d *Domain) getEmptyNameAccount() (*types.Account, *crud.Cursor) {
	cursor, err := d.accounts.Query().Where().Index(types.AccountDomainIndex).Equals([]byte(d.domain.Name)).Do()
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
