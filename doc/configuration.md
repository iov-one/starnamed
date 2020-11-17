# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [x/configuration/types/types.proto](#x/configuration/types/types.proto)
    - [Config](#.Config)
    - [Fees](#.Fees)
  
- [x/configuration/types/msgs.proto](#x/configuration/types/msgs.proto)
    - [MsgUpdateConfig](#.MsgUpdateConfig)
    - [MsgUpdateFees](#.MsgUpdateFees)
  
- [Scalar Value Types](#scalar-value-types)



<a name="x/configuration/types/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/configuration/types/types.proto



<a name=".Config"></a>

### Config
Config is the configuration of the network


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Configurer | [bytes](#bytes) |  | Configurer is the configuration owner, the addresses allowed to handle fees and register domains with no superuser |
| ValidDomainName | [string](#string) |  | ValidDomainName defines a regexp that determines if a domain name is valid or not |
| ValidAccountName | [string](#string) |  | ValidAccountName defines a regexp that determines if an account name is valid or not |
| ValidURI | [string](#string) |  | ValidURI defines a regexp that determines if resource uri is valid or not |
| ValidResource | [string](#string) |  | ValidResource determines a regexp for a resource content |
| DomainRenewalPeriod | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainRenewalPeriod defines the duration of the domain renewal period in seconds |
| DomainRenewalCountMax | [uint32](#uint32) |  | DomainRenewalCountMax defines maximum number of domain renewals a user can do |
| DomainGracePeriod | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainGracePeriod defines the grace period for a domain deletion in seconds |
| AccountRenewalPeriod | [google.protobuf.Duration](#google.protobuf.Duration) |  | AccountRenewalPeriod defines the duration of the account renewal period in seconds |
| AccountRenewalCountMax | [uint32](#uint32) |  | AccountRenewalCountMax defines maximum number of account renewals a user can do |
| AccountGracePeriod | [google.protobuf.Duration](#google.protobuf.Duration) |  | DomainGracePeriod defines the grace period for a domain deletion in seconds |
| ResourcesMax | [uint32](#uint32) |  | ResourcesMax defines maximum number of resources could be saved under an account |
| CertificateSizeMax | [uint64](#uint64) |  | CertificateSizeMax defines maximum size of a certificate that could be saved under an account |
| CertificateCountMax | [uint32](#uint32) |  | CertificateCountMax defines maximum number of certificates that could be saved under an account |
| MetadataSizeMax | [uint64](#uint64) |  | MetadataSizeMax defines maximum size of metadata that could be saved under an account |






<a name=".Fees"></a>

### Fees
Fees contains different type of fees
to calculate coins to detract when
processing different messages


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| FeeCoinDenom | [string](#string) |  | FeeCoinDenom defines the denominator of the coin used to process fees |
| FeeCoinPrice | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | FeeCoinPrice defines the price of the coin |
| FeeDefault | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | FeeDefault is the parameter defining the default fee |
| RegisterAccountClosed | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterAccountClosed is the fee to be paid to register an account in a closed domain |
| RegisterAccountOpen | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterAccountOpen is the fee to be paid to register an account in an open domain |
| TransferAccountClosed | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | TransferAccountClosed is the fee to be paid to register an account in a closed domain |
| TransferAccountOpen | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | TransferAccountOpen is the fee to be paid to register an account in an open domain |
| ReplaceAccountResources | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | ReplaceAccountResources is the fee to be paid to replace account&#39;s resources |
| AddAccountCertificate | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | AddAccountCertificate is the fee to be paid to add a certificate to an account |
| DelAccountCertificate | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | DelAccountCertificate is the feed to be paid to delete a certificate in an account |
| SetAccountMetadata | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | SetAccountMetadata is the fee to be paid to set account&#39;s metadata |
| RegisterDomain1 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain1 is the fee to be paid to register a domain with one character |
| RegisterDomain2 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain2 is the fee to be paid to register a domain with two characters |
| RegisterDomain3 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain3 is the fee to be paid to register a domain with three characters |
| RegisterDomain4 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain4 is the fee to be paid to register a domain with four characters |
| RegisterDomain5 | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomain5 is the fee to be paid to register a domain with five characters |
| RegisterDomainDefault | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomainDefault is the fee to be paid to register a domain with more than five characters |
| RegisterOpenDomainMultiplier | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | RegisterDomainMultiplier is the multiplication applied to fees in register domain operations if they&#39;re of open type |
| TransferDomainClosed | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | TransferDomainClosed is the fee to be paid to transfer a closed domain |
| TransferDomainOpen | [cosmos.base.v1beta1.DecProto](#cosmos.base.v1beta1.DecProto) |  | TransferDomainOpen is the fee to be paid to transfer open domains |
| RenewDomainOpen | [string](#string) |  | RenewDomainOpen is the fee to be paid to renew an open domain |





 

 

 

 



<a name="x/configuration/types/msgs.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/configuration/types/msgs.proto



<a name=".MsgUpdateConfig"></a>

### MsgUpdateConfig
MsgUpdateConfig is used to update
configuration using a multisig strategy


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Signer | [bytes](#bytes) |  | Signer is the address of the entity who is doing the transaction |
| NewConfiguration | [Config](#Config) |  | NewConfiguration contains the new configuration data |






<a name=".MsgUpdateFees"></a>

### MsgUpdateFees
MsgUpdateFees is used to update
the product fees required when interacting
with the starname module.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| Fees | [Fees](#Fees) |  |  |
| Configurer | [bytes](#bytes) |  |  |





 

 

 

 



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

