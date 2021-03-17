package types

import (
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// MsgWithFeePayer abstracts the Msg type to support a fee payer
// which takes care of handling product fees
type MsgWithFeePayer interface {
	sdk.Msg
	FeePayer() sdk.AccAddress
}

// MsgAddAccountCertificateInternal embeds MsgTransferDomain and adds sdk.Address properties for Owner and Payer
type MsgAddAccountCertificateInternal struct {
	MsgAddAccountCertificate
	Owner sdk.AccAddress
	Payer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgAddAccountCertificateInternal struct corresponding to the method receiver
func (m MsgAddAccountCertificate) ToInternal() *MsgAddAccountCertificateInternal {
	var err error
	var owner sdk.AccAddress = nil
	var payer sdk.AccAddress = nil

	if m.Owner != "" {
		owner, err = sdk.AccAddressFromBech32(m.Owner)
		if err != nil {
			panic(err)
		}
	}

	if m.Payer != "" {
		payer, err = sdk.AccAddressFromBech32(m.Payer)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgAddAccountCertificateInternal{
		MsgAddAccountCertificate: m,
		Owner:                    owner,
		Payer:                    payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgAddAccountCertificateInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgAddAccountCertificateInternal) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgAddAccountCertificate) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgAddAccountCertificate) Type() string {
	return "add_certificates_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgAddAccountCertificate) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewCertificate == nil {
		return errors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgAddAccountCertificate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgAddAccountCertificate) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}

	if m.Payer == "" {
		return []sdk.AccAddress{owner}
	}

	payer, err := sdk.AccAddressFromBech32(m.Payer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{payer, owner}
}

var _ MsgWithFeePayer = (*MsgDeleteAccountCertificate)(nil)

// FeePayer implements FeePayer interface
func (m *MsgDeleteAccountCertificate) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgDeleteAccountCertificate) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteAccountCertificate) Type() string {
	return "delete_certificate_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteAccountCertificate) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.DeleteCertificate == nil {
		return errors.Wrap(ErrInvalidRequest, "certificate is empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccountCertificate) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Owner}
	}
	return []sdk.AccAddress{m.Payer, m.Owner}
}

var _ MsgWithFeePayer = (*MsgDeleteAccount)(nil)

// FeePayer implements FeePayer interface
func (m *MsgDeleteAccount) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgDeleteAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteAccount) Type() string {
	return "delete_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteAccount) ValidateBasic() error {
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrOpEmptyAcc, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteAccount) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Owner}
	}
	return []sdk.AccAddress{m.Payer, m.Owner}
}

// MsgDeleteDomainInternal embeds MsgDeleteDomain and adds sdk.Address properties for Owner and Payer
type MsgDeleteDomainInternal struct {
	MsgDeleteDomain
	Owner sdk.AccAddress
	Payer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgDeleteDomainInternal struct corresponding to the method receiver
func (m MsgDeleteDomain) ToInternal() *MsgDeleteDomainInternal {
	var err error
	var owner sdk.AccAddress = nil
	var payer sdk.AccAddress = nil

	if m.Owner != "" {
		owner, err = sdk.AccAddressFromBech32(m.Owner)
		if err != nil {
			panic(err)
		}
	}

	if m.Payer != "" {
		payer, err = sdk.AccAddressFromBech32(m.Payer)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgDeleteDomainInternal{
		MsgDeleteDomain: m,
		Owner:           owner,
		Payer:           payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgDeleteDomainInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgDeleteDomainInternal) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgDeleteDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgDeleteDomain) Type() string {
	return "delete_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgDeleteDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgDeleteDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgDeleteDomain) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}

	if m.Payer == "" {
		return []sdk.AccAddress{owner}
	}

	payer, err := sdk.AccAddressFromBech32(m.Payer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{payer, owner}
}

var _ MsgWithFeePayer = (*MsgRegisterAccount)(nil)

// FeePayer implements FeePayer interface
func (m *MsgRegisterAccount) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Registerer
}

// Route implements sdk.Msg
func (m *MsgRegisterAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRegisterAccount) Type() string {
	return "register_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRegisterAccount) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner.Empty() {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.Registerer.Empty() {
		return errors.Wrap(ErrInvalidRegisterer, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterAccount) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Registerer}
	}
	return []sdk.AccAddress{m.Payer, m.Registerer}
}

