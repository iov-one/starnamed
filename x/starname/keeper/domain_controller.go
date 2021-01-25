package keeper

import (
	"regexp"
	"time"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

// DomainControllerFunc is the function signature for domain validation functions
type DomainControllerFunc func(controller *DomainController) error

// DomainController is the domain controller
type DomainController struct {
	validators []DomainControllerFunc
	domainName string
	ctx        sdk.Context
	domain     *types.Domain
	conf       *configuration.Config
	store      *crud.Store
}

// NewDomainController is the constructor for Domain
// everything is processed sequentially, a wrong order of the sequence
// is forbidden, example: asserting domain expiration on a non existing
// domain causes a panic as it violates the condition scope of action.
func NewDomainController(ctx sdk.Context, domain string) *DomainController {
	return &DomainController{
		domainName: domain,
		ctx:        ctx,
	}
}

// WithDomains allows to specify a cached config
func (c *DomainController) WithDomains(store *crud.Store) *DomainController {
	c.store = store
	return c
}

// WithConfiguration allows to specify a cached config
func (c *DomainController) WithConfiguration(cfg configuration.Config) *DomainController {
	c.conf = &cfg
	return c
}

// WithDomain creates a domain controller with a cached domain
func (c *DomainController) WithDomain(dom types.Domain) *DomainController {
	c.domain = &dom
	c.domainName = dom.Name
	return c
}

// ---------------------- VALIDATION -----------------------------

// Validate validates a domain based on the provided checks
func (c *DomainController) Validate() error {
	for _, check := range c.validators {
		if err := check(c); err != nil {
			return err
		}
	}
	return nil
}

// Admin asserts that the domain owner is the provided address
func (c *DomainController) Admin(addr sdk.AccAddress) *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.isAdmin(addr)
	})
	return c
}

// NotExpired asserts that the domain has not expired
func (c *DomainController) NotExpired() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.notExpired()
	})
	return c
}

// Type makes sure the domain type is set to the provided condition
func (c *DomainController) Type(Type types.DomainType) *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.dType(Type)
	})
	return c
}

// MustExist checks if the provided domain exists
func (c *DomainController) MustExist() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.mustExist()
	})
	return c
}

// MustNotExist checks if the provided domain does not exist
func (c *DomainController) MustNotExist() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.mustNotExist()
	})
	return c
}

// ValidName checks if the name of the domain is valid
func (c *DomainController) ValidName() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.validName()
	})
	return c
}

// DeletableBy checks if the domain can be deleted by the provided address
func (c *DomainController) DeletableBy(addr sdk.AccAddress) *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.deletableBy(addr)
	})
	return c
}

// Transferable asserts the domain is transferable with the given flag
func (c *DomainController) Transferable(flag types.TransferFlag) *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.transferable(flag)
	})
	return c
}

// Renewable checks if the domain is allowed to be renewed
func (c *DomainController) Renewable() *DomainController {
	c.validators = append(c.validators, func(controller *DomainController) error {
		return controller.renewable()
	})
	return c
}

// expired returns nil if domain expired, otherwise ErrDomainNotExpired
func (c *DomainController) expired() error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("conditions check not allowed on non existing domain")
	}
	expireTime := utils.SecondsToTime(c.domain.ValidUntil)
	// if expire time is before block time means domain expired
	if expireTime.Before(c.ctx.BlockTime()) {
		return nil
	}

	return sdkerrors.Wrapf(types.ErrDomainNotExpired, "domain %s has not expired", c.domain.Name)
}

// gracePeriodFinished is the condition that checks if given domain's grace period has finished
func (c *DomainController) gracePeriodFinished() error {
	// require configuration
	c.requireConfiguration()
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("condition check not allowed on non existing domain")
	}
	// get grace period and expiration time
	gracePeriod := c.conf.DomainGracePeriod
	expireTime := utils.SecondsToTime(c.domain.ValidUntil)
	if c.ctx.BlockTime().After(expireTime.Add(gracePeriod)) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrDomainGracePeriodNotFinished, "domain %s grace period has not finished", c.domain.Name)
}

// isAdmin makes sure the domain is owned by the provided address
func (c *DomainController) isAdmin(addr sdk.AccAddress) error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check if admin matches addr
	if c.domain.Admin.Equals(addr) {
		return nil
	}
	return sdkerrors.Wrapf(types.ErrUnauthorized, "%s is not allowed to perform an operation in a domain owned by %s", addr, c.domain.Admin)
}

func (c *DomainController) notExpired() error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	// check if domain has expired
	expireTime := utils.SecondsToTime(c.domain.ValidUntil)
	// if block time is before expiration, return nil
	if c.ctx.BlockTime().Before(expireTime) {
		return nil
	}
	// if it has expired return error
	return sdkerrors.Wrapf(types.ErrDomainExpired, "%s has expired", c.domainName)
}

