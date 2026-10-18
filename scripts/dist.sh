#!/bin/sh
set -eu

ROOT=$(CDPATH= cd -- "$(dirname -- "$0")/.." && pwd)
cd "$ROOT"

VERSION=${VERSION:-dev}

case "$VERSION" in
	dev)
		;;
	v[0-9]*)
		case "$VERSION" in
			*[!A-Za-z0-9._+-]*)
				printf 'error: VERSION must be dev or match v[0-9]* using only letters, digits, ".", "_", "+", and "-"\n' >&2
				exit 1
				;;
		esac
		;;
	*)
		printf 'error: VERSION must be dev or match v[0-9]*\n' >&2
		exit 1
		;;
esac

if ! command -v zip >/dev/null 2>&1; then
	printf 'error: zip is required to package Windows release assets\n' >&2
	exit 1
fi

if command -v sha256sum >/dev/null 2>&1; then
	checksum_tool=sha256sum
elif command -v shasum >/dev/null 2>&1; then
	checksum_tool=shasum
else
	printf 'error: sha256sum or shasum is required to write dist/SHA256SUMS\n' >&2
	exit 1
fi

if ! command -v tar >/dev/null 2>&1; then
	printf 'error: tar is required to package release assets\n' >&2
	exit 1
fi

tar_version=$(tar --version 2>/dev/null || true)
case "$tar_version" in
	*"GNU tar"*)
		tar_owner_flags='--owner=0 --group=0'
		;;
	*)
		tar_owner_flags='--uid 0 --gid 0'
		;;
esac

rm -rf dist
mkdir dist

build() {
	goos=$1
	goarch=$2
	extension=
	if [ "$goos" = windows ]; then
		extension=.exe
	fi
	output="dist/mimamori-${goos}-${goarch}${extension}"
	CGO_ENABLED=0 GOOS="$goos" GOARCH="$goarch" GOARM=6 \
		go build -trimpath -buildvcs=false -ldflags "-s -w -X main.version=$VERSION" -o "$output" .

	package_dir="dist/.package-${goos}-${goarch}"
	mkdir "$package_dir"
	if [ "$goos" = windows ]; then
		archive="mimamori-${VERSION}-windows-${goarch}.zip"
		cp "$output" "$package_dir/mimamori.exe"
		cp LICENSE README.md "$package_dir/"
		if ! (cd "$package_dir" && zip -q -X "$ROOT/dist/$archive" mimamori.exe LICENSE README.md); then
			rm -rf "$package_dir"
			printf 'error: could not create %s\n' "$archive" >&2
			exit 1
		fi
	else
		archive="mimamori-${VERSION}-${goos}-${goarch}.tar.gz"
		cp "$output" "$package_dir/mimamori"
		chmod 0755 "$package_dir/mimamori"
		cp LICENSE README.md "$package_dir/"
		if ! (cd "$package_dir" && tar $tar_owner_flags -czf "$ROOT/dist/$archive" mimamori LICENSE README.md); then
			rm -rf "$package_dir"
			printf 'error: could not create %s\n' "$archive" >&2
			exit 1
		fi
	fi
	rm -rf "$package_dir" "$output"
}

build windows amd64
build windows arm64
build windows 386
build linux amd64
build linux arm64
build linux arm
build darwin amd64
build darwin arm64

if [ "$checksum_tool" = sha256sum ]; then
	(cd dist && sha256sum *.zip *.tar.gz) > dist/SHA256SUMS
else
	(cd dist && shasum -a 256 *.zip *.tar.gz) > dist/SHA256SUMS
fi
