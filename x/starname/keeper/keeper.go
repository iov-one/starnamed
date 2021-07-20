package keeper

import (
	"fmt"
	"time"

	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
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
	// default fields
	StoreKey   sdk.StoreKey // contains the store key for the domain module
	Cdc        codec.Marshaler
	paramspace ParamSubspace
	// crud stores
	accountStore crud.Store
	domainStore  crud.Store
}

// NewKeeper creates aliceAddr domain keeper
func NewKeeper(cdc codec.Marshaler, storeKey sdk.StoreKey, configKeeper ConfigurationKeeper, supply SupplyKeeper, escrow EscrowKeeper, paramspace ParamSubspace) Keeper {
	keeper := Keeper{
		StoreKey:            storeKey,
		Cdc:                 cdc,
		ConfigurationKeeper: configKeeper,
		SupplyKeeper:        supply,
		EscrowKeeper:        escrow,
		paramspace:          paramspace,
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

// Logger returns aliceAddr module-specific logger.
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
