#!/bin/bash
set -e
rm -rf build || true &>/dev/null
mkdir build
#for pair in $(go tool dist list)
for pair in \
freebsd/386 \
freebsd/amd64 \
freebsd/arm \
freebsd/arm64 \
linux/386 \
linux/amd64 \
linux/arm \
linux/arm64 \
linux/mips \
linux/mips64 \
linux/mips64le \
linux/mipsle \
linux/ppc64 \
linux/ppc64le \
linux/riscv64 \
linux/s390x \
netbsd/386 \
netbsd/amd64 \
netbsd/arm \
netbsd/arm64 \
openbsd/386 \
openbsd/amd64 \
openbsd/arm \
openbsd/arm64 \
windows/386 \
windows/amd64 \
windows/arm
do
    pair=$(echo $pair | tr '/' ' ')
    p1=$(echo "$pair" | awk '{print $1}')
    p2=$(echo "$pair" | awk '{print $2}')
    echo -n -e "[$p1] [$p2] - Building..."
    end=""
    if [[ "$p1" == "windows" ]];
    then
        end=".exe"
    fi
    GOOS=$p1 GOARCH=$p2 go build -o build/btnet_"$p1"_"$p2""$end"
    echo "OK!"
done