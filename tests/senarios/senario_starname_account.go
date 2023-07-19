package senarios

import (
	"context"
	"testing"

	utils "github.com/iov-one/starnamed/tests/starnamedutils"
	tools "github.com/iov-one/starnamed/tests/starnametesttools"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
)

// Module: Starname - Account
// Tests Senarios: Accounts
// *each version

// - Generic Tests: inputs (open and close domain)
//     - Create Account - Same onwer
//     - Delete Account - Same onwer
//     - Transfer - New onwer

// - Open domain
// 	   - Create Account - Domain owner X; User Y create account;
//     - Fail -> Domain owner X; User Y create account; Domain owner X try to delete;

// - Create Close domain
//     - Domain owner X; User X create account; Transfer to User Y; Domain owner X try to delete;
//     - Fail -> Domain onwer X; User Y create account;

func ModuleStarnameTestSenarioAccount(t *testing.T, ctx context.Context, rep *testreporter.Reporter, chain *cosmos.CosmosChain) {
	t.Run("Account", func(t *testing.T) {

		rep.TrackTest(t)

		// Create a new wallet
		wallets := utils.IBCWalletFactory(t, ctx, "StarnameTestSenarioAccount", 2, defaultTokenAmount, chain)

		owner := wallets[0]
		nonOwner := wallets[1]

		// Domains:
		domainIsOpen := []bool{true, false}

		t.Run("Generic", func(t *testing.T) {
			rep.TrackTest(t)
			t.Run("Create Account", func(t *testing.T) {
				rep.TrackTest(t)

				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					_, err = utils.NewStarnameAccount(command, ctx, chain, "", owner, *domain)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

			})

			t.Run("Transfer Account", func(t *testing.T) {
				rep.TrackTest(t)

				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					account, err := utils.NewStarnameAccount(command, ctx, chain, "", owner, *domain)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					err = account.Command(command, ctx).TransferOwnership(nonOwner)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

			})

			t.Run("Delete Account", func(t *testing.T) {
				rep.TrackTest(t)

				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					account, err := utils.NewStarnameAccount(command, ctx, chain, "", owner, *domain)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					err = account.Command(command, ctx).Delete()

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

			})
		})

	})

}
