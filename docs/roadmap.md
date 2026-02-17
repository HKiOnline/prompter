# Prompter Roadmap

This document outlines Prompter MCP-server's roadmap.

The roadmap has three categories. The **next up** section highlights features which are the most likely to be done next or are already in being implemented. The **completed** are features which are already in the main branch of the project and released. The **potentials** are features that might be implemented. This is more a *list of ideas* than a list of promises. Their order do not indicate priority or the order which they would be implemented. 

Priority of the items are defined by ephemeral, transient ***vibes***. 

## Next Up

- [ ] Support argument based templating i.e. sending arguments via prompt/get
- [ ] Support prompt list update (i.e. the clients get notified the list has changes)

## Completed

- [x] File based storage manager for prompts (fsProvider)
- [x] Support stdio based JSON-RPC -transport
- [x] List prompts
- [x] Get prompt with name
- [x] List tools
- [x] Create new prompt tool (uses storage provider to save the prompt)
- [x] Tested using [OpenCode](https://opencode.ai/)
- [x] Support first built-in template function (date)
- [x] Support Streamable HTTP based JSON-RPC transport

## Potential

- [ ] Architectural documentation
- [ ] Automated Github actions to produce downloadable binaries
- [ ] Automated Github releases
- [ ] Support for homebrew install for MacOS
- [ ] Support for deb install package for Debian-based Linux distributions
- [ ] Support for rpm install package for Redhat-based Linux distributions
- [ ] Storage provider for sqlite database
- [ ] Storage provider for git
- [ ] CLI command to list prompts
- [ ] CLI command to get a prompt
- [ ] CLI command to get a sample prompt
- [ ] Support multiple directories in fsProvider

