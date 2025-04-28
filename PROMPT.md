**NOTE** This is a vibe coding prompt that I actually used to bootstrap this repo. I did a few tweaks afterwards, but mostly it was automatic.


# Prompt

You are a senior software engineer leading a workshop on creating MCP servers and you are going to build a Go codebase from scratch with a basic tool, that provides one argument, you will also server a markdown file as an MCP resource, and add an example prompt. Make heavy use of the GitHub MCP Server tools if they can help you to find relevent context.

# Instructions

The example package should use the latest https://github.com/mark3labs/mcp-go that will be the library for all MCP server code. You can look up its readme and code on GitHub if you need to, but you will also have it locally once it is installed via go get.

## CLI command

Use cobra and viper to create a main.go cmd

The main command to run the server should be called `stdio`, and should take env vars that start with MCP_ using viper. Initially we want to take MCP_GREETING and MCP_SECRET with a default value of `Hello`.

## Server

The server.go file should initialize a basic mcp server with a name, version, plus tool, resource and prompt functionality enabled. Use the type `*server.MCPServer` fot the server.

## Tools

We want a hello world tool exposed that will take a name argument, that is returned using a ServerTool from MCP Go, to combine the handler and the mcp.Tool. 

## Resources

We want to server the a markdown file in a subfolder of the source file (so you can emebd it, you cannot use ../ syntax for embedding) of the repository as a resource in the server (you will need to embed it as a string as we won't always run it in the same folder). You must import `_ "embed"` to actually embed the example.md as a string in the application, don't just copy the text into a string.

Store this file in the pkg/example/resources folder

## .vscode/mcp.json 

Should contain a json configuration for vscode insiders. The GitHub Repo github/github-mcp-server has a good example of the format, and should set inputs MCP_GREETING and MCP_SECRET, and the secret input should be a password type, the greeting doesn't need to be.

We want to use the go run command in there for rapid development.

## Readme

The README.md file should contain clear instructions for `go run cmd/mcp/main.go stdio` as well as a `go build` version, but first is should contain an instruction to add a full path for the go run command in the .vscode/mcp.json file as then you should be able to run it automatically.

## Codespaces devcontainer config

Create a basic devcontainer config with go tools enabled (could be from the base image or from a feature if available), that means people can develop in a codespace.

# Reasoning Steps

- Think through your plan
- Actually get the readme files of all the libraries you plan on using, and refer to them for more context. Also crawl the entire mcp-go examples directory and read each source file. YOU MUST DO THIS with github mcp get contents tools. DO NOT SKIP.
- Build a working cmd first and let it just fm.Println 
- Use this as a checkpoint to validate go run works
- Next build the basic mcp-go server
- Then add the tool and register with the server
- Next add a test for the tool 
- Then add the resource and register with the server
- Next add a test for the resource
- Then add the prompts and register with the server
- Next add a test for the prompt 

Always fix linter problems, or type problems rather than just re-running tests really think about how syntax issues can be fixed. Actualy check the exported types from the libraries using tools like the GitHub MCP Server tool get file contents, especially for README files, but also for limited numbers of source files. You should fix errors by thinking carefully and checking sources. Don't just keep retrying without having a quality plan to fix them.

Finally, build up the readme and the other files we will need. If you are unsure or need input during any steps then please do. You also have access to the GitHub MCP Server, so use the tools and resources as you need (but be careful not to overload your context so be sparing with whole files)

# Output Format

Create go code files with:

- The main package cmd/mcp/main.go
- The example package in pkg/example (with a server.go, tools.go, resources.go and prompts.go all in the same package)
- Any other files, scripts or makefiles you need to keep things simple

**NOTE** ALL MCP TOOLS, RESOURCES AND PROMPTS SHOULD HAVE SEPARATE FUNCTIONS FOR DEFINING THEMSELVES WITH THEIR HANDLERS AND FUNCTIONS TO REGISTER THEM WITH THE SERVER, DO NOT COMBINE REGISTRATION WITH CREATION.

# Examples
## Very basic server code

```
func main() {
    // Create MCP server
    s := server.NewMCPServer(
        "Demo 🚀",
        "1.0.0",
    )

    // Add tool
    tool := mcp.NewTool("hello_world",
        mcp.WithDescription("Say hello to someone"),
        mcp.WithString("name",
            mcp.Required(),
            mcp.Description("Name of the person to greet"),
        ),
    )

    // Add tool handler
    s.AddTool(tool, helloHandler)

    // Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, ok := request.Params.Arguments["name"].(string)
    if !ok {
        return nil, errors.New("name must be a string")
    }

    return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
```


## Basic cobra example

```
var rootCmd = &cobra.Command{
  Use:   "hugo",
  Short: "Hugo is a very fast static site generator",
  Long: `A Fast and Flexible Static Site Generator built with
                love by spf13 and friends in Go.
                Complete documentation is available at https://gohugo.io/documentation/`,
  Run: func(cmd *cobra.Command, args []string) {
    // Do Stuff Here
  },
}

