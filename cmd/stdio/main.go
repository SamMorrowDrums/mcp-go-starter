// MCP Go Starter - stdio Transport
//
// This entrypoint runs the MCP server using stdio transport,
// which is ideal for local development and CLI tool integration.
//
// Usage:
//
//	go run ./cmd/stdio
//
// Documentation: https://modelcontextprotocol.io/docs/develop/transports#stdio
package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/SamMorrowDrums/mcp-go-starter/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Failed to connect: %v", err)
		os.Exit(1)
	}
}

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	// Handle shutdown signals
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigChan
		cancel()
	}()

	// Create the MCP server
	srv := server.NewServer()
	server.SetGlobalServer(srv)

	// Create stdio transport
	transport := &mcp.StdioTransport{}

	// Connect and run
	log.SetOutput(os.Stderr) // Don't interfere with stdio protocol
	log.Println("MCP Go Starter running on stdio")

	if _, err := srv.Connect(ctx, transport, nil); err != nil {
		return err
	}

	// Wait for context cancellation
	<-ctx.Done()
	log.Println("Server shutting down")
	return nil
}
