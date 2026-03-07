# Prompt List Change Notification

## Purpose
Define how the server notifies clients when the prompt list changes due to successful prompt creation.

## Requirements
### Requirement: Notify clients when prompt list changes after successful save
The server SHALL send an MCP prompts list changed notification after `tools/saveNewPrompt` successfully persists a new prompt.

#### Scenario: Notification emitted after successful save
- **WHEN** a client calls `tools/saveNewPrompt` and the prompt is stored successfully
- **THEN** the server sends a prompts list changed notification to connected clients

### Requirement: Do not notify on failed prompt save
The server MUST NOT send a prompts list changed notification if `tools/saveNewPrompt` fails to persist the prompt.

#### Scenario: Save fails before persistence
- **WHEN** a client calls `tools/saveNewPrompt` and storage returns an error
- **THEN** the server does not send a prompts list changed notification

### Requirement: Preserve save operation compatibility when notification delivery fails
If prompt persistence succeeds but notification delivery fails, the server SHALL preserve the successful save response behavior and treat notification failure as a non-fatal side effect.

#### Scenario: Prompt saved but notification transport fails
- **WHEN** `tools/saveNewPrompt` stores the prompt successfully and notification emission fails
- **THEN** the save operation still reports success according to existing API behavior

### Requirement: Emit list changed notification only when capability is declared
The server SHALL emit prompts list changed notifications only when it advertises prompt list change support in MCP capabilities.

#### Scenario: Capability declared
- **WHEN** the server advertises prompt list change capability and a prompt is saved successfully
- **THEN** the server emits a prompts list changed notification

#### Scenario: Capability not declared
- **WHEN** the server does not advertise prompt list change capability and a prompt is saved successfully
- **THEN** the server does not emit a prompts list changed notification
