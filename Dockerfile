FROM golang:1.12 AS builder

COPY . /go/src/github.com/Difrex/zoorest
WORKDIR /go/src/github.com/Difrex/zoorest

ENV CGO_ENABLED=0

RUN go get -t -v ./...
RUN go build -o /zoorest .

FROM alpine AS final

MAINTAINER Denis Zheleztsov <difrex.punk@gmail.com>

COPY --from=builder /zoorest /bin/zoorest

RUN echo -ne "zoorest\nzoorest\n" | adduser zoorest

USER zoorest

EXPOSE 8889

ENTRYPOINT ["/bin/zoorest"]
