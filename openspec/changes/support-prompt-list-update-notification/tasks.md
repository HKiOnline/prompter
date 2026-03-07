## 1. SDK and Capability Wiring

- [ ] 1.1 Verify how `modelcontextprotocol/go-sdk` exposes prompts list changed notifications and identify the correct API to call.
- [ ] 1.2 Confirm where prompt list change capability is declared in server capabilities and ensure this change does not alter existing initialization behavior.

## 2. Save Flow Integration

- [ ] 2.1 Update the `tools/saveNewPrompt` success path to emit prompts list changed notification only after prompt persistence succeeds.
- [ ] 2.2 Guard notification emission so it runs only when prompt list change capability is advertised.
- [ ] 2.3 Handle notification emission errors as non-fatal side effects (log/record error without changing successful save response semantics).

## 3. Verification and Regression Coverage

- [ ] 3.1 Add or update tests covering notification emission on successful save and no emission on failed save.
- [ ] 3.2 Add or update tests for capability-gated behavior (emit when declared, do not emit when not declared).
- [ ] 3.3 Add or update tests for notification-delivery failure to verify save remains successful and backward compatible.
- [ ] 3.4 Run project test suite (`make test`) and resolve any regressions.
