package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

const (
	QueryEscrow = "escrow" // query an escrow
)

//QueryEscrowParams defines the params to query an escrow
type QueryEscrowParams struct {
	ID tmbytes.HexBytes
}
