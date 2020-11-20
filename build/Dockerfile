# go build
FROM golang:1.15.5-alpine3.12 as builder

WORKDIR /go/src
ENV GO111MODULE=on

COPY ./ /go/src/coredns-api/

RUN apk add --no-cache alpine-sdk \
    && cd coredns-api \
    && go get -u github.com/swaggo/swag/cmd/swag \
    && sh scripts/code_build.sh

# image build
FROM alpine:latest

COPY --from=builder /go/src/coredns-api /coredns-api

CMD ["/coredns-api/build/coredns-api"]
