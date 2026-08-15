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

See `docs/process/branching-strategy.md`.

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

## Linux and Docker Parity Policy

1. Linux CI is the source of truth and must be green on every PR.
2. Docker build and run checks are introduced early, not postponed.
3. Container smoke checks include health endpoint and SSE stream behavior.
4. Windows scripts remain supported, but cannot bypass Linux gating.

## Proposed PR Sequence

### PR 1: Foundation contract skeleton

Includes:

1. Contract test harness structure in validation/contracts.
2. Shared test helpers for contract assertions.
3. Initial API contract test stubs for health/version/chat schema.

Must pass:

1. Contract tests compile and run.
2. No implementation yet beyond compile-safe placeholders.

### PR 2: Inference interface contracts

Includes:

1. Inference interface definitions.
2. Contract tests for ListModels and ChatStream semantics.

Must pass:

1. Contract tests fail first, then pass with minimal adapter mock.
2. Unit tests for inference request/response mapping helpers.

### PR 3: Orchestrator interface contracts

Includes:

1. ChatService and orchestration contracts.
2. Dependency interaction contracts using inference/store doubles.

Must pass:

1. Unit tests for orchestration decision paths.
2. Contract tests for mapping and error behavior.

### PR 4: API interface contracts and SSE behavior

Includes:

1. HTTP request/response schema contracts.
2. SSE framing, event ordering, and terminal behavior contracts.

Must pass:

1. Unit tests for request validation and response encoding.
2. Contract tests asserting message, done, and error event semantics.

### PR 5: Minimal concrete vertical slice

Includes:

1. cmd wiring.
2. api to orchestrator to inference flow.
3. Health and version endpoints.

Must pass:

1. Existing contract tests.
2. New integration test for end-to-end chat happy path.

### PR 6: Storage module contracts and implementation

Includes:

1. Session persistence contract tests.
2. Minimal storage implementation.
3. Orchestrator integration with optional session behavior.

Must pass:

1. Unit tests for storage edge cases.
2. Integration tests for session save/load behavior.

### PR 7: Discovery contracts and implementation

Includes:

1. Discovery start/stop/status contracts.
2. Minimal LAN advertisement behavior.

Must pass:

1. Unit tests for lifecycle and idempotency.
2. Integration tests with startup wiring.

### PR 8: Functional baseline and smoke

Includes:

1. validation/smoke startup and chat stream checks.
2. validation/security baseline checks.
3. Linux Docker smoke execution path.

Must pass:

1. CI gates for lint, unit, contracts, integration, smoke.
2. Docker build plus container smoke checks.

## Definition of Done Per PR

1. Interface and behavior documented or explicitly unchanged.
2. Contract tests included first for new behavior.
3. Unit tests included with implementation.
4. Linux CI green.
5. If runtime behavior changes, Docker smoke still green.
6. No unresolved TODO that affects security boundaries.

## Review Checklist

1. Does this PR implement only one logical concern?
2. Are contract tests present and meaningful?
3. Are unit tests table-driven with edge cases?
4. Are error messages safe and non-sensitive?
5. Is layer dependency direction preserved?
6. Is SSE behavior deterministic and test-covered where relevant?

## Labels

Branch naming and PR rules: `docs/process/branching-strategy.md`.

Suggested labels:

1. layer:api
2. layer:orchestrator
3. layer:inference
4. layer:storage
5. layer:discovery
6. test:contracts
7. test:integration
8. test:smoke
9. platform:linux
10. platform:docker

## Immediate Execution Tasks

1. Approve and freeze the interface specification document.
2. Create contract test skeleton directories and helper package.
3. Open PR 1 for foundation contract skeleton.
4. Add CI job separation for contract, unit, integration, smoke, and docker checks.
