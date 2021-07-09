package types

const (
	QueryEscrow  = "escrow"  // query an escrow
	QueryEscrows = "escrows" // query multiple escrows
)

//QueryEscrowParams defines the params to query an escrow
type QueryEscrowParams struct {
	Id string
}

//QueryEscrowsParams defines the parameters to query multiple escrows
type QueryEscrowsParams struct {
	Seller                            string
	State                             string
	ObjectKey                         string
	PaginationStart, PaginationLength uint64
}
