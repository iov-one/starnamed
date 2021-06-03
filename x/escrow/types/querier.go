package types

import (
	tmbytes "github.com/tendermint/tendermint/libs/bytes"
)

const (
	QueryEscrow = "escrow" // query an escrow
)

//QueryEscrowParams defines the params to query an escrow
type QueryEscrowParams struct {
	//TODO: shouldn't this be a string ?
	ID tmbytes.HexBytes
}
