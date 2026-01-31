# MCP SDK Implementation

This document describes the MCP Go SDK implementation in Prompter MCP Server.

## Overview

Prompter uses the official [MCP Go SDK](https://github.com/modelcontextprotocol/go-sdk) to implement the [Model Context Protocol](https://modelcontextprotocol.io/docs). This provides a standardized, well-tested foundation for MCP server functionality.

## Architecture

The SDK-based implementation follows these key components:

### 1. Server Initialization

The `SDKServer` struct in `internal/server/server.go` is the main entry point for the SDK-based server:

```go
type Server struct {
    config   *config.Config
    logger   plog.Logger
    db       promptsdb.Provider
}
```

The server struct can be extended as needed. The server is started in the main.go file.

### 2. Tool Handlers

Prompter implements one MCP-tool and number of MCP-prompt calls:

**Tools**:
- **tools/saveNewPrompt**: Creates and saves a new prompt

**Prompts**:
- **prompts/list**: Lists all available prompts
- **prompts/get**: Retrieves a specific prompt by name

Tool handlers are defined in `internal/tools/tools.go` and follow the MCP SDK's `ToolHandlerFor` pattern. Prompt handlers are defined in `internal/prompts/prompts.go`. Each feature of MCP should its own directory under the `internal` directory.

### 3. Protocol Compliance

The SDK ensures compliance with:
- JSON-RPC 2.0 over stdio
- MCP protocol specification (version 2024-10-07)
- Proper error handling and logging

### Breaking Changes

- RPC method names have changed to match MCP conventions
- Response formats are now standardized by the SDK
- Some error messages may differ slightly

### Benefits

- Better protocol compliance
- Access to SDK features and updates
- Improved compatibility with MCP clients
- Easier maintenance and debugging

## Testing

Each server feature has its own directory under internal. These directories also include the unit tests for each feature. Test coverage should be at least 80%.

## Resources

- [MCP Protocol Specification](https://modelcontextprotocol.io)
- [MCP Go SDK Documentation](https://github.com/modelcontextprotocol/go-sdk)
- [MCP Tool Handbook](https://modelcontextprotocol.io/docs/concepts/tools)
