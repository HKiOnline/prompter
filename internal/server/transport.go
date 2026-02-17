package server

import (
	"context"
	"fmt"
	"net/http"

	"github.com/hkionline/prompter/internal/configuration"
	"github.com/modelcontextprotocol/go-sdk/mcp"
)

// transport interface defines the methods that a transport must implement
type transport interface {
	start(ctx context.Context, server *mcp.Server, config *configuration.Configuration) error
}

// stdioTransport implements the Transport interface for stdio transport
type stdioTransport struct{}

func (t *stdioTransport) start(ctx context.Context, server *mcp.Server, config *configuration.Configuration) error {
	transport := mcp.NewStdioTransport()
	return server.Run(ctx, transport)
}

// httpTransport implements the Transport interface for Streamable HTTP transport
type httpTransport struct {
	httpServer *http.Server
}

func (t *httpTransport) start(ctx context.Context, server *mcp.Server, config *configuration.Configuration) error {
	// Create Streamable HTTP handler using the Go MCP SDK
	handler := mcp.NewStreamableHTTPHandler(
		func(request *http.Request) *mcp.Server {
			return server
		},
		&mcp.StreamableHTTPOptions{},
	)

	// Create HTTP server
	t.httpServer = &http.Server{
		Addr:    fmt.Sprintf(":%d", config.HTTP.Port),
		Handler: handler,
	}

	// Start the HTTP server in a goroutine
	go func() {
		if err := t.httpServer.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("HTTP server error: %v\n", err)
		}
	}()

	// Wait for context cancellation
	<-ctx.Done()
	return t.httpServer.Shutdown(ctx)
}

// newTransport creates a new transport instance based on the configuration
func newTransport(transportType string) (transport, error) {
	switch transportType {
	case "stdio":
		return &stdioTransport{}, nil
	case "streamable_http":
		return &httpTransport{}, nil
	default:
		return nil, fmt.Errorf("unknown transport type: %s", transportType)
	}
}
