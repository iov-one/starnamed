# Protocol Documentation
<a name="top"></a>

## Table of Contents

- [x/starname/types/tx.proto](#x/starname/types/tx.proto)
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
  
- [x/starname/types/types.proto](#x/starname/types/types.proto)
    - [Account](#starnamed.x.starname.v1beta1.Account)
    - [Domain](#starnamed.x.starname.v1beta1.Domain)
    - [Resource](#starnamed.x.starname.v1beta1.Resource)
  
- [x/starname/types/query.proto](#x/starname/types/query.proto)
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
  
    - [Query](#starnamed.x.starname.v1beta1.Query)
  
- [x/starname/types/genesis.proto](#x/starname/types/genesis.proto)
    - [GenesisState](#starnamed.x.starname.v1beta1.GenesisState)
  
- [Scalar Value Types](#scalar-value-types)



<a name="x/starname/types/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## x/starname/types/tx.proto



<a name="starnamed.x.starname.v1beta1.MsgAddAccountCertificate"></a>

### MsgAddAccountCertificate
MsgAddAccountCertificate is the message used when a user wants to add new certificates to his account


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| owner | [bytes](#bytes) |  | Owner is the owner of the account |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_certificate | [bytes](#bytes) |  | NewCertificate is the new certificate to add |






<a name="starnamed.x.starname.v1beta1.MsgAddAccountCertificateResponse"></a>

### MsgAddAccountCertificateResponse
MsgAddAccountCertificateResponse returns an empty response.






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
| domain | [string](#string) |  | Domain is the domain of the account |
| owner | [string](#string) |  | Owner is the owner of the account |
| payer | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgDeleteDomainResponse"></a>

### MsgDeleteDomainResponse
MsgDeleteDomainResponse returns an empty response.






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






<a name="starnamed.x.starname.v1beta1.MsgRegisterAccountResponse"></a>

### MsgRegisterAccountResponse
MsgRegisterAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRegisterDomain"></a>

### MsgRegisterDomain
MsgRegisterDomain is the request used to register new domains


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| name | [string](#string) |  |  |
| admin | [string](#string) |  |  |
| payer | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| broker | [string](#string) |  | Broker is the account that facilitated the transaction |
| domain_type | [string](#string) |  | DomainType defines the type of the domain |






<a name="starnamed.x.starname.v1beta1.MsgRegisterDomainResponse"></a>

### MsgRegisterDomainResponse
MsgRegisterDomainResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRenewAccount"></a>

### MsgRenewAccount
MsgRenewAccount is the request model used to renew accounts


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| name | [string](#string) |  | Name is the name of the account |
| signer | [bytes](#bytes) |  | Signer is the signer of the request |
| payer | [bytes](#bytes) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgRenewAccountResponse"></a>

### MsgRenewAccountResponse
MsgRenewAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgRenewDomain"></a>

### MsgRenewDomain
MsgRenewDomain is the request model used to renew a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the domain of the account |
| signer | [string](#string) |  | Signer is the signer of the request |
| payer | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |






<a name="starnamed.x.starname.v1beta1.MsgRenewDomainResponse"></a>

### MsgRenewDomainResponse
MsgRegisterDomain returns an empty response.






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






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountMetadataResponse"></a>

### MsgReplaceAccountMetadataResponse
MsgReplaceAccountMetadataResponse returns an empty response.






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






<a name="starnamed.x.starname.v1beta1.MsgReplaceAccountResourcesResponse"></a>

### MsgReplaceAccountResourcesResponse
MsgReplaceAccountResourcesResponse






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






<a name="starnamed.x.starname.v1beta1.MsgTransferAccountResponse"></a>

### MsgTransferAccountResponse
MsgTransferAccountResponse returns an empty response.






<a name="starnamed.x.starname.v1beta1.MsgTransferDomain"></a>

### MsgTransferDomain
MsgTransferDomain is the request model used to transfer a domain


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the name of the domain |
| owner | [string](#string) |  | Owner is the owner of the domain |
| payer | [string](#string) |  | Payer is the address of the entity that pays the product and transaction fees |
| new_admin | [string](#string) |  | NewAdmin is the new owner of the domain |
| transfer_flag | [int64](#int64) |  | TransferFlag controls the operations that occurs on a domain&#39;s accounts |






<a name="starnamed.x.starname.v1beta1.MsgTransferDomainResponse"></a>

### MsgTransferDomainResponse
MsgTransferDomainResponse returns an empty response.





 

 

 


<a name="starnamed.x.starname.v1beta1.Msg"></a>

### Msg
Msg defines the starname Msg service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| AddAccountCertificate | [MsgAddAccountCertificate](#starnamed.x.starname.v1beta1.MsgAddAccountCertificate) | [MsgAddAccountCertificateResponse](#starnamed.x.starname.v1beta1.MsgAddAccountCertificateResponse) | AddAccountCertificate adds a certificate to an Account |
| DeleteAccount | [MsgDeleteAccount](#starnamed.x.starname.v1beta1.MsgDeleteAccount) | [MsgDeleteAccountResponse](#starnamed.x.starname.v1beta1.MsgDeleteAccountResponse) | DeleteAccount registers a Domain |
| DeleteAccountCertificate | [MsgDeleteAccountCertificate](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificate) | [MsgDeleteAccountCertificateResponse](#starnamed.x.starname.v1beta1.MsgDeleteAccountCertificateResponse) | DeleteAccountCertificate deletes a certificate from an account |
| DeleteDomain | [MsgDeleteDomain](#starnamed.x.starname.v1beta1.MsgDeleteDomain) | [MsgDeleteDomainResponse](#starnamed.x.starname.v1beta1.MsgDeleteDomainResponse) | DeleteDomain registers a Domain |
| RegisterAccount | [MsgRegisterAccount](#starnamed.x.starname.v1beta1.MsgRegisterAccount) | [MsgRegisterAccountResponse](#starnamed.x.starname.v1beta1.MsgRegisterAccountResponse) | RegisterAccount registers an Account |
| RegisterDomain | [MsgRegisterDomain](#starnamed.x.starname.v1beta1.MsgRegisterDomain) | [MsgRegisterDomainResponse](#starnamed.x.starname.v1beta1.MsgRegisterDomainResponse) | RegisterDomain registers a Domain |
| RenewAccount | [MsgRenewAccount](#starnamed.x.starname.v1beta1.MsgRenewAccount) | [MsgRenewAccountResponse](#starnamed.x.starname.v1beta1.MsgRenewAccountResponse) | RenewAccount registers a Domain |
| RenewDomain | [MsgRenewDomain](#starnamed.x.starname.v1beta1.MsgRenewDomain) | [MsgRenewDomainResponse](#starnamed.x.starname.v1beta1.MsgRenewDomainResponse) | RenewDomain registers a Domain |
| ReplaceAccountMetadata | [MsgReplaceAccountMetadata](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadata) | [MsgReplaceAccountMetadataResponse](#starnamed.x.starname.v1beta1.MsgReplaceAccountMetadataResponse) | ReplaceAccountMetadata registers a Domain |
| ReplaceAccountResources | [MsgReplaceAccountResources](#starnamed.x.starname.v1beta1.MsgReplaceAccountResources) | [MsgReplaceAccountResourcesResponse](#starnamed.x.starname.v1beta1.MsgReplaceAccountResourcesResponse) | ReplaceAccountResources registers a Domain |
| TransferAccount | [MsgTransferAccount](#starnamed.x.starname.v1beta1.MsgTransferAccount) | [MsgTransferAccountResponse](#starnamed.x.starname.v1beta1.MsgTransferAccountResponse) | TransferAccount registers a Domain |
| TransferDomain | [MsgTransferDomain](#starnamed.x.starname.v1beta1.MsgTransferDomain) | [MsgTransferDomainResponse](#starnamed.x.starname.v1beta1.MsgTransferDomainResponse) | TransferDomain registers a Domain |

 



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



<a name="starnamed.x.starname.v1beta1.QueryBrokerAccountsRequest"></a>

### QueryBrokerAccountsRequest
QueryBrokerAccountsRequest is the request type for the Query/BrokerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| broker | [string](#string) |  | Broker is the broker of accounts. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryBrokerAccountsResponse"></a>

### QueryBrokerAccountsResponse
QueryBrokerAccountsResponse is the response type for the Query/BrokerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| accounts | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts is the accounts associated with broker. |
| page | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryBrokerDomainsRequest"></a>

### QueryBrokerDomainsRequest
QueryBrokerDomainsRequest is the request type for the Query/BrokerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| broker | [string](#string) |  | Broker is the broker of accounts. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryBrokerDomainsResponse"></a>

### QueryBrokerDomainsResponse
QueryBrokerDomainsResponse is the response type for the Query/BrokerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domains | [Domain](#starnamed.x.starname.v1beta1.Domain) | repeated | Accounts is the accounts associated with broker. |
| page | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryDomainAccountsRequest"></a>

### QueryDomainAccountsRequest
QueryDomainAccountsRequest is the request type for the Query/DomainAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domain | [string](#string) |  | Domain is the name of the domain. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryDomainAccountsResponse"></a>

### QueryDomainAccountsResponse
QueryDomainAccountsResponse is the response type for the Query/DomainAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| accounts | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts is the accounts associated with the domain. |
| page | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






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






<a name="starnamed.x.starname.v1beta1.QueryOwnerAccountsRequest"></a>

### QueryOwnerAccountsRequest
QueryOwnerAccountsRequest is the request type for the Query/OwnerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner | [string](#string) |  | Owner is the owner of accounts. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryOwnerAccountsResponse"></a>

### QueryOwnerAccountsResponse
QueryOwnerAccountsResponse is the response type for the Query/OwnerAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| accounts | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts is the accounts associated with owner. |
| page | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryOwnerDomainsRequest"></a>

### QueryOwnerDomainsRequest
QueryOwnerDomainsRequest is the request type for the Query/OwnerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| owner | [string](#string) |  | Owner is the owner of accounts. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryOwnerDomainsResponse"></a>

### QueryOwnerDomainsResponse
QueryOwnerDomainsResponse is the response type for the Query/OwnerDomains RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| domains | [Domain](#starnamed.x.starname.v1beta1.Domain) | repeated | Accounts is the accounts associated with owner. |
| page | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryResourceAccountsRequest"></a>

### QueryResourceAccountsRequest
QueryResourceAccountsRequest is the request type for the Query/ResourceAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| uri | [string](#string) |  | Uri is the uri of the resource. query.pb.gw.to doesn&#39;t respect gogoproto.customname, so we&#39;re stuck with Uri. |
| resource | [string](#string) |  | Resource is the resource of interest. |
| pagination | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryResourceAccountsResponse"></a>

### QueryResourceAccountsResponse
QueryResourceAccountsResponse is the response type for the Query/ResourceAccounts RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| accounts | [Account](#starnamed.x.starname.v1beta1.Account) | repeated | Accounts are the accounts associated with the resource. |
| page | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  |  |






<a name="starnamed.x.starname.v1beta1.QueryStarnameRequest"></a>

### QueryStarnameRequest
QueryStarnameRequest is the request type for the Query/Starname RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| starname | [string](#string) |  | Starname is the of the form account*domain. |






<a name="starnamed.x.starname.v1beta1.QueryStarnameResponse"></a>

### QueryStarnameResponse
QueryStarnameResponse is the response type for the Query/Starname RPC method.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| account | [Account](#starnamed.x.starname.v1beta1.Account) |  | Account is the information associated with the starname. |





 

 

 


<a name="starnamed.x.starname.v1beta1.Query"></a>

### Query
Query defines the gRPC querier service.

| Method Name | Request Type | Response Type | Description |
| ----------- | ------------ | ------------- | ------------|
| Domain | [QueryDomainRequest](#starnamed.x.starname.v1beta1.QueryDomainRequest) | [QueryDomainResponse](#starnamed.x.starname.v1beta1.QueryDomainResponse) | Domain gets a starname&#39;s domain info. |
| DomainAccounts | [QueryDomainAccountsRequest](#starnamed.x.starname.v1beta1.QueryDomainAccountsRequest) | [QueryDomainAccountsResponse](#starnamed.x.starname.v1beta1.QueryDomainAccountsResponse) | DomainAccounts gets accounts associated with a given domain. |
| Starname | [QueryStarnameRequest](#starnamed.x.starname.v1beta1.QueryStarnameRequest) | [QueryStarnameResponse](#starnamed.x.starname.v1beta1.QueryStarnameResponse) | Starname gets accounts associated with a given domain. |
| OwnerAccounts | [QueryOwnerAccountsRequest](#starnamed.x.starname.v1beta1.QueryOwnerAccountsRequest) | [QueryOwnerAccountsResponse](#starnamed.x.starname.v1beta1.QueryOwnerAccountsResponse) | OwnerAccounts gets accounts associated with a given owner. |
| OwnerDomains | [QueryOwnerDomainsRequest](#starnamed.x.starname.v1beta1.QueryOwnerDomainsRequest) | [QueryOwnerDomainsResponse](#starnamed.x.starname.v1beta1.QueryOwnerDomainsResponse) | OwnerDomains gets domains associated with a given owner. |
| ResourceAccounts | [QueryResourceAccountsRequest](#starnamed.x.starname.v1beta1.QueryResourceAccountsRequest) | [QueryResourceAccountsResponse](#starnamed.x.starname.v1beta1.QueryResourceAccountsResponse) | ResourceAccounts gets accounts associated with a given resource. |
| BrokerAccounts | [QueryBrokerAccountsRequest](#starnamed.x.starname.v1beta1.QueryBrokerAccountsRequest) | [QueryBrokerAccountsResponse](#starnamed.x.starname.v1beta1.QueryBrokerAccountsResponse) | BrokerAccounts gets accounts associated with a given broker. |
| BrokerDomains | [QueryBrokerDomainsRequest](#starnamed.x.starname.v1beta1.QueryBrokerDomainsRequest) | [QueryBrokerDomainsResponse](#starnamed.x.starname.v1beta1.QueryBrokerDomainsResponse) | BrokerDomains gets domains associated with a given broker. |

 



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