func Execute() {
  if err := rootCmd.Execute(); err != nil {
    fmt.Fprintln(os.Stderr, err)
    os.Exit(1)
  }
}
```

## Basic viper ENV VAR Docs

%%%EXAMPLE STARTS HERE%%%

### Working with Environment Variables

Viper has full support for environment variables. This enables 12 factor
applications out of the box. There are five methods that exist to aid working
with ENV:

 * `AutomaticEnv()`
 * `BindEnv(string...) : error`
 * `SetEnvPrefix(string)`
 * `SetEnvKeyReplacer(string...) *strings.Replacer`
 * `AllowEmptyEnv(bool)`

_When working with ENV variables, it’s important to recognize that Viper
treats ENV variables as case sensitive._

Viper provides a mechanism to try to ensure that ENV variables are unique. By
using `SetEnvPrefix`, you can tell Viper to use a prefix while reading from
the environment variables. Both `BindEnv` and `AutomaticEnv` will use this
prefix.

`BindEnv` takes one or more parameters. The first parameter is the key name, the
rest are the name of the environment variables to bind to this key. If more than
one are provided, they will take precedence in the specified order. The name of
the environment variable is case sensitive. If the ENV variable name is not provided, then
Viper will automatically assume that the ENV variable matches the following format: prefix + "_" + the key name in ALL CAPS. When you explicitly provide the ENV variable name (the second parameter),
it **does not** automatically add the prefix. For example if the second parameter is "id",
Viper will look for the ENV variable "ID".

One important thing to recognize when working with ENV variables is that the
value will be read each time it is accessed. Viper does not fix the value when
the `BindEnv` is called.

`AutomaticEnv` is a powerful helper especially when combined with
`SetEnvPrefix`. When called, Viper will check for an environment variable any
time a `viper.Get` request is made. It will apply the following rules. It will
check for an environment variable with a name matching the key uppercased and
prefixed with the `EnvPrefix` if set.

`SetEnvKeyReplacer` allows you to use a `strings.Replacer` object to rewrite Env
keys to an extent. This is useful if you want to use `-` or something in your
`Get()` calls, but want your environmental variables to use `_` delimiters. An
example of using it can be found in `viper_test.go`.

Alternatively, you can use `EnvKeyReplacer` with `NewWithOptions` factory function.
Unlike `SetEnvKeyReplacer`, it accepts a `StringReplacer` interface allowing you to write custom string replacing logic.

By default empty environment variables are considered unset and will fall back to
the next configuration source. To treat empty environment variables as set, use
the `AllowEmptyEnv` method.

#### Env example

```go
SetEnvPrefix("spf") // will be uppercased automatically
BindEnv("id")

os.Setenv("SPF_ID", "13") // typically done outside of the app

id := Get("id") // 13
```

### Working with Flags

Viper has the ability to bind to flags. Specifically, Viper supports `Pflags`
as used in the [Cobra](https://github.com/spf13/cobra) library.

Like `BindEnv`, the value is not set when the binding method is called, but when
it is accessed. This means you can bind as early as you want, even in an
`init()` function.

For individual flags, the `BindPFlag()` method provides this functionality.

Example:

```go
serverCmd.Flags().Int("port", 1138, "Port to run Application server on")
viper.BindPFlag("port", serverCmd.Flags().Lookup("port"))
```

%%%EXAMPLE ENDS HERE%%%

## Example tool from the github MCP server

```
// GetIssue creates a tool to get details of a specific issue in a GitHub repository.
func GetIssue(getClient GetClientFn, t translations.TranslationHelperFunc) (tool mcp.Tool, handler server.ToolHandlerFunc) {
	return mcp.NewTool("get_issue",
			mcp.WithDescription(t("TOOL_GET_ISSUE_DESCRIPTION", "Get details of a specific issue in a GitHub repository.")),
			mcp.WithToolAnnotation(mcp.ToolAnnotation{
				Title:        t("TOOL_GET_ISSUE_USER_TITLE", "Get issue details"),
				ReadOnlyHint: true,
			}),
			mcp.WithString("owner",
				mcp.Required(),
				mcp.Description("The owner of the repository"),
			),
			mcp.WithString("repo",
				mcp.Required(),
				mcp.Description("The name of the repository"),
			),
			mcp.WithNumber("issue_number",
				mcp.Required(),
				mcp.Description("The number of the issue"),
			),
		),
		func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
			owner, err := requiredParam[string](request, "owner")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			repo, err := requiredParam[string](request, "repo")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}
			issueNumber, err := RequiredInt(request, "issue_number")
			if err != nil {
				return mcp.NewToolResultError(err.Error()), nil
			}

			client, err := getClient(ctx)
			if err != nil {
				return nil, fmt.Errorf("failed to get GitHub client: %w", err)
			}
			issue, resp, err := client.Issues.Get(ctx, owner, repo, issueNumber)
			if err != nil {
				return nil, fmt.Errorf("failed to get issue: %w", err)
			}
			defer func() { _ = resp.Body.Close() }()

			if resp.StatusCode != http.StatusOK {
				body, err := io.ReadAll(resp.Body)
				if err != nil {
					return nil, fmt.Errorf("failed to read response body: %w", err)
				}
				return mcp.NewToolResultError(fmt.Sprintf("failed to get issue: %s", string(body))), nil
			}

			r, err := json.Marshal(issue)
			if err != nil {
				return nil, fmt.Errorf("failed to marshal issue: %w", err)
			}

			return mcp.NewToolResultText(string(r)), nil
		}
}

