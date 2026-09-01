# Linden Layer Interface Specification

## Status
Draft for implementation guidance.

## Purpose
Define the implementation target for each runtime layer, the interface boundaries between layers, and the contracts that must be validated before implementation code is accepted.

## Scope
This specification covers MVP and near-MVP runtime behavior for:

1. cmd
2. api
3. orchestrator
4. inference
5. storage
6. discovery
7. web

It does not define marketplace, mobile, hardware image, or distributed deployment concerns.

## Architectural Rules

1. Dependencies are one-way and follow layer order.
2. Consumers depend on interfaces, not concrete implementations.
3. Each interface must have contract tests under validation/contracts before implementation is merged.
4. Unit tests must be added with each implementation PR.
5. Integration and functional tests are introduced only after enough interfaces are implemented to execute real flows.
6. No user content in logs or user-facing error bodies.
7. Linux parity is strict: all CI gates must pass on Linux and Docker runtime behavior must be testable early.

## Dependency Graph

cmd -> api -> orchestrator -> inference
cmd -> api -> orchestrator -> storage
cmd -> discovery
api serves static output from web build artifacts

Any layer may import `src/errs`. It is a shared kernel, not a layer, and it
imports nothing itself.

Forbidden:

1. api importing inference or storage directly.
2. inference importing api or orchestrator.
3. storage importing api or orchestrator.
4. discovery importing api or orchestrator.
5. Any layer importing cmd.

## Shared Domain Contracts

### Chat message contract

- role: system | user | assistant
- content: string
- timestamp: RFC3339 UTC string

### Chat request contract

- model: string
- messages: array of chat messages, minimum length 1
- stream: boolean
- session_id: optional string

### Chat response contract (non-stream)

- id: string
- model: string
- output_text: string
- finish_reason: string
- usage: optional token counts

### SSE stream contract

- content type: text/event-stream
- events:
  - message: incremental text chunk
  - metadata: model or usage updates
  - done: terminal completion status
  - error: safe client-facing error event

SSE behavior requirements:

1. Event order is preserved.
2. done or error is terminal.
3. Stream closes after terminal event.
4. Keep-alive comments may be sent for long generations.

### Error taxonomy

Implemented by `src/errs`. The package and this list are one contract: adding a
code requires changing both in the same PR. See
`docs/adr/0004-shared-error-taxonomy.md`.

- invalid_argument
- unauthenticated
- permission_denied
- not_found
- conflict
- resource_exhausted
- unavailable
- deadline_exceeded
- internal

Rules:

1. API maps internal errors to stable error codes.
2. Internal details remain server-side only.
3. Correlation/request ID included in server logs and optional response header.
4. An error carrying no code is treated as `internal`.
5. Any layer may import `src/errs`. No layer imports another layer for error
   definitions.

## Layer Interfaces

### cmd layer
Responsibility:

1. Read config.
2. Wire concrete implementations.
3. Start and stop processes.

Public boundary:

- main entrypoint only.

Rules:

1. No business logic.
2. No direct HTTP handler logic.

### api layer
Responsibility:

1. Expose HTTP routes.
2. Validate request shape and limits.
3. Convert HTTP to orchestrator calls.
4. Return JSON and SSE responses.

Required interface dependency:

- ChatService from orchestrator.

Expected interface shape:

- Health endpoint handler.
- Version endpoint handler.
- Chat handler supporting JSON and SSE.
- Model list handler.

### orchestrator layer
Responsibility:

1. Validate workflow-level semantics.
2. Assemble context and policy decisions.
3. Call inference and storage.
4. Transform provider output into product output schema.

Required interface dependencies:

- InferenceClient
- ConversationStore (initially minimal)

Expected interface shape:

- Chat: accepts chat request and streaming writer/callback.
- ListModels: pass-through plus policy filtering.

### inference layer
Responsibility:

1. Provide model listing.
2. Provide chat generation streaming and non-streaming.
3. Map provider transport errors to inference-level errors.

