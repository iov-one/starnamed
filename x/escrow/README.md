# escrow

A trustless swap implementation built as a Cosmos SDK module

## Description

This module provides swap functionality between objects and tokens without having to trust a third party.

It's a generic module, that can be used to sell objects for tokens reliably, within any blockchain of the Cosmos SDK. 
It requires the [cosmos-sdk-crud][cosmos-sdk-crud-github] package and has, for now, a small dependency 
on the `starnamed/x/configuration` module.

## Concepts

This module creates and manages `escrow`s, which are structures that own an object that can be bought for a certain price.
Each escrow has a unique 16 characters ID and contains the price of the object, an expiration date (a deadline), and the seller address.

The object held by an escrow must implement the `TransferableObject` interface, which specifies how an object can be transferred from one owner to another. Every object that implements that interface must be a protocol buffer message (implementing `codec.ProtoMarshaler`), must have a type identifier, must have a unique key in the domain of the objects of the same type, and must provide the `IsOwnedBy` and `Transfer` functions that provides the ownership transfer logic.

## Integration in an app

To start using escrow, you just need to have an object that satisfies the `TransferableObject` interface. This object has to be registered in the application codec, which is done via the `RegisterInterfaces` method in case of a module : 
```go
    registry.RegisterImplementations(
		(*escrowtypes.TransferableObject)(nil),
		&MyTransferableObject{},
	)
```
Now an escrow can be created with this object, and you can provide the end-users a custom CLI tx command to create escrows holding this type of object. The escrow module does not contain a creation command but it offers utility functions to help create such command.

An example of a creation function (typically used in the `RunE` field of a `cobra.Command` definition) is shown below :
```go
  func(cmd *cobra.Command, args []string) (err error) {
    clientCtx, err := client.GetClientTxContext(cmd)
    if err != nil {
    return err
    }
    
    // Retrieve or create the transferableObject from the cmd.Flags() or/and the args array
    transferableObject := ...
    
    // Create the message from the object and the command flags
    msg, err := escrowcli.NewMsgCreateEscrow(clientCtx, cmd, transferableObject)
    if err != nil {
    return err
    }
    // broadcast request
    return tx.GenerateOrBroadcastTxCLI(clientCtx, cmd.Flags(), msg)
  }
```
Don't forget to add the required flags to the command flag set using the `escrowcli.AddCreateEscrowFlags(cmd)` helper function.

## Usage

The escrow module expose 4 simple methods to manage escrows : `CreateEscrow`, `UpdateEscrow`, `TransferToEscrow` and `RefundEscrow`
```go
    // CreateEscrow creates an escrow and transfer the object to the escrow account.
    // The deadline must be included in the interval ]now, now + escrow_max_period].
    // The price must be in fee_coin_denom denomination.
    // The escrow is created with the predefined escrow_broker broker address and
    // escrow_commission commission.
    // The returned string is the 16 character escrow ID
    CreateEscrow(ctx sdk.Context, seller sdk.AccAddress, price sdk.Coins, object types.TransferableObject, deadline uint64) 
        (string, error)

    // TransferToEscrow transfers coins from the buyer to the escrow account.
    // The specified amount must be greater than or equal to the escrow price.
    // The actual transferred coins match the price of the escrow, the amount provided there is just a security to limit
    // the coins the buyer accepts to spend.
    // The coins will be transferred to the escrow account and then the object is transferred to the buyer and the coins
    // are sent to the seller. The escrow is then marked as completed and removed.
    // If the object or the coin transfer from the escrow account fail, this function panics.
    TransferToEscrow(ctx sdk.Context, buyer sdk.AccAddress,  id string,  amount sdk.Coins) 
        error
    
    // UpdateEscrow perform an escrow update, the updater must be the current escrow owner (seller)
    // The escrow must be in the open state to be updated.
    // If no changes are to be made for a specific parameter, it must be a zero value.
    // An empty update (with all parameter being nil/zero) will fail and return an error.
    // The new deadline must be in the interval [oldDeadline; now + max_escrow_period].
    UpdateEscrow(ctx sdk.Context, id string, updater sdk.AccAddress, newSeller sdk.AccAddress, newPrice sdk.Coins, newDeadline uint64) 
        error
    
    // RefundEscrow refunds the specified escrow, returning the object to the seller and removing the escrow.
    // An escrow can only be refunded by its owner (the seller) or by anybody when it is expired
    RefundEscrow(ctx sdk.Context, sender sdk.AccAddress, id string) 
        error 
```
These methods are accessible as-is from the escrow keeper, or can be called from a gRPC or CLI query for an end user (except the create method which is not available as an escrow CLI command, as the escrow is agnostic of the object it owns, so it cannot create or retrieve one)

Additionally, the keeper offer various query commands, and two gRPC/REST/CLI queries :
* Single escrow query : queries an escrow by its unique ID | `Escrow` / `GET /escrow/escrow/{id}` / `query escrow escrow [id]`
* Multiple escrow query : queries escrows by their attributes (seller, object key, state), if an attribute is not specified then no filtering is done for this attribute | `Escrows` / `GET /escrow/escrows?seller={}&state={}&object={}` / `query escrow escrows [--seller seller][--object objectKey][--state open|expired]`

## Further customization

You can further customize the behavior of the escrow module, and how it handles objects.

You can register custom data (typically in your module keeper constructor or in the app initialization function) that will be passed as argument to your object `Transfer` method using the `keeper.RegisterCustomData(TypeID, CustomData)` method.
You can register a unique custom structure per type ID, calls to the `RegisterCustomData` method will overwrite any previously defined data for this type ID.

If you want to customize creation fees for an object, it can implement the `ObjectWithCustomFees` interface. When creating an object, the default creation fees will be used if the object does not implement this interface, otherwise the value returned by `GetCreationFees()` will be used.

If you want an object to be able to validate the deadline of the escrow, it can implement the `ObjectWithTimeConstraint` interface.
If an object implement this interface, the `ValidateDeadline` method is called as an additional check for the escrow expiration upon creation and update.

[cosmos-sdk-crud-github]:https://github.com/iov-one/cosmos-sdk-crud
