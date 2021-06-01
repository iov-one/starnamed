package types

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
	// Keys for store prefixes
	EscrowKey = []byte{0x01} // prefix for escrow
)

// GetEscrowKey returns the key for the escrow with the specified id
func GetEscrowKey(id []byte) []byte {
	return append(EscrowKey, id...)
}
