---
name: "Implement From Contract"
description: "Guide through implementing a source layer to satisfy existing contract tests. Use when contract tests exist and it is time to write the implementation."
---
# Implement From Contract

## Inputs

- Layer to implement (api, orchestrator, inference, storage, discovery)
- Contract test file(s) to satisfy

## Procedure

1. Read the contract test file(s) in full before writing any code.
2. Read the interface spec section for this layer in `docs/architecture/layer-interface-spec.md`.
3. Check `.context/lessons_learned.md` for relevant rules.
4. Implement the minimal code in `src/<layer>/` that makes the contract tests pass.
5. Verify layer dependency direction — no forbidden imports.
6. Run `go test -tags=contracts ./validation/contracts/...` — all tests must pass.
7. Delegate unit test writing to the Unit Test Agent before closing.

## Output

- Files created or modified
- Contract tests now passing
- Dependency direction confirmed
- Unit test delegation handed off
