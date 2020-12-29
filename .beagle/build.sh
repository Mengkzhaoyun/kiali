#!/bin/bash 

set -ex

mkdir -p dist

export GOARCH=amd64
make build
mv /go/bin/kiali dist/kiali-$GOARCH

export GOARCH=arm64
make build
mv /go/bin/kiali dist/kiali-$GOARCH

export GOARCH=ppc64le
make build
mv /go/bin/kiali dist/kiali-$GOARCH

export GOARCH=mips64le
make build
mv /go/bin/kiali dist/kiali-$GOARCH