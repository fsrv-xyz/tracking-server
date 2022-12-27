#!/usr/bin/env bash

which podman &>/dev/null && UTIL="podman"
which docker &>/dev/null && UTIL="docker"

${UTIL} run -v "${PWD}":/project -it -e DEBIAN_FRONTEND=noninteractive golang:latest bash -c '
apt-get update && apt-get install -y protobuf-compiler || exit 1
go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
cd /project && protoc --go_out=./ --go_opt=paths=source_relative --go-grpc_out=./ --go-grpc_opt=paths=source_relative pkg/proto/ingest.proto
'
