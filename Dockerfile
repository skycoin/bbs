## Skycoin BBS Binaries
FROM golang:1.9-alpine AS build-go

COPY . $GOPATH/src/github.com/skycoin/bbs

RUN cd $GOPATH/src/github.com/skycoin/bbs && \
    CGO_ENABLED=0 GOOS=linux go install -a -installsuffix cgo ./...

## Skycoin BBS GUI
FROM node:8.9 AS build-node

COPY . /bbs

# `unsafe` flag used as work around to prevent infinite loop in Docker
# see https://github.com/nodejs/node-gyp/issues/1236
RUN cd /bbs/static && \
    rm -rf package-lock.json && \
    npm install --unsafe -g && \
    npm run build

## Skycoin BBS Image
FROM alpine:3.7

ENV DATA_DIR=/data \
    MESSENGER_ADDR=35.227.102.45:8005 \
    RPC_PORT=8996 \
    CXO_PORT=8998 \
    WEB_PORT=8080

RUN adduser -D skycoin
RUN mkdir $DATA_DIR
RUN chown skycoin: $DATA_DIR

USER skycoin

# copy binaries & static files
COPY --from=build-go /go/bin/* /usr/bin/
COPY --from=build-node /bbs/static/dist /usr/local/bbs/static
# volumes

VOLUME $DATA_DIR
WORKDIR /

EXPOSE $RPC_PORT $CXO_PORT $WEB_PORT

CMD bbsnode \
    --public=true \
    --memory=false \
    --config-dir=$DATA_DIR \
    --rpc=true \
    --rpc-port=$RPC_PORT \
    --cxo-port=$CXO_PORT \
    --enforced-messenger-addresses=$MESSENGER_ADDR \
    --web-port=$WEB_PORT \
    --web-gui-dir=/usr/local/bbs/static