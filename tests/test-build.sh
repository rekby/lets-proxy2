#!/bin/bash

set -eux

REALPATH=$(realpath "$0")
ROOT=$(dirname "$REALPATH")
ROOT=$(dirname "$ROOT")
TESTBUILD_DIR="$ROOT/test-build"

rm -rf "$TESTBUILD_DIR"

export GOPATH="$TESTBUILD_DIR"
DEST="$GOPATH/src/github.com/rekby/lets-proxy2"
mkdir -p "$TESTBUILD_DIR/src/github.com/rekby"
ln -Ts "$ROOT" "$DEST"
cd "$DEST"

go version
go env
go build -v ./...
