#!/bin/bash
set -e
function ok {
    echo "OK"
}

root=$(dirname $0)
cd "$root"
vcode="+git"$(date +%Y%m%d%H%M)"."$(git log -n 1 | tr " " "\n" | head -2 | tail -1 | head -c 7)
echo "Building BTnet - version: $vcode";
rm -rf out || true
mkdir out
cd out
builddir=$(pwd)
mkdir bin

for addon in ../../helpers/*;
do
    cd "$builddir"
    name=$(basename $addon)
    BINNAME="btnet-$name"
    echo "/ Linux builds - daemon."
    echo -n -e "|- bin/$BINNAME-$name""_linux_386 "
    CGO_ENABLED=1 CC=i686-linux-gnu-gcc CXX=i686-linux-gnu-g++ GOOS=linux GOARCH=386 go build -o bin/"$BINNAME"_linux_386 $addon && ok
    echo -n -e "|- bin/$BINNAME-$name""_linux_amd64 "
    CGO_ENABLED=1 GOOS=linux GOARCH=amd64   go build -o bin/"$BINNAME"_linux_amd64 -tags gui ../../ && ok
    echo -n -e "|- bin/$BINNAME-$name""_linux_arm "
    CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -o bin/"$BINNAME"_linux_arm $addon && ok
    echo -n -e "\_ bin/$BINNAME-$name""_linux_arm64 "
    CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ GOOS=linux GOARCH=arm64 go build -o bin/"$BINNAME"_linux_arm64 $addon && ok
    echo "/ Windows builds - daemon"
    echo -n -e "|- bin/$BINNAME-$name""_windows_386.exe "
    CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -o bin/"$BINNAME"_windows_386.exe $addon && ok
    echo -n -e "|- bin/$BINNAME-$name""_windows_amd64.exe "
    CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o bin/"$BINNAME"_windows_amd64.exe $addon && ok
    for arch in arm64 arm amd64 386
    do
        case "$arch" in
        "386") ar="i386" ;;
        "amd64") ar="amd64" ;;
        "arm") ar="armhf" ;;
        "arm64") ar="aarch64" ;;
        esac        
        cp "$builddir/"../debian debian-deb-$name-$arch -r
        cd debian-deb-$name-$arch
        pwd
        pwd
        addon="../$addon"
        cat $addon/description-pak > description-pak
        GOOS="linux" GOARCH=$arch BINNAME="btnet-$name" checkinstall --install=no \
            --pkgname="btnet-$name" \
            --pkgversion=1.0.0"$vcode" \
            --pkgarch="$ar" \
            --pkgrelease=1 \
            --pkgsource="git.mrcyjanek.net/mrcyjanek/btnet" \
            --pakdir="../bin" \
            --maintainer="cyjan@mrcyjanek.net" \
            --provides="btnet-$name" \
            -D \
            -y
    done
done
cd $builddir

BINNAME="btnet"
echo "/ Linux builds - daemon."
echo -n -e "|- bin/$BINNAME""_linux_386 "
CGO_ENABLED=1 CC=i686-linux-gnu-gcc CXX=i686-linux-gnu-g++ GOOS=linux GOARCH=386 go build -o bin/"$BINNAME"_linux_386 ../../ && ok
echo -n -e "|- bin/$BINNAME""_linux_amd64 "
CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -o bin/jwstudy_linux_amd64 ../../ && ok
echo -n -e "|- bin/$BINNAME""_linux_arm "
CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc CXX=arm-linux-gnueabi-g++ GOOS=linux GOARCH=arm go build -o bin/"$BINNAME"_linux_arm ../../ && ok
echo -n -e "|_ bin/$BINNAME""_linux_arm64 "
CGO_ENABLED=1 CC=aarch64-linux-gnu-gcc CXX=aarch64-linux-gnu-g++ GOOS=linux GOARCH=arm64 go build -o bin/"$BINNAME"_linux_arm64 ../../ && ok
echo "/ Windows builds - daemon"
echo -n -e "|- bin/$BINNAME""_windows_386.exe "
CGO_ENABLED=1 CC=i686-w64-mingw32-gcc CXX=i686-w64-mingw32-g++ GOOS=windows GOARCH=386 go build -o bin/"$BINNAME"_windows_386.exe ../../ && ok
echo -n -e "|_ bin/$BINNAME""_windows_amd64.exe "
CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -o bin/"$BINNAME"_windows_amd64.exe ../../ && ok

echo "/ Packaging for debian"
for arch in arm64 arm amd64 386
do
    cd "$builddir"
    cp ../debian debian-deb-$arch -r
    cd debian-deb-$arch
    GOOS=$p1 GOARCH=$p2 BINNAME=btnet checkinstall --install=no \
        --pkgname="btnet" \
        --pkgversion=1.0.0"$vcode" \
        --pkgarch="$arch" \
        --pkgrelease=1 \
        --pkgsource="git.mrcyjanek.net/mrcyjanek/btnet" \
        --pakdir="../bin" \
        --maintainer="cyjan@mrcyjanek.net" \
        --provides="btnet" \
        -D \
        -y
done