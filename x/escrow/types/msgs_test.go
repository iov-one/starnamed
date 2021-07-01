package types_test

import (
	"encoding/hex"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type MsgTestSuite struct {
	suite.Suite
	msgCreate   types.MsgCreateEscrow
	msgRefund   types.MsgRefundEscrow
	msgTransfer types.MsgTransferToEscrow
	msgUpdate   types.MsgUpdateEscrow
	sender      sdk.AccAddress
	gen         *test.EscrowGenerator
}

func (suite *MsgTestSuite) SetupTest() {
	suite.gen = test.NewEscrowGenerator(200)

	suite.sender = suite.gen.NewAccAddress()
	validId := hex.EncodeToString(sdk.Uint64ToBigEndian(42))
	validPrice := sdk.NewCoins(sdk.NewCoin("denom", sdk.NewInt(50)))
	validObject := suite.gen.NewTestObject(suite.sender)

	suite.msgCreate = types.NewMsgCreateEscrow(suite.sender.String(), "", validObject, validPrice, suite.gen.NowAfter(0))
	suite.msgRefund = types.MsgRefundEscrow{
		Id:     validId,
		Sender: suite.sender.String(),
	}
	suite.msgUpdate = types.MsgUpdateEscrow{
		Id:       validId,
		Updater:  suite.sender.String(),
		Seller:   "",
		Price:    validPrice,
		Deadline: suite.gen.NowAfter(0),
	}
	suite.msgTransfer = types.MsgTransferToEscrow{
		Id:     validId,
		Sender: suite.sender.String(),
		Amount: validPrice,
	}
}

func (suite *MsgTestSuite) TestMsgValidate() {
	invalidBech32Addr := "star15d5e8f5544f5e5fe654"
	invalidPrefixAddr := "cosmos1cqfse93m6r7fr3vx07du5yfmsltca60gyadygf"
	negativePrice := sdk.Coins{sdk.Coin{Denom: "tiov", Amount: sdk.NewInt(-10)}}
	invalidIDHexa := "123456789abcdefg"
	invalidIDLength := "123456789abcdef0123"
	randomOwnerObj := test.MustPackToAny(suite.gen.NewTestObject(suite.gen.NewAccAddress()))
	invalidInterfaceObj := test.MustPackToAny(&suite.msgUpdate)

	completeMsgCreate := func(msg types.MsgCreateEscrow) *types.MsgCreateEscrow {
		if msg.Price == nil {
			msg.Price = suite.msgCreate.Price
		}
		if msg.Deadline == 0 {
			msg.Deadline = suite.msgCreate.Deadline
		}
		if len(msg.Seller) == 0 {
			msg.Seller = suite.msgCreate.Seller
		}
		if msg.Object == nil {
			msg.Object = suite.msgCreate.Object
		}
		return &msg
	}
	completeMsgUpdate := func(msg types.MsgUpdateEscrow) *types.MsgUpdateEscrow {
		if len(msg.Id) == 0 {
			msg.Id = suite.msgUpdate.Id
		}
		if len(msg.Updater) == 0 {
			msg.Updater = suite.msgUpdate.Updater
		}
		return &msg
	}
	completeMsgTransfer := func(msg types.MsgTransferToEscrow) *types.MsgTransferToEscrow {
		if len(msg.Id) == 0 {
			msg.Id = suite.msgTransfer.Id
		}
		if len(msg.Sender) == 0 {
			msg.Sender = suite.msgTransfer.Sender
		}
		if msg.Amount == nil {
			msg.Amount = suite.msgTransfer.Amount
		}
		return &msg
	}
	completeMsgRefund := func(msg types.MsgRefundEscrow) *types.MsgRefundEscrow {
		if len(msg.Id) == 0 {
			msg.Id = suite.msgRefund.Id
		}
		if len(msg.Sender) == 0 {
			msg.Sender = suite.msgRefund.Sender
		}
		return &msg
	}

	testCases := []struct {
		name string
		msg  sdk.Msg
	}{
		{
			name: "create: valid",
			msg:  &suite.msgCreate,
		},
		{
			name: "create: valid with fee payer",
			msg:  completeMsgCreate(types.MsgCreateEscrow{FeePayer: suite.sender.String()}),
		},
		{
			name: "create: invalid seller address: invalid bech32",
			msg:  completeMsgCreate(types.MsgCreateEscrow{Seller: invalidBech32Addr}),
		},
		{
			name: "create: invalid seller address: invalid prefix",
			msg:  completeMsgCreate(types.MsgCreateEscrow{Seller: invalidPrefixAddr}),
		},
		{
			name: "create: invalid fee payer: invalid bech32",
			msg:  completeMsgCreate(types.MsgCreateEscrow{FeePayer: invalidBech32Addr}),
		},
		{
			name: "create: invalid fee payer: invalid prefix",
			msg:  completeMsgCreate(types.MsgCreateEscrow{FeePayer: invalidPrefixAddr}),
		},
		{
			name: "create: invalid price: negative",
			msg:  completeMsgCreate(types.MsgCreateEscrow{Price: negativePrice}),
		},
		{
			name: "create: invalid object: does not belong to seller",
			msg:  completeMsgCreate(types.MsgCreateEscrow{Object: randomOwnerObj}),
		},
		{
			name: "create: invalid object: not a TransferableObject",
			msg:  completeMsgCreate(types.MsgCreateEscrow{Object: invalidInterfaceObj}),
		},
		{
			name: "update: valid",
			msg:  &suite.msgUpdate,
		},
		{
			name: "update: valid with fee payer",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{FeePayer: suite.sender.String(), Seller: suite.sender.String()}),
		},
		{
			name: "update: invalid empty update",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{}),
		},
		{
			name: "update: invalid seller address: invalid bech32",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Seller: invalidBech32Addr}),
		},
		{
			name: "update: invalid seller address: invalid prefix",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Seller: invalidPrefixAddr}),
		},
		{
			name: "update: invalid updater address: invalid bech32",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Updater: invalidBech32Addr}),
		},
		{
			name: "update: invalid updater address: invalid prefix",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Updater: invalidPrefixAddr}),
		},
		{
			name: "update: invalid fee payer: invalid bech32",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{FeePayer: invalidBech32Addr}),
		},
		{
			name: "update: invalid fee payer: invalid prefix",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{FeePayer: invalidPrefixAddr}),
		},
		{
			name: "update: invalid price: negative",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Price: negativePrice}),
		},
		{
			name: "update: invalid escrow ID: not hexadecimal",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Id: invalidIDHexa}),
		},
		{
			name: "update: invalid escrow ID: invalid length",
			msg:  completeMsgUpdate(types.MsgUpdateEscrow{Id: invalidIDLength}),
		},
		{
			name: "transfer: valid",
			msg:  &suite.msgTransfer,
		},
		{
			name: "transfer: valid with fee payer",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{FeePayer: suite.sender.String()}),
		},
		{
			name: "transfer: invalid sender: invalid bech32",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{Sender: invalidBech32Addr}),
		},
		{
			name: "transfer: invalid sender: invalid prefix",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{Sender: invalidPrefixAddr}),
		},
		{
			name: "transfer: invalid fee payer: invalid bech32",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{FeePayer: invalidBech32Addr}),
		},
		{
			name: "transfer: invalid fee payer: invalid prefix",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{FeePayer: invalidPrefixAddr}),
		},
		{
			name: "transfer: invalid amount: negative",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{Amount: negativePrice}),
		},
		{
			name: "transfer: invalid escrow ID: not hexadecimal",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{Id: invalidIDHexa}),
		},
		{
			name: "transfer: invalid escrow ID: invalid length",
			msg:  completeMsgTransfer(types.MsgTransferToEscrow{Id: invalidIDLength}),
		},
		{
			name: "refund: valid",
			msg:  &suite.msgRefund,
		},
		{
			name: "refund: valid with fee payer",
			msg:  completeMsgRefund(types.MsgRefundEscrow{FeePayer: suite.sender.String()}),
		},
		{
			name: "refund: invalid seller: not bech32",
			msg:  completeMsgRefund(types.MsgRefundEscrow{Sender: invalidBech32Addr}),
		},
		{
			name: "refund: invalid seller: invalid prefix",
			msg:  completeMsgRefund(types.MsgRefundEscrow{Sender: invalidPrefixAddr}),
		},
		{
			name: "refund: invalid fee payer: invalid bech32",
			msg:  completeMsgRefund(types.MsgRefundEscrow{FeePayer: invalidBech32Addr}),
		},
		{
			name: "refund: invalid fee payer: invalid prefix",
			msg:  completeMsgRefund(types.MsgRefundEscrow{FeePayer: invalidPrefixAddr}),
		},
		{
			name: "refund: invalid escrow ID: not hexadecimal",
			msg:  completeMsgRefund(types.MsgRefundEscrow{Id: invalidIDHexa}),
		},
		{
			name: "refund: invalid escrow ID: invalid length",
			msg:  completeMsgRefund(types.MsgRefundEscrow{Id: invalidIDLength}),
		},
	}

	for _, tc := range testCases {
		test.EvaluateTest(suite.T(), tc.name, func(*testing.T) error { return tc.msg.ValidateBasic() })
	}
}

