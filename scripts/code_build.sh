#!/bin/sh

wire cmd/web/infrastructure/wire.go
go build -o build/coredns-api cmd/web/main.go
swag init -d cmd/web/ -g main.go
