## Context

The Prompter MCP Server currently uses stdio transport exclusively. This design introduces Streamable HTTP transport as an alternative, enabling web-based deployment scenarios while maintaining backward compatibility.

## Goals / Non-Goals

**Goals:**
- Add Streamable HTTP transport support using the official Go MCP SDK
- Maintain stdio as the default transport for backward compatibility
- Allow configuration-based transport selection
- Ensure both transports can coexist and work with existing RPC handlers
- Follow MCP protocol specifications for Streamable HTTP transport

**Non-Goals:**
- Replace stdio transport entirely
- Support other transport protocols beyond stdio and Streamable HTTP
- Modify existing prompt storage or management functionality
- Change the RPC method signatures or behavior

## Decisions

### Use Official Go MCP SDK for Streamable HTTP
**Decision**: Use the official Go MCP SDK's Streamable HTTP implementation rather than building from scratch.
**Rationale**: The SDK provides a tested, specification-compliant implementation that handles the complex streaming protocol requirements. This reduces development time and ensures compatibility with other MCP implementations.
**Alternatives Considered**: Building custom implementation, but this would require extensive testing and might not be fully specification-compliant.

### Configuration-Based Transport Selection
**Decision**: Add a `transport` field to the configuration that allows selecting between "stdio" and "streamable_http".
**Rationale**: This provides flexibility without breaking existing configurations. Existing installations will continue to use stdio by default.
**Alternatives Considered**: Command-line flag, but configuration file approach is more maintainable for production deployments.

### Separate HTTP Server Implementation
**Decision**: Create a new HTTP server implementation that adapts the existing RPC handlers to work with Streamable HTTP transport.
**Rationale**: This maintains separation of concerns and allows each transport to handle its specific protocol requirements while reusing the core business logic.
**Alternatives Considered**: Modifying existing stdio server to support both, but this would create a tightly coupled design.

### Transport Interface Abstraction
**Decision**: Create a transport interface within the server package that both stdio and HTTP transports implement, allowing prompter.Run to start the appropriate transport without conditional logic.
**Rationale**: This follows Go interfaces best practices and makes it easier to add additional transports in the future. Keeping this within the server package maintains cohesion since transport is tightly coupled with server functionality.
**Alternatives Considered**: 
- Creating a separate transport package - rejected because transport is too closely tied to server initialization
- Direct conditional logic in prompter.Run - rejected as less extensible

## Risks / Trade-offs

**[Dependency on External SDK]** → The Go MCP SDK is officially maintained by the MCP organization, reducing risk of abandonment. We'll vendor the dependency to ensure build reproducibility.

**[Increased Complexity]** → The addition of a new transport layer adds complexity. We'll mitigate this by keeping the transports isolated and well-documented.

**[Performance Differences]** → HTTP transport may have different performance characteristics than stdio. We'll document these differences and provide guidance on when to use each transport.

**[Configuration Migration]** → Existing users might need to understand the new transport option. We'll maintain stdio as default and provide clear documentation about the new option.

**[Protocol Compatibility]** → Ensuring both transports work identically with all RPC methods. We'll implement comprehensive integration tests that verify both transports produce the same results.

## Migration Plan

1. **Implementation Phase**: Add the new transport implementation alongside existing stdio transport
2. **Testing Phase**: Verify both transports work correctly with all existing functionality
3. **Documentation Phase**: Update documentation to explain the new transport option and when to use it
4. **Release Phase**: Release as a new version with stdio remaining the default
5. **Monitoring Phase**: Monitor for any issues with the new transport in production use

**Rollback Strategy**: If issues arise, users can simply revert to stdio transport by changing the configuration or using the previous version.

## Open Questions

- Should we support running both transports simultaneously for debugging purposes?
- What should be the default port for HTTP transport?
- Should we add transport-specific configuration options (e.g., HTTP timeout settings)?
- How should we handle transport-specific errors differently?