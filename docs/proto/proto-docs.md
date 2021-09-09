<!-- This file is auto-generated. Please do not modify it yourself. -->
# Protobuf Documentation
<a name="top"></a>

## Table of Contents

- [cosmwasm/wasm/v1/types.proto](#cosmwasm/wasm/v1/types.proto)
    - [AbsoluteTxPosition](#cosmwasm.wasm.v1.AbsoluteTxPosition)
    - [AccessConfig](#cosmwasm.wasm.v1.AccessConfig)
    - [AccessTypeParam](#cosmwasm.wasm.v1.AccessTypeParam)
    - [CodeInfo](#cosmwasm.wasm.v1.CodeInfo)
    - [ContractCodeHistoryEntry](#cosmwasm.wasm.v1.ContractCodeHistoryEntry)
    - [ContractInfo](#cosmwasm.wasm.v1.ContractInfo)
    - [Model](#cosmwasm.wasm.v1.Model)
    - [Params](#cosmwasm.wasm.v1.Params)

    - [AccessType](#cosmwasm.wasm.v1.AccessType)
    - [ContractCodeHistoryOperationType](#cosmwasm.wasm.v1.ContractCodeHistoryOperationType)

- [cosmwasm/wasm/v1/tx.proto](#cosmwasm/wasm/v1/tx.proto)
    - [MsgClearAdmin](#cosmwasm.wasm.v1.MsgClearAdmin)
    - [MsgClearAdminResponse](#cosmwasm.wasm.v1.MsgClearAdminResponse)
    - [MsgExecuteContract](#cosmwasm.wasm.v1.MsgExecuteContract)
    - [MsgExecuteContractResponse](#cosmwasm.wasm.v1.MsgExecuteContractResponse)
    - [MsgInstantiateContract](#cosmwasm.wasm.v1.MsgInstantiateContract)
    - [MsgInstantiateContractResponse](#cosmwasm.wasm.v1.MsgInstantiateContractResponse)
    - [MsgMigrateContract](#cosmwasm.wasm.v1.MsgMigrateContract)
    - [MsgMigrateContractResponse](#cosmwasm.wasm.v1.MsgMigrateContractResponse)
    - [MsgStoreCode](#cosmwasm.wasm.v1.MsgStoreCode)
    - [MsgStoreCodeResponse](#cosmwasm.wasm.v1.MsgStoreCodeResponse)
    - [MsgUpdateAdmin](#cosmwasm.wasm.v1.MsgUpdateAdmin)
    - [MsgUpdateAdminResponse](#cosmwasm.wasm.v1.MsgUpdateAdminResponse)

    - [Msg](#cosmwasm.wasm.v1.Msg)

- [cosmwasm/wasm/v1/genesis.proto](#cosmwasm/wasm/v1/genesis.proto)
    - [Code](#cosmwasm.wasm.v1.Code)
    - [Contract](#cosmwasm.wasm.v1.Contract)
    - [GenesisState](#cosmwasm.wasm.v1.GenesisState)
    - [GenesisState.GenMsgs](#cosmwasm.wasm.v1.GenesisState.GenMsgs)
    - [Sequence](#cosmwasm.wasm.v1.Sequence)

- [cosmwasm/wasm/v1/ibc.proto](#cosmwasm/wasm/v1/ibc.proto)
    - [MsgIBCCloseChannel](#cosmwasm.wasm.v1.MsgIBCCloseChannel)
    - [MsgIBCSend](#cosmwasm.wasm.v1.MsgIBCSend)

- [cosmwasm/wasm/v1/proposal.proto](#cosmwasm/wasm/v1/proposal.proto)
    - [ClearAdminProposal](#cosmwasm.wasm.v1.ClearAdminProposal)
    - [InstantiateContractProposal](#cosmwasm.wasm.v1.InstantiateContractProposal)
    - [MigrateContractProposal](#cosmwasm.wasm.v1.MigrateContractProposal)
    - [PinCodesProposal](#cosmwasm.wasm.v1.PinCodesProposal)
    - [StoreCodeProposal](#cosmwasm.wasm.v1.StoreCodeProposal)
    - [UnpinCodesProposal](#cosmwasm.wasm.v1.UnpinCodesProposal)
    - [UpdateAdminProposal](#cosmwasm.wasm.v1.UpdateAdminProposal)

- [cosmwasm/wasm/v1/query.proto](#cosmwasm/wasm/v1/query.proto)
    - [CodeInfoResponse](#cosmwasm.wasm.v1.CodeInfoResponse)
    - [QueryAllContractStateRequest](#cosmwasm.wasm.v1.QueryAllContractStateRequest)
    - [QueryAllContractStateResponse](#cosmwasm.wasm.v1.QueryAllContractStateResponse)
    - [QueryCodeRequest](#cosmwasm.wasm.v1.QueryCodeRequest)
    - [QueryCodeResponse](#cosmwasm.wasm.v1.QueryCodeResponse)
    - [QueryCodesRequest](#cosmwasm.wasm.v1.QueryCodesRequest)
    - [QueryCodesResponse](#cosmwasm.wasm.v1.QueryCodesResponse)
    - [QueryContractHistoryRequest](#cosmwasm.wasm.v1.QueryContractHistoryRequest)
    - [QueryContractHistoryResponse](#cosmwasm.wasm.v1.QueryContractHistoryResponse)
    - [QueryContractInfoRequest](#cosmwasm.wasm.v1.QueryContractInfoRequest)
    - [QueryContractInfoResponse](#cosmwasm.wasm.v1.QueryContractInfoResponse)
    - [QueryContractsByCodeRequest](#cosmwasm.wasm.v1.QueryContractsByCodeRequest)
    - [QueryContractsByCodeResponse](#cosmwasm.wasm.v1.QueryContractsByCodeResponse)
    - [QueryRawContractStateRequest](#cosmwasm.wasm.v1.QueryRawContractStateRequest)
    - [QueryRawContractStateResponse](#cosmwasm.wasm.v1.QueryRawContractStateResponse)
    - [QuerySmartContractStateRequest](#cosmwasm.wasm.v1.QuerySmartContractStateRequest)
    - [QuerySmartContractStateResponse](#cosmwasm.wasm.v1.QuerySmartContractStateResponse)

    - [Query](#cosmwasm.wasm.v1.Query)

- [Scalar Value Types](#scalar-value-types)



<a name="cosmwasm/wasm/v1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/wasm/v1/types.proto



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.AbsoluteTxPosition"></a>
=======
<a name="cosmwasm.wasm.v1.AbsoluteTxPosition"></a>
>>>>>>> upstream/master

### AbsoluteTxPosition
AbsoluteTxPosition is a unique transaction position that allows for global
ordering of transactions.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `block_height` | [uint64](#uint64) |  | BlockHeight is the block the contract was created at |
| `tx_index` | [uint64](#uint64) |  | TxIndex is a monotonic counter within the block (actual transaction index, or gas consumed) |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.AccessConfig"></a>
=======
<a name="cosmwasm.wasm.v1.AccessConfig"></a>
>>>>>>> upstream/master

### AccessConfig
AccessConfig access control type.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `permission` | [AccessType](#starnamed.x.wasm.v1beta1.AccessType) |  |  |
=======
| `permission` | [AccessType](#cosmwasm.wasm.v1.AccessType) |  |  |
>>>>>>> upstream/master
| `address` | [string](#string) |  |  |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.AccessTypeParam"></a>
=======
<a name="cosmwasm.wasm.v1.AccessTypeParam"></a>
>>>>>>> upstream/master

### AccessTypeParam
AccessTypeParam


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `value` | [AccessType](#starnamed.x.wasm.v1beta1.AccessType) |  |  |
=======
| `value` | [AccessType](#cosmwasm.wasm.v1.AccessType) |  |  |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.CodeInfo"></a>
=======
<a name="cosmwasm.wasm.v1.CodeInfo"></a>
>>>>>>> upstream/master

### CodeInfo
CodeInfo is data for the uploaded contract WASM code


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_hash` | [bytes](#bytes) |  | CodeHash is the unique identifier created by wasmvm |
| `creator` | [string](#string) |  | Creator address who initially stored the code |
<<<<<<< HEAD
| `source` | [string](#string) |  | Source is a valid absolute HTTPS URI to the contract's source code, optional |
| `builder` | [string](#string) |  | Builder is a valid docker image name with tag, optional |
| `instantiate_config` | [AccessConfig](#starnamed.x.wasm.v1beta1.AccessConfig) |  | InstantiateConfig access control to apply on contract creation, optional |
=======
| `instantiate_config` | [AccessConfig](#cosmwasm.wasm.v1.AccessConfig) |  | InstantiateConfig access control to apply on contract creation, optional |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.ContractCodeHistoryEntry"></a>
=======
<a name="cosmwasm.wasm.v1.ContractCodeHistoryEntry"></a>
>>>>>>> upstream/master

### ContractCodeHistoryEntry
ContractCodeHistoryEntry metadata to a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `operation` | [ContractCodeHistoryOperationType](#starnamed.x.wasm.v1beta1.ContractCodeHistoryOperationType) |  |  |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `updated` | [AbsoluteTxPosition](#starnamed.x.wasm.v1beta1.AbsoluteTxPosition) |  | Updated Tx position when the operation was executed. |
=======
| `operation` | [ContractCodeHistoryOperationType](#cosmwasm.wasm.v1.ContractCodeHistoryOperationType) |  |  |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `updated` | [AbsoluteTxPosition](#cosmwasm.wasm.v1.AbsoluteTxPosition) |  | Updated Tx position when the operation was executed. |
>>>>>>> upstream/master
| `msg` | [bytes](#bytes) |  |  |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.ContractInfo"></a>
=======
<a name="cosmwasm.wasm.v1.ContractInfo"></a>
>>>>>>> upstream/master

### ContractInfo
ContractInfo stores a WASM contract instance


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored Wasm code |
| `creator` | [string](#string) |  | Creator address who initially instantiated the contract |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a contract instance. |
<<<<<<< HEAD
| `created` | [AbsoluteTxPosition](#starnamed.x.wasm.v1beta1.AbsoluteTxPosition) |  | Created Tx position when the contract was instantiated. This data should kept internal and not be exposed via query results. Just use for sorting |
=======
| `created` | [AbsoluteTxPosition](#cosmwasm.wasm.v1.AbsoluteTxPosition) |  | Created Tx position when the contract was instantiated. This data should kept internal and not be exposed via query results. Just use for sorting |
>>>>>>> upstream/master
| `ibc_port_id` | [string](#string) |  |  |
| `extension` | [google.protobuf.Any](#google.protobuf.Any) |  | Extension is an extension point to store custom metadata within the persistence model. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Model"></a>
=======
<a name="cosmwasm.wasm.v1.Model"></a>
>>>>>>> upstream/master

### Model
Model is a struct that holds a KV pair


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `key` | [bytes](#bytes) |  | hex-encode key to read it better (this is often ascii) |
| `value` | [bytes](#bytes) |  | base64-encode raw value |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Params"></a>
=======
<a name="cosmwasm.wasm.v1.Params"></a>
>>>>>>> upstream/master

### Params
Params defines the set of wasm parameters.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `code_upload_access` | [AccessConfig](#starnamed.x.wasm.v1beta1.AccessConfig) |  |  |
| `instantiate_default_permission` | [AccessType](#starnamed.x.wasm.v1beta1.AccessType) |  |  |
=======
| `code_upload_access` | [AccessConfig](#cosmwasm.wasm.v1.AccessConfig) |  |  |
| `instantiate_default_permission` | [AccessType](#cosmwasm.wasm.v1.AccessType) |  |  |
>>>>>>> upstream/master
| `max_wasm_code_size` | [uint64](#uint64) |  |  |





 <!-- end messages -->


<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.AccessType"></a>
=======
<a name="cosmwasm.wasm.v1.AccessType"></a>
>>>>>>> upstream/master

### AccessType
AccessType permission types

| Name | Number | Description |
| ---- | ------ | ----------- |
| ACCESS_TYPE_UNSPECIFIED | 0 | AccessTypeUnspecified placeholder for empty value |
| ACCESS_TYPE_NOBODY | 1 | AccessTypeNobody forbidden |
| ACCESS_TYPE_ONLY_ADDRESS | 2 | AccessTypeOnlyAddress restricted to an address |
| ACCESS_TYPE_EVERYBODY | 3 | AccessTypeEverybody unrestricted |



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.ContractCodeHistoryOperationType"></a>
=======
<a name="cosmwasm.wasm.v1.ContractCodeHistoryOperationType"></a>
>>>>>>> upstream/master

### ContractCodeHistoryOperationType
ContractCodeHistoryOperationType actions that caused a code change

| Name | Number | Description |
| ---- | ------ | ----------- |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_UNSPECIFIED | 0 | ContractCodeHistoryOperationTypeUnspecified placeholder for empty value |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_INIT | 1 | ContractCodeHistoryOperationTypeInit on chain contract instantiation |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_MIGRATE | 2 | ContractCodeHistoryOperationTypeMigrate code migration |
| CONTRACT_CODE_HISTORY_OPERATION_TYPE_GENESIS | 3 | ContractCodeHistoryOperationTypeGenesis based on genesis data |


 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/wasm/v1/tx.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/wasm/v1/tx.proto



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgClearAdmin"></a>
=======
<a name="cosmwasm.wasm.v1.MsgClearAdmin"></a>
>>>>>>> upstream/master

### MsgClearAdmin
MsgClearAdmin removes any admin stored for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgClearAdminResponse"></a>
=======
<a name="cosmwasm.wasm.v1.MsgClearAdminResponse"></a>
>>>>>>> upstream/master

### MsgClearAdminResponse
MsgClearAdminResponse returns empty data






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgExecuteContract"></a>
=======
<a name="cosmwasm.wasm.v1.MsgExecuteContract"></a>
>>>>>>> upstream/master

### MsgExecuteContract
MsgExecuteContract submits the given message data to a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract |
| `funds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Funds coins that are transferred to the contract on execution |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgExecuteContractResponse"></a>
=======
<a name="cosmwasm.wasm.v1.MsgExecuteContractResponse"></a>
>>>>>>> upstream/master

### MsgExecuteContractResponse
MsgExecuteContractResponse returns execution result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgInstantiateContract"></a>
=======
<a name="cosmwasm.wasm.v1.MsgInstantiateContract"></a>
>>>>>>> upstream/master

### MsgInstantiateContract
MsgInstantiateContract create a new smart contract instance for the given
code id.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a contract instance. |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on instantiation |
| `funds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Funds coins that are transferred to the contract on instantiation |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgInstantiateContractResponse"></a>
=======
<a name="cosmwasm.wasm.v1.MsgInstantiateContractResponse"></a>
>>>>>>> upstream/master

### MsgInstantiateContractResponse
MsgInstantiateContractResponse return instantiation result data


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | Address is the bech32 address of the new contract instance. |
| `data` | [bytes](#bytes) |  | Data contains base64-encoded bytes to returned from the contract |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgMigrateContract"></a>
=======
<a name="cosmwasm.wasm.v1.MsgMigrateContract"></a>
>>>>>>> upstream/master

### MsgMigrateContract
MsgMigrateContract runs a code upgrade/ downgrade for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `code_id` | [uint64](#uint64) |  | CodeID references the new WASM code |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on migration |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgMigrateContractResponse"></a>
=======
<a name="cosmwasm.wasm.v1.MsgMigrateContractResponse"></a>
>>>>>>> upstream/master

### MsgMigrateContractResponse
MsgMigrateContractResponse returns contract migration result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains same raw bytes returned as data from the wasm contract. (May be empty) |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgStoreCode"></a>
=======
<a name="cosmwasm.wasm.v1.MsgStoreCode"></a>
>>>>>>> upstream/master

### MsgStoreCode
MsgStoreCode submit Wasm code to the system


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |
<<<<<<< HEAD
| `source` | [string](#string) |  | Source is a valid absolute HTTPS URI to the contract's source code, optional |
| `builder` | [string](#string) |  | Builder is a valid docker image name with tag, optional |
| `instantiate_permission` | [AccessConfig](#starnamed.x.wasm.v1beta1.AccessConfig) |  | InstantiatePermission access control to apply on contract creation, optional |
=======
| `instantiate_permission` | [AccessConfig](#cosmwasm.wasm.v1.AccessConfig) |  | InstantiatePermission access control to apply on contract creation, optional |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgStoreCodeResponse"></a>
=======
<a name="cosmwasm.wasm.v1.MsgStoreCodeResponse"></a>
>>>>>>> upstream/master

### MsgStoreCodeResponse
MsgStoreCodeResponse returns store result data.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgUpdateAdmin"></a>
=======
<a name="cosmwasm.wasm.v1.MsgUpdateAdmin"></a>
>>>>>>> upstream/master

### MsgUpdateAdmin
MsgUpdateAdmin sets a new admin for a smart contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `sender` | [string](#string) |  | Sender is the that actor that signed the messages |
| `new_admin` | [string](#string) |  | NewAdmin address to be set |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgUpdateAdminResponse"></a>
=======
<a name="cosmwasm.wasm.v1.MsgUpdateAdminResponse"></a>
>>>>>>> upstream/master

### MsgUpdateAdminResponse
MsgUpdateAdminResponse returns empty data





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Msg"></a>
=======
<a name="cosmwasm.wasm.v1.Msg"></a>
>>>>>>> upstream/master

### Msg
Msg defines the wasm Msg service.

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
<<<<<<< HEAD
| `StoreCode` | [MsgStoreCode](#starnamed.x.wasm.v1beta1.MsgStoreCode) | [MsgStoreCodeResponse](#starnamed.x.wasm.v1beta1.MsgStoreCodeResponse) | StoreCode to submit Wasm code to the system | |
| `InstantiateContract` | [MsgInstantiateContract](#starnamed.x.wasm.v1beta1.MsgInstantiateContract) | [MsgInstantiateContractResponse](#starnamed.x.wasm.v1beta1.MsgInstantiateContractResponse) | Instantiate creates a new smart contract instance for the given code id. | |
| `ExecuteContract` | [MsgExecuteContract](#starnamed.x.wasm.v1beta1.MsgExecuteContract) | [MsgExecuteContractResponse](#starnamed.x.wasm.v1beta1.MsgExecuteContractResponse) | Execute submits the given message data to a smart contract | |
| `MigrateContract` | [MsgMigrateContract](#starnamed.x.wasm.v1beta1.MsgMigrateContract) | [MsgMigrateContractResponse](#starnamed.x.wasm.v1beta1.MsgMigrateContractResponse) | Migrate runs a code upgrade/ downgrade for a smart contract | |
| `UpdateAdmin` | [MsgUpdateAdmin](#starnamed.x.wasm.v1beta1.MsgUpdateAdmin) | [MsgUpdateAdminResponse](#starnamed.x.wasm.v1beta1.MsgUpdateAdminResponse) | UpdateAdmin sets a new admin for a smart contract | |
| `ClearAdmin` | [MsgClearAdmin](#starnamed.x.wasm.v1beta1.MsgClearAdmin) | [MsgClearAdminResponse](#starnamed.x.wasm.v1beta1.MsgClearAdminResponse) | ClearAdmin removes any admin stored for a smart contract | |
=======
| `StoreCode` | [MsgStoreCode](#cosmwasm.wasm.v1.MsgStoreCode) | [MsgStoreCodeResponse](#cosmwasm.wasm.v1.MsgStoreCodeResponse) | StoreCode to submit Wasm code to the system | |
| `InstantiateContract` | [MsgInstantiateContract](#cosmwasm.wasm.v1.MsgInstantiateContract) | [MsgInstantiateContractResponse](#cosmwasm.wasm.v1.MsgInstantiateContractResponse) | Instantiate creates a new smart contract instance for the given code id. | |
| `ExecuteContract` | [MsgExecuteContract](#cosmwasm.wasm.v1.MsgExecuteContract) | [MsgExecuteContractResponse](#cosmwasm.wasm.v1.MsgExecuteContractResponse) | Execute submits the given message data to a smart contract | |
| `MigrateContract` | [MsgMigrateContract](#cosmwasm.wasm.v1.MsgMigrateContract) | [MsgMigrateContractResponse](#cosmwasm.wasm.v1.MsgMigrateContractResponse) | Migrate runs a code upgrade/ downgrade for a smart contract | |
| `UpdateAdmin` | [MsgUpdateAdmin](#cosmwasm.wasm.v1.MsgUpdateAdmin) | [MsgUpdateAdminResponse](#cosmwasm.wasm.v1.MsgUpdateAdminResponse) | UpdateAdmin sets a new admin for a smart contract | |
| `ClearAdmin` | [MsgClearAdmin](#cosmwasm.wasm.v1.MsgClearAdmin) | [MsgClearAdminResponse](#cosmwasm.wasm.v1.MsgClearAdminResponse) | ClearAdmin removes any admin stored for a smart contract | |
>>>>>>> upstream/master

 <!-- end services -->



<a name="cosmwasm/wasm/v1/genesis.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/wasm/v1/genesis.proto



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Code"></a>
=======
<a name="cosmwasm.wasm.v1.Code"></a>
>>>>>>> upstream/master

### Code
Code struct encompasses CodeInfo and CodeBytes


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  |  |
<<<<<<< HEAD
| `code_info` | [CodeInfo](#starnamed.x.wasm.v1beta1.CodeInfo) |  |  |
=======
| `code_info` | [CodeInfo](#cosmwasm.wasm.v1.CodeInfo) |  |  |
>>>>>>> upstream/master
| `code_bytes` | [bytes](#bytes) |  |  |
| `pinned` | [bool](#bool) |  | Pinned to wasmvm cache |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Contract"></a>
=======
<a name="cosmwasm.wasm.v1.Contract"></a>
>>>>>>> upstream/master

### Contract
Contract struct encompasses ContractAddress, ContractInfo, and ContractState


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contract_address` | [string](#string) |  |  |
<<<<<<< HEAD
| `contract_info` | [ContractInfo](#starnamed.x.wasm.v1beta1.ContractInfo) |  |  |
| `contract_state` | [Model](#starnamed.x.wasm.v1beta1.Model) | repeated |  |
=======
| `contract_info` | [ContractInfo](#cosmwasm.wasm.v1.ContractInfo) |  |  |
| `contract_state` | [Model](#cosmwasm.wasm.v1.Model) | repeated |  |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.GenesisState"></a>
=======
<a name="cosmwasm.wasm.v1.GenesisState"></a>
>>>>>>> upstream/master

### GenesisState
GenesisState - genesis state of x/wasm


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `params` | [Params](#starnamed.x.wasm.v1beta1.Params) |  |  |
| `codes` | [Code](#starnamed.x.wasm.v1beta1.Code) | repeated |  |
| `contracts` | [Contract](#starnamed.x.wasm.v1beta1.Contract) | repeated |  |
| `sequences` | [Sequence](#starnamed.x.wasm.v1beta1.Sequence) | repeated |  |
| `gen_msgs` | [GenesisState.GenMsgs](#starnamed.x.wasm.v1beta1.GenesisState.GenMsgs) | repeated |  |
=======
| `params` | [Params](#cosmwasm.wasm.v1.Params) |  |  |
| `codes` | [Code](#cosmwasm.wasm.v1.Code) | repeated |  |
| `contracts` | [Contract](#cosmwasm.wasm.v1.Contract) | repeated |  |
| `sequences` | [Sequence](#cosmwasm.wasm.v1.Sequence) | repeated |  |
| `gen_msgs` | [GenesisState.GenMsgs](#cosmwasm.wasm.v1.GenesisState.GenMsgs) | repeated |  |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.GenesisState.GenMsgs"></a>
=======
<a name="cosmwasm.wasm.v1.GenesisState.GenMsgs"></a>
>>>>>>> upstream/master

### GenesisState.GenMsgs
GenMsgs define the messages that can be executed during genesis phase in
order. The intention is to have more human readable data that is auditable.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `store_code` | [MsgStoreCode](#starnamed.x.wasm.v1beta1.MsgStoreCode) |  |  |
| `instantiate_contract` | [MsgInstantiateContract](#starnamed.x.wasm.v1beta1.MsgInstantiateContract) |  |  |
| `execute_contract` | [MsgExecuteContract](#starnamed.x.wasm.v1beta1.MsgExecuteContract) |  |  |
=======
| `store_code` | [MsgStoreCode](#cosmwasm.wasm.v1.MsgStoreCode) |  |  |
| `instantiate_contract` | [MsgInstantiateContract](#cosmwasm.wasm.v1.MsgInstantiateContract) |  |  |
| `execute_contract` | [MsgExecuteContract](#cosmwasm.wasm.v1.MsgExecuteContract) |  |  |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Sequence"></a>
=======
<a name="cosmwasm.wasm.v1.Sequence"></a>
>>>>>>> upstream/master

### Sequence
Sequence key and value of an id generation counter


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `id_key` | [bytes](#bytes) |  |  |
| `value` | [uint64](#uint64) |  |  |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/wasm/v1/ibc.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/wasm/v1/ibc.proto



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgIBCCloseChannel"></a>
=======
<a name="cosmwasm.wasm.v1.MsgIBCCloseChannel"></a>
>>>>>>> upstream/master

### MsgIBCCloseChannel
MsgIBCCloseChannel port and channel need to be owned by the contract


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel` | [string](#string) |  |  |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MsgIBCSend"></a>
=======
<a name="cosmwasm.wasm.v1.MsgIBCSend"></a>
>>>>>>> upstream/master

### MsgIBCSend
MsgIBCSend


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `channel` | [string](#string) |  | the channel by which the packet will be sent |
| `timeout_height` | [uint64](#uint64) |  | Timeout height relative to the current block height. The timeout is disabled when set to 0. |
| `timeout_timestamp` | [uint64](#uint64) |  | Timeout timestamp (in nanoseconds) relative to the current block timestamp. The timeout is disabled when set to 0. |
| `data` | [bytes](#bytes) |  | data is the payload to transfer |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/wasm/v1/proposal.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/wasm/v1/proposal.proto



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.ClearAdminProposal"></a>
=======
<a name="cosmwasm.wasm.v1.ClearAdminProposal"></a>
>>>>>>> upstream/master

### ClearAdminProposal
ClearAdminProposal gov proposal content type to clear the admin of a
contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.InstantiateContractProposal"></a>
=======
<a name="cosmwasm.wasm.v1.InstantiateContractProposal"></a>
>>>>>>> upstream/master

### InstantiateContractProposal
InstantiateContractProposal gov proposal content type to instantiate a
contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `run_as` | [string](#string) |  | RunAs is the address that is passed to the contract's environment as sender |
| `admin` | [string](#string) |  | Admin is an optional address that can execute migrations |
| `code_id` | [uint64](#uint64) |  | CodeID is the reference to the stored WASM code |
| `label` | [string](#string) |  | Label is optional metadata to be stored with a constract instance. |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on instantiation |
| `funds` | [cosmos.base.v1beta1.Coin](#cosmos.base.v1beta1.Coin) | repeated | Funds coins that are transferred to the contract on instantiation |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.MigrateContractProposal"></a>
=======
<a name="cosmwasm.wasm.v1.MigrateContractProposal"></a>
>>>>>>> upstream/master

### MigrateContractProposal
MigrateContractProposal gov proposal content type to migrate a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `run_as` | [string](#string) |  | RunAs is the address that is passed to the contract's environment as sender |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |
| `code_id` | [uint64](#uint64) |  | CodeID references the new WASM code |
| `msg` | [bytes](#bytes) |  | Msg json encoded message to be passed to the contract on migration |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.PinCodesProposal"></a>
=======
<a name="cosmwasm.wasm.v1.PinCodesProposal"></a>
>>>>>>> upstream/master

### PinCodesProposal
PinCodesProposal gov proposal content type to pin a set of code ids in the
wasmvm cache.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `code_ids` | [uint64](#uint64) | repeated | CodeIDs references the new WASM codes |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.StoreCodeProposal"></a>
=======
<a name="cosmwasm.wasm.v1.StoreCodeProposal"></a>
>>>>>>> upstream/master

### StoreCodeProposal
StoreCodeProposal gov proposal content type to submit WASM code to the system


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `run_as` | [string](#string) |  | RunAs is the address that is passed to the contract's environment as sender |
| `wasm_byte_code` | [bytes](#bytes) |  | WASMByteCode can be raw or gzip compressed |
<<<<<<< HEAD
| `source` | [string](#string) |  | Source is a valid absolute HTTPS URI to the contract's source code, optional |
| `builder` | [string](#string) |  | Builder is a valid docker image name with tag, optional |
| `instantiate_permission` | [AccessConfig](#starnamed.x.wasm.v1beta1.AccessConfig) |  | InstantiatePermission to apply on contract creation, optional |
=======
| `instantiate_permission` | [AccessConfig](#cosmwasm.wasm.v1.AccessConfig) |  | InstantiatePermission to apply on contract creation, optional |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.UnpinCodesProposal"></a>
=======
<a name="cosmwasm.wasm.v1.UnpinCodesProposal"></a>
>>>>>>> upstream/master

### UnpinCodesProposal
UnpinCodesProposal gov proposal content type to unpin a set of code ids in
the wasmvm cache.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `code_ids` | [uint64](#uint64) | repeated | CodeIDs references the WASM codes |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.UpdateAdminProposal"></a>
=======
<a name="cosmwasm.wasm.v1.UpdateAdminProposal"></a>
>>>>>>> upstream/master

### UpdateAdminProposal
UpdateAdminProposal gov proposal content type to set an admin for a contract.


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `title` | [string](#string) |  | Title is a short summary |
| `description` | [string](#string) |  | Description is a human readable text |
| `new_admin` | [string](#string) |  | NewAdmin address to be set |
| `contract` | [string](#string) |  | Contract is the address of the smart contract |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->

 <!-- end services -->



<a name="cosmwasm/wasm/v1/query.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## cosmwasm/wasm/v1/query.proto



<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.CodeInfoResponse"></a>
=======
<a name="cosmwasm.wasm.v1.CodeInfoResponse"></a>
>>>>>>> upstream/master

### CodeInfoResponse
CodeInfoResponse contains code meta data from CodeInfo


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | id for legacy support |
| `creator` | [string](#string) |  |  |
| `data_hash` | [bytes](#bytes) |  |  |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryAllContractStateRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryAllContractStateRequest"></a>
>>>>>>> upstream/master

### QueryAllContractStateRequest
QueryAllContractStateRequest is the request type for the
Query/AllContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryAllContractStateResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryAllContractStateResponse"></a>
>>>>>>> upstream/master

### QueryAllContractStateResponse
QueryAllContractStateResponse is the response type for the
Query/AllContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `models` | [Model](#starnamed.x.wasm.v1beta1.Model) | repeated |  |
=======
| `models` | [Model](#cosmwasm.wasm.v1.Model) | repeated |  |
>>>>>>> upstream/master
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryCodeRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryCodeRequest"></a>
>>>>>>> upstream/master

### QueryCodeRequest
QueryCodeRequest is the request type for the Query/Code RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | grpc-gateway_out does not support Go style CodID |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryCodeResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryCodeResponse"></a>
>>>>>>> upstream/master

### QueryCodeResponse
QueryCodeResponse is the response type for the Query/Code RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `code_info` | [CodeInfoResponse](#starnamed.x.wasm.v1beta1.CodeInfoResponse) |  |  |
=======
| `code_info` | [CodeInfoResponse](#cosmwasm.wasm.v1.CodeInfoResponse) |  |  |
>>>>>>> upstream/master
| `data` | [bytes](#bytes) |  |  |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryCodesRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryCodesRequest"></a>
>>>>>>> upstream/master

### QueryCodesRequest
QueryCodesRequest is the request type for the Query/Codes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryCodesResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryCodesResponse"></a>
>>>>>>> upstream/master

### QueryCodesResponse
QueryCodesResponse is the response type for the Query/Codes RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `code_infos` | [CodeInfoResponse](#starnamed.x.wasm.v1beta1.CodeInfoResponse) | repeated |  |
=======
| `code_infos` | [CodeInfoResponse](#cosmwasm.wasm.v1.CodeInfoResponse) | repeated |  |
>>>>>>> upstream/master
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryContractHistoryRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryContractHistoryRequest"></a>
>>>>>>> upstream/master

### QueryContractHistoryRequest
QueryContractHistoryRequest is the request type for the Query/ContractHistory
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract to query |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryContractHistoryResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryContractHistoryResponse"></a>
>>>>>>> upstream/master

### QueryContractHistoryResponse
QueryContractHistoryResponse is the response type for the
Query/ContractHistory RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
<<<<<<< HEAD
| `entries` | [ContractCodeHistoryEntry](#starnamed.x.wasm.v1beta1.ContractCodeHistoryEntry) | repeated |  |
=======
| `entries` | [ContractCodeHistoryEntry](#cosmwasm.wasm.v1.ContractCodeHistoryEntry) | repeated |  |
>>>>>>> upstream/master
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryContractInfoRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryContractInfoRequest"></a>
>>>>>>> upstream/master

### QueryContractInfoRequest
QueryContractInfoRequest is the request type for the Query/ContractInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract to query |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryContractInfoResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryContractInfoResponse"></a>
>>>>>>> upstream/master

### QueryContractInfoResponse
QueryContractInfoResponse is the response type for the Query/ContractInfo RPC
method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
<<<<<<< HEAD
| `contract_info` | [ContractInfo](#starnamed.x.wasm.v1beta1.ContractInfo) |  |  |
=======
| `contract_info` | [ContractInfo](#cosmwasm.wasm.v1.ContractInfo) |  |  |
>>>>>>> upstream/master






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryContractsByCodeRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryContractsByCodeRequest"></a>
>>>>>>> upstream/master

### QueryContractsByCodeRequest
QueryContractsByCodeRequest is the request type for the Query/ContractsByCode
RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `code_id` | [uint64](#uint64) |  | grpc-gateway_out does not support Go style CodID |
| `pagination` | [cosmos.base.query.v1beta1.PageRequest](#cosmos.base.query.v1beta1.PageRequest) |  | pagination defines an optional pagination for the request. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryContractsByCodeResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryContractsByCodeResponse"></a>
>>>>>>> upstream/master

### QueryContractsByCodeResponse
QueryContractsByCodeResponse is the response type for the
Query/ContractsByCode RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `contracts` | [string](#string) | repeated | contracts are a set of contract addresses |
| `pagination` | [cosmos.base.query.v1beta1.PageResponse](#cosmos.base.query.v1beta1.PageResponse) |  | pagination defines the pagination in the response. |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryRawContractStateRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QueryRawContractStateRequest"></a>
>>>>>>> upstream/master

### QueryRawContractStateRequest
QueryRawContractStateRequest is the request type for the
Query/RawContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `query_data` | [bytes](#bytes) |  |  |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QueryRawContractStateResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QueryRawContractStateResponse"></a>
>>>>>>> upstream/master

### QueryRawContractStateResponse
QueryRawContractStateResponse is the response type for the
Query/RawContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains the raw store data |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QuerySmartContractStateRequest"></a>
=======
<a name="cosmwasm.wasm.v1.QuerySmartContractStateRequest"></a>
>>>>>>> upstream/master

### QuerySmartContractStateRequest
QuerySmartContractStateRequest is the request type for the
Query/SmartContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `address` | [string](#string) |  | address is the address of the contract |
| `query_data` | [bytes](#bytes) |  | QueryData contains the query data passed to the contract |






<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.QuerySmartContractStateResponse"></a>
=======
<a name="cosmwasm.wasm.v1.QuerySmartContractStateResponse"></a>
>>>>>>> upstream/master

### QuerySmartContractStateResponse
QuerySmartContractStateResponse is the response type for the
Query/SmartContractState RPC method


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
| `data` | [bytes](#bytes) |  | Data contains the json data returned from the smart contract |





 <!-- end messages -->

 <!-- end enums -->

 <!-- end HasExtensions -->


<<<<<<< HEAD
<a name="starnamed.x.wasm.v1beta1.Query"></a>
=======
<a name="cosmwasm.wasm.v1.Query"></a>
>>>>>>> upstream/master

### Query
Query provides defines the gRPC querier service

| Method Name | Request Type | Response Type | Description | HTTP Verb | Endpoint |
| ----------- | ------------ | ------------- | ------------| ------- | -------- |
<<<<<<< HEAD
| `ContractInfo` | [QueryContractInfoRequest](#starnamed.x.wasm.v1beta1.QueryContractInfoRequest) | [QueryContractInfoResponse](#starnamed.x.wasm.v1beta1.QueryContractInfoResponse) | ContractInfo gets the contract meta data | GET|/wasm/v1beta1/contract/{address}|
| `ContractHistory` | [QueryContractHistoryRequest](#starnamed.x.wasm.v1beta1.QueryContractHistoryRequest) | [QueryContractHistoryResponse](#starnamed.x.wasm.v1beta1.QueryContractHistoryResponse) | ContractHistory gets the contract code history | GET|/wasm/v1beta1/contract/{address}/history|
| `ContractsByCode` | [QueryContractsByCodeRequest](#starnamed.x.wasm.v1beta1.QueryContractsByCodeRequest) | [QueryContractsByCodeResponse](#starnamed.x.wasm.v1beta1.QueryContractsByCodeResponse) | ContractsByCode lists all smart contracts for a code id | GET|/wasm/v1beta1/code/{code_id}/contracts|
| `AllContractState` | [QueryAllContractStateRequest](#starnamed.x.wasm.v1beta1.QueryAllContractStateRequest) | [QueryAllContractStateResponse](#starnamed.x.wasm.v1beta1.QueryAllContractStateResponse) | AllContractState gets all raw store data for a single contract | GET|/wasm/v1beta1/contract/{address}/state|
| `RawContractState` | [QueryRawContractStateRequest](#starnamed.x.wasm.v1beta1.QueryRawContractStateRequest) | [QueryRawContractStateResponse](#starnamed.x.wasm.v1beta1.QueryRawContractStateResponse) | RawContractState gets single key from the raw store data of a contract | GET|/wasm/v1beta1/contract/{address}/raw/{query_data}|
| `SmartContractState` | [QuerySmartContractStateRequest](#starnamed.x.wasm.v1beta1.QuerySmartContractStateRequest) | [QuerySmartContractStateResponse](#starnamed.x.wasm.v1beta1.QuerySmartContractStateResponse) | SmartContractState get smart query result from the contract | GET|/wasm/v1beta1/contract/{address}/smart/{query_data}|
| `Code` | [QueryCodeRequest](#starnamed.x.wasm.v1beta1.QueryCodeRequest) | [QueryCodeResponse](#starnamed.x.wasm.v1beta1.QueryCodeResponse) | Code gets the binary code and metadata for a singe wasm code | GET|/wasm/v1beta1/code/{code_id}|
| `Codes` | [QueryCodesRequest](#starnamed.x.wasm.v1beta1.QueryCodesRequest) | [QueryCodesResponse](#starnamed.x.wasm.v1beta1.QueryCodesResponse) | Codes gets the metadata for all stored wasm codes | GET|/wasm/v1beta1/code|

 <!-- end services -->



<a name="iov/configuration/v1beta1/types.proto"></a>
<p align="right"><a href="#top">Top</a></p>

## iov/configuration/v1beta1/types.proto



<a name="starnamed.x.configuration.v1beta1.Config"></a>

### Config
Config is the configuration of the network


| Field | Type | Label | Description |
| ----- | ---- | ----- | ----------- |
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
=======
| `ContractInfo` | [QueryContractInfoRequest](#cosmwasm.wasm.v1.QueryContractInfoRequest) | [QueryContractInfoResponse](#cosmwasm.wasm.v1.QueryContractInfoResponse) | ContractInfo gets the contract meta data | GET|/cosmwasm/wasm/v1/contract/{address}|
| `ContractHistory` | [QueryContractHistoryRequest](#cosmwasm.wasm.v1.QueryContractHistoryRequest) | [QueryContractHistoryResponse](#cosmwasm.wasm.v1.QueryContractHistoryResponse) | ContractHistory gets the contract code history | GET|/cosmwasm/wasm/v1/contract/{address}/history|
| `ContractsByCode` | [QueryContractsByCodeRequest](#cosmwasm.wasm.v1.QueryContractsByCodeRequest) | [QueryContractsByCodeResponse](#cosmwasm.wasm.v1.QueryContractsByCodeResponse) | ContractsByCode lists all smart contracts for a code id | GET|/cosmwasm/wasm/v1/code/{code_id}/contracts|
| `AllContractState` | [QueryAllContractStateRequest](#cosmwasm.wasm.v1.QueryAllContractStateRequest) | [QueryAllContractStateResponse](#cosmwasm.wasm.v1.QueryAllContractStateResponse) | AllContractState gets all raw store data for a single contract | GET|/cosmwasm/wasm/v1/contract/{address}/state|
| `RawContractState` | [QueryRawContractStateRequest](#cosmwasm.wasm.v1.QueryRawContractStateRequest) | [QueryRawContractStateResponse](#cosmwasm.wasm.v1.QueryRawContractStateResponse) | RawContractState gets single key from the raw store data of a contract | GET|/wasm/v1/contract/{address}/raw/{query_data}|
| `SmartContractState` | [QuerySmartContractStateRequest](#cosmwasm.wasm.v1.QuerySmartContractStateRequest) | [QuerySmartContractStateResponse](#cosmwasm.wasm.v1.QuerySmartContractStateResponse) | SmartContractState get smart query result from the contract | GET|/wasm/v1/contract/{address}/smart/{query_data}|
| `Code` | [QueryCodeRequest](#cosmwasm.wasm.v1.QueryCodeRequest) | [QueryCodeResponse](#cosmwasm.wasm.v1.QueryCodeResponse) | Code gets the binary code and metadata for a singe wasm code | GET|/cosmwasm/wasm/v1/code/{code_id}|
| `Codes` | [QueryCodesRequest](#cosmwasm.wasm.v1.QueryCodesRequest) | [QueryCodesResponse](#cosmwasm.wasm.v1.QueryCodesResponse) | Codes gets the metadata for all stored wasm codes | GET|/cosmwasm/wasm/v1/code|
>>>>>>> upstream/master

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
