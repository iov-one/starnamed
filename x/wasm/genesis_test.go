package wasm_test

import (
	"testing"

	keepertest "github.com/iov-one/starnamed/testutil/keeper"
	"github.com/iov-one/starnamed/x/wasm"
	"github.com/iov-one/starnamed/x/wasm/types"
	"github.com/stretchr/testify/require"
)

<<<<<<< HEAD
func TestGenesis(t *testing.T) {
	genesisState := types.GenesisState{
		// this line is used by starport scaffolding # genesis/test/state
=======
func TestInitGenesis(t *testing.T) {
	data := setupTest(t)

	deposit := sdk.NewCoins(sdk.NewInt64Coin("denom", 100000))
	topUp := sdk.NewCoins(sdk.NewInt64Coin("denom", 5000))
	creator := data.faucet.NewFundedRandomAccount(data.ctx, deposit.Add(deposit...)...)
	fred := data.faucet.NewFundedRandomAccount(data.ctx, topUp...)

	h := data.module.Route().Handler()
	q := data.module.LegacyQuerierHandler(nil)

	msg := MsgStoreCode{
		Sender:       creator.String(),
		WASMByteCode: testContract,
>>>>>>> tags/v0.11.6
	}

	k, ctx := keepertest.WasmKeeper(t)
	wasm.InitGenesis(ctx, *k, genesisState)
	got := wasm.ExportGenesis(ctx, *k)
	require.NotNil(t, got)

<<<<<<< HEAD
	// this line is used by starport scaffolding # genesis/test/assert
=======
	_, _, bob := keyPubAddr()
	initMsg := initMsg{
		Verifier:    fred,
		Beneficiary: bob,
	}
	initMsgBz, err := json.Marshal(initMsg)
	require.NoError(t, err)

	initCmd := MsgInstantiateContract{
		Sender: creator.String(),
		CodeID: firstCodeID,
		Msg:    initMsgBz,
		Funds:  deposit,
		Label:  "testing",
	}
	res, err = h(data.ctx, &initCmd)
	require.NoError(t, err)
	contractBech32Addr := parseInitResponse(t, res.Data)

	execCmd := MsgExecuteContract{
		Sender:   fred.String(),
		Contract: contractBech32Addr,
		Msg:      []byte(`{"release":{}}`),
		Funds:    topUp,
	}
	res, err = h(data.ctx, &execCmd)
	require.NoError(t, err)
	// from https://github.com/CosmWasm/cosmwasm/blob/master/contracts/hackatom/src/contract.rs#L167
	assertExecuteResponse(t, res.Data, []byte{0xf0, 0x0b, 0xaa})

	// ensure all contract state is as after init
	assertCodeList(t, q, data.ctx, 1)
	assertCodeBytes(t, q, data.ctx, 1, testContract)

	assertContractList(t, q, data.ctx, 1, []string{contractBech32Addr})
	assertContractInfo(t, q, data.ctx, contractBech32Addr, 1, creator)
	assertContractState(t, q, data.ctx, contractBech32Addr, state{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
		Funder:      creator.String(),
	})

	// export into genstate
	genState := ExportGenesis(data.ctx, &data.keeper)

	// create new app to import genstate into
	newData := setupTest(t)
	q2 := newData.module.LegacyQuerierHandler(nil)

	// initialize new app with genstate
	InitGenesis(newData.ctx, &newData.keeper, *genState, newData.stakingKeeper, newData.module.Route().Handler())

	// run same checks again on newdata, to make sure it was reinitialized correctly
	assertCodeList(t, q2, newData.ctx, 1)
	assertCodeBytes(t, q2, newData.ctx, 1, testContract)

	assertContractList(t, q2, newData.ctx, 1, []string{contractBech32Addr})
	assertContractInfo(t, q2, newData.ctx, contractBech32Addr, 1, creator)
	assertContractState(t, q2, newData.ctx, contractBech32Addr, state{
		Verifier:    fred.String(),
		Beneficiary: bob.String(),
		Funder:      creator.String(),
	})
>>>>>>> tags/v0.11.6
}
