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
	sender, buyer, feePayer sdk.AccAddress
}

func (s *MsgServerTestSuite) SetupTest() {
	s.generator = test.NewEscrowGenerator(uint64(s.ctx.BlockTime().Unix()))
	s.sender = s.generator.NewAccAddress()
	s.feePayer = s.generator.NewAccAddress()
	s.buyer = s.generator.NewAccAddress()
	s.Setup([]sdk.AccAddress{s.sender, s.feePayer, s.buyer}, true)
}

func (s *MsgServerTestSuite) TestAll() {

	validAddress := s.sender

	type testCase struct {
		name     string
		seller   string
		feePayer string
		obj      interface{}
	}

	msgs := map[string]sdk.Msg{
		"creation": &types.MsgCreateEscrow{},
		"transfer": &types.MsgTransferToEscrow{},
		"update":   &types.MsgUpdateEscrow{},
		"refund":   &types.MsgRefundEscrow{},
	}

	getFeePayerBalance := func(tc testCase) sdk.Coins {
		payer := tc.feePayer
		if len(payer) == 0 {
			payer = tc.seller
		}
		return s.balances[payer]
	}
	checkFees := func(tc testCase, operation string, oldBalance sdk.Coins) {
		expectedDelta := s.keeper.ComputeFees(s.ctx, msgs[operation])

		newBalance := getFeePayerBalance(tc)
		s.Assert().Equal(
			expectedDelta[0].Amount.Int64(),                                     // Only one denom for the fee
			oldBalance.Sub(newBalance).AmountOf(expectedDelta[0].Denom).Int64(), // Only one denom for the fee
			"Invalid fee paid for test "+tc.name+"/"+operation)

	}

	checkNoFees := func(tc testCase, operation string, oldBalance sdk.Coins) {
		newBalance := getFeePayerBalance(tc)
		s.Assert().Equal(
			int64(0),
			oldBalance.Sub(newBalance).AmountOf(test.Denom).Int64(), // Only one denom for the fee
			"Invalid fee paid for test "+tc.name+"/"+operation)

	}

	commonTestCases := []testCase{
		{
			name:   "normal scenario",
			seller: validAddress.String(),
			obj:    s.generator.NewTestObject(validAddress),
		},
		{
			name:     "normal scenario with fee payer",
			seller:   validAddress.String(),
			feePayer: s.feePayer.String(),
			obj:      s.generator.NewTestObject(validAddress),
		},
		{
			name:     "invalid seller : fee payer but no seller",
			feePayer: validAddress.String(),
			obj:      s.generator.NewNotPossessedTestObject(),
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
				Id:                  obj.Id,
				Owner:               append([]byte{}, obj.Owner...),
				NumAllowedTransfers: obj.NumAllowedTransfers,
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
		oldFeePayerBalance := getFeePayerBalance(t)
		resp, err := s.msgServer.CreateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgCreateEscrow{
			Seller:   t.seller,
			FeePayer: t.feePayer,
			Object:   test.MustPackToAny(copiedObj),
			Price:    sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50))),
			Deadline: s.generator.NowAfter(100),
		})
		test.CheckError(s.T(), t.name+" creation", shouldFail, err)

		if !shouldFail {
			checkFees(t, "creation", oldFeePayerBalance)
		} else {
			checkNoFees(t, "creation", oldFeePayerBalance)
		}

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
		oldFeePayerBalance := getFeePayerBalance(t)
		updatedPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(20)))
		_, err := s.msgServer.UpdateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgUpdateEscrow{
			Id:       id,
			Updater:  t.seller,
			FeePayer: t.feePayer,
			Seller:   t.seller,
			Price:    updatedPrice,
			Deadline: s.generator.NowAfter(100),
		})
		test.CheckError(s.T(), t.name+" update", shouldFail, err)
		if !shouldFail {
			checkFees(t, "update", oldFeePayerBalance)
		} else {
			checkNoFees(t, "update", oldFeePayerBalance)
		}

		// check transferring
		tCopy := t
		tCopy.seller = s.buyer.String()
		oldFeePayerBalance = getFeePayerBalance(tCopy)
		_, err = s.msgServer.TransferToEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgTransferToEscrow{
			Id:       id,
			Sender:   tCopy.seller,
			FeePayer: t.feePayer,
			Amount:   updatedPrice,
		})
		test.CheckError(s.T(), t.name+" transfer", shouldFail, err)
		if !shouldFail {
			if len(t.feePayer) == 0 {
				oldFeePayerBalance = oldFeePayerBalance.Sub(updatedPrice)
			}
			checkFees(tCopy, "transfer", oldFeePayerBalance)
		} else {
			checkNoFees(tCopy, "transfer", oldFeePayerBalance)
		}

		// check refunding (after creating a new one)
		id = createAndCheck(t)
		oldFeePayerBalance = getFeePayerBalance(t)
		_, err = s.msgServer.RefundEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgRefundEscrow{
			Id:       id,
			Sender:   t.seller,
			FeePayer: t.feePayer,
		})
		test.CheckError(s.T(), t.name+" refund", shouldFail, err)
		if !shouldFail {
			checkFees(t, "refund", oldFeePayerBalance)
		} else {
			checkNoFees(t, "refund", oldFeePayerBalance)
		}

	}
}

func TestMsgServer(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}
