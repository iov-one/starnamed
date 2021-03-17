package keeper

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

// FeeController defines the fee controller behaviour
// exists only in order to avoid devs creating a fee
// controller without using the constructor function
type FeeController interface {
	GetFee(msg sdk.Msg) sdk.Coin
	WithAccounts(store *crud.Store) FeeController
	WithDomain(domain *types.Domain) FeeController
}

// NewFeeController returns a new fee controller
func NewFeeController(ctx sdk.Context, fees *configuration.Fees) FeeController {
	return &feeApplier{
		moduleFees: *fees,
		ctx:        ctx,
	}
}

type feeApplier struct {
	moduleFees configuration.Fees
	ctx        sdk.Context
	store      *crud.Store
	domain     *types.Domain
}

// WithAccounts sets the feeApplier cached crud store
func (f *feeApplier) WithAccounts(store *crud.Store) FeeController {
	f.store = store
	return f
}

// WithDomain sets the feeApplier domain
func (f *feeApplier) WithDomain(domain *types.Domain) FeeController {
	f.domain = domain
	return f
}

func (f feeApplier) requireDomain() {
	if f.domain == nil {
		panic("domain is missing")
	}
}

func (f feeApplier) registerDomain() sdk.Dec {
	var registerDomainFee sdk.Dec
	f.requireDomain()
	level := len(f.domain.Name)
	switch level {
	case 1:
		registerDomainFee = f.moduleFees.RegisterDomain1
	case 2:
		registerDomainFee = f.moduleFees.RegisterDomain2
	case 3:
		registerDomainFee = f.moduleFees.RegisterDomain3
	case 4:
		registerDomainFee = f.moduleFees.RegisterDomain4
	case 5:
		registerDomainFee = f.moduleFees.RegisterDomain5
	default:
		registerDomainFee = f.moduleFees.RegisterDomainDefault
	}
	// if domain is open then we multiply
	if f.domain.Type == types.OpenDomain {
		registerDomainFee = registerDomainFee.Mul(f.moduleFees.RegisterOpenDomainMultiplier)
	}
	return registerDomainFee
}

func (f feeApplier) transferDomain() sdk.Dec {
	f.requireDomain()
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.TransferDomainOpen
	case types.ClosedDomain:
		return f.moduleFees.TransferDomainClosed
	}
	return f.moduleFees.FeeDefault
}

func (f feeApplier) renewDomain() sdk.Dec {
	f.requireDomain()
	if f.domain.Type == types.OpenDomain {
		return f.moduleFees.RenewDomainOpen
	}
	if f.store == nil {
		panic("store is missing")
	}
	var accountN int64
	cursor, err := (*f.store).Query().Where().Index(types.AccountDomainIndex).Equals(f.domain.PrimaryKey()).Do()
	if err != nil {
		panic(err)
	}
	for ; cursor.Valid(); cursor.Next() {
		accountN++
	}
	fee := f.moduleFees.RegisterAccountClosed
	fee = f.registerDomain().Add(fee.MulInt64(accountN))
	return fee
}

func (f feeApplier) registerAccount() sdk.Dec {
	f.requireDomain()
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.RegisterAccountOpen
	case types.ClosedDomain:
		return f.moduleFees.RegisterAccountClosed
	}
	return f.moduleFees.FeeDefault
}

func (f feeApplier) transferAccount() sdk.Dec {
	f.requireDomain()
	switch f.domain.Type {
	case types.ClosedDomain:
		return f.moduleFees.TransferAccountClosed
	case types.OpenDomain:
		return f.moduleFees.TransferAccountOpen
	}
	return f.moduleFees.FeeDefault
}

func (f feeApplier) renewAccount() sdk.Dec {
	f.requireDomain()
	switch f.domain.Type {
	case types.OpenDomain:
		return f.moduleFees.RegisterAccountOpen
	case types.ClosedDomain:
		return f.moduleFees.RegisterAccountClosed
	}
	return f.moduleFees.FeeDefault
}

func (f feeApplier) replaceResources() sdk.Dec {
	return f.moduleFees.ReplaceAccountResources
}

func (f feeApplier) addCert() sdk.Dec {
	return f.moduleFees.AddAccountCertificate
}

func (f feeApplier) delCert() sdk.Dec {
	return f.moduleFees.DelAccountCertificate
}

func (f feeApplier) setMetadata() sdk.Dec {
	return f.moduleFees.SetAccountMetadata
}

func (f feeApplier) defaultFee() sdk.Dec {
	return f.moduleFees.FeeDefault
}

func (f feeApplier) getFeeParam(msg sdk.Msg) sdk.Dec {
	switch msg.(type) {
	case *types.MsgTransferDomainInternal:
		return f.transferDomain()
	case *types.MsgRegisterDomainInternal:
		return f.registerDomain()
	case *types.MsgRenewDomainInternal:
		return f.renewDomain()
	case *types.MsgRegisterAccount:
		return f.registerAccount()
	case *types.MsgTransferAccount:
		return f.transferAccount()
	case *types.MsgRenewAccount:
		return f.renewAccount()
	case *types.MsgReplaceAccountResources:
		return f.replaceResources()
	case *types.MsgDeleteAccountCertificate:
		return f.delCert()
	case *types.MsgAddAccountCertificateInternal:
		return f.addCert()
	case *types.MsgReplaceAccountMetadata:
		return f.setMetadata()
	default:
		return f.defaultFee()
	}
}

// GetFee returns a fee based on the provided message
func (f feeApplier) GetFee(msg sdk.Msg) sdk.Coin {
	// get current price
	currentPrice := f.moduleFees.FeeCoinPrice
	// get fee parameter
	fee := f.getFeeParam(msg)
	// if fee is smaller than default fee, use default fee
	if fee.LT(f.defaultFee()) {
		fee = f.defaultFee()
	}
	// divide fee with current price
	toPay := fee.Quo(currentPrice)
	var feeAmount sdk.Int
	// get fee amount
	feeAmount = toPay.TruncateInt()
	// get coin denom
	coinDenom := f.moduleFees.FeeCoinDenom
	// generate coins to pay
	coinsToPay := sdk.NewCoin(coinDenom, feeAmount)
	return coinsToPay
}
