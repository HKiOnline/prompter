## Why

The templating functionality needs a built-in function to provide the current date in YYYY-MM-DD format. This enables prompts to include dynamic date information without requiring external processing or manual input.

## What Changes

- Add a new built-in template function `date` that returns the current date in YYYY-MM-DD format
- Update the templating engine to register and make this function available to all prompts
- Ensure the function is safe and doesn't expose system information beyond the date

## Capabilities

### New Capabilities
- `templating-date-function`: Built-in function to provide current date in prompts

### Modified Capabilities
- None

## Impact

- Affects the templating engine in `internal/templa/` directory
- No breaking changes to existing functionality
- Adds new capability to prompt templates
- Minimal performance impact (single date formatting operation per template execution)