.PHONY: all install test

# make sure we turn on go modules
export GO111MODULE := on

all: test install

install:
	go install .

build:
	GOARCH=amd64 GOOS=linux go build .

test:
	go vet -mod=readonly ./...
	go test -mod=readonly -race ./...
