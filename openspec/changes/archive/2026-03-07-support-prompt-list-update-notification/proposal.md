## Why

MCP servers that declare prompt list change support should notify clients when the prompt list changes, so clients can refresh state without polling or restarting. Prompter currently saves prompts but does not emit this notification, leaving clients unaware of new prompts until they manually re-request the list.

## What Changes

- Emit the MCP prompt list changed notification after `tools/saveNewPrompt` successfully stores a new prompt.
- Investigate and use the notification mechanism provided by `modelcontextprotocol/go-sdk` for this server-side event.
- Keep behavior unchanged when save fails (no notification emitted on failed writes).
- Keep the existing prompt save API behavior and response contract intact.

## Capabilities

### New Capabilities
- `prompt-list-change-notification`: Server emits a prompt list changed notification to connected clients after successfully adding a prompt.

### Modified Capabilities
- (none)

## Impact

- Affected code: MCP tool handling for saving prompts and server notification plumbing in the RPC layer.
- API/protocol impact: Adds standards-aligned MCP prompt list change notifications; no breaking API changes.
- Dependencies: Relies on existing `modelcontextprotocol/go-sdk` notification support (or equivalent provided mechanism).
- Compatibility: Backward compatible for existing clients and prompt management workflows.
