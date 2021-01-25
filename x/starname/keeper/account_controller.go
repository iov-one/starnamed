package keeper

import (
	"bytes"
	"regexp"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/controllers/domain"
	"github.com/iov-one/starnamed/x/starname/types"
)

// accountControllerFunc is the function signature used by account controllers
type accountControllerFunc func(ctrl *AccountController) error

// AccountController is an account controller, it caches information
// in order to avoid useless query to state to get the same
// information. Order of execution of controllers matters
// if the correct order is not followed the controller will
// panic because of bad operation flow.
// Errors returned are wrapped sdk.Error types.
type AccountController struct {
	validators []accountControllerFunc

	name, domain string
	account      *types.Account
	conf         *configuration.Config

	ctx        sdk.Context
	store      *crud.Store
	domainCtrl *domain.Domain
}

// Validate verifies the account against the order of provided controllers
func (a *AccountController) Validate() error {
	for _, check := range a.validators {
		if err := check(a); err != nil {
			return err
		}
	}
	return nil
}

// MustExist asserts that the given account exists
func (a *AccountController) MustExist() *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.mustExist()
	})
	return a
}

// MustNotExist asserts that the given account does not exist
func (a *AccountController) MustNotExist() *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.mustNotExist()
	})
	return a
}

// ValidName asserts the account name is valid
func (a *AccountController) ValidName() *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return a.validName()
	})
	return a
}

// NotExpired asserts the account is not expired
func (a *AccountController) NotExpired() *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.notExpired()
	})
	return a
}

// Renewable asserts that the account is renewable
func (a *AccountController) Renewable() *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.renewable()
	})
	return a
}

// OwnedBy asserts that the account is owned by the provided address
func (a *AccountController) OwnedBy(addr sdk.AccAddress) *AccountController {
	f := func(ctrl *AccountController) error {
		return ctrl.ownedBy(addr)
	}
	a.validators = append(a.validators, f)
	return a
}

// CertificateSizeNotExceeded asserts that the size of a cert is not beyond the limits
func (a *AccountController) CertificateSizeNotExceeded(cert []byte) *AccountController {
	f := func(ctrl *AccountController) error {
		return ctrl.certSizeNotExceeded(cert)
	}
	a.validators = append(a.validators, f)
	return a
}

// CertificateLimitNotExceeded asserts that the numbers of certificates in an account was not exceeded
func (a *AccountController) CertificateLimitNotExceeded() *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.certLimitNotExceeded()
	})
	return a
}

// DeletableBy checks if the account can be deleted by the provided address
func (a *AccountController) DeletableBy(addr sdk.AccAddress) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.deletableBy(addr)
	})
	return a
}

// CertificateExists asserts that the provided certificate
// exists and if it does the index is saved in the provided pointer
// if certIndex pointer is nil the certificate index will not be saved
func (a *AccountController) CertificateExists(cert []byte, certIndex *int) *AccountController {
	f := func(ctrl *AccountController) error {
		err := ctrl.certNotExist(cert, certIndex)
		if err == nil {
			return sdkerrors.Wrapf(types.ErrCertificateDoesNotExist, "%x", cert)
		}
		return nil
	}
	a.validators = append(a.validators, f)
	return a
}

// ValidResources verifies that the provided resources are valid for the account
func (a *AccountController) ValidResources(resources []*types.Resource) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.validResources(resources)
	})
	return a
}

// TransferableBy checks if the account can be transferred by the provided address
func (a *AccountController) TransferableBy(addr sdk.AccAddress) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.transferableBy(addr)
	})
	return a
}

// ResettableBy checks if the account attributes resettable by the provided address
func (a *AccountController) ResettableBy(addr sdk.AccAddress, reset bool) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.resettableBy(addr, reset)
	})
	return a
}

// ResourceLimitNotExceeded checks if the number of elements in the provided resource array exceeds the configuration limit
func (a *AccountController) ResourceLimitNotExceeded(resources []*types.Resource) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.resourceLimitNotExceeded(resources)
	})
	return a
}

// MetadataSizeNotExceeded asserts that the metadata size of an account was not exceeded
func (a *AccountController) MetadataSizeNotExceeded(metadata string) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.metadataSizeNotExceeded(metadata)
	})
	return a
}

// RegistrableBy asserts that an account can be registered by the provided address
func (a *AccountController) RegistrableBy(addr sdk.AccAddress) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.registrableBy(addr)
	})
	return a
}

// CertificateNotExist asserts the provided certificate
// does not exist in the account already
func (a *AccountController) CertificateNotExist(cert []byte) *AccountController {
	a.validators = append(a.validators, func(ctrl *AccountController) error {
		return ctrl.certNotExist(cert, nil)
	})
	return a
}

// NewAccountController is Account constructor
func NewAccountController(ctx sdk.Context, domain, name string) *AccountController {
	return &AccountController{
		name:   name,
		domain: domain,
		ctx:    ctx,
	}
}