func (c *DomainController) dType(Type types.DomainType) error {
	// assert domain exists
	if err := c.requireDomain(); err != nil {
		panic("validation check is not allowed on a non existing domain")
	}
	if c.domain.Type != Type {
		return sdkerrors.Wrapf(types.ErrInvalidDomainType, "operation not allowed on invalid domain type %s, expected %s", c.domain.Type, Type)
	}
	return nil
}

// requireDomain tries to find the domain by name
// if it is not found then an error is returned
func (c *DomainController) requireDomain() error {
	if c.domain != nil {
		return nil
	}
	if c.store == nil {
		panic("store is missing")
	}
	domain := new(types.Domain)
	if err := (*c.store).Read((&types.Domain{Name: c.domainName}).PrimaryKey(), domain); err != nil {
		return sdkerrors.Wrapf(types.ErrDomainDoesNotExist, "not found: %s", c.domainName)
	}
	c.domain = domain
	return nil
}

// mustExist checks if a domain exists
func (c *DomainController) mustExist() error {
	return c.requireDomain()
}

// mustNotExist asserts that a domain does not exist
func (c *DomainController) mustNotExist() error {
	err := c.requireDomain()
	if err == nil {
		return sdkerrors.Wrapf(types.ErrDomainAlreadyExists, c.domainName)
	}
	return nil
}

// validName checks if the name of the domain is valid
func (c *DomainController) validName() error {
	if c.conf == nil {
		panic("configuration is missing")
	}
	// get valid domain regexp
	validator := regexp.MustCompile(c.conf.ValidDomainName)
	// assert domain name validity
	if !validator.MatchString(c.domainName) {
		return sdkerrors.Wrap(types.ErrInvalidDomainName, c.domainName)
	}
	// success
	return nil
}

// requireConfiguration updates the configuration
// if it is not already set, and caches it after
func (c *DomainController) requireConfiguration() {
	if c.conf == nil {
		panic("configuration is missing")
	}
}

// deletableBy is the underlying operation used by DeletableBy controller
func (c *DomainController) deletableBy(addr sdk.AccAddress) error {
	if err := c.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	switch c.domain.Type {
	case types.ClosedDomain:
		// check if either domain is owned by provided address or if grace period is finished
		if err := c.gracePeriodFinished(); err != nil {
			if err := c.isAdmin(addr); err != nil {
				return sdkerrors.Wrap(types.ErrUnauthorized, "only admin delete domain before grace period is finished")
			}
		}
	case types.OpenDomain:
		if err := c.gracePeriodFinished(); err != nil {
			return sdkerrors.Wrap(types.ErrDomainGracePeriodNotFinished, "cannot delete open domain before grace period is finished")
		}
	}
	return nil
}

func (c *DomainController) transferable(flag types.TransferFlag) error {
	if err := c.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	switch c.domain.Type {
	case types.OpenDomain:
		if flag != types.TransferResetNone {
			return sdkerrors.Wrapf(types.ErrUnauthorized, "unable to transfer open domain %s with flag %d", c.domainName, flag)
		}
		return nil
	default:
		return nil
	}
}

func (c *DomainController) renewable() error {
	c.requireConfiguration()
	if err := c.requireDomain(); err != nil {
		panic("validation check not allowed on a non existing domain")
	}
	// check if current time is not beyond grace period
	renewalDeadline := utils.SecondsToTime(c.domain.ValidUntil).Add(c.conf.DomainGracePeriod)
	if c.ctx.BlockTime().After(renewalDeadline) {
		return sdkerrors.Wrapf(types.ErrRenewalDeadlineExceeded, "the deadline for renewal was: %s, current time is: %s, please delete the domain and re-register", renewalDeadline, c.ctx.BlockTime())
	}
	// do calculations
	newValidUntil := utils.SecondsToTime(c.domain.ValidUntil).Add(c.conf.DomainRenewalPeriod) // set new expected valid until
	// renew count bumped because domain count is already at count 1 when created
	renewCount := c.conf.DomainRenewalCountMax + 1
	maximumValidUntil := c.ctx.BlockTime().Add(c.conf.DomainRenewalPeriod * time.Duration(renewCount))
	// check if new valid until is after maximum allowed
	if newValidUntil.After(maximumValidUntil) {
		return sdkerrors.Wrapf(types.ErrUnauthorized, "unable to renew domain, domain %s renewal period would be after maximum allowed: %s", c.domainName, maximumValidUntil)
	}
	// success
	return nil
}

// Domain returns a copy the domain, panics if the operation is done without
// doing validity checks on domain existence as it is not an allowed op
func (c *DomainController) Domain() types.Domain {
	if err := c.requireDomain(); err != nil {
		panic("get domain without running existence checks is not allowed")
	}
	return *c.domain
}
