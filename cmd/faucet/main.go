package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"strconv"
	"time"

	"github.com/cosmos/cosmos-sdk/crypto/keyring"
	"github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/gorilla/mux"
	"github.com/iov-one/starnamed/cmd/faucet/pkg"
	"github.com/tendermint/crypto/ssh/terminal"
	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"google.golang.org/grpc"
)

func main() {
	conf, err := pkg.ParseConfiguration()
	if err != nil {
		log.Fatalf("Invalid argument : %v", err)
	}

	// setup context for gRPC connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// setup Tendermint RPC and gRPC connection
	grpcClient, err := grpc.DialContext(ctx, conf.GRPCEndpoint, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("grpc connection: %s", err)
	}

	rpcClient, err := rpchttp.New(conf.TendermintRPCEndpoint, "/websocket")
	if err != nil {
		log.Fatalf("Connection to RPC node failed: %v", err)
	}

	keys := keyring.NewInMemory()

	armor, err := ioutil.ReadFile(conf.ArmorFile)
	if err != nil {
		panic(errors.Wrapf(err, "Cannot read the private key: %v", conf.ArmorFile))
	}

	var passphrase string
	if len(conf.Passphrase) == 0 {
		fmt.Printf("[%v key] Passphrase : ", conf.ArmorFile)
		passphraseBytes, err := terminal.ReadPassword(int(os.Stdin.Fd()))
		if err != nil {
			panic(err)
		}
		fmt.Print("\n")
		passphrase = string(passphraseBytes)
		// Garbage collect the passphrase bytes
		passphraseBytes = nil
	} else {
		passphrase = conf.Passphrase
	}

	if err := keys.ImportPrivKey("faucet", string(armor), passphrase); err != nil {
		log.Fatalf("keybase: %v", err)
	}
	// Erase the passphrase
	passphrase = ""

	// setup tx manager
	txManager := pkg.NewTxManager(conf, grpcClient, rpcClient).WithKeybase(keys)
	if err := txManager.Init(); err != nil {
		log.Fatalf("tx manager: %v", err)
	}

	// Wait for ListenAndServe goroutine to close.
	r := mux.NewRouter()
	faucet := pkg.NewFaucetHandler(txManager)
	r.Handle("/credit", faucet)
	r.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		_, err := w.Write([]byte("OK"))
		if err != nil {
			log.Printf("Error while handling a '/health' request: %v\n", err)
		}
		return
	})
	server := &http.Server{Addr: ":" + strconv.FormatUint(uint64(conf.Port), 10), Handler: r}

	go func() {
		log.Print("server started")
		if err := server.ListenAndServe(); err != nil {
			log.Fatalf("http server: %s", err)
		}
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)

	<-stop

	ctx, cancel = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(ctx)
}
