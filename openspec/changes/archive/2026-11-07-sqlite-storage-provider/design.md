# SQLite Storage Provider Design

## Context

The Prompter MCP Server currently uses a file system-based storage provider (fsProvider) to store prompts as individual markdown files. This design adds a SQLite-based storage provider as an alternative option, allowing users to choose between the two based on their specific needs.

### Current State
- fsProvider stores prompts as .md files in a directory structure
- Prompts contain YAML frontmatter with metadata and markdown content
- Provider interface defines CRUD operations for prompts
- Configuration system selects which provider to use

### Constraints
- Must maintain backward compatibility with existing prompts
- fsProvider remains the default storage option
- No breaking changes to RPC interface or existing functionality
- SQLite dependency should be optional (only loaded when using SQLite provider)

## Goals / Non-Goals

**Goals:**
- Provide SQLite as an alternative storage provider alongside fsProvider
- Enable efficient querying and searching of prompts
- Support transactions for atomic operations
- Maintain same API interface as fsProvider
- Allow users to choose storage provider via configuration

**Non-Goals:**
- Replace fsProvider as the default storage option
- Modify existing prompt data structures or RPC interface
- Implement complex query languages beyond basic filtering
- Add versioning or history tracking (future enhancement)

## Decisions

### 1. Provider Interface Implementation
**Decision:** Implement SQLiteProvider struct that satisfies the existing Provider interface.

**Rationale:**
- Maintains consistency with fsProvider
- Allows seamless switching between providers
- No changes required to RPC layer or other components

**Alternatives Considered:**
- Creating a new interface: Would require changes throughout the codebase
- Extending existing interface: Would break existing implementations

### 2. Database Schema Design
**Decision:** Store prompts in a single table with JSON for metadata.

```sql
CREATE TABLE prompts (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL,
    title TEXT,
    description TEXT,
    content TEXT NOT NULL,
    metadata TEXT,  -- JSON string of YAML frontmatter
    created_at DATETIME NOT NULL,
    updated_at DATETIME NOT NULL,
    tags TEXT
);
```

**Rationale:**
- Simple schema matches existing fsProvider structure
- JSON metadata preserves flexibility of YAML format
- Easy migration path from fsProvider

**Alternatives Considered:**
- Normalized schema: More complex, harder to migrate existing data
- Separate tables for metadata fields: Less flexible for future changes

### 3. Configuration Approach
**Decision:** Add storage.provider configuration option.

```yaml
storage:
  provider: sqlite  # or "fs" (default)
  sqlite:
    path: ./prompts.db
```

**Rationale:**
- Clear and explicit provider selection
- Allows provider-specific configuration
- Backward compatible (defaults to fsProvider)

**Alternatives Considered:**
- Environment variables: Less discoverable, harder to document
- Command-line flags: Not suitable for long-running server process

### 4. Migration Strategy
**Decision:** Provide optional migration tool as a separate command.

```bash
prompter migrate sqlite ./existing-prompts-dir ./prompts.db
```

**Rationale:**
- Migration is optional (users can start fresh with SQLite)
- Separate command doesn't affect main server functionality
- Clear and explicit migration process

**Alternatives Considered:**
- Automatic migration on first run: Too risky, could lose data
- In-place conversion: Would require complex logic and testing

## Risks / Trade-offs

### [SQLite dependency bloat] → Mitigation: Use build tags to make SQLite optional
- Adding sqlite3 dependency increases binary size
- Only load the dependency when using SQLite provider
- Use build constraints to exclude from default builds

### [Data migration complexity] → Mitigation: Provide clear documentation and examples
- Migrating existing prompts requires careful handling
- Create comprehensive migration guide with error handling
- Provide dry-run option to validate before actual migration

### [Performance characteristics] → Mitigation: Benchmark against fsProvider and document expectations
- SQLite may have different performance profile than filesystem
- Document expected read/write patterns
- Provide configuration options for tuning (e.g., WAL mode)
