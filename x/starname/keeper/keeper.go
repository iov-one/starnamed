package keeper

import (
	"fmt"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	crudtypes "github.com/iov-one/cosmos-sdk-crud/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
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

// DistributionKeeper is used to estimate the yield of the chain
type DistributionKeeper interface {
	GetCommunityTax(ctx sdk.Context) sdk.Dec
}

// StakingKeeper is used to estimate the yield of the chain
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
	// GetValidDomainNameRegexp returns the regular expression that a domain name must match
	// in order to be valid
	GetValidDomainNameRegexp(ctx sdk.Context) string
	// GetDomainRenewDuration returns the default duration of a domain renewal
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
	// Used for block fees queries
	cms sdk.CommitMultiStore
}

// NewKeeper creates a domain keeper
func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, configKeeper ConfigurationKeeper, supply SupplyKeeper, auth AuthKeeper, distrib DistributionKeeper, staking StakingKeeper, paramspace ParamSubspace, cms sdk.CommitMultiStore) Keeper {
	keeper := Keeper{
		StoreKey:            storeKey,
		Cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		SupplyKeeper:        supply,
		AuthKeeper:          auth,
		DistributionKeeper:  distrib,
		StakingKeeper:       staking,
		paramspace:          paramspace,
		cms:                 cms,
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

// TODO: we should maybe move this in a separate module

//TODO: this cannot be persisted in the keeper as a new one is used for each query
// Find anoter way of persisting this data, without a global variable
var slidingSum struct {
	feesSum            sdk.Coins
	feesSumCount       uint64
	lastComputedHeight uint64
}

// RefreshBlockSumCache refreshes the sliding sum value if it has been previously computed
func (k Keeper) RefreshBlockSumCache(ctx sdk.Context) {
	if slidingSum.feesSumCount != 0 {
		// Ignore errors
		_, _ = k.GetLastWeekBlockFeesSum(ctx)
	}
}

func (k Keeper) addToFeeSum(ctx sdk.Context, height uint64) {
	fees, err := k.GetBlockFees(ctx, height)
	if err != nil {
		panic(fmt.Sprintf("Cannot retrieve fees for the block at height %v", height))
	}
	slidingSum.feesSum = slidingSum.feesSum.Add(fees...)
	slidingSum.feesSumCount++
}

func (k Keeper) removeFromFeeSum(ctx sdk.Context, height uint64) {
	fees, err := k.GetBlockFees(ctx, height)
	if err != nil {
		panic(fmt.Sprintf("Cannot retrieve fees for the block at height %v", height))
	}
	slidingSum.feesSum = slidingSum.feesSum.Sub(fees)
	slidingSum.feesSumCount--
}

// GetLastWeekBlockFeesSum retrieves the current value for the sum of the blocks of last weeks (100k blocks)
func (k Keeper) GetLastWeekBlockFeesSum(ctx sdk.Context) (sdk.Coins, error) {
	//FIXME: the block height is not updated when querying at a different height (only the stores are)
	// So this line prevent to query from a different height (and will make the cms panic)
	// Querying at previous heights also cause problems for the cached sliding sum
	currentHeight := uint64(ctx.BlockHeight())

	// Solves a bug where currentHeight is less than the last computed height, we always want to
	// have the latest available data no matter what
	if currentHeight < slidingSum.lastComputedHeight {
		currentHeight = slidingSum.lastComputedHeight
	}

	// If we don't have enough blocks, we cannot compute the sum and return an error
	if currentHeight < NumBlocksInAWeek {
		return nil, fmt.Errorf("not enough data to estimate yield: current height %v is smaller than %v",
			currentHeight, NumBlocksInAWeek)
	}

	// if lastComputedHeight is too far behind we discard the current sliding sum and reset it
	if slidingSum.lastComputedHeight+NumBlocksInAWeek <= currentHeight {
		slidingSum.lastComputedHeight = currentHeight - NumBlocksInAWeek
		slidingSum.feesSum = sdk.Coins{}
		slidingSum.feesSumCount = 0
	}

	// Remove fees we don't want anymore
	oldestWantedBlock := currentHeight - NumBlocksInAWeek + 1
	oldestCachedBlock := slidingSum.lastComputedHeight - slidingSum.feesSumCount + 1
	for h := oldestCachedBlock; h < oldestWantedBlock; h++ {
		k.removeFromFeeSum(ctx, h)
	}

	// Add fees from news blocks
	newestWantedBlock := currentHeight
	newestCachedBlock := slidingSum.lastComputedHeight
	for h := newestCachedBlock + 1; h <= newestWantedBlock; h++ {
		k.addToFeeSum(ctx, h)
	}
	slidingSum.lastComputedHeight = currentHeight

	// Make sure we have not messed up with the number of elements
	if slidingSum.feesSumCount != NumBlocksInAWeek {
		panic(fmt.Errorf("inconsistent fees sum: %v blocks included when expecting %v", slidingSum.feesSumCount, NumBlocksInAWeek))
	}

	return slidingSum.feesSum, nil
}

// GetBlockFees returns the fees collected at a specific height
// It will return an error if the node has not an history of the given height
func (k Keeper) GetBlockFees(ctx sdk.Context, height uint64) (sdk.Coins, error) {
	ctxWithCurrentHeight := ctx
	if int64(height) != ctxWithCurrentHeight.BlockHeight() {
		cms, err := k.cms.CacheMultiStoreWithVersion(int64(height))
		if err != nil {
			return nil, err
		}
		ctxWithCurrentHeight = ctxWithCurrentHeight.WithMultiStore(cms)
	}

	fees := k.SupplyKeeper.GetAllBalances(ctxWithCurrentHeight, k.AuthKeeper.GetModuleAddress(authtypes.FeeCollectorName))

	return fees, nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
