# Stage A: Contract-First

## Status
Ready to implement.

## Purpose

Build the chat path one interface at a time, with contract tests written before
the code that satisfies them. Stage 0's exemption from contract-first ends here.

Branching rules: `docs/project/branching-strategy.md`.
Delivery policy: `docs/project/contract-first-delivery.md`.

## Prerequisite

`validation/` sits outside the `src` module, so `go test ../validation/...` from
`src` fails with `directory prefix does not contain main module`. A.1 creates
`validation/go.mod` as a separate module with a `replace` directive pointing at
`../src`.

Both the Makefile and the CI contract job already guard on `validation/go.mod`
existing and skip until it does. A.1 is what turns those gates on.

## Shared Decisions

| Decision | Value |
|----------|-------|
| Validation module path | `github.com/dawsonyoung/linden/validation` |
| Source dependency | `replace github.com/dawsonyoung/linden => ../src` |
| Contract build tag | `//go:build contracts` |
| Integration build tag | `//go:build integration` |
| Test helper location | Same package as the tests that use them |

## What a Contract Test Is Here

A contract test asserts what a consumer may depend on. A unit test asserts how a
layer behaves internally. They overlap for trivial handlers; the distinction is
the audience.

- Contract: the published surface in `docs/product/spec/40-api-surface.md` and
  the interface shapes in `docs/architecture/layer-interface-spec.md`.
- Unit: branches, error wrapping, edge cases, co-located with the source.

If a change would break a client, a contract test should fail. If a change would
only break an internal assumption, a unit test should fail.

## Sequencing

```
A.1 test/contract-harness
     ↓
A.2 test/inference-contracts ──→ A.3 feat/inference-ollama-adapter
                                      ↓
A.4 test/orchestrator-contracts ──→ A.5 feat/orchestrator-chat-service
                                      ↓
A.6 test/api-contracts-sse ──────→ A.7 feat/api-chat-handlers
```

Each contract branch merges before its paired implementation branch. The
implementation PR is not ready for review until the contract tests it satisfies
are already on `main`.

---

## A.1 `test/contract-harness`

Intent: create the validation module and prove the contract gate works, using the
endpoints that already exist.

Agents: Contract Test Agent, then Validation Agent, then Review Agent.

### Tasks

1. `validation/go.mod` — module `github.com/dawsonyoung/linden/validation`, go
   1.22, `replace github.com/dawsonyoung/linden => ../src`.
2. `validation/contracts/harness_test.go` — shared assertion helpers, build tag
   `contracts`. Helpers live in the package that uses them, not a separate
   package, until a second consumer exists.
3. `validation/contracts/api_contract_test.go` — contract tests for the
   endpoints published in `docs/product/spec/40-api-surface.md`.

### Coverage Required

Assert the published contract, not the implementation:

- `GET /health` returns 200, `application/json`, and a status field
- `GET /version` returns 200, `application/json`, and version, commit, and Go
  version fields
- Every response carries `X-Request-ID`
- A client-supplied `X-Request-ID` is preserved when safe
- An unsafe `X-Request-ID` is replaced rather than echoed
- An unknown path returns 404
- A wrong method on a known path returns 405

### Acceptance Criteria

```sh
cd validation && go test -tags=contracts -race ./contracts/...   # exits 0
cd validation && go vet -tags=contracts ./...                    # exits 0
make validate-contracts                                          # runs, not skips
```

The CI `contract` job must stop reporting a skip notice and run the suite.

### Gate

Contract. Lint and unit must remain green.

Estimated size: 200–300 lines.

---

## A.2 `test/inference-contracts`

Intent: specify the inference interface before any provider code exists.

Agents: Design Agent if the interface shape needs refinement, then Contract Test
Agent, then Review Agent.

### Tasks

1. Define the `Client` interface in `src/inference` with no implementation:
   `ListModels(ctx) ([]Model, error)` and
   `ChatStream(ctx, Request, func(Chunk) error) (Result, error)`.
2. Define the shared types the interface exposes: `Model`, `Request`, `Chunk`,
   `Result`, and the error taxonomy mapping.
3. `validation/contracts/inference_contract_test.go` — contract tests driven by a
   test double implementing the interface.

### Coverage Required

- `ListModels` returns the provider's models mapped to the product shape
- `ListModels` maps a transport failure to `unavailable`
- `ChatStream` delivers chunks in order
- `ChatStream` reports a terminal result exactly once
- `ChatStream` stops when the callback returns an error
- `ChatStream` honors context cancellation
- A provider timeout maps to `deadline_exceeded`
- An unknown model maps to `not_found`

### Gate

Contract.

Estimated size: 250–400 lines.

---

## A.3 `feat/inference-ollama-adapter`

Intent: implement the Ollama provider against the A.2 contract.

Agents: Implementation Agent, then Unit Test Agent, then Validation Agent, then
Review Agent.

### Tasks

1. `src/inference/ollama.go` — HTTP client for Ollama's model list and chat
   streaming endpoints, with explicit timeouts.
2. Error mapping from transport and HTTP status to the error taxonomy.
3. Unit tests using `httptest`, including malformed stream payloads.

### Obligation

Carried from `docs/adr/0002-role-of-docker-in-distribution.md`: document the
supported container topologies for reaching Ollama — on the host versus as a
sibling service — in this PR.

### Gate

Contract plus unit. No live Ollama in CI; the adapter is tested against
`httptest`.

Estimated size: 400–600 lines.

---

## A.4 through A.7

Specified when A.3 completes, following the same shape: contract branch, then
implementation branch, each with the coverage list written before the code.

A.6 carries the SSE contract, which is the first place the `WriteTimeout` TODO in
`src/api/server.go` must be resolved — a 30 second write deadline will terminate
a streaming response.

## Stage A Exit Criteria

1. A chat request completes end to end against a local model.
2. Replies stream to the client in order, with a single terminal event.
3. Contract tests exist for every interface method on the chat path.
4. Unit tests cover error, timeout, and cancellation paths in each layer.
5. The `contract` CI job runs a real suite rather than reporting a skip.
6. `docs/product/spec/` describes chat as available, updated in the PR that
   ships it.
