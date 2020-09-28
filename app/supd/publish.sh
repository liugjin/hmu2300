#!/bin/sh

## build for arm
# GOARCH=arm sup publish

## build for window
# GOOS=windows GOARCH=amd64 sup publish

## For hmu2000
# GOARCH=mipsle sup publish

## For hmu2200
# GOARCH=arm sup publish

## For hmu2500
# GOARCh=arm64 sup publish

if [ $# -gt 0 ];then
    export GOARCH=$1
fi

echo $GOARCH

make

