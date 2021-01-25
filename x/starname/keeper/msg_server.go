package keeper

import (
	"context"
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/starname/controllers/account"
	"github.com/iov-one/starnamed/x/starname/controllers/domain"
	"github.com/iov-one/starnamed/x/starname/controllers/fees"
	"github.com/iov-one/starnamed/x/starname/keeper/executor"
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
	ctx := sdk.UnwrapSDKContext(goCtx)
	k := m.keeper

	// perform domain checks
	domainCtrl := domain.NewController(ctx, msg.Domain)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithDomainController(domainCtrl)

	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		CertificateLimitNotExceeded().
		CertificateSizeNotExceeded(msg.NewCertificate).
		CertificateNotExist(msg.NewCertificate).
		Validate(); err != nil {
		return nil, err
	}
	feeConf := m.keeper.ConfigurationKeeper.GetFees(ctx)
	feeCtrl := fees.NewController(ctx, feeConf, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}
	// add certificate
	ex := executor.NewAccount(ctx, accountCtrl.Account())
	ex.AddCertificate(msg.NewCertificate)
	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(sdk.AttributeKeySender, msg.Owner.String()),
			sdk.NewAttribute(sdk.AttributeKeyAction, msg.Type()),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyAccountName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyNewCertificate, fmt.Sprintf("%x", msg.NewCertificate)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)

	return &types.MsgAddAccountCertificateResponse{}, nil
}
