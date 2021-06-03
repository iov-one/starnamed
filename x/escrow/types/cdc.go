package types

import (
	fmt "fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the account types and interface.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateEscrow{}, fmt.Sprintf("%s/CreateEscrow", ModuleName), nil)
	cdc.RegisterConcrete(&MsgUpdateEscrow{}, fmt.Sprintf("%s/UpdateEscrow", ModuleName), nil)
	cdc.RegisterConcrete(&MsgTransferToEscrow{}, fmt.Sprintf("%s/TransferToEscrow", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRefundEscrow{}, fmt.Sprintf("%s/RefundEscrow", ModuleName), nil)
}

// RegisterInterfaces registers implementations for the protobuf marshaler.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgCreateEscrow{},
		&MsgUpdateEscrow{},
		&MsgTransferToEscrow{},
		&MsgRefundEscrow{},
	)
	msgservice.RegisterMsgServiceDesc(registry, &_Msg_serviceDesc)
}

var (
	amino = codec.NewLegacyAmino()

	// ModuleCdc references the global x/starname module codec.
	ModuleCdc = codec.NewAminoCodec(amino)
)

func init() {
	RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
}
