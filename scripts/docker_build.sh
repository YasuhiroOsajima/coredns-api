#!/bin/sh

docker build -t coredns-api:0.1 -f build/Dockerfile .
for i in `docker images | awk /none/'{print $3}'`; do docker rmi $i; done

