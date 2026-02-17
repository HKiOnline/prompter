## ADDED Requirements

### Requirement: Transport Configuration Field
The system SHALL support a transport configuration field in the configuration file.

#### Scenario: Default Transport
- **WHEN** no transport is specified in configuration
- **THEN** system uses "stdio" as default transport

#### Scenario: Streamable HTTP Transport
- **WHEN** transport is set to "streamable_http" in configuration
- **THEN** system uses Streamable HTTP transport

#### Scenario: Invalid Transport
- **WHEN** invalid transport type is specified
- **THEN** system returns configuration error

### Requirement: Configuration File Format
The system SHALL maintain backward compatibility with existing configuration files.

#### Scenario: Existing Configuration
- **WHEN** existing configuration file without transport field is used
- **THEN** system defaults to stdio transport

#### Scenario: New Configuration Format
- **WHEN** new configuration file with transport field is used
- **THEN** system uses specified transport type