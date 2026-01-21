// Package server defines and configures the MCP server with all features.
//
// This package creates a feature-complete MCP server demonstrating:
//   - Tools with annotations, sampling, progress, and dynamic loading
//   - Resources (static and dynamic)
//   - Resource Templates
//   - Prompts with completions
//
// Documentation: https://modelcontextprotocol.io/
package server

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Global reference to server for dynamic tool loading.
var globalServer *mcp.Server

// ServerInstructions provides guidance for AI assistants on how to use this server.
const ServerInstructions = `# MCP Go Starter Server

A demonstration MCP server showcasing Go SDK capabilities.

## Available Tools

### Greeting & Demos
- **hello**: Simple greeting - use to test connectivity
- **get_weather**: Returns simulated weather data
- **long_task**: Demonstrates progress reporting (takes ~5 seconds)

### LLM Interaction
- **ask_llm**: Invoke LLM sampling to ask questions (requires client support)

### Dynamic Features
- **load_bonus_tool**: Dynamically adds a calculator tool at runtime
- **bonus_calculator**: Available after calling load_bonus_tool

### Elicitation (User Input)
- **confirm_action**: Demonstrates schema elicitation - requests user confirmation
- **get_feedback**: Demonstrates URL elicitation - opens feedback form in browser

## Available Resources

- **about://server**: Server information
- **doc://example**: Sample document
- **greeting://{name}**: Personalized greeting template
- **item://{id}**: Item data by ID

## Available Prompts

- **greet**: Generates a personalized greeting
- **code_review**: Structured code review prompt

## Recommended Workflows

1. **Testing Connection**: Call hello with your name to verify the server is responding
2. **Weather Demo**: Call get_weather with a location to see structured output
3. **Progress Demo**: Call long_task to see progress notifications
4. **Dynamic Loading**: Call load_bonus_tool, then refresh tools to see bonus_calculator
5. **Elicitation Demo**: Call confirm_action to see user confirmation flow
6. **URL Elicitation**: Call get_feedback to open a feedback form

## Tool Annotations

All tools include annotations indicating:
- Whether they modify state (ReadOnlyHint)
- If they're safe to retry (IdempotentHint)
- Whether they access external systems (OpenWorldHint)

Use these hints to make informed decisions about tool usage.`

// NewServer creates and configures the MCP server with all features.
func NewServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-go-starter",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			HasTools:     true,
			HasResources: true,
			HasPrompts:   true,
			Instructions: ServerInstructions,
		},
	)

	registerTools(server)
	registerResources(server)
	registerPrompts(server)

	return server
}

// SetGlobalServer sets the global server reference for dynamic tool loading.
func SetGlobalServer(s *mcp.Server) {
	globalServer = s
}

// extractParam extracts a parameter from a URI by removing the prefix.
func extractParam(uri, prefix string) string {
	if len(uri) > len(prefix) {
		return uri[len(prefix):]
	}
	return ""
}
