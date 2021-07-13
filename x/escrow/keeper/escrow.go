package keeper

import (
	"encoding/hex"
	"reflect"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/pkg/errors"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// CreateEscrow creates an escrow
func (k Keeper) CreateEscrow(
	ctx sdk.Context,
	seller sdk.AccAddress,
	price sdk.Coins,
	object types.TransferableObject,
	deadline uint64,
) (
	string,
	error,
) {
	id := k.FetchNextId(ctx)

	// check if the escrow already exists
	if k.HasEscrow(ctx, id) {
		return "", sdkerrors.Wrap(types.ErrIdNotUnique, id)
	}

	// Retrieve the store for this object
	objectStore, err := k.getStoreForID(ctx, object.GetType())
	if err != nil {
		return "", err
	}
	// Check the object is in the store and is equal to the store's version
	err = k.checkObjectWithStore(objectStore, object)
	if err != nil {
		return "", err
	}

	// Check the deadline validity
	if deadline > uint64(ctx.BlockTime().Unix())+uint64(k.GetMaximumEscrowDuration(ctx).Seconds()) {
		return "", sdkerrors.Wrap(types.ErrInvalidDeadline, "The deadline exceeds the maximum escrow duration")
	}

	//TODO: use cosmos-sdk params instead of the configuration

	// Get the configuration
	config := k.configurationKeeper.GetConfiguration(ctx)

	// Create and validate the escrow
	escrow := types.NewEscrow(
		id, seller, price, object, deadline, config.EscrowBroker, config.EscrowCommission,
	)
	err = escrow.Validate(k.GetEscrowPriceDenom(ctx), k.GetLastBlockTime(ctx))
	if err != nil {
		return "", err
	}

	// transfer ownership of the object to the escrow account
	err = k.doObjectTransferWithStore(seller, k.GetEscrowAddress(), object, objectStore)
	if err != nil {
		return "", errors.Wrap(err, "Cannot transfer the object to the module account")
	}

	// save the escrow
	k.SaveEscrow(ctx, escrow)
	k.NextId(ctx)
	//TODO: Emit event

	return id, nil
}

func (k Keeper) UpdateEscrow(
	ctx sdk.Context,
	id string,
	updater sdk.AccAddress,
	seller sdk.AccAddress,
	price sdk.Coins,
	deadline uint64,
) error {
	// check that the escrow exists
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return sdkerrors.Wrap(types.ErrEscrowNotFound, id)
	}

	// check that the escrow is open
	if escrow.State != types.EscrowState_Open {
		return sdkerrors.Wrap(types.ErrEscrowExpired, id)
	}

	// Check if there is an actual update
	if seller == nil && price == nil && deadline == 0 {
		return types.ErrEmptyUpdate
	}

	oldSeller, err := sdk.AccAddressFromBech32(escrow.Seller)
	if err != nil {
		// this should be always valid because escrow is guaranteed to be in a valid state when created/updated
		panic(sdkerrors.Wrapf(err, "Invalid seller address : %v", escrow.Seller))
	}

	// If the updater is not the seller he cannot update the escrow
	if !updater.Equals(oldSeller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Only the seller can update an escrow")
	}

	// Update seller, price and deadline if provided
	if seller != nil {
		escrow.Seller = seller.String()
	}
	if price != nil {
		if err := types.ValidatePrice(price, k.GetEscrowPriceDenom(ctx)); err != nil {
			return err
		}
		escrow.Price = price
	}
	if deadline != 0 {
		if err := types.ValidateDeadline(deadline, k.GetLastBlockTime(ctx)); err != nil {
			return err
		}
		//TODO: is that the behavior we want ?

		// We can only delay the deadline and no more than the maximum duration from now
		if deadline < escrow.Deadline {
			return sdkerrors.Wrap(types.ErrInvalidDeadline, "The new deadline cannot be before the old one")
		} else if deadline > uint64(ctx.BlockTime().Unix())+uint64(k.GetMaximumEscrowDuration(ctx).Seconds()) {
			return sdkerrors.Wrap(types.ErrInvalidDeadline, "The new deadline exceeds the maximum escrow duration")
		}
		// We are modifying the deadline, get rid of old deadline indexing
		k.deleteEscrowFromDeadlineStore(ctx, escrow)
		// The new deadline indexing will be added when we save the escrow
		escrow.Deadline = deadline
	}

	k.SaveEscrow(ctx, escrow)
	return nil
}

