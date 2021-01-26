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
	return addAccountCertificate(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}

func (m msgServer) DeleteDomain(goCtx context.Context, msg *types.MsgDeleteDomain) (*types.MsgDeleteDomainResponse, error) {
	return deleteDomain(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}

func (m msgServer) RegisterAccount(goCtx context.Context, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	return registerAccount(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}

func (m msgServer) RegisterDomain(goCtx context.Context, msg *types.MsgRegisterDomain) (*types.MsgRegisterDomainResponse, error) {
	return registerDomain(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg)
}
