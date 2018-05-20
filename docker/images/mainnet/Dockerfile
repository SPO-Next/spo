# spo build
# reference https://github.com/spo-next/spo
FROM golang:1.10 AS build

ARG SPO_VERSION

COPY . $GOPATH/src/github.com/spo-next/spo

# This code checks if SPO_VERSION is set and checkouts to that version if
# so. The git stash line prevents the build to fail if there are any uncommited
# changes in the working copy. It won't affect the host working copy.
RUN sh -c \
    'if test ${SPO_VERSION};then \
        echo "Revision is set to: "${SPO_VERSION}; \
        cd $GOPATH/src/github.com/spo-next/spo; \
        git stash; \
        git checkout ${SPO_VERSION}; \
     fi'

RUN cd $GOPATH/src/github.com/spo-next/spo && \
  CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...

# spo image
FROM alpine:3.7

ENV COIN="spo"
ENV RPC_ADDR="0.0.0.0:8650" \
    DATA_DIR="/data/.$COIN" \
    WALLET_DIR="/wallet" \
    WALLET_NAME="$COIN_cli.wlt"

# copy all the binaries
COPY --from=build /go/bin/* /usr/bin/

# copy assets
COPY --from=build /go/src/github.com/spo-next/spo/src/gui/static /usr/local/spo/src/gui/static

# copy launcher
COPY docker_launcher.sh /usr/local/bin

# volumes
VOLUME $WALLET_DIR
VOLUME $DATA_DIR

EXPOSE 8848 8620 8630

ENTRYPOINT ["docker_launcher.sh", "--web-interface-addr=0.0.0.0", "--rpc-interface-addr=0.0.0.0", "--gui-dir=/usr/local/spo/src/gui/static"]
