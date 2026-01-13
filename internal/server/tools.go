package server

import (
	"context"
	"encoding/json"
	"fmt"
	"math/rand"
	"time"

	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// Weather represents weather data returned by the get_weather tool.
type Weather struct {
	Location    string `json:"location"`
	Temperature int    `json:"temperature"`
	Unit        string `json:"unit"`
	Conditions  string `json:"conditions"`
	Humidity    int    `json:"humidity"`
}

// Track if bonus tool is loaded.
var bonusToolLoaded = false

// Tool input types

type helloInput struct {
	Name string `json:"name" jsonschema:"The name to greet"`
}

type weatherInput struct {
	Location string `json:"location" jsonschema:"City name or coordinates"`
}

type askLLMInput struct {
	Prompt    string `json:"prompt" jsonschema:"The question or prompt for the LLM"`
	MaxTokens int    `json:"maxTokens,omitempty" jsonschema:"Maximum tokens in response"`
}

type longTaskInput struct {
	TaskName string `json:"taskName" jsonschema:"Name for this task"`
}

type calculatorInput struct {
	A         float64 `json:"a" jsonschema:"First number"`
	B         float64 `json:"b" jsonschema:"Second number"`
	Operation string  `json:"operation" jsonschema:"enum=add,enum=subtract,enum=multiply,enum=divide"`
}

type confirmActionInput struct {
	Action string `json:"action" jsonschema:"The action to confirm with the user"`
}

type feedbackInput struct {
	Topic string `json:"topic,omitempty" jsonschema:"Optional topic for the feedback"`
}

// =============================================================================
// Tool Annotations - Every tool SHOULD have annotations for AI assistants
//
// WHY ANNOTATIONS MATTER:
// Annotations enable MCP client applications to understand the risk level of
// tool calls. Clients can use these hints to implement safety policies, such as:
//   - Prompting users for confirmation before executing destructive operations
//   - Auto-approving read-only tools while requiring approval for writes
//   - Warning users when tools access external systems (openWorldHint)
//   - Optimizing retry logic for idempotent operations
//
// ANNOTATION FIELDS:
// - ReadOnlyHint: Tool only reads data, doesn't modify state (bool)
// - DestructiveHint: Tool can permanently delete or modify data (*bool)
// - IdempotentHint: Repeated calls with same args have same effect (bool)
// - OpenWorldHint: Tool accesses external systems (web, APIs, etc.) (*bool)
// =============================================================================

// Helper to create a bool pointer for annotation fields that use *bool
func boolPtr(b bool) *bool {
	return &b
}

func registerTools(server *mcp.Server) {
	mcp.AddTool(server, &mcp.Tool{
		Name:        "hello",
		Title:       "Say Hello",
		Description: "A friendly greeting tool that says hello to someone",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Say Hello",
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(false),
		},
		Icons: []mcp.Icon{
			{
				Source:   WAVING_HAND_ICON,
				MIMEType: "image/png",
				Sizes:    []string{"256x256"},
			},
		},
	}, helloHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_weather",
		Title:       "Get Weather",
		Description: "Get current weather for a location (simulated)",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Get Weather",
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  false,          // Simulated - results vary
			OpenWorldHint:   boolPtr(false), // Not real external call
		},
		Icons: []mcp.Icon{
			{
				Source:   SUN_BEHIND_CLOUD_ICON,
				MIMEType: "image/png",
				Sizes:    []string{"256x256"},
			},
		},
	}, weatherHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "ask_llm",
		Title:       "Ask LLM",
		Description: "Ask the connected LLM a question using sampling",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Ask LLM",
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  false, // LLM responses vary
			OpenWorldHint:   boolPtr(false),
		},
		Icons: []mcp.Icon{
			{
				Source:   ROBOT_ICON,
				MIMEType: "image/png",
				Sizes:    []string{"256x256"},
			},
		},
	}, askLLMHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "long_task",
		Title:       "Long Running Task",
		Description: "A task that takes 5 seconds and reports progress along the way",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Long Running Task",
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true,
			OpenWorldHint:   boolPtr(false),
		},
		Icons: []mcp.Icon{
			{
				Source:   HOURGLASS_ICON,
				MIMEType: "image/png",
				Sizes:    []string{"256x256"},
			},
		},
	}, longTaskHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "load_bonus_tool",
		Title:       "Load Bonus Tool",
		Description: "Dynamically loads a bonus tool that wasn't available at startup",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Load Bonus Tool",
			ReadOnlyHint:    false, // Modifies server state
			DestructiveHint: boolPtr(false),
			IdempotentHint:  true, // Safe to call multiple times
			OpenWorldHint:   boolPtr(false),
		},
		Icons: []mcp.Icon{
			{
				Source:   PACKAGE_ICON,
				MIMEType: "image/png",
				Sizes:    []string{"256x256"},
			},
		},
	}, loadBonusToolHandler)

	// =============================================================================
	// Elicitation Tools - Request user input during tool execution
	//
	// WHY ELICITATION MATTERS:
	// Elicitation allows tools to request additional information from users
	// mid-execution, enabling interactive workflows. This is essential for:
	//   - Confirming destructive actions before they happen
	//   - Gathering missing parameters that weren't provided upfront
	//   - Implementing approval workflows for sensitive operations
	//   - Collecting feedback or additional context during execution
	//
	// TWO ELICITATION MODES:
	// - Form (schema): Display a structured form with typed fields in the client
	// - URL: Open a web page (e.g., OAuth flow, feedback form, documentation)
	// =============================================================================

	mcp.AddTool(server, &mcp.Tool{
		Name:        "confirm_action",
		Title:       "Confirm Action",
		Description: "Demonstrates elicitation - requests user confirmation before proceeding",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Confirm Action",
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  false, // User response varies
			OpenWorldHint:   boolPtr(false),
		},
	}, confirmActionHandler)

	mcp.AddTool(server, &mcp.Tool{
		Name:        "get_feedback",
		Title:       "Get Feedback",
		Description: "Demonstrates URL elicitation - opens a feedback form in the browser",
		Annotations: &mcp.ToolAnnotations{
			Title:           "Get Feedback",
			ReadOnlyHint:    true,
			DestructiveHint: boolPtr(false),
			IdempotentHint:  false, // User response varies
			OpenWorldHint:   boolPtr(true), // Opens external URL
		},
	}, getFeedbackHandler)
}

