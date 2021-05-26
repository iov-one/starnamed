package pkg

import (
	"context"
	"github.com/iov-one/starnamed/app"
	clientcodec "github.com/iov-one/starnamed/x/wasm/client/codec"
	"os"
	"sync"

	"github.com/cosmos/cosmos-sdk/client"
	"github.com/cosmos/cosmos-sdk/client/flags"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/types/tx/signing"
	authsigning "github.com/cosmos/cosmos-sdk/x/auth/signing"

	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"

	"github.com/cosmos/cosmos-sdk/types/errors"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client"
)

type TxManager struct {
	conf      Configuration
	node      rpchttp.Client
	keys      keyring.Keyring
	faucetAcc authtypes.AccountI
	mux       sync.Mutex
	clientCtx client.Context
}

func (tm *TxManager) queryWithData(path string, data []byte) ([]byte, int64, error) {
	res, err := tm.node.ABCIQuery(context.TODO(), path, data)
	if err != nil {
		return nil, 0, err
	}
	return res.Response.Value, res.Response.Height, nil
}

func NewTxManager(conf Configuration, node rpchttp.Client) *TxManager {
	return &TxManager{node: node, conf: conf}
}

func (tm *TxManager) WithKeybase(keys keyring.Keyring) *TxManager {
	tm.keys = keys
	return tm
}

func (tm *TxManager) Init() error {
	tm.clientCtx = client.Context{}.
		// WithJSONMarshaler(ModuleCdc.Marshaler).
		// The following line was taken from wasmd implementation
		WithJSONMarshaler(clientcodec.NewProtoCodec(ModuleCdc.Marshaler, ModuleCdc.InterfaceRegistry)).
		WithInterfaceRegistry(ModuleCdc.InterfaceRegistry).
		WithTxConfig(ModuleCdc.TxConfig).
		WithLegacyAmino(ModuleCdc.Amino).
		WithInput(os.Stdin).
		WithAccountRetriever(authtypes.AccountRetriever{}).
		WithBroadcastMode(flags.BroadcastSync).
		WithHomeDir(app.DefaultNodeHome).
		WithClient(tm.node)

	info, err := tm.keys.Key("faucet")
	if err != nil {
		return err
	}
	// fetch account info
	// we need these step to fetch account sequence on live chain
	acc, err := tm.fetchAccount(info.GetAddress())
	if err != nil {
		return err
	}
	tm.faucetAcc = acc
	return nil
}

func (tm *TxManager) fetchAccount(addr sdk.AccAddress) (authtypes.AccountI, error) {
	// This retrieves the read-only version of an account
	// account, err := tm.clientCtx.AccountRetriever.GetAccount(tm.clientCtx, addr)

	param := authtypes.QueryAccountRequest{Address: addr.String()}
	paramBytes, err := ModuleCdc.Marshaler.MarshalJSON(&param)

	if err != nil {
		return nil, errors.Wrap(err, "Marshal account address")
	}
	req := "/cosmos/auth/v1beta1/accounts/" + addr.String()
	query, err := tm.clientCtx.Client.ABCIQuery(context.TODO(), req, paramBytes)
	if err != nil {
		return nil, errors.Wrap(err, "RPC account request")
	}

	var account authtypes.AccountI
	err = ModuleCdc.Marshaler.UnmarshalInterfaceJSON(query.Response.Value, account)
	return account, err
}

func (tm *TxManager) BroadcastTx(tx []byte) (*sdk.TxResponse, error) {
	return tm.clientCtx.BroadcastTx(tx)
}

func (tm *TxManager) BuildAndSignTx(targetAcc sdk.AccAddress) ([]byte, error) {
	/* CONTRACT
	a faucet wallet must be used by single actor otherwise successful tx will bump
	account sequence on chain.
	*/
	//

	tm.mux.Lock()
	seq := tm.faucetAcc.GetSequence()
	err := tm.faucetAcc.SetSequence(seq + 1)
	tm.mux.Unlock()
	if err != nil {
		return nil, err
	}

	txconfig := ModuleCdc.TxConfig

	signMode := txconfig.SignModeHandler().DefaultMode()

	txBuilder := txconfig.NewTxBuilder()
	signature := signing.SignatureV2{
		PubKey: tm.faucetAcc.GetPubKey(),
		Data: &signing.SingleSignatureData{
			SignMode: signMode,
		},
		Sequence: seq,
	}

	if err := txBuilder.SetSignatures(signature); err != nil {
		return nil, errors.Wrap(err, "Signing transaction")
	}
	txBuilder.SetMemo(tm.conf.Memo)

	sendMsg := bank.MsgSend{
		FromAddress: tm.faucetAcc.GetAddress().String(),
		ToAddress:   targetAcc.String(),
		Amount: sdk.Coins{
			sdk.NewInt64Coin(tm.conf.CoinDenom, tm.conf.SendAmount),
		},
	}
	signerData := authsigning.SignerData{
		ChainID:       tm.conf.ChainID,
		AccountNumber: tm.faucetAcc.GetAccountNumber(),
		Sequence:      tm.faucetAcc.GetSequence(),
	}
	signBytes, err := txconfig.SignModeHandler().GetSignBytes(signMode, signerData, txBuilder.GetTx())
	if err != nil {
		panic(err)
	}
	sig, _, err := tm.keys.Sign("faucet", signBytes)
	if err != nil {
		panic(err)
	}
	signature.Data.(*signing.SingleSignatureData).Signature = sig
	err = txBuilder.SetSignatures(signature)
	if err != nil {
		panic(err)
	}

	// adjust gas
	err = txBuilder.SetMsgs(&sendMsg)
	if err != nil {
		return nil, errors.Wrap(err, "tx gas adjustment failed")
	}

	transaction := txBuilder.GetTx()
	// TODO : should we provide flagSet ?
	fact := clienttx.NewFactoryCLI(tm.clientCtx, nil)
	_, adjusted, err := clienttx.CalculateGas(tm.queryWithData, fact, transaction.GetMsgs()...)
	txBuilder.SetGasLimit(adjusted)
	// TODO : Compute fee amount from gas prices and gas limit
	//txBuilder.SetFeeAmount(tm.conf.GasPrices.Amount.Uint64() * adjusted)

	return txconfig.TxEncoder()(transaction)
}