// MsgRegisterDomainInternal embeds MsgRegisterDomain and adds sdk.Address properties for Admin, Broker, and Payer
type MsgRegisterDomainInternal struct {
	MsgRegisterDomain
	Admin  sdk.AccAddress
	Broker sdk.AccAddress
	Payer  sdk.AccAddress
}

// ToInternal returns a pointer to the MsgRegisterDomainInternal struct corresponding to the method receiver
func (m MsgRegisterDomain) ToInternal() *MsgRegisterDomainInternal {
	var err error
	var admin sdk.AccAddress = nil
	var broker sdk.AccAddress = nil
	var payer sdk.AccAddress = nil

	if m.Admin != "" {
		admin, err = sdk.AccAddressFromBech32(m.Admin)
		if err != nil {
			panic(err)
		}
	}

	if m.Broker != "" {
		broker, err = sdk.AccAddressFromBech32(m.Broker)
		if err != nil {
			panic(err)
		}
	}

	if m.Payer != "" {
		payer, err = sdk.AccAddressFromBech32(m.Payer)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgRegisterDomainInternal{
		MsgRegisterDomain: m,
		Admin:             admin,
		Broker:            broker,
		Payer:             payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgRegisterDomainInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgRegisterDomainInternal) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Admin
}

// Route implements sdk.Msg
func (m *MsgRegisterDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRegisterDomain) Type() string {
	return "register_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRegisterDomain) ValidateBasic() error {
	if m.Admin == "" {
		return errors.Wrap(ErrInvalidRequest, "admin is missing")
	}
	if err := ValidateDomainType(m.DomainType); err != nil {
		return err
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRegisterDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRegisterDomain) GetSigners() []sdk.AccAddress {
	admin, err := sdk.AccAddressFromBech32(m.Admin)
	if err != nil {
		panic(err)
	}

	if m.Payer == "" {
		return []sdk.AccAddress{admin}
	}

	payer, err := sdk.AccAddressFromBech32(m.Payer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{payer, admin}
}

var _ MsgWithFeePayer = (*MsgRenewAccount)(nil)

// FeePayer implements FeePayer interface
func (m *MsgRenewAccount) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Signer
}

// Route implements sdk.Msg
func (m *MsgRenewAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRenewAccount) Type() string {
	return "renew_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRenewAccount) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRenewAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRenewAccount) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Signer}
	}
	return []sdk.AccAddress{m.Payer, m.Signer}
}

// MsgRenewDomainInternal embeds MsgDeleteDomain and adds sdk.Address properties for Owner and Payer
type MsgRenewDomainInternal struct {
	MsgRenewDomain
	Signer sdk.AccAddress
	Payer  sdk.AccAddress
}

// ToInternal returns a pointer to the MsgRenewDomainInternal struct corresponding to the method receiver
func (m MsgRenewDomain) ToInternal() *MsgRenewDomainInternal {
	var err error
	var signer sdk.AccAddress = nil
	var payer sdk.AccAddress = nil

	if m.Signer != "" {
		signer, err = sdk.AccAddressFromBech32(m.Signer)
		if err != nil {
			panic(err)
		}
	}

	if m.Payer != "" {
		payer, err = sdk.AccAddressFromBech32(m.Payer)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgRenewDomainInternal{
		MsgRenewDomain: m,
		Signer:         signer,
		Payer:          payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgRenewDomainInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgRenewDomainInternal) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Signer
}

// Route implements sdk.Msg
func (m *MsgRenewDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgRenewDomain) Type() string {
	return "renew_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgRenewDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgRenewDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgRenewDomain) GetSigners() []sdk.AccAddress {
	signer, err := sdk.AccAddressFromBech32(m.Signer)
	if err != nil {
		panic(err)
	}

	if m.Payer == "" {
		return []sdk.AccAddress{signer}
	}

	payer, err := sdk.AccAddressFromBech32(m.Payer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{payer, signer}
}

var _ MsgWithFeePayer = (*MsgReplaceAccountResources)(nil)

// FeePayer implements FeePayer interface
func (m *MsgReplaceAccountResources) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgReplaceAccountResources) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgReplaceAccountResources) Type() string {
	return "replace_account_resources"
}

// ValidateBasic implements sdk.Msg
func (m *MsgReplaceAccountResources) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgReplaceAccountResources) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountResources) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Owner}
	}
	return []sdk.AccAddress{m.Payer, m.Owner}
}

var _ MsgWithFeePayer = (*MsgReplaceAccountMetadata)(nil)

