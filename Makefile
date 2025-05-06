# Binary output path
BIN := build/jntool

# CLI source path (where main.go lives)
CMD_PATH := ./cmd/

# Default output directory for docs
DOCS_DIR := docs/commands

# Determine install directory (GOBIN or GOPATH/bin)
INSTALL_DIR := $(shell go env GOBIN)
ifeq ($(INSTALL_DIR),)
INSTALL_DIR := $(shell go env GOPATH)/bin
endif

.PHONY: docs prepare build run test


# tidy and vendor dependencies
prepare:
	go mod tidy
	go mod vendor

# Build the CLI binary
build: prepare
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) $(CMD_PATH)

# Copy the built binary to $(INSTALL_DIR) for direct use
install: build
	mkdir -p $(INSTALL_DIR)
	cp $(BIN) $(INSTALL_DIR)/jntool
	@echo "Installed jntool to $(INSTALL_DIR)/jntool"

#   make run ARGS="helm values values.yaml -o json"
build-run: build
	@$(BIN) $(ARGS)

# Run the CLI with help (for debugging)
run:
	go run ./cmd/main.go --help

# Generate markdown docs for all commands
docs:
	go run ./cmd/main.go utils docs -output $(DOCS_DIR)

# Run all tests
test:
	go test ./... -v