package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/iov-one/starnamed/x/escrow/types"
)

type msgServer struct {
	Keeper
}

var _ types.MsgServer = msgServer{}

// NewMsgServerImpl returns an implementation of the escrow MsgServer interface
// for the provided Keeper.
func NewMsgServerImpl(keeper Keeper) types.MsgServer {
	return &msgServer{Keeper: keeper}
}

func (m msgServer) CreateEscrow(ctx context.Context, msg *types.MsgCreateEscrow) (*types.MsgCreateEscrowResponse, error) {
	sdkCtx := sdk.UnwrapSDKContext(ctx)

	// Extract and check seller address
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid seller address : %v", msg.Seller)
	}

	// Check that we are not using blocked (e.g module) accounts
	if m.isBlockedAddr(msg.Seller) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Seller)
	}

	obj := msg.Object.GetCachedValue().(types.TransferableObject)
	// Create the escrow
	id, err := m.Keeper.CreateEscrow(sdkCtx, seller, msg.Price, obj, msg.Deadline)
	if err != nil {
		return nil, err
	}

	// Collect fees
	if err := m.Keeper.CollectFees(sdkCtx, msg); err != nil {
		return nil, err
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventCreatedEscrow{
		Id:               id,
		Seller:           msg.Seller,
		FeePayer:         msg.FeePayer,
		BrokerAddress:    m.Keeper.GetBrokerAddress(sdkCtx),
		BrokerCommission: m.Keeper.GetBrokerCommission(sdkCtx),
		Price:            msg.Price,
		Object:           msg.Object,
		Deadline:         msg.Deadline,
		Fees:             m.Keeper.ComputeFees(sdkCtx, msg),
	}); err != nil {
		return nil, err
	}

	return &types.MsgCreateEscrowResponse{Id: id}, nil
}

func (m msgServer) UpdateEscrow(ctx context.Context, msg *types.MsgUpdateEscrow) (*types.MsgUpdateEscrowResponse, error) {

	// Extract and check updater and seller addresses
	updater, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid updater address : %v", msg.Updater)
	}
	// Check that we are not using blocked (e.g module) accounts
	if m.isBlockedAddr(msg.Updater) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Updater)
	}

	// The seller address is optional
	var seller sdk.AccAddress
	if len(msg.Seller) != 0 {
		seller, err = sdk.AccAddressFromBech32(msg.Seller)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "Invalid seller address : %v", msg.Seller)
		}
		// Check we are not using blocked addresses
		if m.isBlockedAddr(msg.Seller) {
			return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Seller)
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.UpdateEscrow(sdkCtx, msg.Id, updater, seller, msg.Price, msg.Deadline)
	if err != nil {
		return nil, err
	}

	// Collect fees
	if err := m.Keeper.CollectFees(sdkCtx, msg); err != nil {
		return nil, err
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventUpdatedEscrow{
		Id:          msg.Id,
		Updater:     msg.Updater,
		FeePayer:    msg.FeePayer,
		NewPrice:    msg.Price,
		NewSeller:   msg.Seller,
		NewDeadline: msg.Deadline,
		Fees:        m.Keeper.ComputeFees(sdkCtx, msg),
	}); err != nil {
		return nil, err
	}

	return &types.MsgUpdateEscrowResponse{}, nil
}

func (m msgServer) TransferToEscrow(ctx context.Context, msg *types.MsgTransferToEscrow) (*types.MsgTransferToEscrowResponse, error) {

	// Check and extract sender address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid sender address : %v", msg.Sender)
	}
	// Check that we are not using blocked (e.g module) accounts
	if m.isBlockedAddr(msg.Sender) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Sender)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.TransferToEscrow(sdkCtx, sender, msg.Id, msg.Amount)
	if err != nil {
		return nil, err
	}

	// Collect fees
	if err := m.Keeper.CollectFees(sdkCtx, msg); err != nil {
		return nil, err
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventCompletedEscrow{
		Id:       msg.Id,
		FeePayer: msg.FeePayer,
		Buyer:    msg.Sender,
		Fees:     m.Keeper.ComputeFees(sdkCtx, msg),
	}); err != nil {
		return nil, err
	}

	return &types.MsgTransferToEscrowResponse{}, nil
}

func (m msgServer) RefundEscrow(ctx context.Context, msg *types.MsgRefundEscrow) (*types.MsgRefundEscrowResponse, error) {

	// Check and extract the seller (who sent this message) address
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid sender address : %v", msg.Sender)
	}
	// Check that we are not using blocked (e.g module) accounts
	if m.isBlockedAddr(msg.Sender) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Sender)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.RefundEscrow(sdkCtx, sender, msg.Id)
	if err != nil {
		return nil, err
	}

	// Collect fees
	if err := m.Keeper.CollectFees(sdkCtx, msg); err != nil {
		return nil, err
	}

	// Emit event
	if err := sdkCtx.EventManager().EmitTypedEvent(&types.EventRefundedEscrow{
		Id:       msg.Id,
		FeePayer: msg.FeePayer,
		Sender:   msg.Sender,
		Fees:     m.Keeper.ComputeFees(sdkCtx, msg),
	}); err != nil {
		return nil, err
	}

	return &types.MsgRefundEscrowResponse{}, nil
}
