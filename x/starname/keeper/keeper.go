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
	//TODO: find a way to persist this
	feesSum            sdk.Coins
	feesSumCount       uint64
	lastComputedHeight uint64
	// TODO: this should maybe moved out of this keeper
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

// GetBlockFeesSum retrieves the current value for the sum of the last n blocks
func (k Keeper) GetBlockFeesSum(ctx sdk.Context) (sdk.Coins, uint64) {
	// TODO: this shouldn't be an hardcoded value
	const MaxBlocksInSum = 50000

	//TODO: should we offset this to take into account that fee collectors funds are transferred on block start ?
	currentHeight := ctx.BlockHeight()

	//fmt.Printf("Begin block fees sum before %v until %v\n", k.lastComputedHeight, currentHeight)

	// if lastComputedHeight is too far behind we discard the current sliding sum and reset it
	if k.lastComputedHeight+MaxBlocksInSum < uint64(currentHeight) {
		k.lastComputedHeight = uint64(currentHeight) - MaxBlocksInSum
		k.feesSum = sdk.Coins{}
		k.feesSumCount = 0
	}

	// If we need to, we update the value of the sliding sum
	for ; k.lastComputedHeight < uint64(currentHeight); k.lastComputedHeight++ {

		// We store the new value in the sliding sum
		fees, err := k.GetBlockFees(ctx, uint64(k.lastComputedHeight))
		if err != nil {
			panic(fmt.Sprintf("Cannot retrieve fees for the block at height %v", currentHeight))
		}
		k.feesSum = k.feesSum.Add(fees...)

		// If we reached the maximum number of block, we remove the last block from the sliding sum
		if k.feesSumCount == MaxBlocksInSum {
			feesOutgoingBlock, err := k.GetBlockFees(ctx, k.lastComputedHeight-MaxBlocksInSum)
			if err != nil {
				panic(fmt.Sprintf("Cannot retrieve fees for the block at height %v", currentHeight))
			}
			k.feesSum = k.feesSum.Sub(feesOutgoingBlock)
		} else { // Else we just increment the number of block included in the sliding sum
			k.feesSumCount++
		}
	}

	return k.feesSum, k.feesSumCount
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

// Logger returns aliceAddr module-specific logger.
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("x/%s", types.ModuleName))
}
