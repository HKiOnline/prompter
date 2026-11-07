# Tasks

## 1. Embedded semantic version metadata

- [x] 1.1 Add root version.json with all six fields, initialized to the development version 0.8.0 (date = build-time baseline, empty build); verify the version matches the [0.8.0] CHANGELOG heading and the dev/v0.8.0 branch. The changelog has no Unreleased section — 0.8.0 work is recorded directly under [0.8.0].
- [x] 1.2 Implement pure metadata parsing, validation, and formatting in internal/version, and embed the root file from main.go; verify table-driven tests cover stable, pre-release-only, build-only, combined suffixes, zero components, required/null fields, invalid types/JSON/dates, SemVer identifier edge cases, and stamped build identifiers (`lb.20261107130122`, `ci.12345678.20261107130122`).
- [x] 1.3 Add the SemVer and synchronized release-update guidance to AGENTS.md (no Unreleased section; changes go under the current development version heading) and a metadata/changelog consistency test; verify the embedded version matches the latest versioned changelog heading ([0.8.0]) and that no Unreleased heading exists.

## 2. CLI version flag and MCP identity

- [x] 2.1 Add early standard-library flag handling for --version/-version, print the fully composed SemVer on the flag, and pass the composed version without build metadata (MAJOR.MINOR.PATCH plus any -preRelease) to server.New instead of the literal; verify focused tests for exact output (including the +build suffix when stamped), exit status, invalid metadata diagnostics, help, invalid arguments, and --version=false behavior.
- [x] 2.2 Add subprocess tests that run the built executable outside the source tree with missing/altered runtime version.json and unusable user configuration; verify both flag forms exit promptly with expected stdout, empty stderr, no created logs/storage, and no stdin interaction.
- [x] 2.3 Add an MCP initialization regression test for the executable; verify serverInfo.version equals its CLI version minus any +build suffix (pre-release kept when set in version.json) and normal stdio startup emits no version banner.
- [x] 2.4 Document the flag, output format, field semantics, and build-stamped date/build behavior in README.md; ensure the feature is listed under the CHANGELOG.md [0.8.0] section and verify examples against the built executable.

## 3. Build-path and integration verification

- [x] 3.1 Run make test, go test ./..., and go test -tags sqlite ./...; verify all tests pass, including MCP integration and version metadata validation.
- [x] 3.2 Build via go build main.go and the root package into temporary output paths, and run make build; verify all report version 0.8.0, that make build stamps a lb.{YYYYMMDDHHMMSS} build identifier and build-time date while a bare go build reports the committed baseline, and that the SQLite-enabled Makefile build succeeds.
- [x] 3.3 Run goreleaser check and a non-publishing single-target snapshot build; verify its executable embeds the same JSON-derived version components and stamps ci.{pipeline-id}.{YYYYMMDDHHMMSS} (or lb.{YYYYMMDDHHMMSS} when no pipeline ID is available) via ldflags, with no runtime version.json dependency.
