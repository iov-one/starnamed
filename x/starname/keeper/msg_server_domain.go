package keeper

import (
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

	// success TODO maybe emit event?
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

	// success TODO think here, can we emit any useful event
	return &types.MsgRegisterDomainResponse{}, nil
}

// renewDomain renews a domain
func renewDomain(ctx sdk.Context, k Keeper, msg *types.MsgRenewDomain) (*types.MsgRenewDomainResponse, error) {
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

	// success TODO emit event
	return &types.MsgRenewDomainResponse{}, nil
}

func transferDomain(ctx sdk.Context, k Keeper, msg *types.MsgTransferDomain) (*types.MsgTransferDomainResponse, error) {
	// do precondition and authorization checks
	domains := k.DomainStore(ctx)
	c := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	err := c.
		MustExist().
		Admin(msg.Owner).
		NotExpired().
		Transferable(msg.TransferFlag).
		Validate()
	if err != nil {
		return nil, err
	}

	// collect fees
	domain := new(types.Domain)
	if err := domains.Read([]byte(msg.Domain), domain); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}
	if err := k.CollectProductFee(ctx, msg, domain); err != nil {
		return nil, sdkerrors.Wrapf(err, "unable to collect fees")
	}

	// transfer
	accounts := k.AccountStore(ctx)
	ex := NewDomainExecutor(ctx, c.Domain()).WithDomains(&domains).WithAccounts(&accounts)
	ex.Transfer(msg.TransferFlag, msg.NewAdmin)

	// success; TODO emit event?
	return &types.MsgTransferDomainResponse{}, nil
}
