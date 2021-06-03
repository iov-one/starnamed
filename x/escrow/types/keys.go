package types

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
	EscrowStoreKey = []byte{0x01} // prefix for escrow
)
