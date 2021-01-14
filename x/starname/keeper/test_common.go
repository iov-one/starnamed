package keeper

import (
	"bytes"
	"fmt"
	"testing"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	cdctypes "github.com/cosmos/cosmos-sdk/codec/types"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/mock"
	"github.com/iov-one/starnamed/x/configuration"
	"github.com/iov-one/starnamed/x/starname/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"
)

type DullMsg struct {
	signer sdk.AccAddress
}

func (m *DullMsg) FeePayer() sdk.AccAddress {
	return m.signer
}

func (m *DullMsg) Route() string {
	return "dull"
}

func (m *DullMsg) Type() string {
	return "dull_msg"
}

func (m *DullMsg) ValidateBasic() error {
	return nil
}

// GetSignBytes implements sdk.Msg
func (m *DullMsg) GetSignBytes() []byte {
	return nil
}

// GetSigners implements sdk.Msg
func (m *DullMsg) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{m.signer}
}

// NewTestCodec generates aliceAddr mock codec for keeper module
func NewTestCodec() *codec.ProtoCodec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	interfaceRegistry := cdctypes.NewInterfaceRegistry()
	interfaceRegistry.RegisterInterface("cosmos.base.v1beta1.Msg",
		(*sdk.Msg)(nil),
		&types.MsgRegisterDomain{},
		&types.MsgTransferDomain{},
		&types.MsgTransferAccount{},
		&types.MsgAddAccountCertificates{},
		&types.MsgDeleteAccountCertificate{},
		&types.MsgDeleteAccount{},
		&types.MsgDeleteDomain{},
		&types.MsgRegisterAccount{},
		&types.MsgRenewDomain{},
		&types.MsgReplaceAccountResources{},
		&types.MsgReplaceAccountMetadata{},
	)
	cdc := codec.NewProtoCodec(interfaceRegistry)
	return cdc
}

type Mocks struct {
	Supply *mock.SupplyKeeperMock
}

// NewTestKeeper generates aliceAddr keeper and aliceAddr context from it
func NewTestKeeper(t testing.TB, isCheckTx bool) (Keeper, sdk.Context, *Mocks) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(configuration.StoreKey) // configuration module store key
	domainStoreKey := sdk.NewKVStoreKey(types.DomainStoreKey)          // domain module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	ms.MountStoreWithDB(domainStoreKey, sdk.StoreTypeIAVL, mdb)        // mount domain module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create Mocks
	mocks := new(Mocks)
	// create mock supply keeper
	mocks.Supply = mock.NewSupplyKeeper()
	// create config keeper
	confKeeper := configuration.NewKeeper(cdc, configurationStoreKey, nil)
	// create context
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now()}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, domainStoreKey, confKeeper, mocks.Supply.Mock(), nil), ctx, mocks
}

// CompareAccounts compares two accounts
func CompareAccounts(got, want *types.Account) error {
	if got.Domain != want.Domain {
		return fmt.Errorf("got Domain '%s', want '%s'", got.Domain, want.Domain)
	}
	if *got.Name != *want.Name {
		return fmt.Errorf("got Name '%s', want '%s'", *got.Name, *want.Name)
	}
	if !got.Owner.Equals(want.Owner) {
		return fmt.Errorf("got Owner '%s', want '%s'", got.Owner.String(), want.Owner.String())
	}
	if !got.Broker.Equals(want.Broker) {
		return fmt.Errorf("got Broker '%s', want '%s'", got.Broker.String(), want.Broker.String())
	}
	if got.ValidUntil != want.ValidUntil {
		return fmt.Errorf("got ValidUntil '%d', want '%d'", got.ValidUntil, want.ValidUntil)
	}
	if got.Resources != nil {
		for i, goti := range got.Resources {
			wanti := want.Resources[i]
			if goti.URI == wanti.URI {
				return fmt.Errorf("got URI '%s', want '%s'", goti.URI, wanti.URI)
			}
			if goti.Resource == wanti.Resource {
				return fmt.Errorf("got Resource '%s', want '%s'", goti.Resource, wanti.Resource)
			}
		}
	}
	if got.Certificates != nil {
		for i, goti := range got.Certificates {
			wanti := want.Certificates[i]

			if bytes.Compare(goti, wanti) != 0 {
				return fmt.Errorf("got Certificate '%s', want '%s'", string(goti), string(wanti))
			}
		}
	}
	if got.MetadataURI != want.MetadataURI {
		return fmt.Errorf("got MetadataURI '%s', want '%s'", got.MetadataURI, want.MetadataURI)
	}
	return nil
}
