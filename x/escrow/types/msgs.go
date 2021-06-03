package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

const (
	// TypeMsgCreateEscrow is the type for MsgCreateEscrow
	TypeMsgCreateEscrow = "create_escrow"
	// TypeMsgRefundEscrow is the type for MsgRefundEscrow
	TypeMsgRefundEscrow = "refund_escrow"
	// TypeMsgUpdateEscrow is the type for MsgUpdateEscrow
	TypeMsgUpdateEscrow = "update_escrow"
	// TypeMsgTransferToEscrow is the type for MsgTransferToEscrow
	TypeMsgTransferToEscrow = "transfer_to_escrow"
)

var (
	_ sdk.Msg = &MsgCreateEscrow{}
	_ sdk.Msg = &MsgRefundEscrow{}
	_ sdk.Msg = &MsgUpdateEscrow{}
	_ sdk.Msg = &MsgTransferToEscrow{}
)

// NewMsgCreateEscrow creates a new MsgCreateEscrow instance
func NewMsgCreateEscrow(
	seller string,
	buyer string,
	object TransferableObject,
	price sdk.Coins,
	deadline uint64,
) MsgCreateEscrow {
	packedObj, err := codectypes.NewAnyWithValue(object)
	if err != nil {
		panic(err)
	}
	return MsgCreateEscrow{
		Seller:   seller,
		Buyer:    buyer,
		Object:   packedObj,
		Price:    price,
		Deadline: deadline,
	}
}

// Route implements Msg
func (msg MsgCreateEscrow) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgCreateEscrow) Type() string { return TypeMsgCreateEscrow }

// ValidateBasic implements Msg
func (msg MsgCreateEscrow) ValidateBasic() error {
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}

	if _, err := sdk.AccAddressFromBech32(msg.Buyer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}

	if err := ValidatePrice(msg.Price); err != nil {
		return err
	}

	if err := ValidateObject(msg.Object.GetCachedValue().(TransferableObject), seller); err != nil {
		return err
	}

	return ValidateDeadline(msg.Deadline)
}

// GetSignBytes implements Msg
func (msg MsgCreateEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgCreateEscrow) GetSigners() []sdk.AccAddress {
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

// -----------------------------------------------------------------------------

// Route implements Msg
func (msg MsgRefundEscrow) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgRefundEscrow) Type() string { return TypeMsgRefundEscrow }

// ValidateBasic implements Msg
func (msg MsgRefundEscrow) ValidateBasic() error {

	if err := ValidateID(msg.Id); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	return nil
}

// GetSignBytes implements Msg
func (msg MsgRefundEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgRefundEscrow) GetSigners() []sdk.AccAddress {
	seller, err := sdk.AccAddressFromBech32(msg.Seller)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{seller}
}

// -----------------------------------------------------------------------------

// Route implements Msg
func (msg MsgUpdateEscrow) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgUpdateEscrow) Type() string { return TypeMsgUpdateEscrow }

// ValidateBasic implements Msg
func (msg MsgUpdateEscrow) ValidateBasic() error {

	if err := ValidateID(msg.Id); err != nil {
		return err
	}

	if _, err := sdk.AccAddressFromBech32(msg.Updater); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid updater address (%s)", err)
	}

	if len(msg.Seller) != 0 {
		_, err := sdk.AccAddressFromBech32(msg.Seller)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
		}
	}

	if len(msg.Buyer) != 0 {
		if _, err := sdk.AccAddressFromBech32(msg.Buyer); err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
		}
	}

	if msg.Price != nil {
		if err := ValidatePrice(msg.Price); err != nil {
			return err
		}
	}

	if msg.Deadline != 0 {
		return ValidateDeadline(msg.Deadline)
	}
	return nil
}

// GetSignBytes implements Msg
func (msg MsgUpdateEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgUpdateEscrow) GetSigners() []sdk.AccAddress {
	updater, err := sdk.AccAddressFromBech32(msg.Updater)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{updater}
}

// -----------------------------------------------------------------------------

// Route implements Msg
func (msg MsgTransferToEscrow) Route() string { return RouterKey }

// Type implements Msg
func (msg MsgTransferToEscrow) Type() string { return TypeMsgTransferToEscrow }

// ValidateBasic implements Msg
func (msg MsgTransferToEscrow) ValidateBasic() error {
	if err := ValidateID(msg.Id); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return ValidatePrice(msg.Amount)
}

// GetSignBytes implements Msg
func (msg MsgTransferToEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetSigners implements Msg
func (msg MsgTransferToEscrow) GetSigners() []sdk.AccAddress {
	sender, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{sender}
}
