# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

## [0.6.0] - 2026-03-07

### Added
- MCP prompt list change notification support after successful `saveNewPrompt` operations.
- Server-side log entry when prompt list change notification is sent.
- OpenSpec change artifacts for prompt list change notification implementation:
  - `proposal.md`
  - `design.md`
  - `specs/prompt-list-change-notification/spec.md`
  - `tasks.md`

### Changed
- Tool handling updated to support pluggable prompt list change notification callback wiring from server runtime.

### Fixed
- Prompt file loading now parses YAML front matter separately from prompt body content.
- Resolved YAML parsing failures for valid prompt bodies containing colon-delimited lines.
- Added regression/unit coverage for notification behavior and prompt parsing edge cases.

## [0.5.1] - 2026-03-07

### Added
- New client configuration documentation in `docs/client-configuration.md`.

### Changed
- Updated roadmap documentation.
- Removed deprecated OpenCode command files under `.opencode/command/`.
- Updated ignore rules and README references.

## [0.5.0] - 2026-02-17

### Added
- Streamable HTTP transport support for the MCP server.
- Expanded transport configuration support and related OpenSpec specifications.
- Additional integration coverage and transport-focused test improvements.

### Changed
- Refined server transport and startup flow for multiple transport modes.
- Updated configuration defaults and tests for transport handling.
- Documentation updates for MCP SDK implementation and roadmap progress.

### Removed
- Stale dependencies from module configuration.

## [0.4.0] - 2026-02-14

### Added
- Template function for date formatting in prompts
- Documentation on using date function in prompt templates
- OpenSpec project context configuration
- Updated OpenSpec commands and workflows
- New test cases for template functions and date variations

### Changed
- Updated AGENTS.md to reflect OpenSpec changes
- Updated roadmap with dynamic prompt list update
- Updated test prompt template to use new date function
- Reorganized OpenSpec command structure

## [0.3.0] - 2026-02-07

### Added
- Documentation about the prompt file format in the docs directory

### Changed
- Prompt file extension changed from `.yaml` to `.md` (**BREAKING**)

## [0.2.0] - 2026-02-07

### Added
- OpenSpec for structured change proposals
- Improved test coverage tracking and error handling
- Test server script for integration testing

### Changed
- Replaced custom RPC implementation with MCP Go SDK (**BREAKING**)
- Refactored internal architecture to use MCP SDK components

### Fixed
- Enhanced error handling across all RPC methods

### Removed
- Custom JSON-RPC processor and related structures

## [0.1.1] - 2026-01-26

### Added
- Comprehensive test coverage for plog package
- Additional tests for file system provider

### Changed
- Replaced slog implementation with custom plog package

### Fixed
- Improved test reliability and coverage

## [0.1.0] - 2025-11-18

### Added
- Initial implementation of the Prompter MCP server
- Core functionality for prompt management
- File system-based storage provider
- Go templating engine support for dynamic prompts
- JSON-RPC communication protocol over stdio
- Basic RPC methods: initialize, ping, prompts/list, prompts/get, tools/saveNewPrompt 
