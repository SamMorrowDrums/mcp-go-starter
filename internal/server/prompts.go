package server

import (
	"context"
	"fmt"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func registerPrompts(server *mcp.Server) {
	server.AddPrompt(&mcp.Prompt{
		Name:        "greet",
		Title:       "Greeting Prompt",
		Description: "Generate a greeting message",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "name",
				Title:       "Name",
				Description: "Name of the person to greet",
				Required:    true,
			},
			{
				Name:        "style",
				Title:       "Style",
				Description: "Greeting style (formal/casual)",
			},
		},
	}, greetPromptHandler)

	server.AddPrompt(&mcp.Prompt{
		Name:        "code_review",
		Title:       "Code Review",
		Description: "Review code for potential improvements",
		Arguments: []*mcp.PromptArgument{
			{
				Name:        "code",
				Title:       "Code",
				Description: "The code to review",
				Required:    true,
			},
		},
	}, codeReviewPromptHandler)
}

func greetPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	name := req.Params.Arguments["name"]
	style := req.Params.Arguments["style"]
	if style == "" {
		style = "casual"
	}

	styles := map[string]string{
		"formal":       fmt.Sprintf("Please compose a formal, professional greeting for %s.", name),
		"casual":       fmt.Sprintf("Write a casual, friendly hello to %s.", name),
		"enthusiastic": fmt.Sprintf("Create an excited, enthusiastic greeting for %s!", name),
	}

	text, ok := styles[style]
	if !ok {
		text = styles["casual"]
	}

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}

func codeReviewPromptHandler(_ context.Context, req *mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
	code := req.Params.Arguments["code"]

	text := fmt.Sprintf("Please review the following code:\n\n```\n%s\n```", code)

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}
