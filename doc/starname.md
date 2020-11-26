# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [x/starname/types/types.proto](#x/starname/types/types.proto)
    - [Account](#starnamed.x.starname.v1beta1.Account)
    - [Domain](#starnamed.x.starname.v1beta1.Domain)
    - [Resource](#starnamed.x.starname.v1beta1.Resource)
  
- [x/starname/types/query.proto](#x/starname/types/query.proto)
    - [QueryDomainRequest](#starnamed.x.starname.v1beta1.QueryDomainRequest)
    - [QueryDomainResponse](#starnamed.x.starname.v1beta1.QueryDomainResponse)
  
    - [Query](#starnamed.x.starname.v1beta1.Query)
  
- [x/starname/types/msgs.proto](#x/starname/types/msgs.proto)
    - [MsgAddAccountCertificates](#starnamed.x.starname.v1beta1.MsgAddAccountCertificates)
    - [MsgDeleteAccount](#starnamed.x.starname.v1beta1.MsgDeleteAccount)
    - [MsgDeleteAccountCertificate](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificate)
    - [MsgDeleteDomain](#starnamed.x.starname.v1beta1.MsgDeleteDomain)
    - [MsgRegisterAccount](#starnamed.x.starname.v1beta1.MsgRegisterAccount)
    - [MsgRegisterDomain](#starnamed.x.starname.v1beta1.MsgRegisterDomain)
    - [MsgRenewAccount](#starnamed.x.starname.v1beta1.MsgRenewAccount)
    - [MsgRenewDomain](#starnamed.x.starname.v1beta1.MsgRenewDomain)
    - [MsgReplaceAccountMetadata](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadata)
    - [MsgReplaceAccountResources](#starnamed.x.starname.v1beta1.MsgReplaceAccountResources)
    - [MsgTransferAccount](#starnamed.x.starname.v1beta1.MsgTransferAccount)
    - [MsgTransferDomain](#starnamed.x.starname.v1beta1.MsgTransferDomain)
  
- [x/starname/types/genesis.proto](#x/starname/types/genesis.proto)
    - [GenesisState](#starnamed.x.starname.v1beta1.GenesisState)
  
- [Scalar Value Types](#scalar-value-types)



<a name="x/starname/types/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/starname/types/types.proto



<a name="starnamed.x.starname.v1beta1.Account"></a>

### Account
Account defines an account that belongs to a domain
NOTE: It should not be confused with cosmos-sdk auth account
github.com/cosmos/cosmos-sdk/x/auth.Account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain references the domain this account belongs to |
| name | [google.protobuf.StringValue](#google.protobuf.StringValue) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the address that owns the account |
| broker | [bytes](#bytes) |  | Broker identifies an entity that facilitated the transaction of the account and can be empty |
| valid_until | [int64](#int64) |  | ValidUntil defines a unix timestamp of the expiration of the account in seconds |
| resources | [Resource](#starnamed.x.starname.v1beta1.Resource) | repeated | Resources is the list of resources an account resolves to |
| certificates | [bytes](#bytes) | repeated | Certificates contains the list of certificates to identify the account owner |
| metadata_uri | [string](#string) |  | MetadataURI contains a link to extra information regarding the account |






<a name="starnamed.x.starname.v1beta1.Domain"></a>

### Domain
Domain defines a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is the name of the domain |
| admin | [bytes](#bytes) |  | Admin is the owner of the domain |
| broker | [bytes](#bytes) |  |  |
| valid_until | [int64](#int64) |  | ValidUntil is a unix timestamp defines the time when the domain will become invalid in seconds |
| type | [string](#string) |  | Type defines the type of the domain |






<a name="starnamed.x.starname.v1beta1.Resource"></a>

### Resource
Resource defines a resource owned by an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uri | [string](#string) |  | URI defines the ID of the resource |
| resource | [string](#string) |  | Resource is the resource |





 

 

 

 



<a name="x/starname/types/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/starname/types/query.proto



<a name="starnamed.x.starname.v1beta1.QueryDomainRequest"></a>

### QueryDomainRequest
QueryDomainRequest is the request type for the Query/Domain RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  | Name is the name of the domain. |






<a name="starnamed.x.starname.v1beta1.QueryDomainResponse"></a>

### QueryDomainResponse
QueryDomainResponse is the response type for the Query/Domain RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [Domain](#starnamed.x.starname.v1beta1.Domain) |  | Domain is the information associated with the domain. |





 

 

 


<a name="starnamed.x.starname.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Domain | [QueryDomainRequest](#starnamed.x.starname.v1beta1.QueryDomainRequest) | [QueryDomainResponse](#starnamed.x.starname.v1beta1.QueryDomainResponse) | Domain gets a starname&#39;s domain. |

 



<a name="x/starname/types/msgs.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/starname/types/msgs.proto



<a name="starnamed.x.starname.v1beta1.MsgAddAccountCertificates"></a>

### MsgAddAccountCertificates
MsgAddAccountCertificates is the message used when a user wants to add new certificates to his account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_certificate | [bytes](#bytes) |  | NewCertificate is the new certificate to add |






<a name="starnamed.x.starname.v1beta1.MsgDeleteAccount"></a>

### MsgDeleteAccount
MsgDeleteAccount is the request model used to delete an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgDeleteAccountCertificate"></a>

### MsgDeleteAccountCertificate
MsgDeleteAccountCertificate is the request model used to remove certificates from an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| delete_certificate | [bytes](#bytes) |  | DeleteCertificate is the certificate to delete |






<a name="starnamed.x.starname.v1beta1.MsgDeleteDomain"></a>

### MsgDeleteDomain
MsgDeleteDomain is the request model to delete a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgRegisterAccount"></a>

### MsgRegisterAccount
MsgRegisterAccount is the request model used to register new accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| broker | [bytes](#bytes) |  | Broker is the account that facilitated the transaction |
| registerer | [bytes](#bytes) |  | Registerer is the user who registers this account |
| resources | [Resource](#starnamed.x.starname.v1beta1.Resource) | repeated | Resources are the blockchain addresses of the account |






<a name="starnamed.x.starname.v1beta1.MsgRegisterDomain"></a>

### MsgRegisterDomain
MsgRegisterDomain is the request used to register new domains


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| admin | [bytes](#bytes) |  |  |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| broker | [bytes](#bytes) |  | Broker is the account that facilitated the transaction |
| domain_type | [string](#string) |  | DomainType defines the type of the domain |






<a name="starnamed.x.starname.v1beta1.MsgRenewAccount"></a>

### MsgRenewAccount
MsgRenewAccount is the request model used to renew accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| signer | [bytes](#bytes) |  | Signer is the signer of the request |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgRenewDomain"></a>

### MsgRenewDomain
MsgRenewDomain is the request model used to renew a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| signer | [bytes](#bytes) |  | Signer is the signer of the request |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountMetadata"></a>

### MsgReplaceAccountMetadata
MsgReplaceAccountMetadata is the function used to set accounts metadata


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_metadata_uri | [string](#string) |  | NewMetadataURI is the metadata URI of the account we want to update or insert |






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountResources"></a>

### MsgReplaceAccountResources
MsgReplaceAccountResources is the request model used to renew resources associated with an account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_resources | [Resource](#starnamed.x.starname.v1beta1.Resource) | repeated | NewResources are the new resources |






<a name="starnamed.x.starname.v1beta1.MsgTransferAccount"></a>

### MsgTransferAccount
MsgTransferAccount is the request model used to transfer accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_owner | [bytes](#bytes) |  | NewOwner is the new owner of the account |
| reset | [bool](#bool) |  | ToReset if true, removes all old data from account |






<a name="starnamed.x.starname.v1beta1.MsgTransferDomain"></a>

### MsgTransferDomain
MsgTransferDomain is the request model used to transfer a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the name of the domain |
| owner | [bytes](#bytes) |  | Owner is the owner of the domain |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_admin | [bytes](#bytes) |  | NewAdmin is the new owner of the domain |
| transfer_flag | [int64](#int64) |  | TransferFlag controls the operations that occurs on a domain&#39;s accounts |





 

 

 

 



<a name="x/starname/types/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/starname/types/genesis.proto



<a name="starnamed.x.starname.v1beta1.GenesisState"></a>

### GenesisState
GenesisState - genesis state of x/starname


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domains | [Domain](#starnamed.x.starname.v1beta1.Domain) | repeated |  |
| accounts | [Account](#starnamed.x.starname.v1beta1.Account) | repeated |  |





 

 

 

 



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

