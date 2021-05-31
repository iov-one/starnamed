package pkg

import (
	"flag"
	"fmt"
	"math"
	"os"

	"github.com/iov-one/starnamed/app"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/pkg/errors"
)

type Configuration struct {
	GRPCEndpoint          string
	TendermintRPCEndpoint string
	Port                  uint
	ChainID               string
	ArmorFile             string
	Memo                  string
	SendAmount            sdk.Coin
	GasPrices             sdk.DecCoin
	GasAdjust             float64
}

func ParseConfiguration() (*Configuration, error) {

	sdk.GetConfig().SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)

	sendAmountPtr := flag.String("send-amount", "100tiov", "Coin to send when receiving a credit request")
	grpcEndpointPtr := flag.String("grpc-endpoint", "localhost:9090", "The address and port of a tendermint node gRPC")
	rpcEndpointPtr := flag.String("rpc-endpoint", "http://localhost:26657", "A full address, with protocol and port, of a tendermint node RPC")
	portPtr := flag.Uint("listen-port", 8080, "The port the faucet HTTP server will listen to")
	memoPtr := flag.String("memo", "Sent with love by IOV", "The message associated with the transaction")
	chainIdPtr := flag.String("chain-id", "integration-test", "The chain ID")
	armorFilePtr := flag.String("faucet-armor-file", ".faucet_key", "The faucet private key file")
	gasPricePtr := flag.String("gas-price", "0.000001tiov", "The gas price")
	gasAdjustPtr := flag.Float64("gas-adjust", 1.2, "The gas adjustement")

	flag.Parse()

	// Validate server listening port
	if *portPtr > math.MaxUint16 {
		return nil, fmt.Errorf("invalid port number : %v", *portPtr)
	}

	_, err := os.Stat(*armorFilePtr)
	if err != nil {
		return nil, errors.Wrapf(err, "provide a valid faucet private key file: %v", *armorFilePtr)
	}

	// Parse and validate send amount
	amt, err := sdk.ParseCoinNormalized(*sendAmountPtr)
	if err != nil {
		return nil, errors.Wrapf(err, "provide a valid coin amount")
	}

	if amt.IsNegative() {
		return nil, errors.Wrapf(err, "could not send a negative amount")
	}

	gasPrice, err := sdk.ParseDecCoin(*gasPricePtr)
	if err != nil {
		return nil, errors.Wrapf(err, "provide a valid gas price")
	}

	if gasPrice.IsNegative() {
		return nil, errors.Wrapf(err, "could not use a negative gas price")
	}

	return &Configuration{
		GRPCEndpoint:          *grpcEndpointPtr,
		TendermintRPCEndpoint: *rpcEndpointPtr,
		Port:                  *portPtr,
		ChainID:               *chainIdPtr,
		ArmorFile:             *armorFilePtr,
		Memo:                  *memoPtr,
		SendAmount:            amt,
		GasPrices:             gasPrice,
		GasAdjust:             *gasAdjustPtr,
	}, nil
}
