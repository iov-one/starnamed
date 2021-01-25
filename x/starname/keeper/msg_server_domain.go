package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/starname/controllers/fees"
	"github.com/iov-one/starnamed/x/starname/keeper/executor"
	"github.com/iov-one/starnamed/x/starname/types"
)

func handlerMsgDeleteDomain(ctx sdk.Context, k Keeper, msg *types.MsgDeleteDomain) (*sdk.Result, error) {
	domains := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	ctrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains).WithConfiguration(conf)
	// do precondition and authorization checks
	if err := ctrl.
		MustExist().
		DeletableBy(msg.Owner).
		Validate(); err != nil {
		return nil, err
	}
	// operation is allowed
	feeConf := k.ConfigurationKeeper.GetFees(ctx)
	feeCtrl := fees.NewController(ctx, feeConf, ctrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// all checks passed delete domain
	accounts := k.AccountStore(ctx)
	executor.NewDomain(ctx, ctrl.Domain()).WithDomains(&domains).WithAccounts(&accounts).Delete()
	// success TODO maybe emit event?
	return &sdk.Result{}, nil
}

// handleMsgRegisterDomain handles the domain registration process
func handleMsgRegisterDomain(ctx sdk.Context, k Keeper, msg *types.MsgRegisterDomain) (resp *sdk.Result, err error) {
	domains := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	ctrl := NewDomainController(ctx, msg.Name).WithDomains(&domains).WithConfiguration(conf)
	err = ctrl.
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
	feeConf := k.ConfigurationKeeper.GetFees(ctx)
	feeCtrl := fees.NewController(ctx, feeConf, d)
	fee := feeCtrl.GetFee(msg)
	// collect fees
	if err := k.CollectFees(ctx, msg, fee); err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// save domain
	accounts := k.AccountStore(ctx)
	ex := executor.NewDomain(ctx, d).WithDomains(&domains).WithAccounts(&accounts)
	ex.Create()
	// success TODO think here, can we emit any useful event
	return &sdk.Result{}, nil
}

// handlerMsgRenewDomain renews a domain
func handlerMsgRenewDomain(ctx sdk.Context, k Keeper, msg *types.MsgRenewDomain) (*sdk.Result, error) {
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
	feeConf := k.ConfigurationKeeper.GetFees(ctx)
	accounts := k.AccountStore(ctx)
	feeCtrl := fees.NewController(ctx, feeConf, ctrl.Domain()).WithAccounts(&accounts)
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err = k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	// update domain
	executor.NewDomain(ctx, ctrl.Domain()).WithDomains(&domains).WithAccounts(&accounts).WithConfiguration(conf).Renew()
	// success TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgTransferDomain(ctx sdk.Context, k Keeper, msg *types.MsgTransferDomain) (*sdk.Result, error) {
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
	feeConf := k.ConfigurationKeeper.GetFees(ctx)
	feeCtrl := fees.NewController(ctx, feeConf, c.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err = k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, sdkerrors.Wrap(err, "unable to collect fees")
	}
	accounts := k.AccountStore(ctx)
	ex := executor.NewDomain(ctx, c.Domain()).WithDomains(&domains).WithAccounts(&accounts)
	ex.Transfer(msg.TransferFlag, msg.NewAdmin)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}
