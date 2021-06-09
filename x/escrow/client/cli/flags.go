// nolint
package cli

import (
	flag "github.com/spf13/pflag"
)

const (
	FlagBuyer    = "buyer"
	FlagSeller   = "seller"
	FlagPrice    = "price"
	FlagDeadline = "expiration"
)

var (
	FsEscrow = flag.NewFlagSet("escrow", flag.PanicOnError)
)

func init() {
	FsEscrow.String(FlagSeller, "", "Bech32 encoding address of the new seller for the escrow")
	FsEscrow.String(FlagBuyer, "", "Bech32 encoding address of the buyer")
	FsEscrow.String(FlagPrice, "", "Price of the object")
	FsEscrow.String(FlagDeadline, "", "Expiration date of the escrow, in the RFC3339 time format")
}
