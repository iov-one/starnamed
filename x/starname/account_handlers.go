package starname

import (
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/controllers/fees"
	"github.com/iov-one/starnamed/x/starname/keeper/executor"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/x/starname/controllers/account"
	"github.com/iov-one/starnamed/x/starname/controllers/domain"
	"github.com/iov-one/starnamed/x/starname/keeper"
	"github.com/iov-one/starnamed/x/starname/types"
)

func handlerMsgAddAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg *types.MsgAddAccountCertificate) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	as := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl).WithConfiguration(conf)

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
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}
	// add certificate
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.AddCertificate(msg.NewCertificate)
	// success; TODO emit event
	return &sdk.Result{}, nil
}

func handlerMsgDeleteAccountCertificate(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccountCertificate) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}
	// perform account checks, save certificate index
	as := k.AccountStore(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl)
	certIndex := new(int)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		CertificateExists(msg.DeleteCertificate, certIndex).
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// delete cert
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.DeleteCertificate(*certIndex)
	// success; TODO emit event?
	return &sdk.Result{}, nil
}

// handlerMsgDelete account deletes the account from the system
func handlerMsgDeleteAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgDeleteAccount) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.MustExist().Validate(); err != nil {
		return nil, err
	}
	// perform account checks
	as := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		DeletableBy(msg.Owner).
		Validate(); err != nil {
		return nil, err
	}
	// collect fees
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// delete account
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.Delete()
	// success; todo can we emit event?
	return &sdk.Result{}, nil
}

// handleMsgRegisterAccount registers the account
func handleMsgRegisterAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRegisterAccount) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}
	as := k.AccountStore(ctx)
	d := domainCtrl.Domain()
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		ValidName().
		MustNotExist().
		ValidResources(msg.Resources).
		RegistrableBy(msg.Registerer).
		Validate(); err != nil {
		return nil, err
	}

	a := types.Account{
		Domain:       msg.Domain,
		Name:         utils.StrPtr(msg.Name),
		Owner:        msg.Owner,
		Resources:    msg.Resources,
		Certificates: nil,
		Broker:       msg.Broker,
	}
	switch d.Type {
	case types.ClosedDomain:
		a.ValidUntil = types.MaxValidUntil
	case types.OpenDomain:
		a.ValidUntil = ctx.BlockTime().Add(conf.AccountRenewalPeriod).Unix()
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	ex := executor.NewAccount(ctx, k, a)
	ex.Create()
	return &sdk.Result{}, nil
}

func handlerMsgRenewAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgRenewAccount) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds).WithConfiguration(conf)
	if err := domainCtrl.MustExist().Type(types.OpenDomain).Validate(); err != nil {
		return nil, err
	}
	// validate account
	as := k.AccountStore(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		Renewable().
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// renew account
	// account valid until is extended here
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.Renew()
	// get grace period and expiration time
	d := domainCtrl.Domain()
	dgp := conf.DomainGracePeriod
	domainGracePeriodUntil := utils.SecondsToTime(d.ValidUntil).Add(dgp)
	accNewValidUntil := utils.SecondsToTime(ex.State().ValidUntil)
	if domainGracePeriodUntil.Before(accNewValidUntil) {
		dex := executor.NewDomain(ctx, k, domainCtrl.Domain())
		dex.Renew(accNewValidUntil.Unix())
	}
	// success; todo emit event??
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountResources replaces account resources
func handlerMsgReplaceAccountResources(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountResources) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}
	// perform account checks
	as := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		ValidResources(msg.NewResources).
		ResourceLimitNotExceeded(msg.NewResources).
		Validate(); err != nil {
		return nil, err
	}
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// replace accounts resources
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.ReplaceResources(msg.NewResources)
	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// handlerMsgReplaceAccountMetadata takes care of setting account metadata
func handlerMsgReplaceAccountMetadata(ctx sdk.Context, k keeper.Keeper, msg *types.MsgReplaceAccountMetadata) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}
	// perform account checks
	as := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		MetadataSizeNotExceeded(msg.NewMetadataURI).
		Validate(); err != nil {
		return nil, err
	}
	// collect fees
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// save to store
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.UpdateMetadata(msg.NewMetadataURI)
	// success TODO emit event
	return &sdk.Result{}, nil
}

// handlerMsgTransferAccount transfers account to a new owner
// after clearing resources and certificates
func handlerMsgTransferAccount(ctx sdk.Context, k keeper.Keeper, msg *types.MsgTransferAccount) (*sdk.Result, error) {
	// perform domain checks
	ds := k.DomainStore(ctx)
	domainCtrl := domain.NewController(ctx, msg.Domain).WithStore(&ds)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}
	// check if account exists
	as := k.AccountStore(ctx)
	accountCtrl := account.NewController(ctx, msg.Domain, msg.Name).WithStore(&as).WithDomainController(domainCtrl)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		TransferableBy(msg.Owner).
		ResettableBy(msg.Owner, msg.ToReset).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	feeCtrl := fees.NewController(ctx, k, domainCtrl.Domain())
	fee := feeCtrl.GetFee(msg)
	// collect fees
	err := k.CollectFees(ctx, msg, fee)
	if err != nil {
		return nil, errors.Wrap(err, "unable to collect fees")
	}
	// transfer account
	ex := executor.NewAccount(ctx, k, accountCtrl.Account())
	ex.Transfer(msg.NewOwner, msg.ToReset)
	// success, todo emit event?
	return &sdk.Result{}, nil
}
