#!/usr/bin/env bash

DIR="$(mktemp -d)"
function remove-on-exit() {
  rm -rf "$DIR"
}
trap remove-on-exit EXIT

export GOBIN="$DIR"
go get google.golang.org/grpc@v1.22.1
go install google.golang.org/grpc

protoc *.proto --go_out=plugins=grpc:. --plugin=grpc="$DIR/grpc"
