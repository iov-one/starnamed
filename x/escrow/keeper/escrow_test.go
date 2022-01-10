package keeper_test

import (
	"encoding/hex"
	"testing"

	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	crud "github.com/iov-one/cosmos-sdk-crud"
	"github.com/stretchr/testify/suite"

	"github.com/iov-one/starnamed/x/escrow/keeper"
	"github.com/iov-one/starnamed/x/escrow/test"
	"github.com/iov-one/starnamed/x/escrow/types"
)

type EscrowTestSuite struct {
	BaseKeeperSuite
	buyer                      sdk.AccAddress
	seller                     sdk.AccAddress
	expiredEscrowId            string
	refundedEscrowId           string
	completedEscrowId          string
	timeConstrainedObjectsData *testWithCustomTimeConstrainData
}

type assetState struct {
	sellerBalance sdk.Coins
	buyerBalance  sdk.Coins
	brokerBalance sdk.Coins
	objectOwner   sdk.AccAddress
}

func (s *EscrowTestSuite) getState(seller sdk.AccAddress, buyer sdk.AccAddress, broker sdk.AccAddress, escrow types.Escrow, escrowFound bool) assetState {
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
	brokerBalance := balance(broker)
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
		brokerBalance: brokerBalance,
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

func (s *EscrowTestSuite) createErroredObjectEscrow(price sdk.Coins, isAuction bool) string {
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
		isAuction,
	)
	if err != nil {
		panic(err)
	}
	return id
}

// Used for context-aware object validation testing
type testWithCustomTimeConstrainData struct {
	storeKey  sdk.StoreKey
	crudStore crud.Store
}

func (d testWithCustomTimeConstrainData) getDeadlineStore(ctx sdk.Context) sdk.KVStore {
	return prefix.NewStore(ctx.KVStore(d.storeKey), []byte{0x42})
}

func (d testWithCustomTimeConstrainData) GetCrudStore() crud.Store {
	return d.crudStore
}

func (d testWithCustomTimeConstrainData) GetDeadlineOrDefault(ctx sdk.Context, obj types.TransferableObject, defaultDeadline uint64) uint64 {
	bytes := d.getDeadlineStore(ctx).Get(obj.GetUniqueKey())
	if bytes == nil {
		return defaultDeadline
	}
	return sdk.BigEndianToUint64(bytes)
}

func (d testWithCustomTimeConstrainData) SetDeadline(ctx sdk.Context, obj types.TransferableObject, deadline uint64) {
	d.getDeadlineStore(ctx).Set(obj.GetUniqueKey(), sdk.Uint64ToBigEndian(deadline))
}

func (s *EscrowTestSuite) SetupTest() {
	test.SetConfig()
	s.generator = test.NewEscrowGenerator(uint64(test.TimeNow.Unix()))
	s.seller = s.generator.NewAccAddress()
	s.buyer = s.generator.NewAccAddress()
	var storeKey sdk.StoreKey
	s.keeper, s.ctx, s.store, s.balances, storeKey, s.configKeeper = test.NewTestKeeper([]sdk.AccAddress{s.buyer}, true)

	s.msgServer = keeper.NewMsgServerImpl(s.keeper)

	getNextID := func(gen *test.EscrowGenerator) string {
		return hex.EncodeToString(sdk.Uint64ToBigEndian(gen.GetNextId()))
	}

	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	escrow, _ := s.generator.NewTestEscrow(s.keeper.GetEscrowAddress(getNextID(s.generator)), price, s.generator.NowAfter(0)-10)
	escrow.State = types.EscrowState_Expired
	s.expiredEscrowId = escrow.Id
	s.keeper.SaveEscrow(s.ctx, escrow)

	escrow, _ = s.generator.NewTestEscrow(s.keeper.GetEscrowAddress(getNextID(s.generator)), price, s.generator.NowAfter(10))
	escrow.State = types.EscrowState_Refunded
	s.refundedEscrowId = escrow.Id
	s.keeper.SaveEscrow(s.ctx, escrow)

	escrow, _ = s.generator.NewTestEscrow(s.keeper.GetEscrowAddress(getNextID(s.generator)), price, s.generator.NowAfter(10))
	escrow.State = types.EscrowState_Completed
	s.completedEscrowId = escrow.Id
	s.keeper.SaveEscrow(s.ctx, escrow)

	s.timeConstrainedObjectsData = &testWithCustomTimeConstrainData{
		storeKey:  storeKey,
		crudStore: s.store,
	}

	s.keeper.RegisterCustomData(types.TypeIDTestTimeConstrainedObject, s.timeConstrainedObjectsData)

	s.keeper.ImportNextID(s.ctx, s.generator.GetNextId())
}

