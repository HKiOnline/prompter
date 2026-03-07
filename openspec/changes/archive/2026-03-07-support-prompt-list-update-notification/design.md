## Context

Prompter already supports prompt creation via `tools/saveNewPrompt`, but clients do not receive a protocol-level signal that the prompt list changed after a successful save. MCP defines a prompt list changed notification for servers that advertise list-change support, and this change adds that behavior while keeping current save semantics intact.

The implementation should remain localized to the prompt-save path and existing RPC/server plumbing. The design must preserve backward compatibility for clients that do not consume notifications.

## Goals / Non-Goals

**Goals:**
- Emit a prompt list changed notification when `tools/saveNewPrompt` successfully persists a prompt.
- Use the notification API/mechanism provided by `modelcontextprotocol/go-sdk` where available.
- Keep current request/response behavior unchanged for prompt saving.
- Ensure notification emission failures do not corrupt prompt persistence results.

**Non-Goals:**
- No changes to prompt file format, storage provider contracts, or prompt retrieval APIs.
- No attempt to emit notifications for failed save operations.
- No broad event system redesign beyond the specific MCP prompts list change notification.

## Decisions

1. Trigger notification only after successful persistence
- Decision: Emit the list-changed notification only after storage returns success for `saveNewPrompt`.
- Rationale: Prevents false-positive refresh events and aligns protocol semantics with actual state change.
- Alternative considered: Emit before save attempt for lower latency. Rejected because it can notify on writes that later fail.

2. Keep notification as a best-effort side effect
- Decision: Treat notification emission as non-critical relative to persistence success; return save success even if notifying fails, while logging/recording the notify error.
- Rationale: Persistence is the primary operation; a transient transport/session issue should not turn a successful write into a user-visible failure.
- Alternative considered: Fail the tool call when notification fails. Rejected because it breaks backward compatibility and can create duplicate writes on client retries.

3. Route through SDK-native MCP notification support
- Decision: Use the go-sdk's server/session notification facility for prompts list changed instead of custom JSON-RPC writes.
- Rationale: Keeps implementation aligned with MCP abstractions and avoids protocol drift.
- Alternative considered: Manual JSON-RPC notification construction. Rejected due to maintainability and higher risk of schema mismatches.

4. Keep capability advertisement and emission behavior coherent
- Decision: Ensure notifications are sent only in server modes/configurations where prompt list change support is declared.
- Rationale: Matches MCP expectations and avoids sending undeclared protocol behavior.
- Alternative considered: Always emit notifications regardless of advertised capability. Rejected due to protocol inconsistency.

## Risks / Trade-offs

- [SDK API uncertainty] Notification helper naming/signature may differ from assumptions -> Mitigation: confirm exact go-sdk API during implementation and encapsulate call in a small helper.
- [Silent notify failures] Clients may not refresh promptly if notification delivery fails -> Mitigation: add structured logging and integration coverage for successful notification path.
- [Behavior drift across transports] Notification behavior may differ by client/transport implementation -> Mitigation: validate in existing integration tests and test-server flow.
- [Concurrency timing] Rapid consecutive saves may send multiple notifications -> Mitigation: accept as protocol-valid list change signals; avoid dedup unless required by observed client issues.
