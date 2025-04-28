package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"

	"github.com/mark3labs/mcp-go/server"
	example "github.com/sammorrowdrums/mcp-go-starter/pkg/example"
)

func main() {
	rootCmd := &cobra.Command{
		Use:   "mcp",
		Short: "MCP server starter",
		Long:  `A starter MCP server using cobra and viper.`,
	}

	stdioCmd := &cobra.Command{
		Use:   "stdio",
		Short: "Run the MCP server in stdio mode",
		Run: func(cmd *cobra.Command, args []string) {
			greeting := viper.GetString("greeting")
			_ = viper.GetString("secret") // secret is loaded for completeness, but not used directly here

			s := example.NewServer()
			example.RegisterTools(s)
			example.RegisterResources(s)
			example.RegisterPrompts(s)

			fmt.Fprintf(os.Stderr, "[INFO] Starting MCP server with MCP_GREETING: %s\n", greeting)
			if err := server.ServeStdio(s); err != nil {
				fmt.Fprintf(os.Stderr, "Server error: %v\n", err)
				os.Exit(1)
			}
		},
	}

	viper.SetEnvPrefix("MCP")
	viper.AutomaticEnv()
	viper.SetDefault("greeting", "Hello")
	viper.SetDefault("secret", "Hello")

	rootCmd.AddCommand(stdioCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
