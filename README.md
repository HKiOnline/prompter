# Prompter MCP Server

Simple Model Context Protocol Server for storing and calling your custom prompts, powered by Go's own templating language.

***Note***: *This is very much work in progress. I have decided to push to a dev branch my uncleaned, at times horrid untested code, that will eventually become the first 0.1.0 tagged and cleaned release. See project milestones.*

## Project goals

- Simple Model Context Protocol (MCP) Server implementation, no external modules
- Offers only prompts, no tools or resources
- Prompts are stored in text files
- Prompt files can use Go's templating functionality

## Project milestones

Milestones are very waterflowy, not agile at all thing. But that's what you are getting here. As a single (in more than one way), mono-cerebral, non-ambidexterous developer I do not scale. So...milestones. Cool? Cool.

**Things to do for 0.1.0 version:**

1. Create simple stdio server: listens stdin, outputs to stdout ✅
2. Implement JSON RPC MCP server-client intialization ✅
3. Implement JSON RPC MCP server-client shutdown ✅
4. Implement basic prompts capability (can be a static mock at this stage) ✅
5. Implement prompt storage and prompt loading from files, i.e. load files from a given path ✅
6. Clean up and refactor ✅

**Things to do for 0.2.0 version:**

1. Unit tests for the existing code
2. Architectural direction and choises
3. Introduce Go templating to the prompt files

**Things to do for 0.3.0 version:**

1. Github actions to produce downloadable binaries
2. Github releases
3. Support for brew install if possible

**Things to do for 0.4.0 version:**

1. Support for additional install method

**Things to do for 1.0.0 version:**

1. "Things" from 0.x done


## Project motivations

- Learn to implement a basic MCP server from scratch
- Learn about JSON RPC
- Opportunity to use Go's templating system creatively


## Installation

Currently installing prompter requires Go 1.24 or later to be installed on your system. Prompter supports all OS and architetures wich you can compile it to. Development of prompter is done in MacOS so some "unixy" emphaises might exist.

To install clone the repo and build the project. The project roots main.go should be the target.


If you see a tagged version you can start using the following command to install prompter.
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

For example in LM Studio, configuration for stdio based prompter MCP-server would be as follows:

```json
{
  "mcpServers": {
    "prompter": {
      "command": "/usr/local/bin/prompter"
    }
  }
}
```

Note that the path to the command depends on where prompter is installed.
