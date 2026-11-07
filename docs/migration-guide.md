# Migration Guide: Filesystem to SQLite

This guide explains how to migrate existing prompts from the filesystem storage provider to the SQLite storage provider.

## Prerequisites

- Prompter built with the `sqlite` build tag: `go build -tags sqlite -o prompter main.go`
- The migrate tool built: `go build -tags sqlite -o prompter-migrate ./cmd/migrate`

## Migration Steps

### 1. Build the migration tool

```bash
go build -tags sqlite -o prompter-migrate ./cmd/migrate
```

### 2. Dry-run (recommended first step)

Validate what would be migrated without writing to the database:

```bash
./prompter-migrate -src ~/.config/prompter/prompts -dst ~/.config/prompter/prompts.db --dry-run --verbose
```

This shows which prompt files would be migrated and which would be skipped (e.g., corrupted files).

### 3. Run the migration

```bash
./prompter-migrate -src ~/.config/prompter/prompts -dst ~/.config/prompter/prompts.db --verbose
```

### 4. Update configuration

Edit `~/.config/prompter/prompter.yaml`:

```yaml
prompter:
  storage:
    provider: "sqlite"
    sqlite:
      path: "~/.config/prompter/prompts.db"
```

### 5. Verify

Start prompter and verify your prompts are available via your MCP client.

## Migration Tool Flags

| Flag | Description |
|---|---|
| `-src` | Source directory containing prompt `.md` files (required) |
| `-dst` | Destination SQLite database file path (required) |
| `--dry-run` | Validate without writing to the database |
| `--incremental` | Only migrate prompts not already present in the target database |
| `--verbose` | Print detailed progress for each file |

## Behavior

- **Source directory is never modified** — migration is read-only on the source
- **Corrupted files are skipped** with warnings in the report
- **Name conflicts** are skipped by default (existing prompts in the target are not overwritten)
- **Incremental mode** (`--incremental`) skips prompts that already exist in the target database
- A **migration report** is printed at the end showing totals, conflicts, warnings, and errors

## Rollback

If you need to go back to filesystem storage, simply change `storage.provider` back to `"filesystem"` in your configuration. The original prompt files in the source directory are untouched.
