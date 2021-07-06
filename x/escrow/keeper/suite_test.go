package keeper_test

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type BaseKeeperSuite struct {
	suite.Suite
	keeper    keeper.Keeper
	msgServer types.MsgServer
	ctx       sdk.Context
	generator *test.EscrowGenerator
	store     crud.Store
	storeKey  sdk.StoreKey
	balances  map[string]sdk.Coins
}

func (s *BaseKeeperSuite) Setup(coinHolders []sdk.AccAddress) {
	test.SetConfig()
	s.keeper, s.ctx, s.store, s.balances, s.storeKey = test.NewTestKeeper(coinHolders)
	s.keeper.ImportNextID(s.ctx, 1)
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.generator = test.NewEscrowGenerator(uint64(s.ctx.BlockTime().Unix()))
}
