package types

import (
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
)

// NewEscrow constructs a new escrow instance
func NewEscrow(
	id tmbytes.HexBytes,
	seller sdk.AccAddress,
	buyer sdk.AccAddress,
	price sdk.Coins,
	object TransferableObject,
	//TODO: do we need a state ?
	state EscrowState,
	deadline uint64,
) Escrow {
	objectAny, err := codectypes.NewAnyWithValue(object)
	if err != nil {
		panic(err)
	}
	return Escrow{
		Id:       id.String(),
		Seller:   seller.String(),
		Buyer:    buyer.String(),
		Object:   objectAny,
		Price:    price,
		State:    state,
		Deadline: deadline,
	}
}

// Validate validates the escrow
func (e Escrow) Validate() error {
	if err := ValidateID(e.Id); err != nil {
		return err
	}
	// Validate seller and buyer accounts
	seller, err := sdk.AccAddressFromBech32(e.Seller)
	if err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(e.Buyer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}
	// Validate object valid and possessed by the seller
	if err := ValidateObject(e.GetObject(), seller); err != nil {
		return err
	}
	// Validate price
	if err := ValidatePrice(e.Price); err != nil {
		return err
	}
	// TODO: Validate seller and buyer account exist ??
	// Validate state
	if err := ValidateState(e.State); err != nil {
		return err
	}
	// Validate deadline
	if err := ValidateDeadline(e.Deadline); err != nil {
		return err
	}
	return nil
}

func (e *Escrow) GetObject() TransferableObject {
	return e.Object.GetCachedValue().(TransferableObject)
}
