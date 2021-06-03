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

	// Extract and check seller and buyer addresses
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid seller address : %v", msg.Seller)
	}

	buyer, err := sdk.AccAddressFromBech32(msg.Buyer)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid buyer address : %v", msg.Buyer)
	}

	// Check that we are not using blocked (e.g module) accounts
	if m.isBlockedAddr(msg.Seller) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Seller)
	}
	if m.isBlockedAddr(msg.Buyer) {
		return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Buyer)
	}

	obj := msg.Object.GetCachedValue().(types.TransferableObject)
	// Create the escrow
	id, err := m.Keeper.CreateEscrow(sdkCtx, seller, buyer, msg.Price, obj, msg.Deadline)
	return &types.MsgCreateEscrowResponse{Id: id}, nil
}

func (m msgServer) UpdateEscrow(ctx context.Context, msg *types.MsgUpdateEscrow) (*types.MsgUpdateEscrowResponse, error) {

	// Extract and check updater, seller and buyer address
	updater, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid updater address : %v", msg.Updater)
	}

	// The sender and buyer are optional
	var seller, buyer sdk.AccAddress
	if len(msg.Seller) != 0 {
		seller, err = sdk.AccAddressFromBech32(msg.Seller)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "Invalid seller address : %v", msg.Seller)
		}
		if m.isBlockedAddr(msg.Seller) {
			return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Seller)
		}
	}
	if len(msg.Buyer) != 0 {
		buyer, err = sdk.AccAddressFromBech32(msg.Buyer)
		if err != nil {
			return nil, sdkerrors.Wrapf(err, "Invalid buyer address : %v", msg.Buyer)
		}
		if m.isBlockedAddr(msg.Buyer) {
			return nil, sdkerrors.Wrap(types.ErrInvalidAccount, msg.Buyer)
		}
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.UpdateEscrow(sdkCtx, msg.Id, updater, seller, buyer, msg.Price, msg.Deadline)
	if err != nil {
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

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.TransferToEscrow(sdkCtx, sender, msg.Id, msg.Amount)
	if err != nil {
		return nil, err
	}

	return &types.MsgTransferToEscrowResponse{}, nil
}

func (m msgServer) RefundEscrow(ctx context.Context, msg *types.MsgRefundEscrow) (*types.MsgRefundEscrowResponse, error) {

	// Check and extract the seller (who sent this message) address
	sender, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "Invalid seller address : %v", msg.Seller)
	}

	sdkCtx := sdk.UnwrapSDKContext(ctx)
	err = m.Keeper.RefundEscrow(sdkCtx, sender, msg.Id)
	if err != nil {
		return nil, err
	}

	return &types.MsgRefundEscrowResponse{}, nil
}
