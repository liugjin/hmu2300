#!/bin/sh

GOPROXY=https://goproxy.cn GOARCH=arm go build -ldflags "-s -w"
