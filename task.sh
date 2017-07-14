#!/usr/bin/env bash

set -ea

OUTPUT_DIR=${OUTPUT_DIR:-$PWD}

set -u

export GOPATH=/gopath
export PATH=$PATH:$GOPATH/bin
go get github.com/onsi/ginkgo/ginkgo
go get github.com/onsi/gomega
go get github.com/sclevine/agouti
service dbus restart

xvfb-run ginkgo