func (suite *MsgTestSuite) TestMsgGetSignersAndGetFeePayer() {
	senderInArray := []sdk.AccAddress{suite.sender}
	t := suite.T()
	require.Equal(t, senderInArray, suite.msgRefund.GetSigners(), "Invalid refund message signers")
	require.Equal(t, senderInArray, suite.msgTransfer.GetSigners(), "Invalid transfer message signers")
	require.Equal(t, senderInArray, suite.msgCreate.GetSigners(), "Invalid create message signers")
	require.Equal(t, senderInArray, suite.msgUpdate.GetSigners(), "Invalid update message signers")
	require.Equal(t, suite.sender, suite.msgRefund.GetFeePayer(), "Invalid refund message fee payer")
	require.Equal(t, suite.sender, suite.msgTransfer.GetFeePayer(), "Invalid transfer message fee payer")
	require.Equal(t, suite.sender, suite.msgCreate.GetFeePayer(), "Invalid create message fee payer")
	require.Equal(t, suite.sender, suite.msgUpdate.GetFeePayer(), "Invalid update message fee payer")

	feePayer := suite.gen.NewAccAddress()
	senderWithFeePayer := []sdk.AccAddress{suite.sender, feePayer}

	msgRefundWithFeePayer := suite.msgRefund
	msgRefundWithFeePayer.FeePayer = feePayer.String()
	msgCreateWithFeePayer := suite.msgCreate
	msgCreateWithFeePayer.FeePayer = feePayer.String()
	msgUpdateWithFeePayer := suite.msgUpdate
	msgUpdateWithFeePayer.FeePayer = feePayer.String()
	msgTransferWithFeePayer := suite.msgTransfer
	msgTransferWithFeePayer.FeePayer = feePayer.String()

	require.Equal(t, senderWithFeePayer, msgRefundWithFeePayer.GetSigners(), "Invalid refund message signers with fee payer")
	require.Equal(t, senderWithFeePayer, msgTransferWithFeePayer.GetSigners(), "Invalid transfer message signers with fee payer")
	require.Equal(t, senderWithFeePayer, msgCreateWithFeePayer.GetSigners(), "Invalid create message signers with fee payer")
	require.Equal(t, senderWithFeePayer, msgUpdateWithFeePayer.GetSigners(), "Invalid update message signers with fee payer")
	require.Equal(t, feePayer, msgRefundWithFeePayer.GetFeePayer(), "Invalid refund message fee payer with fee payer")
	require.Equal(t, feePayer, msgTransferWithFeePayer.GetFeePayer(), "Invalid transfer message fee payer with fee payer")
	require.Equal(t, feePayer, msgCreateWithFeePayer.GetFeePayer(), "Invalid create message fee payer with fee payer")
	require.Equal(t, feePayer, msgUpdateWithFeePayer.GetFeePayer(), "Invalid update message fee payer with fee payer")

}

func TestMsgSuite(t *testing.T) {
	suite.Run(t, new(MsgTestSuite))
}
