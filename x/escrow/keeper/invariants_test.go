package keeper_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type InvariantsSuite struct {
	BaseKeeperSuite
}

func (s *InvariantsSuite) saveObject(obj *TestObject) {
	// Escrows are already created, so all your objects are belong to us
	obj.Owner = s.keeper.GetEscrowAddress()
	// Save the object to the store
	s.Assert().Nil(s.store.Create(obj), "Error while saving a generated object")
}

func (s *InvariantsSuite) reinitEscrows() {
	s.Setup()
	// Add 100 valid escrow
	for i := 0; i < 100; i++ {
		escrow, obj := s.generator.NewRandomTestEscrow()
		s.keeper.SaveEscrow(s.ctx, escrow)
		s.saveObject(obj)
	}
	// Add 20 valid expired escrows
	for i := 0; i < 20; i++ {
		escrow, obj := s.generator.NewRandomTestEscrow()
		escrow.State = types.EscrowState_Expired
		s.keeper.SaveEscrow(s.ctx, escrow)
		s.saveObject(obj)
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
	escrow, obj := s.generator.NewRandomTestEscrow()
	// Escrow is not actually expired
	escrow.State = types.EscrowState_Expired
	s.keeper.SaveEscrow(s.ctx, escrow)
	s.saveObject(obj)

	s.expects(invariant, true)
	s.reinitEscrows()

	escrow, obj = s.generator.NewRandomTestEscrow()
	// Escrow should be expired
	escrow.Deadline = s.generator.NowAfter(0) - 1
	s.keeper.SaveEscrow(s.ctx, escrow)
	s.saveObject(obj)

	s.expects(invariant, true)
	s.reinitEscrows()

	// Test completed and refunded escrows
	escrow, obj = s.generator.NewRandomTestEscrow()
	escrow.State = types.EscrowState_Completed
	s.keeper.SaveEscrow(s.ctx, escrow)
	s.saveObject(obj)

	s.expects(invariant, true)
	s.reinitEscrows()

	escrow.State = types.EscrowState_Refunded
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	//Test object not in crud store
	escrow, obj = s.generator.NewRandomTestEscrow()
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	// Test object not owned by module / not owned by module in store
	escrow, obj = s.generator.NewRandomTestEscrow()
	obj.Owner = s.generator.NewAccAddress()
	s.keeper.SaveEscrow(s.ctx, escrow)
	s.Assert().Nil(s.store.Create(obj))

	s.expects(invariant, true)
	s.reinitEscrows()

	// Object is valid but not his version in store
	obj.Owner = s.keeper.GetEscrowAddress()
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.expects(invariant, true)
	s.reinitEscrows()

	//Test id on the future
	escrow, obj = s.generator.NewRandomTestEscrow()
	escrow.Id = hex.EncodeToString(sdk.Uint64ToBigEndian(s.generator.nextId + 500))
	s.keeper.SaveEscrow(s.ctx, escrow)
	s.saveObject(obj)

	s.expects(invariant, true)
}
