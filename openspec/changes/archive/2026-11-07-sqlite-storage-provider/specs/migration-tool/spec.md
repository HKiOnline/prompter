# Migration Tool Specification

## REQUIREMENTS

The migration tool must provide a way to migrate existing prompts from fsProvider to SQLite storage.

### Core Functionality
- [ ] Support migration from fsProvider directory structure to SQLite database
- [ ] Preserve all prompt metadata during migration
- [ ] Handle various file formats and encodings
- [ ] Provide progress reporting during migration
- [ ] Support dry-run mode to validate before actual migration
- [ ] Generate migration report with statistics and warnings

### Data Transformation
- [ ] Convert markdown files to SQLite records
- [ ] Parse YAML frontmatter into structured metadata
- [ ] Handle special characters and escaping in content
- [ ] Preserve file modification timestamps as created_at/updated_at
- [ ] Extract tags from frontmatter or filename patterns

### Error Handling
- [ ] Skip corrupted files with appropriate warnings
- [ ] Continue migration after non-critical errors
- [ ] Provide detailed error report at completion
- [ ] Handle permission issues gracefully
- [ ] Validate output database integrity after migration

## BEHAVIOR

### Scenario 1: Full Migration
**Given**: Existing fsProvider directory with prompts
**When**: `migrate` command is executed
**Then**:
- All valid prompts are migrated to SQLite
- Invalid files are skipped with warnings
- Migration report is generated
- Source directory remains unchanged

### Scenario 2: Dry Run
**Given**: Existing fsProvider directory with prompts
**When**: `migrate --dry-run` is executed
**Then**:
- No actual migration occurs
- Validation checks are performed
- Report shows what would be migrated
- Source directory remains unchanged

### Scenario 3: Incremental Migration
**Given**: Previously migrated prompts with new additions
**When**: `migrate --incremental` is executed
**Then**:
- Only new/changed files are migrated
- Existing prompts are not duplicated
- Report shows only incremental changes

### Scenario 4: Migration with Conflicts
**Given**: Prompts that already exist in target database
**When**: Migration encounters name conflicts
**Then**:
- Conflicts are reported in output
- User can choose to skip, overwrite, or rename
- Default behavior is to skip conflicts

## VALIDATION

### Unit Tests
- [ ] Test file parsing and metadata extraction
- [ ] Test data transformation logic
- [ ] Test error handling for various scenarios
- [ ] Test progress reporting mechanisms

### Integration Tests
- [ ] Test full migration from real fsProvider directory
- [ ] Verify data integrity after migration
- [ ] Test dry-run mode with various inputs
- [ ] Test incremental migration scenarios

### Acceptance Criteria
- [ ] All prompts migrated successfully (or with documented exceptions)
- [ ] No data loss or corruption
- [ ] Migration report is comprehensive and accurate
- [ ] Source directory unchanged after migration
- [ ] Target database is valid and functional

## DEPENDENCIES

### External Dependencies
- SQLite storage provider (target)
- fsProvider (source)
- Standard library file I/O packages

### Internal Dependencies
- Prompt data structure and validation
- YAML parsing library
- Error reporting framework
