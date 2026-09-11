# Linden Contract-First Delivery Plan

## Goal
Implement Linden in reviewable chunks with a clean PR history by writing interface tests first, then implementing to satisfy those tests, and only later enabling broader integration and functional suites.

## Development Workflow Contract

Every implementation unit follows this sequence:

1. Interface proposal update in architecture spec.
2. Contract tests added first in validation/contracts.
3. Implementation added in source layer.
4. Unit tests added in the layer package.
5. Integration tests expanded when two or more concrete modules are wired.
6. Functional and smoke tests expanded after full vertical slice is runnable.

No implementation PR is merged unless steps 2 and 4 are present.

## PR Size Rules

See `docs/project/branching-strategy.md`.

## Test Progression Model

### Stage A: Contract-first (initial)

Primary focus:

1. validation/contracts
2. Layer unit tests

Do not require:

1. Full end-to-end flow
2. Functional UX assertions

### Stage B: Modular integration (middle)

Trigger:

- api, orchestrator, and inference have concrete MVP implementations.

Primary focus:

1. validation/integration for vertical chat path
2. Cross-layer error and timeout behavior
3. Unit tests continue for each incremental feature

### Stage C: Functional and smoke (later)

Trigger:

- Chat path, model list, and storage session path are wired.

Primary focus:

1. validation/smoke for startup and chat stream
2. validation/security for boundary checks and leak prevention
3. Browser/API functional tests for core user journeys

#### LAN Testing
When Stage C is complete, testing from a mobile device or separate computer on the local network is supported and automated.
- **Automated Check:** `make docker-smoke` will automatically dynamically determine the host's LAN IP, boot the container on `0.0.0.0`, and run the Playwright integration suite against the LAN interface to prove it is accessible.
- **Manual Verification:** Developers can run `docker compose up --build` and connect from their phone or another device by navigating to `http://<LAN_IP>:8080`.

## Linux and Docker Parity Policy

1. Linux CI is the source of truth and must be green on every PR.
2. Docker build and run checks are introduced early, not postponed.
3. Container smoke checks include health endpoint and SSE stream behavior.
4. Windows scripts remain supported, but cannot bypass Linux gating.

## Standard Capability PR Sequence

For any new feature or capability (e.g. Track C document retrieval), the contract-first lifecycle proceeds through these distinct stages:

### Phase 1: Interface & Contract Specification
1. Update or create the normative technical specification in `docs/architecture/interfaces/`.
2. Define exact Go interface methods and types.
3. Author isolated contract tests under `validation/contracts/` using `//go:build contracts`.
4. Run `cd tools && go run ./docgen -verify` to confirm documentation and AST parity.
5. Contract tests should fail or compile against minimal doubles before implementation is written.

### Phase 2: Layer Implementation & Unit Testing
1. Implement the interface within the target layer package (`src/orchestrator/`, `src/inference/`, etc.).
2. Write comprehensive, table-driven unit tests in the same package (`Test_FunctionName_Scenario_Expected`).
3. Run `cd validation && go test -tags=contracts ./contracts/...` to prove the implementation satisfies the contract.
4. Verify edge cases: nil contexts, network timeouts, empty inputs, disk exhaustion.

### Phase 3: Cross-Layer Integration & Wiring
1. Wire the new layer implementation into `src/cmd/main.go`.
2. Expand cross-layer tests in `validation/integration/` using real (non-mocked) internal components.
3. Assert that errors bubble up correctly through the shared `errs` taxonomy.

### Phase 4: Runtime & Security Validation
1. Execute the full local validation suite:
   ```sh
   make lint && make test && make validate
   ```
2. Verify container build and smoke tests:
   ```sh
   make docker-smoke
   ```
3. Audit for privacy: zero user chat data in logs, sanitized API error responses.

---

## Definition of Done Per PR

1. Interface and behavior documented in `docs/architecture/interfaces/`.
2. Contract tests included and passing in `validation/contracts/`.
3. Table-driven unit tests included with the implementation.
4. Linux CI gates 100% green.
5. If runtime behavior or endpoints change, container smoke tests pass.
6. No user content in error messages or server logs.

## Review Checklist

1. Does this PR implement only one logical concern?
2. Are contract tests present and meaningful?
3. Are unit tests table-driven with edge cases?
4. Are error messages safe and non-sensitive?
5. Is layer dependency direction preserved?
6. Is SSE behavior deterministic and test-covered where relevant?

## Labels

Branch naming and PR rules: `docs/project/branching-strategy.md`.

Suggested labels:

1. `layer:api`
2. `layer:orchestrator`
3. `layer:inference`
4. `layer:storage`
5. `layer:discovery`
6. `layer:mcp`
7. `test:contracts`
8. `test:integration`
9. `test:smoke`
10. `platform:linux`
11. `platform:docker`
