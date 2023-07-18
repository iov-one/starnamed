package starnametesttools

import (
	"context"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

// ############################## Commands ##############################

type Command struct {
	chain      *cosmos.CosmosChain
	log_format string
	bin        string
	pre_args   []string
	args       []string
	post_args  []string
	env        []string
}

func CommandBuilder(chain *cosmos.CosmosChain, default_args bool) (cmd *Command) {

	cmd = &Command{
		chain:      chain,
		log_format: "json",
		bin:        "starnamed",
	}

	if default_args {
		cmd.Default_args()
	}

	return
}

func (c *Command) Default_args() *Command {
	c.pre_args = []string{
		"--keyring-backend", keyring.BackendTest,
		"--node", c.chain.GetRPCAddress(),
		"--gas-prices", c.chain.Config().GasPrices,
		"--home", c.chain.HomeDir(),
		"--chain-id", c.chain.Config().ChainID,
	}

	return c
}

func (c *Command) SetPreArgs(args ...string) *Command {
	c.pre_args = args
	return c
}

func (c *Command) SetArgs(args ...string) *Command {
	c.args = args
	return c
}

func (c *Command) SetPostArgs(args ...string) *Command {
	c.post_args = args
	return c
}

func (c *Command) SetEnv(env ...string) *Command {
	c.env = env
	return c
}

func (c *Command) build() []string {

	var cmd []string = []string{
		c.bin,
	}

	cmd = append(cmd, c.pre_args...)
	cmd = append(cmd, c.args...)
	cmd = append(cmd, c.post_args...)

	return cmd
}

func (c *Command) Exec(ctx context.Context) (std_out []byte, std_err []byte, err error) {

	cmd := c.build()

	std_out, std_err, err = c.chain.Exec(ctx, cmd, c.env)

	return
}

// ############################## Cosmos Commands ##############################

// Cosmos Commands

func (c *Command) Tx(user ibc.Wallet, autoAcceptTx bool) *Command {
	c.args = append(c.args, "tx")

	c.pre_args = append(c.pre_args, "--from", user.KeyName())

	if autoAcceptTx {
		c.post_args = append(c.post_args, "--yes")
	}

	return c
}

func (c *Command) Query() *Command {
	c.args = append(c.args, "query")
	return c
}
