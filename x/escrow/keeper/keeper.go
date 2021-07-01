package keeper

import (
	"encoding/hex"
	"fmt"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"

	crud "github.com/iov-one/cosmos-sdk-crud"

	"github.com/tendermint/tendermint/libs/log"

	"github.com/iov-one/starnamed/x/escrow/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

var (
	// Keys for store prefixes
	EscrowStoreKey   = []byte{0x01} // prefix for escrow
	DeadlineStoreKey = []byte{0x02} // prefix for escrow stored by expiration date
	ParamsStoreKey   = []byte{0x03} // prefix for the keeper parameters

	// Keys for the parameters store
	paramsStoreLastBlockTime = []byte{0x01}
	paramsStoreNextId        = []byte{0x02}
)

// Keeper defines the escrow keeper
type Keeper struct {
	storeKey            sdk.StoreKey
	cdc                 codec.Marshaler
	paramSpace          paramstypes.Subspace
	accountKeeper       types.AccountKeeper
	bankKeeper          types.BankKeeper
	configurationKeeper types.ConfigurationKeeper
	storeHolders        map[types.TypeID]types.StoreHolder
	blockedAddrs        map[string]bool
}

// NewKeeper creates a new escrow Keeper instance
func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	configurationKeeper types.ConfigurationKeeper,
	blockedAddrs map[string]bool,
) Keeper {
	// ensure the escrow module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:            key,
		cdc:                 cdc,
		paramSpace:          paramSpace,
		accountKeeper:       accountKeeper,
		bankKeeper:          bankKeeper,
		configurationKeeper: configurationKeeper,
		storeHolders:        make(map[types.TypeID]types.StoreHolder),
		blockedAddrs:        blockedAddrs,
	}
}

// AddStore registers a simple store holder for the specified type id, which always retrieves the given store
func (k Keeper) AddStore(id types.TypeID, store crud.Store) {
	k.AddStoreHolder(id, types.NewSimpleStoreHolder(func(sdk.Context) crud.Store { return store }))
}

// AddStoreHolder registers a store holder for the specified type id
func (k Keeper) AddStoreHolder(id types.TypeID, store types.StoreHolder) {
	if _, alreadyPresent := k.storeHolders[id]; alreadyPresent {
		panic(fmt.Errorf("cannot register a store holder for type id %v because it is already registered", id))
	}
	k.storeHolders[id] = store
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetEscrowAddress returns the escrow module address
func (k Keeper) GetEscrowAddress() sdk.AccAddress {
	return k.accountKeeper.GetModuleAddress(types.ModuleName)
}

func (k Keeper) ImportNextID(ctx sdk.Context, nextID uint64) {
	k.getParamStore(ctx).Set(paramsStoreNextId, sdk.Uint64ToBigEndian(nextID))
}

func (k Keeper) FetchNextId(ctx sdk.Context) string {
	return hex.EncodeToString(k.getParamStore(ctx).Get(paramsStoreNextId))
}

func (k Keeper) GetNextIDForExport(ctx sdk.Context) uint64 {
	return sdk.BigEndianToUint64(k.getParamStore(ctx).Get(paramsStoreNextId))
}

func (k Keeper) NextId(ctx sdk.Context) {
	next := k.GetNextIDForExport(ctx) + 1
	k.getParamStore(ctx).Set(paramsStoreNextId, sdk.Uint64ToBigEndian(next))
}

func (k Keeper) getStoreForID(ctx sdk.Context, id types.TypeID) (crud.Store, error) {
	store, ok := k.storeHolders[id]
	if !ok {
		return nil, types.ErrUnknownTypeID
	}
	return store.GetCRUDStore(ctx), nil
}

func (k Keeper) getStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), EscrowStoreKey)
}

func (k Keeper) getDeadlineStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), DeadlineStoreKey)
}

func (k Keeper) getParamStore(ctx sdk.Context) store.KVStore {
	return prefix.NewStore(ctx.KVStore(k.storeKey), ParamsStoreKey)
}

func (k Keeper) SetLastBlockTime(ctx sdk.Context, date uint64) {
	k.getParamStore(ctx).Set(paramsStoreLastBlockTime, sdk.Uint64ToBigEndian(date))
}

func (k Keeper) GetLastBlockTime(ctx sdk.Context) uint64 {
	return sdk.BigEndianToUint64(k.getParamStore(ctx).Get(paramsStoreLastBlockTime))
}

func (k Keeper) isBlockedAddr(address string) bool {
	return k.blockedAddrs[address]
}
