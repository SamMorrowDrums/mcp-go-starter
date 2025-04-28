package example

import (
	"context"
	_ "embed"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

//go:embed resources/example.md
var exampleMarkdown string

func NewExampleResource() (mcp.Resource, server.ResourceHandlerFunc) {
	resource := mcp.NewResource(
		"docs://example",
		"Example Markdown Resource",
		mcp.WithResourceDescription("A sample markdown file served as an MCP resource."),
		mcp.WithMIMEType("text/markdown"),
	)
	handler := func(ctx context.Context, req mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
		return []mcp.ResourceContents{
			mcp.TextResourceContents{
				URI:      "docs://example",
				MIMEType: "text/markdown",
				Text:     exampleMarkdown,
			},
		}, nil
	}
	return resource, handler
}

func RegisterResources(s *server.MCPServer) {
	resource, handler := NewExampleResource()
	s.AddResource(resource, handler)
}
