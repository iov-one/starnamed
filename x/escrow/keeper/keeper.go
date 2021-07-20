package keeper

import (
	"encoding/hex"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/store"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	crudtypes "github.com/iov-one/cosmos-sdk-crud/types"

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
	customData          map[types.TypeID]types.CustomData
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
		customData:          make(map[types.TypeID]types.CustomData),
		blockedAddrs:        blockedAddrs,
	}
}

// RegisterCustomData registers custom data to be given to the Transfer function of a certain type of TransferableObject
func (k Keeper) RegisterCustomData(id types.TypeID, data types.CustomData) {
	k.customData[id] = data
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

func (k Keeper) getCustomDataForType(id types.TypeID) types.CustomData {
	return k.customData[id]
}

func (k Keeper) getEscrowStore(ctx sdk.Context) crud.Store {
	return crudtypes.NewStore(k.cdc, ctx.KVStore(k.storeKey), EscrowStoreKey)
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

// GetMaximumEscrowDuration returns the maximum allowed duration of an escrow
func (k Keeper) GetMaximumEscrowDuration(ctx sdk.Context) time.Duration {
	return k.configurationKeeper.GetConfiguration(ctx).EscrowMaxPeriod
}

// GetEscrowPriceDenom returns the denomination of the allowed token for the price fo an escrow
func (k Keeper) GetEscrowPriceDenom(ctx sdk.Context) string {
	return k.configurationKeeper.GetFees(ctx).FeeCoinDenom
}

// GetBrokerAddress returns the escrow broker address
func (k Keeper) GetBrokerAddress(ctx sdk.Context) string {
	return k.configurationKeeper.GetConfiguration(ctx).EscrowBroker
}

// GetBrokerCommission returns the escrow broker commission
func (k Keeper) GetBrokerCommission(ctx sdk.Context) sdk.Dec {
	return k.configurationKeeper.GetConfiguration(ctx).EscrowCommission

}

func (k Keeper) isBlockedAddr(address string) bool {
	return k.blockedAddrs[address]
}
