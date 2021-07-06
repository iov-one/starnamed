package keeper_test

import (
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type EscrowTestSuite struct {
	BaseKeeperSuite
	buyer             sdk.AccAddress
	seller            sdk.AccAddress
	expiredEscrowId   string
	refundedEscrowId  string
	completedEscrowId string
}

type assetState struct {
	sellerBalance sdk.Coins
	buyerBalance  sdk.Coins
	objectOwner   sdk.AccAddress
}

func (s *EscrowTestSuite) getState(seller sdk.AccAddress, buyer sdk.AccAddress, escrow types.Escrow, escrowFound bool) assetState {
	balance := func(address sdk.AccAddress) sdk.Coins {
		if address == nil {
			return sdk.NewCoins()
		}
		balance, found := s.balances[address.String()]
		if !found {
			return sdk.NewCoins()
		}
		return balance
	}
	sellerBalance := balance(seller)
	buyerBalance := balance(buyer)
	var owner sdk.AccAddress
	if !escrowFound {
		owner = nil
	} else {
		var obj types.TestObject
		// Refresh object value with store if possible
		if err := s.store.Read(escrow.GetObject().(*types.TestObject).PrimaryKey(), &obj); err != nil {
			owner = escrow.GetObject().(*types.TestObject).Owner
		} else {
			owner = obj.Owner
		}
	}

	return assetState{
		sellerBalance: sellerBalance,
		buyerBalance:  buyerBalance,
		objectOwner:   owner,
	}

}

func newSavedObject(generator *test.EscrowGenerator, seller sdk.AccAddress, store crud.Store) *types.TestObject {
	obj := generator.NewTestObject(seller)
	if err := store.Create(obj); err != nil {
		panic(err)
	}

	return obj
}

func (s *EscrowTestSuite) createErroredObjectEscrow(price sdk.Coins) string {
	// Create an object whose second transfer would fail
	obj := s.generator.NewErroredTestObject(1)
	if err := s.store.Create(obj); err != nil {
		panic(err)
	}

	id, err := s.keeper.CreateEscrow(
		s.ctx,
		s.seller,
		price,
		obj,
		s.generator.NowAfter(10),
	)
	if err != nil {
		panic(err)
	}
	return id
}

func (s *EscrowTestSuite) SetupTest() {
	s.generator = test.NewEscrowGenerator(uint64(s.ctx.BlockTime().Unix()))
	s.seller = s.generator.NewAccAddress()
	s.buyer = s.generator.NewAccAddress()
	s.keeper, s.ctx, s.store, s.balances, _ = test.NewTestKeeper([]sdk.AccAddress{s.buyer})

	s.msgServer = keeper.NewMsgServerImpl(s.keeper)

	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	escrow, _ := s.generator.NewTestEscrow(s.keeper.GetEscrowAddress(), price, s.generator.NowAfter(0)-10)
	escrow.State = types.EscrowState_Expired
	s.expiredEscrowId = escrow.Id
	s.keeper.SaveEscrow(s.ctx, escrow)

	escrow, _ = s.generator.NewTestEscrow(s.keeper.GetEscrowAddress(), price, s.generator.NowAfter(10))
	escrow.State = types.EscrowState_Refunded
	s.refundedEscrowId = escrow.Id
	s.keeper.SaveEscrow(s.ctx, escrow)

	escrow, _ = s.generator.NewTestEscrow(s.keeper.GetEscrowAddress(), price, s.generator.NowAfter(10))
	escrow.State = types.EscrowState_Completed
	s.completedEscrowId = escrow.Id
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.keeper.ImportNextID(s.ctx, s.generator.GetNextId())
}

func (s *EscrowTestSuite) TestCreate() {

	validAddress := s.seller
	defaultPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100)))
	defaultDeadline := s.generator.NowAfter(1)

	obj := s.generator.NewTestObject(s.seller)
	if err := s.store.Create(obj); err != nil {
		panic(err)
	}
	invalidObj := s.generator.NewTestObject(s.generator.NewAccAddress())
	if err := s.store.Create(invalidObj); err != nil {
		panic(err)
	}
	erroredObj := s.generator.NewErroredTestObject(0)
	if err := s.store.Create(erroredObj); err != nil {
		panic(err)
	}

	modifiedIdObj := *obj
	modifiedIdObj.Id += 10
	modifiedOwnerObj := *obj
	invalidAddr := s.generator.NewAccAddress()
	modifiedOwnerObj.Owner = append([]byte(nil), invalidAddr...)

	negativePrice := sdk.NewCoin(test.DenomAux, sdk.NewInt(5))
	negativePrice.Amount = negativePrice.Amount.SubRaw(10)

	testCases := []struct {
		name     string
		seller   sdk.AccAddress
		obj      types.TransferableObject
		price    sdk.Coins
		deadline uint64
	}{
		{
			name:     "valid scenario",
			seller:   validAddress,
			obj:      obj,
			price:    defaultPrice,
			deadline: defaultDeadline,
		},
		{
			name:     "invalid deadline: zero",
			seller:   validAddress,
			obj:      obj,
			price:    defaultPrice,
			deadline: 0,
		},
		{
			name:     "invalid deadline: in the past",
			seller:   validAddress,
			obj:      obj,
			price:    defaultPrice,
			deadline: s.generator.NowAfter(0) - 1,
		},
		{
			name:   "invalid price: negative",
			seller: validAddress,
			obj:    obj,
			price: sdk.Coins{
				sdk.NewCoin(test.Denom, sdk.NewInt(50)),
				negativePrice,
			},
			deadline: defaultDeadline,
		},
		{
			name:     "invalid object: does not belong to seller",
			seller:   validAddress,
			obj:      invalidObj,
			price:    defaultPrice,
			deadline: defaultDeadline,
		},
		{
			name:     "invalid object: not in sync with store",
			seller:   validAddress,
			obj:      &modifiedIdObj,
			price:    defaultPrice,
			deadline: defaultDeadline,
		},
		{
			name:     "invalid object: store version has different owner",
			seller:   invalidAddr,
			obj:      &modifiedOwnerObj,
			price:    defaultPrice,
			deadline: defaultDeadline,
		},
		{
			name:     "invalid transfer : object cannot be transferred",
			seller:   validAddress,
			obj:      erroredObj,
			price:    defaultPrice,
			deadline: defaultDeadline,
		},
	}

	for _, t := range testCases {
		create := func(*testing.T) error {
			_, err := s.keeper.CreateEscrow(
				s.ctx,
				t.seller,
				t.price,
				t.obj,
				t.deadline,
			)
			return err
		}

		test.EvaluateTest(s.T(), t.name, create)
	}
}

