package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

// RegisterLegacyAminoCodec registers concrete types that will appear in
// interface fields/elements to be encoded/decoded by go-amino.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(MsgUpdateFees{}, fmt.Sprintf("%s/MsgUpdateFees", ModuleName), nil)
	cdc.RegisterConcrete(MsgUpdateConfig{}, fmt.Sprintf("%s/MsgUpdateConfig", ModuleName), nil)
}

// RegisterInterfaces registers implementations on registry.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgUpdateConfig{},
		&MsgUpdateFees{},
	)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/configuration module codec.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
