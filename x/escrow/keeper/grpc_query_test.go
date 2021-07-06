package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type QueryTestSuite struct {
	BaseKeeperSuite
}

func (s *QueryTestSuite) SetupTest() {
	s.Setup(nil)
}

func TestQueryTestSuite(t *testing.T) {
	suite.Run(t, new(QueryTestSuite))
}

func (s *QueryTestSuite) TestQueryEscrow() {
	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	existingEscrow, _ := s.generator.NewTestEscrow(s.generator.NewAccAddress(), price, s.generator.NowAfter(10))
	expiredEscrow, _ := s.generator.NewTestEscrow(s.generator.NewAccAddress(), price, s.generator.NowAfter(0)-10)

	s.keeper.SaveEscrow(s.ctx, existingEscrow)
	s.keeper.SaveEscrow(s.ctx, expiredEscrow)

	existingEscrowId := existingEscrow.Id
	expiredEscrowId := expiredEscrow.Id
	nonExistingEscrowId := "AABBCCDDEEFF1122"

	// Normal
	resp, err := s.keeper.Escrow(sdk.WrapSDKContext(s.ctx), &types.QueryEscrowRequest{Id: existingEscrowId})

	s.Assert().Nil(err)
	s.Assert().Equal(existingEscrow, *resp.Escrow)

	// Expired
	resp, err = s.keeper.Escrow(sdk.WrapSDKContext(s.ctx), &types.QueryEscrowRequest{Id: expiredEscrowId})

	s.Assert().Nil(err)
	s.Assert().Equal(expiredEscrow, *resp.Escrow)

	// Does not exist
	resp, err = s.keeper.Escrow(sdk.WrapSDKContext(s.ctx), &types.QueryEscrowRequest{Id: nonExistingEscrowId})

	s.Assert().NotNil(err)
	s.Assert().ErrorIs(err, types.ErrEscrowNotFound)
}
