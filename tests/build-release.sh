#!/bin/bash

set -eux

REALPATH=$(realpath "$0")
ROOT=$(dirname "$REALPATH")
ROOT=$(dirname "$ROOT")
TESTBUILD_DIR="$ROOT/test-build"
OUTPUT="$ROOT/output"

rm -rf "$TESTBUILD_DIR"

export GOPATH="$TESTBUILD_DIR"
export PATH="$GOPATH/bin:$PATH"
DEST="$GOPATH/src/github.com/rekby/lets-proxy2"
mkdir -p "$TESTBUILD_DIR/src/github.com/rekby"
ln -Ts "$ROOT" "$DEST"
cd "$DEST"

go version
go env

go get github.com/mitchellh/gox

if [ "$GO_VERSION" = "1.10" ]; then
  OS_ARCH_BUILDS="linux/386 linux/amd64 linux/arm freebsd/386 freebsd/amd64 freebsd/arm windows/386 windows/amd64"
else
  OS_ARCH_BUILDS="linux/386 linux/amd64 linux/arm freebsd/386 freebsd/amd64 freebsd/arm windows/386 windows/amd64 darwin/amd64"
fi

GITHUB_TAG="${GITHUB_REF/refs\/tags\//}"
BUILD_TIME=$(TZ=UTC date --rfc-3339=seconds)
CGO_ENABLED=0 gox --mod=vendor -osarch "$OS_ARCH_BUILDS" --ldflags "-X \"main.VERSION=$GITHUB_TAG+build-$GITHUB_RUN_ID, Build time $BUILD_TIME, commit $GITHUB_SHA, $GO_VERSION\"" --output="$OUTPUT/lets-proxy_{{.OS}}_{{.Arch}}" -verbose --rebuild ./cmd/
