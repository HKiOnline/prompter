# Prompt File Format

## Overview

Prompters file system storage provider (fsProvider) uses markdown files with YAML frontmatter to store and manage prompts. This format allows for both structured metadata and rich content in a single file. The same format is also used when outputting prompts to MCP-clients.

## File Structure

A prompt file consists of two main parts:

1. **YAML Frontmatter**: Structured metadata about the prompt
2. **Content**: The actual prompt text

### Example Prompt File

```markdown
---
name: "describe_tampere"
title: "Tell me about Tampere"
description: "Prompt to get information about Tampere"
arguments:
tags:
  - tampere
  - finland
  - sample
---
Tell me about Tampere, Finland.
```

### Example with Template Functions

```markdown
---
name: "daily_report"
title: "Daily Report Template"
description: "Template for generating daily reports with current date"
arguments:
  - name
tags:
  - report
  - daily
  - template
---
Daily Report for {{date}}

Hello {{.name}},

Here's your daily report for {{date}}.
```

## YAML Frontmatter Fields

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `name` | string | Yes | Unique computer-readable identifier for the prompt (used as filename) |
| `title` | string | No | Human-readable title of the prompt |
| `description` | string | No | Detailed explanation of what the prompt does |
| `arguments` | array of strings | No | Arguments that can be passed to the prompt when invoked |
| `tags` | array of strings | No | Tags for categorization and search (used in completion suggestions) |

## Content Section

After the YAML frontmatter (separated by `---`), you can include any text content. This is where you write your actual prompt instructions.

### Features Supported in Content

- **Markdown**: Standard markdown syntax is supported
  - Headers, lists, code blocks
  - Links and images
  - Tables
  
- **Go Templates**: You can use Go template syntax for dynamic content
  ```
  Hello {{.name}}, you are {{.age}} years old.
  ```

### Built-in Template Functions

The templating system provides built-in functions that can be used in your prompt templates. These functions are automatically available and don't require any special configuration.

#### Using Built-in Functions

Built-in functions are called using the standard Go template syntax:
```
{{functionName}}
```

For functions that require arguments, use the standard Go template argument syntax:
```
{{functionName arg1 arg2}}
```

#### Available Built-in Functions

##### `date` Function

Returns the current date in YYYY-MM-DD format.

**Signature:** `date() string`

**Usage:**
```
Today is {{date}}.
```

**Example:**
```markdown
---
name: "daily_briefing"
title: "Daily Briefing"
description: "Generate a daily briefing with current date"
arguments:
  - name
tags:
  - briefing
  - daily
---
Daily Briefing - {{date}}

Hello {{.name}},

Here's your briefing for {{date}}.
```

**Output:**
```
Daily Briefing - 2024-01-15

Hello John,

Here's your briefing for 2024-01-15.
```

**Notes:**
- Always returns the current date at the time of template execution
- Format is always YYYY-MM-DD (e.g., "2024-01-15")
- Safe to use - doesn't expose any system information beyond the date
- Can be used multiple times in the same template

## File Extension

Prompt files use the `.md` extension to reflect their markdown-based format with YAML frontmatter.

## Best Practices

1. **Naming**: Use descriptive, lowercase names with hyphens for spaces (e.g., `describe-tampere.md`)
2. **Tags**: Use relevant tags to make prompts easier to discover
3. **Arguments**: Document any expected arguments in the description
4. **Content**: Be clear and specific about what you want the AI to generate
