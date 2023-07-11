package cosmos_test

import (
	"context"
	"testing"

	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

const (
	NumberOfStarnameUsers int   = 100
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
			ChainConfig: ibc.ChainConfig{ //TODO: Enforce the docker image
				Denom:         "uiov",
				ModifyGenesis: cosmos.ModifyGenesis(shortVoteGenesis),
			},
		},
	})

	chains, err := cf.Chains(t.Name())
	require.NoError(t, err)

	chain := chains[0].(*cosmos.CosmosChain)

	client, network := interchaintest.DockerSetup(t)
	ic := interchaintest.NewInterchain().
		AddChain(chain)

	ctx := context.Background()

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
	users := make([]*StarnameIBCWallet, NumberOfStarnameUsers)
	for i := 0; i < NumberOfStarnameUsers; i++ {
		starnameWallet := StarnameIBCWallet{
			wallet:       interchaintest.GetAndFundTestUsers(t, ctx, "starname_module", userFunds, chain)[0],
			chain:        chain,
			StarDomains:  make([]string, 0),
			StarAccounts: make([]string, 0),
		}
		users[i] = &starnameWallet
	}
	// Domain creation

	for i := 0; i < NumberOfStarnameUsers; i++ {
		Starname_CreateDomain(chain, users[i], "")
	}

}
