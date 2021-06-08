package types

import (
	"encoding/hex"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/errors"
)

// TODO: not sure if this is the correct place
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

var (
	//TODO: move this in the keeper, nothing to do there

	// Keys for store prefixes
	EscrowStoreKey   = []byte{0x01} // prefix for escrow
	DeadlineStoreKey = []byte{0x02} // prefix for escrow stored by expiration date
	ParamsStoreKey   = []byte{0x03} // prefix for the keeper parameters

	// Keys for the parameters store
	ParamsStoreLastBlockTime = []byte{0x01}
	ParamsStoreNextId        = []byte{0x02}
)

func GetEscrowKey(id string) []byte {
	key, err := hex.DecodeString(id)
	if err != nil {
		//TODO: should we panic here ?
		panic(errors.Wrap(err, "Invalid escrow key format"))
	}
	return key
}

func GetDeadlineKey(deadline uint64, id string) []byte {
	return append(sdk.Uint64ToBigEndian(deadline), GetEscrowKey(id)...)
}
