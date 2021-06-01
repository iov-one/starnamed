package types

import (
	"github.com/iov-one/starnamed/x/escrow/types"
	"github.com/tendermint/tendermint/crypto/tmhash"
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
	object types.TransferableObject,
	state EscrowState,
	deadline uint64,
) Escrow {
	return Escrow{
		Id:       id.String(),
		Seller:   seller.String(),
		Buyer:    buyer.String(),
		Object:   object,
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
	if _, err := sdk.AccAddressFromBech32(e.Seller); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid seller address (%s)", err)
	}
	if _, err := sdk.AccAddressFromBech32(e.Buyer); err != nil {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidAddress, "invalid buyer address (%s)", err)
	}
	// Validate accounts not module accounts
	// Validate starname
	// Validate starname possesed by seller
	// Validate seller and buyer account exist ??
	// Validate state
	return nil
}

//TODO: find a way to compute unique ID
func GetID(
	sender sdk.AccAddress,
	to sdk.AccAddress,
	amount sdk.Coins,
	hashLock tmbytes.HexBytes,
) tmbytes.HexBytes {
	return tmhash.Sum(
		append(append(append(hashLock, sender...), to...), []byte(amount.Sort().String())...),
	)
}