Expected interface shape:

- ListModels(ctx) -> []Model, error
- ChatStream(ctx, request, onChunk func(Chunk) error) -> Result, error

Rules:

1. Provider-specific details are hidden behind interface.
2. Timeouts and retry policy are explicit and testable.
3. The callback returns an error so a consumer can abort a stream. When it
   returns non-nil, ChatStream stops delivering, does no further provider work,
   and returns that error unwrapped.
4. Context cancellation is user-initiated stop; a callback error is
   consumer-side failure. Both are supported and they are distinct.
5. Errors use `src/errs` codes: provider unreachable maps to `unavailable`,
   provider timeout to `deadline_exceeded`, unknown model to `not_found`, and a
   request carrying no messages to `invalid_argument` before any provider call.
6. The callback is invoked on the calling goroutine, in order, and never after
   ChatStream returns. Consumers need no synchronization.
7. An implementation checks `ctx.Err()` before each delivery and before
   classifying any transport failure. Cancellation takes precedence over
   transport classification.
8. Result is meaningful only when the error is nil or is a context error. On
   every other error path it is the zero value.

Every implementation must pass the conformance suite in
`validation/contracts/inference_contract_test.go`.

#### Supported Ollama Container Topologies

Per `docs/adr/0002-role-of-docker-in-distribution.md`, three network topologies are supported for communicating with Ollama:

1. **Host-native Ollama (GPU-accelerated)**:
   - Ollama runs natively on the host OS for direct GPU hardware acceleration.
   - The Linden container connects using `OLLAMA_URL=http://host.docker.internal:11434` (Docker Desktop / host gateway).
2. **Compose Sibling Container**:
   - Ollama runs as a container service on the same Docker network as Linden.
   - Linden connects using `OLLAMA_URL=http://ollama:11434`.
3. **Local Process (Bare Metal / Development)**:
   - Linden and Ollama run directly as local processes on the host.
   - Linden connects using `OLLAMA_URL=http://127.0.0.1:11434`.

### storage layer
Responsibility:

1. Persist sessions and settings.
2. Persist retrieval metadata in later phase.

Expected interface shape for MVP:

- SaveSession(ctx, session) -> error
- LoadSession(ctx, sessionID) -> session, error
- ListSessions(ctx, filter) -> []session, error

Rules:

1. Storage errors map to stable error taxonomy.
2. Sensitive fields have explicit serialization rules.

### discovery layer
Responsibility:

1. Advertise service on LAN.
2. Support pairing and revocation hooks in later phase.

Expected interface shape:

- Start(ctx) -> error
- Stop(ctx) -> error
- Status() -> discovery state

### web layer
Responsibility:

1. Render chat and settings views.
2. Consume API JSON and SSE endpoints.
3. Surface privacy controls and status to user.

Rules:

1. Strict TypeScript.
2. No any.
3. SSE reconnection strategy is explicit.

## Cross-Cutting Non-Functional Requirements

1. Linux parity is mandatory in CI.
2. Docker image path must work early and remain green.
3. Health endpoint is required for runtime and container health checks.
4. Request timeouts and payload limits are required at API boundary.
5. Observability baseline includes structured logs and request IDs.

## Contract Test Coverage Matrix

Minimum contract tests before implementation merge:

1. api contracts
   - request validation and error mapping
   - SSE framing and terminal event behavior
2. orchestrator contracts
   - dependency invocation order
   - fallback and error mapping behavior
3. inference contracts
   - model list mapping
   - stream chunk, done, and error mapping
4. storage contracts
   - persistence and retrieval behavior
   - not found and conflict semantics
5. discovery contracts
   - start/stop idempotency and status behavior

## Exit Criteria For Interface Spec Stability

This document is considered stable enough for implementation when:

1. API request and response schemas are frozen for MVP.
2. SSE event contract is accepted.
3. Error taxonomy is accepted.
4. Each layer has at least one contract test file stub committed.
