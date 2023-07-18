package cosmos_test

import (
	"context"
	"testing"

	starnametesttools "github.com/iov-one/starnamed/tests/starnameTestTools"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	NumberOfStarnameUsers int   = 10
	userFunds             int64 = int64(10_000_000_000)
	version
)

func TestStarnameModule(t *testing.T) {
	t.Parallel()
	// SDK v45 params
	shortVoteGenesis := []cosmos.GenesisKV{
		{
			Key:   "app_state.gov.voting_params.voting_period",
			Value: votingPeriod,
		},
		{
			Key:   "app_state.gov.deposit_params.max_deposit_period",
			Value: maxDepositPeriod,
		},
		{
			Key:   "app_state.gov.deposit_params.min_deposit.0.denom",
			Value: "uiov",
		},
	}

	cf := interchaintest.NewBuiltinChainFactory(zaptest.NewLogger(t), []*interchaintest.ChainSpec{
		{
			Name:      "starname",
			ChainName: "starname",
			Version:   "v0.12.0",
			ChainConfig: ibc.ChainConfig{
				Denom:         "uiov",
				ModifyGenesis: cosmos.ModifyGenesis(shortVoteGenesis),
				Images: []ibc.DockerImage{
					{
						Repository: "starnamed",
						Version:    "v0.12.0",
					},
				},
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)
	ctx := context.Background()

	client, network := interchaintest.DockerSetup(t)
	chain.UpgradeVersion(ctx, client, "starnamed", "v0.12.0")

	ic := interchaintest.NewInterchain().
		AddChain(chain)

	rep := testreporter.NewNopReporter()

	require.NoError(t, ic.Build(ctx, rep.RelayerExecReporter(t), interchaintest.InterchainBuildOptions{
		TestName:  t.Name(),
		Client:    client,
		NetworkID: network,
		// BlockDatabaseFile: interchaintest.DefaultBlockDatabaseFilepath(),
		SkipPathCreation: false,
	}))
	t.Cleanup(func() {
		_ = ic.Close()
	})

	ModuleStarnameTests(t, chain)
}

func ModuleStarnameTests(t *testing.T, chain *cosmos.CosmosChain) {
	ctx := context.Background()
	users := make([]*(starnametesttools.StarnameIBCWallet), NumberOfStarnameUsers)
	for i := 0; i < NumberOfStarnameUsers; i++ {
		starnameWallet := starnametesttools.StarnameIBCWallet{
			Wallet:       interchaintest.GetAndFundTestUsers(t, ctx, "starname_module", userFunds, chain)[0],
			Chain:        chain,
			StarDomains:  make([]string, 0),
			StarAccounts: make([]string, 0),
		}
		users[i] = &starnameWallet
	}

	t.Run("Starname module - Domain creation", func(t *testing.T) {

		command := starnametesttools.CommandBuilder(chain, true)

		for i := 0; i < NumberOfStarnameUsers; i++ {
			cmd := (*(starnametesttools.StarnameCommandTx))(command.Tx(users[i], true))
			ctx := context.Background()
			_, _, err := cmd.Starname().DomainRegister("").Exec(ctx)

			require.NoError(t, err, "Domain creation failed")
		}
	})

	t.Run("Starname module - Accounts", func(t *testing.T) {})

	t.Run("Starname module - Escrows", func(t *testing.T) {})

}
