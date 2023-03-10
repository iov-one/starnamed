// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSeller           = "seller"
	FlagPrice            = "price"
	FlagDeadline         = "expiration"
	FlagFeePayer         = "fee-payer"
	FlagObjectKey        = "object"
	FlagState            = "state"
	FlagPaginationStart  = "pagination-start"
	FlagPaginationLength = "pagination-length"
)

var (
	FsEscrow       = flag.NewFlagSet("escrow", flag.PanicOnError)
	FsQueryEscrows = flag.NewFlagSet("query_escrows", flag.PanicOnError)
)

func init() {
	FsEscrow.String(FlagSeller, "", "Bech32 encoded address of the new seller for the escrow")
	FsEscrow.String(FlagPrice, "", "Price of the object")
	FsEscrow.String(FlagDeadline, "", "Expiration date of the escrow, in the RFC3339 time format")
	addCommonFlags(FsEscrow)

	FsQueryEscrows.String(FlagSeller, "", "Bech32 encoded address of the seller of the escrow")
	FsQueryEscrows.String(FlagState, "", "State of the escrow, can be open or expired")
	FsQueryEscrows.String(FlagObjectKey, "", "Primary key of the escrow's object, encoded in hexadecimal")
	FsQueryEscrows.Uint64(FlagPaginationStart, 0, "Pagination starting index")
	FsQueryEscrows.Uint64(FlagPaginationLength, 0, "Maximal number of escrows to fetch, 0 to fetch them all")

}

func addCommonFlags(flagSet *flag.FlagSet) {
	flagSet.String(FlagFeePayer, "", "Bech32 encoded address ot the fee payer for this operation")
}
