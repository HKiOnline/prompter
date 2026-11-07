# Architecture Overview

## Components

```
main.go
  └── configuration.New()          # Load & validate prompter.yaml + env overrides
  └── factory.New()                # Create storage provider based on config
        ├── promptsdb.NewPromptsFsProvider()   # Filesystem provider (default)
        └── sqlite.NewSQLiteProvider()         # SQLite provider (build tag: sqlite)
  └── server.New()                 # MCP server with stdio/HTTP transport
        └── prompts (RPC handlers) # prompts/list, prompts/get
        └── tools  (RPC handlers) # tools/saveNewPrompt
```

## Storage Providers

Both providers implement the `promptsdb.Provider` interface:

```go
type Provider interface {
    Create(prompt Prompt) error
    Read(promptId string) (Prompt, error)
    Update(prompt Prompt) error
    Delete(promptId string) error
    List(query PromptQuery) ([]Prompt, error)
}
```

### Filesystem Provider (`internal/promptsdb/fsProvider.go`)

- Stores each prompt as a `.md` file with YAML frontmatter
- Caches all prompts in memory on startup
- Default provider; always available

### SQLite Provider (`internal/promptsdb/sqlite/`)

- Stores prompts in a single SQLite database table
- Only compiled with `-tags sqlite` build tag
- Supports transactions, connection pooling, and indexed queries
- Schema versioned via `schema_version` table

## Configuration Flow

```
prompter.yaml ──► koanf ──► Configuration struct ──► factory.New() ──► Provider
                    ▲
                    └── Environment variables (PROMPTER_*)
```

## Build Tags

| Tag | Effect |
|---|---|
| (none) | Filesystem provider only; SQLite stub returns error |
| `sqlite` | Full SQLite support; migrate command available |

## Migration Tool

`cmd/migrate/` — standalone CLI to migrate prompts from filesystem directory to SQLite database.
