package keeper

import (
	"bytes"
	"reflect"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	crud "github.com/iov-one/cosmos-sdk-crud"
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
	string,
	error,
) {
	id := k.FetchNextId()

	//TODO: shouldn't this be hex-encoding ?
	stringId := string(id)

	// check if the escrow already exists
	if k.HasEscrow(ctx, stringId) {
		return "", sdkerrors.Wrap(types.ErrIdNotUnique, id.String())
	}

	// Create and validate the escrow
	escrow := types.NewEscrow(
		id, buyer, seller, price, object, deadline,
	)
	err := escrow.Validate()
	if err != nil {
		return "", err
	}

	// Retrieve the store for this object
	objectStore, err := k.getStoreForID(object.GetType())
	if err != nil {
		return "", err
	}

	// Check the object is in the store and is equal to the store's version
	err = k.checkObjectWithStore(objectStore, object)
	if err != nil {
		return "", err
	}

	// transfer ownership of the object to the escrow account
	err = k.doObjectTransferWithStore(seller, k.GetEscrowAccount(ctx).GetAddress(), object, objectStore)
	if err != nil {
		return "", errors.Wrap(err, "Cannot transfer the object to the module account")
	}

	// save the modified object
	err = objectStore.Update(object.GetObject())
	if err != nil {
		return "", err
	}

	// save the escrow
	k.saveEscrow(ctx, escrow)
	k.NextId()
	//TODO: Emit event

	return stringId, nil
}

func (k Keeper) UpdateEscrow(
	ctx sdk.Context,
	id string,
	updater sdk.AccAddress,
	seller sdk.AccAddress,
	buyer sdk.AccAddress,
	price sdk.Coins,
	deadline uint64,
) error {
	// check if the escrow exists
	_, found := k.GetEscrow(ctx, id)
	if !found {
		return sdkerrors.Wrap(types.ErrEscrowNotFound, id)
	}

	//TODO: take care of deadline store if modifying deadline !
	//TODO: implement the update method
	return nil
}

// TransferToEscrow transfers coins from the buyer to the escrow account
func (k Keeper) TransferToEscrow(
	ctx sdk.Context,
	sender sdk.AccAddress,
	id string,
	amount sdk.Coins,

) error {
	// check if the escrow exists
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return sdkerrors.Wrap(types.ErrEscrowNotFound, id)
	}

	// check if the escrow is open
	if escrow.State != types.EscrowState_Open {
		return sdkerrors.Wrap(types.ErrEscrowNotOpen, escrow.Id)
	}

	//TODO: check if deadline here is needed or if BeginBlocker is a sufficient warranty

	buyer, err := sdk.AccAddressFromBech32(escrow.Buyer)
	//TODO: this should be always valid because escrow is guaranteed to be in a valid state when created/updated
	if err != nil {
		return sdkerrors.Wrapf(err, "Invalid buyer address : %v", escrow.Buyer)
	}

	// Ensure that the sender is the designated buyer of this escrow
	if !bytes.Equal(buyer, sender) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Only the designated buyer can transfer coins to the escrow")
	}

	seller, err := sdk.AccAddressFromBech32(escrow.Seller)
	//TODO: this should be always valid because escrow is guaranteed to be in a valid state when created/updated
	if err != nil {
		return sdkerrors.Wrapf(err, "Invalid seller address : %v", escrow.Seller)
	}
	// Check if the provided amount is valid
	if !amount.IsValid() {
		return types.ErrInvalidAmount
	}

	// Check if the amount is greater or equal than the price
	if !amount.IsAllGTE(escrow.Price) {
		return types.ErrTransferAmountTooLow
	}

	// Send the price to the module
	err = k.bankKeeper.SendCoinsFromAccountToModule(ctx, buyer, types.ModuleName, escrow.Price)
	if err != nil {
		return sdkerrors.Wrap(err, "Cannot send the coins to the escrow")
	}

	// Do the exchange
	err = k.doSwap(ctx, escrow, buyer, seller)
	//TODO: should we try to recover ? (e.g. sending back coins / object and closing the escrow on transfer failure)
	if err != nil {
		panic(err)
	}

	escrow.State = types.EscrowState_Completed
	k.deleteEscrow(ctx, escrow)
	//TODO: Emit event

	return nil
}

func (k Keeper) doSwap(ctx sdk.Context, escrow types.Escrow, buyer, seller sdk.AccAddress) error {

	// Transfer the object from the module to the buyer
	err := k.doObjectTransfer(k.GetEscrowAccount(ctx).GetAddress(), buyer, escrow.GetObject())
	if err != nil {
		return sdkerrors.Wrap(err, "Cannot send the object to the buyer")
	}

	// Transfer the coins
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, seller, escrow.Price)
	if err != nil {
		return sdkerrors.Wrap(err, "Cannot send the coins to the seller")
	}
	return nil
}

// RefundEscrow refunds the specified escrow
func (k Keeper) RefundEscrow(ctx sdk.Context, sender sdk.AccAddress, id string) error {

	// check if the escrow exists
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return sdkerrors.Wrap(types.ErrEscrowNotFound, id)
	}

	if escrow.State != types.EscrowState_Open {
		return sdkerrors.Wrap(types.ErrEscrowNotOpen, escrow.Id)
	}

	//TODO: check if deadline here is needed or if BeginBlocker is a sufficient warranty

	seller, err := sdk.AccAddressFromBech32(escrow.Seller)
	if err != nil {
		return err
	}

	// Ensure the seller is the one asking for a refund
	if !bytes.Equal(sender, seller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Only the designated seller can trigger a refund")
	}

	if err := k.refundEscrow(ctx, escrow, seller); err != nil {
		return err
	}

	//TODO: Emit event

	return nil
}