// TransferToEscrow transfers coins from the buyer to the escrow account
func (k Keeper) TransferToEscrow(
	ctx sdk.Context,
	buyer sdk.AccAddress,
	id string,
	amount sdk.Coins,

) error {
	// check that the escrow exists
	escrow, found := k.GetEscrow(ctx, id)
	if !found {
		return sdkerrors.Wrap(types.ErrEscrowNotFound, id)
	}

	// check that the escrow is open
	if escrow.State != types.EscrowState_Open {
		return sdkerrors.Wrap(types.ErrEscrowNotOpen, escrow.Id)
	}

	seller, err := sdk.AccAddressFromBech32(escrow.Seller)
	if err != nil {
		//this should be always valid because the escrow is guaranteed to be in a valid state when created/updated
		panic(sdkerrors.Wrapf(err, "Invalid seller address : %v", escrow.Seller))
	}

	broker, err := sdk.AccAddressFromBech32(escrow.BrokerAddress)
	if err != nil {
		//this should be always valid because the escrow is guaranteed to be in a valid state when created
		panic(sdkerrors.Wrapf(err, "Invalid broker address : %v", escrow.BrokerAddress))
	}

	// Ensure that the buyer is not the seller of this escrow
	if buyer.Equals(seller) {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "The owner of the escrow cannot transfer coins to the escrow")
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
	err = k.doSwap(ctx, escrow, buyer, seller, broker)
	// If an error occurs here, the buyer have sent the coins and :
	// - The buyer can have received the object or not
	// - The seller has not received the coins
	// This case should never happen because the escrow should possess the coins and the object
	//TODO: or maybe some other external problems could make this transaction fail
	// how can we recover in that case ?
	if err != nil {
		panic(err)
	}

	escrow.State = types.EscrowState_Completed
	k.deleteEscrow(ctx, escrow)
	//TODO: Emit event

	return nil
}

func (k Keeper) doSwap(ctx sdk.Context, escrow types.Escrow, buyer, seller sdk.AccAddress, broker sdk.AccAddress) error {

	// Transfer the object from the module to the buyer
	err := k.doObjectTransfer(ctx, k.GetEscrowAddress(), buyer, escrow.GetObject())
	if err != nil {
		return sdkerrors.Wrap(err, "Cannot send the object to the buyer")
	}

	// Transfer the coins, making sure that brokerCoins + sellerCoins = escrow.Price
	brokerCoins, _ := sdk.NewDecCoinsFromCoins(escrow.Price...).MulDec(escrow.BrokerCommission).TruncateDecimal()
	sellerCoins := escrow.Price.Sub(brokerCoins)

	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, broker, brokerCoins)
	if err != nil {
		return sdkerrors.Wrap(err, "Cannot send the coins to the broker")
	}
	err = k.bankKeeper.SendCoinsFromModuleToAccount(ctx, types.ModuleName, seller, sellerCoins)
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

	if escrow.State != types.EscrowState_Open && escrow.State != types.EscrowState_Expired {
		return sdkerrors.Wrap(types.ErrEscrowNotOpen, escrow.Id)
	}

	seller, err := sdk.AccAddressFromBech32(escrow.Seller)
	if err != nil {
		return err
	}

	// Ensure the seller is the one asking for a refund or that escrow is expired
	if !sender.Equals(seller) && escrow.State != types.EscrowState_Expired {
		return sdkerrors.Wrap(sdkerrors.ErrUnauthorized, "Only the escrow owner can trigger a refund if the escrow is not expired")
	}

	// We refund the object to the seller
	if err := k.refundEscrow(ctx, escrow, seller); err != nil {
		return err
	}

	//TODO: Emit event

	return nil
}

func (k Keeper) refundEscrow(ctx sdk.Context, escrow types.Escrow, seller sdk.AccAddress) error {

	// Transfer the object back to the seller
	err := k.doObjectTransfer(ctx, k.GetEscrowAddress(), seller, escrow.GetObject())
	if err != nil {
		return sdkerrors.Wrap(err, "Error while transferring the object back to the seller")

	}

	// update the state of the escrow
	escrow.State = types.EscrowState_Refunded
	// delete escrow
	k.deleteEscrow(ctx, escrow)
	return nil
}

func (k Keeper) doObjectTransfer(ctx sdk.Context, from, to sdk.AccAddress, object types.TransferableObject) error {
	// Retrieve the object store
	objectStore, err := k.getStoreForID(ctx, object.GetType())
	if err != nil {
		return err
	} // Retrieve the store for this object

	return k.doObjectTransferWithStore(from, to, object, objectStore)
}

