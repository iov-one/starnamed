package keeper

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/pkg/utils"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
)

// decFromStr is a helper to convert string decimals such as 0.12311 easily
func decFromStr(str string) sdk.Dec {
	dec, err := sdk.NewDecFromStr(str)
	if err != nil {
		panic(err)
	}
	return dec
}

func Test_FeeApplier(t *testing.T) {
	fee := configuration.Fees{
		FeeCoinDenom:                 "tiov",
		FeeCoinPrice:                 decFromStr("2"),
		FeeDefault:                   decFromStr("3"),
		RegisterAccountClosed:        decFromStr("5"),
		RegisterAccountOpen:          sdk.NewDec(7),
		TransferAccountClosed:        sdk.NewDec(11),
		TransferAccountOpen:          sdk.NewDec(13),
		ReplaceAccountResources:      sdk.NewDec(17),
		AddAccountCertificate:        sdk.NewDec(19),
		DelAccountCertificate:        sdk.NewDec(23),
		SetAccountMetadata:           sdk.NewDec(31),
		RegisterDomain1:              sdk.NewDec(37),
		RegisterDomain2:              sdk.NewDec(41),
		RegisterDomain3:              sdk.NewDec(43),
		RegisterDomain4:              sdk.NewDec(47),
		RegisterDomain5:              sdk.NewDec(53),
		RegisterDomainDefault:        sdk.NewDec(59),
		RegisterOpenDomainMultiplier: sdk.NewDec(61),
		TransferDomainClosed:         sdk.NewDec(67),
		TransferDomainOpen:           sdk.NewDec(71),
		RenewDomainOpen:              sdk.NewDec(73),
	}
	cases := map[string]struct {
		Msg         sdk.Msg
		Domain      types.Domain
		ExpectedFee sdk.Dec
	}{
		"register closed domain 5": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "test1",
			},
			ExpectedFee: sdk.NewDec(26),
		},
		"register closed domain 4": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "test",
			},
			ExpectedFee: sdk.NewDec(23),
		},
		"register closed domain 3": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "tes",
			},
			ExpectedFee: sdk.NewDec(21),
		},
		"register closed domain 2": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "te",
			},
			ExpectedFee: sdk.NewDec(20),
		},
		"register closed domain 1": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "t",
			},
			ExpectedFee: sdk.NewDec(18),
		},
		"register open domain 5": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "test1",
				Type: types.OpenDomain,
			},
			ExpectedFee: sdk.NewDec(1616),
		},
		"register open domain default": {
			Msg: &types.MsgRegisterDomainInternal{},
			Domain: types.Domain{
				Name: "test12",
				Type: types.ClosedDomain,
			},
			ExpectedFee: sdk.NewDec(29),
		},
		"transfer domain open": {
			Msg:         &types.MsgTransferDomainInternal{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(35),
		},
		"transfer domain closed": {
			Msg:         &types.MsgTransferDomainInternal{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(33),
		},
		"set metadata": {
			Msg:         &types.MsgReplaceAccountMetadata{},
			ExpectedFee: sdk.NewDec(15),
		},
		"delete certs": {
			Msg:         &types.MsgDeleteAccountCertificate{},
			ExpectedFee: sdk.NewDec(11),
		},
		"add certs": {
			Msg:         &types.MsgAddAccountCertificateInternal{},
			ExpectedFee: sdk.NewDec(9),
		},
		"replace resources": {
			Msg:         &types.MsgReplaceAccountResources{},
			ExpectedFee: sdk.NewDec(8),
		},
		"transfer account closed": {
			Msg:         &types.MsgTransferAccount{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(5),
		},
		"transfer account open": {
			Msg:         &types.MsgTransferAccount{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(6),
		},
		"register account open": {
			Msg:         &types.MsgRegisterAccount{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(3),
		},
		"register account closed": {
			Msg:         &types.MsgRegisterAccount{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(2),
		},
		"renew account closed": {
			Msg:         &types.MsgRenewAccount{},
			Domain:      types.Domain{Type: types.ClosedDomain},
			ExpectedFee: sdk.NewDec(2),
		},
		"renew account open": {
			Msg:         &types.MsgRenewAccount{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(3),
		},
		"renew domain open": {
			Msg:         &types.MsgRenewDomainInternal{},
			Domain:      types.Domain{Type: types.OpenDomain},
			ExpectedFee: sdk.NewDec(36),
		},
		"renew domain closed": {
			Msg:         &types.MsgRenewDomainInternal{},
			Domain:      types.Domain{Type: types.ClosedDomain, Name: "renew"},
			ExpectedFee: sdk.NewDec(34), // 5-char domain + its three accounts-> "", "1", "2"; so 53/2 + 5/2*3=34
		},
		/* TODO: FIXME
		"default fee unknown message": {
			Msg:         &DullMsg{},
			ExpectedFee: sdk.NewDec(1),
		},
		*/
		"use default fee": {
			Msg:         &types.MsgRenewDomainInternal{},
			Domain:      types.Domain{Type: types.ClosedDomain, Name: "default-fee"},
			ExpectedFee: sdk.NewDec(32), // 6+char domain + one account ""; so 59/2 + 5/2 * 1 = 32
		},
	}
	k, ctx, _ := NewTestKeeper(t, true)
	ds := k.DomainStore(ctx)
	as := k.AccountStore(ctx)
	ds.Create(&types.Domain{Name: "renew", Admin: AliceKey})
	as.Create(&types.Account{Domain: "renew", Name: utils.StrPtr(types.EmptyAccountName), Owner: AliceKey})
	as.Create(&types.Account{Domain: "renew", Name: utils.StrPtr("1"), Owner: AliceKey})
	as.Create(&types.Account{Domain: "renew", Name: utils.StrPtr("2"), Owner: AliceKey})
	ds.Create(&types.Domain{Name: "default-fee", Admin: AliceKey})
	as.Create(&types.Account{Domain: "default-fee", Name: utils.StrPtr(types.EmptyAccountName), Owner: AliceKey})

	k.ConfigurationKeeper.(ConfigurationSetter).SetFees(ctx, &fee)
	for name, c := range cases {
		t.Run(name, func(t *testing.T) {
			ctrl := NewFeeController(ctx, &fee).WithDomain(&c.Domain).WithAccounts(&as)
			got := ctrl.GetFee(c.Msg)
			if !got.Amount.Equal(c.ExpectedFee.RoundInt()) {
				t.Fatalf("expected fee: %s, got %s", c.ExpectedFee, got.Amount)
			}
		})
	}
}
