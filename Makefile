# Build variables
VERSION := 0.2.0
COMMIT := $(shell git rev-parse --short HEAD 2>/dev/null || echo "unknown")
DATE := $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")
LDFLAGS := -X command/cmd.Version=$(VERSION) -X command/cmd.Commit=$(COMMIT) -X command/cmd.Date=$(DATE)

# Binary name
BINARY := command

.PHONY: all
all: build

.PHONY: build
build:
	go build -ldflags "$(LDFLAGS)" -o $(BINARY) .

.PHONY: install
install:
	go install -ldflags "$(LDFLAGS)"

.PHONY: test
test:
	go test -v ./...

.PHONY: test-coverage
test-coverage:
	go test -v -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out -o coverage.html

.PHONY: clean
clean:
	rm -f $(BINARY) coverage.out coverage.html

.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: lint
lint:
	golangci-lint run

.PHONY: version
version:
	@echo "Version: $(VERSION)"
	@echo "Commit: $(COMMIT)"
	@echo "Date: $(DATE)"

.PHONY: help
help:
	@echo "Available targets:"
	@echo "  make build         - Build the binary with version info"
	@echo "  make install       - Install the binary with version info"
	@echo "  make test          - Run tests"
	@echo "  make test-coverage - Run tests with coverage report"
	@echo "  make clean         - Remove build artifacts"
	@echo "  make fmt           - Format code"
	@echo "  make lint          - Run linter"
	@echo "  make version       - Show version information"
	@echo "  make help          - Show this help message"
