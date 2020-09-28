#!/bin/sh

protoc -I ../apppb/ --go_out=plugins=grpc:../apppb appmanager.proto
