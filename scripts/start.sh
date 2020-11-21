#!/bin/sh

PWD=$(pwd)
export SERVER=172.28.21.40
export PORT=8080
export CONF_PATH=${PWD}/coredns_conf/coredns.conf
export HOSTS_DIR=${PWD}/coredns_conf/hosts/

./build/coredns-api