// FeePayer implements FeePayer interface
func (m *MsgReplaceAccountMetadata) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgReplaceAccountMetadata) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgReplaceAccountMetadata) Type() string {
	return "set_account_metadata"
}

// ValidateBasic implements sdk.Msg
func (m *MsgReplaceAccountMetadata) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrapf(ErrInvalidDomainName, "empty")
	}
	if m.Owner.Empty() {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgReplaceAccountMetadata) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgReplaceAccountMetadata) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Owner}
	}
	return []sdk.AccAddress{m.Payer, m.Owner}
}

var _ MsgWithFeePayer = (*MsgTransferAccount)(nil)

// FeePayer implements FeePayer interface
func (m *MsgTransferAccount) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgTransferAccount) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgTransferAccount) Type() string {
	return "transfer_account"
}

// ValidateBasic implements sdk.Msg
func (m *MsgTransferAccount) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Name == "" {
		return errors.Wrap(ErrOpEmptyAcc, "empty")
	}
	if m.Owner == nil {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewOwner == nil {
		return errors.Wrap(ErrInvalidOwner, "new owner is empty")
	}
	// success
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgTransferAccount) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgTransferAccount) GetSigners() []sdk.AccAddress {
	if m.Payer.Empty() {
		return []sdk.AccAddress{m.Owner}
	}
	return []sdk.AccAddress{m.Payer, m.Owner}
}

// TransferFlag defines the type of domain transfer
type TransferFlag int

const (
	// TransferFlush clears all domain account data, except empty account)
	TransferFlush = iota
	// TransferOwned transfers only accounts owned by the current owner
	TransferOwned
	// TransferResetNone leaves things as they are except for empty account
	TransferResetNone
	// TransferAll is not available is here only for tests backwards compatibility and will be removed. TODO deprecate
	TransferAll
)

// MsgTransferDomainInternal embeds MsgTransferDomain and adds sdk.Address properties for Owner, Payer, and NewAdmin
type MsgTransferDomainInternal struct {
	MsgTransferDomain
	Owner    sdk.AccAddress
	Payer    sdk.AccAddress
	NewAdmin sdk.AccAddress
}

// ToInternal returns a pointer to the MsgTransferDomainInternal struct corresponding to the method receiver
func (m MsgTransferDomain) ToInternal() *MsgTransferDomainInternal {
	var err error
	var owner sdk.AccAddress = nil
	var payer sdk.AccAddress = nil
	var newAdmin sdk.AccAddress = nil

	if m.Owner != "" {
		owner, err = sdk.AccAddressFromBech32(m.Owner)
		if err != nil {
			panic(err)
		}
	}

	if m.Payer != "" {
		payer, err = sdk.AccAddressFromBech32(m.Payer)
		if err != nil {
			panic(err)
		}
	}

	if m.NewAdmin != "" {
		newAdmin, err = sdk.AccAddressFromBech32(m.NewAdmin)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgTransferDomainInternal{
		MsgTransferDomain: m,
		Owner:             owner,
		Payer:             payer,
		NewAdmin:          newAdmin,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgTransferDomainInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgTransferDomainInternal) FeePayer() sdk.AccAddress {
	if !m.Payer.Empty() {
		return m.Payer
	}
	return m.Owner
}

// Route implements sdk.Msg
func (m *MsgTransferDomain) Route() string {
	return RouterKey
}

// Type implements sdk.Msg
func (m *MsgTransferDomain) Type() string {
	return "transfer_domain"
}

// ValidateBasic implements sdk.Msg
func (m *MsgTransferDomain) ValidateBasic() error {
	if m.Domain == "" {
		return errors.Wrap(ErrInvalidDomainName, "empty")
	}
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewAdmin == "" {
		return errors.Wrap(ErrInvalidRequest, "new admin is empty")
	}
	switch m.TransferFlag {
	case TransferOwned:
	case TransferResetNone:
	case TransferFlush:
	default:
		return errors.Wrapf(ErrInvalidRequest, "unknown reset flag: %d", m.TransferFlag)
	}
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *MsgTransferDomain) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(m))
}

// GetSigners implements sdk.Msg
func (m *MsgTransferDomain) GetSigners() []sdk.AccAddress {
	owner, err := sdk.AccAddressFromBech32(m.Owner)
	if err != nil {
		panic(err)
	}

	if m.Payer == "" {
		return []sdk.AccAddress{owner}
	}

	payer, err := sdk.AccAddressFromBech32(m.Payer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{payer, owner}
}
