# Version Reporting Specification

## Purpose

Provide a self-contained, consistent version identity for installed Prompter binaries, accessible to users through the CLI and to clients through MCP initialization.

## Requirements

### Requirement: Embedded version metadata
The binary SHALL embed the root version.json at build time. The file SHALL contain nonnegative integer major, minor, and patch fields and string preRelease, build, and date fields. Empty preRelease and build strings SHALL indicate absent SemVer suffixes. The date SHALL be a valid ISO 8601 UTC build timestamp in the format YYYY-MM-DDTHH:mm:ssZ, separate from the semantic version; build tooling SHALL stamp it with the actual build time.

#### Scenario: Self-contained executable
- **WHEN** a built binary is run from a directory without version.json or beside an altered version.json
- **THEN** it reports the metadata embedded during its build rather than reading a runtime file

#### Scenario: Invalid metadata
- **WHEN** embedded metadata contains malformed JSON, missing required fields, invalid field types, a negative version component, an invalid SemVer identifier, or an invalid date
- **THEN** the executable exits nonzero with a descriptive error on stderr before configuration, storage, or transport initialization
- **AND** it emits no version output or protocol data on stdout

### Requirement: Semantic version composition
The system SHALL compose major.minor.patch[-preRelease][+build] according to SemVer 2.0.0, without a leading v. Pre-release and build identifiers SHALL be nonempty dot-separated ASCII alphanumeric/hyphen identifiers when present. Numeric pre-release identifiers SHALL NOT have leading zeroes; numeric build identifiers MAY have leading zeroes.

#### Scenario: Stable version
- **WHEN** metadata specifies 1, 2, 3 and empty preRelease and build strings
- **THEN** the composed version is `1.2.3`

#### Scenario: Pre-release and build metadata
- **WHEN** metadata specifies 1, 2, 3, preRelease `rc.1`, and build `abc.001`
- **THEN** the composed version is `1.2.3-rc.1+abc.001`

#### Scenario: Independent optional suffixes
- **WHEN** only preRelease or only build is nonempty
- **THEN** the output includes only that suffix with its respective `-` or `+` separator

### Requirement: Build identification
Local builds SHALL stamp the build field as `lb.{YYYYMMDDHHMMSS}` (UTC) and the date field as the build timestamp. CI builds SHALL stamp the build field as `ci.{pipeline-id}.{YYYYMMDDHHMMSS}` when a pipeline identifier is available from the environment (for example GITHUB_RUN_ID or CI_PIPELINE_ID), and SHALL fall back to the `lb.` form otherwise. Each dot-separated build identifier SHALL be an alphanumeric/hyphen identifier or a numeric identifier; numeric build identifiers MAY have leading zeroes.

#### Scenario: Local build identifier
- **WHEN** a binary is built locally via make build
- **THEN** its --version output contains `+lb.` followed by the UTC build timestamp in YYYYMMDDHHMMSS form

#### Scenario: CI build identifier
- **WHEN** a binary is built in CI with a pipeline ID available in the environment
- **THEN** its --version output contains `+ci.{pipeline-id}.{YYYYMMDDHHMMSS}`
- **AND WHEN** no pipeline ID is available
- **THEN** it falls back to `+lb.{YYYYMMDDHHMMSS}`

### Requirement: Version flag output and early exit
Invoking `prompter --version` or `prompter -version` SHALL print exactly `prompter <semver> (<date>)` followed by a newline to stdout and exit successfully, where `<semver>` is the composed version including any `+build` suffix and `<date>` is the ISO 8601 build timestamp. It SHALL NOT read user configuration, create logs or storage, start a transport, or wait for stdin. Successful version output SHALL leave stderr empty.

#### Scenario: Version inspection with unusable configuration
- **WHEN** either version flag is invoked with invalid user configuration and no available MCP client
- **THEN** the executable promptly prints the embedded version and exits with status 0 without accessing the configuration or creating files

### Requirement: Consistent MCP identity
Normal no-argument startup SHALL continue to run the configured MCP transport without a version banner on stdout. MCP initialization SHALL report the composed semantic version without build metadata as serverInfo.version: MAJOR.MINOR.PATCH, plus the `-preRelease` suffix when preRelease is nonempty in version.json. The `+build` suffix SHALL NOT appear in serverInfo.version.

#### Scenario: MCP version excludes the build stamp
- **WHEN** an MCP client initializes a build-stamped binary whose --version prints `prompter 0.8.0+lb.20261107130122 (...)`
- **THEN** serverInfo.version is `0.8.0`
- **AND** stdout contains only protocol output for stdio transport

#### Scenario: MCP version keeps a declared pre-release
- **WHEN** version.json sets preRelease `rc.1` and the binary is build-stamped, so --version prints `prompter 0.9.0-rc.1+ci.12345678.20261107130122 (...)`
- **THEN** serverInfo.version is `0.9.0-rc.1`

### Requirement: Release metadata synchronization guidance
AGENTS.md SHALL specify SemVer 2.0.0 and require version.json and CHANGELOG.md to be updated together when a release is cut. The composed version without build metadata SHALL match the corresponding changelog heading. The changelog SHALL NOT use an Unreleased section: development changes SHALL be recorded directly under the heading of the version under development (currently [0.8.0]), whose date is confirmed when that release ships. Release tags SHALL use v followed by the composed version.

#### Scenario: Preparing a release
- **WHEN** a contributor follows the release guidance
- **THEN** they choose the appropriate semantic version increment, update version.json and confirm the matching changelog heading and its date, and use the matching version tag

#### Scenario: Development changes
- **WHEN** a contributor adds a feature on a development branch
- **THEN** guidance directs them to document it under the current development version heading (for example [0.8.0]) without changing historical release entries
