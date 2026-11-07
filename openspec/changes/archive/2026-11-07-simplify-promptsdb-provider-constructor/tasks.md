## 1. Provider registration and constructor in promptsdb

- [x] 1.1 Add `RegisterProvider(name string, ctor func(ProviderConfiguration, string) (Provider, error))` and a package-level registry to `internal/promptsdb/new.go`
- [x] 1.2 Implement `New(dbProvider string, config ProviderConfiguration, logfile string) (Provider, error)` in `internal/promptsdb/new.go`: filesystem for `FILE_SYSTEM_PROVIDER`/`""`/unknown, registered constructor for `SQLITE_PROVIDER`, error if unregistered
- [x] 1.3 Keep the `FILE_SYSTEM_PROVIDER` / `SQLITE_PROVIDER` constants in `new.go`

## 2. SQLite package registration

- [x] 2.1 Add `init()` to `internal/promptsdb/sqlite/sqlite.go` (`//go:build sqlite`) registering a constructor that maps `ProviderConfiguration.SQLite` to `SQLiteProviderConfiguration` and calls `NewSQLiteProvider`
- [x] 2.2 Add `init()` to `internal/promptsdb/sqlite/sqlite_stub.go` (`//go:build !sqlite`) registering a constructor that returns the "built without SQLite support" error

## 3. Move tests out of the factory package

- [x] 3.1 Create `internal/promptsdb/new_test.go` (`package promptsdb_test`, blank-import the sqlite package, `//go:build sqlite` for sqlite cases) with the factory tests adapted to call `promptsdb.New`
- [x] 3.2 Add a filesystem/fallback test for `New` that runs without the sqlite build tag

## 4. Update main.go and delete the factory package

- [x] 4.1 In `main.go`: remove the `factory` import, add blank import `_ "github.com/hkionline/prompter/internal/promptsdb/sqlite"`, change the call to `db, err := promptsdb.New(providerName, config.Storage, config.LogFile)`
- [x] 4.2 Delete `internal/promptsdb/factory/factory.go` and `internal/promptsdb/factory/factory_test.go`

## 5. Verify

- [x] 5.1 `go build ./...` and `go build -tags sqlite ./...` both succeed
- [x] 5.2 `go vet ./...` is clean
- [x] 5.3 `make test` passes (unit + integration), including with `-tags sqlite` if covered by the Makefile
