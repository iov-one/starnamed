package cosmos_test

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
)

func randomString(length int) string {
	rand.NewSource(time.Now().UnixNano())

	var alphabet string = "abcdefghijklmnopqrstuvwxyz"
	var sb strings.Builder

	l := len(alphabet)

	for i := 0; i < length; i++ {
		c := alphabet[rand.Intn(l)]
		sb.WriteByte(c)
	}

	return sb.String()
}

type StarnameIBCWallet struct {
	wallet       ibc.Wallet
	chain        *cosmos.CosmosChain
	StarDomains  []string
	StarAccounts []string
}

// ############################## Wallet ##############################
func (w *StarnameIBCWallet) KeyName() string {
	return w.wallet.KeyName()
}

func (w *StarnameIBCWallet) FormattedAddress() string {
	return w.wallet.FormattedAddress()
}

func (w *StarnameIBCWallet) Mnemonic() string {
	return w.wallet.Mnemonic()
}

func (w *StarnameIBCWallet) Address() []byte {
	return w.wallet.Address()
}

func (w *StarnameIBCWallet) GetIBCWallet() ibc.Wallet {
	return w.wallet
}

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

// ############################## Starname ##############################

type StarnameCommand Command

func (c *StarnameCommand) starname() *StarnameCommand {
	c.args = append(c.args, "starname")
	return c
}

// If the domainName is not set, the command will generate a random domain name
func (c *StarnameCommand) DomainRegister(domainName string, flags ...string) *StarnameCommand {

	if domainName == "" {
		domainName = randomString(10)
	}

	c.args = append(c.args, "domain-register", "--domain", domainName)

	if len(flags) > 0 {
		c.args = append(c.args, flags...)
	}

	return c
}

func (c *StarnameCommand) Exec(ctx context.Context) (std_out []byte, std_err []byte, err error) {
	return (*Command)(c).Exec(ctx)
}
