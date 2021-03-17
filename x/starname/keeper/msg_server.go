package keeper

import (
	"context"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/x/starname/types"
)

var _ types.MsgServer = msgServer{}

type msgServer struct {
	keeper *Keeper
}

// NewMsgServerImpl returns a msgServer implementation
func NewMsgServerImpl(k *Keeper) types.MsgServer {
	return &msgServer{keeper: k}
}

func (m msgServer) AddAccountCertificate(goCtx context.Context, msg *types.MsgAddAccountCertificate) (*types.MsgAddAccountCertificateResponse, error) {
	return addAccountCertificate(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) DeleteAccount(goCtx context.Context, msg *types.MsgDeleteAccount) (*types.MsgDeleteAccountResponse, error) {
	return deleteAccount(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) DeleteAccountCertificate(goCtx context.Context, msg *types.MsgDeleteAccountCertificate) (*types.MsgDeleteAccountCertificateResponse, error) {
	return deleteAccountCertificate(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) DeleteDomain(goCtx context.Context, msg *types.MsgDeleteDomain) (*types.MsgDeleteDomainResponse, error) {
	return deleteDomain(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	return registerAccount(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) RegisterDomain(goCtx context.Context, msg *types.MsgRegisterDomain) (*types.MsgRegisterDomainResponse, error) {
	return registerDomain(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) RenewAccount(goCtx context.Context, msg *types.MsgRenewAccount) (*types.MsgRenewAccountResponse, error) {
	return renewAccount(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) RenewDomain(goCtx context.Context, msg *types.MsgRenewDomain) (*types.MsgRenewDomainResponse, error) {
	return renewDomain(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}

func (m msgServer) ReplaceAccountMetadata(goCtx context.Context, msg *types.MsgReplaceAccountMetadata) (*types.MsgReplaceAccountMetadataResponse, error) {
	return replaceAccountMetadata(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}

func (m msgServer) ReplaceAccountResources(goCtx context.Context, msg *types.MsgReplaceAccountResources) (*types.MsgReplaceAccountResourcesResponse, error) {
	return replaceAccountResources(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}

func (m msgServer) TransferAccount(goCtx context.Context, msg *types.MsgTransferAccount) (*types.MsgTransferAccountResponse, error) {
	return transferAccount(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}

func (m msgServer) TransferDomain(goCtx context.Context, msg *types.MsgTransferDomain) (*types.MsgTransferDomainResponse, error) {
	return transferDomain(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg.ToInternal())
}
