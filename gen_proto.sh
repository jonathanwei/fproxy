#!/bin/bash

rm -f proto/*.pb.go
protoc --proto_path=. --go_out=plugins=grpc:. proto/*.proto
