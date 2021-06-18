package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type InvariantsSuite struct {
	BaseKeeperSuite
}

func (s *InvariantsSuite) reinitEscrows() {
	s.Setup()
	// Add 100 valid escrow
	for i := 0; i < 100; i++ {
		escrow, _ := s.generator.NewRandomTestEscrow()
		s.keeper.SaveEscrow(s.ctx, escrow)
	}
	// Add 20 valid expired escrows
	for i := 0; i < 20; i++ {
		escrow, _ := s.generator.NewRandomTestEscrow()
		escrow.State = types.EscrowState_Expired
		s.keeper.SaveEscrow(s.ctx, escrow)
	}

	//Invariant should hold
	s.expects(keeper.StateInvariant(s.keeper), false)
}

func (s *InvariantsSuite) SetupTest() {
	s.reinitEscrows()
}

func TestInvariants(t *testing.T) {
	suite.Run(t, new(InvariantsSuite))

}

func (s *InvariantsSuite) expects(invariant sdk.Invariant, shouldBeBroken bool) {
	errorMsg, broken := invariant(s.ctx)

	if broken != shouldBeBroken {
		errorStr := "no "
		if shouldBeBroken {
			errorStr = ""
		}
		s.Failf("Expected %{s}error but got %s", errorStr, errorMsg)
	}
}

func (s *InvariantsSuite) TestStateInvariant() {
	invariant := keeper.StateInvariant(s.keeper)

	// Test invalid expired escrows
	escrow, _ := s.generator.NewRandomTestEscrow()
	// Escrow is not actually expired
	escrow.State = types.EscrowState_Expired
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	escrow, _ = s.generator.NewRandomTestEscrow()
	// Escrow should be expired
	escrow.Deadline = s.generator.NowAfter(0) - 1
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	// Test completed and refunded escrows
	escrow, _ = s.generator.NewRandomTestEscrow()
	escrow.State = types.EscrowState_Completed
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	escrow.State = types.EscrowState_Refunded
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	//Test object not in crud store / not owned by module / not owned by module in store
	//TODO

	//Test id on the future
	//TODOg

}