// WithDomainController allows to specify a cached domain controller
func (a *AccountController) WithDomainController(dom *domain.Domain) *AccountController {
	a.domainCtrl = dom
	return a
}

// WithConfiguration allows to specify a cached config
func (a *AccountController) WithConfiguration(cfg configuration.Config) *AccountController {
	a.conf = &cfg
	return a
}

// WithAccounts allows to specify a cached crud store
func (a *AccountController) WithAccounts(store *crud.Store) *AccountController {
	a.store = store
	return a
}

// WithAccount allows to specify a cached account
func (a *AccountController) WithAccount(acc types.Account) *AccountController {
	a.account = &acc
	a.domain = acc.Domain
	a.name = *acc.Name
	return a
}

// requireDomain builds the domain controller after asserting domain existence
func (a *AccountController) requireDomain() error {
	if a.domainCtrl == nil {
		panic("missing domain controller")
	}
	return a.domainCtrl.MustExist().Validate()
}

// requireAccount finds the accounts and caches it, so future
// queries will always use the same account first found account
func (a *AccountController) requireAccount() error {
	if a.account != nil {
		return nil
	}
	if a.store == nil {
		panic("store is missing")
	}
	account := new(types.Account)
	err := (*a.store).Read((&types.Account{Domain: a.domain, Name: utils.StrPtr(a.name)}).PrimaryKey(), account)
	if err != nil {
		return sdkerrors.Wrapf(types.ErrAccountDoesNotExist, "%s was not found in domain %s", a.name, a.domain)
	}
	a.account = account
	return nil
}

// mustExist makes sure an account exist
func (a *AccountController) mustExist() error {
	return a.requireAccount()
}

// mustNotExist is the unexported function executed by MustNotExist
func (a *AccountController) mustNotExist() error {
	err := a.requireAccount()
	if err != nil {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrAccountExists, "account %s already exists in domain %s", a.name, a.domain)
}

// requireConfiguration updates the configuration
// if it is not already set, and caches it after
func (a *AccountController) requireConfiguration() {
	if a.conf == nil {
		panic("configuration is missing")
	}
}

// validName is the unexported function used by ValidAccountName
func (a *AccountController) validName() error {
	a.requireConfiguration()
	if !regexp.MustCompile(a.conf.ValidAccountName).MatchString(a.name) {
		return sdkerrors.Wrapf(types.ErrInvalidAccountName, "invalid name: %s", a.name)
	}
	return nil
}

// notExpired is the unexported function used by NotExpired
func (a *AccountController) notExpired() error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	if err := a.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	switch a.domainCtrl.Domain().Type {
	// if domain is closed type then skip the expiration validation checks
	case types.ClosedDomain:
		return nil
	}
	// check if account has expired
	expireTime := utils.SecondsToTime(a.account.ValidUntil)
	if !expireTime.Before(a.ctx.BlockTime()) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrAccountExpired, "account %s in domain %s has expired", a.name, a.domain)
}

func (a *AccountController) renewable() error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()

	// do calculations
	newValidUntil := utils.SecondsToTime(a.account.ValidUntil).Add(a.conf.AccountRenewalPeriod)
	// renew count bumped because domain is already at count 1 when created
	renewCount := a.conf.AccountRenewalCountMax + 1
	// set new expected valid until
	maximumValidUntil := a.ctx.BlockTime().Add(a.conf.AccountRenewalPeriod * time.Duration(renewCount))
	// check if new valid until is after maximum allowed
	if newValidUntil.After(maximumValidUntil) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "unable to renew account %s in domain %s, maximum account renewal has exceeded: %s", *a.account.Name, a.domain, maximumValidUntil)
	}

	// if it has expired return error
	return nil
}

// ownedBy is the unexported function used by Owner
func (a *AccountController) ownedBy(addr sdk.AccAddress) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	// check if admin matches at least one address
	if a.account.Owner.Equals(addr) {
		return nil
	}

	return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to perform operation in the account owned by %s", addr, a.account.Owner)
}

// certNotExist is the unexported function used by CertificateNotExist
// and CertificateExists, it saves the index of the found certificate
// in indexPointer if it is not nil
func (a *AccountController) certNotExist(newCert []byte, indexPointer *int) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	// check if certificate is already present in account
	for i, cert := range a.account.Certificates {
		if bytes.Equal(cert, newCert) {
			if indexPointer != nil {
				*indexPointer = i
			}
			return sdkerrors.Wrapf(types.ErrCertificateExists, "certificate is already present")
		}
	}
	return nil
}

func (a *AccountController) certSizeNotExceeded(newCert []byte) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint64(len(newCert)) > a.conf.CertificateSizeMax {
		return sdkerrors.Wrapf(types.ErrCertificateSizeExceeded, "max certificate size %d exceeded", a.conf.CertificateSizeMax)
	}
	return nil
}

func (a *AccountController) certLimitNotExceeded() error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint32(len(a.account.Certificates)) >= a.conf.CertificateCountMax {
		return sdkerrors.Wrapf(types.ErrCertificateLimitReached, "max certificate limit %d reached, cannot add more", a.conf.CertificateCountMax)
	}
	return nil
}

