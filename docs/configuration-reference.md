# Configuration Reference

Complete reference for all Prompter configuration options.

## Configuration File Location

```
~/.config/prompter/prompter.yaml
```

## Full Configuration Schema

```yaml
prompter:
  transport:
    type: "stdio"              # "stdio" or "streamable_http"
    streamable_http:
      port: 8080               # HTTP port for streamable_http transport

  logFile: "~/.config/prompter/prompter.log"

  storage:
    provider: "filesystem"     # "filesystem" (default) or "sqlite"

    filesystem:
      prompts_directory: "~/.config/prompter/prompts"

    sqlite:
      path: "~/.config/prompter/prompts.db"  # SQLite database file path
      timeout: 30                            # Connection timeout in seconds
      foreign_keys: true                     # Enable foreign key constraints
      max_open_conns: 10                     # Maximum open database connections
      max_idle_conns: 5                      # Maximum idle database connections
```

## Storage Providers

### filesystem (default)

Stores each prompt as a separate `.md` file with YAML frontmatter in the configured directory.

| Option | Default | Description |
|---|---|---|
| `storage.filesystem.prompts_directory` | `~/.config/prompter/prompts` | Directory containing prompt `.md` files |

### sqlite

Stores all prompts in a single SQLite database file. Requires building with `-tags sqlite`.

| Option | Default | Description |
|---|---|---|
| `storage.sqlite.path` | `~/.config/prompter/prompts.db` | Path to the SQLite database file |
| `storage.sqlite.timeout` | `30` | Connection timeout in seconds |
| `storage.sqlite.foreign_keys` | `true` | Enable SQLite foreign key constraints |
| `storage.sqlite.max_open_conns` | `10` | Maximum number of open database connections |
| `storage.sqlite.max_idle_conns` | `5` | Maximum number of idle database connections |

## Environment Variable Overrides

Environment variables take precedence over configuration file values:

| Variable | Overrides |
|---|---|
| `PROMPTER_STORAGE_PROVIDER` | `storage.provider` |
| `PROMPTER_SQLITE_PATH` | `storage.sqlite.path` |
| `PROMPTER_SQLITE_TIMEOUT` | `storage.sqlite.timeout` |
| `PROMPTER_SQLITE_FOREIGN_KEYS` | `storage.sqlite.foreign_keys` (set to `true`) |
| `PROMPTER_LOG_FILE` | `logFile` |

## Examples

### Filesystem storage (default)

```yaml
prompter:
  transport:
    type: "stdio"
  storage:
    provider: "filesystem"
    filesystem:
      prompts_directory: "/Users/me/my-prompts"
```

### SQLite storage

```yaml
prompter:
  transport:
    type: "stdio"
  storage:
    provider: "sqlite"
    sqlite:
      path: "/Users/me/my-prompts.db"
      timeout: 60
      foreign_keys: true
      max_open_conns: 20
      max_idle_conns: 10
```

### Using environment variables

```bash
export PROMPTER_STORAGE_PROVIDER=sqlite
export PROMPTER_SQLITE_PATH=/Users/me/prompts.db
prompter
```
