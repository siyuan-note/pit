#!/bin/bash

echo 'Building'

go version
export GO111MODULE=on
export GOPROXY=https://goproxy.io

export GOOS=windows
export GOARCH=amd64
go build -v -o "pit-win.exe" -ldflags "-s -w -H=windowsgui"

export GOOS=darwin
export GOARCH=amd64
go build -v -o "pit-darwin" -ldflags "-s -w"

export GOOS=linux
export GOARCH=amd64
go build -v -o "pit-linux" -ldflags "-s -w"