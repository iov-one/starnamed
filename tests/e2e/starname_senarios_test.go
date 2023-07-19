package cosmos_test

import (
	"context"
	"testing"

	senarios "github.com/iov-one/starnamed/tests/senarios"
	"github.com/strangelove-ventures/interchaintest/v7"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/ibc"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
	"github.com/stretchr/testify/require"
	"go.uber.org/zap/zaptest"
)

func TestStarnameSenarios(t *testing.T) {
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

	// Run the tests
	senarios.ModuleStarnameTestSenarioDomain(t, ctx, rep, chain)
}
