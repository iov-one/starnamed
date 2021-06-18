package keeper_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/suite"
)

type KeeperSuite struct {
	BaseKeeperSuite
}

func (s *KeeperSuite) SetupTest() {
	s.Setup()
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
