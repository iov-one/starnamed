package keeper_test

import (
	"strings"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	"github.com/gogo/protobuf/proto"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/app"
	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type MsgServerTestSuite struct {
	BaseKeeperSuite
}

func (s *MsgServerTestSuite) SetupTest() {
	s.Setup()
}

func (s *MsgServerTestSuite) TestAll() {

	validAddress := s.generator.NewAccAddress()

	type testCase struct {
		name   string
		seller string
		obj    interface{}
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
			obj:    new(types.GenesisState),
		},
	}

	createTestCase := []testCase{
		commonTestCases[len(commonTestCases)-1],
	}

	createAndCheck := func(t testCase) string {
		copiedObj := t.obj.(proto.Message)
		shouldFail := strings.Contains(t.name, "invalid")
		switch obj := t.obj.(type) {
		case *types.TestObject:
			// Make a copy of the object, to avoid to alter the original one
			cpy := &types.TestObject{
				Id:    obj.Id,
				Owner: append([]byte{}, obj.Owner...),
			}
			// returns ErrAlreadyExists if the object already exists and no other possible error
			err := s.store.Create(cpy)
			if err != nil {
				if err := s.store.Delete(cpy.PrimaryKey()); err != nil {
					panic(err)
				}
				_ = s.store.Create(cpy)
			}
			copiedObj = cpy
		}
		resp, err := s.msgServer.CreateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgCreateEscrow{
			Seller:   t.seller,
			Object:   test.MustPackToAny(copiedObj),
			Price:    sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(50))),
			Deadline: s.generator.NowAfter(100),
		})
		test.CheckError(s.T(), t.name+" creation", shouldFail, err)
		if resp == nil {
			return ""
		} else {
			return resp.Id
		}
	}

	for _, t := range createTestCase {
		createAndCheck(t)
	}

	for _, t := range commonTestCases {
		shouldFail := strings.Contains(t.name, "invalid")

		// check creation
		id := createAndCheck(t)

		// check updating
		_, err := s.msgServer.UpdateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgUpdateEscrow{
			Id:       id,
			Updater:  t.seller,
			Seller:   t.seller,
			Price:    sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(20))),
			Deadline: s.generator.NowAfter(100),
		})
		test.CheckError(s.T(), t.name+" update", shouldFail, err)

		// check transfering
		_, err = s.msgServer.TransferToEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgTransferToEscrow{
			Id:     id,
			Sender: s.generator.NewAccAddress().String(),
			Amount: sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(20))),
		})
		test.CheckError(s.T(), t.name+" transfer", shouldFail, err)

		// check refunding (after creating a new one)
		id = createAndCheck(t)
		_, err = s.msgServer.RefundEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgRefundEscrow{
			Id:     id,
			Sender: t.seller,
		})
		test.CheckError(s.T(), t.name+" refund", shouldFail, err)
	}
}

func TestMsgServer(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}
