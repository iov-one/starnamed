package mock

import (
	"fmt"
	"os"

	"github.com/cosmos/cosmos-sdk/crypto/hd"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
)

func Addresses() (sdk.AccAddress, sdk.AccAddress) {
	keyBase := keyring.NewInMemory()
	addr1, _, err := keyBase.NewMnemonic("alice", keyring.English, "", "", hd.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	addr2, _, err := keyBase.NewMnemonic("bob", keyring.English, "", "", hd.Secp256k1)
	if err != nil {
		fmt.Println("unable to generate mock addresses " + err.Error())
		os.Exit(1)
	}
	return addr1.GetAddress(), addr2.GetAddress()
}
