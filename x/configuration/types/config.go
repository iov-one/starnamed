package types

import (
	"fmt"
	"regexp"

	"github.com/cosmos/cosmos-sdk/types"
)

// Validate validates the Config object.
func (c Config) Validate() error {
	if c.Configurer == "" {
		return fmt.Errorf("empty configurer")
	}
	if c.DomainRenewalPeriod < 0 {
		return fmt.Errorf("empty domain renew")
	}
	if c.DomainGracePeriod < 0 {
		return fmt.Errorf("empty domain grace period")
	}
	if c.AccountRenewalPeriod < 0 {
		return fmt.Errorf("empty account renew")
	}
	if c.AccountGracePeriod < 0 {
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
	if c.EscrowMaxPeriod < 0 {
		return fmt.Errorf("empty escrow maximum duration")
	}
	if c.EscrowBroker == "" {
		return fmt.Errorf("empty escrow broker")
	}
	if _, err := types.AccAddressFromBech32(c.EscrowBroker); err != nil {
		return fmt.Errorf("invalid escrow broker address : %s", err)
	}
	if c.EscrowCommission.LT(types.ZeroDec()) || c.EscrowCommission.GT(types.OneDec()) {
		return fmt.Errorf("invalid escrow commission: not in interval [0;1]")
	}

	return nil
}
