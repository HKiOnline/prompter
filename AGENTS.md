<!-- OPENSPEC:START -->
# OpenSpec Instructions

These instructions are for AI assistants working in this project.

Always open `@/openspec/AGENTS.md` when the request:
- Mentions planning or proposals (words like proposal, spec, change, plan)
- Introduces new capabilities, breaking changes, architecture shifts, or big performance/security work
- Sounds ambiguous and you need the authoritative spec before coding

Use `@/openspec/AGENTS.md` to learn:
- How to create and apply change proposals
- Spec format and conventions
- Project structure and guidelines

Keep this managed block so 'openspec update' can refresh the instructions.

<!-- OPENSPEC:END -->

# Prompter MCP Server for AI Agents

This document provides instructions for AI agents on how to interact with and contribute to the Prompter MCP server project.

## Project Overview

The Prompter MCP server is a Go application that provides a simple and efficient way to store and manage prompts for language models. The server implements the Model Context Protocol (MCP) and communicates with clients using JSON-RPC over stdio.

The server's core functionality includes:

*   **Prompt Management:** Creating, reading, updating, and deleting prompts.
*   **Prompt Storage:** Storing prompts in a file system-based database.
*   **Go Templating:** Using Go's templating engine to create dynamic prompts.

## Getting Started

To get started with the Prompter MCP server, you will need to have Go 1.24 or later installed on your system.

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

## Interacting with the Server

To interact with the server, you will need to use a JSON-RPC client that supports stdio. You can find an example of how to interact with the server in the `tests` directory.

## Extending the Server

The Prompter MCP server is designed to be extensible. You can add new functionality to the server by creating new storage providers or RPC methods.

To create a new storage provider, you will need to implement the `Provider` interface in the `internal/promptsdb` directory.

To create a new RPC method, you will need to add a new function to the `internal/rpc` directory and register it with the RPC server.
