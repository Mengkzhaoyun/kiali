#!/bin/bash 

set -ex

mkdir -p dist

go get -u github.com/mjibson/esc

export GOARCH=amd64
make build-all-in-one
mv cmd/all-in-one/all-in-one-linux-$GOARCH dist/all-in-one-linux-$GOARCH

export GOARCH=arm64
make build-all-in-one
mv cmd/all-in-one/all-in-one-linux-$GOARCH dist/all-in-one-linux-$GOARCH

export GOARCH=ppc64le
make build-all-in-one
mv cmd/all-in-one/all-in-one-linux-$GOARCH dist/all-in-one-linux-$GOARCH

export GOARCH=mips64le
make build-all-in-one
mv cmd/all-in-one/all-in-one-linux-$GOARCH dist/all-in-one-linux-$GOARCH