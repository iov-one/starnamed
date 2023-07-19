package senarios

import (
	"context"
	"testing"

	utils "github.com/iov-one/starnamed/tests/starnamedutils"
	tools "github.com/iov-one/starnamed/tests/starnametesttools"
	"github.com/strangelove-ventures/interchaintest/v7/chain/cosmos"
	"github.com/strangelove-ventures/interchaintest/v7/testreporter"
)

// Module: Starname - Escrow
// Tests Senarios: Escrow
// *each version

// * Must Set have Seccond IBC chain;
// * Set in configuration the other token;

// inputs: domain, account

//   - Generic Tests:
//     inputs:
//     > Tradable: open Domain, close Domain, open Account, Close Account
//   - Create Escrow with user X
//   - Refound:
//   	- Create Escrow with user X; User Y Buy Escrows

//   - IBC Tests:
//     > Denom: Default Chain Token, Custom IBC Token
//   - Create Escrow with user X; User Y Buy Escrows; User X Refound
func ModuleStarnameTestSenarioEscrow(t *testing.T, ctx context.Context, rep *testreporter.Reporter, chain *cosmos.CosmosChain) {
	t.Run("Escrow", func(t *testing.T) {

		rep.TrackTest(t)

		// Create a new wallet
		wallets := utils.IBCWalletFactory(t, ctx, "StarnameTestSenarioEscrow", 2, defaultTokenAmount, chain)

		owner := wallets[0]
		nonOwner := wallets[1]

		// Domains:
		domainIsOpen := []bool{true, false}

		t.Run("Generic", func(t *testing.T) {
			rep.TrackTest(t)
			t.Run("Create Escrow", func(t *testing.T) {
				rep.TrackTest(t)

				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					_, err = domain.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

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
					_, err = account.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

			})

			t.Run("Buy Escrow", func(t *testing.T) {
				rep.TrackTest(t)

				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					escrow, err := domain.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					err = escrow.Buy(nonOwner)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

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
					escrow, err := account.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					err = escrow.Buy(nonOwner)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

			})

			t.Run("Delete Escrow", func(t *testing.T) {
				rep.TrackTest(t)

				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					escrow, err := domain.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					err = escrow.Delete()

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

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
					escrow, err := account.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					err = escrow.Delete()

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}

			})

			t.Run("Update Escrow", func(t *testing.T) {
				rep.TrackTest(t)
				const (
					customEscrowPrice = int64(5_000)
				)
				for i := 0; i < len(domainIsOpen); i++ {
					command := tools.CommandBuilder(chain, true)
					domain, err := utils.NewStarnameDomain(command, ctx, chain, "", owner, domainIsOpen[i])

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					command = tools.CommandBuilder(chain, true)
					escrow, err := domain.Command(command, ctx).Escrow(defaultEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

					err = escrow.UpdatePrice(customEscrowPrice, chain.Config().Denom)

					if err != nil {
						t.Fatalf("Domain Kind %s, error: %s", domain.DomainKind(), err)
					}

				}
			})
		})

	})

}
