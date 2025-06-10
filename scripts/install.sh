#!/usr/bin/env bash
set -euo pipefail
BIN_DIR=${PREFIX:-/usr/local/bin}

# build for the local platform
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

echo "Building cmd for $GOOS/$GOARCH..."
tmp=$(mktemp -t cmd.XXXXXX)
go build -o "$tmp" ./
install -m 755 "$tmp" "$BIN_DIR/cmd"
rm "$tmp"

echo "Installed cmd to $BIN_DIR/cmd"
