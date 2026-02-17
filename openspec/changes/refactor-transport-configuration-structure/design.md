## Context

The current configuration mixes transport concerns across unrelated top-level keys: a scalar `transport` selector and a separate `http` object for streamable HTTP options. This differs from the provider-style nesting already used for storage and makes transport-specific options harder to discover and evolve.

The change introduces a structured transport object so selection (`type`) and transport-specific settings (`streamable_http`) live in one namespace. The implementation touches configuration defaults/loading/validation, server transport startup, and user-facing examples/tests.

## Goals / Non-Goals

**Goals:**
- Align transport configuration shape with existing provider-style patterns.
- Make streamable HTTP options discoverable under transport scope.
- Preserve runtime behavior for valid configurations while updating schema and tests.
- Keep implementation small and localized to configuration and server startup paths.

**Non-Goals:**
- Adding new transport types.
- Changing streamable HTTP protocol behavior.
- Introducing automatic migration tooling for old config files.

## Decisions

### 1) Transport becomes a structured object
Use:
- `prompter.transport.type`
- `prompter.transport.streamable_http.port`

Rationale: mirrors storage provider configuration style and provides a stable namespace for future transport-specific options.

Alternative considered: keep `transport` scalar and only rename top-level `http` to `streamable_http`. Rejected because it keeps transport concerns split across top-level keys.

### 2) Rename HTTP config struct and YAML key to streamable HTTP naming
Use `StreamableHTTPConfiguration` and YAML key `streamable_http` under `transport`.

Rationale: explicit naming aligned with transport type string (`streamable_http`) avoids ambiguous `http` references.

Alternative considered: keep struct/key as `HTTP` while nesting under transport. Rejected to avoid mixed terminology and future confusion.

### 3) Keep validation strict on accepted transport types
Validation remains constrained to `stdio` and `streamable_http`, now applied to `transport.type`.

Rationale: fail fast on invalid config and preserve current safety behavior.

Alternative considered: permissive validation with runtime fallback. Rejected because silent fallback could mask operator errors.

### 4) Treat schema as breaking without compatibility fallback
Do not support legacy top-level `http` and scalar `transport` in the same change.

Rationale: reduces complexity and avoids dual-shape parsing logic in an early-stage project.

Alternative considered: temporary backward-compat parser. Rejected to keep configuration semantics explicit and avoid extra maintenance.

## Risks / Trade-offs

- [Breaking user configs] → Document migration clearly in examples and release notes.
- [Missed test fixture updates] → Update configuration and server tests that build config literals or YAML snippets.
- [Terminology drift in docs] → Update README and example YAML together with code.

## Migration Plan

1. Update config structs/defaults and parser validation.
2. Update server code to consume nested transport fields.
3. Update tests and example configuration files.
4. Validate with `make test`.

Rollback strategy: revert this change set to restore old schema.

## Open Questions

- None currently; transport-specific options beyond port are intentionally deferred.