func (a *AccountController) deletableBy(addr sdk.AccAddress) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check on a non existing domain is not allowed")
	}
	// get cached domain
	d := a.domainCtrl.Domain()
	if err := a.requireAccount(); err != nil {
		panic("validation check on a non existing account is not allowed")
	}
	switch d.Type {
	case types.ClosedDomain:
		if err := a.domainCtrl.
			Admin(addr).
			NotExpired().
			Validate(); err != nil {
			return err
		}
	case types.OpenDomain:
		if a.gracePeriodFinished() != nil {
			if a.ownedBy(addr) != nil {
				return sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to delete the account before grace period", a.account.Owner)
			}
		}
	}
	return nil
}

// validResources validates different resources
func (a *AccountController) validResources(resources []*types.Resource) error {
	a.requireConfiguration()
	validURI := regexp.MustCompile(a.conf.ValidURI)
	validResource := regexp.MustCompile(a.conf.ValidResource)
	// create resources set to identify duplicates
	sets := make(map[string]struct{}, len(resources))
	// iterate over resources to check their validity
	for _, resource := range resources {
		if resource == nil {
			return sdkerrors.Wrap(types.ErrInvalidResource, "nil resource provided")
		}
		// check if URI was already specified
		if _, ok := sets[resource.URI]; ok {
			return sdkerrors.Wrapf(types.ErrInvalidResource, "duplicate URI %s", resource.URI)
		}
		sets[resource.URI] = struct{}{}
		// is uri valid?
		if !validURI.MatchString(resource.URI) {
			return sdkerrors.Wrapf(types.ErrInvalidResource, "%s is not a valid URI", resource.URI)
		}
		// is resource valid?
		if !validResource.MatchString(resource.Resource) {
			return sdkerrors.Wrapf(types.ErrInvalidResource, "%s is not a valid resource", resource.Resource)
		}
	}
	// success
	return nil
}

func (a *AccountController) transferableBy(addr sdk.AccAddress) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	// check if domain has super user
	switch a.domainCtrl.Domain().Type {
	// if it has a super user then only domain admin can transfer accounts
	case types.ClosedDomain:
		if a.domainCtrl.
			Admin(addr).
			Validate() != nil {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "only domain admin %s is allowed to transfer accounts", a.domainCtrl.Domain().Admin)
		}
	// if it has not a super user then only account owner can transfer the account
	case types.OpenDomain:
		if a.ownedBy(addr) != nil {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "only account owner %s is allowed to transfer the account", a.account.Owner)
		}
	}
	return nil
}

func (a *AccountController) resettableBy(addr sdk.AccAddress, reset bool) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	d := a.domainCtrl.Domain()
	switch d.Type {
	case types.OpenDomain:
		if reset {
			if d.Admin.Equals(addr) {
				return sdkerrors.Wrapf(types.ErrUnauthorized, "domain admin is not authorized to reset account contents on open domains")
			}
		}
	case types.ClosedDomain:
	}
	return nil
}

// gracePeriodFinished is the condition that checks if given account's grace period has finished
func (a *AccountController) gracePeriodFinished() error {
	// require configuration
	a.requireConfiguration()
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("condition check not allowed on non existing account ")
	}
	// get grace period and expiration time
	gracePeriod := a.conf.AccountGracePeriod
	expireTime := utils.SecondsToTime(a.account.ValidUntil)
	if a.ctx.BlockTime().After(expireTime.Add(gracePeriod)) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrAccountGracePeriodNotFinished, "account %s grace period has not finished", *a.account.Name)
}

func (a *AccountController) resourceLimitNotExceeded(resources []*types.Resource) error {
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint32(len(resources)) > a.conf.ResourcesMax {
		return sdkerrors.Wrapf(types.ErrResourceLimitExceeded, "resource limit: %d", a.conf.ResourcesMax)
	}
	return nil
}

func (a *AccountController) metadataSizeNotExceeded(metadata string) error {
	// assert domain exists
	if err := a.requireAccount(); err != nil {
		panic("validation check is not allowed on a non existing account")
	}
	a.requireConfiguration()
	if uint64(len(metadata)) > a.conf.MetadataSizeMax {
		return sdkerrors.Wrapf(types.ErrMetadataSizeExceeded, "max metadata size %d exceeded", a.conf.MetadataSizeMax)
	}
	return nil
}

func (a *AccountController) registrableBy(addr sdk.AccAddress) error {
	if err := a.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check domain type
	switch a.domainCtrl.Domain().Type {
	// if domain is closed then the registerer must be domain owner
	case types.ClosedDomain:
		return a.domainCtrl.
			Admin(addr).
			Validate()
	default:
		return nil
	}
}

// Account returns the cached account, if the account existence
// was not asserted before, it panics.
func (a *AccountController) Account() types.Account {
	if err := a.requireAccount(); err != nil {
		panic("getting an account is not allowed before existence checks")
	}
	return *a.account
}
