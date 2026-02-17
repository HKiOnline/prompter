# Transport Configuration

## Purpose
Define how transport is configured, validated, and interpreted by the server.

## Requirements
### Requirement: Structured transport object in configuration
The system MUST support transport configuration as an object under `prompter.transport`, and this object MUST include a `type` field used to select the transport implementation.

#### Scenario: Parse nested transport type
- **WHEN** configuration defines `prompter.transport.type` as `stdio`
- **THEN** the server SHALL load configuration successfully and select the stdio transport

#### Scenario: Parse nested streamable HTTP transport type
- **WHEN** configuration defines `prompter.transport.type` as `streamable_http`
- **THEN** the server SHALL load configuration successfully and select the streamable HTTP transport

### Requirement: Streamable HTTP options are scoped under transport
The system MUST read streamable HTTP transport options from `prompter.transport.streamable_http`, and transport startup MUST use values from that nested object.

#### Scenario: Use nested streamable HTTP port
- **WHEN** configuration sets `prompter.transport.type` to `streamable_http` and `prompter.transport.streamable_http.port` to a custom value
- **THEN** the HTTP transport SHALL bind using that configured port

### Requirement: Transport type validation
The system MUST validate `prompter.transport.type` and reject values other than `stdio` and `streamable_http` with a clear configuration error.

#### Scenario: Reject unknown transport type
- **WHEN** configuration sets `prompter.transport.type` to an unsupported value
- **THEN** configuration loading SHALL fail with an error describing valid transport type values

### Requirement: Legacy top-level transport keys are unsupported
The system MUST treat legacy top-level transport keys as invalid for this change, including scalar `prompter.transport` and top-level `prompter.http`.

#### Scenario: Legacy scalar transport is rejected
- **WHEN** configuration uses `prompter.transport` as a scalar string instead of an object
- **THEN** configuration loading SHALL fail and indicate that the nested transport object format is required

#### Scenario: Legacy top-level http does not configure streamable HTTP transport
- **WHEN** configuration sets top-level `prompter.http` but omits `prompter.transport.streamable_http`
- **THEN** streamable HTTP transport SHALL use only supported nested transport fields and SHALL ignore unsupported legacy keys
