package types

import (
	"testing"

	"github.com/cosmos/cosmos-sdk/codec"
	codectypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFees_Validate(t *testing.T) {
	type fields struct {
		FeeCoinDenom                 string
		FeeCoinPrice                 types.Dec
		FeeDefault                   types.Dec
		RegisterAccountClosed        types.Dec
		RegisterAccountOpen          types.Dec
		TransferAccountClosed        types.Dec
		TransferAccountOpen          types.Dec
		ReplaceAccountResources      types.Dec
		AddAccountCertificate        types.Dec
		DelAccountCertificate        types.Dec
		SetAccountMetadata           types.Dec
		RegisterDomain1              types.Dec
		RegisterDomain2              types.Dec
		RegisterDomain3              types.Dec
		RegisterDomain4              types.Dec
		RegisterDomain5              types.Dec
		RegisterDomainDefault        types.Dec
		RegisterOpenDomainMultiplier types.Dec
		TransferDomainClosed         types.Dec
		TransferDomainOpen           types.Dec
		RenewDomainOpen              types.Dec
		CreateEscrow                 types.Dec
		UpdateEscrow                 types.Dec
		TransferToEscrow             types.Dec
		RefundEscrow                 types.Dec
	}
	tests := []struct {
		name    string
		fields  fields
		wantErr bool
	}{
		{
			name: "success",
			fields: func() fields {
				fees := NewFees()
				fees.SetDefaults("test")
				return fields(*fees)
			}(),
			wantErr: false,
		},
		{
			name:    "fail missing fee",
			fields:  fields{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &Fees{
				FeeCoinDenom:                 tt.fields.FeeCoinDenom,
				FeeCoinPrice:                 tt.fields.FeeCoinPrice,
				FeeDefault:                   tt.fields.FeeDefault,
				RegisterAccountClosed:        tt.fields.RegisterAccountClosed,
				RegisterAccountOpen:          tt.fields.RegisterAccountOpen,
				TransferAccountClosed:        tt.fields.TransferAccountClosed,
				TransferAccountOpen:          tt.fields.TransferAccountOpen,
				ReplaceAccountResources:      tt.fields.ReplaceAccountResources,
				AddAccountCertificate:        tt.fields.AddAccountCertificate,
				DelAccountCertificate:        tt.fields.DelAccountCertificate,
				SetAccountMetadata:           tt.fields.SetAccountMetadata,
				RegisterDomain1:              tt.fields.RegisterDomain1,
				RegisterDomain2:              tt.fields.RegisterDomain2,
				RegisterDomain3:              tt.fields.RegisterDomain3,
				RegisterDomain4:              tt.fields.RegisterDomain4,
				RegisterDomain5:              tt.fields.RegisterDomain5,
				RegisterDomainDefault:        tt.fields.RegisterDomainDefault,
				RegisterOpenDomainMultiplier: tt.fields.RegisterOpenDomainMultiplier,
				TransferDomainClosed:         tt.fields.TransferDomainClosed,
				TransferDomainOpen:           tt.fields.TransferDomainOpen,
				RenewDomainOpen:              tt.fields.RenewDomainOpen,
				CreateEscrow:                 tt.fields.CreateEscrow,
				UpdateEscrow:                 tt.fields.UpdateEscrow,
				TransferToEscrow:             tt.fields.TransferToEscrow,
				RefundEscrow:                 tt.fields.RefundEscrow,
			}
			if err := f.Validate(); (err != nil) != tt.wantErr {
				t.Errorf("Validate() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestFees_UnmarshalJson(t *testing.T) {
	specs := map[string]struct {
		src string
		exp Fees
	}{
		"defaults": {
			src: `{
				"fee_coin_denom": "tiov",
				"fee_coin_price": "10.000000000000000000",
				"fee_default": "10.000000000000000000",
				"register_account_closed": "10.000000000000000000",
				"register_account_open": "10.000000000000000000",
				"transfer_account_closed": "10.000000000000000000",
				"transfer_account_open": "10.000000000000000000",
				"replace_account_resources": "10.000000000000000000",
				"add_account_certificate": "10.000000000000000000",
				"del_account_certificate": "10.000000000000000000",
				"set_account_metadata": "10.000000000000000000",
				"register_domain_1": "10.000000000000000000",
				"register_domain_2": "10.000000000000000000",
				"register_domain_3": "10.000000000000000000",
				"register_domain_4": "10.000000000000000000",
				"register_domain_5": "10.000000000000000000",
				"register_domain_default": "10.000000000000000000",
				"register_open_domain_multiplier": "2.000000000000000000",
				"transfer_domain_closed": "10.000000000000000000",
				"transfer_domain_open": "10.000000000000000000",
				"renew_domain_open": "10.000000000000000000",
				"create_escrow": "10.000000000000000000",
				"update_escrow": "10.000000000000000000",
				"transfer_to_escrow": "10.000000000000000000",
				"refund_escrow": "10.000000000000000000"
			}`,
			exp: func() Fees {
				fees := NewFees()
				fees.SetDefaults("tiov")
				return *fees
			}(),
		},
	}
	for msg, spec := range specs {
		t.Run(msg, func(t *testing.T) {
			interfaceRegistry := codectypes.NewInterfaceRegistry()
			marshaler := codec.NewProtoCodec(interfaceRegistry)

			var fees Fees
			err := marshaler.UnmarshalJSON([]byte(spec.src), &fees)
			require.NoError(t, err)
			assert.Equal(t, spec.exp, fees)
		})
	}
}
