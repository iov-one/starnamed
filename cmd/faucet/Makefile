.PHONY: all install test

IMAGE_FAUCET ?= "iov1/faucet:latest"
WASMVM := $(shell go list -f '{{ .Dir }}' -m github.com/CosmWasm/wasmvm)

# make sure we turn on go modules
export GO111MODULE := on

all: test install

install:
	go install .

build:
	GOARCH=amd64 GOOS=linux go build .

docker: build
	cp -af $(WASMVM)/api/libwasmvm.so .
	docker build --pull --tag ${IMAGE_FAUCET} .

test:
	go vet -mod=readonly ./...
	go test -mod=readonly -race ./...
