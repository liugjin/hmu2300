#!/bin/sh

## build for arm
# GOARCH=arm sup publish

## build for window
# GOOS=windows GOARCH=amd64 sup publish

## For hmu2000
# GOOS=linux GOARCH=mipsle sup publish

## For hmu2200
# GOOS=linux GOARCH=arm sup publish

## For hmu2500
# GOOS=linux GOARCh=arm64 sup publish

if [ $# -gt 0 ];then
    export GOARCH=$1
fi

echo $GOARCH

make

