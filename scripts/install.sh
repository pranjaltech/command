#!/usr/bin/env bash
set -euo pipefail
BIN_DIR=${PREFIX:-/usr/local/bin}

# build for the local platform
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

echo "Building cmd for $GOOS/$GOARCH..."
tmp=$(mktemp -t cmd.XXXXXX)
go build -o "$tmp" ./
if [ ! -w "$BIN_DIR" ]; then
    echo "Installing to $BIN_DIR requires elevated permissions" >&2
    if command -v sudo >/dev/null 2>&1; then
        sudo install -m 755 "$tmp" "$BIN_DIR/cmd"
    else
        echo "sudo not found. Rerun with PREFIX set to a writable directory." >&2
        exit 1
    fi
else
    install -m 755 "$tmp" "$BIN_DIR/cmd"
fi
rm "$tmp"

echo "Installed cmd to $BIN_DIR/cmd"
