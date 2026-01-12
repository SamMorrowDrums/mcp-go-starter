// MCP Go Starter - HTTP Transport
//
// This entrypoint runs the MCP server using HTTP with SSE streams,
// which is ideal for remote deployment and web-based clients.
//
// Usage:
//
//	go run ./cmd/http
//	PORT=8080 go run ./cmd/http
//
// Documentation: https://modelcontextprotocol.io/docs/develop/transports#streamable-http
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/SamMorrowDrums/mcp-go-starter/internal/server"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

func main() {
	if err := run(); err != nil {
		log.Printf("Server error: %v", err)
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

	port := os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}

	// Create HTTP handler for MCP
	handler := mcp.NewStreamableHTTPHandler(func(_ *http.Request) *mcp.Server {
		srv := server.NewServer()
		server.SetGlobalServer(srv)
		return srv
	}, nil)

	// Set up routes
	mux := http.NewServeMux()
	mux.Handle("/mcp", handler)
	mux.HandleFunc("/health", healthHandler)

	addr := fmt.Sprintf(":%s", port)

	// Start server
	log.Printf("MCP Go Starter running on http://localhost%s", addr)
	log.Printf("  MCP endpoint: http://localhost%s/mcp", addr)
	log.Printf("  Health check: http://localhost%s/health", addr)
	log.Println("Press Ctrl+C to exit")

	httpServer := &http.Server{
		Addr:    addr,
		Handler: mux,
	}

	// Graceful shutdown
	go func() {
		<-ctx.Done()
		log.Println("Shutting down server...")
		_ = httpServer.Shutdown(context.Background())
	}()

	if err := httpServer.ListenAndServe(); err != http.ErrServerClosed {
		return err
	}
	return nil
}

func healthHandler(w http.ResponseWriter, _ *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	_, _ = w.Write([]byte(`{"status":"ok","server":"mcp-go-starter","version":"1.0.0"}`))
}
