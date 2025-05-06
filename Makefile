# Binary output path
BIN := build/jntool

# CLI source path (where main.go lives)
CMD_PATH := ./cmd/

# Default output directory for docs
DOCS_DIR := docs/commands


.PHONY: docs prepare build run test


# tidy and vendor dependencies
prepare:
	go mod tidy
	go mod vendor

# Build the CLI binary
build: prepare
	mkdir -p build
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o $(BIN) $(CMD_PATH)


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