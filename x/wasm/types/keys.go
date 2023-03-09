package types

const (
	// ModuleName defines the module name
	ModuleName = "wasm"

	// StoreKey defines the primary module store key
	StoreKey = ModuleName

	// RouterKey is the message route for slashing
	RouterKey = ModuleName

    // QuerierRoute defines the module's query routing key
    QuerierRoute = ModuleName

	// MemStoreKey defines the in-memory store key
	MemStoreKey = "mem_wasm"

    
)

<<<<<<< HEAD
=======
// nolint
var (
	CodeKeyPrefix                                  = []byte{0x01}
	ContractKeyPrefix                              = []byte{0x02}
	ContractStorePrefix                            = []byte{0x03}
	SequenceKeyPrefix                              = []byte{0x04}
	ContractCodeHistoryElementPrefix               = []byte{0x05}
	ContractByCodeIDAndCreatedSecondaryIndexPrefix = []byte{0x06}
	PinnedCodeIndexPrefix                          = []byte{0x07}
	TXCounterPrefix                                = []byte{0x08}
>>>>>>> tags/v0.11.6


<<<<<<< HEAD
func KeyPrefix(p string) []byte {
    return []byte(p)
=======
// GetCodeKey constructs the key for retreiving the ID for the WASM code
func GetCodeKey(codeID uint64) []byte {
	contractIDBz := sdk.Uint64ToBigEndian(codeID)
	return append(CodeKeyPrefix, contractIDBz...)
}

// GetContractAddressKey returns the key for the WASM contract instance
func GetContractAddressKey(addr sdk.AccAddress) []byte {
	return append(ContractKeyPrefix, addr...)
}

// GetContractStorePrefix returns the store prefix for the WASM contract instance
func GetContractStorePrefix(addr sdk.AccAddress) []byte {
	return append(ContractStorePrefix, addr...)
}

// GetContractByCreatedSecondaryIndexKey returns the key for the secondary index:
// `<prefix><codeID><created/last-migrated><contractAddr>`
func GetContractByCreatedSecondaryIndexKey(contractAddr sdk.AccAddress, c ContractCodeHistoryEntry) []byte {
	prefix := GetContractByCodeIDSecondaryIndexPrefix(c.CodeID)
	prefixLen := len(prefix)
	contractAddrLen := len(contractAddr)
	r := make([]byte, prefixLen+AbsoluteTxPositionLen+contractAddrLen)
	copy(r[0:], prefix)
	copy(r[prefixLen:], c.Updated.Bytes())
	copy(r[prefixLen+AbsoluteTxPositionLen:], contractAddr)
	return r
}

// GetContractByCodeIDSecondaryIndexPrefix returns the prefix for the second index: `<prefix><codeID>`
func GetContractByCodeIDSecondaryIndexPrefix(codeID uint64) []byte {
	prefixLen := len(ContractByCodeIDAndCreatedSecondaryIndexPrefix)
	const codeIDLen = 8
	r := make([]byte, prefixLen+codeIDLen)
	copy(r[0:], ContractByCodeIDAndCreatedSecondaryIndexPrefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(codeID))
	return r
}

// GetContractCodeHistoryElementKey returns the key a contract code history entry: `<prefix><contractAddr><position>`
func GetContractCodeHistoryElementKey(contractAddr sdk.AccAddress, pos uint64) []byte {
	prefix := GetContractCodeHistoryElementPrefix(contractAddr)
	prefixLen := len(prefix)
	r := make([]byte, prefixLen+8)
	copy(r[0:], prefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(pos))
	return r
}

// GetContractCodeHistoryElementPrefix returns the key prefix for a contract code history entry: `<prefix><contractAddr>`
func GetContractCodeHistoryElementPrefix(contractAddr sdk.AccAddress) []byte {
	prefixLen := len(ContractCodeHistoryElementPrefix)
	contractAddrLen := len(contractAddr)
	r := make([]byte, prefixLen+contractAddrLen)
	copy(r[0:], ContractCodeHistoryElementPrefix)
	copy(r[prefixLen:], contractAddr)
	return r
}

// GetPinnedCodeIndexPrefix returns the key prefix for a code id pinned into the wasmvm cache
func GetPinnedCodeIndexPrefix(codeID uint64) []byte {
	prefixLen := len(PinnedCodeIndexPrefix)
	r := make([]byte, prefixLen+8)
	copy(r[0:], PinnedCodeIndexPrefix)
	copy(r[prefixLen:], sdk.Uint64ToBigEndian(codeID))
	return r
}

// ParsePinnedCodeIndex converts the serialized code ID back.
func ParsePinnedCodeIndex(s []byte) uint64 {
	return sdk.BigEndianToUint64(s)
>>>>>>> tags/v0.11.6
}
