## Why

The Prompter MCP Server currently only supports stdio transport, which limits its deployment flexibility. Adding Streamable HTTP transport support will enable the server to be used in web-based environments and remote scenarios, expanding its usability and compatibility with modern AI applications.

## What Changes

- Add Streamable HTTP transport implementation using the official Go MCP SDK
- Add configuration option for "streamable_http" transport type
- Modify prompter.Run to check transport type and start appropriate server
- Maintain stdio as the default transport for backward compatibility
- Add new HTTP server implementation that handles MCP protocol over HTTP
- Update configuration structure to support transport selection

## Capabilities

### New Capabilities
- `streamable-http-transport`: Support for Streamable HTTP transport method alongside existing stdio transport
- `transport-configuration`: Configuration system for selecting between transport methods

### Modified Capabilities
- `core-server`: Update the main server initialization to support multiple transport types

## Impact

- Core server initialization logic in `prompter.Run`
- Configuration structure and parsing
- New HTTP server implementation files
- Potential updates to RPC handler registration to work with both transports
- Documentation updates for new transport option