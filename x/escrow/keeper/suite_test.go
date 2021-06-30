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
}

func (s *BaseKeeperSuite) Setup() {
	test.SetConfig()
	s.keeper, s.ctx, s.store, _, s.storeKey = test.NewTestKeeper(nil)
	s.keeper.ImportNextID(s.ctx, 1)
	s.msgServer = keeper.NewMsgServerImpl(s.keeper)
	s.generator = test.NewEscrowGenerator(uint64(s.ctx.BlockTime().Unix()))
}
