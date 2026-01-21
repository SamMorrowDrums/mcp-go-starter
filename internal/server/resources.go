package server

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// ItemData represents an item in our example data store.
type ItemData struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

// Example data for resources.
var itemsData = map[string]ItemData{
	"1": {ID: "1", Name: "Widget", Description: "A useful widget"},
	"2": {ID: "2", Name: "Gadget", Description: "A fancy gadget"},
	"3": {ID: "3", Name: "Gizmo", Description: "A mysterious gizmo"},
}

func registerResources(server *mcp.Server) {
	server.AddResource(&mcp.Resource{
		Name:        "About",
		Description: "Information about this MCP server",
		MIMEType:    "text/plain",
		URI:         "about://server",
	}, aboutResourceHandler)

	server.AddResource(&mcp.Resource{
		Name:        "Example Document",
		Description: "An example document resource",
		MIMEType:    "text/plain",
		URI:         "doc://example",
	}, exampleFileHandler)

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "Personalized Greeting",
		Description: "A personalized greeting for a specific person",
		MIMEType:    "text/plain",
		URITemplate: "greeting://{name}",
	}, greetingTemplateHandler)

	server.AddResourceTemplate(&mcp.ResourceTemplate{
		Name:        "Item Data",
		Description: "Data for a specific item by ID",
		MIMEType:    "application/json",
		URITemplate: "item://{id}",
	}, itemTemplateHandler)
}

func aboutResourceHandler(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "about://server",
				MIMEType: "text/plain",
				Text: `MCP Go Starter v1.0.0

This is a feature-complete MCP server demonstrating:
- Tools with annotations and structured output
- Resources (static and dynamic)
- Resource templates
- Prompts with completions
- Sampling, progress updates, and dynamic tool loading

For more information, visit: https://modelcontextprotocol.io`,
			},
		},
	}, nil
}

func exampleFileHandler(_ context.Context, _ *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      "doc://example",
				MIMEType: "text/plain",
				Text: `# Example Document

This is an example markdown document served as an MCP resource.

## Features

- **Bold text** and *italic text*
- Lists and formatting
- Code blocks

` + "```go\nhello := \"world\"\n```" + `

## Links

- [MCP Documentation](https://modelcontextprotocol.io)
- [Go SDK](https://github.com/modelcontextprotocol/go-sdk)`,
			},
		},
	}, nil
}

func greetingTemplateHandler(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	name := extractParam(req.Params.URI, "greeting://")

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "text/plain",
				Text:     fmt.Sprintf("Hello, %s! This greeting was generated just for you.", name),
			},
		},
	}, nil
}

func itemTemplateHandler(_ context.Context, req *mcp.ReadResourceRequest) (*mcp.ReadResourceResult, error) {
	id := extractParam(req.Params.URI, "item://")

	item, ok := itemsData[id]
	if !ok {
		return nil, fmt.Errorf("item not found: %s", id)
	}

	jsonBytes, _ := json.MarshalIndent(item, "", "  ")

	return &mcp.ReadResourceResult{
		Contents: []*mcp.ResourceContents{
			{
				URI:      req.Params.URI,
				MIMEType: "application/json",
				Text:     string(jsonBytes),
			},
		},
	}, nil
}
