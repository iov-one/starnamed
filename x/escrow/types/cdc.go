package types

import (
	"fmt"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the message for the legacy amino codec, used in the legacy REST handlers
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgCreateEscrow{}, fmt.Sprintf("%s/CreateEscrow", ModuleName), nil)
	cdc.RegisterConcrete(&MsgUpdateEscrow{}, fmt.Sprintf("%s/UpdateEscrow", ModuleName), nil)
	cdc.RegisterConcrete(&MsgTransferToEscrow{}, fmt.Sprintf("%s/TransferToEscrow", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRefundEscrow{}, fmt.Sprintf("%s/RefundEscrow", ModuleName), nil)

	cdc.RegisterInterface((*TransferableObject)(nil), nil)
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
	// ModuleCdc references the global x/escrow module codec.
	ModuleCdc = codec.NewAminoCodec(amino)
)