func (s *EscrowTestSuite) TestCreate() {

	validAddress := s.seller
	defaultPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100)))
	defaultDeadline := s.generator.NowAfter(1)

	createAndSaveObject := func() *types.TestObject {
		return newSavedObject(s.generator, s.seller, s.store)
	}
	createAndSaveTimeConstrainedObject := func() *types.TestTimeConstrainedObject {
		obj := s.generator.NewTimeConstrainedObject(s.seller, s.generator.NowAfter(10))
		if err := s.store.Create(obj); err != nil {
			panic(err)
		}
		return obj
	}
	createAndSaveTimeConstrainedObjectWithContextDeadline := func(deadline uint64) *types.TestTimeConstrainedObject {
		obj := createAndSaveTimeConstrainedObject()
		s.timeConstrainedObjectsData.SetDeadline(s.ctx, obj, deadline)
		return obj
	}
	invalidObj := s.generator.NewTestObject(s.generator.NewAccAddress())
	if err := s.store.Create(invalidObj); err != nil {
		panic(err)
	}
	erroredObj := s.generator.NewErroredTestObject(0)
	if err := s.store.Create(erroredObj); err != nil {
		panic(err)
	}

	obj := createAndSaveObject()
	modifiedIdObj := *obj
	modifiedIdObj.Id += 10
	modifiedOwnerObj := *obj
	invalidAddr := s.generator.NewAccAddress()
	modifiedOwnerObj.Owner = append([]byte(nil), invalidAddr...)
	maxDuration := s.keeper.GetMaximumEscrowDuration(s.ctx)

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
			obj:      createAndSaveObject(),
			price:    defaultPrice,
			deadline: defaultDeadline,
		},
		{
			name:     "valid scenario with max duration",
			seller:   validAddress,
			obj:      createAndSaveObject(),
			price:    defaultPrice,
			deadline: s.generator.NowAfter(uint64(maxDuration.Seconds())),
		},
		{
			name:     "valid scenario with time constrained object",
			seller:   validAddress,
			obj:      createAndSaveTimeConstrainedObject(),
			price:    defaultPrice,
			deadline: s.generator.NowAfter(5),
		},
		{
			name:     "valid scenario with context-aware time constrained object",
			seller:   validAddress,
			obj:      createAndSaveTimeConstrainedObjectWithContextDeadline(s.generator.NowAfter(10)),
			price:    defaultPrice,
			deadline: s.generator.NowAfter(5),
		},
		{
			name:     "invalid price: zero",
			seller:   validAddress,
			obj:      createAndSaveObject(),
			price:    sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.ZeroInt())),
			deadline: s.generator.NowAfter(uint64(maxDuration.Seconds())),
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
			name:     "invalid deadline: maximum duration exceeded",
			seller:   validAddress,
			obj:      obj,
			price:    defaultPrice,
			deadline: s.generator.NowAfter(uint64(maxDuration.Seconds()) + 1),
		},
		{
			name:     "invalid deadline: not validated by object",
			seller:   validAddress,
			obj:      createAndSaveTimeConstrainedObject(),
			price:    defaultPrice,
			deadline: s.generator.NowAfter(15),
		},
		{
			name:     "invalid deadline: not validated by object with context-aware validation",
			seller:   validAddress,
			obj:      createAndSaveTimeConstrainedObjectWithContextDeadline(s.generator.NowAfter(2)),
			price:    defaultPrice,
			deadline: s.generator.NowAfter(5),
		},
		{
			name:     "invalid deadline: not validated by object with both validations",
			seller:   validAddress,
			obj:      createAndSaveTimeConstrainedObjectWithContextDeadline(s.generator.NowAfter(10)),
			price:    defaultPrice,
			deadline: s.generator.NowAfter(15),
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
			name:     "invalid price: empty",
			seller:   validAddress,
			obj:      obj,
			price:    sdk.Coins{},
			deadline: defaultDeadline,
		},
		{
			name:     "invalid price: not the correct token denomination",
			seller:   validAddress,
			obj:      obj,
			price:    sdk.NewCoins(sdk.NewCoin("abcd", sdk.OneInt())),
			deadline: defaultDeadline,
		},
		{
			name:   "invalid price: not the correct token denomination with multiple coin types",
			seller: validAddress,
			obj:    obj,
			price: sdk.NewCoins(
				sdk.NewCoin(test.Denom, sdk.OneInt()),
				sdk.NewCoin("abcd", sdk.OneInt()),
			),
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
				false,
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
	escrowID, err := s.keeper.CreateEscrow(
		s.ctx,
		s.seller,
		price,
		newSavedObject(s.generator, s.seller, s.store),
		escrowDeadline,
		false,
	)
	if err != nil {
		panic(err)
	}

	auctionId, err := s.keeper.CreateEscrow(
		s.ctx,
		s.seller,
		price,
		newSavedObject(s.generator, s.seller, s.store),
		escrowDeadline,
		true,
		)
	if err != nil {
		panic(err)
	}

	timeConstrainedObj := s.generator.NewTimeConstrainedObject(s.seller, s.generator.NowAfter(10))
	if err := s.store.Create(timeConstrainedObj); err != nil {
		panic(err)
	}
	escrowWithTimeConstrainedObjectID, err := s.keeper.CreateEscrow(
		s.ctx, s.seller, price, timeConstrainedObj, s.generator.NowAfter(5), false,
	)

	contextAwaretimeConstrainObj := s.generator.NewTimeConstrainedObject(s.seller, s.generator.NowAfter(10))
	s.timeConstrainedObjectsData.SetDeadline(s.ctx, contextAwaretimeConstrainObj, s.generator.NowAfter(8))
	if err := s.store.Create(contextAwaretimeConstrainObj); err != nil {
		panic(err)
	}
	escrowWithContextAwareTimeConstrainedObjectID, err := s.keeper.CreateEscrow(
		s.ctx, s.seller, price, contextAwaretimeConstrainObj, s.generator.NowAfter(5), false,
	)

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
			name:     "deadline update with time constrained object",
			updater:  s.seller,
			deadline: s.generator.NowAfter(8),
			id:       escrowWithTimeConstrainedObjectID,
		},
		{
			name:     "deadline update with context aware time constrained object",
			updater:  s.seller,
			deadline: s.generator.NowAfter(7),
			id:       escrowWithContextAwareTimeConstrainedObjectID,
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
			name:     "invalid deadline: more than maximum duration",
			updater:  s.seller,
			deadline: s.generator.NowAfter(10 + uint64(s.keeper.GetMaximumEscrowDuration(s.ctx).Seconds())),
		},
		{
			name:     "invalid deadline: not validated by the object",
			updater:  s.seller,
			deadline: s.generator.NowAfter(100),
			id:       escrowWithTimeConstrainedObjectID,
		},
		{
			name:     "invalid deadline: not validated by an object with a context time constraint",
			updater:  s.seller,
			deadline: s.generator.NowAfter(9),
			id:       escrowWithContextAwareTimeConstrainedObjectID,
		},
		{
			name:    "invalid update: empty",
			updater: s.seller,
		},
		{
			name:    "invalid update: updating the price of an auction",
			updater: s.seller,
			id: auctionId,
			price: price.Add(sdk.NewCoin(test.Denom, sdk.OneInt())),
		},
		{
			name:    "invalid update: updating the deadline of an auction",
			updater: s.seller,
			id: auctionId,
			deadline: s.generator.NowAfter(15),
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
		{
			name:    "invalid price: negative",
			updater: s.seller,
			price: sdk.Coins{
				sdk.Coin{Denom: test.Denom, Amount: sdk.NewInt(-10)},
			},
		},
		{
			name:    "invalid price: empty",
			updater: s.seller,
			price:   sdk.Coins{},
		},
		{
			name:    "invalid price: not the correct token denomination",
			updater: s.seller,
			price:   sdk.NewCoins(sdk.NewCoin("abcd", sdk.OneInt())),
		},
		{
			name:    "invalid price: not the correct token denomination with multiple coin types",
			updater: s.seller,
			price: sdk.NewCoins(
				sdk.NewCoin(test.Denom, sdk.OneInt()),
				sdk.NewCoin("abcd", sdk.OneInt()),
			),
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
	defaultConfig := s.configKeeper.GetConfiguration(s.ctx)

	var testEscrows = make(map[string]string)
	defaultPriceAmt := int64(50)
	defaultPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(defaultPriceAmt)))
	brokerAddr := s.generator.NewAccAddress()
	prices := map[string]sdk.Coins{
		"default":   defaultPrice,
		"expensive": sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100000000000))),
	}

	createAndSaveEscrow := func(seller sdk.AccAddress, price sdk.Coins) string {
		id, err := s.keeper.CreateEscrow(
			s.ctx,
			seller,
			price,
			newSavedObject(s.generator, seller, s.store),
			s.generator.NowAfter(10),
			false,
		)
		if err != nil {
			panic(err)
		}
		return id
	}

	createAndSaveEscrowWithBroker := func(seller sdk.AccAddress, price sdk.Coins, commission sdk.Dec) string {
		config := defaultConfig
		config.EscrowBroker = brokerAddr.String()
		config.EscrowCommission = commission
		s.configKeeper.SetConfig(s.ctx, config)
		id := createAndSaveEscrow(seller, price)
		s.configKeeper.SetConfig(s.ctx, defaultConfig)
		return id
	}

	for name, price := range prices {
		testEscrows[name] = createAndSaveEscrow(s.seller, price)
	}

	invalidObjectEscrowId := s.createErroredObjectEscrow(defaultPrice, false)

	checkDefaultValidTransfer := func(before, after assetState, name string, expectedCommission sdk.Coins, id string) {
		price := defaultPrice.Sub(expectedCommission)
		expectedBrokerBalance := sdk.NewCoins(before.brokerBalance.Add(expectedCommission...)...)

		s.Assert().Equal(before.buyerBalance, after.buyerBalance.Add(defaultPrice...), "Buyer balance on test %s", name)
		s.Assert().Equal(before.sellerBalance.Add(price...), after.sellerBalance, "Seller balance on test %s", name)
		s.Assert().Equal(expectedBrokerBalance, after.brokerBalance, "Broker balance on test %s", name)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
		s.Assert().Equal(after.objectOwner, s.buyer, "Object owner on test %s", name)
	}

	checkInvalidTransferWithoutObject := func(before, after assetState, name string, _ sdk.Coins, id string) {
		s.Assert().Equal(before.buyerBalance, after.buyerBalance, "Buyer balance on test %s", name)
		s.Assert().Equal(before.sellerBalance, after.sellerBalance, "Seller balance on test %s", name)
		s.Assert().Equal(before.brokerBalance, after.brokerBalance, "Broker balance on test %s", name)
	}

	checkInvalidTransfer := func(before, after assetState, name string, _ sdk.Coins, id string) {
		checkInvalidTransferWithoutObject(before, after, name, nil, id)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
		s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
	}

	testCases := []struct {
		name               string
		buyer              sdk.AccAddress
		amount             sdk.Coins
		id                 string
		expectedCommission sdk.Coins
		broker             sdk.AccAddress
		check              func(before, after assetState, name string, expectedCommision sdk.Coins, id string)
	}{
		{
			name:   "valid transfer: exact coins",
			buyer:  s.buyer,
			amount: defaultPrice,
			check:  checkDefaultValidTransfer,
			id:     createAndSaveEscrow(s.seller, defaultPrice),
		},
		{
			name:   "valid transfer: too much coin",
			buyer:  s.buyer,
			amount: defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(20))),
			check:  checkDefaultValidTransfer,
			id:     createAndSaveEscrow(s.seller, defaultPrice),
		},
		{
			name:   "invalid transfer: multiple coins",
			buyer:  s.buyer,
			amount: defaultPrice.Add(sdk.NewCoin(test.DenomAux, sdk.NewInt(30))),
			check:  checkInvalidTransfer,
			id:     createAndSaveEscrow(s.seller, defaultPrice),
		},
		{
			name:               "valid transfer with broker",
			buyer:              s.buyer,
			amount:             defaultPrice,
			check:              checkDefaultValidTransfer,
			broker:             brokerAddr,
			id:                 createAndSaveEscrowWithBroker(s.seller, defaultPrice, sdk.NewDec(1).Quo(sdk.NewDec(10))),
			expectedCommission: sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(defaultPriceAmt/10))),
		},
		{
			name:               "invalid transfer with broker and a composite price",
			buyer:              s.buyer,
			amount:             defaultPrice.Add(sdk.NewCoin(test.DenomAux, sdk.NewInt(30))),
			check:              checkInvalidTransfer,
			broker:             brokerAddr,
			id:                 createAndSaveEscrowWithBroker(s.seller, defaultPrice, sdk.NewDec(1).Quo(sdk.NewDec(5))),
			expectedCommission: sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(defaultPriceAmt/5))),
		},
		{
			name:               "valid transfer with broker and decimal escrow commission",
			buyer:              s.buyer,
			amount:             defaultPrice,
			check:              checkDefaultValidTransfer,
			broker:             brokerAddr,
			id:                 createAndSaveEscrowWithBroker(s.seller, defaultPrice, sdk.NewDec(1).Quo(sdk.NewDec(3))),
			expectedCommission: sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(defaultPriceAmt/3))),
		},
		{
			name:   "invalid transfer: not enough coin",
			buyer:  s.buyer,
			amount: defaultPrice.Sub(sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(1)))),
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
			check = func(_, _ assetState, _ string, _ sdk.Coins, _ string) {}
		}

		transfer := func(*testing.T) error {
			escrow, found := s.keeper.GetEscrow(s.ctx, id)
			before := s.getState(s.seller, t.buyer, t.broker, escrow, found)
			err := s.keeper.TransferToEscrow(s.ctx, t.buyer, id, t.amount)
			check(before, s.getState(s.seller, t.buyer, t.broker, escrow, found), t.name, t.expectedCommission, id)
			return err
		}

		test.EvaluateTest(s.T(), t.name, transfer)
	}
}

