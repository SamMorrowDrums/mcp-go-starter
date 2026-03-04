// Package server defines and configures the MCP server with all features.
//
// ARCHITECTURE OVERVIEW:
// An MCP server exposes three core primitives to AI clients:
//   - Tools:     Functions the AI can call (like API endpoints)
//   - Resources: Read-only data the AI can access (like GET endpoints)
//   - Prompts:   Reusable prompt templates for common workflows
//
// Each primitive is registered with the server and advertised to clients
// during capability negotiation. See: https://modelcontextprotocol.io/
package server

import "github.com/modelcontextprotocol/go-sdk/mcp"

// Global reference to server for dynamic tool loading.
var globalServer *mcp.Server

// ServerInstructions provides guidance for AI assistants on how to use this server.
const ServerInstructions = "# MCP Go Starter Server\n\n" +
	"A demonstration MCP server showcasing Go SDK capabilities.\n\n" +
	"## Recommended Workflows\n\n" +
	"1. **Test connectivity** → Call `hello` to verify the server responds\n" +
	"2. **Structured output** → Call `get_weather` to see typed response data\n" +
	"3. **Progress reporting** → Call `long_task` to observe real-time progress notifications\n" +
	"4. **Dynamic tools** → Call `load_bonus_tool`, then re-list tools to see `bonus_calculator` appear\n" +
	"5. **LLM sampling** → Call `ask_llm` to have the server request a completion from the client\n" +
	"6. **Elicitation** → Call `confirm_action` (form-based) or `get_feedback` (URL-based) to request user input\n\n" +
	"## Multi-Tool Flows\n\n" +
	"- **Full demo**: `hello` → `get_weather` → `long_task` → `load_bonus_tool` → `bonus_calculator`\n" +
	"- **Dynamic loading**: `load_bonus_tool` triggers a `tools/list_changed` notification — refresh your tool list to see `bonus_calculator`\n" +
	"- **User interaction**: `confirm_action` demonstrates schema elicitation, `get_feedback` demonstrates URL elicitation\n\n" +
	"## Notes\n\n" +
	"- All tools include annotations (readOnlyHint, idempotentHint, openWorldHint) to guide safe usage\n" +
	"- Resources and prompts are available for context and templating — use `resources/list` and `prompts/list` to discover them"

// NewServer creates and configures the MCP server with all features.
//
// CAPABILITIES tell the client what this server supports. During the MCP
// handshake, the client reads these to know which features are available.
func NewServer() *mcp.Server {
	server := mcp.NewServer(
		&mcp.Implementation{
			Name:    "mcp-go-starter",
			Version: "1.0.0",
		},
		&mcp.ServerOptions{
			Instructions: ServerInstructions,
			Capabilities: &mcp.ServerCapabilities{
				Experimental: map[string]any{},
				Resources: &mcp.ResourceCapabilities{
					// ListChanged: false — our resources are static, so we never
					// need to notify clients that the resource list has changed.
					ListChanged: false,
					Subscribe:   false,
				},
				Tools: &mcp.ToolCapabilities{
					// ListChanged: true — because load_bonus_tool adds tools
					// dynamically at runtime. When a tool is added, the server
					// sends a tools/list_changed notification so clients refresh.
					ListChanged: true,
				},
			},
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
