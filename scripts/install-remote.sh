#!/usr/bin/env bash
set -euo pipefail

REPO="pranjaltech/command"
BINARY_NAME="cmd"
INSTALL_DIR="${PREFIX:-/usr/local/bin}"

cleanup() {
    [ -n "${TMPDIR:-}" ] && rm -rf "$TMPDIR"
}
trap cleanup EXIT

# Detect OS
OS="$(uname -s | tr '[:upper:]' '[:lower:]')"
case "$OS" in
    linux|darwin) ;;
    *) echo "Unsupported OS: $OS" >&2; exit 1 ;;
esac

# Detect architecture
ARCH="$(uname -m)"
case "$ARCH" in
    x86_64)          ARCH="amd64" ;;
    aarch64|arm64)   ARCH="arm64" ;;
    *) echo "Unsupported architecture: $ARCH" >&2; exit 1 ;;
esac

# Determine version
if [ -z "${VERSION:-}" ]; then
    echo "Fetching latest release..."
    VERSION="$(curl -sSL "https://api.github.com/repos/${REPO}/releases/latest" \
        | grep '"tag_name"' | sed -E 's/.*"v?([^"]+)".*/\1/')" || true
    if [ -z "$VERSION" ]; then
        echo "Could not determine latest version. Set VERSION=x.y.z and retry." >&2
        exit 1
    fi
fi

URL="https://github.com/${REPO}/releases/download/v${VERSION}/${BINARY_NAME}_${VERSION}_${OS}_${ARCH}.tar.gz"
echo "Downloading ${BINARY_NAME} v${VERSION} for ${OS}/${ARCH}..."

TMPDIR="$(mktemp -d)"
curl -sSL -o "${TMPDIR}/cmd.tar.gz" "$URL"
tar -xzf "${TMPDIR}/cmd.tar.gz" -C "$TMPDIR"

mkdir -p "$INSTALL_DIR"
if [ -w "$INSTALL_DIR" ]; then
    install -m 755 "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
else
    echo "Installing to ${INSTALL_DIR} requires elevated permissions" >&2
    if command -v sudo >/dev/null 2>&1; then
        sudo install -m 755 "${TMPDIR}/${BINARY_NAME}" "${INSTALL_DIR}/${BINARY_NAME}"
    else
        echo "sudo not found. Set PREFIX to a writable directory and retry." >&2
        exit 1
    fi
fi

echo "Installed ${BINARY_NAME} v${VERSION} to ${INSTALL_DIR}/${BINARY_NAME}"
"${INSTALL_DIR}/${BINARY_NAME}" --version
