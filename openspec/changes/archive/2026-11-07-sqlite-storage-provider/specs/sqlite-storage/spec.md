# SQLite Storage Provider Specification

## REQUIREMENTS

The SQLite storage provider must implement the existing Provider interface with SQLite as the backend.

### Core Functionality
- [ ] Implement all methods from the Provider interface:
  - `ListPrompts()`: Return all prompts stored in SQLite database
  - `GetPrompt(name string)`: Retrieve a prompt by name from SQLite
  - `SavePrompt(prompt *Prompt)`: Store or update a prompt in SQLite
  - `DeletePrompt(name string)`: Remove a prompt from SQLite
- [ ] Store prompts in a single table with the following structure:
  - `id`: Primary key (UUID)
  - `name`: Unique identifier for the prompt
  - `title`: Human-readable title
  - `description`: Optional description
  - `content`: Markdown content of the prompt
  - `metadata`: JSON representation of YAML frontmatter
  - `created_at`: Timestamp when prompt was created
  - `updated_at`: Timestamp when prompt was last updated
  - `tags`: Comma-separated list of tags
- [ ] Support transactions for atomic operations across multiple prompts
- [ ] Implement proper indexing for efficient querying by name and tags

### Data Integrity
- [ ] Validate prompt data before storage (name, content required)
- [ ] Handle concurrent access with appropriate locking
- [ ] Prevent duplicate prompt names
- [ ] Maintain referential integrity for all relationships

### Performance
- [ ] Optimize queries with appropriate indexes
- [ ] Implement connection pooling for database connections
- [ ] Support batch operations for bulk inserts/updates

## BEHAVIOR

### Scenario 1: Storing a New Prompt
**Given**: A valid prompt with unique name
**When**: `SavePrompt()` is called
**Then**: 
- Prompt is inserted into the database
- `created_at` and `updated_at` timestamps are set
- Return error if prompt name already exists

### Scenario 2: Updating an Existing Prompt
**Given**: An existing prompt in the database
**When**: `SavePrompt()` is called with same name
**Then**:
- Prompt content and metadata are updated
- `updated_at` timestamp is refreshed
- Return error if prompt does not exist

### Scenario 3: Listing Prompts
**Given**: Multiple prompts in the database
**When**: `ListPrompts()` is called
**Then**:
- Return all prompts sorted by name (ascending)
- Include all metadata fields
- Return empty list if no prompts exist

### Scenario 4: Deleting a Prompt
**Given**: An existing prompt in the database
**When**: `DeletePrompt(name)` is called
**Then**:
- Prompt is removed from the database
- Return nil error on success
- Return error if prompt does not exist

### Scenario 5: Transaction Rollback
**Given**: Multiple prompts being saved in a transaction
**When**: One operation fails
**Then**:
- All changes are rolled back
- Database remains in consistent state
- Error is propagated to caller

## VALIDATION

### Unit Tests
- [ ] Test each Provider method with various inputs
- [ ] Verify database schema and constraints
- [ ] Test error handling for invalid data
- [ ] Test concurrent access scenarios

### Integration Tests
- [ ] Test with actual RPC methods that use the provider
- [ ] Verify data persistence across restarts
- [ ] Test migration from fsProvider to SQLite

### Acceptance Criteria
- [ ] All Provider interface methods work correctly
- [ ] Data integrity maintained in all scenarios
- [ ] Performance meets expectations (TBD)
- [ ] No breaking changes to existing functionality

## DEPENDENCIES

### External Dependencies
- `github.com/mattn/go-sqlite3`: SQLite driver for Go
- Standard library database/sql package

### Internal Dependencies
- Existing Provider interface definition
- Prompt data structure and validation logic
- Configuration system for database path
