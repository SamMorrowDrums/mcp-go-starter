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
	Name string `json:"name" jsonschema:"description=The name to greet"`
}

type weatherInput struct {
	Location string `json:"location" jsonschema:"description=City name or coordinates"`
}

type askLLMInput struct {
	Prompt    string `json:"prompt" jsonschema:"description=The question or prompt for the LLM"`
	MaxTokens int    `json:"maxTokens,omitempty" jsonschema:"description=Maximum tokens in response"`
}

type longTaskInput struct {
	TaskName string `json:"taskName" jsonschema:"description=Name for this task"`
}

type calculatorInput struct {
	A         float64 `json:"a" jsonschema:"description=First number"`
	B         float64 `json:"b" jsonschema:"description=Second number"`
	Operation string  `json:"operation" jsonschema:"enum=add,enum=subtract,enum=multiply,enum=divide"`
}

// =============================================================================
// Tool Annotations - Every tool MUST have annotations for AI assistants
//
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
				Source:      "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Waving%20hand/3D/waving_hand_3d.png",
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
				Source:      "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Sun%20behind%20cloud/3D/sun_behind_cloud_3d.png",
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
				Source:      "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Robot/3D/robot_3d.png",
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
				Source:      "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Hourglass%20not%20done/3D/hourglass_not_done_3d.png",
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
				Source:      "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Package/3D/package_3d.png",
				MIMEType: "image/png",
				Sizes:    []string{"256x256"},
			},
		},
	}, loadBonusToolHandler)
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
					Source:      "https://raw.githubusercontent.com/microsoft/fluentui-emoji/main/assets/Abacus/3D/abacus_3d.png",
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