func (k Keeper) refundEscrow(ctx sdk.Context, escrow types.Escrow, seller sdk.AccAddress) error {

	// Transfer the object back to the seller
	err := k.doObjectTransfer(k.GetEscrowAccount(ctx).GetAddress(), seller, escrow.GetObject())
	if err != nil {
		return sdkerrors.Wrap(err, "Error while transferring the object back to the seller")

	}

	// update the state of the escrow
	escrow.State = types.EscrowState_Refunded
	// delete escrow
	k.deleteEscrow(ctx, escrow)
	return nil
}

func (k Keeper) doObjectTransfer(from, to sdk.AccAddress, object types.TransferableObject) error {
	// Retrieve the object store
	objectStore, err := k.getStoreForID(object.GetType())
	if err != nil {
		return err
	}
	return k.doObjectTransferWithStore(from, to, object, objectStore)
}

func (k Keeper) doObjectTransferWithStore(from, to sdk.AccAddress, object types.TransferableObject, objectStore crud.Store) error {
	// Transfer the object
	err := object.Transfer(from, to)
	if err != nil {
		return err
	}

	// Save the object in its store
	err = objectStore.Update(object.GetObject())
	if err != nil {
		return err
	}
	return nil
}

func (k Keeper) getStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.EscrowStoreKey)
}

func (k Keeper) getDeadlineStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), types.DeadlineStoreKey)
}

func (k Keeper) addEscrowToDeadlineStore(ctx sdk.Context, escrow types.Escrow) {
	k.getDeadlineStore(ctx).Set(types.GetDeadlineKey(escrow.Deadline, escrow.Id), types.GetEscrowKey(escrow.Id))
}

func (k Keeper) deleteEscrowFromDeadlineStore(ctx sdk.Context, escrow types.Escrow) {
	k.getDeadlineStore(ctx).Delete(types.GetDeadlineKey(escrow.Deadline, escrow.Id))
}

func (k Keeper) checkObjectWithStore(objectStore crud.Store, object types.TransferableObject) error {
	var objInStore crud.Object
	err := objectStore.Read(object.GetObject().PrimaryKey(), objInStore)
	if err != nil {
		return types.ErrUnknownObject
	}

	if !reflect.DeepEqual(objInStore, object) {
		return sdkerrors.Wrap(types.ErrUnknownObject, "The object and his stored version does not match")
	}

	return nil
}

// HasEscrow checks if the given escrow exists
func (k Keeper) HasEscrow(ctx sdk.Context, id string) bool {
	return k.getStore(ctx).Has(types.GetEscrowKey(id))
}

// saveEscrow sets the given escrow
func (k Keeper) saveEscrow(ctx sdk.Context, escrow types.Escrow) {
	bz := k.cdc.MustMarshalBinaryBare(&escrow)
	k.getStore(ctx).Set(types.GetEscrowKey(escrow.Id), bz)
	k.addEscrowToDeadlineStore(ctx, escrow)
}

func (k Keeper) deleteEscrow(ctx sdk.Context, escrow types.Escrow) {
	if escrow.State == types.EscrowState_Open {
		panic("Attempted to delete an open escrow")
	}

	k.getStore(ctx).Delete(types.GetEscrowKey(escrow.Id))
	k.deleteEscrowFromDeadlineStore(ctx, escrow)
}

// GetEscrow retrieves the specified escrow
func (k Keeper) GetEscrow(ctx sdk.Context, id string) (escrow types.Escrow, found bool) {
	return k.getEscrowByKey(ctx, types.GetEscrowKey(id))
}

// getEscrowImpl retrieves the specified escrow with its key
func (k Keeper) getEscrowByKey(ctx sdk.Context, key []byte) (escrow types.Escrow, found bool) {
	bz := k.getStore(ctx).Get(key)
	if bz == nil {
		return escrow, false
	}
	k.cdc.MustUnmarshalBinaryBare(bz, &escrow)
	return escrow, true
}

func (k Keeper) IterateExpiredEscrows(ctx sdk.Context, date uint64, op func(types.Escrow) bool) {
	store := k.getDeadlineStore(ctx)
	// TODO : check if that's actually valid
	end := sdk.Uint64ToBigEndian(date + 1)
	iterator := store.Iterator(nil, end)
	defer iterator.Close()

	for ; iterator.Valid(); iterator.Next() {
		escrow, found := k.getEscrowByKey(ctx, iterator.Value())
		if !found {
			panic("Inconstancy in expired escrows store : escrow not found")
		}
		if stop := op(escrow); stop {
			break
		}
	}

}

func (k Keeper) RefundExpiredEscrows(ctx sdk.Context, date uint64) {
	k.IterateExpiredEscrows(ctx, date,
		func(e types.Escrow) (stop bool) {
			//TODO: check if allowed because we modify expired store (refund -> delete escrow -> delete escrow from expired store)
			// while iterating over it

			// refund escrow
			seller, err := sdk.AccAddressFromBech32(e.Seller)
			if err != nil {
				panic(err)
			}
			err = k.refundEscrow(ctx, e, seller)
			if err != nil {
				panic(err)
			}
			return false
		})
}

// IterateEscrows iterates through the escrows
func (k Keeper) IterateEscrows(
	ctx sdk.Context,
	op func(tmbytes.HexBytes, types.Escrow) bool,
) {
	store := k.getStore(ctx)

	iterator := store.Iterator(nil, nil)
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
