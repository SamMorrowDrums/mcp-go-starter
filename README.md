# MCP Go Starter

[![CI](https://github.com/SamMorrowDrums/mcp-go-starter/actions/workflows/ci.yml/badge.svg)](https://github.com/SamMorrowDrums/mcp-go-starter/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/SamMorrowDrums/mcp-go-starter)](https://goreportcard.com/report/github.com/SamMorrowDrums/mcp-go-starter)
[![Go Version](https://img.shields.io/badge/Go-1.26.8%2B-00ADD8?logo=go&logoColor=white)](https://go.dev/)
[![License: MIT](https://img.shields.io/badge/License-MIT-yellow.svg)](https://opensource.org/licenses/MIT)
[![MCP](https://img.shields.io/badge/MCP-Model%20Context%20Protocol-purple)](https://modelcontextprotocol.io/)

A feature-complete Model Context Protocol (MCP) server template in Go using the official go-sdk. This starter demonstrates all major MCP features with clean, idiomatic Go code.

## 📚 Documentation

- [Model Context Protocol](https://modelcontextprotocol.io/)
- [Go SDK](https://github.com/modelcontextprotocol/go-sdk)
- [Building MCP Servers](https://modelcontextprotocol.io/docs/develop/build-server)

## ✨ Features

| Category | Feature | Description |
|----------|---------|-------------|
| **Tools** | `hello` | Basic tool with annotations |
| | `get_weather` | Tool returning structured data |
| | `ask_llm` | Tool that invokes LLM sampling |
| | `long_task` | Tool with 5-second progress updates |
| | `load_bonus_tool` | Dynamically loads a new tool |
| **Resources** | `info://about` | Static informational resource |
| | `file://example.md` | File-based markdown resource |
| **Templates** | `greeting://{name}` | Personalized greeting |
| | `data://items/{id}` | Data lookup by ID |
| **Prompts** | `greet` | Greeting in various styles |
| | `code_review` | Code review with focus areas |

MCP sampling (`ask_llm`) is deprecated as of protocol version 2026-07-28
([SEP-2577](https://modelcontextprotocol.io/seps/2577-deprecate-roots-sampling-and-logging)).
This workshop keeps the example during the supported deprecation window; new
production servers should call LLM provider APIs directly.

## 🚀 Quick Start

### Prerequisites

- [Go 1.26.8+](https://go.dev/dl/)
- (Optional) [air](https://github.com/air-verse/air) for live reload
- (Optional) [golangci-lint](https://golangci-lint.run/welcome/install/) for linting

Air and goimports are pinned in `tools/go.mod`, including their transitive
dependencies, and updated by Dependabot alongside the server modules.
Install them with `go -C tools install tool`.

### Installation

```bash
# Clone the repository
git clone https://github.com/SamMorrowDrums/mcp-go-starter.git
cd mcp-go-starter

# Download dependencies
go mod download
```

### Running the Server

**stdio transport** (for local development):
```bash
go run ./cmd/stdio
# Or: make run-stdio
```

**HTTP transport** (for remote/web deployment):
```bash
go run ./cmd/http
# Or: make run-http
# Server runs on http://localhost:3000
```

The HTTP endpoint requires JSON requests and rejects cross-origin browser
requests. Native MCP clients without an `Origin` header and same-origin browser
clients remain supported.

### Building Binaries

```bash
make build
# Creates bin/stdio and bin/http
```

## 🔧 VS Code Integration

This project includes VS Code configuration for seamless development:

1. Open the project in VS Code
2. The MCP configuration is in `.vscode/mcp.json`
3. Build with `Ctrl+Shift+B` (or `Cmd+Shift+B` on Mac)
4. Test the server using VS Code's MCP tools

### Using DevContainers

1. Install the [Dev Containers extension](https://marketplace.visualstudio.com/items?itemName=ms-vscode-remote.remote-containers)
2. Open command palette: "Dev Containers: Reopen in Container"
3. Everything is pre-configured and ready to use!

## 📁 Project Structure

```
.
├── cmd/
│   ├── stdio/
│   │   └── main.go        # stdio transport entrypoint
│   └── http/
│       └── main.go        # HTTP transport entrypoint
├── internal/
│   └── server/
│       ├── server.go      # Server orchestration
│       ├── tools.go       # Tool definitions (hello, get_weather, etc.)
│       ├── resources.go   # Resource and template definitions
│       └── prompts.go     # Prompt definitions
├── .vscode/
│   ├── mcp.json           # MCP server configuration
│   ├── tasks.json         # Build/run tasks
│   └── extensions.json
├── .devcontainer/
│   └── devcontainer.json
├── .air.toml              # Live reload configuration
├── .golangci.yml          # Linter configuration
├── go.mod
├── Makefile
└── README.md
```

## 🛠️ Development

```bash
# Development with live reload (recommended)
make dev
# Requires air: go install github.com/air-verse/air@latest

# Run without live reload
make run-stdio

# Run tests
make test

# Format code
make fmt

# Lint code
make lint

# Install all dev tools
make install-tools

# Clean build artifacts
make clean
```

### Live Reload

Install [air](https://github.com/air-verse/air) for automatic rebuilds:
```bash
go install github.com/air-verse/air@latest
make dev
```
Changes to any `.go` file will automatically rebuild and restart the server.

## 🔍 MCP Inspector

The [MCP Inspector](https://modelcontextprotocol.io/docs/tools/inspector) is an essential development tool for testing and debugging MCP servers.

### Running Inspector

```bash
npx @modelcontextprotocol/inspector -- go run ./cmd/stdio/main.go
```

### What Inspector Provides

- **Tools Tab**: List and invoke all registered tools with parameters
- **Resources Tab**: Browse and read resources and templates
- **Prompts Tab**: View and test prompt templates
- **Logs Tab**: See JSON-RPC messages between client and server
- **Schema Validation**: Verify tool input/output schemas

### Debugging Tips

1. Start Inspector before connecting your IDE/client
2. Use the "Logs" tab to see exact request/response payloads
3. Test tool annotations (ToolAnnotations) are exposed correctly
4. Verify progress notifications appear for `long_task`
5. Check that sampling works with `ask_llm`

## 📖 Feature Examples

### Tool with Annotations

```go
mcp.AddTool(server, &mcp.Tool{
    Name:        "hello",
    Title:       "Say Hello",
    Description: "A friendly greeting tool",
    Annotations: &mcp.ToolAnnotations{
        ReadOnlyHint: ptr(true),
    },
}, helloHandler)

func helloHandler(ctx context.Context, req *mcp.CallToolRequest, input helloInput) (*mcp.CallToolResult, any, error) {
    return &mcp.CallToolResult{
        Content: []mcp.Content{
            &mcp.TextContent{Text: fmt.Sprintf("Hello, %s!", input.Name)},
        },
    }, nil, nil
}
```

### Resource Template

```go
server.AddResourceTemplate(&mcp.ResourceTemplate{
    Name:        "Personalized Greeting",
    URITemplate: "greeting://{name}",
    MIMEType:    "text/plain",
}, greetingHandler)
```

### Tool with Progress Updates

```go
func longTaskHandler(ctx context.Context, req *mcp.CallToolRequest, input longTaskInput) (*mcp.CallToolResult, any, error) {
    progressToken := req.Params.GetProgressToken()
    
    for i := 0; i < 5; i++ {
        req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
            ProgressToken: progressToken,
            Progress:      float64(i) / 5.0,
            Total:         1.0,
        })
        time.Sleep(time.Second)
    }
    
    return &mcp.CallToolResult{
        Content: []mcp.Content{&mcp.TextContent{Text: "Done!"}},
    }, nil, nil
}
```

### Tool with Sampling

```go
result, err := req.Session.CreateMessage(ctx, &mcp.CreateMessageParams{
    Messages: []*mcp.SamplingMessage{
        {Role: "user", Content: &mcp.TextContent{Text: prompt}},
    },
    MaxTokens: 100,
})
```

## 🔐 Environment Variables

Copy `.env.example` to `.env` and configure:

```bash
cp .env.example .env
```

| Variable | Description | Default |
|----------|-------------|---------|
| `PORT` | HTTP server port | `3000` |

## 🤝 Contributing

Contributions welcome! Please ensure your changes maintain feature parity with other language starters.

## 📄 License

MIT License - see [LICENSE](LICENSE) for details.