func helloHandler(_ context.Context, _ *mcp.CallToolRequest, input helloInput) (*mcp.CallToolResult, any, error) {
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Hello, %s! Welcome to MCP.", input.Name)},
		},
	}, nil, nil
}

func weatherHandler(_ context.Context, _ *mcp.CallToolRequest, input weatherInput) (*mcp.CallToolResult, any, error) {
	conditions := []string{"sunny", "cloudy", "rainy", "windy"}
	weather := Weather{
		Location:    input.Location,
		Temperature: 15 + rand.Intn(20),
		Unit:        "celsius",
		Conditions:  conditions[rand.Intn(len(conditions))],
		Humidity:    40 + rand.Intn(40),
	}

	jsonBytes, _ := json.MarshalIndent(weather, "", "  ")
	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: string(jsonBytes)},
		},
	}, weather, nil
}

func askLLMHandler(ctx context.Context, req *mcp.CallToolRequest, input askLLMInput) (*mcp.CallToolResult, any, error) {
	maxTokens := input.MaxTokens
	if maxTokens == 0 {
		maxTokens = 100
	}

	result, err := req.Session.CreateMessage(ctx, &mcp.CreateMessageParams{
		Messages: []*mcp.SamplingMessage{
			{
				Role:    "user",
				Content: &mcp.TextContent{Text: input.Prompt},
			},
		},
		MaxTokens: int64(maxTokens),
	})
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Sampling not supported or failed: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	text := "[non-text response]"
	if tc, ok := result.Content.(*mcp.TextContent); ok {
		text = tc.Text
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("LLM Response: %s", text)},
		},
	}, nil, nil
}

func longTaskHandler(ctx context.Context, req *mcp.CallToolRequest, input longTaskInput) (*mcp.CallToolResult, any, error) {
	steps := 5
	progressToken := req.Params.GetProgressToken()

	for i := 0; i < steps; i++ {
		if progressToken != nil {
			_ = req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
				ProgressToken: progressToken,
				Progress:      float64(i) / float64(steps),
				Total:         1.0,
				Message:       fmt.Sprintf("Step %d/%d", i+1, steps),
			})
		}
		time.Sleep(1 * time.Second)
	}

	if progressToken != nil {
		_ = req.Session.NotifyProgress(ctx, &mcp.ProgressNotificationParams{
			ProgressToken: progressToken,
			Progress:      1.0,
			Total:         1.0,
			Message:       "Complete!",
		})
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("Task %q completed successfully after %d steps!", input.TaskName, steps)},
		},
	}, nil, nil
}

