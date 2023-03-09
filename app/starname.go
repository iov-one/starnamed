package app

import (
	"fmt"
	"strconv"
)

func parseCoinType() uint32 {
	if parsed, err := strconv.ParseInt(CoinTypeStr, 10, 32); err != nil {
		panic(err)
	} else {
		return uint32(parsed)
	}
}

var (
	CoinTypeStr string = "234"
	// CoinType is the type as defined in SLIP44 (https://github.com/satoshilabs/slips/blob/master/slip-0044.md)
	CoinType = parseCoinType()
	// FullFundraiserPath is the parts of the BIP44 HD path
	FullFundraiserPath = fmt.Sprintf("m/44'/%d'/0'/0/0", CoinType)
)
