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

	defaultBroker := gen.NewAccAddress().String()
	defaultCommission := sdk.ZeroDec()
	defaultId := "0123456789abcdef"
	defaultSellerBytes := gen.NewAccAddress()
	defaultSeller := defaultSellerBytes.String()
	defaultObj := gen.NewTestObject(defaultSellerBytes)
	defaultDeadline := gen.NowAfter(10)
	defaultPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	negativePrice := sdk.Coins{
		sdk.Coin{
			Denom:  test.Denom,
			Amount: sdk.NewInt(-20),
		},
	}
	denom := test.Denom

	testCases := []struct {
		name       string
		id         string
		seller     string
		state      types.EscrowState
		price      sdk.Coins
		obj        types.TransferableObject
		deadline   uint64
		commission sdk.Dec
		broker     string
	}{
		{
			name: "valid escrow: simple price",
		},
		{
			name:  "invalid escrow: composed price",
			price: defaultPrice.Add(sdk.NewCoin(test.DenomAux, sdk.NewInt(20))),
		},
		{
			name:  "invalid escrow: negative price",
			price: negativePrice,
		},
		{
			name:  "invalid escrow: empty price",
			price: sdk.Coins{},
		},
		{
			name:  "invalid escrow: invalid denom price",
			price: sdk.NewCoins(sdk.NewCoin("abcde", sdk.OneInt())),
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
			name:   "invalid seller: invalid bech32",
			seller: "star14894684ded56f",
			obj:    gen.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid seller: bech32 of another network",
			seller: "cosmos1cqfse93m6r7fr3vx07du5yfmsltca60gyadygf",
			obj:    gen.NewNotPossessedTestObject(),
		},
		{
			name:   "invalid broker: invalid bech32",
			broker: "star14894684ded56f",
		},
		{
			name:   "invalid broker: bech32 of another network",
			broker: "cosmos1cqfse93m6r7fr3vx07du5yfmsltca60gyadygf",
		},
		{
			name:       "invalid commission: negative",
			commission: sdk.NewDec(-1),
		},
		{
			name:       "invalid commission: over 1",
			commission: sdk.NewDec(2),
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
		broker := tc.broker
		if len(broker) == 0 {
			broker = defaultBroker
		}
		commission := tc.commission
		if commission.IsNil() {
			commission = defaultCommission
		}

		test.EvaluateTest(t, tc.name, func(t *testing.T) error {
			escrow := types.Escrow{
				Id:               id,
				Seller:           seller,
				Object:           test.MustPackToAny(obj),
				Price:            price,
				State:            state,
				Deadline:         deadline,
				BrokerAddress:    broker,
				BrokerCommission: commission,
			}
			err1 := escrow.Validate(denom, gen.NowAfter(0))
			err2 := escrow.ValidateWithoutDeadlineAndObject(denom)
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
		}.Validate(denom, gen.NowAfter(0))
	})
	test.EvaluateTest(t, "invalid escrow: just passed deadline", func(*testing.T) error {
		return types.Escrow{
			Id:       defaultId,
			Seller:   defaultSeller,
			Object:   test.MustPackToAny(defaultObj),
			Price:    defaultPrice,
			Deadline: gen.NowAfter(0),
		}.Validate(denom, gen.NowAfter(0))
	})
	test.EvaluateTest(t, "invalid escrow: passed deadline", func(*testing.T) error {
		return types.Escrow{
			Id:       defaultId,
			Seller:   defaultSeller,
			Object:   test.MustPackToAny(defaultObj),
			Price:    defaultPrice,
			Deadline: gen.NowAfter(0) - 20,
		}.Validate(denom, gen.NowAfter(0))
	})
	// TODO test valid escrow object deadline without context but invalid with context
}
