package senarios

import (
	"context"
	"testing"

	utils "github.com/iov-one/starnamed/tests/starnamedutils"
	tools "github.com/iov-one/starnamed/tests/starnametesttools"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
)

// Module: Starname - Domain
// Tests Senarios: Domain
// *each version of starname chain
//
// - Generic Tests: inputs (open and close domain)
//     - Create
//     - Delete
//     - Transfer
//
// - Open domain
//     - Fail -> Create one account;
//              Delete Domain;

func ModuleStarnameTestSenarioDomain(t *testing.T, ctx context.Context, rep *testreporter.Reporter, chain *cosmos.CosmosChain) {
	t.Run("Domain", func(t *testing.T) {

		rep.TrackTest(t)

		// Create a new wallet
		wallets := utils.IBCWalletFactory(t, ctx, "StarnameTestSenarioDomain", 2, defaultTokenAmount, chain)

		owner := wallets[0]
		nonOwner := wallets[1]

		// Domains:
		domainIsOpen := []bool{true, false}

		t.Run("Create Domain", func(t *testing.T) {
			rep.TrackTest(t)

			for i := 0; i < len(domainIsOpen); i++ {
				command := tools.CommandBuilder(chain, true)
				_, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])
				if err != nil {
					t.Fatalf("error: %s", err)
				}

			}

		})

		t.Run("Transfer Domain", func(t *testing.T) {
			rep.TrackTest(t)

			// Domain to delete
			for i := 0; i < len(domainIsOpen); i++ {
				command := tools.CommandBuilder(chain, true)

				domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])
				if err != nil {
					t.Fatalf("error: %s", err)
				}

				command = tools.CommandBuilder(chain, true)

				domain = domain.Command(command, ctx)
				err = domain.TransferOwnership(nonOwner)

				if err != nil {
					t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
				}
			}

		})

		t.Run("Delete Domain", func(t *testing.T) {
			rep.TrackTest(t)

			// Domain to delete
			for i := 0; i < len(domainIsOpen); i++ {
				command := tools.CommandBuilder(chain, true)
				domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])
				if err != nil {
					t.Fatalf("error: %s", err)
				}

				command = tools.CommandBuilder(chain, true)
				domain = domain.Command(command, ctx)
				err = domain.Delete()

				if err != nil {
					t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
				}
			}

		})

	})

}
