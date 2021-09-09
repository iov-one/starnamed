package keeper

import (
	"encoding/binary"
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	crudtypes "github.com/iov-one/cosmos-sdk-crud/types"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
	"github.com/tendermint/tendermint/libs/log"
)

// ParamSubspace is a placeholder
type ParamSubspace interface {
}

// list expected keepers

// SupplyKeeper defines the behaviour
// of the supply keeper used to collect
// and then distribute the fees
type SupplyKeeper interface {
	SendCoinsFromAccountToModule(ctx sdk.Context, addr sdk.AccAddress, moduleName string, coins sdk.Coins) error
	GetAllBalances(ctx sdk.Context, addr sdk.AccAddress) sdk.Coins
}

// AuthKeeper defines the behavious of the auth keeper used to get module addresses
type AuthKeeper interface {
	GetModuleAddress(module string) sdk.AccAddress
}

// DistributionKeeper is used to estimate the yield for the delegators
type DistributionKeeper interface {
	GetCommunityTax(ctx sdk.Context) sdk.Dec
}

// StakingKeeper is used to estimate the yield for delegators
type StakingKeeper interface {
	GetLastTotalPower(ctx sdk.Context) sdk.Int
}

// ConfigurationKeeper defines the behaviour of the configuration state checks
type ConfigurationKeeper interface {
	// GetFees gets the fees
	GetFees(ctx sdk.Context) *configuration.Fees
	// GetConfiguration returns the configuration
	GetConfiguration(ctx sdk.Context) configuration.Config
	// IsOwner returns if the provided address is an owner or not
	IsOwner(ctx sdk.Context, addr sdk.AccAddress) bool
	// GetValidDomainNameRegexp returns the regular expression that aliceAddr domain name must match
	// in order to be valid
	GetValidDomainNameRegexp(ctx sdk.Context) string
	// GetDomainRenewDuration returns the default duration of aliceAddr domain renewal
	GetDomainRenewDuration(ctx sdk.Context) time.Duration
	// GetDomainGracePeriod returns the grace period duration
	GetDomainGracePeriod(ctx sdk.Context) time.Duration
}

// Keeper of the domain store
// TODO split this keeper in sub-struct in order to avoid possible mistakes with keys and not clutter the exposed methods
type Keeper struct {
	// external keepers
	ConfigurationKeeper ConfigurationKeeper
	SupplyKeeper        SupplyKeeper
	AuthKeeper          AuthKeeper
	StakingKeeper       StakingKeeper
	DistributionKeeper  DistributionKeeper
	// default fields
	StoreKey   sdk.StoreKey // contains the store key for the domain module
	Cdc        codec.Marshaler
	paramspace ParamSubspace
	// crud stores
	accountStore crud.Store
	domainStore  crud.Store
}

// NewKeeper creates a domain keeper
func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, configKeeper ConfigurationKeeper, supply SupplyKeeper, auth AuthKeeper, distrib DistributionKeeper, staking StakingKeeper, paramspace ParamSubspace) Keeper {
	keeper := Keeper{
		StoreKey:            storeKey,
		Cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		SupplyKeeper:        supply,
		AuthKeeper:          auth,
		DistributionKeeper:  distrib,
		StakingKeeper:       staking,
		paramspace:          paramspace,
	}
	return keeper
}

// AccountStore returns the crud.Store used to interact with account objects
func (k Keeper) AccountStore(ctx sdk.Context) crud.Store {
	return crudtypes.NewStore(k.Cdc, ctx.KVStore(k.StoreKey), []byte{0x1})
}

// DomainStore returns the crud.Store used to interact with domain objects
func (k Keeper) DomainStore(ctx sdk.Context) crud.Store {
	return crudtypes.NewStore(k.Cdc, ctx.KVStore(k.StoreKey), []byte{0x2})
}

// TODO: move this in a separate module
// feesStore return the store storing the fees for each block
func (k Keeper) feesStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(k.StoreKey), []byte{0x3})
}

// StoreBlockFees stores the fees for the current block height
func (k Keeper) StoreBlockFees(ctx sdk.Context, fees sdk.Coins) {
	//TODO: duplicated code for key retrieval
	key := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))

	bytes := sdk.Uint64ToBigEndian(uint64(len(fees)))
	for _, fee := range fees {
		bytes = append(bytes, k.Cdc.MustMarshalBinaryLengthPrefixed(&fee)...)
	}
	k.feesStore(ctx).Set(key, bytes)
}

func (k Keeper) GetBlockFees(ctx sdk.Context, height uint64) (sdk.Coins, error) {
	key := sdk.Uint64ToBigEndian(uint64(ctx.BlockHeight()))
	bytes := k.feesStore(ctx).Get(key)
	if bytes == nil {
		return nil, fmt.Errorf("No fees were registered for block %v", height)
	}

	feesLength := sdk.BigEndianToUint64(bytes[0:8])
	//FIXME: this seems a bit hacky (using UVarInt is implementation specific) and maybe we could use a simple fee instead
	// Or find a better way to serialize coins
	var fees sdk.Coins
	last := uint64(8)
	for i := uint64(0); i < feesLength; i++ {
		var fee sdk.Coin
		bytesSize, lengthSize := binary.Uvarint(bytes[last:])
		totalSize := bytesSize + uint64(lengthSize)
		k.Cdc.MustUnmarshalBinaryLengthPrefixed(bytes[last:last+totalSize], &fee)
		last += totalSize
		fees = append(fees, fee)
	}

	return fees, nil
}

// Logger returns aliceAddr module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
