package configuration

import (
	"testing"
	"time"

	"github.com/iov-one/wasmd/x/configuration/types"
	"github.com/cosmos/cosmos-sdk/codec"
	cryptocodec "github.com/cosmos/cosmos-sdk/crypto/codec"
	"github.com/cosmos/cosmos-sdk/store"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
	tmproto "github.com/tendermint/tendermint/proto/tendermint/types"
	db "github.com/tendermint/tm-db"
)

// NewTestCodec generates aliceAddr mock codec for keeper module
func NewTestCodec() *codec.AminoCodec {
	// we should register this codec for all the modules
	// that are used and referenced by domain module
	amino := codec.NewLegacyAmino()
	cdc := codec.NewAminoCodec(amino)
	types.RegisterLegacyAminoCodec(amino)
	cryptocodec.RegisterCrypto(amino)
	amino.Seal()
	return cdc
}

// NewTestKeeper generates aliceAddr keeper and aliceAddr context from it
func NewTestKeeper(t testing.TB, isCheckTx bool) (Keeper, sdk.Context) {
	cdc := NewTestCodec()
	// generate store
	mdb := db.NewMemDB()
	// generate multistore
	ms := store.NewCommitMultiStore(mdb)
	// generate store keys
	configurationStoreKey := sdk.NewKVStoreKey(StoreKey) // configuration module store key
	// generate sub store for each module referenced by the keeper
	ms.MountStoreWithDB(configurationStoreKey, sdk.StoreTypeIAVL, mdb) // mount configuration module
	// test no errors
	require.Nil(t, ms.LoadLatestVersion())
	// create context
	ctx := sdk.NewContext(ms, tmproto.Header{Time: time.Now()}, isCheckTx, log.NewNopLogger())
	// create domain.Keeper
	return NewKeeper(cdc, configurationStoreKey, nil), ctx
}
