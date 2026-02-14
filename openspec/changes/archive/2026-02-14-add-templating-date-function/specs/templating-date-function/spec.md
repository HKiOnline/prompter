## ADDED Requirements

### Requirement: Date function availability
The templating system SHALL provide a built-in `date` function that returns the current date.

#### Scenario: Date function exists
- **WHEN** a template contains `{{date}}`
- **THEN** the template execution succeeds
- **AND** the `date` function is called

#### Scenario: Date function returns correct format
- **WHEN** a template contains `{{date}}`
- **THEN** the rendered output contains a date string in YYYY-MM-DD format
- **AND** the date represents the current day

### Requirement: Date function format
The `date` function SHALL return dates in YYYY-MM-DD format (e.g., "2024-01-15").

#### Scenario: Consistent date format
- **WHEN** a template calls the `date` function multiple times
- **THEN** all calls return dates in the same YYYY-MM-DD format
- **AND** the format uses 4-digit year, 2-digit month, and 2-digit day

### Requirement: Date function safety
The `date` function SHALL not expose system information beyond the current date.

#### Scenario: No system information leakage
- **WHEN** a template calls the `date` function
- **THEN** the output contains only the date string
- **AND** no system paths, environment variables, or other sensitive information is exposed

### Requirement: Date function availability in all templates
The `date` function SHALL be available in all prompt templates without additional configuration.

#### Scenario: Function available in new prompts
- **WHEN** a new prompt is created with template syntax
- **THEN** the `date` function is available for use

#### Scenario: Function available in existing prompts
- **WHEN** an existing prompt is updated to use template syntax
- **THEN** the `date` function is available for use