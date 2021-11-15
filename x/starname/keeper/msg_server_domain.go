package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/iov-one/starnamed/x/starname/types"
)

// deleteDomain deletes a domain from the store
func deleteDomain(ctx sdk.Context, k Keeper, msg *types.MsgDeleteDomainInternal) (*types.MsgDeleteDomainResponse, error) {
	// do precondition and authorization checks
	domains := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	ctrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains).WithConfiguration(conf)
	if err := ctrl.
		MustExist().
		DeletableBy(msg.Owner).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}

	// all checks passed delete domain
	accounts := k.AccountStore(ctx)
	NewDomainExecutor(ctx, ctrl.Domain()).WithDomains(&domains).WithAccounts(&accounts).Delete()

	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &types.MsgDeleteDomainResponse{}, nil
}

// registerDomain handles the domain registration process
func registerDomain(ctx sdk.Context, k Keeper, msg *types.MsgRegisterDomainInternal) (*types.MsgRegisterDomainResponse, error) {
	// do precondition and authorization checks
	domains := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	ctrl := NewDomainController(ctx, msg.Name).WithDomains(&domains).WithConfiguration(conf)
	err := ctrl.
		MustNotExist().
		ValidName().
		Validate()
	if err != nil {
		return nil, err
	}

	// create new domain
	d := types.Domain{
		Name:       msg.Name,
		Admin:      msg.Admin,
		ValidUntil: ctx.BlockTime().Add(k.ConfigurationKeeper.GetDomainRenewDuration(ctx)).Unix(),
		Type:       msg.DomainType,
		Broker:     msg.Broker,
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg, &d); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}

	// save domain
	accounts := k.AccountStore(ctx)
	ex := NewDomainExecutor(ctx, d).WithDomains(&domains).WithAccounts(&accounts)
	ex.Create()

	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Name),
			sdk.NewAttribute(types.AttributeKeyDomainType, (string)(msg.DomainType)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Admin.String()),
			sdk.NewAttribute(types.AttributeKeyBroker, msg.Broker.String()),
		),
	)
	return &types.MsgRegisterDomainResponse{}, nil
}

// renewDomain renews a domain
func renewDomain(ctx sdk.Context, k Keeper, msg *types.MsgRenewDomainInternal) (*types.MsgRenewDomainResponse, error) {
	// do precondition and authorization checks
	domains := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	ctrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains).WithConfiguration(conf)
	err := ctrl.
		MustExist().
		Renewable().
		Validate()
	if err != nil {
		return nil, err
	}

	// collect fees
	domain := new(types.Domain)
	if err := domains.Read([]byte(msg.Domain), domain); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}
	if err := k.CollectProductFee(ctx, msg, domain, k.AccountStore); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}

	// update domain
	accounts := k.AccountStore(ctx)
	NewDomainExecutor(ctx, ctrl.Domain()).WithDomains(&domains).WithAccounts(&accounts).WithConfiguration(conf).Renew()

	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Signer.String()),
		),
	)
	return &types.MsgRenewDomainResponse{}, nil
}

func transferDomain(ctx sdk.Context, k Keeper, msg *types.MsgTransferDomainInternal) (*types.MsgTransferDomainResponse, error) {
	// do checks and domain transfer
	if err := k.DoDomainTransfer(ctx, msg.Domain, msg.Owner, msg.NewAdmin, msg.TransferFlag); err != nil {
		return nil, err
	}

	// collect fees
	domain := new(types.Domain)
	if err := k.DomainStore(ctx).Read([]byte(msg.Domain), domain); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}
	if err := k.CollectProductFee(ctx, msg, domain); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}

	// success
	ctx.EventManager().EmitEvent(
		sdk.NewEvent(
			sdk.EventTypeMessage,
			sdk.NewAttribute(sdk.AttributeKeyModule, types.ModuleName),
			sdk.NewAttribute(types.AttributeKeyDomainName, msg.Domain),
			sdk.NewAttribute(types.AttributeKeyTransferDomainNewOwner, msg.NewAdmin.String()),
			sdk.NewAttribute(types.AttributeKeyTransferDomainFlag, fmt.Sprintf("%d", msg.TransferFlag)),
			sdk.NewAttribute(types.AttributeKeyOwner, msg.Owner.String()),
		),
	)
	return &types.MsgTransferDomainResponse{}, nil
}
