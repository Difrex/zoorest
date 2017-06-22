#!/bin/bash

cat > Dockerfile.builder <<EOF
FROM golang

MAINTAINER Denis Zheleztsov <difrex.punk@gmail.com>

ENV GOPATH /usr

RUN go get github.com/Difrex/zoorest/rest
RUN cd /usr/src/github.com/Difrex/zoorest && go get -t -v ./...

WORKDIR /usr/src/github.com/Difrex/zoorest

ENTRYPOINT go build -ldflags "-linkmode external -extldflags -static" && mv zoorest /out
EOF

# Build builder
docker build --no-cache -t zoorest_builder -f Dockerfile.builder .
# Build bin
docker run -v $(pwd)/out:/out zoorest_builder

case $1 in alpine)
               docker build -t zoorest -f Dockerfile .
               ;;
            *)
               ;;
esac

