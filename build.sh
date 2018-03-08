#!/bin/bash

export COLLECTD_SRC="./collectd-5.8.0"
export CGO_CPPFLAGS="-I${COLLECTD_SRC}/src/daemon -I${COLLECTD_SRC}/src"
go build -buildmode=c-shared -o ddb.so

