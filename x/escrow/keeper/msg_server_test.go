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
			"Invalid fee payed for test "+tc.name+"/"+operation)

	}

	checkNoFees := func(tc testCase, operation string, oldBalance sdk.Coins) {
		newBalance := getFeePayerBalance(tc)
		s.Assert().Equal(
			int64(0),
			oldBalance.Sub(newBalance).AmountOf(test.Denom).Int64(), // Only one denom for the fee
			"Invalid fee payed for test "+tc.name+"/"+operation)

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
			name:     "invalid feePayer : invalid bech32",
			feePayer: validAddress.String() + "5",
			seller:   validAddress.String(),
			obj:      s.generator.NewNotPossessedTestObject(),
		},
		{
			name:     "invalid feePayer : invalid prefix",
			seller:   validAddress.String(),
			feePayer: strings.ReplaceAll(validAddress.String(), app.Bech32Prefix, "cosmos"),
			obj:      s.generator.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid seller: module address",
			seller: authtypes.NewModuleAddress(types.ModuleName).String(),
			obj:    s.generator.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid object: not a transferable object",
			seller: validAddress.String(),
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

		// check transfering
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

func (s *MsgServerTestSuite) TestCompleteAuctionFees() {

	validAddress := s.generator.NewAccAddress()
	s.balances[validAddress.String()] = sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100)))

	type testCase struct {
		name     string
		sender   string
		feePayer string
	}

	getFeePayerBalance := func(tc testCase) sdk.Coins {
		payer := tc.feePayer
		if len(payer) == 0 {
			payer = tc.sender
		}
		return s.balances[payer]
	}
	checkFees := func(tc testCase, oldBalance sdk.Coins) {
		expectedDelta := s.keeper.ComputeFees(s.ctx, &types.MsgCompleteAuction{})

		newBalance := getFeePayerBalance(tc)
		s.Assert().Equal(
			expectedDelta[0].Amount.Int64(),                                     // Only one denom for the fee
			oldBalance.Sub(newBalance).AmountOf(expectedDelta[0].Denom).Int64(), // Only one denom for the fee
			"Invalid fee payed for test "+tc.name)

	}

	checkNoFees := func(tc testCase, oldBalance sdk.Coins) {
		newBalance := getFeePayerBalance(tc)
		s.Assert().Equal(
			int64(0),
			oldBalance.Sub(newBalance).AmountOf(test.Denom).Int64(), // Only one denom for the fee
			"Invalid fee payed for test "+tc.name)

	}

	commonTestCases := []testCase{
		{
			name:   "normal scenario",
			sender: validAddress.String(),
		},
		{
			name:     "normal scenario with fee payer",
			sender:   validAddress.String(),
			feePayer: s.feePayer.String(),
		},
	}

	createAuction := func(t testCase) string {
		obj := s.generator.NewNotPossessedTestObject()
		err := s.store.Create(obj)
		if err != nil {
			panic(err)
		}

		passedDeadline := s.generator.NowAfter(0) - 2
		s.keeper.SetLastBlockTime(s.ctx, passedDeadline - 2)

		resp, err := s.msgServer.CreateEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgCreateEscrow{
			Seller:    s.sender.String(),
			Object:    test.MustPackToAny(obj),
			Price:     sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50))),
			IsAuction: true,
			Deadline:  passedDeadline,
		})

		if err != nil {
			panic(err)
		}

		_, err = s.msgServer.TransferToEscrow(sdk.WrapSDKContext(s.ctx), &types.MsgTransferToEscrow{
			Id:       resp.Id,
			Sender:   s.buyer.String(),
			Amount:    sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(60))),
		})

		s.keeper.SetLastBlockTime(s.ctx, s.generator.NowAfter(0))
		s.keeper.MarkExpiredEscrows(s.ctx, s.generator.NowAfter(0))

		if resp == nil {
			return ""
		} else {
			return resp.Id
		}
	}

	for _, t := range commonTestCases {
		shouldFail := strings.Contains(t.name, "invalid")

		// check completing auction
		id := createAuction(t)
		oldFeePayerBalance := getFeePayerBalance(t)
		_, err := s.msgServer.CompleteAuction(sdk.WrapSDKContext(s.ctx), &types.MsgCompleteAuction{
			Id:       id,
			Sender:   t.sender,
			FeePayer: t.feePayer,
		})
		test.CheckError(s.T(), t.name, shouldFail, err)
		if !shouldFail {
			checkFees(t, oldFeePayerBalance)
		} else {
			checkNoFees(t, oldFeePayerBalance)
		}
	}
}

func TestMsgServer(t *testing.T) {
	suite.Run(t, new(MsgServerTestSuite))
}