func (s *EscrowTestSuite) TestUpdate() {
	newSeller := s.generator.NewAccAddress()
	escrowDeadline := s.generator.NowAfter(10)
	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	id, err := s.keeper.CreateEscrow(
		s.ctx,
		s.seller,
		price,
		newSavedObject(s.generator, s.seller, s.store),
		escrowDeadline,
	)
	if err != nil {
		panic(err)
	}
	escrowID := id
	testCases := []struct {
		name     string
		id       string
		seller   sdk.AccAddress
		updater  sdk.AccAddress
		price    sdk.Coins
		deadline uint64
	}{
		{
			name:    "price update",
			updater: s.seller,
			price:   sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(10))),
		},
		{
			name:     "deadline update",
			updater:  s.seller,
			deadline: s.generator.NowAfter(100000),
		},
		{
			name:     "multiple fields update",
			updater:  s.seller,
			price:    sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100))),
			deadline: s.generator.NowAfter(1000000),
		},
		{
			name:     "invalid deadline: earlier than before",
			updater:  s.seller,
			deadline: escrowDeadline - 1,
		},
		{
			name:     "invalid deadline: past deadline",
			updater:  s.seller,
			deadline: 1,
		},
		{
			name:    "invalid update: empty",
			updater: s.seller,
		},
		{
			name:    "invalid updater: not the escrow seller",
			updater: newSeller,
			seller:  newSeller,
		},
		{
			name:    "invalid updater: not the escrow seller (2)",
			updater: newSeller,
			seller:  s.seller,
		},
		{
			name:    "invalid escrow: non existing",
			updater: s.seller,
			seller:  newSeller,
			id:      "AAAAAAAAAAAAAAFA",
		},
		{
			name:    "invalid escrow: expired",
			updater: s.seller,
			id:      s.expiredEscrowId,
		},
		{
			name:    "invalid escrow: completed",
			updater: s.seller,
			id:      s.completedEscrowId,
		},
		{
			name:    "invalid escrow: refunded",
			updater: s.seller,
			id:      s.refundedEscrowId,
		},
		// Put this test at the end so it does not mess with the other tests
		{
			name:    "seller update",
			updater: s.seller,
			seller:  s.generator.NewAccAddress(),
		},
	}

	for _, t := range testCases {
		id := t.id
		if len(id) == 0 {
			id = escrowID
		}

		update := func(*testing.T) error {
			return s.keeper.UpdateEscrow(s.ctx, id, t.updater, t.seller, t.price, t.deadline)
		}

		test.EvaluateTest(s.T(), t.name, update)
	}
}

