---
name: validation-gates
description: "Execute and interpret Linden quality gates. Use when running CI checks locally, diagnosing a gate failure, or producing gate evidence for a PR."
---
# Validation Gates

## When to Use

- Before packaging a PR.
- After a Debug Agent fix — confirm the gate now passes.
- When CI fails and the failure needs local reproduction.

## Gate Reference

See `prompts/run-validation-gates.prompt.md` for the full gate command table.

## Interpreting Failures

### Lint failure (`go vet`)

Check for: unused variables, unreachable code, incorrect format verbs, suspicious constructs.
Fix in the implementation file, then re-run.

### Unit failure

Read the test name — it encodes `FunctionName_Scenario_Expected`.
Check: did the implementation change break an existing behavior? Or is a new case missing?

### Contract failure

Check: does the implementation satisfy the interface shape the test asserts?
Check: are all required error cases handled and returning the correct error type?

### Docker build failure

Check: multi-stage build paths, missing COPY sources, base image availability.

### Health check failure

Check: server startup time, port binding, `/health` route registration.

## Producing Gate Evidence

After a full pass, output:

```
| Gate        | Status | Notes |
|-------------|--------|-------|
| Lint        | PASS   |       |
| Unit        | PASS   | N tests |
| Contract    | PASS   | N tests |
| Integration | PASS   | N tests |
| Docker      | PASS   |       |
| Health      | PASS   | HTTP 200 |
```
