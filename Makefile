.PHONY: build build-stdio build-http run-stdio run-http dev test clean fmt lint check

# Build all binaries
build: build-stdio build-http

build-stdio:
	go build -o bin/stdio ./cmd/stdio

build-http:
	go build -o bin/http ./cmd/http

# Run commands
run-stdio:
	go run ./cmd/stdio

run-http:
	go run ./cmd/http

# Development with live reload (requires air: go install github.com/air-verse/air@latest)
dev:
	air

# Test
test:
	go test -v ./...

# Format code
fmt:
	gofmt -w -s .
	goimports -w .

# Lint code
lint:
	golangci-lint run ./...

# Check all (for CI)
check: fmt lint test

# Clean
clean:
	rm -rf bin/ tmp/
	go clean

# Development
deps:
	go mod download
	go mod tidy

# Install linter (if needed)
install-tools:
	go install golang.org/x/tools/cmd/goimports@latest
	go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	go install github.com/air-verse/air@latest
