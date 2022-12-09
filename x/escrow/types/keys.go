package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

const (
	// ModuleName is the name of the escrow module
	ModuleName = "escrow"

	// StoreKey is the string store representation
	StoreKey string = ModuleName

	// QuerierRoute is the querier route for the escrow module
	QuerierRoute string = ModuleName

	// RouterKey is the msg router key for the escrow module
	RouterKey string = ModuleName
)

// GetEscrowKey returns a byte array that can be used as a unique key from an escrow id, used in the escrow store
func GetEscrowKey(id string) []byte {
	key, err := hex.DecodeString(id)
	if err != nil {
		panic(errors.Wrap(err, "Invalid escrow key format"))
	}
	return key
}

// GetDeadlineKey returns a byte array that can be used as a unique key from an escrow id and its deadline,
// prefixing the deadline so it can be used to iterate through escrows by deadline
func GetDeadlineKey(deadline uint64, id string) []byte {
	return append(sdk.Uint64ToBigEndian(deadline), GetEscrowKey(id)...)
}

// This is not used and can be used if we want to query escrow by object differentiating between object types

// GetEscrowObjectKey returns a byte array that can be used as a unique key from a TransferableObject
func GetEscrowObjectKey(obj TransferableObject) []byte {
	return append(sdk.Uint64ToBigEndian(uint64(obj.GetObjectTypeID())), obj.GetUniqueKey()...)
}

// ContructEscrowObjectKey returns a byte array that can be used as a unique key from a TransferableObject
func ContructEscrowObjectKey(id TypeID, key []byte) []byte {
	return append(sdk.Uint64ToBigEndian(uint64(id)), key...)
}
