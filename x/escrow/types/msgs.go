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

func validateFeePayer(feePayer string) error {
	if len(feePayer) == 0 {
		return nil
	}
	_, err := sdk.AccAddressFromBech32(feePayer)
	if err != nil {
		return sdkerrors.Wrap(err, "invalid fee payer address")
	}
	return err
}

func getSigners(sender, feePayer string) []sdk.AccAddress {
	senderAddr, err := sdk.AccAddressFromBech32(sender)
	if err != nil {
		panic(err)
	}
	if len(feePayer) == 0 {
		return []sdk.AccAddress{senderAddr}
	}
	feePayerAddr, err := sdk.AccAddressFromBech32(feePayer)
	if err != nil {
		panic(err)
	}
	return []sdk.AccAddress{senderAddr, feePayerAddr}
}

func getFeePayer(sender, feePayer string) sdk.AccAddress {
	var err error
	var feePayerAddr sdk.AccAddress
	if len(feePayer) == 0 {
		feePayerAddr, err = sdk.AccAddressFromBech32(sender)
	} else {
		feePayerAddr, err = sdk.AccAddressFromBech32(feePayer)
	}
	if err != nil {
		panic(err)
	}
	return feePayerAddr
}

// NewMsgCreateEscrow creates a new MsgCreateEscrow instance
func NewMsgCreateEscrow(
	seller string,
	feePayer string,
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
		FeePayer: feePayer,
		Object:   packedObj,
		Price:    price,
		Deadline: deadline,
	}
}

// UnpackInterfaces make sure the Anys included in MsgCreateEscrow are unpacked (e.g the object field)
func (msg *MsgCreateEscrow) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if msg.Object != nil {
		var obj TransferableObject
		return unpacker.UnpackAny(msg.Object, &obj)
	}

	return nil
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

	if err := validateFeePayer(msg.FeePayer); err != nil {
		return err
	}

	if err := ValidatePrice(msg.Price, ""); err != nil {
		return err
	}

	switch msg.Object.GetCachedValue().(type) {
	case TransferableObject:
		break
	default:
		return sdkerrors.Wrapf(
			ErrUnknownObject,
			"The object should be of type TransferableObject but is of type %T",
			msg.Object.GetCachedValue(),
		)
	}

	obj := msg.Object.GetCachedValue().(TransferableObject)
	if err := ValidateObjectDeadlineBasic(obj, msg.Deadline); err != nil {
		return err
	}

	return ValidateObject(obj, seller)
}

// GetSignBytes implements Msg
func (msg MsgCreateEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetFeePayer implements MsgWithFeePayer
func (msg MsgCreateEscrow) GetFeePayer() sdk.AccAddress {
	return getFeePayer(msg.Seller, msg.FeePayer)
}

// GetSigners implements Msg
func (msg MsgCreateEscrow) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Seller, msg.FeePayer)
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

	if err := validateFeePayer(msg.FeePayer); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}
	return nil
}

// GetSignBytes implements Msg
func (msg MsgRefundEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetFeePayer implements MsgWithFeePayer
func (msg MsgRefundEscrow) GetFeePayer() sdk.AccAddress {
	return getFeePayer(msg.Sender, msg.FeePayer)
}

// GetSigners implements Msg
func (msg MsgRefundEscrow) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Sender, msg.FeePayer)
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

	if err := validateFeePayer(msg.FeePayer); err != nil {
		return err
	}

	var hasUpdate = false
	if len(msg.Seller) != 0 {
		hasUpdate = true
		_, err := sdk.AccAddressFromBech32(msg.Seller)
		if err != nil {
			return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
		}
	}

	if msg.Price != nil {
		hasUpdate = true
		if err := ValidatePrice(msg.Price, ""); err != nil {
			return err
		}
	}

	if msg.Deadline != 0 {
		hasUpdate = true
	}

	if !hasUpdate {
		return ErrEmptyUpdate
	}

	return nil
}

// GetSignBytes implements Msg
func (msg MsgUpdateEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetFeePayer implements MsgWithFeePayer
func (msg MsgUpdateEscrow) GetFeePayer() sdk.AccAddress {
	return getFeePayer(msg.Updater, msg.FeePayer)
}

// GetSigners implements Msg
func (msg MsgUpdateEscrow) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Updater, msg.FeePayer)
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

	if err := validateFeePayer(msg.FeePayer); err != nil {
		return err
	}

	_, err := sdk.AccAddressFromBech32(msg.Sender)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid sender address (%s)", err)
	}

	return ValidatePrice(msg.Amount, "")
}

// GetSignBytes implements Msg
func (msg MsgTransferToEscrow) GetSignBytes() []byte {
	bz := ModuleCdc.MustMarshalJSON(&msg)
	return sdk.MustSortJSON(bz)
}

// GetFeePayer implements MsgWithFeePayer
func (msg MsgTransferToEscrow) GetFeePayer() sdk.AccAddress {
	return getFeePayer(msg.Sender, msg.FeePayer)
}

// GetSigners implements Msg
func (msg MsgTransferToEscrow) GetSigners() []sdk.AccAddress {
	return getSigners(msg.Sender, msg.FeePayer)
}
