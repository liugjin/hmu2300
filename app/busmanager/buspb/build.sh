#!/bin/sh

protoc -I ../buspb/ --go_out=plugins=grpc:../buspb busmanager.proto
