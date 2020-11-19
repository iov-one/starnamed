package types

import (
	"fmt"
	"regexp"
)

// Validate validates the Config object.
func (c Config) Validate() error {
	if c.Configurer == nil {
		return fmt.Errorf("empty configurer")
	}
	if c.DomainRenewalPeriod.Seconds < 0 {
		return fmt.Errorf("empty domain renew")
	}
	if c.DomainGracePeriod.Seconds < 0 {
		return fmt.Errorf("empty domain grace period")
	}
	if c.AccountRenewalPeriod.Seconds < 0 {
		return fmt.Errorf("empty account renew")
	}
	if c.AccountGracePeriod.Seconds < 0 {
		return fmt.Errorf("empty account grace period")
	}
	if _, err := regexp.Compile(c.ValidAccountName); err != nil {
		return err
	}
	if _, err := regexp.Compile(c.ValidResource); err != nil {
		return err
	}
	if _, err := regexp.Compile(c.ValidURI); err != nil {
		return err
	}
	if _, err := regexp.Compile(c.ValidDomainName); err != nil {
		return err
	}

	return nil
}
