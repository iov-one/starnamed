# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [x/configuration/types/types.proto](#x/configuration/types/types.proto)
    - [Config](#wasmd.x.configuration.v1beta1.Config)
    - [Fees](#wasmd.x.configuration.v1beta1.Fees)
    - [GenesisState](#wasmd.x.configuration.v1beta1.GenesisState)
  
- [x/configuration/types/query.proto](#x/configuration/types/query.proto)
    - [QueryConfigRequest](#wasmd.x.configuration.v1beta1.QueryConfigRequest)
    - [QueryConfigResponse](#wasmd.x.configuration.v1beta1.QueryConfigResponse)
    - [QueryFeesRequest](#wasmd.x.configuration.v1beta1.QueryFeesRequest)
    - [QueryFeesResponse](#wasmd.x.configuration.v1beta1.QueryFeesResponse)
  
    - [Query](#wasmd.x.configuration.v1beta1.Query)
  
- [x/configuration/types/msgs.proto](#x/configuration/types/msgs.proto)
    - [MsgUpdateConfig](#wasmd.x.configuration.v1beta1.MsgUpdateConfig)
    - [MsgUpdateFees](#wasmd.x.configuration.v1beta1.MsgUpdateFees)
  
- [Scalar Value Types](#scalar-value-types)



<a name="x/configuration/types/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/configuration/types/types.proto



<a name="wasmd.x.configuration.v1beta1.Config"></a>

### Config
Config is the configuration of the network


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| configurer | [bytes](#bytes) |  | Configurer is the configuration owner, the addresses allowed to handle fees and register domains with no superuser |
| valid_domain_name | [string](#string) |  | ValidDomainName defines a regexp that determines if a domain name is valid or not |
| valid_account_name | [string](#string) |  | ValidAccountName defines a regexp that determines if an account name is valid or not |
| valid_uri | [string](#string) |  | ValidURI defines a regexp that determines if resource uri is valid or not |
| valid_resource | [string](#string) |  | ValidResource determines a regexp for a resource content |
| domain_renewal_period | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainRenewalPeriod defines the duration of the domain renewal period in seconds |
| domain_renewal_count_max | [uint32](#uint32) |  | DomainRenewalCountMax defines maximum number of domain renewals a user can do |
| domain_grace_period | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainGracePeriod defines the grace period for a domain deletion in seconds |
| account_renewal_period | [google.protobuf.Duration](#google.protobuf.Duration) |  | AccountRenewalPeriod defines the duration of the account renewal period in seconds |
| account_renewal_count_max | [uint32](#uint32) |  | AccountRenewalCountMax defines maximum number of account renewals a user can do |
| account_grace_period | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainGracePeriod defines the grace period for a domain deletion in seconds |
| resources_max | [uint32](#uint32) |  | ResourcesMax defines maximum number of resources could be saved under an account |
| certificate_size_max | [uint64](#uint64) |  | CertificateSizeMax defines maximum size of a certificate that could be saved under an account |
| certificate_count_max | [uint32](#uint32) |  | CertificateCountMax defines maximum number of certificates that could be saved under an account |
| metadata_size_max | [uint64](#uint64) |  | MetadataSizeMax defines maximum size of metadata that could be saved under an account |






<a name="wasmd.x.configuration.v1beta1.Fees"></a>

### Fees
Fees contains different type of fees to calculate coins to detract when processing different messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fee_coin_denom | [string](#string) |  | FeeCoinDenom defines the denominator of the coin used to process fees |
| fee_coin_price | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | FeeCoinPrice defines the price of the coin |
| fee_default | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | FeeDefault is the parameter defining the default fee |
| register_account_closed | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterAccountClosed is the fee to be paid to register an account in a closed domain |
| register_account_open | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterAccountOpen is the fee to be paid to register an account in an open domain |
| transfer_account_closed | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | TransferAccountClosed is the fee to be paid to register an account in a closed domain |
| transfer_account_open | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | TransferAccountOpen is the fee to be paid to register an account in an open domain |
| replace_account_resources | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | ReplaceAccountResources is the fee to be paid to replace account&#39;s resources |
| add_account_certificate | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | AddAccountCertificate is the fee to be paid to add a certificate to an account |
| del_account_certificate | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | DelAccountCertificate is the feed to be paid to delete a certificate in an account |
| set_account_metadata | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | SetAccountMetadata is the fee to be paid to set account&#39;s metadata |
| register_domain_1 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain1 is the fee to be paid to register a domain with one character |
| register_domain_2 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain2 is the fee to be paid to register a domain with two characters |
| register_domain_3 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain3 is the fee to be paid to register a domain with three characters |
| register_domain_4 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain4 is the fee to be paid to register a domain with four characters |
| register_domain_5 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain5 is the fee to be paid to register a domain with five characters |
| register_domain_default | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomainDefault is the fee to be paid to register a domain with more than five characters |
| register_open_domain_multiplier | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | register_open_domain_multiplier is the multiplication applied to fees in register domain operations if they&#39;re of open type |
| transfer_domain_closed | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | transfer_domain_closed is the fee to be paid to transfer a closed domain |
| transfer_domain_open | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | transfer_domain_open is the fee to be paid to transfer open domains |
| renew_domain_open | [string](#string) |  | renew_domain_open is the fee to be paid to renew an open domain |






<a name="wasmd.x.configuration.v1beta1.GenesisState"></a>

### GenesisState
GenesisState - genesis state of x/configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [Config](#wasmd.x.configuration.v1beta1.Config) |  |  |
| fees | [Fees](#wasmd.x.configuration.v1beta1.Fees) |  |  |





 

 

 

 



<a name="x/configuration/types/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/configuration/types/query.proto



<a name="wasmd.x.configuration.v1beta1.QueryConfigRequest"></a>

### QueryConfigRequest
QueryConfigRequest is the request type for the Query/Configuration RPC method.






<a name="wasmd.x.configuration.v1beta1.QueryConfigResponse"></a>

### QueryConfigResponse
QueryConfigResponse is the response type for the Query/Configuration RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| config | [Config](#wasmd.x.configuration.v1beta1.Config) |  | Configuration is the starname configuration. |






<a name="wasmd.x.configuration.v1beta1.QueryFeesRequest"></a>

### QueryFeesRequest
QueryFeesRequest is the request type for the Query/Configuration RPC method.






<a name="wasmd.x.configuration.v1beta1.QueryFeesResponse"></a>

### QueryFeesResponse
QueryFeesResponse is the response type for the Query/Fees RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fees | [Fees](#wasmd.x.configuration.v1beta1.Fees) |  | Fees is the starname product fee object. |





 

 

 


<a name="wasmd.x.configuration.v1beta1.Query"></a>

### Query
Query provides defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Config | [QueryConfigRequest](#wasmd.x.configuration.v1beta1.QueryConfigRequest) | [QueryConfigResponse](#wasmd.x.configuration.v1beta1.QueryConfigResponse) | Config gets starname configuration. |
| Fees | [QueryFeesRequest](#wasmd.x.configuration.v1beta1.QueryFeesRequest) | [QueryFeesResponse](#wasmd.x.configuration.v1beta1.QueryFeesResponse) | Fees gets starname product fees. |

 



<a name="x/configuration/types/msgs.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/configuration/types/msgs.proto



<a name="wasmd.x.configuration.v1beta1.MsgUpdateConfig"></a>

### MsgUpdateConfig
MsgUpdateConfig is used to update starname configuration


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| signer | [bytes](#bytes) |  | Signer is the address of the entity who is doing the transaction |
| new_configuration | [Config](#wasmd.x.configuration.v1beta1.Config) |  | NewConfiguration contains the new configuration data |






<a name="wasmd.x.configuration.v1beta1.MsgUpdateFees"></a>

### MsgUpdateFees
MsgUpdateFees is used to update the starname product fees in the starname module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| fees | [Fees](#wasmd.x.configuration.v1beta1.Fees) |  |  |
| configurer | [bytes](#bytes) |  |  |





 

 

 

 



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

