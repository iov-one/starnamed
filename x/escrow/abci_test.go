package escrow_test

import (
	"reflect"
	"testing"

	"github.com/iov-one/starnamed/x/escrow"
	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

func TestBeginBlocker(t *testing.T) {
	keeper, ctx, _, _, _, _ := test.NewTestKeeper(nil)
	gen := test.NewEscrowGenerator(uint64(ctx.BlockTime().Unix()))

	normalEscrow, _ := gen.NewRandomTestEscrow()
	keeper.SaveEscrow(ctx, normalEscrow)

	expiredEscrow, _ := gen.NewRandomTestEscrow()
	expiredEscrow.Deadline = gen.NowAfter(0) - 10
	expiredEscrow.State = types.EscrowState_Expired
	keeper.SaveEscrow(ctx, expiredEscrow)

	expiringEscrow, _ := gen.NewRandomTestEscrow()
	expiringEscrow.Deadline = expiredEscrow.Deadline
	keeper.SaveEscrow(ctx, expiringEscrow)

	// Just to test if everything works as expected independently of position
	anotherNormalEscrow, _ := gen.NewRandomTestEscrow()
	keeper.SaveEscrow(ctx, anotherNormalEscrow)

	anotherExpiringEscrow, _ := gen.NewRandomTestEscrow()
	anotherExpiringEscrow.Deadline = expiringEscrow.Deadline
	keeper.SaveEscrow(ctx, anotherExpiringEscrow)

	escrow.BeginBlocker(ctx, keeper)

	if keeper.GetLastBlockTime(ctx) != gen.NowAfter(0) {
		t.Fatalf("Invalid last block time : expected %v, got %v", gen.NowAfter(0), keeper.GetLastBlockTime(ctx))
	}

	normalEscrows := []types.Escrow{normalEscrow, expiredEscrow, anotherNormalEscrow}
	for _, expected := range normalEscrows {
		actual, found := keeper.GetEscrow(ctx, expected.Id)
		if !found {
			t.Fatalf("Escrow %v not found after BeginBlocker call", expected.Id)
		}
		if !reflect.DeepEqual(expected, actual) {
			t.Fatalf("BeginBlocker modified a normal escrow : expected %v, got %v", expected, actual)
		}
	}

	expiringEscrows := []types.Escrow{expiringEscrow, anotherExpiringEscrow}
	for _, escrow := range expiringEscrows {
		actual, found := keeper.GetEscrow(ctx, escrow.Id)
		if !found {
			t.Fatalf("Escrow %v not found after BeginBlocker call", escrow.Id)
		}
		if actual.State != types.EscrowState_Expired {
			t.Fatalf("Escrow %v has not been marked expired when it should have, its state is %v", escrow.Id, actual.State)
		}
	}
}
