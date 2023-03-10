package burner

import (
	"encoding/json"
	"fmt"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/types/module"
	"github.com/gorilla/mux"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/spf13/cobra"
	abci "github.com/tendermint/tendermint/abci/types"

	"github.com/iov-one/starnamed/x/burner/types"
)

var (
	_ module.AppModule      = AppModule{}
	_ module.AppModuleBasic = AppModuleBasic{}
)

// AppModuleBasic defines the basic application module used by the burner module.
type AppModuleBasic struct {
	cdc codec.Codec
}

// RegisterLegacyAminoCodec registers the amino codec.
func (b AppModuleBasic) RegisterLegacyAminoCodec(*codec.LegacyAmino) {
}

// RegisterGRPCGatewayRoutes registers the query handler client.
func (b AppModuleBasic) RegisterGRPCGatewayRoutes(client.Context, *runtime.ServeMux) {
}

// Name returns the burner module's name.
func (AppModuleBasic) Name() string { return types.ModuleName }

// DefaultGenesis returns default genesis state as raw bytes for the burner module.
func (AppModuleBasic) DefaultGenesis(codec.JSONCodec) json.RawMessage {
	return nil
}

// ValidateGenesis performs genesis state validation for the burner module.
func (b AppModuleBasic) ValidateGenesis(_ codec.JSONCodec, _ client.TxEncodingConfig, genesisData json.RawMessage) error {
	if len(genesisData) > 0 {
		return fmt.Errorf("invalid genesis data for module burner: should be empty")
	}
	return nil
}

// RegisterRESTRoutes registers the REST routes for this module.
func (AppModuleBasic) RegisterRESTRoutes(client.Context, *mux.Router) {
}

// GetQueryCmd returns no root query command for this module.
func (AppModuleBasic) GetQueryCmd() *cobra.Command {
	return nil
}

// GetTxCmd returns the root tx command for this module.
func (AppModuleBasic) GetTxCmd() *cobra.Command {
	return nil
}

// RegisterInterfaces implements InterfaceModule
func (b AppModuleBasic) RegisterInterfaces(cdctypes.InterfaceRegistry) {
}

// AppModule implements an application module for the burner module.
type AppModule struct {
	AppModuleBasic
	supplyKeeper  types.SupplyKeeper
	accountKeeper types.AccountKeeper
}

// NewAppModule creates a new AppModule object.
func NewAppModule(supplyKeeper types.SupplyKeeper, accountKeeper types.AccountKeeper) AppModule {
	return AppModule{
		AppModuleBasic: AppModuleBasic{},
		supplyKeeper:   supplyKeeper,
		accountKeeper:  accountKeeper,
	}
}

// Name returns the burner module's name.
func (am AppModule) Name() string { return am.AppModuleBasic.Name() }

// RegisterServices allows a module to register services
func (am AppModule) RegisterServices(module.Configurator) {
}

// LegacyQuerierHandler provides an sdk.Querier object that uses the legacy amino codec.
func (AppModule) LegacyQuerierHandler(*codec.LegacyAmino) sdk.Querier {
	return func(ctx sdk.Context, path []string, req abci.RequestQuery) ([]byte, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrUnknownRequest, "%s", path[0])
	}
}

// RegisterInvariants registers the burner module invariants.
func (AppModule) RegisterInvariants(sdk.InvariantRegistry) {}

// Route returns the message routing key for the burner module.
func (am AppModule) Route() sdk.Route {
	return sdk.NewRoute(types.ModuleName, func(sdk.Context, sdk.Msg) (*sdk.Result, error) {
		return nil, sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "unknown request")
	})
}

// QuerierRoute returns the burner module's querier route name.
func (AppModule) QuerierRoute() string { return types.ModuleName }

// BeginBlock returns the begin blocker for the burner module.
func (AppModule) BeginBlock(sdk.Context, abci.RequestBeginBlock) {}

// EndBlock returns the end blocker for the burner module. It returns no validator updates.
func (am AppModule) EndBlock(ctx sdk.Context, _ abci.RequestEndBlock) []abci.ValidatorUpdate {
	EndBlocker(ctx, am.supplyKeeper, am.accountKeeper)
	return []abci.ValidatorUpdate{}
}

// InitGenesis performs genesis initialization for the burner module. It returns no validator updates.
func (am AppModule) InitGenesis(ctx sdk.Context, _ codec.JSONCodec, _ json.RawMessage) []abci.ValidatorUpdate {
	return []abci.ValidatorUpdate{}
}

// ExportGenesis returns the exported genesis state as raw bytes for the burner module.
func (am AppModule) ExportGenesis(sdk.Context, codec.JSONCodec) json.RawMessage {
	return am.DefaultGenesis(nil)
}

// ConsensusVersion implements AppModule/ConsensusVersion.
func (AppModule) ConsensusVersion() uint64 { return 1 }
