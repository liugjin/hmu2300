#!/bin/sh

GOARCH=mipsle go build -ldflags="-s -w"
