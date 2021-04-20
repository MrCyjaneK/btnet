#!/bin/bash
echo "[DEPRECATED] - Do not use this build script... use dist/build.sh"
sleep 15
set -e
rm -rf build || true &>/dev/null
mkdir build
GITVERSION="+git"$(date +%Y%m%d%H%M)"."$(git log -n 1 | tr " " "\n" | head -2 | tail -1 | head -c 7)
#for pair in $(go tool dist list)
for pair in \
linux/386 \
linux/amd64 \
linux/arm \
linux/arm64 \
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
    if [[ "$p1" == "linux" ]];
    then
        case "$p2" in
        "386") arch="i386" ;;
        "amd64") arch="amd64" ;;
        "arm") arch="armhf" ;;
        "arm64") arch="aarch64" ;;
        esac        
        a=$(pwd)
        mkdir -p build/$p2
        cd dist/debian
        GOOS=$p1 GOARCH=$p2 BINNAME=btnet checkinstall --install=no \
        --pkgname="btnet" \
        --pkgversion=1.0.0"$GITVERSION" \
        --pkgarch="$arch" \
        --pkgrelease=1 \
        --pkgsource="git.mrcyjanek.net/mrcyjanek/btnet" \
        --pakdir="../$p2" \
        --maintainer="cyjan@mrcyjanek.net" \
        --provides="btnet" \
        -D \
        -y
        cd "$a"
    fi
    echo "OK!"
done