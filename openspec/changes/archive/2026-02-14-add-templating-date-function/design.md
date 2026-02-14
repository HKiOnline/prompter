## Context

The templating system currently supports basic Go template syntax but lacks built-in functions for common operations like date formatting. Users need to manually provide dates or use external processing. This design adds a safe, built-in date function.

## Goals / Non-Goals

**Goals:**
- Add a `date` function that returns current date in YYYY-MM-DD format
- Ensure the function is safe and doesn't expose system information
- Make the function available to all prompt templates
- Maintain backward compatibility with existing templates

**Non-Goals:**
- Support for custom date formats (beyond YYYY-MM-DD)
- Timezone-aware date formatting
- Date arithmetic or manipulation functions
- Localization/internationalization of date formats

## Decisions

**Function Signature and Behavior:**
- Function name: `date`
- Returns: string in "2006-01-02" format (YYYY-MM-DD)
- No parameters (always uses current time)
- Uses `time.Now()` from Go standard library
- Rationale: Simple, predictable interface that covers the most common use case

**Implementation Location:**
- Add function to `internal/templa/builtins.go`
- Register in template initialization in `internal/templa/templa.go`
- Rationale: Follows existing pattern for template functions, keeps templating logic centralized

**Safety Considerations:**
- No access to environment variables or system information
- No file system or network access
- Deterministic output format
- Rationale: Maintain security posture, prevent information leakage

**Error Handling:**
- No error returns (date formatting is infallible with fixed format)
- Panic recovery in template execution
- Rationale: Simplifies template usage, errors would be programming bugs

## Risks / Trade-offs

**Performance Impact:**
- [Risk] Additional function call overhead per template execution → Minimal, date formatting is fast
- [Risk] Memory allocation for date string → Acceptable, strings are small

**Compatibility:**
- [Risk] Function name collision with future user-defined functions → Unlikely, "date" is a common name but reserved namespace for builtins
- [Risk] Breaking changes to template API → None, purely additive

**Maintenance:**
- [Risk] Additional code to maintain → Minimal, simple function with no dependencies
- [Risk] Need to support this function long-term → Acceptable, dates are fundamental requirement

**Testing:**
- [Risk] Date-dependent tests may be flaky → Use fixed time in tests via dependency injection
- [Risk] Timezone issues in testing → Use UTC consistently