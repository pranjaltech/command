#!/usr/bin/env bash
set -euo pipefail

export CGO_ENABLED=0

go install github.com/spf13/cobra-cli@latest

go install honnef.co/go/tools/cmd/staticcheck@latest

go install github.com/segmentio/golines@latest

go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
