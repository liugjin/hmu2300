#!/bin/sh

protoc -I ../portpb/ --go_out=plugins=grpc:../portpb portmanager.proto
