#!/bin/bash

VER=$(date "+%Y%m%d.%H%M")
VER_FILE=lib/version.go

echo "building version: ${VER}"

echo "package lib" > ${VER_FILE}
echo  >> ${VER_FILE}
echo "var Version = \"salesforcedoc/force-cli ${VER}\"" >> ${VER_FILE}
### build for macos-x64
env GOOS=darwin GOARCH=amd64 go build -o force-macos-x64 main.go
## build for windows-x64
env GOOS=windows GOARCH=amd64 go build -o force-windows-x64.exe main.go
### build for linux-x64
env GOOS=linux GOARCH=amd64 go build -o force-linux-x64 main.go
### install to macos bin path
#cp -pR force-macos-x64 /usr/local/Cellar/go/1.13/bin/force
cp -pR force-macos-x64 /opt/homebrew/Cellar/go/1.22.0/bin/force
###
force version