```

## An example of prompts from MCP Go

```

// Simple greeting prompt
s.AddPrompt(mcp.NewPrompt("greeting",
    mcp.WithPromptDescription("A friendly greeting prompt"),
    mcp.WithArgument("name",
        mcp.ArgumentDescription("Name of the person to greet"),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    name := request.Params.Arguments["name"]
    if name == "" {
        name = "friend"
    }
    
    return mcp.NewGetPromptResult(
        "A friendly greeting",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleAssistant,
                mcp.NewTextContent(fmt.Sprintf("Hello, %s! How can I help you today?", name)),
            ),
        },
    ), nil
})

// Code review prompt with embedded resource
s.AddPrompt(mcp.NewPrompt("code_review",
    mcp.WithPromptDescription("Code review assistance"),
    mcp.WithArgument("pr_number",
        mcp.ArgumentDescription("Pull request number to review"),
        mcp.RequiredArgument(),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    prNumber := request.Params.Arguments["pr_number"]
    if prNumber == "" {
        return nil, fmt.Errorf("pr_number is required")
    }
    
    return mcp.NewGetPromptResult(
        "Code review assistance",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleSystem,
                mcp.NewTextContent("You are a helpful code reviewer. Review the changes and provide constructive feedback."),
            ),
            mcp.NewPromptMessage(
                mcp.RoleAssistant,
                mcp.NewEmbeddedResource(mcp.ResourceContents{
                    URI: fmt.Sprintf("git://pulls/%s/diff", prNumber),
                    MIMEType: "text/x-diff",
                }),
            ),
        },
    ), nil
})

// Database query builder prompt
s.AddPrompt(mcp.NewPrompt("query_builder",
    mcp.WithPromptDescription("SQL query builder assistance"),
    mcp.WithArgument("table",
        mcp.ArgumentDescription("Name of the table to query"),
        mcp.RequiredArgument(),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    tableName := request.Params.Arguments["table"]
    if tableName == "" {
        return nil, fmt.Errorf("table name is required")
    }
    
    return mcp.NewGetPromptResult(
        "SQL query builder assistance",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleSystem,
                mcp.NewTextContent("You are a SQL expert. Help construct efficient and safe queries."),
            ),
            mcp.NewPromptMessage(
                mcp.RoleAssistant,
                mcp.NewEmbeddedResource(mcp.ResourceContents{
                    URI: fmt.Sprintf("db://schema/%s", tableName),
                    MIMEType: "application/json",
                }),
            ),
        },
    ), nil
})
```

# Context

Build a starter mcp server repo that will be a joy to use, and ensure that you use mcp-go, cobra and viper effectively and any other tools you need.

This file is your bible: https://github.com/mark3labs/mcp-go/blob/main/README.md YOU MUST READ IT IN ITS ENTIRITY AND UNDERSTAND EXACTLY HOW TO APPLY IT TO THIS PROJECT. YOU SHOULD THEN ALSO SEE ALL THE CODE IN https://github.com/mark3labs/mcp-go/tree/main/examples DIRECTORY AS **THAT** IS GOING TO REALLY SHOWCASE HOW THIS SHOULD WORK. DO THIS. The repository is mark3labs/mcp-go and the folder is /examples - you can use template resources via the github mcp server to get the contents of that folder in main, and that will provide all the files and folders within. CRAWL THIS EXHAUSTIVELY. DO NOT SKIP, IF YOU FAIL - TRY AGAIN WITH DIFFERENT TOOLS OR RESOURCES FROM GITHUB MCP SERVER TOOLS HAVE AVAILABLE AND DO NOT GIVE UP. YOU CAN DO THIS, IF YOU GIVE UP YOU WILL NOT DO THE BEST JOB, SO KEEP TRYING UNTIL YOU FIND THE RIGHT WAY TO GET THESE FILES FOR CONTEXT.

# Final instructions

This is a lot of tasks, I want you to carry them out proactively without checking too much but I also want you to update your plans as you go, think, clarify, check documentation if needed and proceed. Generally use the latest versions of go libraries and remember that this should be simple, clear well documented and easy to get started on for developers with a wide range of abilities and knowledge of Golang.

Get started please!