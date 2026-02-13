# Change: Change prompt file extension from YAML to MD

## Why
The current implementation uses `.yaml` file extension for prompts, but these files are largely markdown files with YAML frontmatter at the beginning. Using `.md` extension would be more accurate and intuitive for users.

## What Changes
- **BREAKING**: Change prompt file extension from `.yaml` to `.md`
- Add documentation about the prompt file format in the docs directory

## Impact
- Affected specs: prompts
- Affected code: internal/promptsdb/fsProvider.go (savePrompt function)