func (s *EscrowTestSuite) TestTransferTo() {
	var testEscrows = make(map[string]string)
	defaultPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	prices := map[string]sdk.Coins{
		"default":   defaultPrice,
		"default1":  defaultPrice,
		"default2":  defaultPrice,
		"default3":  defaultPrice,
		"expensive": sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100000000000))),
		"multi":     sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)), sdk.NewCoin(test.DenomAux, sdk.NewInt(50))),
	}
	for name, price := range prices {
		id, err := s.keeper.CreateEscrow(
			s.ctx,
			s.seller,
			price,
			newSavedObject(s.generator, s.seller, s.store),
			s.generator.NowAfter(10),
		)
		if err != nil {
			panic(err)
		}
		testEscrows[name] = id
	}

	invalidObjectEscrowId := s.createErroredObjectEscrow(defaultPrice)

	checkDefaultValidTransfer := func(before, after assetState, name string) {
		s.Assert().Equal(before.buyerBalance, after.buyerBalance.Add(defaultPrice...), "Buyer balance on test %s", name)
		s.Assert().Equal(before.sellerBalance.Add(defaultPrice...), after.sellerBalance, "Seller balance on test %s", name)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(), "Object owner on test %s", name)
		s.Assert().Equal(after.objectOwner, s.buyer, "Object owner on test %s", name)
	}

	checkInvalidTransferWithoutObject := func(before, after assetState, name string) {
		s.Assert().Equal(before.buyerBalance, after.buyerBalance, "Buyer balance on test %s", name)
		s.Assert().Equal(before.sellerBalance, after.sellerBalance, "Seller balance on test %s", name)

	}

	checkInvalidTransfer := func(before, after assetState, name string) {
		checkInvalidTransferWithoutObject(before, after, name)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(), "Object owner on test %s", name)
		s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress(), "Object owner on test %s", name)
	}

	testCases := []struct {
		name   string
		buyer  sdk.AccAddress
		amount sdk.Coins
		id     string
		check  func(before, after assetState, name string)
	}{
		{
			name:   "valid transfer: exact coins",
			buyer:  s.buyer,
			amount: defaultPrice,
			check:  checkDefaultValidTransfer,
			id:     testEscrows["default1"],
		},
		{
			name:   "valid transfer: too much coin",
			buyer:  s.buyer,
			amount: defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(20))),
			check:  checkDefaultValidTransfer,
			id:     testEscrows["default2"],
		},
		{
			name:   "valid transfer: too much coins",
			buyer:  s.buyer,
			amount: defaultPrice.Add(sdk.NewCoin(test.DenomAux, sdk.NewInt(30))),
			check:  checkDefaultValidTransfer,
			id:     testEscrows["default3"],
		},
		{
			name:   "invalid transfer: not enough coin",
			buyer:  s.buyer,
			amount: defaultPrice.Sub(sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(1)))),
			check:  checkInvalidTransfer,
		},
		{
			name:   "invalid transfer: not enough coin on one coin",
			buyer:  s.buyer,
			amount: prices["multi"].Sub(sdk.NewCoins(sdk.NewCoin(test.DenomAux, sdk.NewInt(10)))),
			id:     testEscrows["multi"],
			check:  checkInvalidTransfer,
		},
		{
			name:  "invalid buyer: it is the seller",
			buyer: s.seller,
			check: checkInvalidTransfer,
		},
		{
			name:   "invalid escrow: non existing escrow",
			buyer:  s.buyer,
			amount: defaultPrice,
			id:     "AABBCCDDEEFF1122",
			check:  checkInvalidTransferWithoutObject,
		},
		{
			name:   "invalid escrow: expired escrow",
			buyer:  s.buyer,
			amount: defaultPrice,
			id:     s.expiredEscrowId,
			check:  checkInvalidTransfer,
		},
		{
			name:   "invalid escrow: refunded escrow",
			buyer:  s.buyer,
			amount: defaultPrice,
			id:     s.refundedEscrowId,
			check:  checkInvalidTransfer,
		},
		{
			name:   "invalid escrow: completed escrow",
			buyer:  s.buyer,
			amount: defaultPrice,
			id:     s.completedEscrowId,
			check:  checkInvalidTransfer,
		},
		{
			name:   "invalid transfer: not enough coins on buyer account",
			buyer:  s.buyer,
			amount: prices["expensive"],
			id:     testEscrows["expensive"],
			check:  checkInvalidTransfer,
		},
		{
			name:   "panic when error on object transfer",
			buyer:  s.buyer,
			amount: defaultPrice,
			id:     invalidObjectEscrowId,
			check:  checkInvalidTransfer,
		},
	}

	for _, t := range testCases {

		id := t.id
		if len(id) == 0 {
			id = testEscrows["default"]
		}

		check := t.check
		if check == nil {
			check = func(before, after assetState, name string) {}
		}

		transfer := func(*testing.T) error {
			escrow, found := s.keeper.GetEscrow(s.ctx, id)
			before := s.getState(s.seller, t.buyer, escrow, found)
			err := s.keeper.TransferToEscrow(s.ctx, t.buyer, id, t.amount)
			check(before, s.getState(s.seller, t.buyer, escrow, found), t.name)
			return err
		}

		test.EvaluateTest(s.T(), t.name, transfer)
	}

}

