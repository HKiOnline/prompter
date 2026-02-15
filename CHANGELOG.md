# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.0.0/)
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

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
