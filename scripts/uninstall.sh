#!/usr/bin/env bash
set -euo pipefail
BIN_DIR=${PREFIX:-/usr/local/bin}

if [ -f "$BIN_DIR/cmd" ]; then
    rm "$BIN_DIR/cmd"
    echo "Removed $BIN_DIR/cmd"
else
    echo "$BIN_DIR/cmd not found"
fi
