
# Prompter MCP Server for AI Agents

This document provides instructions for AI agents on how to interact with and contribute to the Prompter MCP server project. For detailed thematic documentation see the docs-directory.

## Project Overview

The Prompter MCP server is a Go application that provides a simple and efficient way to store and manage prompts for language models. The server implements the Model Context Protocol (MCP) and communicates with clients using JSON-RPC over stdio.

The server's core functionality includes:

*   **Prompt Management:** Creating, reading, updating, and deleting prompts.
*   **Prompt Storage:** Storing prompts in a file system-based database.
*   **Go Templating:** Using Go's templating engine to create dynamic prompts.

## Getting Started

To get started with the Prompter MCP server, you will need to have Go 1.26.0 or later installed on your system (see `go.mod`).

To build the project, clone the repository and run the following command:

```bash
go build main.go
```

This will create an executable file named `main` in the project's root directory.

To run the server, simply execute the following command:

```bash
./main
```

The server will then start listening for incoming connections on stdin.

## Core Concepts

### Prompts

A prompt is a piece of text that is used to instruct a language model to perform a specific task. In the Prompter MCP server, a prompt is represented by a YAML file with the following structure:

```yaml
name: my-prompt
title: My Prompt
description: A simple prompt that demonstrates how to use the Prompter MCP server.
arguments:
  - name
  - age
tags:
  - example
  - simple
---
Hello, my name is {{.name}} and I am {{.age}} years old.
```

The prompt's metadata is defined in the YAML front matter, and the prompt's content is defined in the body of the file.

### Storage Providers

The Prompter MCP server uses a storage provider to store and manage prompts. The default storage provider is the file system provider, which stores prompts in a directory on the local file system.

You can configure the storage provider in the `prompter.yaml` file.

### Communication Protocol

The Prompter MCP server communicates with clients using the JSON-RPC protocol over stdio. The server implements the following RPC methods:

*   `initialize`: Initializes the server.
*   `ping`: Pings the server to check if it is alive.
*   `prompts/list`: Lists all available prompts.
*   `prompts/get`: Gets a specific prompt.
*   `tools/saveNewPrompt`: Saves a new prompt.

## Versioning and Releases

The project adheres to [Semantic Versioning 2.0.0](https://semver.org/spec/v2.0.0.html). The single source of version truth is the root `version.json`, embedded into the binary at build time:

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

*   `major`/`minor`/`patch` are nonnegative integers; `preRelease`/`build`/`date` are strings. All six fields are required. Empty `preRelease`/`build` mean the suffix is absent.
*   `date` is the **build time** in ISO 8601 UTC (`YYYY-MM-DDTHH:mm:ssZ`), stamped by the build tooling (Makefile/GoReleaser via `-ldflags -X`); a bare `go build main.go` reports the committed baseline.
*   `build` is dot-separated SemVer build metadata stamped per build: `lb.{YYYYMMDDHHMMSS}` for local builds, `ci.{pipeline-id}.{YYYYMMDDHHMMSS}` in CI when a pipeline ID is available.
*   `prompter --version` prints `prompter <semver> (<date>)` with the full composed SemVer; MCP `serverInfo.version` reports the composed version without build metadata (MAJOR.MINOR.PATCH plus `-preRelease` when set).

**Changelog rules:** `CHANGELOG.md` has **no Unreleased section**. On a development branch, record changes directly under the heading of the version under development (currently `## [0.8.0] - <date>`). When a release is cut on `main`, confirm that heading's date, update `version.json` so its version matches the heading, and tag the release as `v` followed by the composed version (e.g. `v0.8.0`). A consistency test in `internal/version` enforces that the embedded version matches the latest changelog heading and that no Unreleased heading exists.

Release preparation and non-publishing validation commands are documented in `docs/releases.md`. The tag workflow validates the tag against committed version metadata and the latest changelog heading, tests both storage build variants, and passes curated release notes to GoReleaser before publishing.

## Testing

The project uses `make test` as the default command for running all tests. This command:
- Builds the project
- Runs unit tests for all packages
- Executes integration tests including the test server
- Ensures all JSON-RPC responses are error-free

To run tests, execute:
```bash
make test
```

All unit tests should pass and the test server's JSON output should not contain any error messages.

## Interacting with the Server

To interact with the server, you will need to use a JSON-RPC client that supports stdio. You can find an example of how to interact with the server in the `tests` directory.

## Extending the Server

The Prompter MCP server is designed to be extensible. You can add new functionality to the server by creating new storage providers or RPC methods.

To create a new storage provider, you will need to implement the `Provider` interface in the `internal/promptsdb` directory.

To create a new RPC method, you will need to add a new function to the `internal/rpc` directory and register it with the RPC server.
