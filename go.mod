module github.com/iov-one/starnamed

go 1.15

require (
	github.com/CosmWasm/wasmvm v0.12.0
	github.com/cosmos/cosmos-sdk v0.40.0-rc3
	github.com/cosmos/iavl v0.15.0-rc5
	github.com/dvsekhvalnov/jose2go v0.0.0-20200901110807-248326c1351b
	github.com/fatih/structs v1.1.0
	github.com/gogo/protobuf v1.3.1
	github.com/golang/protobuf v1.4.3
	github.com/google/gofuzz v1.0.0
	github.com/gorilla/mux v1.8.0
	github.com/grpc-ecosystem/grpc-gateway v1.16.0
	github.com/iov-one/cosmos-sdk-crud v0.0.0-00010101000000-000000000000
	github.com/pkg/errors v0.9.1
	github.com/rakyll/statik v0.1.7
	github.com/snikch/goodman v0.0.0-20171125024755-10e37e294daa
	github.com/spf13/cast v1.3.1
	github.com/spf13/cobra v1.1.1
	github.com/spf13/pflag v1.0.5
	github.com/spf13/viper v1.7.1
	github.com/stretchr/testify v1.6.1
	github.com/tendermint/tendermint v0.34.0
	github.com/tendermint/tm-db v0.6.3
	google.golang.org/genproto v0.0.0-20201111145450-ac7456db90a6
	google.golang.org/grpc v1.33.2
	gopkg.in/yaml.v2 v2.3.0
)

replace github.com/gogo/protobuf => github.com/regen-network/protobuf v1.3.2-alpha.regen.4

// TODO: drop the following when possible
replace github.com/iov-one/cosmos-sdk-crud => ../cosmos-sdk-crud
