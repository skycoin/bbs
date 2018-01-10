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
RUN npm install -g --unsafe @angular/cli && \
    cd /bbs/static && \
    yarn && \
    npm run build

## Skycoin BBS Image
FROM alpine:3.7

RUN adduser -D skycoin
RUN mkdir /data
RUN chown skycoin: /data

USER skycoin

# copy binaries & static files
COPY --from=build-go /go/bin/* /usr/bin/
COPY --from=build-node /bbs/static/dist /usr/local/bbs/static
# volumes

VOLUME /data
WORKDIR /

EXPOSE 8080 8998

CMD [ \
    "bbsnode", \
    "--public=true", \
    "--memory=false", \
    "--config-dir=/data", \
    "--cxo-port=8998", \
    "--enforced-messenger-addresses=35.227.102.45:8005", \
    "--enforced-subscriptions=03588a2c8085e37ece47aec50e1e856e70f893f7f802cb4f92d52c81c4c3212742", \
    "--web-port=8080", \
    "--web-gui-dir=/usr/local/bbs/static" \
]