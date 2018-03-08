#!/bin/bash

export COLLECTD_SRC="${PWD}/collectd-5.8.0"
export CGO_CPPFLAGS="-I${COLLECTD_SRC}/src/daemon -I${COLLECTD_SRC}/src"
go build -buildmode=c-shared -o ddb.so plugin/main.go
go build -buildmode=default -o ddb ddb_util/main.go