func (k Keeper) doObjectTransferWithStore(from, to sdk.AccAddress, object types.TransferableObject, objectStore crud.Store) error {
	// Transfer the object
	err := object.Transfer(from, to)
	if err != nil {
		return err
	}

	// Save the object in its store
	return objectStore.Update(object.GetObject())
}

func (k Keeper) addEscrowToDeadlineStore(ctx sdk.Context, escrow types.Escrow) {
	k.getDeadlineStore(ctx).Set(types.GetDeadlineKey(escrow.Deadline, escrow.Id), types.GetEscrowKey(escrow.Id))
}

func (k Keeper) deleteEscrowFromDeadlineStore(ctx sdk.Context, escrow types.Escrow) {
	k.getDeadlineStore(ctx).Delete(types.GetDeadlineKey(escrow.Deadline, escrow.Id))
}

//FIXME: this function syncs the object with the version in store but does not check anything
func (k Keeper) checkObjectWithStore(objectStore crud.Store, object types.TransferableObject) error {
	var objInStore = object.GetObject()

	err := objectStore.Read(objInStore.PrimaryKey(), objInStore)
	if err != nil {
		return types.ErrUnknownObject
	}

	//FIXME: this will be always true
	if !reflect.DeepEqual(objInStore, object.GetObject()) {
		return sdkerrors.Wrap(types.ErrUnknownObject, "The object and his stored version does not match")
	}

	return nil
}

// HasEscrow checks if the given escrow exists
func (k Keeper) HasEscrow(ctx sdk.Context, id string) bool {
	var escrow types.Escrow
	err := k.getEscrowStore(ctx).Read(types.GetEscrowKey(id), &escrow)
	return err == nil
}

// SaveEscrow sets the given escrow
func (k Keeper) SaveEscrow(ctx sdk.Context, escrow types.Escrow) {
	escrow.SyncObject()
	if k.HasEscrow(ctx, escrow.Id) {
		if err := k.getEscrowStore(ctx).Update(&escrow); err != nil {
			panic(err)
		}
	} else {
		if err := k.getEscrowStore(ctx).Create(&escrow); err != nil {
			panic(err)
		}
	}
	k.addEscrowToDeadlineStore(ctx, escrow)
}

func (k Keeper) deleteEscrow(ctx sdk.Context, escrow types.Escrow) {
	if escrow.State == types.EscrowState_Open {
		panic("Attempted to delete an open escrow")
	}
	if escrow.State == types.EscrowState_Expired {
		panic("Attempted to delete an expired escrow without refunding it")
	}

	if err := k.getEscrowStore(ctx).Delete(escrow.PrimaryKey()); err != nil {
		panic(err)
	}
	k.deleteEscrowFromDeadlineStore(ctx, escrow)
}

// GetEscrow retrieves the specified escrow
func (k Keeper) GetEscrow(ctx sdk.Context, id string) (escrow types.Escrow, found bool) {
	return k.getEscrowByKey(ctx, types.GetEscrowKey(id))
}

func consumeEscrowCursor(cursor crud.Cursor) ([]types.Escrow, error) {
	var escrows []types.Escrow
	for ; cursor.Valid(); cursor.Next() {
		var escrow types.Escrow
		if err := cursor.Read(&escrow); err != nil {
			return nil, err
		}
		escrows = append(escrows, escrow)
	}
	return escrows, nil
}

func (k Keeper) QueryEscrows(ctx sdk.Context, filter func(crud.QueryStatement) crud.ValidQuery) ([]types.Escrow, error) {
	cursor, err := filter(k.getEscrowStore(ctx).Query()).Do()
	if err != nil {
		return nil, err
	}
	return consumeEscrowCursor(cursor)
}

func (k Keeper) QueryEscrowsWithRange(ctx sdk.Context, filter func(crud.QueryStatement) crud.ValidQuery, start, end uint64) ([]types.Escrow, error) {
	return k.QueryEscrows(ctx, func(query crud.QueryStatement) crud.ValidQuery {
		return filter(query).WithRange().Start(start).End(end)
	})
}

func (k Keeper) GetEscrowsBySeller(ctx sdk.Context, seller string, start, end uint64) ([]types.Escrow, error) {
	sellerAddr, err := sdk.AccAddressFromBech32(seller)
	if err != nil {
		return nil, err
	}
	return k.QueryEscrowsWithRange(ctx, func(query crud.QueryStatement) crud.ValidQuery {
		return query.Where().Index(types.SellerIndex).Equals(sellerAddr)
	}, start, end)

}

