# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [0.8.0] - 2026-11-07

### Added
- Embedded version metadata: root `version.json` with a `--version`/`-version` flag that prints `prompter <semver> (<build date>)` and exits before startup, the composed version without build metadata (MAJOR.MINOR.PATCH plus any pre-release) reported as MCP `serverInfo.version`, and build-time stamping of `date` and the SemVer `build` identifier (`lb.{YYYYMMDDHHMMSS}` for local builds, `ci.{pipeline-id}.{YYYYMMDDHHMMSS}` in CI when a pipeline ID is available).
- GoReleaser Homebrew cask publishing to `hkionline/homebrew-tap` for macOS, with automatic quarantine removal for unsigned binaries and stable-release-only tap updates.

### Changed
- GoReleaser release builds and `make build` now include SQLite support by default; filesystem remains the default storage backend.
- Release workflow selects the Go version from `go.mod` and uses a dedicated Homebrew tap token.
- GoReleaser project and binary names standardized to `prompter`, with package metadata configured on the Homebrew cask.
- Homebrew casks now select only macOS Intel and Apple Silicon archives; GitHub Releases also provide explicitly targeted Linux and Windows binaries, with ZIP archives for Windows and stripped executables.
- GitHub releases use curated changelog entries and Homebrew-first installation instructions. The release workflow validates version tags, runs filesystem and SQLite tests, and automatically classifies prereleases.
- Installation documentation now covers Homebrew, prebuilt downloads, Go 1.26 source builds, and MCP executable paths; added a release preparation guide.

### Fixed
- Integration smoke tests now use a valid prompts-list request and fail on missing, malformed, or JSON-RPC error responses, with bounded execution and reliable server cleanup.
- `make test` runs both storage build variants and propagates failures instead of masking test errors.

## [0.7.0] - 2026-11-07

### Added
- SQLite storage provider as an alternative to the filesystem provider, implementing the same `Provider` interface with CRUD operations, indexing, transactions, and connection pooling. Built with `-tags sqlite` build tag; filesystem remains the default.
- `storage.provider` configuration option (`filesystem` or `sqlite`) with provider-specific settings under `storage.sqlite.*` (path, timeout, foreign_keys, max_open_conns, max_idle_conns).
- Environment variable overrides: `PROMPTER_STORAGE_PROVIDER`, `PROMPTER_SQLITE_PATH`, `PROMPTER_SQLITE_TIMEOUT`, `PROMPTER_SQLITE_FOREIGN_KEYS`, `PROMPTER_LOG_FILE`.
- Migration tool (`cmd/migrate`) to copy prompts from a filesystem directory to a SQLite database, with `--dry-run`, `--incremental`, and `--verbose` flags plus a migration report.
- Provider constructor registry in `internal/promptsdb` for selecting the storage backend, with filesystem fallback for empty or unknown provider names. Selecting SQLite in a binary built without SQLite support returns an explicit error.
- Documentation: configuration reference (`docs/configuration-reference.md`), migration guide (`docs/migration-guide.md`), and architecture overview (`docs/architecture.md`).
- Unit tests covering CRUD, schema, validation, concurrency, transactions, and persistence for the SQLite provider.
- Integration tests for migration and provider interface compliance.
- Configuration tests for SQLite options, env var overrides, and validation.

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

## [0.4.0] - 2026-02-15

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

## [0.3.1] - 2026-02-13

### Added
- Project roadmap in `docs/roadmap.md`.

### Changed
- Reorganized README installation, project status, and motivation sections, and made the MCP client executable path example generic.
- Clarified where the Markdown prompt format is used and added documentation references to AGENTS.md and OpenSpec project context.

## [0.3.0] - 2026-02-13

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
- Removed the unused slog handler implementation from the existing custom plog package

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
