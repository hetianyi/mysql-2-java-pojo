#!/bin/sh
export GOPROXY=https://goproxy.io
echo "building..."
go build -i -o bin/generator generator.go
echo "build success!"