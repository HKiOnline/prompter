# Prompt File Format

## ADDED Requirements
### Requirement: Prompt File Documentation
The system SHALL provide documentation about the prompt file format in the docs directory.

#### Scenario: User reads prompt file documentation
- **WHEN** user navigates to docs/prompt-file-format.md
- **THEN** they can read detailed information about the prompt file structure and format

## MODIFIED Requirements
### Requirement: Prompt File Extension
The system SHALL use `.md` extension for prompt files instead of `.yaml`.

#### Scenario: Creating a new prompt
- **WHEN** user creates a new prompt via MCP tool
- **THEN** the prompt is saved as a .md file in the prompts directory

#### Scenario: Reading an existing prompt
- **WHEN** user reads a prompt via MCP prompt call
- **THEN** the system returns the prompt content regardless of file extension
