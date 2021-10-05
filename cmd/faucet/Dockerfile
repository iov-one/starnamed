FROM alpine:3.14

RUN apk update && apk upgrade && apk --no-cache add curl gcompat libgcc

COPY libwasmvm.so /usr/lib/
COPY faucet /usr/bin/

CMD ["faucet"]