func (s *EscrowTestSuite) TestRefund() {
	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	lastBlockTime := s.keeper.GetLastBlockTime(s.ctx)

	states := []types.EscrowState{types.EscrowState_Open, types.EscrowState_Expired}
	var notExpiredEscrowIds, expiredEscrowIds [2]string
	escrowIds := []*[2]string{&notExpiredEscrowIds, &expiredEscrowIds}
	for i, state := range states {
		for j := 0; j < 2; j++ {
			deadline := lastBlockTime + 5
			if state == types.EscrowState_Expired {
				deadline -= 50
			}
			// Cheat a little to be able to create the escrows in the past
			s.keeper.SetLastBlockTime(s.ctx, deadline-1)
			id, err := s.keeper.CreateEscrow(
				s.ctx,
				s.seller,
				price,
				newSavedObject(s.generator, s.seller, s.store),
				deadline,
			)

			s.keeper.MarkExpiredEscrows(s.ctx, lastBlockTime)

			if err != nil {
				panic(err)
			}
			escrowIds[i][j] = id
		}
	}

	s.keeper.SetLastBlockTime(s.ctx, lastBlockTime)

	invalidObjectEscrowId := s.createErroredObjectEscrow(price)

	validRefund := func(before, after assetState) {
		s.Assert().Equal(before.sellerBalance, after.sellerBalance)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress())
		s.Assert().Equal(after.objectOwner, s.seller)
	}

	invalidRefund := func(before, after assetState) {
		s.Assert().Equal(before.sellerBalance, after.sellerBalance)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress())
		s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress())
	}

	testCases := []struct {
		name   string
		id     string
		sender sdk.AccAddress
		check  func(before, after assetState)
	}{
		{
			name:   "valid refund: triggered by seller before expiration",
			sender: s.seller,
			id:     notExpiredEscrowIds[0],
			check:  validRefund,
		},
		{
			name:   "valid refund: triggered by seller after expiration",
			sender: s.seller,
			id:     expiredEscrowIds[0],
			check:  validRefund,
		},
		{
			name:   "valid refund: triggered by random address after expiration",
			sender: s.generator.NewAccAddress(),
			id:     expiredEscrowIds[1],
			check:  validRefund,
		},
		{
			name:   "invalid refund: triggered by random address before expiration",
			sender: s.generator.NewAccAddress(),
			id:     notExpiredEscrowIds[1],
			check:  invalidRefund,
		},
		{
			name:   "invalid state: escrow already completed",
			sender: s.seller,
			id:     s.completedEscrowId,
			check:  invalidRefund,
		},
		{
			name:   "invalid state: escrow already refunded",
			sender: s.seller,
			id:     s.refundedEscrowId,
			check:  invalidRefund,
		},
		{
			name:   "invalid escrow: non existing",
			sender: s.seller,
			id:     "AABBCCDDEEFF1122",
			check:  invalidRefund,
		},
		{
			name:   "invalid refund: error on object transfer",
			sender: s.seller,
			id:     invalidObjectEscrowId,
		},
	}

	for _, t := range testCases {
		check := func(before, after assetState) {}

		refund := func(*testing.T) error {
			escrow, found := s.keeper.GetEscrow(s.ctx, t.id)
			before := s.getState(s.seller, nil, escrow, found)
			err := s.keeper.RefundEscrow(s.ctx, t.sender, t.id)
			check(before, s.getState(s.seller, nil, escrow, found))
			return err
		}

		test.EvaluateTest(s.T(), t.name, refund)
	}
}

