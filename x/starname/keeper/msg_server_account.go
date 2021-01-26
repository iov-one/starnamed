package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/starname/types"
)

func addAccountCertificate(ctx sdk.Context, k Keeper, msg *types.MsgAddAccountCertificate) (*types.MsgAddAccountCertificateResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accounts := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl).WithConfiguration(conf)
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

	// collect fees
	if err := k.CollectProductFee(ctx, msg); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// add certificate
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
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

func deleteAccountCertificate(ctx sdk.Context, k Keeper, msg *types.MsgDeleteAccountCertificate) (*types.MsgDeleteAccountCertificateResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}
	// perform account checks, save certificate index
	accounts := k.AccountStore(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl)
	certIndex := new(int)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		CertificateExists(msg.DeleteCertificate, certIndex).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// delete cert
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
	ex.DeleteCertificate(*certIndex)

	// success; TODO emit event?
	return &types.MsgDeleteAccountCertificateResponse{}, nil
}

// deleteAccount account deletes the account from the system
func deleteAccount(ctx sdk.Context, k Keeper, msg *types.MsgDeleteAccount) (*types.MsgDeleteAccountResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.MustExist().Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accounts := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		DeletableBy(msg.Owner).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// delete account
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
	ex.Delete()

	// success; todo can we emit event?
	return &types.MsgDeleteAccountResponse{}, nil
}

// registerAccount registers an account
func registerAccount(ctx sdk.Context, k Keeper, msg *types.MsgRegisterAccount) (*types.MsgRegisterAccountResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.
		MustExist().
		NotExpired().
		Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accounts := k.AccountStore(ctx)
	d := domainCtrl.Domain()
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl).WithConfiguration(conf)
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

	// collect fees
	if err := k.CollectProductFee(ctx, msg, &d); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	ex := NewAccountExecutor(ctx, a).WithAccounts(&accounts)
	ex.Create()

	// success; todo can we emit event?
	return &types.MsgRegisterAccountResponse{}, nil
}

func renewAccount(ctx sdk.Context, k Keeper, msg *types.MsgRenewAccount) (*types.MsgRenewAccountResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains).WithConfiguration(conf)
	if err := domainCtrl.MustExist().Type(types.OpenDomain).Validate(); err != nil {
		return nil, err
	}

	// validate account
	accounts := k.AccountStore(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		Renewable().
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg, domainCtrl.domain); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// renew account
	// account valid until is extended here
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts).WithConfiguration(conf)
	ex.Renew()
	// get grace period and expiration time
	d := domainCtrl.Domain()
	dgp := conf.DomainGracePeriod
	domainGracePeriodUntil := utils.SecondsToTime(d.ValidUntil).Add(dgp)
	accNewValidUntil := utils.SecondsToTime(ex.State().ValidUntil)
	if domainGracePeriodUntil.Before(accNewValidUntil) {
		dex := NewDomainExecutor(ctx, domainCtrl.Domain()).WithDomains(&domains).WithAccounts(&accounts).WithConfiguration(conf)
		dex.Renew(accNewValidUntil.Unix())
	}

	// success; todo emit event??
	return &types.MsgRenewAccountResponse{}, nil
}

// handlerMsgReplaceAccountResources replaces account resources
func handlerMsgReplaceAccountResources(ctx sdk.Context, k Keeper, msg *types.MsgReplaceAccountResources) (*sdk.Result, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accounts := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		ValidResources(msg.NewResources).
		ResourceLimitNotExceeded(msg.NewResources).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// replace accounts resources
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
	ex.ReplaceResources(msg.NewResources)

	// success; TODO emit any useful event?
	return &sdk.Result{}, nil
}

// replaceAccountMetadata sets account metadata
func replaceAccountMetadata(ctx sdk.Context, k Keeper, msg *types.MsgReplaceAccountMetadata) (*types.MsgReplaceAccountMetadataResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}

	// perform account checks
	accounts := k.AccountStore(ctx)
	conf := k.ConfigurationKeeper.GetConfiguration(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl).WithConfiguration(conf)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		OwnedBy(msg.Owner).
		MetadataSizeNotExceeded(msg.NewMetadataURI).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// save to store
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
	ex.UpdateMetadata(msg.NewMetadataURI)

	// success TODO emit event
	return &types.MsgReplaceAccountMetadataResponse{}, nil
}

// transferAccount transfers account to a new owner and may clear resources and certificates
func transferAccount(ctx sdk.Context, k Keeper, msg *types.MsgTransferAccount) (*types.MsgTransferAccountResponse, error) {
	// perform domain checks
	domains := k.DomainStore(ctx)
	domainCtrl := NewDomainController(ctx, msg.Domain).WithDomains(&domains)
	if err := domainCtrl.MustExist().NotExpired().Validate(); err != nil {
		return nil, err
	}

	// check if account exists
	accounts := k.AccountStore(ctx)
	accountCtrl := NewAccountController(ctx, msg.Domain, msg.Name).WithAccounts(&accounts).WithDomainController(domainCtrl)
	if err := accountCtrl.
		MustExist().
		NotExpired().
		TransferableBy(msg.Owner).
		ResettableBy(msg.Owner, msg.ToReset).
		Validate(); err != nil {
		return nil, err
	}

	// collect fees
	if err := k.CollectProductFee(ctx, msg, domainCtrl.domain, k.AccountStore); err != nil {
		return nil, errors.Wrapf(err, "unable to collect fees")
	}

	// transfer account
	ex := NewAccountExecutor(ctx, accountCtrl.Account()).WithAccounts(&accounts)
	ex.Transfer(msg.NewOwner, msg.ToReset)

	// success, todo emit event?
	return &types.MsgTransferAccountResponse{}, nil
}
