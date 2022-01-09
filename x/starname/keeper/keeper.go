package keeper

import (
	"fmt"
	"time"

	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"

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
	TokensFromConsensusPower(ctx sdk.Context, power int64) sdk.Int
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

// EscrowKeeper defines the behaviour of the escrow keeper, used to add stores to the module
// and register custom data for transfer handlers
type EscrowKeeper interface {
	RegisterCustomData(id escrowtypes.TypeID, data escrowtypes.CustomData)
}

// Keeper of the domain store
// TODO split this keeper in sub-struct in order to avoid possible mistakes with keys and not clutter the exposed methods
type Keeper struct {
	// external keepers
	ConfigurationKeeper ConfigurationKeeper
	SupplyKeeper        SupplyKeeper
	EscrowKeeper        EscrowKeeper
	AuthKeeper          AuthKeeper
	StakingKeeper       StakingKeeper
	DistributionKeeper  DistributionKeeper
	// default fields
	StoreKey   sdk.StoreKey // contains the store key for the domain module
	Cdc        codec.Codec
	paramspace ParamSubspace
	// Used for block fees queries
	cms sdk.CommitMultiStore
}

// NewKeeper creates a domain keeper
func NewKeeper(cdc codec.Codec, storeKey sdk.StoreKey, configKeeper ConfigurationKeeper, supply SupplyKeeper, escrow EscrowKeeper, auth AuthKeeper, distrib DistributionKeeper, staking StakingKeeper, paramspace ParamSubspace, cms sdk.CommitMultiStore) Keeper {
	keeper := Keeper{
		StoreKey:            storeKey,
		Cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		SupplyKeeper:        supply,
		EscrowKeeper:        escrow,
		AuthKeeper:          auth,
		DistributionKeeper:  distrib,
		StakingKeeper:       staking,
		paramspace:          paramspace,
		cms:                 cms,
	}
	keeper.ConfigureEscrowModule()
	return keeper
}

// ConfigureEscrowModule wraps this keepers reference in custom data, which is needed by the transfer handlers
func (k Keeper) ConfigureEscrowModule() {
	// Add custom data
	k.EscrowKeeper.RegisterCustomData(types.DomainTypeID, k)
	k.EscrowKeeper.RegisterCustomData(types.AccountTypeID, k)
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
		// FIXME: determine why slidingSum.feesSum can be negative: 3:41PM ERR CONSENSUS FAILURE!!! err="negative coin amount"
		if fees.IsAnyGT(slidingSum.feesSum) {
			fees = slidingSum.feesSum
		}
		slidingSum.feesSum = slidingSum.feesSum.Sub(fees)
	}
}

// GetBlockFeesSum retrieves the current value for the sum of the last n blocks
func (k Keeper) GetBlockFeesSum(ctx sdk.Context, maxBlocksInSum uint64) (sdk.Coins, uint64, error) {
	currentHeight := uint64(ctx.BlockHeight())

	// We force the query height to be greater than the latest computed height of the cached sliding sum otherwise querying
	// at different height would cause high latency as the sum would have to be recomputed.
	if currentHeight < slidingSum.lastComputedHeight {
		return nil, 0, fmt.Errorf("querying at past height is forbidden because of performance issues: queried " +
			"height %v when last known height is %v", currentHeight, slidingSum.lastComputedHeight)
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

	return slidingSum.feesSum, slidingSum.feesSumCount, nil
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

// DoAccountTransfer checks if the account transfer request is valid and then transfer the account and resets the
// associated data if toReset is true
func (k Keeper) DoAccountTransfer(
	ctx sdk.Context,
	name string,
	domain string,
	currentOwner sdk.AccAddress,
	newOwner sdk.AccAddress,
	toReset bool) (*types.Account, *types.Domain, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, domain).WithDomains(&domains)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, nil, err
	}

	// check if account exists
	accounts := k.AccountStore(ctx)
	accountCtrl := NewAccountController(ctx, domain, name).WithAccounts(&accounts).WithDomainController(domainCtrl)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		TransferableBy(currentOwner).
		ResettableBy(currentOwner, toReset).
		Validate(); err != nil {
		return nil, nil, err
	}

	// transfer account
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
	ex.Transfer(newOwner, toReset)
	return accountCtrl.account, domainCtrl.domain, nil
}

// DoDomainTransfer checks if the domain transfer request is valid and then transfer the domain according to the requested
// transferFlag
func (k Keeper) DoDomainTransfer(
	ctx sdk.Context,
	domain string,
	currentOwner sdk.AccAddress,
	newOwner sdk.AccAddress,
	transferFlag types.TransferFlag) error {
	// do precondition and authorization checks
	domains := k.DomainStore(ctx)
	c := NewDomainController(ctx, domain).WithDomains(&domains)
	err := c.
		MustExist().
		Admin(currentOwner).
		NotExpired().
		Transferable(transferFlag).
		Validate()
	if err != nil {
		return err
	}
	// transfer
	accounts := k.AccountStore(ctx)
	ex := NewDomainExecutor(ctx, c.Domain()).WithDomains(&domains).WithAccounts(&accounts)
	ex.Transfer(transferFlag, newOwner)
	return nil
}
