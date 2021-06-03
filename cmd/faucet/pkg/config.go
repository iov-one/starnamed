package pkg

import (
	"fmt"
	"math"
	"os"
	"strings"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/iov-one/starnamed/app"
	"github.com/pkg/errors"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
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
	Passphrase            string
}

func ParseConfiguration() (*Configuration, error) {

	sdk.GetConfig().SetBech32PrefixForAccount(app.Bech32PrefixAccAddr, app.Bech32PrefixAccPub)

	pflag.String("send-amount", "100tiov", "Coin to send when receiving a credit request")
	pflag.String("grpc-endpoint", "localhost:9090", "The address and port of a tendermint node gRPC")
	pflag.String("rpc-endpoint", "http://localhost:26657", "A full address, with protocol and port, of a tendermint node RPC")
	pflag.Uint("listen-port", 8080, "The port the faucet HTTP server will listen to")
	pflag.String("memo", "Sent with love by IOV", "The message associated with the transaction")
	pflag.String("chain-id", "integration-test", "The chain ID")
	pflag.String("armor-file", ".faucet_key", "The faucet private key file")
	pflag.String("gas-price", "0.000001tiov", "The gas price")
	pflag.Float64("gas-adjust", 1.2, "The gas adjustement")

	pflag.Parse()

	err := viper.BindPFlags(pflag.CommandLine)
	if err != nil {
		return nil, errors.Wrap(err, "Could not read command-line arguments")
	}
	// Bind environment variables
	viper.SetEnvPrefix("FAUCET")
	viper.SetEnvKeyReplacer(strings.NewReplacer("-", "_"))
	viper.AutomaticEnv()

	// Validate server listening port
	if viper.GetUint("listen-port") > math.MaxUint16 {
		return nil, fmt.Errorf("invalid port number : %v", viper.GetUint("listen-port"))
	}

	_, err = os.Stat(viper.GetString("armor-file"))
	if err != nil {
		return nil, errors.Wrapf(err, "provide a valid faucet private key file: %v", viper.GetString("armor-file"))
	}

	// Parse and validate send amount
	amt, err := sdk.ParseCoinNormalized(viper.GetString("send-amount"))
	if err != nil {
		return nil, errors.Wrapf(err, "provide a valid coin amount")
	}

	if amt.IsNegative() {
		return nil, errors.Wrapf(err, "could not send a negative amount")
	}

	gasPrice, err := sdk.ParseDecCoin(viper.GetString("gas-price"))
	if err != nil {
		return nil, errors.Wrapf(err, "provide a valid gas price")
	}

	if gasPrice.IsNegative() {
		return nil, errors.Wrapf(err, "could not use a negative gas price")
	}

	return &Configuration{
		GRPCEndpoint:          viper.GetString("grpc-endpoint"),
		TendermintRPCEndpoint: viper.GetString("rpc-endpoint"),
		Port:                  viper.GetUint("listen-port"),
		ChainID:               viper.GetString("chain-id"),
		ArmorFile:             viper.GetString("armor-file"),
		Memo:                  viper.GetString("memo"),
		SendAmount:            amt,
		GasPrices:             gasPrice,
		Passphrase:            viper.GetString("armor-passphrase"),
		GasAdjust:             viper.GetFloat64("gas-adjust"),
	}, nil
}
