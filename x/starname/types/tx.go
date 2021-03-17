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

// MsgDeleteAccountCertificateInternal embeds MsgDeleteAccountCertificate and adds sdk.Address properties for Owner and Payer
type MsgDeleteAccountCertificateInternal struct {
	MsgDeleteAccountCertificate
	Owner sdk.AccAddress
	Payer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgAddAccountCertificateInternal struct corresponding to the method receiver
func (m MsgDeleteAccountCertificate) ToInternal() *MsgDeleteAccountCertificateInternal {
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

	msgi := MsgDeleteAccountCertificateInternal{
		MsgDeleteAccountCertificate: m,
		Owner:                       owner,
		Payer:                       payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgDeleteAccountCertificateInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgDeleteAccountCertificateInternal) FeePayer() sdk.AccAddress {
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
	if m.Owner == "" {
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

// MsgDeleteAccountInternal embeds MsgDeleteDomain and adds sdk.Address properties for Owner and Payer
type MsgDeleteAccountInternal struct {
	MsgDeleteAccount
	Owner sdk.AccAddress
	Payer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgDeleteAccountInternal struct corresponding to the method receiver
func (m MsgDeleteAccount) ToInternal() *MsgDeleteAccountInternal {
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

	msgi := MsgDeleteAccountInternal{
		MsgDeleteAccount: m,
		Owner:            owner,
		Payer:            payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgDeleteAccountInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgDeleteAccountInternal) FeePayer() sdk.AccAddress {
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
	if m.Owner == "" {
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

// MsgRegisterAccountInternal embeds MsgRegisterAccount and adds sdk.Address properties for Admin, Broker, and Payer
type MsgRegisterAccountInternal struct {
	MsgRegisterAccount
	Owner      sdk.AccAddress
	Broker     sdk.AccAddress
	Payer      sdk.AccAddress
	Registerer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgRegisterAccountInternal struct corresponding to the method receiver
func (m MsgRegisterAccount) ToInternal() *MsgRegisterAccountInternal {
	var err error
	var owner sdk.AccAddress = nil
	var broker sdk.AccAddress = nil
	var payer sdk.AccAddress = nil
	var registerer sdk.AccAddress = nil

	if m.Owner != "" {
		owner, err = sdk.AccAddressFromBech32(m.Owner)
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

	if m.Registerer != "" {
		registerer, err = sdk.AccAddressFromBech32(m.Registerer)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgRegisterAccountInternal{
		MsgRegisterAccount: m,
		Owner:              owner,
		Broker:             broker,
		Payer:              payer,
		Registerer:         registerer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgRegisterAccountInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgRegisterAccountInternal) FeePayer() sdk.AccAddress {
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
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.Registerer == "" {
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
	registerer, err := sdk.AccAddressFromBech32(m.Registerer)
	if err != nil {
		panic(err)
	}

	if m.Payer == "" {
		return []sdk.AccAddress{registerer}
	}

	payer, err := sdk.AccAddressFromBech32(m.Payer)
	if err != nil {
		panic(err)
	}

	return []sdk.AccAddress{payer, registerer}
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

// MsgRenewAccountInternal embeds MsgRenewAccount and adds sdk.Address properties for Signer and Payer
type MsgRenewAccountInternal struct {
	MsgRenewAccount
	Signer sdk.AccAddress
	Payer  sdk.AccAddress
}

// ToInternal returns a pointer to the MsgRenewAccountInternal struct corresponding to the method receiver
func (m MsgRenewAccount) ToInternal() *MsgRenewAccountInternal {
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

	msgi := MsgRenewAccountInternal{
		MsgRenewAccount: m,
		Signer:          signer,
		Payer:           payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgRenewAccountInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgRenewAccountInternal) FeePayer() sdk.AccAddress {
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

// MsgRenewDomainInternal embeds MsgDeleteDomain and adds sdk.Address properties for Signer and Payer
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

// MsgReplaceAccountResourcesInternal embeds MsgReplaceAccountResources and adds sdk.Address properties for Owner and Payer
type MsgReplaceAccountResourcesInternal struct {
	MsgReplaceAccountResources
	Owner sdk.AccAddress
	Payer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgReplaceAccountResourcesInternal struct corresponding to the method receiver
func (m MsgReplaceAccountResources) ToInternal() *MsgReplaceAccountResourcesInternal {
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

	msgi := MsgReplaceAccountResourcesInternal{
		MsgReplaceAccountResources: m,
		Owner:                      owner,
		Payer:                      payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgReplaceAccountResourcesInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgReplaceAccountResourcesInternal) FeePayer() sdk.AccAddress {
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
	if m.Owner == "" {
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

// MsgReplaceAccountMetadataInternal embeds MsgReplaceAccountMetadata and adds sdk.Address properties for Owner and Payer
type MsgReplaceAccountMetadataInternal struct {
	MsgReplaceAccountMetadata
	Owner sdk.AccAddress
	Payer sdk.AccAddress
}

// ToInternal returns a pointer to the MsgReplaceAccountMetadataInternal struct corresponding to the method receiver
func (m MsgReplaceAccountMetadata) ToInternal() *MsgReplaceAccountMetadataInternal {
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

	msgi := MsgReplaceAccountMetadataInternal{
		MsgReplaceAccountMetadata: m,
		Owner:                     owner,
		Payer:                     payer,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgReplaceAccountMetadataInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgReplaceAccountMetadataInternal) FeePayer() sdk.AccAddress {
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
	if m.Owner == "" {
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

// MsgTransferAccountInternal embeds MsgTransferAccount and adds sdk.Address properties for Owner, Payer, and NewOwner
type MsgTransferAccountInternal struct {
	MsgTransferAccount
	Owner    sdk.AccAddress
	Payer    sdk.AccAddress
	NewOwner sdk.AccAddress
}

// ToInternal returns a pointer to the MsgTransferAccountInternal struct corresponding to the method receiver
func (m MsgTransferAccount) ToInternal() *MsgTransferAccountInternal {
	var err error
	var owner sdk.AccAddress = nil
	var payer sdk.AccAddress = nil
	var newOwner sdk.AccAddress = nil

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

	if m.NewOwner != "" {
		newOwner, err = sdk.AccAddressFromBech32(m.NewOwner)
		if err != nil {
			panic(err)
		}
	}

	msgi := MsgTransferAccountInternal{
		MsgTransferAccount: m,
		Owner:              owner,
		Payer:              payer,
		NewOwner:           newOwner,
	}

	return &msgi
}

var _ MsgWithFeePayer = (*MsgTransferAccountInternal)(nil)

// FeePayer implements FeePayer interface
func (m *MsgTransferAccountInternal) FeePayer() sdk.AccAddress {
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
	if m.Owner == "" {
		return errors.Wrap(ErrInvalidOwner, "empty")
	}
	if m.NewOwner == "" {
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
