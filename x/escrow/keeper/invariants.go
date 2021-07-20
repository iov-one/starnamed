package keeper

import (
	"fmt"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/escrow/types"
)

// RegisterInvariants registers all escrow invariants
func RegisterInvariants(ir sdk.InvariantRegistry, k Keeper) {
	ir.RegisterRoute(types.ModuleName, "attributes", AttributeInvariant(k))
	ir.RegisterRoute(types.ModuleName, "object", ObjectStateInvariant(k))
	ir.RegisterRoute(types.ModuleName, "state", StateInvariant(k))
}

// AllInvariants runs all invariants of the escrow module
func AllInvariants(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		res, stop := AttributeInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		res, stop = ObjectStateInvariant(k)(ctx)
		if stop {
			return res, stop
		}

		return StateInvariant(k)(ctx)
	}
}

//TODO: This is not the job of an Invariant (AttributeInvariant) and should be removed
// it adds overhead (as invariants are run for each block) and should be useless as these conditions
// are checked at genesis or at escrow creation/update

//AttributeInvariant checks that all escrow attributes are correct
func AttributeInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		invalidIDEscrows := 0
		invalidPriceEscrows := 0
		invalidBrokerAddressEscrows := 0
		invalidBrokerCommissionEscrow := 0
		invalidSellerEscrows := 0

		k.IterateEscrows(ctx, func(escrow types.Escrow) bool {
			// Check that the escrow ID is correct
			if types.ValidateID(escrow.Id) != nil ||
				sdk.BigEndianToUint64(types.GetEscrowKey(escrow.Id)) >= k.GetNextIDForExport(ctx) {
				invalidIDEscrows++
				return false
			}

			if types.ValidateAddress(escrow.Seller) != nil {
				invalidSellerEscrows++
				return false
			}

			if types.ValidatePrice(escrow.Price, k.GetEscrowPriceDenom(ctx)) != nil {
				invalidPriceEscrows++
				return false
			}

			if types.ValidateAddress(escrow.BrokerAddress) != nil {
				invalidBrokerAddressEscrows++
				return false
			}

			if types.ValidateCommission(escrow.BrokerCommission) != nil {
				invalidBrokerCommissionEscrow++
				return false
			}

			return false
		})

		broken := invalidIDEscrows+
			invalidPriceEscrows+
			invalidBrokerAddressEscrows+
			invalidSellerEscrows+
			invalidBrokerCommissionEscrow != 0

		return sdk.FormatInvariant(
				types.ModuleName,
				"escrows state",
				fmt.Sprintf("Number of escrows with invalid ID: %v\n"+
					"Number of escrows with invalid price: %v\n"+
					"Number of escrows with invalid seller address: %v\n"+
					"Number of escrows with invalid broker address: %v\n"+
					"Number of escrow with invalid broker commission: %v\n",
					invalidIDEscrows, invalidPriceEscrows, invalidSellerEscrows,
					invalidBrokerAddressEscrows, invalidBrokerCommissionEscrow),
			),
			broken
	}
}

// ObjectStateInvariant checks that all the object of the escrows are in a valid state
func ObjectStateInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		objNotOwnedByModuleEscrows := 0

		k.IterateEscrows(ctx, func(escrow types.Escrow) bool {
			// Check that the object belongs to the module
			obj := escrow.GetObject()
			if ownedByModule, err := obj.IsOwnedBy(k.GetEscrowAddress()); err != nil || !ownedByModule {
				objNotOwnedByModuleEscrows++
			}

			return false
		})

		broken := objNotOwnedByModuleEscrows != 0

		return sdk.FormatInvariant(
				types.ModuleName,
				"escrows object",
				fmt.Sprintf(
					"Number of escrows with non-existing objects: %v",
					objNotOwnedByModuleEscrows),
			),
			broken
	}
}

// StateInvariant checks that all escrows are in a valid state
func StateInvariant(k Keeper) sdk.Invariant {
	return func(ctx sdk.Context) (string, bool) {
		invalidExpirationEscrows := 0
		completedEscrows := 0
		refundedEscrows := 0
		invalidDeadlineEscrows := 0

		k.IterateEscrows(ctx, func(escrow types.Escrow) bool {

			date := uint64(ctx.BlockTime().Unix())

			// Check that the deadline does not exceeds maximum duration as if it was created now
			if escrow.Deadline > date+uint64(k.GetMaximumEscrowDuration(ctx).Seconds()) {
				invalidDeadlineEscrows++
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

			return false
		})

		broken := invalidExpirationEscrows+
			completedEscrows+
			refundedEscrows+
			invalidDeadlineEscrows != 0

		return sdk.FormatInvariant(
				types.ModuleName,
				"escrows state",
				fmt.Sprintf("Number of escrows with invalid expiration status : %v\n"+
					"Number of escrows with invalid deadline: %v\n"+
					"Number of escrows with invalid state (refunded) : %v\n"+
					"Number of escrows with invalid state (completed): %v",
					invalidExpirationEscrows, invalidDeadlineEscrows, refundedEscrows, completedEscrows,
				),
			),
			broken
	}
}
