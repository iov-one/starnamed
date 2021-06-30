package types_test

import (
	"errors"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

func TestValidate(t *testing.T) {

	test.SetConfig()
	gen := test.NewEscrowGenerator(100)

	defaultId := "0123456789abcdef"
	defaultSellerBytes := gen.NewAccAddress()
	defaultSeller := defaultSellerBytes.String()
	defaultObj := gen.NewTestObject(defaultSellerBytes)
	defaultDeadline := gen.NowAfter(10)
	defaultPrice := sdk.NewCoins(sdk.NewCoin("tiov", sdk.NewInt(50)))
	negativePrice := sdk.Coins{
		sdk.Coin{
			Denom:  "tiov",
			Amount: sdk.NewInt(-20),
		},
	}

	testCases := []struct {
		name     string
		id       string
		seller   string
		state    types.EscrowState
		price    sdk.Coins
		obj      types.TransferableObject
		deadline uint64
	}{
		{
			name: "valid escrow: simple price",
		},
		{
			name:  "valid escrow: composed price",
			price: defaultPrice.Add(sdk.NewCoin("tiov2", sdk.NewInt(20))),
		},
		{
			name:  "invalid escrow: negative price",
			price: negativePrice,
		},
		{
			name: "invalid escrow: non-hexadecimal id",
			id:   "123456789abcdefg",
		},
		{
			name: "invalid escrow: id too short",
			id:   "abcdef",
		},
		{
			name: "invalid escrow: id too long",
			id:   "0123456789abcdefedcba",
		},
		{
			name:   "invalid escrow: invalid bech32",
			seller: "star14894684ded56f",
			obj:    gen.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid escrow: bech32 of another network",
			seller: "cosmos1cqfse93m6r7fr3vx07du5yfmsltca60gyadygf",
			obj:    gen.NewNotPossessedTestObject(),
		},
	}

	for _, tc := range testCases {
		id := tc.id
		if len(id) == 0 {
			id = defaultId
		}
		seller := tc.seller
		if len(seller) == 0 {
			seller = defaultSeller
		}
		obj := tc.obj
		if obj == nil {
			obj = defaultObj
		}
		price := tc.price
		if price == nil {
			price = defaultPrice
		}
		deadline := tc.deadline
		if deadline == 0 {
			deadline = defaultDeadline
		}
		state := tc.state

		test.EvaluateTest(t, tc.name, func(t *testing.T) error {
			escrow := types.Escrow{
				Id:       id,
				Seller:   seller,
				Object:   test.MustPackToAny(obj),
				Price:    price,
				State:    state,
				Deadline: deadline,
			}
			err1 := escrow.Validate(gen.NowAfter(0))
			err2 := escrow.ValidateWithoutDeadlineAndObject()
			if !errors.Is(err1, err2) {
				t.Fatalf("Error, mismatch of validation error between Validate and ValidateWithoutDeadlineAndObject : %v and %v",
					err1,
					err2,
				)
			}
			return err1
		})
	}

	// Test with deadline and owner
	test.EvaluateTest(t, "invalid escrow: object does not belong to seller", func(*testing.T) error {
		return types.Escrow{
			Id:       defaultId,
			Seller:   gen.NewAccAddress().String(),
			Object:   test.MustPackToAny(defaultObj),
			Price:    defaultPrice,
			Deadline: defaultDeadline,
		}.Validate(gen.NowAfter(0))
	})
	test.EvaluateTest(t, "invalid escrow: passed deadline", func(*testing.T) error {
		return types.Escrow{
			Id:       defaultId,
			Seller:   defaultSeller,
			Object:   test.MustPackToAny(defaultObj),
			Price:    defaultPrice,
			Deadline: gen.NowAfter(0) - 20,
		}.Validate(gen.NowAfter(0))
	})
}
