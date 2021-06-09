package types

import (
	"fmt"

	escrowtypes "github.com/iov-one/starnamed/x/escrow/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/codec/types"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/msgservice"
)

// RegisterLegacyAminoCodec registers the account types and interface.
func RegisterLegacyAminoCodec(cdc *codec.LegacyAmino) {
	cdc.RegisterConcrete(&MsgRegisterDomain{}, fmt.Sprintf("%s/RegisterDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgTransferDomain{}, fmt.Sprintf("%s/TransferDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgTransferAccount{}, fmt.Sprintf("%s/TransferAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRenewAccount{}, fmt.Sprintf("%s/RenewAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgAddAccountCertificate{}, fmt.Sprintf("%s/AddAccountCertificate", ModuleName), nil)
	cdc.RegisterConcrete(&MsgDeleteAccountCertificate{}, fmt.Sprintf("%s/DeleteAccountCertificate", ModuleName), nil)
	cdc.RegisterConcrete(&MsgDeleteAccount{}, fmt.Sprintf("%s/DeleteAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgDeleteDomain{}, fmt.Sprintf("%s/DeleteDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRegisterAccount{}, fmt.Sprintf("%s/RegisterAccount", ModuleName), nil)
	cdc.RegisterConcrete(&MsgRenewDomain{}, fmt.Sprintf("%s/RenewDomain", ModuleName), nil)
	cdc.RegisterConcrete(&MsgReplaceAccountResources{}, fmt.Sprintf("%s/ReplaceAccountResources", ModuleName), nil)
	cdc.RegisterConcrete(&MsgReplaceAccountMetadata{}, fmt.Sprintf("%s/SetAccountMetadata", ModuleName), nil)
}

// RegisterInterfaces registers implementations for the protobuf marshaler.
func RegisterInterfaces(registry types.InterfaceRegistry) {
	registry.RegisterImplementations(
		(*sdk.Msg)(nil),
		&MsgAddAccountCertificate{},
		&MsgDeleteAccount{},
		&MsgDeleteAccountCertificate{},
		&MsgDeleteDomain{},
		&MsgRegisterAccount{},
		&MsgRegisterDomain{},
		&MsgRenewAccount{},
		&MsgRenewDomain{},
		&MsgReplaceAccountMetadata{},
		&MsgReplaceAccountResources{},
		&MsgTransferAccount{},
		&MsgTransferDomain{},
	)
	registry.RegisterImplementations(
		(*escrowtypes.TransferableObject)(nil),
		// Register the account object as a TransferableObject implementation to send it in a MsgCreateEscrow
		// TODO: is this the correct way of doing this ?
		&Account{},
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
