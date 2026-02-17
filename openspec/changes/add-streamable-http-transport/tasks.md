## 1. Configuration Updates

- [x] 1.1 Add transport field to configuration structure
- [x] 1.2 Update configuration parsing to handle transport field
- [x] 1.3 Add validation for transport field values
- [x] 1.4 Update default configuration to include transport field

## 2. Transport Interface

- [x] 2.1 Define Transport interface with Start() method
- [x] 2.2 Create stdio transport implementation
- [x] 2.3 Create HTTP transport implementation using Go MCP SDK

## 3. HTTP Transport Implementation

- [x] 3.1 Add Go MCP SDK dependency to project
- [x] 3.2 Implement HTTP server adapter for MCP protocol
- [x] 3.3 Add HTTP request/response handling
- [x] 3.4 Implement error handling for HTTP transport
- [x] 3.5 Add configuration for HTTP server (port, etc.)

## 4. Core Server Updates

- [x] 4.1 Modify prompter.Run to use transport interface
- [x] 4.2 Update server initialization logic
- [x] 4.3 Ensure RPC handlers work with both transports

## 5. Testing

- [x] 5.1 Add unit tests for transport interface implementations
- [x] 5.2 Add integration tests for HTTP transport
- [x] 5.3 Verify backward compatibility with stdio transport
- [x] 5.4 Test configuration parsing and validation

## 6. Documentation

- [x] 6.1 Update README with HTTP transport information
- [x] 6.2 Add configuration examples for both transports
- [x] 6.3 Document when to use each transport type