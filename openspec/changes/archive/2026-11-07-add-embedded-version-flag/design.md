# Design

## Context

See proposal.md for motivation. main.go currently ignores CLI arguments and initializes configuration, logging, and storage before calling server.New with a literal `0.5.0`. internal/server/server.go already passes its constructor version to the MCP SDK, so that API can remain unchanged. Makefile and AGENTS.md use `go build main.go`, while GoReleaser builds the root package. Git tag v0.7.0 points to 2ac02ca dated 2026-11-07; CHANGELOG.md has been corrected to reflect it. Development is on dev/v0.8.0, and the changelog records 0.8.0 work directly under the `[0.8.0]` heading — there is no Unreleased section. The `date` field is interpreted as build time and is stamped by the build tooling.

## Goals / Non-Goals

**Goals:** Keep one reproducible metadata source across build paths; make version inspection independent of runtime setup; preserve clean stdio protocol output; define metadata semantics and contributor release rules.

**Non-Goals:** Introduce a subcommand framework, change storage defaults, version the migration utility, publish a release, or automate tag/changelog editing. Build-time stamping applies only to `date` and `build`; the version components stay committed and static.

## Decisions

### Use a boolean flag through the standard library

Use `flag` for a boolean `version` option; it accepts both `--version` and `-version`. A flag fits this single-purpose daemon with no existing server subcommands and needs no new dependency. Go does not prescribe one universal version interface; a subcommand would also be valid but would add unnecessary dispatch here. Keep normal no-argument startup unchanged. Use standard flag help/error behavior; reject unexpected positional arguments with a nonzero exit and a stderr diagnostic rather than treating them as commands. Do not add a `-v` alias, leaving it available for future verbosity.

Parse flags and validate embedded metadata before user configuration or logging. On a true version flag print one line, `prompter <semver> (<date>)`, and return. On invalid metadata, print a descriptive error to stderr and exit nonzero without startup side effects. `--version=false` proceeds to ordinary startup.

### Embed the root file directly in main.go

Place `//go:embed version.json` and its string/byte variable in main.go with the embed import. This preserves existing single-file builds: placing the embed declaration only in a sibling root Go file would exclude it from `go build main.go`. An internal/version package can own pure parsing, validation, and formatting and be imported by main.go; it must receive the embedded bytes rather than trying to embed `../../version.json`, which Go disallows.

Parse once and derive two forms: the CLI prints the fully composed SemVer (with `-preRelease`/`+build` when present), while server.New receives the composed version without build metadata — MAJOR.MINOR.PATCH plus `-preRelease` when set — for the MCP serverInfo identity. The version components (major/minor/patch/preRelease) come solely from the embedded version.json — no git discovery or other competing source. `date` and `build` are build-stamped fields: the Makefile and GoReleaser pass them via `-ldflags -X` into internal/version, overriding the embedded defaults, while a bare `go build main.go` reports the committed baseline. Local builds stamp `build` as `lb.{YYYYMMDDHHMMSS}` (UTC) and `date` as the build timestamp; CI builds stamp `ci.{pipeline-id}.{YYYYMMDDHHMMSS}` when a pipeline ID is available from the environment (e.g. GITHUB_RUN_ID, CI_PIPELINE_ID), falling back to `lb.`. GoReleaser snapshots therefore identify their CI pipeline; release preparation must still synchronize the tag with the committed version components.

### Define all six fields explicitly

Initial proposed file:

```json
{
  "major": 0,
  "minor": 8,
  "patch": 0,
  "preRelease": "",
  "build": "",
  "date": "2026-11-07T13:01:22Z"
}
```

All fields are required. Major/minor/patch are nonnegative integers, not numeric strings; preRelease/build/date are strings. Distinguish absent/null fields from legitimate zero/empty values while decoding. Enforce SemVer identifier rules, including the numeric pre-release leading-zero restriction and allowance of leading zeroes in build identifiers. Validate date as an actual ISO 8601 UTC timestamp, specifically the format YYYY-MM-DDTHH:mm:ssZ using the standard library. `build` is a string of dot-separated SemVer build identifiers (alphanumeric/hyphen identifiers, or numeric identifiers, which may have leading zeroes), not an automatically incremented number. `date` means build time, not release date.

The baseline is 0.8.0, the version under development on dev/v0.8.0, and its composed version must match the `[0.8.0]` CHANGELOG heading; the changelog has no Unreleased section. The committed file keeps an empty `build` and the initialization timestamp as `date`; Makefile and GoReleaser builds stamp both at build time. README must explain that `date`/`build` are build-time stamps identifying when and where a binary was built, not a guarantee that the source matches a release tag.

### Contributor synchronization rules

Add a focused versioning section to AGENTS.md describing SemVer 2.0.0, coordinated JSON/changelog release updates, optional suffixes, and matching `v<semver>` tags. There is no Unreleased section: on a development branch, changes are recorded directly under the heading of the version under development (currently `[0.8.0]`), and that heading's date is confirmed when the release is cut. Version components in version.json and the changelog heading are updated together. No automatic version bump or publication belongs to this implementation.

Add a small consistency test that checks the embedded version components against the latest versioned changelog heading and rejects an Unreleased heading, catching the drift already observed without a new release management tool. `date` is build time, so the test compares version components only, not dates.

## Risks / Trade-offs

- [Version text could corrupt stdio MCP] → Only emit the human-readable line in explicit version mode; test normal initialization for protocol-only output.
- [Bad committed metadata breaks startup] → Validate through focused unit tests and a test of the actual embedded file before distributing binaries.
- [Single-file build excludes helper source] → Keep embedding and entry-point wiring in main.go; place reusable helpers in an imported package and verify both build forms.
- [Release tag can disagree with committed metadata] → Document the required synchronization; automatic CI tag enforcement is deferred.
- [Bare `go build` shows the committed baseline date/build] → Stamp `date`/`build` in the Makefile and GoReleaser so `make build` and released binaries are identifiable; document that an unstamped `go build main.go` reports the committed baseline.
- [Build-stamped date differs from the changelog heading date] → By design: `date` is build time; the consistency test compares version components only.

## Migration Plan

Add metadata and parser, wire early CLI handling and MCP identity, add Makefile/GoReleaser build stamping, then add documentation and tests. No user configuration or data migration is required. Verify both filesystem-only and SQLite-enabled builds, single-file builds, and Makefile builds. Rollback consists of reverting this change; storage and protocol capabilities are unaffected.
