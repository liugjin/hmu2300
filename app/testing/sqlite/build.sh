#!/bin/sh

# sudo aptitude install libsqlite3-dev
# sudo aptitude install gcc # gcc 6
# sudo aptitude install gcc-mips-linux-gnueabi # gcc 6
# sudo aptitude install gcc-arm-linux-gnueabi  # gcc 6

# build for mipsel
echo "building for mipsel"
GOARCH=mipsle CC=mipsel-linux-gnu-gcc CGO_ENABLED=1 go build -o sqlite_mipsel .

echo "building for arm"
GOARCH=arm CC=arm-linux-gnueabi-gcc CGO_ENABLED=1 go build -o sqlite_arm .

echo "building for current system"
go build 
