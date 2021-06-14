package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"
	tmbytes "github.com/tendermint/tendermint/libs/bytes"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// RegisterInvariants registers all escrow invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "state", StateInvariant(k))
}

// AllInvariants runs all invariants of the escrow module
func AllInvariants(k Keeper) sdk.Invariant {
	return StateInvariant(k)
}

// StateInvariant checks that all escrows are in a valid state
func StateInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		invalidEscrows := 0
		k.IterateEscrows(ctx, func(bytes tmbytes.HexBytes, escrow types.Escrow) bool {

			date := uint64(ctx.BlockTime().Unix())
			// Check that escrow is expired iff deadline is passed
			if (escrow.State == types.EscrowState_Expired) != (escrow.Deadline >= date) {
				invalidEscrows++
				return false
			}
			// Check that the escrow is in a valid state
			if escrow.State == types.EscrowState_Refunded || escrow.State == types.EscrowState_Completed {
				invalidEscrows++
				return false
			}

			// Check that the object belongs to the module
			obj := escrow.GetObject()
			store, err := k.getStoreForID(ctx, obj.GetType())
			if err != nil {
				invalidEscrows++
				return false
			}

			err = k.checkObjectWithStore(store, obj)
			if err != nil {
				invalidEscrows++
				return false
			}

			if ownedByModule, err := obj.IsOwnedBy(k.GetEscrowAccount(ctx).GetAddress()); err != nil || !ownedByModule {
				invalidEscrows++
			}

			return false
		})

		broken := invalidEscrows != 0

		//TODO: print a more useful message
		return sdk.FormatInvariant(types.ModuleName, "escrows state", fmt.Sprintf("Number of invalid escrows : %v", invalidEscrows)), broken
	}
}
