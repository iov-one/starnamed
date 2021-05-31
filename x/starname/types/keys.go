package types

// Module names
const (
	// ModuleName is the name of the module
	ModuleName = "starname"
	// DomainStore key defines the store key used to store domains information
	DomainStoreKey = "starname"
	// RouterKey defines the path used to interact with the domain module
	RouterKey    = ModuleName
	QuerierAlias = "starname"
	// QuerierRoute defines the query path used to interact with the domain module
	QuerierRoute = ModuleName
	// DefaultParamSpace defines the key for the default param space
	DefaultParamSpace = ModuleName
)

// Event attribute keys
const (
	AttributeKeyAccountName             = "account_name"
	AttributeKeyBroker                  = "broker"
	AttributeKeyRegisterer              = "registerer"
	AttributeKeyResources               = "resources"
	AttributeKeyDeletedCertificate      = "deleted_certificate"
	AttributeKeyDomainName              = "domain_name"
	AttributeKeyDomainType              = "domain_type"
	AttributeKeyNewCertificate          = "new_certificate"
	AttributeKeyNewMetadata             = "new_metadata"
	AttributeKeyNewResources            = "new_resources"
	AttributeKeyOwner                   = "owner"
	AttributeKeyTransferAccountNewOwner = "new_account_owner"
	AttributeKeyTransferAccountReset    = "transfer_account_reset"
	AttributeKeyTransferDomainFlag      = "transfer_domain_flag"
	AttributeKeyTransferDomainNewOwner  = "new_domain_owner"
)
