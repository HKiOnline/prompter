# Prompter MCP Client Configuration

This document provides example configurations for common MCP hosts/clients when using Prompter.

All remote examples use a dummy URL (`https://your-prompter.example.com/mcp`). Replace it with your real server URL.

## Table of Contents

- [General Notes](#general-notes)
- [OpenCode](#opencode)
- [VS Code](#vs-code)
- [Claude Code CLI](#claude-code-cli)
- [Cursor](#cursor)
- [Gemini CLI](#gemini-cli)

## General Notes

### Local (stdio) vs Remote (HTTP)

- Local/stdio runs the `prompter` binary directly from the host/client.
- Remote/HTTP connects to a running Prompter server over the network.

### Prompter remote transport setup

To serve Prompter over HTTP, set transport to `streamable_http`:

```yaml
prompter:
  transport:
    type: "streamable_http"
    streamable_http:
      port: 8080
```

### Binary path

Use the real absolute path where `prompter` is installed, for example:

- `/Users/your-user/go/bin/prompter`
- `/usr/local/bin/prompter`

## OpenCode

OpenCode local example (same format as in `README.md`):

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "prompter": {
      "type": "local",
      "enabled": true,
      "command": ["/absolute/path/to/prompter"]
    }
  }
}
```

OpenCode remote example:

```json
{
  "$schema": "https://opencode.ai/config.json",
  "mcp": {
    "prompter": {
      "type": "remote",
      "enabled": true,
      "url": "https://your-prompter.example.com/mcp"
    }
  }
}
```

## VS Code

VS Code local example (`.vscode/mcp.json`):

```json
{
  "servers": {
    "prompter": {
      "type": "stdio",
      "command": "/absolute/path/to/prompter",
      "args": []
    }
  }
}
```

VS Code remote example (`.vscode/mcp.json`):

```json
{
  "servers": {
    "prompter": {
      "type": "http",
      "url": "https://your-prompter.example.com/mcp"
    }
  }
}
```

## Claude Code CLI

Claude Code local example (`.mcp.json`):

```json
{
  "mcpServers": {
    "prompter": {
      "command": "/absolute/path/to/prompter",
      "args": []
    }
  }
}
```

Claude Code remote example (`.mcp.json`):

```json
{
  "mcpServers": {
    "prompter": {
      "url": "https://your-prompter.example.com/mcp"
    }
  }
}
```

## Cursor

Cursor local example (`~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "prompter": {
      "command": "/absolute/path/to/prompter",
      "args": []
    }
  }
}
```

Cursor remote example (`~/.cursor/mcp.json`):

```json
{
  "mcpServers": {
    "prompter": {
      "url": "https://your-prompter.example.com/mcp"
    }
  }
}
```

## Gemini CLI

Gemini CLI local example (`~/.gemini/settings.json`):

```json
{
  "mcpServers": {
    "prompter": {
      "command": "/absolute/path/to/prompter",
      "args": []
    }
  }
}
```

Gemini CLI remote example (`~/.gemini/settings.json`):

```json
{
  "mcpServers": {
    "prompter": {
      "url": "https://your-prompter.example.com/mcp"
    }
  }
}
```

## Notes on client schema differences

Host/client MCP schemas can change over time. If your client version expects a different shape (for example `type`, `transport`, headers, or auth fields), keep the same local/remote idea and adapt keys to that client's latest documentation.