func (s *EscrowTestSuite) TestBid() {
	defaultPriceAmt := int64(50)
	defaultPrice := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(defaultPriceAmt)))
	prices := map[string]sdk.Coins{
		"default":   defaultPrice,
		"expensive": sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(100000000000))),
	}

	createAndSaveAuction := func(seller sdk.AccAddress, price sdk.Coins) string {
		id, err := s.keeper.CreateEscrow(
			s.ctx,
			seller,
			price,
			newSavedObject(s.generator, seller, s.store),
			s.generator.NowAfter(10),
			true,
		)
		if err != nil {
			panic(err)
		}
		return id
	}
	testAuction := createAndSaveAuction(s.seller, prices["default"])


	checkDefaultValidBid:= func(before, after assetState, name string, price sdk.Coins, id string, lastBid sdk.Coins) {
		s.Assert().Equal(before.buyerBalance.Sub(price), after.buyerBalance, "Buyer balance on test %s", name)
		s.Assert().Equal(before.sellerBalance, after.sellerBalance, "Seller balance on test %s", name)
		//FIXME: we use the broker field as the last bidder, and that's not pretty
		// We add .String() to have the same format on both side (nil coins vs empty non-nil coins)
		s.Assert().Equal(before.brokerBalance.Add(lastBid...).String(), after.brokerBalance.String(), "Previous bidder balance on test %s", name)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
		s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
	}
	checkInvalidBid := func(before, after assetState, name string, _ sdk.Coins, id string, _ sdk.Coins) {
		s.Assert().Equal(before.buyerBalance, after.buyerBalance, "Buyer balance on test %s", name)
		s.Assert().Equal(before.sellerBalance, after.sellerBalance, "Seller balance on test %s", name)
		s.Assert().Equal(before.brokerBalance, after.brokerBalance, "Broker balance on test %s", name)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
		s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress(id), "Object owner on test %s", name)
	}

	testAuctionToBid := createAndSaveAuction(s.seller, defaultPrice)
	buyers := []sdk.AccAddress{s.generator.NewAccAddress(), s.generator.NewAccAddress()}

	for i := range buyers {
		s.balances[buyers[i].String()] = sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(10000)))
	}

	testCases := []struct {
		name               string
		bidders            []sdk.AccAddress
		amounts            []sdk.Coins
		id                 string
		check              func(before, after assetState, name string, price sdk.Coins, id string, lastBid sdk.Coins)
	}{
		{
			name:   "valid transfer: 1 bid",
			bidders: []sdk.AccAddress{s.buyer},
			amounts: []sdk.Coins{defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(20)))},
			check:  checkDefaultValidBid,
			id:     testAuctionToBid,
		},
		{
			name:   "invalid bid: bid too low on 2nd bid",
			bidders:  []sdk.AccAddress{s.generator.NewAccAddress()},
			amounts: []sdk.Coins{defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(10)))},
			id: testAuctionToBid,
			check:  checkInvalidBid,
		},
		{
			name:   "valid bid: multiple bids",
			bidders: []sdk.AccAddress{buyers[0], buyers[1], s.buyer},
			amounts: []sdk.Coins{
				defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(10))),
				defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(20))),
				defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(30)))},
			check:  checkDefaultValidBid,
			id:     createAndSaveAuction(s.seller, defaultPrice),
		},
		{
			name:   "invalid bid: multiple coins",
			bidders:  []sdk.AccAddress{s.buyer},
			amounts: []sdk.Coins{defaultPrice.Add(sdk.NewCoin(test.DenomAux, sdk.NewInt(30)))},
			check:  checkInvalidBid,
		},
		{
			name:   "invalid bid: 1 bid exact coins",
			bidders:  []sdk.AccAddress{s.buyer},
			amounts: []sdk.Coins{defaultPrice},
			check:  checkInvalidBid,
		},
		{
			name:   "invalid bid: bid too low",
			bidders:  []sdk.AccAddress{s.buyer},
			amounts: []sdk.Coins{defaultPrice.Sub(sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(1))))},
			check:  checkInvalidBid,
		},
		{
			name:   "invalid bid: not enough coins on buyer account",
			bidders:  []sdk.AccAddress{s.buyer},
			amounts: []sdk.Coins{ prices["expensive"]},
			check:  checkInvalidBid,
		},
		{
			name:  "invalid buyer: it is the seller",
			bidders: []sdk.AccAddress{s.seller},
			amounts: []sdk.Coins{ defaultPrice},
			check: checkInvalidBid,
		},
		//TODO: should we disallow adding a bid when the bidder is already the last bidder ?
		/*
		{
			name:  "invalid buyer: it is the last bidder",
			bidders: []sdk.AccAddress {s.buyer, s.buyer},
			amounts: []sdk.Coins{ defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(10))), defaultPrice.Add(sdk.NewCoin(test.Denom, sdk.NewInt(20)))},
			check: checkInvalidBid,
		},*/
	}

	// Set a commission to ensure bidding is not affected by commissions
	defaultConfig := s.configKeeper.GetConfiguration(s.ctx)
	config := defaultConfig
	config.EscrowCommission = sdk.NewDec(1).Quo(sdk.NewDec(5))
	s.configKeeper.SetConfig(s.ctx, config)

	for _, t := range testCases {

		id := t.id
		if len(id) == 0 {
			id = testAuction
		}

		check := t.check
		if check == nil {
			check = func(assetState, assetState, string, sdk.Coins, string, sdk.Coins) {}
		}

		transfer := func(*testing.T) error {
			auction, found := s.keeper.GetEscrow(s.ctx, id)
			var err error
			var lastBidder sdk.AccAddress = nil
			lastBid := sdk.Coins(nil)
			for i, buyer := range t.bidders {
				before := s.getState(s.seller, buyer, lastBidder, auction, found)
				if err = s.keeper.TransferToEscrow(s.ctx, buyer, id, t.amounts[i]); err != nil {
					break
				}
				check(before, s.getState(s.seller, buyer, lastBidder, auction, found), t.name, t.amounts[i], id, lastBid)
				lastBidder = buyer
				lastBid = t.amounts[i]
			}
			return err
		}

		test.EvaluateTest(s.T(), t.name, transfer)
	}
	// Revert the previously set commission
	s.configKeeper.SetConfig(s.ctx, defaultConfig)
}

