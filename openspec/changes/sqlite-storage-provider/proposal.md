# SQLite Storage Provider

## Why

The current file system-based storage provider (fsProvider) works well for simple prompt storage. However, users may benefit from additional storage options that provide:
- Transaction support for atomic operations
- Efficient querying and searching across prompts
- Versioning and history tracking
- Better performance for large prompt collections

A SQLite storage provider would add this capability as an alternative to fsProvider, giving users the choice based on their specific needs. The fsProvider remains the default and is not replaced.

## What Changes

- **New Capability**: SQLite storage provider implementing the Provider interface (in addition to fsProvider)
- **New Capability**: Configuration option to specify which storage provider to use (fsProvider or SQLite)
- **New Capability**: Optional migration tool for existing prompts from fsProvider to SQLite
- **Modified Capability**: Storage provider selection mechanism (non-breaking, maintains fsProvider as default)

## Capabilities

### New Capabilities
- `sqlite-storage`: SQLite-based prompt storage with CRUD operations, indexing, and query support
- `migration-tool`: Tool to migrate existing prompts from fsProvider to SQLite storage
- `config-option`: Configuration option for database path and connection settings

### Modified Capabilities
- None (no requirement changes to existing capabilities)

## Impact

- **Affected Code**: 
  - `internal/promptsdb/` - New SQLite provider implementation (alongside fsProvider)
  - `internal/config/` - Configuration schema updates to support provider selection
  - `cmd/prompter/main.go` - Provider initialization logic with choice between fsProvider and SQLite
- **APIs**: No breaking changes to RPC interface - all providers implement the same interface
- **Dependencies**: Optional `github.com/mattn/go-sqlite3` for SQLite support (only loaded when using SQLite provider)
- **Systems**: 
  - Configuration system extended with storage provider selection
  - Optional migration tool for users who want to switch from fsProvider to SQLite
