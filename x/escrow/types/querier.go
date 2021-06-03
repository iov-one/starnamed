package types

const (
	QueryEscrow = "escrow" // query an escrow
)

//QueryEscrowParams defines the params to query an escrow
type QueryEscrowParams struct {
	Id string
}
