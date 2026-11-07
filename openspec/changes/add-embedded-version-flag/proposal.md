# Proposal

## Why

Users need a quick way to verify installed GitHub Release and Homebrew binaries without starting the MCP server. The runtime currently reports a hardcoded 0.5.0 while the latest tagged release is 0.7.0, demonstrating the need for one embedded version source and explicit release synchronization rules. The changelog has been reconciled against Git tags and release commits.

## What Changes

- Add `--version` (and Go's equivalent `-version`) to print semantic version information (including any `+build` suffix) and the build date, then exit before configuration or storage initialization.
- Embed a root `version.json` containing major, minor, patch, preRelease, build, and date fields; compose SemVer from these fields. `date` is the build time (ISO 8601 UTC, `YYYY-MM-DDTHH:mm:ssZ`) and `build` holds dot-separated SemVer build identifiers stamped per build: `lb.{YYYYMMDDHHMMSS}` for local builds, `ci.{pipeline-id}.{YYYYMMDDHHMMSS}` in CI when a pipeline ID is available (otherwise `lb.`).
- Report the composed version without build metadata in MCP initialization responses — MAJOR.MINOR.PATCH, plus `-preRelease` when set in version.json; the CLI flag keeps the full composed SemVer including `+build`.
- Document Semantic Versioning and coordinated version.json/CHANGELOG.md updates in AGENTS.md — no Unreleased section; changes are recorded directly under the current development version heading — with user-facing flag documentation in README.md.
- Initialize metadata as 0.8.0, the version under development on `dev/v0.8.0`, with `date` set to the build-time baseline; document this new feature under the `[0.8.0]` changelog section. When 0.8.0 ships on main, the heading date is confirmed and the `v0.8.0` tag is cut.

## Capabilities

### New Capabilities

- `version-reporting`: Embedded version metadata, side-effect-free CLI version reporting, and consistent MCP version identity.

### Modified Capabilities

None. Existing transport startup requirements remain applicable to normal invocation.

## Impact

Touches main.go, a new root version.json, version parsing/formatting tests, MCP integration tests, README.md, CHANGELOG.md, AGENTS.md, plus build stamping of `date`/`build` in the Makefile and .goreleaser.yaml via `-ldflags -X` over the embedded defaults. Uses standard-library embedding, JSON, and flag parsing without a CLI framework. Existing `go build main.go`, Makefile, and GoReleaser builds must embed the same metadata. Automated tag generation and release publication are outside this change.
