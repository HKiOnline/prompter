## 1. Configuration Updates

- [ ] 1.1 Add transport field to configuration structure
- [ ] 1.2 Update configuration parsing to handle transport field
- [ ] 1.3 Add validation for transport field values
- [ ] 1.4 Update default configuration to include transport field

## 2. Transport Interface

- [ ] 2.1 Define Transport interface with Start() method
- [ ] 2.2 Create stdio transport implementation
- [ ] 2.3 Create HTTP transport implementation using Go MCP SDK

## 3. HTTP Transport Implementation

- [ ] 3.1 Add Go MCP SDK dependency to project
- [ ] 3.2 Implement HTTP server adapter for MCP protocol
- [ ] 3.3 Add HTTP request/response handling
- [ ] 3.4 Implement error handling for HTTP transport
- [ ] 3.5 Add configuration for HTTP server (port, etc.)

## 4. Core Server Updates

- [ ] 4.1 Modify prompter.Run to use transport interface
- [ ] 4.2 Update server initialization logic
- [ ] 4.3 Ensure RPC handlers work with both transports

## 5. Testing

- [ ] 5.1 Add unit tests for transport interface implementations
- [ ] 5.2 Add integration tests for HTTP transport
- [ ] 5.3 Verify backward compatibility with stdio transport
- [ ] 5.4 Test configuration parsing and validation

## 6. Documentation

- [ ] 6.1 Update README with HTTP transport information
- [ ] 6.2 Add configuration examples for both transports
- [ ] 6.3 Document when to use each transport type