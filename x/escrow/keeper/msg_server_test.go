package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/app"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type MsgServerTestSuite struct {
	BaseKeeperSuite
}

func (s *MsgServerTestSuite) SetupTest() {
	s.Setup()
}

func (s *MsgServerTestSuite) TestAll() {

	t := s.T()
	validAddress := s.generator.NewAccAddress()

	type testCase struct {
		name   string
		seller string
		obj    proto.Message
	}

	commonTestCases := []testCase{
		{
			name:   "normal scenario",
			seller: validAddress.String(),
			obj:    s.generator.NewTestObject(validAddress),
		},
		{
			name:   "invalid seller : invalid bech32",
			seller: validAddress.String() + "5",
			obj:    s.generator.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid seller : invalid prefix",
			seller: strings.ReplaceAll(validAddress.String(), app.Bech32Prefix, "cosmos"),
			obj:    s.generator.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid seller: module address",
			seller: authtypes.NewModuleAddress(types.ModuleName).String(),
			obj:    s.generator.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid object: not a transferable object",
			seller: authtypes.NewModuleAddress(types.ModuleName).String(),
			obj:    new(types.TestObject),
		},
	}

	createTestCase := []testCase{
		{
			name:   "invalid object: not a transferable object",
			seller: authtypes.NewModuleAddress(types.ModuleName).String(),
			obj:    new(types.TestObject),
		},
	}

	createAndCheck := func(test testCase) string {
		shouldFail := strings.Contains(test.name, "invalid")
		resp, err := s.msgServer.CreateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgCreateEscrow{
			Seller:   test.seller,
			Object:   MustPackToAny(test.obj),
			Price:    sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(50))),
			Deadline: s.generator.NowAfter(100),
		})
		CheckError(t, test.name, shouldFail, err)
		return resp.Id
	}

	for _, test := range createTestCase {
		createAndCheck(test)
	}

	for _, test := range commonTestCases {
		shouldFail := strings.Contains(test.name, "invalid")

		// check creation
		id := createAndCheck(test)

		// check updating
		_, err := s.msgServer.UpdateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgUpdateEscrow{
			Id:       id,
			Updater:  test.seller,
			Seller:   test.seller,
			Price:    sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(80))),
			Deadline: s.generator.NowAfter(100),
		})
		CheckError(t, test.name, shouldFail, err)

		// check transfering
		_, err = s.msgServer.TransferToEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgTransferToEscrow{
			Id:     id,
			Sender: s.generator.NewAccAddress().String(),
			Amount: sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(80))),
		})
		CheckError(t, test.name, shouldFail, err)

		// check refunding (after creating a new one)
		id = createAndCheck(test)
		_, err = s.msgServer.RefundEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgRefundEscrow{
			Id:     id,
			Sender: s.generator.NewAccAddress().String(),
		})
		CheckError(t, test.name, shouldFail, err)
	}
}

func TestMsgServer(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}
