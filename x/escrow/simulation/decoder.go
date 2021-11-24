package simulation

import (
	"bytes"
	"fmt"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/types"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/types/kv"
)

// NewDecodeStore unmarshals the KVPair's Value to the corresponding escrow type
func NewDecodeStore(cdc codec.Codec) func(kvA, kvB kv.Pair) string {
	return func(kvA, kvB kv.Pair) string {
		switch {
		case bytes.Equal(kvA.Key[:1], keeper.EscrowStoreKey):
			var escrow1, escrow2 types.Escrow
			cdc.MustUnmarshal(kvA.Value, &escrow1)
			cdc.MustUnmarshal(kvB.Value, &escrow2)
			return fmt.Sprintf("%v\n%v", escrow1, escrow2)

		case bytes.Equal(kvA.Key[:1], keeper.DeadlineStoreKey):
			return fmt.Sprintf("%v\n%v", kvA.Value, kvB.Value)
		case bytes.Equal(kvA.Key[:1], keeper.ParamsStoreKey):
			//TODO: factor in parameter name
			return fmt.Sprintf("%v\n%v", kvA.Value, kvB.Value)
		default:
			panic(fmt.Sprintf("invalid escrow key prefix %X", kvA.Key[:1]))
		}
	}
}
