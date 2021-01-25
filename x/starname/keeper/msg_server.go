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
	if _, err := handlerMsgAddAccountCertificate(sdk.UnwrapSDKContext(goCtx), *m.keeper, msg); err != nil {
		return nil, err
	}
	return &types.MsgAddAccountCertificateResponse{}, nil
}
