package pkg

import (
	"context"

	"github.com/cosmos/cosmos-sdk/types/errors"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"github.com/iov-one/starnamed/app"
	clientcodec "github.com/iov-one/starnamed/x/wasm/client/codec"
	abci "github.com/tendermint/tendermint/abci/types"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
)

type TxManager struct {
	conf       *Configuration
	keys       keyring.Keyring
	faucetAcc  authtypes.AccountI
	grpcClient *grpc.ClientConn
	rpcClient  *rpchttp.HTTP
	clientCtx  client.Context
	faucetAddr string
}

func NewTxManager(conf *Configuration, conn *grpc.ClientConn, rpcClient *rpchttp.HTTP) *TxManager {
	return &TxManager{grpcClient: conn, conf: conf, rpcClient: rpcClient}
}

func (tm *TxManager) WithKeybase(keys keyring.Keyring) *TxManager {
	tm.keys = keys
	return tm
}

func (tm *TxManager) Init() error {
	encodingCfg := app.MakeEncodingConfig()
	tm.clientCtx = client.Context{}.
		WithJSONMarshaler(clientcodec.NewProtoCodec(encodingCfg.Marshaler, encodingCfg.InterfaceRegistry)).
		WithInterfaceRegistry(encodingCfg.InterfaceRegistry).
		WithTxConfig(encodingCfg.TxConfig).
		WithHomeDir(app.DefaultNodeHome).
		WithClient(tm.rpcClient)

	info, err := tm.keys.Key("faucet")
	if err != nil {
		return errors.Wrap(err, "Cannot retrieve the faucet key")
	}
	// fetch account info
	// we need these step to fetch account sequence on live chain

	tm.faucetAddr, err = bech32.ConvertAndEncode(app.Bech32Prefix, info.GetAddress())
	if err != nil {
		return errors.Wrap(err, "Cannot retrieve the faucet address from the faucet key")
	}

	err = tm.refreshAccount()
	if err != nil {
		return errors.Wrap(err, "Cannot retrieve the faucet account")
	}
	return nil
}

func (tm *TxManager) refreshAccount() error {
	param := authtypes.QueryAccountRequest{Address: tm.faucetAddr}
	query := authtypes.NewQueryClient(tm.grpcClient)
	response, err := query.Account(context.TODO(), &param)
	if err != nil {
		return errors.Wrap(err, "gRPC account request")
	}

	var account authtypes.AccountI
	err = tm.clientCtx.InterfaceRegistry.UnpackAny(response.GetAccount(), &account)
	if err != nil {
		return errors.Wrap(err, "Unpacking account info")
	}
	tm.faucetAcc = account
	return nil
}

func (tm *TxManager) queryWithData(path string, data []byte) ([]byte, int64, error) {
	res, err := tm.clientCtx.Client.ABCIQuery(context.TODO(), path, data)
	if err != nil {
		panic(err)
	}
	if res.Response.Code != abci.CodeTypeOK {
		return nil, 0, errors.New("faucet", 1, "Error while simulating the transaction")
	}
	return res.Response.Value, res.Response.Height, nil
}

func (tm *TxManager) BroadcastTx(tx []byte) (*sdk.TxResponse, error) {
	resp, err := tm.clientCtx.BroadcastTxCommit(tx)
	return resp, err
}

func (tm *TxManager) BuildAndSignTx(targetAddr string) ([]byte, error) {

	// Refresh account to sync sequence number
	// This could be done by incrementing the sequence number, but this would get out of sync if the faucet account
	// emits other transaction in the meantime
	err := tm.refreshAccount()

	seq := tm.faucetAcc.GetSequence()

	txconfig := tm.clientCtx.TxConfig

	fact := clienttx.Factory{}.
		WithMemo(tm.conf.Memo).
		WithChainID(tm.conf.ChainID).
		WithKeybase(tm.keys).
		WithTxConfig(tm.clientCtx.TxConfig).
		WithGasAdjustment(tm.conf.GasAdjust).
		WithSequence(seq).
		WithAccountNumber(tm.faucetAcc.GetAccountNumber())

	txBuilder := txconfig.NewTxBuilder()

	// Set message
	sendMsg := bank.MsgSend{
		FromAddress: tm.faucetAddr,
		ToAddress:   targetAddr,
		Amount:      sdk.NewCoins(tm.conf.SendAmount),
	}
	err = txBuilder.SetMsgs(&sendMsg)
	if err != nil {
		return nil, errors.Wrap(err, "tx gas adjustment failed")
	}

	// Set memo
	txBuilder.SetMemo(tm.conf.Memo)

	// Set gas limit and fee amount
	_, adjusted, err := clienttx.CalculateGas(tm.queryWithData, fact, txBuilder.GetTx().GetMsgs()...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to simulate the transaction")
	}
	txBuilder.SetGasLimit(adjusted)

	feeAmt := sdk.NewInt(int64(adjusted * tm.conf.GasPrices.Amount.Uint64()))
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(tm.conf.GasPrices.Denom, feeAmt)))

	// Sign the transaction
	if err := clienttx.Sign(fact, "faucet", txBuilder, true); err != nil {
		return nil, errors.Wrap(err, "Failed to sign the transaction")
	}

	return txconfig.TxEncoder()(txBuilder.GetTx())
}
