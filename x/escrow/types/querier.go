package types

import codectypes "github.com/cosmos/cosmos-sdk/codec/types"

const (
	QueryEscrow  = "escrow"  // query an escrow
	QueryEscrows = "escrows" // query multiple escrows
)

// QueryEscrowParams defines the params to query an escrow
type QueryEscrowParams struct {
	Id string
}

// QueryEscrowsParams defines the parameters to query multiple escrows
type QueryEscrowsParams struct {
	Seller                            string
	State                             string
	ObjectKey                         string
	PaginationStart, PaginationLength uint64
}

// UnpackInterfaces make sure the Anys included in QueryEscrowResponse are unpacked (e.g the object field)
func (q *QueryEscrowResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if q.Escrow != nil {
		return q.Escrow.UnpackInterfaces(unpacker)
	}

	return nil
}

// UnpackInterfaces make sure the Anys included in the QueryEscrowsResponse are unpacked (e.g the object field)
func (q *QueryEscrowsResponse) UnpackInterfaces(unpacker codectypes.AnyUnpacker) error {
	if q.Escrows != nil {
		for i := range q.Escrows {
			if err := q.Escrows[i].UnpackInterfaces(unpacker); err != nil {
				return err
			}
		}
	}

	return nil
}
