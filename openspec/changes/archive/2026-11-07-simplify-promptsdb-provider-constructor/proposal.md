# Simplify PromptsDB Provider Constructor

## Why

The SQLite storage provider change (v0.7.0) introduced a separate `internal/promptsdb/factory` package whose only job is a `switch` on the provider name. Before that change, `promptsdb.New(dbProvider, config, logfile)` lived directly in `internal/promptsdb/new.go` and `main.go` called it as `promptsdb.New(promptsdb.FILE_SYSTEM_PROVIDER, config.Storage, config.LogFile)` — a more idiomatic Go constructor. The extra package adds an import and indirection for callers without adding value: there is exactly one call site (`main.go`).

## What Changes

- Move the provider-selection logic from `internal/promptsdb/factory/factory.go` back into `internal/promptsdb/new.go` as `promptsdb.New(dbProvider string, config ProviderConfiguration, logfile string) (Provider, error)`.
- `main.go` calls `promptsdb.New(providerName, config.Storage, config.LogFile)` again; the `factory` import is removed.
- Move `factory/factory_test.go` into the `promptsdb` package (e.g. `new_test.go`) and update it to call `promptsdb.New`.
- Delete the `internal/promptsdb/factory` package.
- Because `promptsdb/sqlite` imports `promptsdb` (for the `Provider` interface), `promptsdb` cannot import `promptsdb/sqlite` back (circular import). The SQLite branch is wired via a provider-registration hook: the `sqlite` package registers its constructor in an `init()`, and `main.go` blank-imports the `sqlite` package so registration runs. This is the same pattern as `database/sql` drivers.

## Capabilities

### New Capabilities
- `provider-constructor`: unified `promptsdb.New` constructor that selects and builds a storage `Provider` by name (filesystem or sqlite), with provider registration for backends that live in sub-packages.

### Modified Capabilities
- None. The `main.go` change from `factory.New(...)` to `promptsdb.New(...)` is an internal call-site/import change with no spec-level behavior change, so no delta specs are needed.

## Impact

- **Affected code**:
  - `internal/promptsdb/new.go` — gains the `New` constructor and provider-registration hook
  - `internal/promptsdb/new_test.go` — moved/adapted factory tests (new file)
  - `internal/promptsdb/sqlite/sqlite.go` and `sqlite_stub.go` — `init()` registers the SQLite constructor
  - `main.go` — call site and imports
  - `internal/promptsdb/factory/` — deleted
- **APIs**: internal only; the RPC interface and `Provider` interface are unchanged. `factory.New` was introduced in the unreleased v0.7.0 branch, so removing it is not a breaking change for any release.
- **Dependencies**: no new external dependencies.