func loadBonusToolHandler(_ context.Context, _ *mcp.CallToolRequest, _ any) (*mcp.CallToolResult, any, error) {
	if bonusToolLoaded {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Bonus tool is already loaded! Try calling 'bonus_calculator'."},
			},
		}, nil, nil
	}

	if globalServer != nil {
		mcp.AddTool(globalServer, &mcp.Tool{
			Name:        "bonus_calculator",
			Title:       "Bonus Calculator",
			Description: "A calculator that was dynamically loaded",
			Annotations: &mcp.ToolAnnotations{
				Title:           "Bonus Calculator",
				ReadOnlyHint:    true, // Pure computation
				DestructiveHint: boolPtr(false),
				IdempotentHint:  true, // Same inputs = same outputs
				OpenWorldHint:   boolPtr(false),
			},
			Icons: []mcp.Icon{
				{
					Source:   ABACUS_ICON,
					MIMEType: "image/png",
					Sizes:    []string{"256x256"},
				},
			},
		}, calculatorHandler)
	}
	bonusToolLoaded = true

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: "Bonus tool 'bonus_calculator' has been loaded! Refresh your tools list to see it."},
		},
	}, nil, nil
}

func calculatorHandler(_ context.Context, _ *mcp.CallToolRequest, input calculatorInput) (*mcp.CallToolResult, any, error) {
	var result float64
	switch input.Operation {
	case "add":
		result = input.A + input.B
	case "subtract":
		result = input.A - input.B
	case "multiply":
		result = input.A * input.B
	case "divide":
		if input.B != 0 {
			result = input.A / input.B
		} else {
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: "Error: division by zero"},
				},
				IsError: true,
			}, nil, nil
		}
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			&mcp.TextContent{Text: fmt.Sprintf("%v %s %v = %v", input.A, input.Operation, input.B, result)},
		},
	}, nil, nil
}

// =============================================================================
// Elicitation Handlers
//
// Elicitation requests return one of three actions:
//   - "accept": User provided the requested information
//   - "decline": User explicitly refused to provide information
//   - "cancel": User dismissed the request without responding
//
// Always handle all three cases gracefully.
// =============================================================================

func confirmActionHandler(ctx context.Context, req *mcp.CallToolRequest, input confirmActionInput) (*mcp.CallToolResult, any, error) {
	// Form elicitation: Display a structured form with typed fields
	// The client renders this as a dialog/form based on the JSON schema
	result, err := req.Session.Elicit(ctx, &mcp.ElicitParams{
		Message: fmt.Sprintf("Please confirm: %s", input.Action),
		RequestedSchema: map[string]any{
			"type": "object",
			"properties": map[string]any{
				"confirm": map[string]any{
					"type":        "boolean",
					"title":       "Confirm",
					"description": "Confirm the action",
				},
				"reason": map[string]any{
					"type":        "string",
					"title":       "Reason",
					"description": "Optional reason for your choice",
				},
			},
			"required": []string{"confirm"},
		},
	})
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Elicitation not supported or failed: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	switch result.Action {
	case "accept":
		if confirmed, ok := result.Content["confirm"].(bool); ok && confirmed {
			reason := "No reason provided"
			if r, ok := result.Content["reason"].(string); ok && r != "" {
				reason = r
			}
			return &mcp.CallToolResult{
				Content: []mcp.Content{
					&mcp.TextContent{Text: fmt.Sprintf("Action confirmed: %s\nReason: %s", input.Action, reason)},
				},
			}, nil, nil
		}
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Action declined by user: %s", input.Action)},
			},
		}, nil, nil
	case "decline":
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("User declined to respond for: %s", input.Action)},
			},
		}, nil, nil
	case "cancel":
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("User cancelled elicitation for: %s", input.Action)},
			},
		}, nil, nil
	default:
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("Unexpected elicitation response: %s", result.Action)},
			},
		}, nil, nil
	}
}

func getFeedbackHandler(ctx context.Context, req *mcp.CallToolRequest, input feedbackInput) (*mcp.CallToolResult, any, error) {
	// URL elicitation: Open a web page in the user's browser
	// Useful for OAuth flows, external forms, documentation links, etc.
	feedbackURL := "https://github.com/SamMorrowDrums/mcp-starters/issues/new?template=workshop-feedback.yml"
	if input.Topic != "" {
		feedbackURL += "&title=" + input.Topic
	}

	// Request user to visit URL via URL elicitation
	result, err := req.Session.Elicit(ctx, &mcp.ElicitParams{
		Mode:    "url",
		Message: "Please provide feedback on MCP Starters by completing the form at the URL below:",
		URL:     feedbackURL,
	})
	if err != nil {
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: fmt.Sprintf("URL elicitation not supported or failed: %v", err)},
			},
			IsError: true,
		}, nil, nil
	}

	switch result.Action {
	case "accept":
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Thank you for providing feedback! Your input helps improve MCP Starters."},
			},
		}, nil, nil
	case "decline":
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "No problem! Feel free to provide feedback anytime at: " + feedbackURL},
			},
		}, nil, nil
	case "cancel":
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Feedback request cancelled."},
			},
		}, nil, nil
	default:
		return &mcp.CallToolResult{
			Content: []mcp.Content{
				&mcp.TextContent{Text: "Feedback URL: " + feedbackURL},
			},
		}, nil, nil
	}
}
