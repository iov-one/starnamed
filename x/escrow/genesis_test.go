package escrow_test

import (
	"fmt"
	"testing"

	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/stretchr/testify/suite"
	abci "github.com/tendermint/tendermint/abci/types"
	"github.com/tendermint/tendermint/libs/rand"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/app"
	"github.com/iov-one/starnamed/x/escrow"
	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type GenesisTestSuite struct {
	suite.Suite
	cdc       codec.JSONCodec
	app       *app.WasmApp
	ctx       sdk.Context
	storeKey  sdk.StoreKey
	crudStore crud.Store
	keeper    keeper.Keeper
	gen       *test.EscrowGenerator
}

func (suite *GenesisTestSuite) SetupTest() {
	test.SetConfig()

	suite.keeper, suite.ctx, suite.crudStore, _, suite.storeKey, _ = test.NewTestKeeper(nil, true)
	suite.keeper.ImportNextID(suite.ctx, 1)
	app := app.Setup(false)
	//suite.ctx = app.BaseApp.NewContext(false, tmproto.Header{Height: 1, Time: tmtime.Now()})

	suite.cdc, _ = test.NewTestCodec()
	suite.app = app
	suite.gen = test.NewEscrowGenerator(uint64(suite.ctx.BlockTime().Unix()))
}

func TestGenesisTestSuite(t *testing.T) {
	suite.Run(t, new(GenesisTestSuite))
}

func (suite *GenesisTestSuite) TestGenesisImport() {
	type GenState func(*test.EscrowGenerator) *types.GenesisState

	testCases := []struct {
		name     string
		genState GenState
	}{{
		name: "empty",
		genState: func(gen *test.EscrowGenerator) *types.GenesisState {
			return gen.NewEmptyEscrowGenesis()
		},
	}, {
		name: "import escrows",
		genState: func(gen *test.EscrowGenerator) *types.GenesisState {
			return gen.NewEscrowGenesis(1000)
		},
	},
	}

	for _, tc := range testCases {
		test.EvaluateTest(suite.T(), tc.name, func(t *testing.T) error {
			state := tc.genState(suite.gen)
			stateBytes := suite.cdc.MustMarshalJSON(state)
			suite.app.InitChain(abci.RequestInitChain{
				AppStateBytes: stateBytes,
			})
			return nil
		})
	}
}

func (suite *GenesisTestSuite) TestImportExport() {
	N := 5000

	var escrows []types.Escrow
	for i := 0; i < N; i++ {
		escrow, obj := suite.gen.NewRandomTestEscrow()
		escrow.IsAuction = rand.Bool()

		err := suite.crudStore.Create(obj)
		if err != nil {
			panic(fmt.Errorf("error while saving the escrow's object : %v", err))
		}

		id, err := suite.keeper.CreateEscrow(suite.ctx, obj.Owner, escrow.Price, obj, escrow.Deadline, escrow.IsAuction)
		if err != nil {
			panic(fmt.Errorf("error while creating an escrow : %v", err))
		}

		escrow, _ = suite.keeper.GetEscrow(suite.ctx, id)
		escrows = append(escrows, escrow)
	}

	lastBlockTime := suite.keeper.GetLastBlockTime(suite.ctx)
	nextId := suite.keeper.GetNextIDForExport(suite.ctx)
	genesis := escrow.ExportGenesis(suite.ctx, suite.keeper)

	// Wipe out escrows from db
	for _, e := range escrows {
		test.DeleteEscrow(suite.ctx, suite.storeKey, e.Id)
	}
	// Wipe out keeper params
	suite.keeper.ImportNextID(suite.ctx, 0)
	suite.keeper.SetLastBlockTime(suite.ctx, 0)

	escrow.InitGenesis(suite.ctx, suite.keeper, *genesis)

	suite.Equal(lastBlockTime, suite.keeper.GetLastBlockTime(suite.ctx))
	suite.Equal(nextId, suite.keeper.GetNextIDForExport(suite.ctx))
	for _, expected := range escrows {
		actual, found := suite.keeper.GetEscrow(suite.ctx, expected.Id)
		if !found {
			suite.Fail(fmt.Sprintf("Expected escrow %v not found", expected.Id))
		}

		suite.EqualValues(expected, actual, "Expected escrow different from actual one")
	}

	msg, broken := keeper.StateInvariant(suite.keeper)(suite.ctx)
	if broken {
		suite.Fail(fmt.Sprintf("Invriant broken after export with message : %v", msg))
	}
}
