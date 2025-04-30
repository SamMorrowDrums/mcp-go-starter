package example

import (
	"context"
	"fmt"

	"github.com/spf13/viper"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// NewHelloTool returns the hello tool and its handler
func NewHelloTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("hello_world",
		mcp.WithDescription("Say hello to someone"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Say hello",
			ReadOnlyHint: true,
		}),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		name, ok := request.Params.Arguments["name"].(string)
		if !ok || name == "" {
			return mcp.NewToolResultError("name must be a string"), nil
		}
		greeting := viper.GetString("greeting")
		if greeting == "" {
			greeting = "Hello"
		}
		return mcp.NewToolResultText(fmt.Sprintf("%s, %s!", greeting, name)), nil
	}
	return tool, handler
}

// NewEnumTool returns a tool that demonstrates enum usage
func NewEnumTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("choose_color",
		mcp.WithDescription("Choose a color from a predefined set of options"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Choose a color",
			ReadOnlyHint: true,
		}),
		mcp.WithString("color",
			mcp.Required(),
			mcp.Description("The color to choose"),
			mcp.Enum("red", "green", "blue"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		color, ok := request.Params.Arguments["color"].(string)
		if !ok || color == "" {
			return mcp.NewToolResultError("color must be one of: red, green, blue"), nil
		}
		return mcp.NewToolResultText(fmt.Sprintf("You chose the color: %s", color)), nil
	}
	return tool, handler
}

// NewWeatherTool returns a tool that provides weather information for a given city
func NewWeatherTool() (mcp.Tool, server.ToolHandlerFunc) {
	tool := mcp.NewTool("get_weather",
		mcp.WithDescription("Get weather information for a city"),
		mcp.WithToolAnnotation(mcp.ToolAnnotation{
			Title:        "Get Weather",
			ReadOnlyHint: true,
		}),
		mcp.WithString("city",
			mcp.Required(),
			mcp.Description("Name of the city to get weather for"),
		),
	)
	handler := func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
		city, ok := request.Params.Arguments["city"].(string)
		if !ok || city == "" {
			return mcp.NewToolResultError("city must be a string and cannot be empty"), nil
		}

		// For MVP, return a simple weather report for Seattle area
		return mcp.NewToolResultText(fmt.Sprintf("Weather report for %s: cloudy, windy and rainy", city)), nil
	}
	return tool, handler
}

// RegisterTools registers all tools with the server
func RegisterTools(s *server.MCPServer) {
	tool, handler := NewHelloTool()
	s.AddTool(tool, handler)

	colorTool, colorHandler := NewEnumTool()
	s.AddTool(colorTool, colorHandler)

	weatherTool, weatherHandler := NewWeatherTool()
	s.AddTool(weatherTool, weatherHandler)
}
