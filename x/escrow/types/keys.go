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

func GetEscrowKey(id string) []byte {
	key, err := hex.DecodeString(id)
	if err != nil {
		panic(errors.Wrap(err, "Invalid escrow key format"))
	}
	return key
}

func GetDeadlineKey(deadline uint64, id string) []byte {
	return append(sdk.Uint64ToBigEndian(deadline), GetEscrowKey(id)...)
}
