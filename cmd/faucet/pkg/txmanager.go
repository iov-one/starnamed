package pkg

import (
	"context"
	"math"

	"github.com/cosmos/cosmos-sdk/client"
	clienttx "github.com/cosmos/cosmos-sdk/client/tx"
	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/types/bech32"
	"github.com/cosmos/cosmos-sdk/types/errors"
	authtypes "github.com/cosmos/cosmos-sdk/x/auth/types"
	bank "github.com/cosmos/cosmos-sdk/x/bank/types"
	"google.golang.org/grpc"

	"github.com/iov-one/starnamed/app"
)

type TxManager struct {
	conf       *Configuration
	keys       keyring.Keyring
	faucetAcc  authtypes.AccountI
	grpcClient *grpc.ClientConn
	clientCtx  client.Context
	faucetAddr string
}

func NewTxManager(conf *Configuration, conn *grpc.ClientConn) *TxManager {
	return &TxManager{grpcClient: conn, conf: conf}
}

func (tm *TxManager) WithKeybase(keys keyring.Keyring) *TxManager {
	tm.keys = keys
	return tm
}

func (tm *TxManager) Init() error {
	encodingCfg := app.MakeEncodingConfig()
	tm.clientCtx = client.Context{}.
		WithInterfaceRegistry(encodingCfg.InterfaceRegistry).
		WithTxConfig(encodingCfg.TxConfig).
		WithHomeDir(app.DefaultNodeHome)

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
	_, adjusted, err := clienttx.CalculateGas(tm.grpcClient, fact, txBuilder.GetTx().GetMsgs()...)
	if err != nil {
		return nil, errors.Wrap(err, "Failed to simulate the transaction")
	}
	txBuilder.SetGasLimit(adjusted)

	if adjusted > math.MaxInt64 {
		return nil, errors.New("faucet", 2, "Gas used too high for an integer representation")
	}

	feeAmt := tm.conf.GasPrices.Amount.Mul(sdk.NewDec(int64(adjusted)))
	txBuilder.SetFeeAmount(sdk.NewCoins(sdk.NewCoin(tm.conf.GasPrices.Denom, feeAmt.TruncateInt())))

	// Sign the transaction
	if err := clienttx.Sign(fact, "faucet", txBuilder, true); err != nil {
		return nil, errors.Wrap(err, "Failed to sign the transaction")
	}

	return txconfig.TxEncoder()(txBuilder.GetTx())
}
