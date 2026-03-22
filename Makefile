.PHONY: help build test test-verbose test-integration test-bench lint clean install examples run-example

help:
	@echo "BlurHash Go Library - Available commands:"
	@echo ""
	@echo "  make build          Build CLI tool"
	@echo "  make install        Build and install CLI tool"
	@echo "  make test           Run all tests"
	@echo "  make test-verbose   Run tests with verbose output"
	@echo "  make test-bench     Run benchmarks"
	@echo "  make test-coverage  Run tests with coverage report"
	@echo "  make test-fuzz      Run fuzz tests"
	@echo "  make test-integration Run only integration tests"
	@echo "  make lint           Run linter"
	@echo "  make fmt            Format code"
	@echo "  make clean          Clean build artifacts"
	@echo "  make examples       Run example code"
	@echo "  make cli-encode     Test CLI encode"
	@echo "  make cli-decode     Test CLI decode"
	@echo "  make cli-validate   Test CLI validate"
	@echo ""

build:
	@echo "Building CLI tool..."
	@go build -o cmd/blurhash/blurhash ./cmd/blurhash
	@echo "✓ Built cmd/blurhash/blurhash"

install: build
	@echo "Installing CLI tool..."
	@go install ./cmd/blurhash
	@echo "✓ Installed blurhash CLI"

test:
	@echo "Running tests..."
	@go test ./...

test-verbose:
	@echo "Running tests (verbose)..."
	@go test ./... -v

test-integration:
	@echo "Running integration tests..."
	@go test ./... -run Integration -v

test-bench:
	@echo "Running benchmarks..."
	@go test ./... -bench=. -benchmem -run=^$

test-coverage:
	@echo "Running tests with coverage..."
	@go test ./... -coverprofile=coverage.out
	@go tool cover -html=coverage.out -o coverage.html
	@echo "✓ Coverage report: coverage.html"

test-fuzz:
	@echo "Running fuzz tests..."
	@go test ./... -fuzz=Fuzz -fuzztime=10s

lint:
	@echo "Linting code..."
	@gofmt -l ./...
	@go vet ./...

fmt:
	@echo "Formatting code..."
	@gofmt -w ./...
	@echo "✓ Code formatted"

clean:
	@echo "Cleaning..."
	@rm -f cmd/blurhash/blurhash coverage.out coverage.html
	@go clean -cache
	@echo "✓ Cleaned"

examples:
	@echo "Running examples..."
	@go run examples/example.go

cli-encode: build
	@echo "Testing CLI encode..."
	@./cmd/blurhash/blurhash encode testdata/gradient.png

cli-decode: build
	@echo "Testing CLI decode..."
	@./cmd/blurhash/blurhash decode "LxH2cg2kwzX5l?WGjue:gLfkfQfj" -w 128 -h 128 -out /tmp/test.png
	@file /tmp/test.png 2>&1 || echo "File saved to /tmp/test.png"

cli-validate: build
	@echo "Testing CLI validate..."
	@./cmd/blurhash/blurhash validate "B~LrYI~c{H?b=::k" "" "invalid!"

# Development helpers
dev-watch:
	@echo "Watching for changes and running tests..."
	@while true; do \
		clear; \
		go test ./...; \
		echo "Waiting for changes..."; \
		find . -name "*.go" | entr -r true; \
	done

release:
	@echo "Creating release..."
	@git tag -a v1.0.0 -m "BlurHash Go library v1.0.0" 2>/dev/null || true
	@git push origin v1.0.0 2>/dev/null || echo "Note: Push manually with: git push origin v1.0.0"

.DEFAULT_GOAL := help
