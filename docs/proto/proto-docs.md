<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

<<<<<<< HEAD
- [iov/configuration/v1beta1/types.proto](#iov/configuration/v1beta1/types.proto)
    - [Config](#starnamed.x.configuration.v1beta1.Config)
    - [Fees](#starnamed.x.configuration.v1beta1.Fees)
    - [GenesisState](#starnamed.x.configuration.v1beta1.GenesisState)
  
- [iov/configuration/v1beta1/msgs.proto](#iov/configuration/v1beta1/msgs.proto)
    - [MsgUpdateConfig](#starnamed.x.configuration.v1beta1.MsgUpdateConfig)
    - [MsgUpdateFees](#starnamed.x.configuration.v1beta1.MsgUpdateFees)
  
- [iov/configuration/v1beta1/query.proto](#iov/configuration/v1beta1/query.proto)
    - [QueryConfigRequest](#starnamed.x.configuration.v1beta1.QueryConfigRequest)
    - [QueryConfigResponse](#starnamed.x.configuration.v1beta1.QueryConfigResponse)
    - [QueryFeesRequest](#starnamed.x.configuration.v1beta1.QueryFeesRequest)
    - [QueryFeesResponse](#starnamed.x.configuration.v1beta1.QueryFeesResponse)
  
    - [Query](#starnamed.x.configuration.v1beta1.Query)
  
- [iov/escrow/v1beta1/events.proto](#iov/escrow/v1beta1/events.proto)
    - [EventCompletedEscrow](#starnamed.x.escrow.v1beta1.EventCompletedEscrow)
    - [EventCreatedEscrow](#starnamed.x.escrow.v1beta1.EventCreatedEscrow)
    - [EventRefundedEscrow](#starnamed.x.escrow.v1beta1.EventRefundedEscrow)
    - [EventUpdatedEscrow](#starnamed.x.escrow.v1beta1.EventUpdatedEscrow)
  
- [iov/escrow/v1beta1/types.proto](#iov/escrow/v1beta1/types.proto)
    - [Escrow](#starnamed.x.escrow.v1beta1.Escrow)
  
    - [EscrowState](#starnamed.x.escrow.v1beta1.EscrowState)
  
- [iov/escrow/v1beta1/params.proto](#iov/escrow/v1beta1/params.proto)
    - [Params](#starnamed.x.escrow.v1beta1.Params)
  
- [iov/escrow/v1beta1/genesis.proto](#iov/escrow/v1beta1/genesis.proto)
    - [GenesisState](#starnamed.x.escrow.v1beta1.GenesisState)
  
- [iov/escrow/v1beta1/query.proto](#iov/escrow/v1beta1/query.proto)
    - [QueryEscrowRequest](#starnamed.x.escrow.v1beta1.QueryEscrowRequest)
    - [QueryEscrowResponse](#starnamed.x.escrow.v1beta1.QueryEscrowResponse)
    - [QueryEscrowsRequest](#starnamed.x.escrow.v1beta1.QueryEscrowsRequest)
    - [QueryEscrowsResponse](#starnamed.x.escrow.v1beta1.QueryEscrowsResponse)
  
    - [Query](#starnamed.x.escrow.v1beta1.Query)
  
- [iov/escrow/v1beta1/test.proto](#iov/escrow/v1beta1/test.proto)
    - [TestObject](#starnamed.x.escrow.v1beta1.TestObject)
    - [TestTimeConstrainedObject](#starnamed.x.escrow.v1beta1.TestTimeConstrainedObject)
  
- [iov/escrow/v1beta1/tx.proto](#iov/escrow/v1beta1/tx.proto)
    - [MsgCreateEscrow](#starnamed.x.escrow.v1beta1.MsgCreateEscrow)
    - [MsgCreateEscrowResponse](#starnamed.x.escrow.v1beta1.MsgCreateEscrowResponse)
    - [MsgRefundEscrow](#starnamed.x.escrow.v1beta1.MsgRefundEscrow)
    - [MsgRefundEscrowResponse](#starnamed.x.escrow.v1beta1.MsgRefundEscrowResponse)
    - [MsgTransferToEscrow](#starnamed.x.escrow.v1beta1.MsgTransferToEscrow)
    - [MsgTransferToEscrowResponse](#starnamed.x.escrow.v1beta1.MsgTransferToEscrowResponse)
    - [MsgUpdateEscrow](#starnamed.x.escrow.v1beta1.MsgUpdateEscrow)
    - [MsgUpdateEscrowResponse](#starnamed.x.escrow.v1beta1.MsgUpdateEscrowResponse)
  
    - [Msg](#starnamed.x.escrow.v1beta1.Msg)
  
- [iov/offchain/v1alpha1/offchain.proto](#iov/offchain/v1alpha1/offchain.proto)
    - [ListOfMsgSignData](#cosmos.offchain.v1alpha1.ListOfMsgSignData)
    - [MsgSignData](#cosmos.offchain.v1alpha1.MsgSignData)
  
- [iov/starname/v1beta1/types.proto](#iov/starname/v1beta1/types.proto)
    - [Account](#starnamed.x.starname.v1beta1.Account)
    - [Domain](#starnamed.x.starname.v1beta1.Domain)
    - [Resource](#starnamed.x.starname.v1beta1.Resource)
  
- [iov/starname/v1beta1/genesis.proto](#iov/starname/v1beta1/genesis.proto)
    - [GenesisState](#starnamed.x.starname.v1beta1.GenesisState)
  
- [iov/starname/v1beta1/query.proto](#iov/starname/v1beta1/query.proto)
    - [QueryBrokerAccountsRequest](#starnamed.x.starname.v1beta1.QueryBrokerAccountsRequest)
    - [QueryBrokerAccountsResponse](#starnamed.x.starname.v1beta1.QueryBrokerAccountsResponse)
    - [QueryBrokerDomainsRequest](#starnamed.x.starname.v1beta1.QueryBrokerDomainsRequest)
    - [QueryBrokerDomainsResponse](#starnamed.x.starname.v1beta1.QueryBrokerDomainsResponse)
    - [QueryDomainAccountsRequest](#starnamed.x.starname.v1beta1.QueryDomainAccountsRequest)
    - [QueryDomainAccountsResponse](#starnamed.x.starname.v1beta1.QueryDomainAccountsResponse)
    - [QueryDomainRequest](#starnamed.x.starname.v1beta1.QueryDomainRequest)
    - [QueryDomainResponse](#starnamed.x.starname.v1beta1.QueryDomainResponse)
    - [QueryOwnerAccountsRequest](#starnamed.x.starname.v1beta1.QueryOwnerAccountsRequest)
    - [QueryOwnerAccountsResponse](#starnamed.x.starname.v1beta1.QueryOwnerAccountsResponse)
    - [QueryOwnerDomainsRequest](#starnamed.x.starname.v1beta1.QueryOwnerDomainsRequest)
    - [QueryOwnerDomainsResponse](#starnamed.x.starname.v1beta1.QueryOwnerDomainsResponse)
    - [QueryResourceAccountsRequest](#starnamed.x.starname.v1beta1.QueryResourceAccountsRequest)
    - [QueryResourceAccountsResponse](#starnamed.x.starname.v1beta1.QueryResourceAccountsResponse)
    - [QueryStarnameRequest](#starnamed.x.starname.v1beta1.QueryStarnameRequest)
    - [QueryStarnameResponse](#starnamed.x.starname.v1beta1.QueryStarnameResponse)
    - [QueryYieldRequest](#starnamed.x.starname.v1beta1.QueryYieldRequest)
    - [QueryYieldResponse](#starnamed.x.starname.v1beta1.QueryYieldResponse)
  
    - [Query](#starnamed.x.starname.v1beta1.Query)
  
- [iov/starname/v1beta1/tx.proto](#iov/starname/v1beta1/tx.proto)
    - [MsgAddAccountCertificate](#starnamed.x.starname.v1beta1.MsgAddAccountCertificate)
    - [MsgAddAccountCertificateResponse](#starnamed.x.starname.v1beta1.MsgAddAccountCertificateResponse)
    - [MsgDeleteAccount](#starnamed.x.starname.v1beta1.MsgDeleteAccount)
    - [MsgDeleteAccountCertificate](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificate)
    - [MsgDeleteAccountCertificateResponse](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificateResponse)
    - [MsgDeleteAccountResponse](#starnamed.x.starname.v1beta1.MsgDeleteAccountResponse)
    - [MsgDeleteDomain](#starnamed.x.starname.v1beta1.MsgDeleteDomain)
    - [MsgDeleteDomainResponse](#starnamed.x.starname.v1beta1.MsgDeleteDomainResponse)
    - [MsgRegisterAccount](#starnamed.x.starname.v1beta1.MsgRegisterAccount)
    - [MsgRegisterAccountResponse](#starnamed.x.starname.v1beta1.MsgRegisterAccountResponse)
    - [MsgRegisterDomain](#starnamed.x.starname.v1beta1.MsgRegisterDomain)
    - [MsgRegisterDomainResponse](#starnamed.x.starname.v1beta1.MsgRegisterDomainResponse)
    - [MsgRenewAccount](#starnamed.x.starname.v1beta1.MsgRenewAccount)
    - [MsgRenewAccountResponse](#starnamed.x.starname.v1beta1.MsgRenewAccountResponse)
    - [MsgRenewDomain](#starnamed.x.starname.v1beta1.MsgRenewDomain)
    - [MsgRenewDomainResponse](#starnamed.x.starname.v1beta1.MsgRenewDomainResponse)
    - [MsgReplaceAccountMetadata](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadata)
    - [MsgReplaceAccountMetadataResponse](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadataResponse)
    - [MsgReplaceAccountResources](#starnamed.x.starname.v1beta1.MsgReplaceAccountResources)
    - [MsgReplaceAccountResourcesResponse](#starnamed.x.starname.v1beta1.MsgReplaceAccountResourcesResponse)
    - [MsgTransferAccount](#starnamed.x.starname.v1beta1.MsgTransferAccount)
    - [MsgTransferAccountResponse](#starnamed.x.starname.v1beta1.MsgTransferAccountResponse)
    - [MsgTransferDomain](#starnamed.x.starname.v1beta1.MsgTransferDomain)
    - [MsgTransferDomainResponse](#starnamed.x.starname.v1beta1.MsgTransferDomainResponse)
  
    - [Msg](#starnamed.x.starname.v1beta1.Msg)
  
- [wasm/types.proto](#wasm/types.proto)
    - [AccessConfig](#iovone.starnamed.wasm.AccessConfig)
    - [AccessTypeParam](#iovone.starnamed.wasm.AccessTypeParam)
    - [Params](#iovone.starnamed.wasm.Params)
  
    - [AccessType](#iovone.starnamed.wasm.AccessType)
  
- [wasm/genesis.proto](#wasm/genesis.proto)
    - [GenesisState](#iovone.starnamed.wasm.GenesisState)
  
- [wasm/query.proto](#wasm/query.proto)
    - [Query](#iovone.starnamed.wasm.Query)
  
- [wasm/tx.proto](#wasm/tx.proto)
    - [Msg](#iovone.starnamed.wasm.Msg)
=======
- [iov/offchain/v1alpha1/offchain.proto](#iov/offchain/v1alpha1/offchain.proto)
    - [ListOfMsgSignData](#cosmos.offchain.v1alpha1.ListOfMsgSignData)
    - [MsgSignData](#cosmos.offchain.v1alpha1.MsgSignData)
>>>>>>> tags/v0.11.6
  
- [Scalar Value Types](#scalar-value-types)



<<<<<<< HEAD
<a name="iov/configuration/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/configuration/v1beta1/types.proto



<a name="starnamed.x.configuration.v1beta1.Config"></a>

### Config
Config is the configuration of the network
=======
<a name="iov/offchain/v1alpha1/offchain.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/offchain/v1alpha1/offchain.proto



<a name="cosmos.offchain.v1alpha1.ListOfMsgSignData"></a>

### ListOfMsgSignData
ListOfMsgSignData defines a list of MsgSignData, used to marshal and
unmarshal them in a clean way
>>>>>>> tags/v0.11.6


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `configurer` | [string](#string) |  | Configurer is the configuration owner, the addresses allowed to handle fees and register domains with no superuser |
| `valid_domain_name` | [string](#string) |  | ValidDomainName defines a regexp that determines if a domain name is valid or not |
| `valid_account_name` | [string](#string) |  | ValidAccountName defines a regexp that determines if an account name is valid or not |
| `valid_uri` | [string](#string) |  | ValidURI defines a regexp that determines if resource uri is valid or not |
| `valid_resource` | [string](#string) |  | ValidResource determines a regexp for a resource content |
| `domain_renewal_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainRenewalPeriod defines the duration of the domain renewal period in seconds |
| `domain_renewal_count_max` | [uint32](#uint32) |  | DomainRenewalCountMax defines maximum number of domain renewals a user can do |
| `domain_grace_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainGracePeriod defines the grace period for a domain deletion in seconds |
| `account_renewal_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | AccountRenewalPeriod defines the duration of the account renewal period in seconds |
| `account_renewal_count_max` | [uint32](#uint32) |  | AccountRenewalCountMax defines maximum number of account renewals a user can do |
| `account_grace_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainGracePeriod defines the grace period for a domain deletion in seconds |
| `resources_max` | [uint32](#uint32) |  | ResourcesMax defines maximum number of resources could be saved under an account |
| `certificate_size_max` | [uint64](#uint64) |  | CertificateSizeMax defines maximum size of a certificate that could be saved under an account |
| `certificate_count_max` | [uint32](#uint32) |  | CertificateCountMax defines maximum number of certificates that could be saved under an account |
| `metadata_size_max` | [uint64](#uint64) |  | MetadataSizeMax defines maximum size of metadata that could be saved under an account |
| `escrow_broker` | [string](#string) |  | EscrowBroker defines an address that will receive a commission for completed escrows |
| `escrow_commission` | [string](#string) |  | EscrowCommission defines the commission taken by the broker for a completed escrow, between 0 (no commission) and 1 (100% commission) |
| `escrow_max_period` | [google.protobuf.Duration](#google.protobuf.Duration) |  | EscrowPeriod defines the maximum duration of an escrow in seconds |
=======
| `msgs` | [MsgSignData](#cosmos.offchain.v1alpha1.MsgSignData) | repeated | msgs is a list of messages |
>>>>>>> tags/v0.11.6






<<<<<<< HEAD
<a name="starnamed.x.configuration.v1beta1.Fees"></a>

### Fees
Fees contains different type of fees to calculate coins to detract when processing different messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fee_coin_denom` | [string](#string) |  | FeeCoinDenom defines the denominator of the coin used to process fees |
| `fee_coin_price` | [string](#string) |  | FeeCoinPrice defines the price of the coin |
| `fee_default` | [string](#string) |  | FeeDefault is the parameter defining the default fee |
| `register_account_closed` | [string](#string) |  | RegisterAccountClosed is the fee to be paid to register an account in a closed domain |
| `register_account_open` | [string](#string) |  | RegisterAccountOpen is the fee to be paid to register an account in an open domain |
| `transfer_account_closed` | [string](#string) |  | TransferAccountClosed is the fee to be paid to register an account in a closed domain |
| `transfer_account_open` | [string](#string) |  | TransferAccountOpen is the fee to be paid to register an account in an open domain |
| `replace_account_resources` | [string](#string) |  | ReplaceAccountResources is the fee to be paid to replace account's resources |
| `add_account_certificate` | [string](#string) |  | AddAccountCertificate is the fee to be paid to add a certificate to an account |
| `del_account_certificate` | [string](#string) |  | DelAccountCertificate is the feed to be paid to delete a certificate in an account |
| `set_account_metadata` | [string](#string) |  | SetAccountMetadata is the fee to be paid to set account's metadata |
| `register_domain_1` | [string](#string) |  | RegisterDomain1 is the fee to be paid to register a domain with one character |
| `register_domain_2` | [string](#string) |  | RegisterDomain2 is the fee to be paid to register a domain with two characters |
| `register_domain_3` | [string](#string) |  | RegisterDomain3 is the fee to be paid to register a domain with three characters |
| `register_domain_4` | [string](#string) |  | RegisterDomain4 is the fee to be paid to register a domain with four characters |
| `register_domain_5` | [string](#string) |  | RegisterDomain5 is the fee to be paid to register a domain with five characters |
| `register_domain_default` | [string](#string) |  | RegisterDomainDefault is the fee to be paid to register a domain with more than five characters |
| `register_open_domain_multiplier` | [string](#string) |  | register_open_domain_multiplier is the multiplication applied to fees in register domain operations if they're of open type |
| `transfer_domain_closed` | [string](#string) |  | transfer_domain_closed is the fee to be paid to transfer a closed domain |
| `transfer_domain_open` | [string](#string) |  | transfer_domain_open is the fee to be paid to transfer open domains |
| `renew_domain_open` | [string](#string) |  | renew_domain_open is the fee to be paid to renew an open domain |
| `create_escrow` | [string](#string) |  | create_escrow is the fee to be paid to create an escrow |
| `update_escrow` | [string](#string) |  | update_escrow is the fee to be paid to update an escrow |
| `transfer_to_escrow` | [string](#string) |  | transfer_to_escrow is the fee to be paid to transfer coins to an escrow |
| `refund_escrow` | [string](#string) |  | refund_escrow is the fee to be paid to refund the account or domain placed in an escrow |






<a name="starnamed.x.configuration.v1beta1.GenesisState"></a>

### GenesisState
GenesisState - genesis state of x/configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [Config](#starnamed.x.configuration.v1beta1.Config) |  |  |
| `fees` | [Fees](#starnamed.x.configuration.v1beta1.Fees) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/configuration/v1beta1/msgs.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/configuration/v1beta1/msgs.proto



<a name="starnamed.x.configuration.v1beta1.MsgUpdateConfig"></a>

### MsgUpdateConfig
MsgUpdateConfig is used to update starname configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  | Signer is the address of the entity who is doing the transaction |
| `new_configuration` | [Config](#starnamed.x.configuration.v1beta1.Config) |  | NewConfiguration contains the new configuration data |






<a name="starnamed.x.configuration.v1beta1.MsgUpdateFees"></a>

### MsgUpdateFees
MsgUpdateFees is used to update the starname product fees in the starname module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fees` | [Fees](#starnamed.x.configuration.v1beta1.Fees) |  |  |
| `configurer` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/configuration/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/configuration/v1beta1/query.proto



<a name="starnamed.x.configuration.v1beta1.QueryConfigRequest"></a>

### QueryConfigRequest
QueryConfigRequest is the request type for the Query/Configuration RPC method.






<a name="starnamed.x.configuration.v1beta1.QueryConfigResponse"></a>

### QueryConfigResponse
QueryConfigResponse is the response type for the Query/Configuration RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `config` | [Config](#starnamed.x.configuration.v1beta1.Config) |  | Configuration is the starname configuration. |






<a name="starnamed.x.configuration.v1beta1.QueryFeesRequest"></a>

### QueryFeesRequest
QueryFeesRequest is the request type for the Query/Configuration RPC method.






<a name="starnamed.x.configuration.v1beta1.QueryFeesResponse"></a>

### QueryFeesResponse
QueryFeesResponse is the response type for the Query/Fees RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `fees` | [Fees](#starnamed.x.configuration.v1beta1.Fees) |  | Fees is the starname product fee object. |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="starnamed.x.configuration.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Config` | [QueryConfigRequest](#starnamed.x.configuration.v1beta1.QueryConfigRequest) | [QueryConfigResponse](#starnamed.x.configuration.v1beta1.QueryConfigResponse) | Config gets starname configuration. | GET|/starname/v1beta1/configuration/params|
| `Fees` | [QueryFeesRequest](#starnamed.x.configuration.v1beta1.QueryFeesRequest) | [QueryFeesResponse](#starnamed.x.configuration.v1beta1.QueryFeesResponse) | Fees gets starname product fees. | GET|/starname/v1beta1/configuration/fees|

 <!-- end services -->



<a name="iov/escrow/v1beta1/events.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/events.proto



<a name="starnamed.x.escrow.v1beta1.EventCompletedEscrow"></a>

### EventCompletedEscrow
EventCompletedEscrow is emitted when an escrow is completed


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `buyer` | [string](#string) |  |  |
| `fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="starnamed.x.escrow.v1beta1.EventCreatedEscrow"></a>

### EventCreatedEscrow
EventCreatedEscrow is emitted when an escrow is created


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `seller` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `broker_address` | [string](#string) |  |  |
| `broker_commission` | [string](#string) |  |  |
| `price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `object` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `deadline` | [uint64](#uint64) |  |  |
| `fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="starnamed.x.escrow.v1beta1.EventRefundedEscrow"></a>

### EventRefundedEscrow
EventRefundedEscrow is emitted when an escrow is refunded


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="starnamed.x.escrow.v1beta1.EventUpdatedEscrow"></a>

### EventUpdatedEscrow
EventUpdatedEscrow is emitted when an escrow is updated


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `updater` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `new_seller` | [string](#string) |  |  |
| `new_price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `new_deadline` | [uint64](#uint64) |  |  |
| `fees` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/escrow/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/types.proto



<a name="starnamed.x.escrow.v1beta1.Escrow"></a>

### Escrow
Escrow defines the struct of an escrow


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `seller` | [string](#string) |  |  |
| `object` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | TODO: refactor this to use sdk.Coin instead of sdk.Coins Although the price contains multiple coins, for now we enforce a specific denomination, so there will be only one coin type in a valid escrow |
| `state` | [EscrowState](#starnamed.x.escrow.v1beta1.EscrowState) |  |  |
| `deadline` | [uint64](#uint64) |  |  |
| `broker_address` | [string](#string) |  |  |
| `broker_commission` | [string](#string) |  |  |





 <!-- end messages -->


<a name="starnamed.x.escrow.v1beta1.EscrowState"></a>

### EscrowState
EscrowState defines the state of an escrow

| Name | Number | Description |
| ---- | ------ | ----------- |
| ESCROW_STATE_OPEN | 0 | ESCROW_STATE_OPEN defines an open state. |
| ESCROW_STATE_COMPLETED | 1 | ESCROW_STATE_COMPLETED defines a completed state. |
| ESCROW_STATE_REFUNDED | 2 | ESCROW_STATE_REFUNDED defines a refunded state. |
| ESCROW_STATE_EXPIRED | 3 | ESCROW_STATE_REFUNDED defines an expired state. |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/escrow/v1beta1/params.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/params.proto



<a name="starnamed.x.escrow.v1beta1.Params"></a>

### Params
Params defines the parameters of the escrow module


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `module_enabled` | [bool](#bool) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/escrow/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/genesis.proto



<a name="starnamed.x.escrow.v1beta1.GenesisState"></a>

### GenesisState
GenesisState defines the Escrow module's genesis state


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `escrows` | [Escrow](#starnamed.x.escrow.v1beta1.Escrow) | repeated |  |
| `last_block_time` | [uint64](#uint64) |  |  |
| `next_escrow_id` | [uint64](#uint64) |  |  |
| `params` | [Params](#starnamed.x.escrow.v1beta1.Params) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/escrow/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/query.proto



<a name="starnamed.x.escrow.v1beta1.QueryEscrowRequest"></a>

### QueryEscrowRequest
QueryEscrowRequest is the request type for the Query/Escrow RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="starnamed.x.escrow.v1beta1.QueryEscrowResponse"></a>

### QueryEscrowResponse
QueryEscrowResponse is the response type for the Query/Escrow RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `escrow` | [Escrow](#starnamed.x.escrow.v1beta1.Escrow) |  |  |






<a name="starnamed.x.escrow.v1beta1.QueryEscrowsRequest"></a>

### QueryEscrowsRequest
QueryEscrowsRequest is the request type for the Query/Escrows RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seller` | [string](#string) |  |  |
| `state` | [string](#string) |  |  |
| `object_key` | [string](#string) |  |  |
| `pagination_start` | [uint64](#uint64) |  |  |
| `pagination_length` | [uint64](#uint64) |  |  |






<a name="starnamed.x.escrow.v1beta1.QueryEscrowsResponse"></a>

### QueryEscrowsResponse
QueryEscrowsResponse is the response type for the Query/Escrows RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `escrows` | [Escrow](#starnamed.x.escrow.v1beta1.Escrow) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="starnamed.x.escrow.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Escrow` | [QueryEscrowRequest](#starnamed.x.escrow.v1beta1.QueryEscrowRequest) | [QueryEscrowResponse](#starnamed.x.escrow.v1beta1.QueryEscrowResponse) | Escrow queries the escrow by the specified id | GET|/starnamed/v1beta1/escrow/{id}|
| `Escrows` | [QueryEscrowsRequest](#starnamed.x.escrow.v1beta1.QueryEscrowsRequest) | [QueryEscrowsResponse](#starnamed.x.escrow.v1beta1.QueryEscrowsResponse) | Escrows queries escrows by the specified key-value pairs | GET|/starnamed/v1beta1/escrows|

 <!-- end services -->



<a name="iov/escrow/v1beta1/test.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/test.proto



<a name="starnamed.x.escrow.v1beta1.TestObject"></a>

### TestObject
TestObject defines a transferable object used for testing


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `owner` | [bytes](#bytes) |  |  |
| `num_allowed_transfers` | [int64](#int64) |  |  |






<a name="starnamed.x.escrow.v1beta1.TestTimeConstrainedObject"></a>

### TestTimeConstrainedObject
TestTimeConstrainedObject defines a transferable object with a time constrain used for testing


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [uint64](#uint64) |  |  |
| `owner` | [bytes](#bytes) |  |  |
| `expiration` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/escrow/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/escrow/v1beta1/tx.proto



<a name="starnamed.x.escrow.v1beta1.MsgCreateEscrow"></a>

### MsgCreateEscrow
MsgCreateEscrow defines a message to create an escrow


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `seller` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `object` | [google.protobuf.Any](#google.protobuf.Any) |  |  |
| `price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `deadline` | [uint64](#uint64) |  |  |






<a name="starnamed.x.escrow.v1beta1.MsgCreateEscrowResponse"></a>

### MsgCreateEscrowResponse
MsgCreateEscrowResponse defines the Msg/CreateEscrow response type


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |






<a name="starnamed.x.escrow.v1beta1.MsgRefundEscrow"></a>

### MsgRefundEscrow



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |






<a name="starnamed.x.escrow.v1beta1.MsgRefundEscrowResponse"></a>

### MsgRefundEscrowResponse







<a name="starnamed.x.escrow.v1beta1.MsgTransferToEscrow"></a>

### MsgTransferToEscrow



| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `sender` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `amount` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |






<a name="starnamed.x.escrow.v1beta1.MsgTransferToEscrowResponse"></a>

### MsgTransferToEscrowResponse







<a name="starnamed.x.escrow.v1beta1.MsgUpdateEscrow"></a>

### MsgUpdateEscrow
MsgUpdateEscrow defines a message to update an escrow


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id` | [string](#string) |  |  |
| `updater` | [string](#string) |  |  |
| `fee_payer` | [string](#string) |  |  |
| `seller` | [string](#string) |  |  |
| `price` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated |  |
| `deadline` | [uint64](#uint64) |  |  |






<a name="starnamed.x.escrow.v1beta1.MsgUpdateEscrowResponse"></a>

### MsgUpdateEscrowResponse
MsgUpdateEscrowResponse defines the Msg/UpdateEscrow response type





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="starnamed.x.escrow.v1beta1.Msg"></a>

### Msg
Msg defines the escrow Msg service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `CreateEscrow` | [MsgCreateEscrow](#starnamed.x.escrow.v1beta1.MsgCreateEscrow) | [MsgCreateEscrowResponse](#starnamed.x.escrow.v1beta1.MsgCreateEscrowResponse) | CreateEscrow defines a method for creating an escrow | |
| `UpdateEscrow` | [MsgUpdateEscrow](#starnamed.x.escrow.v1beta1.MsgUpdateEscrow) | [MsgUpdateEscrowResponse](#starnamed.x.escrow.v1beta1.MsgUpdateEscrowResponse) | UpdateEscrow defines a method for updating an escrow | |
| `TransferToEscrow` | [MsgTransferToEscrow](#starnamed.x.escrow.v1beta1.MsgTransferToEscrow) | [MsgTransferToEscrowResponse](#starnamed.x.escrow.v1beta1.MsgTransferToEscrowResponse) | TransferToEscrow defines a method for a buyer to transfer funds to the escrow | |
| `RefundEscrow` | [MsgRefundEscrow](#starnamed.x.escrow.v1beta1.MsgRefundEscrow) | [MsgRefundEscrowResponse](#starnamed.x.escrow.v1beta1.MsgRefundEscrowResponse) | RefundEscrow defines a method for the seller to return the assets locked in the escrow | |

 <!-- end services -->



<a name="iov/offchain/v1alpha1/offchain.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/offchain/v1alpha1/offchain.proto



<a name="cosmos.offchain.v1alpha1.ListOfMsgSignData"></a>

### ListOfMsgSignData
ListOfMsgSignData defines a list of MsgSignData, used to marshal and unmarshal them in a clean way


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `msgs` | [MsgSignData](#cosmos.offchain.v1alpha1.MsgSignData) | repeated | msgs is a list of messages |






<a name="cosmos.offchain.v1alpha1.MsgSignData"></a>

### MsgSignData
MsgSignData defines an arbitrary, general-purpose, off-chain message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `signer` | [string](#string) |  | signer is the bech32 representation of the signer's account address |
| `data` | [bytes](#bytes) |  | data represents the raw bytes of the content that is signed (text, json, etc) |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/starname/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/starname/v1beta1/types.proto



<a name="starnamed.x.starname.v1beta1.Account"></a>

### Account
Account defines an account that belongs to a domain
NOTE: It should not be confused with cosmos-sdk auth account
github.com/cosmos/cosmos-sdk/x/auth.Account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain references the domain this account belongs to |
| `name` | [google.protobuf.StringValue](#google.protobuf.StringValue) |  | Name is the name of the account |
| `owner` | [bytes](#bytes) |  | Owner is the address that owns the account |
| `broker` | [bytes](#bytes) |  | Broker identifies an entity that facilitated the transaction of the account and can be empty |
| `valid_until` | [int64](#int64) |  | ValidUntil defines a unix timestamp of the expiration of the account in seconds |
| `resources` | [Resource](#starnamed.x.starname.v1beta1.Resource) | repeated | Resources is the list of resources an account resolves to |
| `certificates` | [bytes](#bytes) | repeated | Certificates contains the list of certificates to identify the account owner |
| `metadata_uri` | [string](#string) |  | MetadataURI contains a link to extra information regarding the account |






<a name="starnamed.x.starname.v1beta1.Domain"></a>

### Domain
Domain defines a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Name is the name of the domain |
| `admin` | [bytes](#bytes) |  | Admin is the owner of the domain |
| `broker` | [bytes](#bytes) |  |  |
| `valid_until` | [int64](#int64) |  | ValidUntil is a unix timestamp defines the time when the domain will become invalid in seconds |
| `type` | [string](#string) |  | Type defines the type of the domain |






<a name="starnamed.x.starname.v1beta1.Resource"></a>

### Resource
Resource defines a resource owned by an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `uri` | [string](#string) |  | URI defines the ID of the resource |
| `resource` | [string](#string) |  | Resource is the resource |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/starname/v1beta1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/starname/v1beta1/genesis.proto



<a name="starnamed.x.starname.v1beta1.GenesisState"></a>

### GenesisState
GenesisState - genesis state of x/starname


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domains` | [Domain](#starnamed.x.starname.v1beta1.Domain) | repeated |  |
| `accounts` | [Account](#starnamed.x.starname.v1beta1.Account) | repeated |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="iov/starname/v1beta1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/starname/v1beta1/query.proto



<a name="starnamed.x.starname.v1beta1.QueryBrokerAccountsRequest"></a>

### QueryBrokerAccountsRequest
QueryBrokerAccountsRequest is the request type for the Query/BrokerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `broker` | [string](#string) |  | Broker is the broker of accounts. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryBrokerAccountsResponse"></a>

### QueryBrokerAccountsResponse
QueryBrokerAccountsResponse is the response type for the Query/BrokerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts is the accounts associated with broker. |
| `page` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryBrokerDomainsRequest"></a>

### QueryBrokerDomainsRequest
QueryBrokerDomainsRequest is the request type for the Query/BrokerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `broker` | [string](#string) |  | Broker is the broker of accounts. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryBrokerDomainsResponse"></a>

### QueryBrokerDomainsResponse
QueryBrokerDomainsResponse is the response type for the Query/BrokerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domains` | [Domain](#starnamed.x.starname.v1beta1.Domain) | repeated | Accounts is the accounts associated with broker. |
| `page` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryDomainAccountsRequest"></a>

### QueryDomainAccountsRequest
QueryDomainAccountsRequest is the request type for the Query/DomainAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the name of the domain. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryDomainAccountsResponse"></a>

### QueryDomainAccountsResponse
QueryDomainAccountsResponse is the response type for the Query/DomainAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts is the accounts associated with the domain. |
| `page` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryDomainRequest"></a>

### QueryDomainRequest
QueryDomainRequest is the request type for the Query/Domain RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  | Name is the name of the domain. |






<a name="starnamed.x.starname.v1beta1.QueryDomainResponse"></a>

### QueryDomainResponse
QueryDomainResponse is the response type for the Query/Domain RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [Domain](#starnamed.x.starname.v1beta1.Domain) |  | Domain is the information associated with the domain. |






<a name="starnamed.x.starname.v1beta1.QueryOwnerAccountsRequest"></a>

### QueryOwnerAccountsRequest
QueryOwnerAccountsRequest is the request type for the Query/OwnerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  | Owner is the owner of accounts. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryOwnerAccountsResponse"></a>

### QueryOwnerAccountsResponse
QueryOwnerAccountsResponse is the response type for the Query/OwnerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts is the accounts associated with owner. |
| `page` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryOwnerDomainsRequest"></a>

### QueryOwnerDomainsRequest
QueryOwnerDomainsRequest is the request type for the Query/OwnerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `owner` | [string](#string) |  | Owner is the owner of accounts. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryOwnerDomainsResponse"></a>

### QueryOwnerDomainsResponse
QueryOwnerDomainsResponse is the response type for the Query/OwnerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domains` | [Domain](#starnamed.x.starname.v1beta1.Domain) | repeated | Accounts is the accounts associated with owner. |
| `page` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryResourceAccountsRequest"></a>

### QueryResourceAccountsRequest
QueryResourceAccountsRequest is the request type for the Query/ResourceAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `uri` | [string](#string) |  | Uri is the uri of the resource. query.pb.gw.to doesn't respect gogoproto.customname, so we're stuck with Uri. |
| `resource` | [string](#string) |  | Resource is the resource of interest. |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryResourceAccountsResponse"></a>

### QueryResourceAccountsResponse
QueryResourceAccountsResponse is the response type for the Query/ResourceAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `accounts` | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts are the accounts associated with the resource. |
| `page` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryStarnameRequest"></a>

### QueryStarnameRequest
QueryStarnameRequest is the request type for the Query/Starname RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `starname` | [string](#string) |  | Starname is the of the form account*domain. |






<a name="starnamed.x.starname.v1beta1.QueryStarnameResponse"></a>

### QueryStarnameResponse
QueryStarnameResponse is the response type for the Query/Starname RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `account` | [Account](#starnamed.x.starname.v1beta1.Account) |  | Account is the information associated with the starname. |






<a name="starnamed.x.starname.v1beta1.QueryYieldRequest"></a>

### QueryYieldRequest
QueryYieldRequest is the request type for the Query/Yield RPC method.






<a name="starnamed.x.starname.v1beta1.QueryYieldResponse"></a>

### QueryYieldResponse
QueryYieldResponse is the response type for the Query/Yield RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `yield` | [string](#string) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="starnamed.x.starname.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `Domain` | [QueryDomainRequest](#starnamed.x.starname.v1beta1.QueryDomainRequest) | [QueryDomainResponse](#starnamed.x.starname.v1beta1.QueryDomainResponse) | Domain gets a starname's domain info. | GET|/starname/v1beta1/domain/{name}|
| `DomainAccounts` | [QueryDomainAccountsRequest](#starnamed.x.starname.v1beta1.QueryDomainAccountsRequest) | [QueryDomainAccountsResponse](#starnamed.x.starname.v1beta1.QueryDomainAccountsResponse) | DomainAccounts gets accounts associated with a given domain. | GET|/starname/v1beta1/accounts/domain/{domain}|
| `Starname` | [QueryStarnameRequest](#starnamed.x.starname.v1beta1.QueryStarnameRequest) | [QueryStarnameResponse](#starnamed.x.starname.v1beta1.QueryStarnameResponse) | Starname gets all the information associated with a starname. | GET|/starname/v1beta1/account/{starname}|
| `OwnerAccounts` | [QueryOwnerAccountsRequest](#starnamed.x.starname.v1beta1.QueryOwnerAccountsRequest) | [QueryOwnerAccountsResponse](#starnamed.x.starname.v1beta1.QueryOwnerAccountsResponse) | OwnerAccounts gets accounts associated with a given owner. | GET|/starname/v1beta1/accounts/owner/{owner}|
| `OwnerDomains` | [QueryOwnerDomainsRequest](#starnamed.x.starname.v1beta1.QueryOwnerDomainsRequest) | [QueryOwnerDomainsResponse](#starnamed.x.starname.v1beta1.QueryOwnerDomainsResponse) | OwnerDomains gets domains associated with a given owner. | GET|/starname/v1beta1/domains/owner/{owner}|
| `ResourceAccounts` | [QueryResourceAccountsRequest](#starnamed.x.starname.v1beta1.QueryResourceAccountsRequest) | [QueryResourceAccountsResponse](#starnamed.x.starname.v1beta1.QueryResourceAccountsResponse) | ResourceAccounts gets accounts associated with a given resource. | GET|/starname/v1beta1/accounts/resource/{uri}/{resource}|
| `BrokerAccounts` | [QueryBrokerAccountsRequest](#starnamed.x.starname.v1beta1.QueryBrokerAccountsRequest) | [QueryBrokerAccountsResponse](#starnamed.x.starname.v1beta1.QueryBrokerAccountsResponse) | BrokerAccounts gets accounts associated with a given broker. | GET|/starname/v1beta1/accounts/broker/{broker}|
| `BrokerDomains` | [QueryBrokerDomainsRequest](#starnamed.x.starname.v1beta1.QueryBrokerDomainsRequest) | [QueryBrokerDomainsResponse](#starnamed.x.starname.v1beta1.QueryBrokerDomainsResponse) | BrokerDomains gets domains associated with a given broker. | GET|/starname/v1beta1/domains/broker/{broker}|
| `Yield` | [QueryYieldRequest](#starnamed.x.starname.v1beta1.QueryYieldRequest) | [QueryYieldResponse](#starnamed.x.starname.v1beta1.QueryYieldResponse) | Yield estimates and retrieves the annualized yield for delegators | GET|/starname/v1beta1/yield|

 <!-- end services -->



<a name="iov/starname/v1beta1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/starname/v1beta1/tx.proto



<a name="starnamed.x.starname.v1beta1.MsgAddAccountCertificate"></a>

### MsgAddAccountCertificate
MsgAddAccountCertificate is the message used when a user wants to add new certificates to his account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `new_certificate` | [bytes](#bytes) |  | NewCertificate is the new certificate to add |






<a name="starnamed.x.starname.v1beta1.MsgAddAccountCertificateResponse"></a>

### MsgAddAccountCertificateResponse
MsgAddAccountCertificateResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgDeleteAccount"></a>

### MsgDeleteAccount
MsgDeleteAccount is the request model used to delete an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgDeleteAccountCertificate"></a>

### MsgDeleteAccountCertificate
MsgDeleteAccountCertificate is the request model used to remove certificates from an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `delete_certificate` | [bytes](#bytes) |  | DeleteCertificate is the certificate to delete |






<a name="starnamed.x.starname.v1beta1.MsgDeleteAccountCertificateResponse"></a>

### MsgDeleteAccountCertificateResponse
MsgDeleteAccountCertificateResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgDeleteAccountResponse"></a>

### MsgDeleteAccountResponse
MsgDeleteAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgDeleteDomain"></a>

### MsgDeleteDomain
MsgDeleteDomain is the request model to delete a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgDeleteDomainResponse"></a>

### MsgDeleteDomainResponse
MsgDeleteDomainResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRegisterAccount"></a>

### MsgRegisterAccount
MsgRegisterAccount is the request model used to register new accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `broker` | [string](#string) |  | Broker is the account that facilitated the transaction |
| `registerer` | [string](#string) |  | Registerer is the user who registers this account |
| `resources` | [Resource](#starnamed.x.starname.v1beta1.Resource) | repeated | Resources are the blockchain addresses of the account |






<a name="starnamed.x.starname.v1beta1.MsgRegisterAccountResponse"></a>

### MsgRegisterAccountResponse
MsgRegisterAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRegisterDomain"></a>

### MsgRegisterDomain
MsgRegisterDomain is the request used to register new domains


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `name` | [string](#string) |  |  |
| `admin` | [string](#string) |  |  |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `broker` | [string](#string) |  | Broker is the account that facilitated the transaction |
| `domain_type` | [string](#string) |  | DomainType defines the type of the domain |






<a name="starnamed.x.starname.v1beta1.MsgRegisterDomainResponse"></a>

### MsgRegisterDomainResponse
MsgRegisterDomainResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRenewAccount"></a>

### MsgRenewAccount
MsgRenewAccount is the request model used to renew accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `signer` | [string](#string) |  | Signer is the signer of the request |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgRenewAccountResponse"></a>

### MsgRenewAccountResponse
MsgRenewAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRenewDomain"></a>

### MsgRenewDomain
MsgRenewDomain is the request model used to renew a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `signer` | [string](#string) |  | Signer is the signer of the request |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgRenewDomainResponse"></a>

### MsgRenewDomainResponse
MsgRegisterDomain returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountMetadata"></a>

### MsgReplaceAccountMetadata
MsgReplaceAccountMetadata is the function used to set accounts metadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `new_metadata_uri` | [string](#string) |  | NewMetadataURI is the metadata URI of the account we want to update or insert |






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountMetadataResponse"></a>

### MsgReplaceAccountMetadataResponse
MsgReplaceAccountMetadataResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountResources"></a>

### MsgReplaceAccountResources
MsgReplaceAccountResources is the request model used to renew resources associated with an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `new_resources` | [Resource](#starnamed.x.starname.v1beta1.Resource) | repeated | NewResources are the new resources |






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountResourcesResponse"></a>

### MsgReplaceAccountResourcesResponse
MsgReplaceAccountResourcesResponse






<a name="starnamed.x.starname.v1beta1.MsgTransferAccount"></a>

### MsgTransferAccount
MsgTransferAccount is the request model used to transfer accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the domain of the account |
| `name` | [string](#string) |  | Name is the name of the account |
| `owner` | [string](#string) |  | Owner is the owner of the account |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `new_owner` | [string](#string) |  | NewOwner is the new owner of the account |
| `reset` | [bool](#bool) |  | ToReset if true, removes all old data from account |






<a name="starnamed.x.starname.v1beta1.MsgTransferAccountResponse"></a>

### MsgTransferAccountResponse
MsgTransferAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgTransferDomain"></a>

### MsgTransferDomain
MsgTransferDomain is the request model used to transfer a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `domain` | [string](#string) |  | Domain is the name of the domain |
| `owner` | [string](#string) |  | Owner is the owner of the domain |
| `payer` | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| `new_admin` | [string](#string) |  | NewAdmin is the new owner of the domain |
| `transfer_flag` | [int64](#int64) |  | TransferFlag controls the operations that occurs on a domain's accounts |






<a name="starnamed.x.starname.v1beta1.MsgTransferDomainResponse"></a>

### MsgTransferDomainResponse
MsgTransferDomainResponse returns an empty response.





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="starnamed.x.starname.v1beta1.Msg"></a>

### Msg
Msg defines the starname Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
| `AddAccountCertificate` | [MsgAddAccountCertificate](#starnamed.x.starname.v1beta1.MsgAddAccountCertificate) | [MsgAddAccountCertificateResponse](#starnamed.x.starname.v1beta1.MsgAddAccountCertificateResponse) | AddAccountCertificate adds a certificate to an Account | |
| `DeleteAccount` | [MsgDeleteAccount](#starnamed.x.starname.v1beta1.MsgDeleteAccount) | [MsgDeleteAccountResponse](#starnamed.x.starname.v1beta1.MsgDeleteAccountResponse) | DeleteAccount registers a Domain | |
| `DeleteAccountCertificate` | [MsgDeleteAccountCertificate](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificate) | [MsgDeleteAccountCertificateResponse](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificateResponse) | DeleteAccountCertificate deletes a certificate from an account | |
| `DeleteDomain` | [MsgDeleteDomain](#starnamed.x.starname.v1beta1.MsgDeleteDomain) | [MsgDeleteDomainResponse](#starnamed.x.starname.v1beta1.MsgDeleteDomainResponse) | DeleteDomain registers a Domain | |
| `RegisterAccount` | [MsgRegisterAccount](#starnamed.x.starname.v1beta1.MsgRegisterAccount) | [MsgRegisterAccountResponse](#starnamed.x.starname.v1beta1.MsgRegisterAccountResponse) | RegisterAccount registers an Account | |
| `RegisterDomain` | [MsgRegisterDomain](#starnamed.x.starname.v1beta1.MsgRegisterDomain) | [MsgRegisterDomainResponse](#starnamed.x.starname.v1beta1.MsgRegisterDomainResponse) | RegisterDomain registers a Domain | |
| `RenewAccount` | [MsgRenewAccount](#starnamed.x.starname.v1beta1.MsgRenewAccount) | [MsgRenewAccountResponse](#starnamed.x.starname.v1beta1.MsgRenewAccountResponse) | RenewAccount registers a Domain | |
| `RenewDomain` | [MsgRenewDomain](#starnamed.x.starname.v1beta1.MsgRenewDomain) | [MsgRenewDomainResponse](#starnamed.x.starname.v1beta1.MsgRenewDomainResponse) | RenewDomain registers a Domain | |
| `ReplaceAccountMetadata` | [MsgReplaceAccountMetadata](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadata) | [MsgReplaceAccountMetadataResponse](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadataResponse) | ReplaceAccountMetadata registers a Domain | |
| `ReplaceAccountResources` | [MsgReplaceAccountResources](#starnamed.x.starname.v1beta1.MsgReplaceAccountResources) | [MsgReplaceAccountResourcesResponse](#starnamed.x.starname.v1beta1.MsgReplaceAccountResourcesResponse) | ReplaceAccountResources registers a Domain | |
| `TransferAccount` | [MsgTransferAccount](#starnamed.x.starname.v1beta1.MsgTransferAccount) | [MsgTransferAccountResponse](#starnamed.x.starname.v1beta1.MsgTransferAccountResponse) | TransferAccount registers a Domain | |
| `TransferDomain` | [MsgTransferDomain](#starnamed.x.starname.v1beta1.MsgTransferDomain) | [MsgTransferDomainResponse](#starnamed.x.starname.v1beta1.MsgTransferDomainResponse) | TransferDomain registers a Domain | |

 <!-- end services -->



<a name="wasm/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## wasm/types.proto



<a name="iovone.starnamed.wasm.AccessConfig"></a>
=======
<a name="cosmos.offchain.v1alpha1.MsgSignData"></a>
>>>>>>> tags/v0.11.6

### MsgSignData
MsgSignData defines an arbitrary, general-purpose, off-chain message


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `permission` | [AccessType](#iovone.starnamed.wasm.AccessType) |  |  |
| `address` | [string](#string) |  |  |






<a name="iovone.starnamed.wasm.AccessTypeParam"></a>

### AccessTypeParam
AccessTypeParam


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `value` | [AccessType](#iovone.starnamed.wasm.AccessType) |  |  |






<a name="iovone.starnamed.wasm.Params"></a>

### Params
Params defines the set of wasm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_upload_access` | [AccessConfig](#iovone.starnamed.wasm.AccessConfig) |  |  |
| `instantiate_default_permission` | [AccessType](#iovone.starnamed.wasm.AccessType) |  |  |
| `max_wasm_code_size` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<a name="iovone.starnamed.wasm.AccessType"></a>

### AccessType
AccessType permission types

| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_TYPE_UNSPECIFIED | 0 | AccessTypeUnspecified placeholder for empty value |
| ACCESS_TYPE_NOBODY | 1 | AccessTypeNobody forbidden |
| ACCESS_TYPE_ONLY_ADDRESS | 2 | AccessTypeOnlyAddress restricted to an address |
| ACCESS_TYPE_EVERYBODY | 3 | AccessTypeEverybody unrestricted |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="wasm/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## wasm/genesis.proto



<a name="iovone.starnamed.wasm.GenesisState"></a>

### GenesisState
GenesisState defines the wasm module's genesis state.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `params` | [Params](#iovone.starnamed.wasm.Params) |  |  |
=======
| `signer` | [string](#string) |  | signer is the bech32 representation of the signer's account address |
| `data` | [bytes](#bytes) |  | data represents the raw bytes of the content that is signed (text, json, etc) |
>>>>>>> tags/v0.11.6





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

<<<<<<< HEAD
 <!-- end services -->



<a name="wasm/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## wasm/query.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="iovone.starnamed.wasm.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |

 <!-- end services -->



<a name="wasm/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## wasm/tx.proto


 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<a name="iovone.starnamed.wasm.Msg"></a>

### Msg
Msg defines the Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |

=======
>>>>>>> tags/v0.11.6
 <!-- end services -->



## Scalar Value Types

| .proto Type | Notes | C++ | Java | Python | Go | C# | PHP | Ruby |
| ----------- | ----- | --- | ---- | ------ | -- | -- | --- | ---- |
| <a name="double" /> double |  | double | double | float | float64 | double | float | Float |
| <a name="float" /> float |  | float | float | float | float32 | float | float | Float |
| <a name="int32" /> int32 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint32 instead. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="int64" /> int64 | Uses variable-length encoding. Inefficient for encoding negative numbers – if your field is likely to have negative values, use sint64 instead. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="uint32" /> uint32 | Uses variable-length encoding. | uint32 | int | int/long | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="uint64" /> uint64 | Uses variable-length encoding. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum or Fixnum (as required) |
| <a name="sint32" /> sint32 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int32s. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sint64" /> sint64 | Uses variable-length encoding. Signed int value. These more efficiently encode negative numbers than regular int64s. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="fixed32" /> fixed32 | Always four bytes. More efficient than uint32 if values are often greater than 2^28. | uint32 | int | int | uint32 | uint | integer | Bignum or Fixnum (as required) |
| <a name="fixed64" /> fixed64 | Always eight bytes. More efficient than uint64 if values are often greater than 2^56. | uint64 | long | int/long | uint64 | ulong | integer/string | Bignum |
| <a name="sfixed32" /> sfixed32 | Always four bytes. | int32 | int | int | int32 | int | integer | Bignum or Fixnum (as required) |
| <a name="sfixed64" /> sfixed64 | Always eight bytes. | int64 | long | int/long | int64 | long | integer/string | Bignum |
| <a name="bool" /> bool |  | bool | boolean | boolean | bool | bool | boolean | TrueClass/FalseClass |
| <a name="string" /> string | A string must always contain UTF-8 encoded or 7-bit ASCII text. | string | String | str/unicode | string | string | string | String (UTF-8) |
| <a name="bytes" /> bytes | May contain any arbitrary sequence of bytes. | string | ByteString | str | []byte | ByteString | string | String (ASCII-8BIT) |

