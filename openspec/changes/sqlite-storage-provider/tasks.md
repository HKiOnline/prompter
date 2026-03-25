# SQLite Storage Provider Implementation Tasks

## 1. Setup and Dependencies

- [ ] 1.1 Add go-sqlite3 dependency to go.mod
- [ ] 1.2 Create internal/promptsdb/sqlite directory structure
- [ ] 1.3 Set up build tags for optional SQLite support

## 2. SQLite Provider Implementation

- [ ] 2.1 Implement SQLiteProvider struct with Provider interface methods
- [ ] 2.2 Create database schema migration logic
- [ ] 2.3 Implement ListPrompts() with proper indexing
- [ ] 2.4 Implement GetPrompt(name string) with error handling
- [ ] 2.5 Implement SavePrompt(prompt *Prompt) with validation
- [ ] 2.6 Implement DeletePrompt(name string)
- [ ] 2.7 Add transaction support for atomic operations
- [ ] 2.8 Implement connection pooling and timeout handling

## 3. Configuration System Updates

- [ ] 3.1 Add storage.provider configuration option
- [ ] 3.2 Add SQLite-specific configuration options (path, timeout, etc.)
- [ ] 3.3 Update config validation logic
- [ ] 3.4 Add environment variable overrides for sensitive options
- [ ] 3.5 Update configuration documentation

## 4. Provider Factory and Initialization

- [ ] 4.1 Update provider factory to support SQLite
- [ ] 4.2 Implement provider selection based on configuration
- [ ] 4.3 Add fsProvider as default fallback
- [ ] 4.4 Update main.go initialization logic

## 5. Migration Tool Implementation

- [ ] 5.1 Create migrate command structure
- [ ] 5.2 Implement fsProvider to SQLite migration logic
- [ ] 5.3 Add dry-run mode for validation
- [ ] 5.4 Implement progress reporting
- [ ] 5.5 Add migration report generation
- [ ] 5.6 Handle error cases and corrupted files

## 6. Unit Tests

- [ ] 6.1 Write unit tests for SQLiteProvider methods
- [ ] 6.2 Test database schema and constraints
- [ ] 6.3 Test error handling scenarios
- [ ] 6.4 Test concurrent access scenarios
- [ ] 6.5 Test transaction rollback behavior

## 7. Integration Tests

- [ ] 7.1 Test with actual RPC methods using SQLite provider
- [ ] 7.2 Verify data persistence across restarts
- [ ] 7.3 Test migration from fsProvider to SQLite
- [ ] 7.4 Test configuration options at runtime
- [ ] 7.5 Test environment variable overrides

## 8. Documentation Updates

- [ ] 8.1 Update README with SQLite storage option
- [ ] 8.2 Add configuration reference for new options
- [ ] 8.3 Document migration process and tool usage
- [ ] 8.4 Add examples for SQLite provider configuration
- [ ] 8.5 Update architecture documentation with new component

## 9. Final Validation

- [ ] 9.1 Verify backward compatibility with existing prompts
- [ ] 9.2 Test all RPC methods work correctly
- [ ] 9.3 Verify no breaking changes to existing functionality
- [ ] 9.4 Run full test suite including new tests
- [ ] 9.5 Create release notes for new feature
