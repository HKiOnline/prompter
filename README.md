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

Currently installing prompter requires Go 1.24 or later to be installed on your system. Prompter supports all OS and architetures wich you can compile it to. Development of prompter is done in MacOS so some "unixy" emphaises might exist.

To install clone the repo and build the project. The project roots main.go should be the target.

```bash
  go install github.com/hkionline/prompter@latest
```

## Configuration

Prompter has sane default values which means you do not need to add configuration for Prompter. However, if you wish to change those defaults, you can do it via a YAML-file called prompter.yaml. Prompter looks a configuration file from the following path:

```bash
~/.config/prompter/prompter.yaml
```

Default configuration and values are as follows:

```yaml
# Prompter MCP Server Configuration
prompter:
  # The method transporting MCP's JSON-RCP calls
  transport: "stdio"
  
  # The storage where prompts are kept
  storage:
    provider: "filesystem"
    
    # Filesystem specific configurations
    filesystem:
        prompts_directory: "~/.config/prompter/prompts"
```

*Note:* By default, the filesystem storage provider is used. If there is no *~/.config/prompter/prompts* directory, it will be created. If you wish to change the location of the prompts files, you need to define it in the prompter.yaml file.

## Usage

By default the MCP hosts manage clients which in turn manage the lifecycle of MCP server communication. The MCP servers are typically configured in the hosts own configuration, typically called *mcp.json*.

For example in **[OpenCode](https://opencode.ai)**, configuration for stdio based prompter MCP-server would be as follows:

```json

{
  "$schema": "https://opencode.ai/config.json",

  "mcp": {
	"prompter": {
		"type": "local",
		"enabled": true,
		"command": ["/Users/hki/Developer/Tools/Go/bin/prompter"]
	}
  }

}
```

Note that the path to the command depends on where prompter is installed.

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

