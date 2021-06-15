package escrow_test

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/simapp"
	"github.com/stretchr/testify/suite"

	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	tmtime "github.com/tendermint/tendermint/types/time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/app"
	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type GenesisTestSuite struct {
	suite.Suite

	cdc    codec.JSONMarshaler
	app    *app.WasmApp
	ctx    sdk.Context
	keeper *keeper.Keeper
}

func (suite *GenesisTestSuite) SetupTest() {
	app := app.Setup(false)
	suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.cdc = app.AppCodec()
	suite.app = app
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestGenesisState() {
	type GenState func(*EscrowGenerator) *types.GenesisState

	testCases := []struct {
		name       string
		genState   GenState
		expectPass bool
	}{{
		name: "empty",
		genState: func(gen *EscrowGenerator) *types.GenesisState {
			return gen.NewEmptyEscrowGenesis()
		},
		expectPass: true,
	}, {
		name: "import atomic htlcs and asset supplies",
		genState: func() *types.GenesisState {
			gs := NewHTLTGenesis(suite.addrs[0])
			_, addrs := GeneratePrivKeyAddressPairs(2)
			var htlcs []types.HTLC
			var supplies []types.AssetSupply
			for i := 0; i < 2; i++ {
				htlc, supply := loadSwapAndSupply(addrs[i], i)
				htlcs = append(htlcs, htlc)
				supplies = append(supplies, supply)
			}
			gs.Htlcs = htlcs
			gs.Supplies = supplies
			return gs
		},
		expectPass: true,
	},
	}

	for _, tc := range testCases {
		suite.Run(
			tc.name,
			func() {
				if tc.expectPass {
					suite.NotPanics(
						func() {
							suite.app.Init
						},
						tc.name,
					)
				} else {
					suite.Panics(
						func() {
							simapp.SetupWithGenesisHTLC(tc.genState())
						},
						tc.name,
					)
				}
			},
		)
	}
}
