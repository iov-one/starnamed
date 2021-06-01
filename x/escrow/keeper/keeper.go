package keeper

import (
	"fmt"

	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/iov-one/starnamed/x/escrow/types"
	"github.com/tendermint/tendermint/libs/log"

	"github.com/cosmos/cosmos-sdk/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	paramstypes "github.com/cosmos/cosmos-sdk/x/params/types"
)

// Keeper defines the escrow keeper
type Keeper struct {
	storeKey      sdk.StoreKey
	cdc           codec.Marshaler
	paramSpace    paramstypes.Subspace
	accountKeeper types.AccountKeeper
	bankKeeper    types.BankKeeper
	storeHolders  map[types.TypeID]types.StoreHolder
}

// NewKeeper creates a new escrow Keeper instance
func NewKeeper(
	cdc codec.Marshaler,
	key sdk.StoreKey,
	paramSpace paramstypes.Subspace,
	accountKeeper types.AccountKeeper,
	bankKeeper types.BankKeeper,
	storeHolders map[types.TypeID]types.StoreHolder,
) Keeper {
	// ensure the escrow module account is set
	if addr := accountKeeper.GetModuleAddress(types.ModuleName); addr == nil {
		panic(fmt.Sprintf("%s module account has not been set", types.ModuleName))
	}

	return Keeper{
		storeKey:      key,
		cdc:           cdc,
		paramSpace:    paramSpace,
		accountKeeper: accountKeeper,
		bankKeeper:    bankKeeper,
		storeHolders:  storeHolders,
	}
}

// Logger returns a module-specific logger
func (k Keeper) Logger(ctx sdk.Context) log.Logger {
	return ctx.Logger().With("module", fmt.Sprintf("irismod/%s", types.ModuleName))
}

// GetEscrowAccount returns the escrow module account
func (k Keeper) GetEscrowAccount(ctx sdk.Context) authtypes.ModuleAccountI {
	return k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
}

func (k Keeper) GetNextId() tmbytes.HexBytes {
	//TODO: not implemented
	return nil
}

// TODO: check the utility of this
// EnsureModuleAccountPermissions syncs the bep3 module account's permissions with those in the supply keeper.
func (k Keeper) EnsureModuleAccountPermissions(ctx sdk.Context) error {
	maccI := k.accountKeeper.GetModuleAccount(ctx, types.ModuleName)
	macc, ok := maccI.(*authtypes.ModuleAccount)
	if !ok {
		return fmt.Errorf("expected %s account to be a module account type", types.ModuleName)
	}
	_, perms := k.accountKeeper.GetModuleAddressAndPermissions(types.ModuleName)
	macc.Permissions = perms
	k.accountKeeper.SetModuleAccount(ctx, macc)
	return nil
}