func generateExpiringEscrows(generator *test.EscrowGenerator) []types.Escrow {
	deriveEscrow := func(e types.Escrow, state types.EscrowState) types.Escrow {
		e.State = state
		e.Id = generator.NextID()
		return e
	}

	expiredEscrow, _ := generator.NewTestEscrow(
		generator.NewAccAddress(),
		sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(10))),
		generator.NowAfter(0)-5,
	)
	nonExpiredEscrow := expiredEscrow
	nonExpiredEscrow.Deadline = generator.NowAfter(10)

	return []types.Escrow{
		deriveEscrow(expiredEscrow, types.EscrowState_Expired),
		deriveEscrow(expiredEscrow, types.EscrowState_Open),
		deriveEscrow(nonExpiredEscrow, types.EscrowState_Expired),
		deriveEscrow(nonExpiredEscrow, types.EscrowState_Open),
	}
}

func (s *EscrowTestSuite) TestMarkExpiredEscrows() {
	escrows := generateExpiringEscrows(s.generator)
	for _, e := range escrows {
		s.keeper.SaveEscrow(s.ctx, e)
	}

	s.keeper.MarkExpiredEscrows(s.ctx, s.generator.NowAfter(0))

	// Refresh escrow array
	for i, e := range escrows {
		escrows[i], _ = s.keeper.GetEscrow(s.ctx, e.Id)
	}

	s.Assert().Equal(types.EscrowState_Expired, escrows[0].State, "The first escrow should not have been modified")
	s.Assert().Equal(types.EscrowState_Expired, escrows[1].State, "The second escrow should have been marked as expired")
	s.Assert().Equal(types.EscrowState_Expired, escrows[2].State, "The third escrow should not have been modified")
	s.Assert().Equal(types.EscrowState_Open, escrows[3].State, "The fourth escrow should not have been modified")
}

func (s *EscrowTestSuite) TestIterateEscrowsWithPassedDeadline() {
	escrows := generateExpiringEscrows(s.generator)
	for _, e := range escrows {
		s.keeper.SaveEscrow(s.ctx, e)
	}

	s.keeper.IterateEscrowsWithPassedDeadline(s.ctx, s.generator.NowAfter(0), func(escrow types.Escrow) bool {
		s.Assert().LessOrEqual(escrow.Deadline, s.generator.NowAfter(0), "This escrow has not a passed deadline")
		return true
	})
}

func TestEscrow(t *testing.T) {
	suite.Run(t, new(EscrowTestSuite))
}
