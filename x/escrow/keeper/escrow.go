package keeper

import (
	"bytes"
	"encoding/hex"

	"github.com/pkg/errors"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/escrow/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

// CreateEscrow creates an escrow
func (k Keeper) CreateEscrow(
	ctx sdk.Context,
	seller sdk.AccAddress,
	buyer sdk.AccAddress,
	price sdk.Coins,
	object types.TransferableObject,
	deadline uint64,
) (
	tmbytes.HexBytes,
	error,
) {
	// TODO : replace this
	id := k.GetNextId()

	// check if the escrow already exists
	if k.HasEscrow(ctx, id) {
		return id, sdkerrors.Wrap(types.ErrIdNotUnique, id.String())
	}

	//TODO: transfer ownership of starname to escrow
	err := object.Transfer(seller, k.GetEscrowAccount(ctx).GetAddress())
	if err != nil {
		return nil, errors.Wrap(err, "Cannot transfer the object to the module account")
	}

	escrow := types.NewEscrow(
		id, buyer, seller, price, object, types.Open, deadline,
	)

	// set the escrow
	k.SetEscrow(ctx, escrow, id)

	return id, nil
}

// ClaimEscrow claims the specified escrow with the given secret
func (k Keeper) ClaimEscrow(
	ctx sdk.Context,
	id tmbytes.HexBytes,
	secret tmbytes.HexBytes,
) (
	string,
	bool,
	error,
) {
	// query the escrow
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return "", false, types.None, sdkerrors.Wrap(types.ErrUnknownEscrow, id.String())
	}

	// check if the escrow is open
	if escrow.State != types.Open {
		return "", false, types.None, sdkerrors.Wrap(types.ErrEscrowNotOpen, id.String())
	}

	hashLock, _ := hex.DecodeString(escrow.HashLock)

	// check if the secret matches with the hash lock
	if !bytes.Equal(types.GetHashLock(secret, escrow.Timestamp), hashLock) {
		return "", false, types.None, sdkerrors.Wrap(types.ErrInvalidSecret, secret.String())
	}

	to, err := sdk.AccAddressFromBech32(escrow.To)
	if err != nil {
		return "", false, types.None, err
	}

	if escrow.Transfer {
		if err := k.claimHTLT(ctx, escrow); err != nil {
			return "", false, types.None, err
		}
	} else {
		if err := k.claimEscrow(ctx, escrow.Amount, to); err != nil {
			return "", false, types.None, err
		}
	}

	// update the secret and state of the escrow
	escrow.Secret = secret.String()
	escrow.State = types.Completed
	escrow.ClosedBlock = uint64(ctx.BlockHeight())
	k.SetEscrow(ctx, escrow, id)

	// delete from the expiration queue
	k.DeleteEscrowFromExpiredQueue(ctx, escrow.ExpirationHeight, id)

	return escrow.HashLock, escrow.Transfer, escrow.Direction, nil
}

// RefundEscrow refunds the specified escrow
func (k Keeper) RefundEscrow(ctx sdk.Context, h types.Escrow, id tmbytes.HexBytes) error {
	sender, err := sdk.AccAddressFromBech32(h.Sender)
	if err != nil {
		return err
	}

	if h.Transfer {
		if err := k.refundHTLT(ctx, h.Direction, sender, h.Amount); err != nil {
			return err
		}
	} else {
		if err := k.refundEscrow(ctx, sender, h.Amount); err != nil {
			return err
		}
	}

	// update the state of the escrow
	h.State = types.Refunded
	h.ClosedBlock = uint64(ctx.BlockHeight())
	k.SetEscrow(ctx, h, id)

	return nil
}

func (k Keeper) getStore(ctx sdk.Context) sdk.KVStore {
	return ctx.KVStore(k.storeKey)
}

// HasEscrow checks if the given escrow exists
func (k Keeper) HasEscrow(ctx sdk.Context, id tmbytes.HexBytes) bool {
	return k.getStore(ctx).Has(types.GetEscrowKey(id))
}

// SetEscrow sets the given escrow
func (k Keeper) SetEscrow(ctx sdk.Context, escrow types.Escrow, id tmbytes.HexBytes) {
	bz := k.cdc.MustMarshalBinaryBare(&escrow)
	k.getStore(ctx).Set(types.GetEscrowKey(id), bz)
}

// GetEscrow retrieves the specified escrow
func (k Keeper) GetEscrow(ctx sdk.Context, id tmbytes.HexBytes) (escrow types.Escrow, found bool) {
	bz := k.getStore(ctx).Get(types.GetEscrowKey(id))
	if bz == nil {
		return escrow, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &escrow)
	return escrow, true
}

// IterateEscrows iterates through the escrows
func (k Keeper) IterateEscrows(
	ctx sdk.Context,
	op func(id tmbytes.HexBytes, h types.Escrow) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.EscrowKey)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := tmbytes.HexBytes(iterator.Key()[1:])

		var escrow types.Escrow
		k.cdc.MustUnmarshalBinaryBare(iterator.Value(), &escrow)

		if stop := op(id, escrow); stop {
			break
		}
	}
}

// IterateEscrowExpiredQueueByHeight iterates through the escrow expiration queue by the specified height
func (k Keeper) IterateEscrowExpiredQueueByHeight(
	ctx sdk.Context, height uint64,
	op func(id tmbytes.HexBytes, h types.Escrow) (stop bool),
) {
	store := ctx.KVStore(k.storeKey)

	iterator := sdk.KVStorePrefixIterator(store, types.GetEscrowExpiredQueueSubspace(height))
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		id := tmbytes.HexBytes(iterator.Key()[9:])
		escrow, _ := k.GetEscrow(ctx, id)

		if stop := op(id, escrow); stop {
			break
		}
	}
}