func (s *EscrowTestSuite) TestRefund() {
	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(50)))
	lowerPrice := price.Sub(sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.OneInt())))
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
				false,
			)
			s.keeper.MarkExpiredEscrows(s.ctx, lastBlockTime)

			if err != nil {
				panic(err)
			}
			escrowIds[i][j] = id
		}
	}

	passedDeadline :=  s.generator.NowAfter(0) - 10
	s.keeper.SetLastBlockTime(s.ctx, passedDeadline-1)

	createAuction := func(price sdk.Coins, deadline uint64) string {
		id, err := s.keeper.CreateEscrow(
		s.ctx, s.seller, price, newSavedObject(s.generator, s.seller, s.store), deadline, true)
		if err != nil {
			panic(err)
		}
		return id
	}

	expiredAuctionWithoutBidsId := createAuction(price, passedDeadline)
	expiredAuctionWithBidsId := createAuction(lowerPrice, passedDeadline)
	err := s.keeper.TransferToEscrow(s.ctx, s.buyer, expiredAuctionWithBidsId, price)
	if err != nil {
		panic(err)
	}

	s.keeper.MarkExpiredEscrows(s.ctx, lastBlockTime)
	s.keeper.SetLastBlockTime(s.ctx, lastBlockTime)

	auctionWithoutBidsId := createAuction(price, s.generator.NowAfter(10))
	auctionWithBidsId := createAuction(lowerPrice, s.generator.NowAfter(10))
	err = s.keeper.TransferToEscrow(s.ctx, s.buyer, auctionWithBidsId, price)
	if err != nil {
		panic(err)
	}

	invalidObjectEscrowId := s.createErroredObjectEscrow(price, false)

	validRefund := func(before, after assetState, escrowId string) {
		s.Assert().Equal(before.sellerBalance, after.sellerBalance)
		s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(escrowId))
		s.Assert().Equal(after.objectOwner, s.seller)
	}

	invalidRefund := func(before, after assetState, escrowId string) {
		s.Assert().Equal(before.sellerBalance, after.sellerBalance)
		// If this is has sense
		if before.objectOwner != nil {
			s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(escrowId))
			s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress(escrowId))
		}
	}

	testCases := []struct {
		name   string
		id     string
		sender sdk.AccAddress
		check  func(before, after assetState, escrowId string)
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
			name: "valid refund: trigger by seller before expiration in an auction without bids",
			sender: s.seller,
			id: auctionWithoutBidsId,
			check:  validRefund,
		},
		{
			name: "valid refund: trigger by seller after expiration in an auction without bids",
			sender: s.seller,
			id: expiredAuctionWithoutBidsId,
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
			check:  invalidRefund,
		},
		{
			name: "invalid refund: auction with bids",
			sender: s.seller,
			id: auctionWithBidsId,
			check:  invalidRefund,
		},
		{
			name: "invalid refund: expired auction with bids",
			sender: s.seller,
			id: expiredAuctionWithBidsId,
			check:  invalidRefund,
		},
	}

	for _, t := range testCases {
		check := func(before, after assetState, id string) {}
		if t.check != nil {
			check = t.check
		}

		refund := func(*testing.T) error {
			escrow, found := s.keeper.GetEscrow(s.ctx, t.id)
			before := s.getState(s.seller, nil, nil, escrow, found)
			err := s.keeper.RefundEscrow(s.ctx, t.sender, t.id)
			check(before, s.getState(s.seller, nil, nil, escrow, found), t.id)
			return err
		}

		test.EvaluateTest(s.T(), t.name, refund)
	}
}

