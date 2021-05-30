## How to use

Start the server using `faucet [flags]`

Running `faucet -h` display this help message :

```
Usage of faucet:
  -chain-id string
    	The chain ID (default "integration-test")
  -faucet-armor-file string
    	The faucet private key file (default ".faucet_key")
  -gas-adjust float
    	The gas adjustement (default 1.2)
  -gas-price string
    	The gas price (default "0.000001tiov")
  -grpc-endpoint string
    	The address and port of a tendermint node gRPC (default "localhost:9090")
  -listen-port uint
    	The port the faucet HTTP server will listen to (default 8080)
  -memo string
    	The message associated with the transaction (default "Sent with love by IOV")
  -rpc-endpoint string
    	A full address, with protocol and port, of a tendermint node RPC (default "http://localhost:26657")
  -send-amount string
    	Coin to send when receiving a credit request (default "100tiov")
```

When the server is running, you can interact with the HTTP API :
- `http://localhost:8080/health` Sends a 200 OK response if the server is alive
- `http://localhost:8080/credit?address=<targetBech32addr>` Credits the target account with a predefined amount
