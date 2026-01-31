# Tools Implementation Specification

## ADDED Requirements

### Requirement: Tools Save New Prompt Handler
The system SHALL implement a tools/saveNewPrompt handler using MCP SDK's ToolHandlerFor pattern.

#### Scenario: Save Valid Prompt
- **WHEN** a client requests tools/saveNewPrompt with valid prompt data
- **THEN** the server SHALL create and save the new prompt

#### Scenario: Invalid Prompt Data
- **WHEN** a client requests tools/saveNewPrompt with invalid prompt data
- **THEN** the server SHALL return an appropriate error

### Requirement: Tool Input Validation
The system SHALL validate input for tool operations.

#### Scenario: Missing Required Fields
- **WHEN** required fields are missing in tool requests
- **THEN** the server SHALL return an appropriate error

#### Scenario: Invalid Field Types
- **WHEN** field types are invalid in tool requests
- **THEN** the server SHALL return an appropriate error

### Requirement: Tool Error Handling
The system SHALL provide proper error handling for tool operations.

#### Scenario: Database Error During Tool Operation
- **WHEN** a database error occurs during tool operations
- **THEN** the server SHALL return an appropriate error response

#### Scenario: Duplicate Prompt Name
- **WHEN** a tool operation attempts to create a duplicate prompt
- **THEN** the server SHALL return an appropriate error

## MODIFIED Requirements

### Requirement: Tool Response Format
The tool response format SHALL be updated to match MCP SDK standards.

#### Scenario: Tool Success Response
- **WHEN** a tool operation succeeds
- **THEN** the response SHALL use MCP-standardized format

#### Scenario: Tool Error Response
- **WHEN** a tool operation fails
- **THEN** the error response SHALL use MCP-standardized format