func (s *EscrowTestSuite) TestCompleteAuction() {
	var priceAmt int64 = 50
	price := sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(priceAmt)))
	lastBlockTime := s.keeper.GetLastBlockTime(s.ctx)
	defaultConfig := s.configKeeper.GetConfiguration(s.ctx)
	brokerAddr := s.generator.NewAccAddress()

	var expiredAuctionIds [3]string
	for i := 0; i < 3; i++ {
		deadline := lastBlockTime - 5
		// Cheat a little to be able to create the escrows in the past
		s.keeper.SetLastBlockTime(s.ctx, deadline-1)
		id, err := s.keeper.CreateEscrow(
			s.ctx,
			s.seller,
			price,
			newSavedObject(s.generator, s.seller, s.store),
			deadline,
			true,
		)
		if err != nil {
			panic(err)
		}

		for j := 0; j < i; j++ {
			addr := s.generator.NewAccAddress()
			s.balances[addr.String()]= price.Add(price...)
			err := s.keeper.TransferToEscrow(s.ctx, addr, id, price.Add(sdk.NewCoin(test.Denom, sdk.NewInt(int64(j+1)))))
			if err != nil {
				panic(err)
			}
		}

		s.keeper.SetLastBlockTime(s.ctx, lastBlockTime)
		s.keeper.MarkExpiredEscrows(s.ctx, lastBlockTime)

		expiredAuctionIds[i] = id
	}


	invalidObjectAuctionId := s.createErroredObjectEscrow(price, true)

	notExpiredAuctionId, err := s.keeper.CreateEscrow(
		s.ctx, s.seller, price, newSavedObject(s.generator, s.seller, s.store), s.generator.NowAfter(10), true,
		)
	if err != nil {
		panic(err)
	}

	createAndSaveAuctionWithBroker := func(seller sdk.AccAddress, price sdk.Coins, commission sdk.Dec) string {
		config := defaultConfig
		config.EscrowBroker = brokerAddr.String()
		config.EscrowCommission = commission
		s.configKeeper.SetConfig(s.ctx, config)

		passedDeadline := s.generator.NowAfter(0) - 5
		s.keeper.SetLastBlockTime(s.ctx, passedDeadline - 1)

		id, err := s.keeper.CreateEscrow(
			s.ctx,
			s.seller,
			price,
			newSavedObject(s.generator, s.seller, s.store),
			passedDeadline,
			true,
		)
		s.configKeeper.SetConfig(s.ctx, defaultConfig)

		if err != nil {
			panic(err)
		}
		addr := s.generator.NewAccAddress()
		s.balances[addr.String()] = price.Add(price...)
		err = s.keeper.TransferToEscrow(s.ctx, addr, id, price.Add(sdk.NewCoin(test.Denom, sdk.NewInt(int64(1)))))
		if err != nil {
			panic(err)
		}

		s.keeper.SetLastBlockTime(s.ctx, s.generator.NowAfter(0))
		s.keeper.MarkExpiredEscrows(s.ctx, s.generator.NowAfter(0))
		return id
	}

	getValidCompletionCheck := func (price sdk.Coins) func(before, after assetState, escrowId string, _ sdk.Coins, _ sdk.AccAddress) {
		return func(before, after assetState, escrowId string, expectedCommission sdk.Coins, lastBidder sdk.AccAddress) {
			s.Assert().Equal(before.sellerBalance.Add(price...).Sub(expectedCommission), after.sellerBalance)
			s.Assert().Equal(before.buyerBalance, after.buyerBalance)
			s.Assert().Equal(before.brokerBalance.Add(expectedCommission...).String(), after.brokerBalance.String())

			s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(escrowId))
			s.Assert().Equal(after.objectOwner, lastBidder)
		}
	}

	invalidCompletion := func(before, after assetState, escrowId string, _ sdk.Coins, _ sdk.AccAddress) {
		s.Assert().Equal(before.sellerBalance, after.sellerBalance)
		s.Assert().Equal(before.buyerBalance, after.buyerBalance)
		// If this is has sense
		if before.objectOwner != nil {
			s.Assert().Equal(before.objectOwner, s.keeper.GetEscrowAddress(escrowId))
			s.Assert().Equal(after.objectOwner, s.keeper.GetEscrowAddress(escrowId))
		}
	}

	testCases := []struct {
		name   string
		id     string
		broker sdk.AccAddress
		expectedCommission sdk.Coins
		check  func(before, after assetState, escrowId string, expectedCommission sdk.Coins, _ sdk.AccAddress)
	}{
		{
			name:   "invalid completion: no bids",
			id:     expiredAuctionIds[0],
			check: invalidCompletion,
		},
		{
			name:   "valid completion: 1 bid",
			id:     expiredAuctionIds[1],
			check:  getValidCompletionCheck(price.Add(sdk.NewCoin(test.Denom, sdk.OneInt()))),
		},
		{
			name:   "valid completion: n bids",
			id:     expiredAuctionIds[2],
			check:  getValidCompletionCheck(price.Add(sdk.NewCoin(test.Denom, sdk.NewInt(2)))),
		},
		{
			name:               "valid transfer with broker and decimal escrow commission",
			check:              getValidCompletionCheck(price.Add(sdk.NewCoin(test.Denom, sdk.NewInt(1)))),
			broker:             brokerAddr,
			id:                 createAndSaveAuctionWithBroker(s.seller, price, sdk.NewDec(1).Quo(sdk.NewDec(3))),
			expectedCommission: sdk.NewCoins(sdk.NewCoin(test.Denom, sdk.NewInt(priceAmt/3))),
		},
		{
			name:   "invalid completion: auction not reached deadline",
			id:     notExpiredAuctionId,
			check:  invalidCompletion,
		},
		{
			name:   "invalid escrow: non existing",
			id:     "AABBCCDDEEFF1122",
			check:  invalidCompletion,
		},
		{
			name:   "invalid completion: error on object transfer",
			id:     invalidObjectAuctionId,
		},
	}

	for _, t := range testCases {
		check := func(before, after assetState, id string, _ sdk.Coins, _ sdk.AccAddress) {}
		if t.check != nil {
			check = t.check
		}

		completeAuction := func(*testing.T) error {
			escrow, found := s.keeper.GetEscrow(s.ctx, t.id)
			var lastBidder sdk.AccAddress
			var err error
			if len(escrow.LastBidder) != 0 {
				lastBidder, err = sdk.AccAddressFromBech32(escrow.LastBidder)
				if err != nil {
					panic(err)
				}
			}
			before := s.getState(s.seller, nil, t.broker, escrow, found)
			err = s.keeper.CompleteAuction(s.ctx, t.id)
			check(before,
				s.getState(s.seller, nil, t.broker, escrow, found),
				t.id,
				t.expectedCommission,
				lastBidder,
				)
			return err
		}

		test.EvaluateTest(s.T(), t.name, completeAuction)
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
