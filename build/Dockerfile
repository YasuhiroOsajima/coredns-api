# build
FROM golang:1.15.5-alpine3.12 as builder

WORKDIR /go/src
ENV GO111MODULE=on

RUN mkdir /gobuild
RUN apk add --no-cache alpine-sdk git \
    && git clone https://github.com/YasuhiroOsajima/coredns-api.git \
    && cd coredns-api \
    && go get -u github.com/swaggo/swag/cmd/swag \
    && sh scripts/code_build.sh

# run
FROM alpine:latest

COPY --from=builder /go/src/coredns-api /coredns-api

CMD ["/coredns-api/build/coredns-api"]
