---
name: contract-first-implementation
description: "Step-by-step contract-first development workflow for Linden. Use when implementing any layer — guides the full sequence from interface spec to passing contract and unit tests."
---
# Contract-First Implementation

## When to Use

- Starting implementation of any layer interface method.
- Adding a new method to an existing interface.
- Onboarding to the codebase before writing first code.

## Procedure

### Step 1: Read the interface spec

Open `docs/architecture/01-layer-interface-spec.md`.
Locate the target layer section. Note the method signatures and behavior rules.

### Step 2: Write the contract test

Delegate to Contract Test Agent or follow `prompts/create-contract-tests.prompt.md`.
Confirm tests compile with `go test -tags=contracts ./validation/contracts/...`.
Commit the test file before writing implementation.

### Step 3: Implement to the contract

Delegate to Implementation Agent or follow `prompts/implement-from-contract.prompt.md`.
Run contract tests after each method — they must pass before moving on.

### Step 4: Write unit tests

Delegate to Unit Test Agent.
Cover all behavior paths including error, empty input, and timeout.
Run `go test -race ./src/...`.

### Step 5: Confirm gates

Run `prompts/run-validation-gates.prompt.md` for Contract and Unit gates.
Both must be PASS before the PR is packaged.

## Common Pitfalls

- Writing implementation before the contract test exists — always test first.
- Using concrete types instead of interfaces in contract tests.
- Skipping the empty-input and timeout test cases.
