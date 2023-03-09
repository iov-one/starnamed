<<<<<<< HEAD
# docker build . -t iov1/starnamed:latest
# docker run --rm -it iov1/starnamed:latest /bin/sh
FROM golang:1.16.8-alpine3.13 AS go-builder
=======
# docker build . -t cosmwasm/starnamed:latest
# docker run --rm -it cosmwasm/starnamed:latest /bin/sh
FROM golang:1.18-alpine3.15 AS go-builder
ARG arch=x86_64
>>>>>>> tags/v0.11.6

# this comes from standard alpine nightly file
#  https://github.com/rust-lang/docker-rust-nightly/blob/master/alpine3.12/Dockerfile
# with some changes to support our toolchain, etc
RUN set -eux; apk add --no-cache ca-certificates build-base;

RUN apk add git
# NOTE: add these to run with LEDGER_ENABLED=true
# RUN apk add libusb-dev linux-headers

WORKDIR /code
COPY . /code/

# See https://github.com/CosmWasm/wasmvm/releases
<<<<<<< HEAD
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.0.0-beta2/libwasmvm_muslc.a /lib/libwasmvm_muslc.a
RUN sha256sum /lib/libwasmvm_muslc.a | grep 3f5de8df9c6b606b4211f90edd681c84b0ecd870fdbf50678b6d9afd783a571c
=======
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.aarch64.a /lib/libwasmvm_muslc.aarch64.a
ADD https://github.com/CosmWasm/wasmvm/releases/download/v1.1.1/libwasmvm_muslc.x86_64.a /lib/libwasmvm_muslc.x86_64.a
RUN sha256sum /lib/libwasmvm_muslc.aarch64.a | grep 9ecb037336bd56076573dc18c26631a9d2099a7f2b40dc04b6cae31ffb4c8f9a
RUN sha256sum /lib/libwasmvm_muslc.x86_64.a | grep 6e4de7ba9bad4ae9679c7f9ecf7e283dd0160e71567c6a7be6ae47c81ebe7f32

# Copy the library you want to the final location that will be found by the linker flag `-lwasmvm_muslc`
RUN cp /lib/libwasmvm_muslc.${arch}.a /lib/libwasmvm_muslc.a
>>>>>>> tags/v0.11.6

# force it to use static lib (from above) not standard libgo_cosmwasm.so file
RUN LEDGER_ENABLED=false BUILD_TAGS=muslc LINK_STATICALLY=true make build
RUN echo "Ensuring binary is statically linked ..." \
  && (file /code/build/starnamed | grep "statically linked")

# --------------------------------------------------------
<<<<<<< HEAD
FROM alpine:3.13

# add extremely useful commands to the image
RUN apk update && apk upgrade && apk --no-cache add curl jq openssh-client rsync

=======
FROM alpine:3.15

>>>>>>> tags/v0.11.6
COPY --from=go-builder /code/build/starnamed /usr/bin/starnamed

COPY docker/* /opt/
RUN chmod +x /opt/*.sh

WORKDIR /opt

# rest server
EXPOSE 1317
# tendermint p2p
EXPOSE 26656
# tendermint rpc
EXPOSE 26657
# grpc
EXPOSE 9090

<<<<<<< HEAD
CMD ["/usr/bin/starnamed", "version"]
=======
ENTRYPOINT [ "/usr/bin/starnamed" ]
CMD ["version"]
>>>>>>> tags/v0.11.6
