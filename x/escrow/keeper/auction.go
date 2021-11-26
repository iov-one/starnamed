package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/iov-one/starnamed/x/escrow/types"
)

func (k Keeper) completeAuction(ctx sdk.Context, auction types.Escrow) error {

	//TODO: this check seems odd here, maybe an assert-like check with a panic will be more appropriate
	if !auction.IsAuction {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "the specified escrow should be an auction but isn't one")
	}
	if auction.State != types.EscrowState_Open {
		return types.ErrEscrowNotOpen
	}

	seller, err := sdk.AccAddressFromBech32(auction.Seller)
	if err != nil {
		// This should not happen as the seller address is already validated
		panic(err)
	}

	// If no one has made a bid, then just refund the escrow
	if len(auction.LastBidder) == 0 {
		err = k.refundEscrow(ctx, auction, seller)
		if err != nil {
			return err
		}
	} else { // Else, we complete the auction by doing the swap and removing the escrow

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
	}
	return nil
}
