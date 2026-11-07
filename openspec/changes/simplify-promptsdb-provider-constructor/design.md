# Design: Simplify PromptsDB Provider Constructor

## Context

`main.go` currently builds the storage provider via `factory.New(providerName, config.Storage, config.LogFile)`, where `internal/promptsdb/factory` is a thin `switch` over `promptsdb.FILE_SYSTEM_PROVIDER` / `promptsdb.SQLITE_PROVIDER`. Before v0.7.0 the same selection lived in `promptsdb.New(...)` inside `internal/promptsdb/new.go`, called as `promptsdb.New(promptsdb.FILE_SYSTEM_PROVIDER, config.Storage, config.LogFile)`. The factory package has a single call site and adds an import layer without behavior.

Constraint: `internal/promptsdb/sqlite` imports `internal/promptsdb` (it returns `promptsdb.Provider` and uses `promptsdb.Prompt`/`PromptQuery`). Therefore `promptsdb` cannot import `promptsdb/sqlite` — a direct `import` of the sqlite package from `new.go` would be a circular import that the Go compiler rejects, regardless of build tags.

## Goals / Non-Goals

**Goals:**
- Restore `promptsdb.New(dbProvider, config, logfile)` as the single public constructor for storage providers.
- Support both filesystem and sqlite through that one function.
- Remove the `factory` package entirely; move its tests into the `promptsdb` package.
- Keep behavior identical: unknown/empty provider names fall back to the filesystem provider; sqlite without the build tag returns the existing "built without SQLite support" error.

**Non-Goals:**
- No changes to the `Provider` interface, the RPC layer, or prompt semantics.
- No new storage backends.
- No changes to configuration schema or environment variables.

## Decisions

### D1: Provider registration via `init()` instead of a direct import

`promptsdb.New` selects the provider with a `switch` on the name. The filesystem branch calls `NewPromptsFsProvider` directly (same package). The sqlite branch cannot call `sqlite.NewSQLiteProvider` from within `promptsdb` (circular import), so the `sqlite` package registers a constructor function with `promptsdb` in its `init()`:

- `promptsdb` exposes `RegisterProvider(name string, ctor func(ProviderConfiguration, string) (Provider, error))` and keeps a package-level registry map.
- `sqlite.go` (build tag `sqlite`) and `sqlite_stub.go` (tag `!sqlite`) each get an `init()` that registers `promptsdb.SQLITE_PROVIDER`. The stub registers a constructor returning the existing "built without SQLite support" error, so behavior without the tag is unchanged.
- `main.go` blank-imports `_ "github.com/hkionline/prompter/internal/promptsdb/sqlite"` so the registration `init()` runs in both build modes. `main` importing `sqlite` is acyclic: `main → sqlite → promptsdb`, `main → promptsdb`.

**Alternatives considered:**
- *Direct `switch` calling `sqlite.NewSQLiteProvider` from `promptsdb`*: rejected — circular import, compile error.
- *Move `Provider`/`Prompt` types into a separate `types` package*: correct long-term fix for the dependency direction, but touches every provider and the server package; out of scope for this refactor.
- *Keep the switch in `main.go`*: rejected — the user explicitly wants the factory behavior back in `promptsdb.New`, not pushed to the caller.

This is the same pattern as `database/sql` driver registration, so it is idiomatic Go.

### D2: Keep fallback semantics

`New` keeps the current factory semantics: `FILE_SYSTEM_PROVIDER` and `""` → filesystem; unknown names → filesystem fallback (current factory behavior); `SQLITE_PROVIDER` with no registered constructor → error. `main.go` keeps its existing `if providerName == "" { providerName = promptsdb.FILE_SYSTEM_PROVIDER }` normalization, so the empty-string path is exercised in practice.

### D3: Tests move with the code

`factory/factory_test.go` moves to `internal/promptsdb/new_test.go`, package `promptsdb`, calling `New(...)` directly. The sqlite-tagged tests stay behind the existing `//go:build sqlite` tag. Registration makes the sqlite branch testable from within the `promptsdb` package only if the `sqlite` package is linked into the test binary; to avoid an import cycle in tests, the moved tests either import `sqlite` (test files in package `promptsdb` *can* import `sqlite` only from an external test package `promptsdb_test`) — therefore the moved test file uses `package promptsdb_test` and blank-imports the sqlite package, mirroring `main.go`.

## Risks / Trade-offs

- **[Risk]** Registration via `init()` means a binary that never imports the `sqlite` package gets "provider not registered" instead of a compile-time link — mitigated by the blank import in `main.go` and by the stub registering an explicit "built without -tags sqlite" error, preserving today's runtime error.
- **[Risk]** Test placement: `package promptsdb` test files cannot import `sqlite` (cycle); using `package promptsdb_test` for the moved factory tests avoids this — the test file asserts behavior through the public `New` only.
- **[Trade-off]** The registry map is mutable global state; acceptable here because registration happens once in `init()` and lookups are read-only afterwards. A `sync.RWMutex` is unnecessary given init-time writes.

## Migration Plan

Single commit on the feature branch: move constructor + tests, update `main.go`, delete `factory/`. No data migration — this is a pure code refactor on an unreleased branch. Rollback is `git revert` of the commit.

## Open Questions

None. The registration approach is fully determined by the import-cycle constraint.
