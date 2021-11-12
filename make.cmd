@echo off

set GOPROXY=https://goproxy.io
echo building...
go build -i -o bin/generator.exe generator.go
echo build success!

