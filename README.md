# Prompter MCP Server

Simple Model Context Protocol Server for storing and calling your custom prompts, powered by Go's own templating language.

***Note***: *This is still very much work in progress. See the dev-branches for the latest development.*

## Project Goals

- Simple Model Context Protocol Server using the official [Go MCP SDK](https://github.com/modelcontextprotocol/go-sdk).
- Heavy focus on prompts; tools and resources will be used in supporting role
- Prompts are by default stored in text files
- Prompt files can use Go's templating functionality

More information where the project is going can be found from the [roadmap](docs/roadmap.md).

## Installation

### macOS: Homebrew (recommended)

Starting with the 0.8.0 release, install the native Apple Silicon or Intel binary with:

```bash
brew install --cask hkionline/tap/prompter
prompter --version
```

To upgrade:

```bash
brew update
brew upgrade --cask hkionline/tap/prompter
```

The macOS binaries are currently unsigned. The Homebrew cask removes quarantine from the installed executable automatically. Prereleases are available as direct downloads and do not update the stable tap.

### GitHub release downloads

Download a prebuilt archive from [GitHub Releases](https://github.com/hkionline/prompter/releases). No Go installation is needed. Release binaries include SQLite support; filesystem storage remains the default.

| Platform | Architectures | Archive |
|---|---|---|
| macOS (`darwin`) | `arm64` (Apple Silicon), `amd64` (Intel) | `.tar.gz` |
| Linux | `arm64`, `amd64`, `386` (32-bit x86) | `.tar.gz` |
| Windows | `arm64`, `amd64`, `386` (32-bit x86) | `.zip` |

Extract `prompter` (or `prompter.exe`) into a directory on your PATH. GitHub's automatic source archives do not contain prebuilt binaries. Verify your download against the release's SHA-256 checksum file, for example with `shasum -a 256 <archive>` on macOS, `sha256sum <archive>` on Linux, or `Get-FileHash <archive> -Algorithm SHA256` in PowerShell.

### Build from source

Source builds require **Go 1.26.0 or later** (see `go.mod`). To install with SQLite enabled:

```bash
go install -tags sqlite github.com/hkionline/prompter@latest
```

From a clone, `make build` produces `bin/prompter` with SQLite and build metadata. Omit `-tags sqlite` for a filesystem-only `go install` build.

After installation, configure your MCP client to launch Prompter. On macOS/Linux, run `command -v prompter` and copy the absolute path into the client configuration. On Windows, use `(Get-Command prompter.exe).Source` in PowerShell. Desktop clients may not inherit your shell's PATH. See [MCP client configuration](docs/client-configuration.md).

## Version Information

To check the version of an installed binary without starting the MCP server, use the version flag:

```bash
prompter --version
# prompter 0.8.0+lb.20261107132107 (2026-11-07T13:21:07Z)
```

`-version` works the same way. The flag prints one line — `prompter <semver> (<build date>)` — and exits successfully before reading configuration, creating logs or storage, or waiting for stdin, so it works even with a broken configuration.

The version comes from the root `version.json`, which is embedded into the binary at build time:

```json
{
  "major": 0,
  "minor": 8,
  "patch": 0,
  "preRelease": "",
  "build": "",
  "date": "2026-11-07T13:21:07Z"
}
```

- `major`/`minor`/`patch` are the SemVer version components; `preRelease` and `build` are optional SemVer suffixes (empty means absent).
- `date` is the **build time** in ISO 8601 UTC (`YYYY-MM-DDTHH:mm:ssZ`).
- `build` identifies the build with dot-separated SemVer build identifiers, stamped at build time: `lb.{YYYYMMDDHHMMSS}` for local builds (e.g. `make build`), `ci.{pipeline-id}.{YYYYMMDDHHMMSS}` for CI builds when a pipeline ID is available. A plain `go build` reports the committed baseline without a build stamp.

The MCP `serverInfo.version` reported during initialization is the composed version **without** build metadata — `MAJOR.MINOR.PATCH`, plus `-preRelease` when set in `version.json` (e.g. `0.8.0`).

## Prompter Configuration

Prompter has sane default values which means you do not need to add configuration for Prompter. However, if you wish to change those defaults, you can do it via a YAML-file called prompter.yaml. Prompter looks a configuration file from the following path:

```bash
~/.config/prompter/prompter.yaml
```

Default configuration and values are as follows:

```yaml
# Prompter MCP Server Configuration
prompter:
  # The method transporting MCP's JSON-RCP calls
  transport:
    type: "stdio" #set to streamable_http if you want that instead

    # Streamable HTTP specific configuration
    streamable_http:
      port: 8080
  
  # The storage where prompts are kept
  storage:
    provider: "filesystem"  # or "sqlite" for SQLite storage
    
    # Filesystem specific configurations
    filesystem:
        prompts_directory: "~/.config/prompter/prompts"
    
    # SQLite specific configurations (used when provider is "sqlite")
    sqlite:
      path: "~/.config/prompter/prompts.db"
      timeout: 30            # connection timeout in seconds
      foreign_keys: true     # enable foreign key constraints
      max_open_conns: 10     # max open database connections
      max_idle_conns: 5      # max idle database connections
```

*Note:* By default, the filesystem storage provider is used. If there is no *~/.config/prompter/prompts* directory, it will be created. If you wish to change the location of the prompts files, you need to define it in the prompter.yaml file.

### SQLite Storage Provider

Prompter supports SQLite as an alternative storage backend. Homebrew and release binaries already include it; set `storage.provider` to `"sqlite"` in your `prompter.yaml`. For a manual source build, use the `sqlite` build tag:

```bash
go build -tags sqlite -o prompter main.go
```

Without the `sqlite` build tag, selecting SQLite returns an explicit error asking you to rebuild with SQLite support.

**Environment variable overrides:**

| Variable | Description |
|---|---|
| `PROMPTER_STORAGE_PROVIDER` | Override storage provider (`filesystem` or `sqlite`) |
| `PROMPTER_SQLITE_PATH` | Override SQLite database file path |
| `PROMPTER_SQLITE_TIMEOUT` | Override SQLite connection timeout (seconds) |
| `PROMPTER_SQLITE_FOREIGN_KEYS` | Set to `true` to enable foreign keys |
| `PROMPTER_LOG_FILE` | Override log file path |

**Migrating from filesystem to SQLite:**

Use the `migrate` command to copy existing prompts from a filesystem directory into a SQLite database:

```bash
go build -tags sqlite -o prompter-migrate ./cmd/migrate
./prompter-migrate -src ~/.config/prompter/prompts -dst ~/.config/prompter/prompts.db
```

Flags: `--dry-run` (validate only), `--incremental` (skip existing), `--verbose` (detailed progress).

## MCP Client Configuration

By default the MCP hosts manage [clients](docs/client-configuration.md) which in turn manage the lifecycle of MCP server communication. The MCP servers are typically configured in the hosts own configuration, typically called *mcp.json*.

For example in **[OpenCode](https://opencode.ai)**, configuration for stdio based prompter MCP-server would be as follows:

```json

{
  "$schema": "https://opencode.ai/config.json",

  "mcp": {
	"prompter": {
		"type": "local",
		"enabled": true,
		"command": ["/opt/homebrew/bin/prompter"]
	}
  }

}
```

Replace the example Apple Silicon Homebrew path with the absolute path printed by `command -v prompter` on your machine. Intel Homebrew installations typically use `/usr/local/bin/prompter`; source installations usually use `$(go env GOPATH)/bin/prompter`. Do not put unevaluated shell expressions into the client configuration.

Prompter is not limited to OpenCode. It supports any MCP client/host that supports MCP prompts and tools capabilities. Additional client configuration examples are available in the **[client configuration documentation](docs/client-configuration.md)**.

## Project Status

See the projects [roadmap](docs/roadmap.md).

## Project Motivations

- Learn to implement a basic MCP server from scratch (done in v0.1.0, 0.2.0 forwards uses [official SDK](https://github.com/modelcontextprotocol/go-sdk))
- Learn about JSON RPC
- Opportunity to use Go's templating system creatively

## Note on AI and tools

Since this project is an AI tool, it seems almost natural that at least parts of it would be coded using AI assistance. Initial v0.1.0 which already included functional MCP server with support listing prompts and getting stored prompt as well as one tool to save a new prompt was "handcrafted" without any AI assistance. This phase's motivation was to get to know the [protocol](https://modelcontextprotocol.io/) itself intimately.

From v0.2.0 onwards the coding has been AI-assisted. Not blindly "vibe coded" but carefully managed and curated. To bring structure to AI-assisted coding and maintain a more persistent context, I use [OpenSpec](https://openspec.dev/) for light spec driven AI-development. The coding has been done using [Neovim](https://neovim.io/) when done by hand and [OpenCode](https://opencode.ai) when agentic AI-coding has been employed. Much of the AI-generated code has been the work of Mistral's Devstral 2 Small and Regular model. The small model has been ran locally using [LM Studio](https://lmstudio.ai/). The Devstral 2 (regular) has been used via OpenCode and Mistral's own [Vibe CLI](https://mistral.ai/news/devstral-2-vibe-cli) on Mistral's servers.


**Tools list**:

- [Neovim](https://neovim.io/) - for all human friendly editing needs
- [OpenCode](https://opencode.ai) - for agentic coding tasks
- [Mistral Vibe CLI](https://mistral.ai/news/devstral-2-vibe-cli) - for agentic coding tasks
- [LM Studio](https://lmstudio.ai/) - for running models locally
- [Mistral Devstral 2](https://huggingface.co/mistralai/Devstral-2-123B-Instruct-2512) - Mistral's code focused LLM
- [Mistral Devstral 2 Small](https://huggingface.co/mistralai/Devstral-Small-2-24B-Instruct-2512) - Mistral's smaller code focused LLM, runnable locally
- [OpenSpec](https://openspec.dev/) - for giving LLM's persistent context and structured guidelines on implementation
