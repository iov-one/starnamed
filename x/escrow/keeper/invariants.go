package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

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
		invalidIDEscrows := 0
		invalidPriceEscrows := 0
		invalidExpirationEscrows := 0
		completedEscrows := 0
		refundedEscrows := 0
		objNotExistingEscrows := 0
		syncErrWithStoreEscrows := 0
		objNotOwnedByModuleEscrows := 0

		k.IterateEscrows(ctx, func(id string, escrow types.Escrow) bool {

			date := uint64(ctx.BlockTime().Unix())

			// Check that the escrow ID is correct
			if types.ValidateID(escrow.Id) != nil ||
				sdk.BigEndianToUint64(types.GetEscrowKey(escrow.Id)) >= k.GetNextIDForExport(ctx) {
				invalidIDEscrows++
				return false
			}

			// Check that escrow is expired iff deadline is passed
			if (escrow.State == types.EscrowState_Expired) != (date >= escrow.Deadline) {
				invalidExpirationEscrows++
				return false
			}
			// Check that the escrow is in a valid state
			if escrow.State == types.EscrowState_Refunded {
				refundedEscrows++
				return false
			} else if escrow.State == types.EscrowState_Completed {
				completedEscrows++
				return false
			}

			if types.ValidatePrice(escrow.Price) != nil {
				invalidPriceEscrows++
				return false
			}
			// Check that the object belongs to the module
			obj := escrow.GetObject()
			store, err := k.getStoreForID(ctx, obj.GetType())
			if err != nil {
				objNotExistingEscrows++
				return false
			}

			err = k.checkObjectWithStore(store, obj)
			if err != nil {
				syncErrWithStoreEscrows++
				return false
			}

			if ownedByModule, err := obj.IsOwnedBy(k.GetEscrowAddress()); err != nil || !ownedByModule {
				objNotOwnedByModuleEscrows++
			}

			return false
		})

		broken := invalidIDEscrows+
			invalidPriceEscrows+
			invalidExpirationEscrows+
			completedEscrows+
			refundedEscrows+
			objNotExistingEscrows+
			syncErrWithStoreEscrows+
			objNotOwnedByModuleEscrows != 0

		return sdk.FormatInvariant(
				types.ModuleName,
				"escrows state",
				fmt.Sprintf("Number of escrows with invalid ID: %v\n"+
					"Number of escrows with invalid price: %v\n"+
					"Number of escrows with invalid expiration status : %v\n"+
					"Number of escrows with invalid state (refunded) : %v\n"+
					"Number of escrows with invalid state (completed): %v\n"+
					"Number of escrows with non-existing objects: %v\n"+
					"Number of escrows with objects not in sync with store : %v\n"+
					"Number of escrows with objects not owned by module: %v",
					invalidIDEscrows, invalidPriceEscrows, invalidExpirationEscrows, refundedEscrows,
					completedEscrows, objNotExistingEscrows, syncErrWithStoreEscrows, objNotOwnedByModuleEscrows),
			),
			broken
	}
}
