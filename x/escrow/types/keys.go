package types

import sdk "github.com/cosmos/cosmos-sdk/types"

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

	// EscrowStoreKey Key for store prefixes
	EscrowStoreKey   = []byte{0x01} // prefix for escrow
	DeadlineStoreKey = []byte{0x02} // prefix for escrow stored by expiration date
)

func GetEscrowKey(id string) []byte {
	//TODO: shouldn't this be hex-decoding ?
	return []byte(id)
}

func GetDeadlineKey(deadline uint64, id string) []byte {
	return append(sdk.Uint64ToBigEndian(deadline), GetEscrowKey(id)...)
}
