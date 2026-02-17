## Why

Transport configuration is currently split across a top-level scalar (`transport`) and a separate top-level object (`http`), which makes the structure inconsistent with provider-style configuration patterns already used elsewhere. Aligning transport with the same object pattern improves readability, extensibility, and reduces ambiguity as transport-specific options grow.

## What Changes

- Refactor configuration schema so `transport` is an object with a required `type` field.
- Move transport-specific HTTP options under `transport.streamable_http`.
- Rename top-level `http` configuration to `streamable_http` under `transport`.
- Update defaults, configuration loading, validation, and runtime transport selection to use the nested structure.
- Update examples and tests to reflect the new schema.
- **BREAKING**: Existing configs using top-level `transport: <string>` and top-level `http:` must migrate to nested `transport.type` and `transport.streamable_http`.

## Capabilities

### New Capabilities
- `transport-configuration`: Structured transport configuration with `type` and transport-specific nested options.

### Modified Capabilities
- None.

## Impact

- Affected code: `internal/configuration`, `internal/server`, config examples, and tests.
- APIs/config surface: YAML schema for runtime configuration changes in a backward-incompatible way.
- Dependencies: none.
- Systems: startup configuration parsing/validation and transport bootstrapping behavior.
