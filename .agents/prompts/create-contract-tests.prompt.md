---
name: "Create Contract Tests"
description: "Guide through writing a contract test for a Linden layer interface. Use when: starting a new interface, adding a method to an existing interface, or when a contract test is missing."
---
# Create Contract Tests

## Inputs

- Layer name (api, orchestrator, inference, storage, discovery)
- Interface method(s) to cover
- Expected behavior per method

## Procedure

1. Read `docs/architecture/layer-interface-spec.md` for the target layer's interface shape.
2. Create or open `validation/contracts/<layer>_contract_test.go`.
3. Add build tag `//go:build contracts` at the top.
4. For each method, write a table-driven test covering:
   - Happy path with valid inputs
   - Error path (dependency returns error)
   - Empty/zero-value input
   - Timeout (if the method accepts a context)
5. Use interface doubles (never concrete implementations) for dependencies.
6. Run `go test -tags=contracts ./validation/contracts/...` and confirm tests compile and run.

## Output

- File path of created/modified test
- Methods covered
- Test cases per method
- Result of test run
