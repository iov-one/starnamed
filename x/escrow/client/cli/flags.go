// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagSeller   = "seller"
	FlagPrice    = "price"
	FlagDeadline = "expiration"
	FlagFeePayer = "fee-payer"
)

var (
	FsEscrow = flag.NewFlagSet("escrow", flag.PanicOnError)
)

func init() {
	FsEscrow.String(FlagSeller, "", "Bech32 encoded address of the new seller for the escrow")
	FsEscrow.String(FlagPrice, "", "Price of the object")
	FsEscrow.String(FlagDeadline, "", "Expiration date of the escrow, in the RFC3339 time format")
	addCommonFlags(FsEscrow)
}

func addCommonFlags(flagSet *flag.FlagSet) {
	flagSet.String(FlagFeePayer, "", "Bech32 encoded address ot the fee payer for this operation")
}
