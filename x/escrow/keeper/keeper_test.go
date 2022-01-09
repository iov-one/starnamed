package keeper_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/stretchr/testify/suite"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/test"
)

type KeeperSuite struct {
	BaseKeeperSuite
}

func (s *KeeperSuite) SetupTest() {
	s.Setup(nil, true)
}

func TestKeeper(t *testing.T) {
	suite.Run(t, new(KeeperSuite))
}

func (s *KeeperSuite) TestNextID() {
	stringify := func(id uint64) string {
		return hex.EncodeToString(sdk.Uint64ToBigEndian(id))
	}

	const InitialID = uint64(42)
	var id = InitialID
	// Test import id
	s.keeper.ImportNextID(s.ctx, id)

	// Check that FetchNextId and export value match
	s.Assert().Equal(InitialID, s.keeper.GetNextIDForExport(s.ctx))
	s.Assert().Equal(stringify(InitialID), s.keeper.FetchNextId(s.ctx))

	// Go to next value
	s.keeper.NextId(s.ctx)

	// Check thant values match
	s.Assert().Equal(InitialID+1, s.keeper.GetNextIDForExport(s.ctx))
	s.Assert().Equal(stringify(InitialID+1), s.keeper.FetchNextId(s.ctx))

}

func TestModuleDisabled(t *testing.T) {

	k, ctx, _, _, _, _ := test.NewTestKeeper(nil, false)
	gen := test.NewEscrowGenerator(10)

	testCases := []struct {
		name       string
		callMethod func(k keeper.Keeper)
	}{
		{
			name: "should panic upon creation",
			callMethod: func(k keeper.Keeper) {
				escrow, obj := gen.NewRandomTestEscrow()
				seller, _ := sdk.AccAddressFromBech32(escrow.Seller)
				_, _ = k.CreateEscrow(ctx, seller, escrow.Price, obj, escrow.Deadline, rand.Bool())
			},
		},
		{
			name: "should panic upon update",
			callMethod: func(k keeper.Keeper) {
				_ = k.UpdateEscrow(ctx, "0001", gen.NewAccAddress(), gen.NewAccAddress(), nil, 0)
			},
		},
		{
			name: "should panic upon transfer",
			callMethod: func(k keeper.Keeper) {
				_ = k.TransferToEscrow(ctx, gen.NewAccAddress(), "00001", sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.OneInt())))
			},
		},
		{
			name: "should panic upon refund",
			callMethod: func(k keeper.Keeper) {
				_ = k.RefundEscrow(ctx, gen.NewAccAddress(), "00001")
			},
		},
		{
			name: "should panic upon completing auction",
			callMethod: func(k keeper.Keeper) {
				_ = k.CompleteAuction(ctx, "00001")
			},
		},
		{
			name: "should panic upon query",
			callMethod: func(k keeper.Keeper) {
				k.GetEscrow(ctx, "0000000000000001")
			},
		},
		{
			name: "should panic upon read",
			callMethod: func(k keeper.Keeper) {
				_, _ = k.QueryEscrows(ctx, func(statement crud.QueryStatement) crud.ValidQuery { return statement })
			},
		},
	}

	for _, testCase := range testCases {
		test.EvaluateTest(t, testCase.name, func(*testing.T) error {
			testCase.callMethod(k)
			return nil
		})
	}
}
