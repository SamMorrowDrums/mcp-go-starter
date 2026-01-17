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
		Description: "Generate a personalized greeting message with customizable style",
		Arguments: []*mcp.PromptArgument{
			{Name: "name", Description: "Name of the person to greet", Required: true},
			{Name: "style", Description: "The greeting style (formal, casual, enthusiastic)", Required: false},
		},
	}, greetPromptHandler)

	server.AddPrompt(&mcp.Prompt{
		Name:        "code_review",
		Title:       "Code Review",
		Description: "Request a code review with specific focus areas",
		Arguments: []*mcp.PromptArgument{
			{Name: "code", Description: "The code to review", Required: true},
			{Name: "language", Description: "Programming language", Required: true},
			{Name: "focus", Description: "What to focus on (security, performance, readability, all)", Required: false},
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
	language := req.Params.Arguments["language"]
	focus := req.Params.Arguments["focus"]
	if focus == "" {
		focus = "all"
	}

	focusInstructions := map[string]string{
		"security":    "Focus on security vulnerabilities and potential exploits.",
		"performance": "Focus on performance optimizations and efficiency issues.",
		"readability": "Focus on code clarity, naming, and maintainability.",
		"all":         "Provide a comprehensive review covering security, performance, and readability.",
	}

	instruction, ok := focusInstructions[focus]
	if !ok {
		instruction = focusInstructions["all"]
	}

	text := fmt.Sprintf("Please review the following %s code. %s\n\n```%s\n%s\n```", language, instruction, language, code)

	return &mcp.GetPromptResult{
		Messages: []*mcp.PromptMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: text},
			},
		},
	}, nil
}