func (k Keeper) GetEscrowsByState(ctx sdk.Context, state types.EscrowState, start, end uint64) ([]types.Escrow, error) {
	return k.QueryEscrowsWithRange(ctx, func(query crud.QueryStatement) crud.ValidQuery {
		return query.Where().Index(types.StateIndex).Equals(sdk.Uint64ToBigEndian(uint64(state)))
	}, start, end)
}

func (k Keeper) GetEscrowsByObject(ctx sdk.Context, object types.TransferableObject) ([]types.Escrow, error) {
	return k.QueryEscrows(ctx, func(query crud.QueryStatement) crud.ValidQuery {
		return query.Where().Index(types.ObjectIndex).Equals(object.GetObject().PrimaryKey())
	})
}

func (k Keeper) queryEscrowsByAttributes(
	ctx sdk.Context,
	sellerStr string,
	stateStr string,
	objectKeyStr string,
	start uint64,
	length uint64,
) ([]types.Escrow, error) {
	var seller sdk.AccAddress
	if len(sellerStr) != 0 {
		var err error
		seller, err = sdk.AccAddressFromBech32(sellerStr)
		if err != nil {
			return nil, sdkerrors.Wrap(err, "Invalid seller address")
		}
	}

	var state types.EscrowState
	var hasState bool
	if len(stateStr) != 0 {
		hasState = true
		switch strings.ToLower(stateStr) {
		case "open":
			state = types.EscrowState_Open
		case "expired":
			state = types.EscrowState_Expired
		default:
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "The state is invalid, can be open or expired")
		}
	}

	var objectKey []byte
	if len(objectKeyStr) != 0 {
		var err error
		objectKey, err = hex.DecodeString(objectKeyStr)
		if err != nil {
			return nil, sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "Object key must be an hex-encoded byte array : "+err.Error())
		}
	}

	var end uint64
	if length == 0 {
		end = 0
	} else {
		end = start + length
	}

	filter := func(query crud.QueryStatement) crud.ValidQuery {
		getStatement := func(query crud.QueryStatement, previous crud.FinalizedIndexStatement) crud.WhereStatement {
			if previous == nil {
				return query.Where()
			} else {
				return previous.And()
			}
		}

		previousStatement := crud.FinalizedIndexStatement(nil)
		//TODO: maybe optimize by filtering first by object if there is an object filter
		if seller != nil {
			previousStatement = getStatement(query, previousStatement).
				Index(types.SellerIndex).Equals(seller)
		}
		if hasState {
			previousStatement = getStatement(query, previousStatement).
				Index(types.StateIndex).Equals(sdk.Uint64ToBigEndian(uint64(state)))
		}
		if objectKey != nil {
			previousStatement = getStatement(query, previousStatement).
				Index(types.ObjectIndex).Equals(objectKey)
		}

		if previousStatement == nil {
			return query
		} else {
			return previousStatement
		}
	}

	return k.QueryEscrowsWithRange(ctx, filter, start, end)
}

// getEscrowByKey retrieves the specified escrow with its key
func (k Keeper) getEscrowByKey(ctx sdk.Context, key []byte) (escrow types.Escrow, found bool) {
	err := k.getEscrowStore(ctx).Read(key, &escrow)
	if errors.Is(err, crud.ErrNotFound) {
		return escrow, false
	} else if err != nil {
		panic(err)
	}
	return escrow, true
}

func (k Keeper) IterateEscrowsWithPassedDeadline(ctx sdk.Context, date uint64, op func(types.Escrow) bool) {
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
func (k Keeper) MarkExpiredEscrows(ctx sdk.Context, date uint64) {
	k.IterateEscrowsWithPassedDeadline(ctx, date,
		func(e types.Escrow) (stop bool) {
			if e.State == types.EscrowState_Open {
				e.State = types.EscrowState_Expired
				k.SaveEscrow(ctx, e)
			}
			return false
		})
}

func (k Keeper) RefundExpiredEscrows(ctx sdk.Context) {
	k.IterateEscrows(ctx,
		func(e types.Escrow) (stop bool) {
			//TODO: check if allowed because we modify expired store (refund -> delete escrow -> delete escrow from expired store)
			// while iterating over it

			if e.State != types.EscrowState_Expired {
				return false
			}

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
	op func(types.Escrow) bool,
) {
	store := k.getEscrowStore(ctx)
	cursor, err := store.Query().Do()
	if err != nil {
		panic(err)
	}

	for ; cursor.Valid(); cursor.Next() {

		var escrow types.Escrow
		if err := cursor.Read(&escrow); err != nil {
			panic(err)
		}

		if stop := op(escrow); stop {
			break
		}
	}
}
