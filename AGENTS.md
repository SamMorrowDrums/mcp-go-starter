# AGENTS.md

This file provides context for AI coding agents working in this repository.

## Quick Reference

| Task | Command |
|------|---------|
| Install deps | `go mod download` |
| Build | `go build ./...` |
| Test | `go test -v ./...` |
| Vet | `go vet ./...` |
| Format | `gofmt -w .` |
| Format check | `gofmt -l .` (empty = OK) |
| Run (stdio) | `go run ./cmd/stdio` |
| Run (HTTP) | `go run ./cmd/http` |

## Project Overview

**MCP Go Starter** is a feature-complete Model Context Protocol (MCP) server template in Go using the official go-sdk. It demonstrates all major MCP features including tools, resources, resource templates, prompts, sampling, progress updates, and dynamic tool loading.

**Purpose**: Workshop starter template for learning MCP server development.

## Technology Stack

- **Runtime**: Go 1.23
- **MCP SDK**: `github.com/modelcontextprotocol/go-sdk`
- **HTTP Server**: net/http (stdlib)
- **Formatter**: gofmt (built-in)

## Project Structure

```
.
├── go.mod                      # Module definition
├── go.sum                      # Dependency checksums
├── Makefile                    # Build/run commands
├── .golangci.yml               # Linter configuration
├── cmd/
│   ├── stdio/
│   │   └── main.go             # stdio transport entrypoint
│   └── http/
│       └── main.go             # HTTP transport entrypoint
├── internal/
│   └── server/
│       └── server.go           # Main server (tools, resources, prompts)
├── .vscode/
│   ├── mcp.json                # MCP server configuration
│   ├── tasks.json              # Build/run tasks
│   ├── launch.json             # Debug configurations
│   └── extensions.json
└── .devcontainer/
    └── devcontainer.json
```

## Build & Run Commands

```bash
# Download dependencies
go mod download && go mod tidy

# Build all binaries
go build ./...

# Run server (stdio transport)
go run ./cmd/stdio

# Run server (HTTP transport)
go run ./cmd/http
# With custom port: PORT=8080 go run ./cmd/http
```

## Linting & Formatting

```bash
# Format code (required before commit)
gofmt -w .

# Check formatting (CI uses this - empty output = OK)
gofmt -l .

# Vet code for suspicious constructs
go vet ./...
```

## Testing

```bash
go test -v ./...
```

## Key Files to Modify

- **Add/modify tools**: `internal/server/server.go` → `registerTools()` function
- **Add/modify resources**: `internal/server/server.go` → `registerResources()` function
- **Add/modify prompts**: `internal/server/server.go` → `registerPrompts()` function
- **Server configuration**: `internal/server/server.go` → `NewServer()` function
- **HTTP config**: `cmd/http/main.go`
- **Module name**: `go.mod`

## MCP Features Implemented

| Feature | Location | Description |
|---------|----------|-------------|
| `hello` tool | `server.go` | Basic tool with annotations |
| `get_weather` tool | `server.go` | Structured JSON output |
| `ask_llm` tool | `server.go` | Sampling/LLM invocation |
| `long_task` tool | `server.go` | Progress updates |
| `load_bonus_tool` | `server.go` | Dynamic tool loading |
| Resources | `server.go` | Static `info://about`, `file://example.md` |
| Templates | `server.go` | `greeting://{name}`, `data://items/{id}` |
| Prompts | `server.go` | `greet`, `code_review` with arguments |

## Environment Variables

- `PORT` - HTTP server port (default: 3000)

## Conventions

- Use `mcp.AddTool()` to register tools with typed input structs
- Use jsonschema struct tags for input validation
- Follow standard Go project layout (`cmd/`, `internal/`)
- Run `gofmt -w .` before committing
- Run `go vet ./...` before PRs

## Code Quality Tools

- **gofmt**: Standard Go formatter (built-in)

## Tool Input Pattern

```go
type myToolInput struct {
    Name string `json:"name" jsonschema:"description=The name parameter"`
}

mcp.AddTool(server, &mcp.Tool{
    Name:        "my_tool",
    Description: "Tool description",
}, func(ctx context.Context, req *mcp.CallToolRequest, input myToolInput) (*mcp.CallToolResult, any, error) {
    // Handler logic
})
```

## Documentation Links

- [MCP Specification](https://modelcontextprotocol.io/)
- [Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Building Servers](https://modelcontextprotocol.io/docs/develop/build-server)
