# Change: Migrate from Custom MCP Implementation to Official MCP Go SDK

## Why
The current implementation uses a custom MCP protocol implementation that needs to be replaced with the official MCP Go SDK. This migration will provide better protocol compliance, access to SDK features and updates, improved compatibility with MCP clients, and easier maintenance and debugging.

## What Changes
- Replace custom RPC implementation with MCP Go SDK
- Create new server implementation using SDK patterns
- Implement tool handlers using SDK's ToolHandlerFor pattern
- Implement prompt handlers using SDK patterns
- Update main.go to use new SDK-based server
- Keep existing storage provider (promptsdb) and configuration unchanged
- Maintain all existing functionality (prompts/list, prompts/get, tools/saveNewPrompt)

## Impact
- Affected specs: None (new capabilities)
- Affected code:
  - internal/rpc/ (will be replaced)
  - internal/server/ (new implementation)
  - internal/prompts/ (new implementation)
  - internal/tools/ (new implementation)
  - main.go (updated to use new server)
  - internal/transport/stdio/ (updated to use new server)
- Breaking changes: RPC method names and response formats will change to match MCP conventions
- Backward compatibility: Existing prompts and storage will remain compatible

## Dependencies
- MCP Go SDK (github.com/modelcontextprotocol/go-sdk)
- Existing promptsdb, configuration, and plog packages will remain unchanged