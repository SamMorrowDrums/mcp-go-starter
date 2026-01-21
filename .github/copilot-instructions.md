# MCP Go Starter - Copilot Coding Agent Instructions

## Building and Testing

- **Download dependencies:**
  ```bash
  go mod download
  ```

- **Build the project:**
  ```bash
  go build -o mcp-go-starter ./cmd/mcp-go-starter
  ```

- **Run the server:**
  ```bash
  go run ./cmd/mcp-go-starter
  ```

- **Run tests:**
  ```bash
  go test ./...
  ```

- **Format code:**
  ```bash
  goimports -l -w <files you changed>
  ```

- **Vet code:**
  ```bash
  go vet ./...
  ```

## Code Conventions

- Follow Go naming conventions: use `ID` not `Id`, `API` not `Api`, `URL` not `Url` for acronyms.
- Wrap errors using `fmt.Errorf("context: %w", err)` to preserve error chains.
- Pass `context.Context` through call chains.

## Before Committing Checklist

1. ✅ Run `goimports -l -w` on modified files
2. ✅ Run `go vet ./...` and fix any errors
3. ✅ Run `go build ./...` to verify compilation
4. ✅ Run `go test ./...` to verify tests pass

