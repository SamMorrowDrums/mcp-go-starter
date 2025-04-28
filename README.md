# mcp-go-starter

A starter repo for building a Go MCP server using [mcp-go](https://github.com/mark3labs/mcp-go), Cobra, and Viper.

## Quick Start

1. **Set the full path for your go-run script in `.vscode/mcp.json`:**

   Run this command in your project root to get the absolute path:

   ```sh
   pwd
   ```

   Then update the `command` field in `.vscode/mcp.json` to use the full path, for example:

   ```jsonc
   "command": "/Users/yourname/path/to/mcp-go-starter/script/go-run"
   ```

   This ensures VS Code can run the server automatically from anywhere.

2. **Run the server (development mode) using the script:**

```sh
./script/go-run
```

   This script ensures you always run the latest code from the correct directory, just like the approach in [github/github-mcp-server#51](https://github.com/github/github-mcp-server/pull/51).

3. **Or build and run the server manually:**

```sh
go build -o mcp-server ./cmd/mcp && ./mcp-server stdio
```

## Environment Variables

- `MCP_GREETING`: Greeting to use for the hello tool and prompt (default: `Hello`)
- `MCP_SECRET`: Secret value for demonstration purposes (default: `Hello`)

## Features

- **Hello World Tool**: Takes a `name` argument and returns a greeting.
- **Markdown Resource**: Serves `pkg/example/resources/example.md` as an MCP resource.
- **Prompt Example**: Simple prompt that greets the user by name.

## Example: Calling the hello_world Tool

You can call the `hello_world` tool from your MCP client (such as VS Code with the MCP extension) or programmatically. Here’s an example using the tool with the argument `name: "Sam"`:

**Request:**
```json
{
  "method": "call_tool",
  "params": {
    "tool": "hello_world",
    "arguments": {
      "name": "Sam"
    }
  },
  "id": 1,
  "jsonrpc": "2.0"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 1,
  "result": {
    "type": "text",
    "text": "Hello, Sam!"
  }
}
```

## Example: Calling the choose_color Tool (Enum)

**Request:**
```json
{
  "method": "call_tool",
  "params": {
    "tool": "choose_color",
    "arguments": {
      "color": "green"
    }
  },
  "id": 2,
  "jsonrpc": "2.0"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 2,
  "result": {
    "type": "text",
    "text": "You chose the color: green"
  }
}
```

## Example: Getting a Prompt

**Request:**
```json
{
  "method": "get_prompt",
  "params": {
    "prompt": "hello_prompt",
    "arguments": {
      "name": "Sam"
    }
  },
  "id": 3,
  "jsonrpc": "2.0"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "title": "A friendly greeting",
    "messages": [
      {
        "role": "assistant",
        "content": {
          "type": "text",
          "text": "Hello, Sam! How can I help you today?"
        }
      }
    ]
  }
}
```

## Example: Getting the Markdown Resource

**Request:**
```json
{
  "method": "read_resource",
  "params": {
    "resource": "docs://example"
  },
  "id": 4,
  "jsonrpc": "2.0"
}
```

**Response:**
```json
{
  "jsonrpc": "2.0",
  "id": 4,
  "result": [
    {
      "uri": "docs://example",
      "mime_type": "text/markdown",
      "text": "# Example Markdown Resource\n\nThis is an example markdown file served as an MCP resource.\n\n- You can edit this file to change the resource content.\n- It is embedded in the Go binary using Go's embed package.\n"
    }
  ]
}
```

## Example: Prompt for an LLM Agent

If you want an LLM agent (such as Copilot or another MCP-compatible assistant) to call your tools and resources, you can use a prompt like this:

> You are an assistant with access to the following tools:
> - `hello_world`: Greets a person by name using the configured greeting.
> - `choose_color`: Lets you select a color from red, green, or blue.
> - `docs://example`: A markdown resource with example content.
>
> When a user asks for a greeting, call the `hello_world` tool with their name. If they ask to pick a color, call the `choose_color` tool with the color they want. If they ask for documentation or help, return the contents of the `docs://example` resource.

This prompt will encourage the LLM to use your tools and resources as intended.

## VS Code Integration

- `.vscode/mcp.json` is preconfigured for rapid development with environment variable inputs.
- Make sure to use the full path to the script in the `command` field for best results.

## Codespaces/Devcontainer

A basic devcontainer is provided for GitHub Codespaces with Go tools enabled.

---

For more details, see the [mcp-go README](https://github.com/mark3labs/mcp-go/blob/main/README.md).

## Additional Examples

### Raw JSON-RPC for Accessing Prompts

Note: Only some clients support prompts. Here's an example using direct JSON-RPC over stdio:

```bash
echo '{"jsonrpc":"2.0","id":3,"params":{"name": "hello_prompt"},"method":"prompts/get"}' | go run cmd/mcp/main.go stdio
```

Output:
```json
{
  "jsonrpc": "2.0",
  "id": 3,
  "result": {
    "description": "A friendly greeting",
    "messages": [
      {
        "role": "assistant",
        "content": {
          "type": "text",
          "text": "Hello, friend! How can I help you today?"
        }
      }
    ]
  }
}
```
