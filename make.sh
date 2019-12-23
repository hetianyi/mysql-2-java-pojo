#!/bin/sh

echo "building..."
go build -i -o bin/generator generator.go
echo "build success!"