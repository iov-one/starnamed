package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// CompleteAuction complete an auction whose deadline has passed.
// The auction must have at least one bid to be completed, otherwise it has to be refunded with RefundEscrow instead
// It sends the object to the last bidder and send the locked tokens to the seller
func (k Keeper) CompleteAuction(ctx sdk.Context, id string) error {
	k.checkThatModuleIsEnabled(ctx)

	auction, found := k.GetEscrow(ctx, id)
	if !found {
		return sdkerrors.Wrap(types.ErrEscrowNotFound, id)
	}

	if !auction.IsAuction {
		return sdkerrors.Wrap(types.ErrNotAnAuction, id)
	}

	if auction.State != types.EscrowState_Expired {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the auction has not finished yet")
	}

	if len(auction.LastBidder) == 0 {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "the auction has no bidder and cannot be completed, call refund instead")
	}

	seller, err := sdk.AccAddressFromBech32(auction.Seller)
	if err != nil {
		// This should not happen as the seller address is already validated
		panic(err)
	}

   // We complete the auction by doing the swap and removing the escrow
	broker, err := sdk.AccAddressFromBech32(auction.BrokerAddress)
	if err != nil {
		// This should not happen as the seller address is already validated
		panic(err)
	}

	lastBidder, err := sdk.AccAddressFromBech32(auction.LastBidder)
	if err != nil {
		// This should not happen as the seller address is already validated
		panic(err)
	}

	// Do the swap between the object and the coins
	err = k.doSwap(ctx, auction, lastBidder, seller, broker)
	if err != nil {
		panic(err)
	}

	// We mark the auction as completed and we delete it
	auction.State = types.EscrowState_Completed
	k.deleteEscrow(ctx, auction)

	return nil
}
