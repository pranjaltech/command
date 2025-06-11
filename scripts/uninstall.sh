#!/usr/bin/env bash
set -euo pipefail
BIN_DIR=${PREFIX:-/usr/local/bin}

# remove with sudo if directory is protected
if [ -f "$BIN_DIR/cmd" ]; then
    if [ ! -w "$BIN_DIR" ] || [ ! -w "$BIN_DIR/cmd" ]; then
        echo "Removing from $BIN_DIR requires elevated permissions" >&2
        if command -v sudo >/dev/null 2>&1; then
            sudo rm "$BIN_DIR/cmd"
        else
            echo "sudo not found. Re-run with PREFIX pointing to a writable directory." >&2
            exit 1
        fi
    else
        rm "$BIN_DIR/cmd"
    fi
    echo "Removed $BIN_DIR/cmd"
else
    echo "$BIN_DIR/cmd not found"
fi
