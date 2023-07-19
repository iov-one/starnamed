package starnametesttools

import (
	"context"

	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

// ############################## StarnameTx ##############################

type StarnameCommandTx Command

func (c *StarnameCommandTx) Starname() *StarnameCommandTx {
	c.args = append(c.args, "starname")
	return c
}

// If the domainName is not set, the command will generate a random domain name
func (c *StarnameCommandTx) DomainRegister(domainName string, flags ...string) *StarnameCommandTx {

	if domainName == "" {
		domainName = randomString(10)
	}

	c.args = append(c.args, "domain-register", "--domain", domainName)

	if len(flags) > 0 {
		c.args = append(c.args, flags...)
	}

	return c
}

func (c *StarnameCommandTx) Exec(ctx context.Context) (std_out []byte, std_err []byte, err error) {
	return (*Command)(c).Exec(ctx)
}

// ############################## StarnameQuery ##############################

type StarnameCommandQuery Command

func (c *StarnameCommandQuery) Starname() *StarnameCommandQuery {
	c.args = append(c.args, "starname")
	return c
}

func (c *StarnameCommandQuery) DomainsByOwner(user ibc.Wallet) *StarnameCommandQuery {
	c.args = append(c.args, "domains-by-owner", "--address", string(user.Address()))
	return c
}

func (c *StarnameCommandQuery) Exec(ctx context.Context) (std_out []byte, std_err []byte, err error) {
	return (*Command)(c).Exec(ctx)
}
