# provider-constructor Specification

## Purpose
Define how a storage Provider is constructed from configuration via the unified `promptsdb.New` constructor, including provider registration for sub-package backends.

## Requirements
### Requirement: Unified provider constructor in promptsdb package
The `promptsdb` package SHALL expose `New(dbProvider string, config ProviderConfiguration, logfile string) (Provider, error)` as the single constructor that selects and builds a storage `Provider` by name.

#### Scenario: Filesystem provider selected
- **WHEN** `New` is called with `dbProvider` equal to `promptsdb.FILE_SYSTEM_PROVIDER`
- **THEN** it returns a filesystem provider rooted at `config.Filesystem.Directory`

#### Scenario: Empty provider name falls back to filesystem
- **WHEN** `New` is called with `dbProvider` equal to `""`
- **THEN** it returns a filesystem provider rooted at `config.Filesystem.Directory`

#### Scenario: Unknown provider name falls back to filesystem
- **WHEN** `New` is called with a `dbProvider` that matches no registered provider
- **THEN** it returns a filesystem provider rooted at `config.Filesystem.Directory`

### Requirement: SQLite provider via registration
The `sqlite` sub-package SHALL register its constructor with `promptsdb` in an `init()` function, and `promptsdb.New` SHALL use the registered constructor when `dbProvider` equals `promptsdb.SQLITE_PROVIDER`.

#### Scenario: SQLite provider selected with build tag
- **WHEN** the binary is built with `-tags sqlite` and `New` is called with `promptsdb.SQLITE_PROVIDER`
- **THEN** it returns a SQLite provider configured from `config.SQLite`

#### Scenario: SQLite provider selected without build tag
- **WHEN** the binary is built without `-tags sqlite` and `New` is called with `promptsdb.SQLITE_PROVIDER`
- **THEN** it returns an error stating the binary was built without SQLite support

#### Scenario: SQLite provider selected but not registered
- **WHEN** `New` is called with `promptsdb.SQLITE_PROVIDER` in a binary that does not import the `sqlite` package
- **THEN** it returns an error indicating no provider is registered for that name

### Requirement: Provider registration API
The `promptsdb` package SHALL expose `RegisterProvider(name string, ctor func(ProviderConfiguration, string) (Provider, error))` so sub-packages can register constructors without creating an import cycle.

#### Scenario: Registering a provider
- **WHEN** a sub-package calls `promptsdb.RegisterProvider("sqlite", ctor)` in its `init()`
- **THEN** `promptsdb.New("sqlite", ...)` invokes `ctor`

### Requirement: Factory package removed
The `internal/promptsdb/factory` package SHALL no longer exist; its constructor logic and tests SHALL live in the `promptsdb` package.

#### Scenario: Callers use promptsdb.New
- **WHEN** server startup code needs a storage provider
- **THEN** it calls `promptsdb.New(providerName, config.Storage, config.LogFile)` with no dependency on a `factory` package

