#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"

VERSION=${VERSION:-dev}
mkdir -p dist

build() {
	goos=$1
	goarch=$2
	extension=
	if [ "$goos" = windows ]; then
		extension=.exe
	fi
	output="dist/mimamori-${goos}-${goarch}${extension}"
	CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" \
		go build -trimpath -ldflags "-s -w -X main.version=$VERSION" -o "$output" .
}

build windows amd64
build windows arm64
build windows 386
build linux amd64
build linux arm64
build linux arm
build darwin amd64
build darwin arm64
