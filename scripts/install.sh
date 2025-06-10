#!/usr/bin/env bash
set -euo pipefail
BIN_DIR=${PREFIX:-/usr/local/bin}

# build for the local platform
GOOS=$(go env GOOS)
GOARCH=$(go env GOARCH)

echo "Building cmd for $GOOS/$GOARCH..."
go build -o cmd
install -m 755 cmd "$BIN_DIR/cmd"
rm cmd

echo "Installed cmd to $BIN_DIR/cmd"
