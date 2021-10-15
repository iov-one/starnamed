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

// RefreshBlockSumCache refresh the sliding sum value if it was
func (k Keeper) RefreshBlockSumCache(ctx sdk.Context, maxBlocksInSum uint64) {
	if slidingSum.feesSumCount != 0 {
		k.GetBlockFeesSum(ctx, maxBlocksInSum)
	}
}

func (k Keeper) addOrRemoveFeesSum(ctx sdk.Context, height uint64, add bool) {
	fees, err := k.GetBlockFees(ctx, height)
	if err != nil {
		panic(fmt.Sprintf("Cannot retrieve fees for the block at height %v", height))
	}
	if add {
		slidingSum.feesSum = slidingSum.feesSum.Add(fees...)
	} else {
		slidingSum.feesSum = slidingSum.feesSum.Sub(fees)
	}
}

// GetBlockFeesSum retrieves the current value for the sum of the last n blocks
func (k Keeper) GetBlockFeesSum(ctx sdk.Context, maxBlocksInSum uint64) (sdk.Coins, uint64) {
	//FIXME: the block height is not updated when querying at a different height (only the stores are)
	// So this line prevent to query from a different height (and will make the cms panic)
	// Querying at previous heights also cause problems for the cached sliding sum
	currentHeight := uint64(ctx.BlockHeight())

	// Solves a bug where currentHeight is less than the last computed height, we always want to
	// have the latest available data no matter what
	if currentHeight < slidingSum.lastComputedHeight {
		currentHeight = slidingSum.lastComputedHeight
	}

	if currentHeight < maxBlocksInSum {
		maxBlocksInSum = currentHeight
	}

	// if lastComputedHeight is too far behind we discard the current sliding sum and reset it
	if slidingSum.lastComputedHeight+maxBlocksInSum <= currentHeight {
		slidingSum.lastComputedHeight = currentHeight - maxBlocksInSum
		slidingSum.feesSum = sdk.Coins{}
		slidingSum.feesSumCount = 0
	}

	// if cached sliding sum has a smaller number of element, we may need to add the fees from block older than the older
	// block of the sliding sum
	// if the sliding sum contains all the elements wanted for this query, the following loop will not execute
	for h := slidingSum.lastComputedHeight - slidingSum.feesSumCount; h > currentHeight-maxBlocksInSum; h-- {
		k.addOrRemoveFeesSum(ctx, h, true)
		slidingSum.feesSumCount++
	}

	// Remove fees from blocks that are too old (if the value of maxBlocksInSum was higher in the previous call)
	for h := slidingSum.lastComputedHeight - slidingSum.feesSumCount + 1; h <= currentHeight-maxBlocksInSum; h++ {
		k.addOrRemoveFeesSum(ctx, h, false)
		slidingSum.feesSumCount--
	}

	// If we need to, we update the value of the sliding sum, starting with the first new block
	for ; slidingSum.lastComputedHeight < currentHeight; slidingSum.lastComputedHeight++ {

		// We store the new value in the sliding sum
		k.addOrRemoveFeesSum(ctx, slidingSum.lastComputedHeight+1, true)

		// If we reached the maximum number of block, we remove the last block from the sliding sum
		if slidingSum.feesSumCount == maxBlocksInSum {
			k.addOrRemoveFeesSum(ctx, slidingSum.lastComputedHeight-maxBlocksInSum+1, false)
		} else { // Else we just increment the number of block included in the sliding sum
			slidingSum.feesSumCount++
		}
	}

	return slidingSum.feesSum, slidingSum.feesSumCount
}

func (k Keeper) GetBlockFees(ctx sdk.Context, height uint64) (sdk.Coins, error) {

	cms, err := k.cms.CacheMultiStoreWithVersion(int64(height))
	if err != nil {
		return nil, err
	}
	ctxWithCurrentHeight := ctx.WithMultiStore(cms)

	fees := k.SupplyKeeper.GetAllBalances(ctxWithCurrentHeight, k.AuthKeeper.GetModuleAddress(authtypes.FeeCollectorName))

	return fees, nil
}

// Logger returns a module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
