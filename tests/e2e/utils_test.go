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
	rand.Seed(time.Now().UnixNano())

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

func baseCMD(chain *cosmos.CosmosChain, user ibc.Wallet, args ...string) []string {
	default_args := []string{
		"--keyring-backend", keyring.BackendTest,
		"--node", chain.GetRPCAddress(),
		"--from", user.KeyName(),
		"--gas-prices", chain.Config().GasPrices,
		"--home", chain.HomeDir(),
		"--chain-id", chain.Config().ChainID,
	}

	cmd := append([]string{chain.Config().Bin}, args...)
	cmd = append(cmd, default_args...)
	return cmd
}

func Starname_CreateDomain(chain *cosmos.CosmosChain, user *StarnameIBCWallet, domain string) error {

	ctx := context.Background()

	if domain == "" {
		domain = randomString(10)
	}

	cmd := baseCMD(chain, user.GetIBCWallet(),
		"tx", "starname", "domain-register",
		"-d", domain,
		"--yes", "--log_format", "json",
	)

	std_out, std_err, err := chain.Exec(ctx, cmd, []string{})

	_, _ = std_out, std_err

	return err
}
